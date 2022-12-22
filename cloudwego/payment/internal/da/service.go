package da

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
	base0 "github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/dao"
	"go.opentelemetry.io/otel/baggage"
	"gorm.io/gorm"
)

type impl struct {
	dbClient *gorm.DB
}

func NewSvc(dbClient *gorm.DB) da.PaymentService {
	return &impl{dbClient: dbClient}
}

func (i *impl) GetAssets(ctx context.Context, req *da.GetAssetsReq) (resp *da.GetAssetsResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	resp = &da.GetAssetsResp{UserAssets: []*base0.UserAsset{}, BaseResp: &base.BaseResp{}}

	userAssetDao := dao.Use(i.dbClient).UserAsset
	items, err := userAssetDao.WithContext(ctx).Where(userAssetDao.UserID.In(req.UserIds...)).Find()
	if err != nil && gorm.ErrRecordNotFound != err {
		klog.CtxErrorf(ctx, "GetAssets err:%s", err.Error())
		resp.BaseResp.SetErrCode(int64(common.Err_PaymentDbInteralError))
		resp.BaseResp.SetErrMsg(common.Err_PaymentDbInteralError.String())
		resp.BaseResp.SetExtra(map[string]string{"err": err.Error()})
		return
	}

	for _, item := range items {
		resp.UserAssets = append(resp.UserAssets, &base0.UserAsset{
			UserId:    item.UserID,
			AssetType: base0.AssetType(item.AssetType),
			AssetCn:   item.AssetCn,
		})
	}

	return
}
