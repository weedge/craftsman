package rmq

import (
	"context"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"go.opentelemetry.io/otel/trace"
)

type UserAssetEventMsgLister struct {
	userAssetEventRepository domain.IUserAssetEventRepository
}

func NewUserAssetEventMsgLister(userAssetEventRepository domain.IUserAssetEventRepository) domain.IUserAssetEventMsgListener {
	return &UserAssetEventMsgLister{userAssetEventRepository: userAssetEventRepository}
}

// Deprecated: ExecuteLocalTransaction do nothing just adapter older interface,
// use primitive.Message.WithTxHandler instead
func (m *UserAssetEventMsgLister) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

// CheckLocalTransaction check user asset event msg is ok
func (m *UserAssetEventMsgLister) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	eventId := msg.GetProperty("eventId")
	userId, _ := strconv.ParseInt(msg.GetProperty("userId"), 10, 64)
	spanStr := msg.GetProperty(constants.MqTraceSpanKey)
	spanConf := trace.SpanContextConfig{}
	sonic.Unmarshal(utils.StringToSliceByte(spanStr), &spanConf)

	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(spanConf))
	res, err := m.userAssetEventRepository.GetUserAssetEventMsg(ctx, userId, eventId)
	if err != nil {
		return primitive.RollbackMessageState
	}

	if res > 0 {
		return primitive.CommitMessageState
	}

	return primitive.UnknowState
}
