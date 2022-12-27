package mysql

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/dao"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
	"gorm.io/gorm"
)

type UserAssetRepository struct {
	db *gorm.DB
}

func NewUserAssetRepository(db *gorm.DB) domain.IUserAssetRepository {
	return &UserAssetRepository{db: db}
}

func (m *UserAssetRepository) GetUserAssets(ctx context.Context, userIds []int64) ([]*model.UserAsset, error) {
	userAssetDao := dao.Use(m.db).UserAsset
	res, err := userAssetDao.WithContext(ctx).Where(userAssetDao.UserID.In(userIds...)).Find()

	return res, err
}
