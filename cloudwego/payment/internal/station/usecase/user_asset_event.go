package usecase

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
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

	err = m.userAssetEventMsgRepos.SendUserAssetChangeMsgTx(ctx, constants.TopicUserAssetChange, event.EventType.String(), changeInfo.UserId, event, func(ctx context.Context) primitive.LocalTransactionState {
		userAsset, err = m.userAssetEventRepository.UserAssetChangeTx(ctx, event.EventId, changeInfo, handle)
		if err != nil {
			return primitive.RollbackMessageState
		}

		return primitive.CommitMessageState
	})
	if err != nil {
		return
	}

	return
}
