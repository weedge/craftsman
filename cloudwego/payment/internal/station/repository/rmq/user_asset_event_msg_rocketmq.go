package rmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"go.opentelemetry.io/otel/trace"
)

type UserAssetEventMsg struct {
	txRmqProducerClient rocketmq.TransactionProducer
	topicName           string
	tagName             string
}

func NewUserAssetEventMsg(txRmqProducerClient rocketmq.TransactionProducer, topicName string, tagName string) domain.IUserAssetEventMsgRepository {
	return &UserAssetEventMsg{
		txRmqProducerClient: txRmqProducerClient,
		topicName:           topicName,
		tagName:             tagName,
	}
}

func (m *UserAssetEventMsg) SendUserAssetChangeMsgTx(ctx context.Context, event *station.BizEventAssetChange) (err error) {

	rawMsg, err := sonic.Marshal(event)
	if err != nil {
		klog.CtxErrorf(ctx, "json Marshal err:%s", err.Error())
		return
	}

	//todo: add otel tracing
	span := trace.SpanFromContext(ctx)
	msg := primitive.NewMessage(m.topicName, rawMsg)
	msg.WithTag(m.tagName)
	msg.WithKeys([]string{event.EventId})
	msg.WithProperties(map[string]string{
		"eventId":              event.EventId,
		"eventType":            event.EventType.String(),
		constants.MqTraceIdKey: span.SpanContext().TraceID().String(),
	})

	res, err := m.txRmqProducerClient.SendMessageInTransaction(ctx, msg)
	if err != nil {
		klog.CtxErrorf(ctx, "SendMessageInTransactionErr msg:%s err:%s", msg.String(), err.Error())
		return domain.ErrorSendAssetChangeEventMsg
	}
	klog.CtxInfof(ctx, "SendMessageInTransaction msg:%s ok, res:%s", msg.String(), res.String())

	return
}
