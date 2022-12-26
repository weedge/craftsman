package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

type UserAssetEvent struct {
	*station.BizEventAssetChange
}

type IUserAssetEventRepository interface {
	ChangeUsersAssetTx(ctx context.Context, event *station.BizEventAssetChange) error
}

type IUserAssetEventUseCase interface {
	ChangeUsersAssetTx(ctx context.Context, event *station.BizEventAssetChange) error
}

type AssetChangeHandler func(oldAssetCn int64) (newAssetCn int64)
