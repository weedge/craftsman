package rmq

import (
	"context"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"go.opentelemetry.io/otel/trace"
)

type UserAssetEventMsg struct {
	txRmqProducerClient      rocketmq.TransactionProducer
	userAssetEventRepository domain.IUserAssetEventRepository
}

func NewUserAssetEventMsg(txRmqProducerClient rocketmq.TransactionProducer, userAssetEventRepository domain.IUserAssetEventRepository) domain.IUserAssetEventMsgRepository {
	return &UserAssetEventMsg{
		txRmqProducerClient:      txRmqProducerClient,
		userAssetEventRepository: userAssetEventRepository,
	}
}

func (m *UserAssetEventMsg) SendUserAssetChangeMsgTx(ctx context.Context, topicName, tagName string, userId int64, event *station.BizEventAssetChange, handler primitive.TxHandler) (err error) {
	rawMsg, err := sonic.Marshal(event)
	if err != nil {
		klog.CtxErrorf(ctx, "json Marshal err:%s", err.Error())
		return
	}

	//todo: add otel tracing
	span := trace.SpanFromContext(ctx)
	spanBytes, _ := span.SpanContext().MarshalJSON()
	msg := primitive.NewMessage(topicName, rawMsg)
	msg.WithProperties(map[string]string{
		"eventId":                event.EventId,
		"userId":                 strconv.FormatInt(userId, 10),
		"eventType":              event.EventType.String(),
		constants.MqTraceSpanKey: utils.SliceByteToString(spanBytes),
	})
	msg.WithTag(tagName)
	msg.WithKeys([]string{event.EventId})
	msg.WithTxHandler(handler)

	res, err := m.txRmqProducerClient.SendMessageInTransaction(ctx, msg)
	if err != nil {
		klog.CtxErrorf(ctx, "SendMessageInTransactionErr msg:%s err:%s", msg.String(), err.Error())
		return domain.ErrorSendAssetChangeEventMsg
	}

	if res.State != primitive.CommitMessageState {
		klog.CtxWarnf(ctx, "SendMessageInTransaction msg:%s ok, res:%s, localTxState:%d; but txMsg don't RC", msg.String(), res.String(), res.State)
	} else {
		klog.CtxInfof(ctx, "SendMessageInTransaction msg:%s ok, res:%s, localTxState:%d; txMsg RC", msg.String(), res.String(), res.State)
	}

	return
}
