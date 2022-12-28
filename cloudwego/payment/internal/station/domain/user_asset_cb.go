package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
)

type IUserAssetCallBack interface {
	GetAssets(ctx context.Context, userIds []int64) (assetDtos []*base.UserAsset, err error)
	SetAsset(ctx context.Context, userId int64, assetType base.AssetType) (err error)
}
