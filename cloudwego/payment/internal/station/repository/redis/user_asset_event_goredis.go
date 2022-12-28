package redis

import (
	"context"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
)

type UserAssetEventRepository struct {
	redisClient redis.UniversalClient
	cb          domain.IUserAssetCallBack
	// cas, lua, function-lua
	changeAssetTxMethod string
}

func NewUserAssetEventRepository(redisClient redis.UniversalClient, cb domain.IUserAssetCallBack, method string) domain.IUserAssetEventRepository {
	return &UserAssetEventRepository{redisClient: redisClient, cb: cb, changeAssetTxMethod: method}
}

func (m *UserAssetEventRepository) UserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) error {
	switch strings.ToLower(m.changeAssetTxMethod) {
	case "cas":
		m.watchUserAssetChangeTx(ctx, eventId, changeInfo, handle)
	case "lua":
		m.userAssetChangeLuaAtomicTx(ctx, eventId, changeInfo, handle)
	default:
		m.userAssetChangeLuaAtomicTx(ctx, eventId, changeInfo, handle)
	}

	return nil
}

func (m *UserAssetEventRepository) watchUserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) error {
	key := constants.GetUserAssetInfoKey(changeInfo.UserId, changeInfo.AssetType.String())
	eventMsgKey := constants.GetUserAssetEventMsgKey(changeInfo.UserId, eventId)

	err := m.cb.SetAsset(ctx, changeInfo.UserId, changeInfo.AssetType)
	if err != nil {
		return err
	}

	// Transactional function.
	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		bytes, err := tx.Get(ctx, key).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}

		assetObj := &base.UserAsset{}
		err = sonic.Unmarshal(bytes, assetObj)
		if err != nil {
			return err
		}

		// Actual operation (local in optimistic lock).
		assetObj.AssetCn += handle(ctx)
		if assetObj.AssetCn < 0 {
			return domain.ErrorNoEnoughAsset
		}
		bytes, err = sonic.Marshal(assetObj)
		if err != nil {
			return err
		}

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, bytes, constants.RedisKeyExpireUserAssetInfo)
			pipe.Set(ctx, eventMsgKey, 1, constants.RedisKeyExpireUserAssetInfo)
			return nil
		})
		return err
	}

	// Retry if the key has been changed.
	for i := 0; i < constants.MaxRetries; i++ {
		err := m.redisClient.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			//println(key, err.Error())
			// Optimistic lock lost. Retry.
			continue
		}

		// Return any other error.
		return err
	}

	return domain.ErrorWatchAssetCasLoopMaxRetry
}

func (m *UserAssetEventRepository) userAssetChangeLuaAtomicTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) error {
	key := constants.GetUserAssetInfoKey(changeInfo.UserId, changeInfo.AssetType.String())
	eventMsgKey := constants.GetUserAssetEventMsgKey(changeInfo.UserId, eventId)

	val, err := m.redisClient.Eval(ctx, assetStringChangeLua, []string{key, eventMsgKey},
		handle(ctx), constants.RedisKeyExpireUserAssetInfo.Seconds(), constants.RedisKeyExpireAssetEventMsg.Seconds()).Result()
	if err != nil {
		return err
	}
	if val == constants.RedisLuaAssetChangeResCodeSuccess {
		return nil
	}
	if val == constants.RedisLuaAssetChangeResCodeNoEnough {
		return domain.ErrorNoEnoughAsset
	}

	if val == constants.RedisLuaAssetChangeResCodeNoExists {
		err = m.cb.SetAsset(ctx, changeInfo.UserId, changeInfo.AssetType)
		if err != nil {
			return err
		}

		val, err = m.redisClient.Eval(ctx, assetStringChangeLua, []string{key, eventMsgKey},
			handle(ctx), constants.RedisKeyExpireUserAssetInfo.Seconds(), constants.RedisKeyExpireAssetEventMsg.Seconds()).Result()
		if err != nil {
			return err
		}
		if val == constants.RedisLuaAssetChangeResCodeSuccess {
			return nil
		}
		if val == constants.RedisLuaAssetChangeResCodeNoEnough {
			return domain.ErrorNoEnoughAsset
		}
	}

	return nil
}
