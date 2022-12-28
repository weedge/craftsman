package consumer

import (
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
)

// Init register event handler
func Init(opts map[string]*subscriber.RmqPushConsumerOptions) (mapSubscribeHandler map[string]subscriber.IRocketMQConsumerSubscribeHandler) {
	return map[string]subscriber.IRocketMQConsumerSubscribeHandler{
		constants.SendGiftConsumerEventName:    &UserAssetChangeEvent{opt: opts[constants.SendGiftConsumerEventName]},
		constants.OrderAppleConsumerEventName:  &UserAssetChangeEvent{opt: opts[constants.SendGiftConsumerEventName]},
		constants.OrderAlipayConsumerEventName: &UserAssetChangeEvent{opt: opts[constants.SendGiftConsumerEventName]},
		constants.OrderWXConsumerEventName:     &UserAssetChangeEvent{opt: opts[constants.SendGiftConsumerEventName]},
		constants.OrderDouyinConsumerEventName: &UserAssetChangeEvent{opt: opts[constants.SendGiftConsumerEventName]},
	}
}
