package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

type AssetIncrHandler func(ctx context.Context) (incrAssetCn int64)

type IUserAssetEventRepository interface {
	UserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle AssetIncrHandler) (*base.UserAsset, error)
	GetUserAssetEventMsg(ctx context.Context, userId int64, eventId string) (int, error)
}

type IUserAssetEventUseCase interface {
	UserAssetChangeTx(ctx context.Context, opUserType int, event *station.BizEventAssetChange, handle AssetIncrHandler) (*base.UserAsset, error)
}
