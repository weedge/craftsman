package rmq

import (
	"context"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
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

var MsgTracePayload struct {
	TraceID trace.TraceID    `json:"ID"`
	SpanID  trace.SpanID     `json:"SpanID"`
	Flags   trace.TraceFlags `json:"TraceFlags"`
}

// CheckLocalTransaction check user asset event msg is ok
func (m *UserAssetEventMsgLister) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	eventId := msg.GetProperty("eventId")
	userId, _ := strconv.ParseInt(msg.GetProperty("userId"), 10, 64)
	spanStr := msg.GetProperty(constants.MqTraceSpanKey)

	ctx := context.Background()
	if len(spanStr) > 0 {
		spanConf := trace.SpanContextConfig{}
		err := sonic.Unmarshal(utils.StringToSliceByte(spanStr), &spanConf)
		if err != nil {
			klog.CtxErrorf(ctx, "sonic.Unmarshal err:%s", err.Error())
		}
		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(spanConf))
	}

	res, err := m.userAssetEventRepository.GetUserAssetEventMsg(ctx, userId, eventId)
	if err != nil {
		klog.CtxErrorf(ctx, "userAssetEventRepository.GetUserAssetEventMsg err:%s", err.Error())
		return primitive.RollbackMessageState
	}

	if res > 0 {
		return primitive.CommitMessageState
	}

	return primitive.UnknowState
}
