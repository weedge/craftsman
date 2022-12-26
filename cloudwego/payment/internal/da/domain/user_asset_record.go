package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

type IUserAssetRecordRepository interface {
	GetRecordsByUserId(ctx context.Context, userId int64) (res []model.UserAssetRecord, nextCursor string, err error)
}

type IUserAssetRecordUseCase interface {
	GetRecordsByUserId(ctx context.Context, userId int64) (res []model.UserAssetRecord, nextCursor string, err error)
}
