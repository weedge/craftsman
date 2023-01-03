package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/dao"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/model"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserAssetEventRepository struct {
	db *gorm.DB
}

func NewUserAssetEventRepository(db *gorm.DB) domain.IUserAssetEventRepository {
	return &UserAssetEventRepository{db: db}
}

func (m *UserAssetEventRepository) ChangeUsersAssetTx(ctx context.Context, event *station.BizEventAssetChange) (err error) {
	if event == nil || event.OpUserAssetChange == nil {
		return domain.ErrInnerNilPointer
	}

	err = dao.Use(m.db).Transaction(func(tx *dao.Query) error {
		qTx := &dao.QueryTx{Query: tx}
		err := m.changeUserAssetTx(ctx, qTx, event.OpUserAssetChange, func(oldAssetCn int64) (newAssetCn int64) {
			return oldAssetCn + int64(event.OpUserAssetChange.Incr)
		})
		if err != nil {
			return err
		}

		err = m.changeUserAssetTx(ctx, &dao.QueryTx{Query: tx}, event.ToUserAssetChange, func(oldAssetCn int64) (newAssetCn int64) {
			return oldAssetCn + int64(event.ToUserAssetChange.Incr)
		})
		if err != nil {
			return err
		}

		records := []*model.UserAssetRecord{
			{
				UserID:     event.OpUserAssetChange.UserId,
				OpUserType: constants.OpUserTypeActive,
				BizID:      event.BizId,
				BizType:    int32(event.BizType),
				ObjID:      event.ObjId,
				EventID:    event.EventId,
				EventType:  event.EventType.String(),
				CreatedAt:  time.Now(),
			},
		}
		if event.ToUserAssetChange != nil {
			records = append(records, &model.UserAssetRecord{
				UserID:     event.ToUserAssetChange.UserId,
				OpUserType: constants.OpUserTypePassive,
				BizID:      event.BizId,
				BizType:    int32(event.BizType),
				ObjID:      event.ObjId,
				EventID:    event.EventId,
				EventType:  event.EventType.String(),
				CreatedAt:  time.Now(),
			})
		}

		err = m.addUserAssetRecordTx(ctx, qTx, records)
		if err != nil {
			return err
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	if err != nil {
		return err
	}

	return nil
}

func (m *UserAssetEventRepository) changeUserAssetTx(ctx context.Context, tx *dao.QueryTx, changeInfo *station.UserAssetChangeInfo, handle domain.AssetChangeHandler) (err error) {
	if changeInfo == nil {
		return
	}

	userAssetModle, err := tx.UserAsset.WithContext(ctx).Select(tx.UserAsset.ALL).Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where(tx.UserAsset.UserID.Eq(changeInfo.UserId), tx.UserAsset.AssetType.Eq(int32(changeInfo.AssetType))).Take()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		userAssetModle = &model.UserAsset{}
	}

	newAssetCn := handle(userAssetModle.AssetCn)
	if newAssetCn < 0 {
		return domain.ErrNoEnoughAsset
	}

	err = tx.UserAsset.WithContext(ctx).Clauses(clause.OnConflict{
		//Columns: []clause.Column{{Name: "userId"}, {Name: "assetType"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"assetCn": newAssetCn,
			"version": userAssetModle.Version + 1,
		}),
	}).Create(&model.UserAsset{
		UserID:    changeInfo.UserId,
		AssetCn:   newAssetCn,
		AssetType: int32(changeInfo.AssetType),
		Version:   1,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	return
}

func (m *UserAssetEventRepository) addUserAssetRecordTx(ctx context.Context, tx *dao.QueryTx, records []*model.UserAssetRecord) (err error) {
	err = tx.UserAssetRecord.WithContext(ctx).Create(records...)

	return
}
