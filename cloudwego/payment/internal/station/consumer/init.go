package consumer

import (
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
)

// Register event handler
func RegisterUserAssetEvent(opts map[string]*subscriber.RmqPushConsumerOptions, userAssetEventUseCase domain.IUserAssetEventUseCase) (mapSubscribeHandler map[string]subscriber.IRocketMQConsumerSubscribeHandler) {
	return map[string]subscriber.IRocketMQConsumerSubscribeHandler{
		constants.SendGiftConsumerEventName: &UserAssetChangeEvent{
			opt:                   opts[constants.SendGiftConsumerEventName],
			userAssetEventUseCase: userAssetEventUseCase,
		},
		constants.OrderAppleConsumerEventName: &UserAssetChangeEvent{
			opt:                   opts[constants.OrderAppleConsumerEventName],
			userAssetEventUseCase: userAssetEventUseCase,
		},
		constants.OrderAlipayConsumerEventName: &UserAssetChangeEvent{
			opt:                   opts[constants.OrderAlipayConsumerEventName],
			userAssetEventUseCase: userAssetEventUseCase,
		},
		constants.OrderWXConsumerEventName: &UserAssetChangeEvent{
			opt:                   opts[constants.OrderWXConsumerEventName],
			userAssetEventUseCase: userAssetEventUseCase,
		},
		constants.OrderDouyinConsumerEventName: &UserAssetChangeEvent{
			opt:                   opts[constants.OrderDouyinConsumerEventName],
			userAssetEventUseCase: userAssetEventUseCase,
		},
	}
}
