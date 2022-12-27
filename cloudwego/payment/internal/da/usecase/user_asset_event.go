package usecase

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/repository/mysql"
)

type UserAssetEventUseCase struct {
	userAssetEventRepos  *mysql.UserAssetEventRepository
	userAssetRecordRepos *mysql.UserAssetRecordRepository
}

func NewUserAssetEventUseCase(userAssetEventRepos *mysql.UserAssetEventRepository, userAssetRecordRepos *mysql.UserAssetRecordRepository) domain.IUserAssetEventUseCase {

	return &UserAssetEventUseCase{
		userAssetEventRepos:  userAssetEventRepos,
		userAssetRecordRepos: userAssetRecordRepos,
	}
}

func (m *UserAssetEventUseCase) ChangeUsersAssetTx(ctx context.Context, event *station.BizEventAssetChange) (err error) {
	records, err := m.userAssetRecordRepos.GetRecordsByUserChangeAssetEvent(ctx, event)
	if err != nil {
		return
	}

	if len(records) > 0 {
		return
	}

	err = m.userAssetEventRepos.ChangeUsersAssetTx(ctx, event)
	if err != nil {
		return
	}

	return
}
