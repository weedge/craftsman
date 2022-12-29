package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

type UserAssetEventMsg struct {
	*station.BizEventAssetChange
}

type IUserAssetEventMsgRepository interface {
	SendUserAssetChangeMsgTx(ctx context.Context, eventMsg *station.BizEventAssetChange) error
}
