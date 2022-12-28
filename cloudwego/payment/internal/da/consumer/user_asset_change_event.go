package consumer

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
)

type UserAssetChangeEvent struct {
	opt                   *subscriber.RmqPushConsumerOptions
	userAssetEventUseCase domain.IUserAssetEventUseCase
}

func (m *UserAssetChangeEvent) SubMsgsHandle(ctx context.Context, msgs ...*primitive.MessageExt) (res consumer.ConsumeResult, err error) {
	concurrentCtx, _ := primitive.GetConcurrentlyCtx(ctx)
	concurrentCtx.DelayLevelWhenNextConsume = m.opt.DelayLevel // only run when return consumer.ConsumeRetryLater

	for _, msg := range msgs {
		if m.opt.LogicMaxRetryCn > 0 && msg.ReconsumeTimes > int32(m.opt.LogicMaxRetryCn) {
			log.Printf("msg ReconsumeTimes > %d. msg: %v return consumer success", m.opt.LogicMaxRetryCn, msg)
			return consumer.ConsumeSuccess, nil
		} else {
			log.Printf("subscribe callback: %v \n", msg)

			// change user asset
			event := &station.BizEventAssetChange{}
			err = sonic.Unmarshal(msg.Body, event)
			if err != nil {
				klog.CtxErrorf(ctx, "json decode err: %s", err.Error())
				return consumer.ConsumeSuccess, nil
			}
			err = m.userAssetEventUseCase.ChangeUsersAssetTx(ctx, event)

			if err == nil {
				return consumer.ConsumeSuccess, nil
			}
		}
	}

	return consumer.ConsumeRetryLater, err

}
