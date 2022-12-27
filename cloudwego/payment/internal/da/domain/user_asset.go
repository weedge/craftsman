package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

type IUserAssetRepository interface {
	GetUserAssets(ctx context.Context, userIds []int64) ([]*model.UserAsset, error)
}
