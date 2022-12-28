package redis

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da/paymentservice"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"golang.org/x/sync/singleflight"
)

type UserAssetCallBack struct {
	daClient       paymentservice.Client
	getAssetLocker *redsync.Redsync
	userAssetRepos domain.IUserAssetRepository
	sfGroup        *singleflight.Group
}

func NewUserAssetCallBack(daClient paymentservice.Client, locker *redsync.Redsync, userAssetRepos domain.IUserAssetRepository, sfGroup *singleflight.Group) domain.IUserAssetCallBack {
	return &UserAssetCallBack{
		daClient:       daClient,
		getAssetLocker: locker,
		userAssetRepos: userAssetRepos,
		sfGroup:        sfGroup,
	}
}

func (m *UserAssetCallBack) GetAssets(ctx context.Context, userIds []int64) (assetDtos []*base.UserAsset, err error) {
	resp, err := m.daClient.GetAssets(ctx, &da.GetAssetsReq{UserIds: userIds})
	if err != nil {
		klog.CtxErrorf(ctx, "userIds:%v GetAssets error: %s", userIds, err.Error())
		return
	}

	assetDtos = resp.UserAssets

	return
}

func (m *UserAssetCallBack) GetUserAssetByType(ctx context.Context, userId int64, assetType base.AssetType) (assetDto *base.UserAsset, err error) {
	userAssets, err := m.GetAssets(ctx, []int64{userId})
	if err != nil {
		return
	}

	for _, item := range userAssets {
		if item.AssetType == assetType {
			assetDto = item
			return
		}
	}

	return
}

func (m *UserAssetCallBack) SetAsset(ctx context.Context, userId int64, assetType base.AssetType) (err error) {
	key := constants.GetUserAssetInfoKey(userId, assetType.String())
	lockerKey := constants.GetUserAssetInfoLockKey(userId, assetType.String())

	_, err = m.userAssetRepos.GetAsset(ctx, key)
	if err == nil {
		return
	}
	if err != nil && err != redis.Nil {
		klog.CtxErrorf(ctx, "key %s userAssetRepos.GetAsset error:%s", key, err.Error())
		return
	}

	// 	_, err, _ = m.sfGroup.Do(key, func() (data interface{}, err error) {
	resCh := m.sfGroup.DoChan(key, func() (data interface{}, err error) {
		// Create a new lock client.
		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		mutex := m.getAssetLocker.NewMutex(lockerKey, redsync.WithTries(constants.DisLockerBlockRetryCn), redsync.WithExpiry(constants.RedisKeyExpireUserAssetInfoLock))

		// Obtain a lock for our given mutex. After this is successful, no one else
		// can obtain the same lock (the same mutex name) until we unlock it.
		if err = mutex.LockContext(ctx); err != nil {
			klog.CtxErrorf(ctx, "key %s lock error:%s", lockerKey, err.Error())
			return
		}
		defer func() {
			// Release the lock so other processes or threads can obtain a lock.
			if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
				klog.CtxErrorf(ctx, "key %s unlock error:%s", lockerKey, err.Error())
			}
		}()

		data, err = m.userAssetRepos.GetAsset(ctx, key)
		if err == nil {
			return
		}
		if err != nil && err != redis.Nil {
			klog.CtxErrorf(ctx, "key %s userAssetRepos.GetAsset error:%s", key, err.Error())
			return
		}

		data, err = m.GetUserAssetByType(ctx, userId, assetType)
		if err != nil {
			klog.CtxErrorf(ctx, "userId:%d assetType:%s GetUserAssetByType error:%s", userId, assetType, err.Error())
			return
		}
		err = m.userAssetRepos.SetAsset(ctx, key, data.(*base.UserAsset))
		if err != nil {
			klog.CtxErrorf(ctx, "key:%s userAssetRepos.SetAsset data:%+v error:%s", key, data, err.Error())
			return
		}
		klog.CtxInfof(ctx, "set asset key %s val %+v is ok", key, data)
		return
	})

	select {
	case <-ctx.Done():
		return domain.ErrorSetAssetDoneOut
	case <-time.After(constants.TimeOutSetAssetFromCB):
		return domain.ErrorSetAssetFromCBTimeOut
	case res := <-resCh:
		return res.Err
	}
}
