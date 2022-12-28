package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

type AssetIncrHandler func(ctx context.Context) (incrAssetCn int64)

type IUserAssetEventRepository interface {
	UserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle AssetIncrHandler) error
}

type IUserAssetEventUseCase interface {
	UserAssetChangeTx(ctx context.Context, eventId string, changeInfo *station.UserAssetChangeInfo, handle AssetIncrHandler) error
}
