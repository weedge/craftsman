package repository

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

type IUserAssetEventRepository interface {
	GetUserAssets(ctx context.Context, userIds []int64) ([]*model.UserAsset, error)
	ChangeUsersAssetTx(ctx context.Context, opUserId int64, toUserId int64) error
}

type UserAssetEventRepository struct {
}

func (m *UserAssetEventRepository) ChangeUsersAssetTx(ctx context.Context, opUserId int64, toUserId int64) error {

	return nil
}
