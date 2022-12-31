package domain

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

type UserAssetEventMsg struct {
	*station.BizEventAssetChange
}

type IUserAssetEventMsgRepository interface {
	SendUserAssetChangeMsgTx(ctx context.Context, topicName, tagName string, eventMsg *station.BizEventAssetChange, handler primitive.TxHandler) error
}

type IUserAssetEventMsgListener interface {
	primitive.TransactionListener
}
