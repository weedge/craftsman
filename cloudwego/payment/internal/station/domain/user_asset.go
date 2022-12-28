package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
)

type IUserAssetRepository interface {
	GetAsset(ctx context.Context, key string) (assetObj *base.UserAsset, err error)
	SetAsset(ctx context.Context, key string, assetObj *base.UserAsset) (err error)
}
