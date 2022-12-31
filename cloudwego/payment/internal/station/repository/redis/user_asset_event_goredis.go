package redis

import (
	"context"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
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

func (m *UserAssetEventRepository) UserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) (userAsset *base.UserAsset, err error) {
	var assetCn int64
	switch strings.ToLower(m.changeAssetTxMethod) {
	case constants.UserAssetTxMethodCas:
		assetCn, err = m.watchUserAssetChangeTx(ctx, eventId, changeInfo, handle)
	case constants.UserAssetTxMethodLua:
		assetCn, err = m.userAssetChangeLuaAtomicTx(ctx, eventId, changeInfo, handle)
	default:
		assetCn, err = m.userAssetChangeLuaAtomicTx(ctx, eventId, changeInfo, handle)
	}
	if err != nil {
		klog.CtxErrorf(ctx, "%s UserAssetChangeTx err:%s", m.changeAssetTxMethod, err.Error())
		return
	}

	userAsset = &base.UserAsset{
		UserId:    changeInfo.UserId,
		AssetType: changeInfo.AssetType,
		AssetCn:   assetCn,
	}

	return
}

func (m *UserAssetEventRepository) watchUserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) (int64, error) {
	key := constants.GetUserAssetInfoKey(changeInfo.UserId, changeInfo.AssetType.String())
	eventMsgKey := constants.GetUserAssetEventMsgKey(changeInfo.UserId, eventId)

	err := m.cb.SetAsset(ctx, changeInfo.UserId, changeInfo.AssetType)
	if err != nil {
		return 0, err
	}

	var assertCn int64
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
		if err != nil {
			return err
		}

		assertCn = assetObj.AssetCn

		return nil
	}

	// Retry if the key has been changed.
	for i := 0; i < constants.MaxRetries; i++ {
		err := m.redisClient.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return assertCn, nil
		}
		if err == redis.TxFailedErr {
			//println(key, err.Error())
			// Optimistic lock lost. Retry.
			continue
		}

		// Return any other error.
		return 0, err
	}

	return 0, domain.ErrorWatchAssetCasLoopMaxRetry
}

func (m *UserAssetEventRepository) userAssetChangeLuaAtomicTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle domain.AssetIncrHandler) (int64, error) {
	key := constants.GetUserAssetInfoKey(changeInfo.UserId, changeInfo.AssetType.String())
	eventMsgKey := constants.GetUserAssetEventMsgKey(changeInfo.UserId, eventId)

	val, err := m.redisClient.Eval(ctx, assetStringChangeLua, []string{key, eventMsgKey},
		handle(ctx), constants.RedisKeyExpireUserAssetInfo.Seconds(), constants.RedisKeyExpireAssetEventMsg.Seconds()).Result()
	if err != nil {
		return 0, err
	}
	if val.(int64) >= constants.RedisLuaAssetChangeResCodeSuccess {
		return val.(int64), nil
	}
	if val == constants.RedisLuaAssetChangeResCodeNoEnough {
		return 0, domain.ErrorNoEnoughAsset
	}

	if val == constants.RedisLuaAssetChangeResCodeNoExists {
		err = m.cb.SetAsset(ctx, changeInfo.UserId, changeInfo.AssetType)
		if err != nil {
			return 0, err
		}

		val, err = m.redisClient.Eval(ctx, assetStringChangeLua, []string{key, eventMsgKey},
			handle(ctx), constants.RedisKeyExpireUserAssetInfo.Seconds(), constants.RedisKeyExpireAssetEventMsg.Seconds()).Result()
		if err != nil {
			return 0, err
		}
		if val.(int64) >= constants.RedisLuaAssetChangeResCodeSuccess {
			return val.(int64), nil
		}
		if val == constants.RedisLuaAssetChangeResCodeNoEnough {
			return 0, domain.ErrorNoEnoughAsset
		}
	}

	return 0, nil
}

func (m *UserAssetEventRepository) GetUserAssetEventMsg(ctx context.Context, userId int64, eventId string) (res int, err error) {
	eventMsgKey := constants.GetUserAssetEventMsgKey(userId, eventId)
	res, eventErr := m.redisClient.Get(ctx, eventMsgKey).Int()

	if eventErr != nil && eventErr != redis.Nil {
		klog.CtxErrorf(ctx, "GetUserAssetEventMsg key:%s err:%s", eventMsgKey, err.Error())
		err = domain.ErrorInternalRedis
		return
	}
	klog.CtxInfof(ctx, "GetUserAssetEventMsg key:%s res:%s ok", eventMsgKey, res)

	return
}
