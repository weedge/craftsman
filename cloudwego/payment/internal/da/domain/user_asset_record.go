package domain

import (
	"context"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
)

type IUserAssetRecordRepository interface {
	GetRecordsByUserId(ctx context.Context, userId int64) (res []*model.UserAssetRecord, nextCursor string, err error)
	GetRecordsByUserChangeAssetEvent(ctx context.Context, event *station.BizEventAssetChange) (res []*model.UserAssetRecord, err error)
}
