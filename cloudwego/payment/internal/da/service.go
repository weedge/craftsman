package da

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
	base0 "github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"go.opentelemetry.io/otel/baggage"
	"gorm.io/gorm"
)

type impl struct {
	userAssetRepos domain.IUserAssetRepository
}

func NewService(userAssetRepos domain.IUserAssetRepository) da.PaymentService {
	return &impl{userAssetRepos: userAssetRepos}
}

func (i *impl) GetAssets(ctx context.Context, req *da.GetAssetsReq) (resp *da.GetAssetsResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	resp = &da.GetAssetsResp{UserAssets: []*base0.UserAsset{}, BaseResp: &base.BaseResp{}}

	if len(req.UserIds) == 0 {
		klog.CtxWarnf(ctx, "GetAssets len(UserIds) = 0 ")
		return
	}
	if len(req.UserIds) > constants.DaGetAssetsMaxUserIdCn {
		klog.CtxWarnf(ctx, "GetAssets len(UserIds) > %d ", constants.DaGetAssetsMaxUserIdCn)
		resp.BaseResp.SetErrCode(int64(common.Err_PaymentBadRequest))
		resp.BaseResp.SetErrMsg(common.Err_PaymentBadRequest.String())
		return
	}

	items, err := i.userAssetRepos.GetUserAssets(ctx, req.UserIds)
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
