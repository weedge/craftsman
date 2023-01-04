package consumer

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/containers"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
)

type UserAssetChangeEvent struct {
	opt                   *subscriber.RmqPushConsumerOptions
	userAssetEventUseCase domain.IUserAssetEventUseCase
}

func (m *UserAssetChangeEvent) SubMsgsHandle(ctx context.Context, msgs ...*primitive.MessageExt) (res consumer.ConsumeResult, err error) {
	concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
	concurrentCtx.DelayLevelWhenNextConsume = m.opt.DelayLevel // only run when return consumer.ConsumeRetryLater

	mapTxErr := containers.NewErrMap()
	for _, msg := range msgs {
		if m.opt.LogicMaxRetryCn > 0 && msg.ReconsumeTimes > int32(m.opt.LogicMaxRetryCn) {
			klog.CtxWarnf(ctx, "msg.ReconsumeTimes %d > %d continue don't consume. msg: %+v", msg.ReconsumeTimes, m.opt.LogicMaxRetryCn, msg)
			continue
		}

		klog.CtxInfof(ctx, "subscribe callback: %s ", msg.Body)

		// decode msg body get event
		event := &station.BizEventAssetChange{}
		err = sonic.Unmarshal(msg.Body, event)
		if err != nil {
			klog.CtxErrorf(ctx, "json decode err: %s", err.Error())
			continue
		}

		// change user asset
		err = m.userAssetEventUseCase.ChangeUsersAssetTx(ctx, event)
		if err != nil {
			if err != domain.ErrInnerNilPointer {
				mapTxErr.Add(msg.MsgId, err)
			} else {
				klog.CtxErrorf(ctx, "userAssetEventUseCase.ChangeUsersAssetTx err: %s", err.Error())
			}
			continue
		}
	}

	if mapTxErr.Len() > 0 {
		klog.CtxErrorf(ctx, "userAssetEventUseCase.ChangeUsersAssetTx err: %s", mapTxErr.String())
		return consumer.ConsumeRetryLater, nil
	}

	klog.CtxDebugf(ctx, "userAssetEventUseCase.ChangeUsersAssetTx msgs:%+v ok", msgs)

	return consumer.ConsumeSuccess, nil
}
