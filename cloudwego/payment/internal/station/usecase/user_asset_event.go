package usecase

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
)

type UserAssetEventUseCase struct {
	userAssetEventMsgRepos   domain.IUserAssetEventMsgRepository
	userAssetEventRepository domain.IUserAssetEventRepository
}

func NewUserAssetEventUseCase(userAssetEventMsgRepos domain.IUserAssetEventMsgRepository, userAssetEventRepository domain.IUserAssetEventRepository) domain.IUserAssetEventUseCase {
	return &UserAssetEventUseCase{
		userAssetEventMsgRepos:   userAssetEventMsgRepos,
		userAssetEventRepository: userAssetEventRepository,
	}
}

func (m *UserAssetEventUseCase) UserAssetChangeTx(ctx context.Context, opUserType int, event *station.BizEventAssetChange, handle domain.AssetIncrHandler) (userAsset *base.UserAsset, err error) {
	topicName := constants.TopicUserAssetChange
	if opUserType == constants.OpUserTypeActive &&
		event.OpUserAssetChange != nil && event.OpUserAssetChange.UserId > 0 &&
		event.ToUserAssetChange != nil && event.ToUserAssetChange.UserId > 0 {
		topicName = constants.TopicCacheUserAssetChange
	}

	changeInfo := event.OpUserAssetChange
	if opUserType == constants.OpUserTypePassive {
		changeInfo = event.ToUserAssetChange
	}

	res, err := m.userAssetEventRepository.GetUserAssetEventMsg(ctx, changeInfo.UserId, event.EventId)
	if err != nil {
		return
	}

	// have event msg, don't do tx again, return done
	if res > 0 {
		err = domain.ErrorEventDone
		return
	}

	var errTx error
	err = m.userAssetEventMsgRepos.SendUserAssetChangeMsgTx(ctx, topicName, event.EventType.String(), changeInfo.UserId, event, func(ctx context.Context) primitive.LocalTransactionState {
		userAsset, errTx = m.userAssetEventRepository.UserAssetChangeTx(ctx, event.EventId, changeInfo, handle)
		if errTx != nil {
			return primitive.RollbackMessageState
		}

		return primitive.CommitMessageState
	})

	// check send msg err
	if err != nil {
		return
	}

	// check tx err
	if errTx != nil {
		err = errTx
		return
	}

	klog.CtxDebugf(ctx, "UserAssetChangeTx ok, return userAsset:%+v ", userAsset)

	return
}
