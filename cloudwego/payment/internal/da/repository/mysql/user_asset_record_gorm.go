package mysql

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/dao"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"gorm.io/gorm"
)

type UserAssetRecordRepository struct {
	db *gorm.DB
}

func NewUserAssetRecordRepository(db *gorm.DB) domain.IUserAssetRecordRepository {
	return &UserAssetRecordRepository{db: db}
}

func (m *UserAssetRecordRepository) GetRecordsByUserId(ctx context.Context, userId int64) (res []*model.UserAssetRecord, nextCursor string, err error) {

	return
}

func (m *UserAssetRecordRepository) GetRecordsByUserChangeAssetEvent(ctx context.Context, event *station.BizEventAssetChange) (res []*model.UserAssetRecord, err error) {
	userAssetRecordDao := dao.Use(m.db).UserAssetRecord
	userAssetRecordDao.WithContext(ctx).Where(
		userAssetRecordDao.UserID.Eq(event.OpUserAssetChange.UserId),
		userAssetRecordDao.OpUserType.Eq(constants.OpUserTypeActive),
		userAssetRecordDao.EventID.Eq(event.EventId),
	).Or(
		userAssetRecordDao.UserID.Eq(event.ToUserAssetChange.UserId),
		userAssetRecordDao.OpUserType.Eq(constants.OpUserTypePassive),
		userAssetRecordDao.EventID.Eq(event.EventId),
	).Find()

	return
}
