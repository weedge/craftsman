package usecase

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
)

type UserAssetEventUseCase struct {
	userAssetEventMsgRepos   domain.IUserAssetEventMsgRepository
	userAssetEventRepository domain.IUserAssetEventRepository
	opUserType               int
}

func NewUserAssetEventUseCase(userAssetEventMsgRepos domain.IUserAssetEventMsgRepository, userAssetEventRepository domain.IUserAssetEventRepository, opUserType int) domain.IUserAssetEventUseCase {
	return &UserAssetEventUseCase{
		userAssetEventMsgRepos:   userAssetEventMsgRepos,
		userAssetEventRepository: userAssetEventRepository,
		opUserType:               opUserType,
	}
}

func (m *UserAssetEventUseCase) UserAssetChangeTx(ctx context.Context, event *station.BizEventAssetChange, handle domain.AssetIncrHandler) (userAsset *base.UserAsset, err error) {
	changeInfo := event.OpUserAssetChange
	if m.opUserType == constants.OpUserTypePassive {
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

	err = m.userAssetEventMsgRepos.SendUserAssetChangeMsgTx(ctx, event)
	if err != nil {
		return
	}

	userAsset, err = m.userAssetEventRepository.UserAssetChangeTx(ctx, event.EventId, changeInfo, handle)
	if err != nil {
		return
	}

	return
}
