package redis

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
)

type UserAssetRepository struct {
	redisClient redis.UniversalClient
}

func NewUserAssetRepository(redisClient redis.UniversalClient) domain.IUserAssetRepository {
	return &UserAssetRepository{redisClient: redisClient}
}

func (m *UserAssetRepository) GetAsset(ctx context.Context, key string) (assetObj *base.UserAsset, err error) {
	assetObj = &base.UserAsset{}
	resCmd := m.redisClient.Get(ctx, key)
	res, err := resCmd.Bytes()
	if err != nil {
		return
	}

	err = sonic.Unmarshal(res, assetObj)

	return
}

func (m *UserAssetRepository) SetAsset(ctx context.Context, key string, assetObj *base.UserAsset) (err error) {
	res, err := sonic.Marshal(*assetObj)
	if err != nil {
		return
	}
	err = m.redisClient.Set(ctx, key, res, constants.RedisKeyExpireUserAssetInfo).Err()

	return
}
