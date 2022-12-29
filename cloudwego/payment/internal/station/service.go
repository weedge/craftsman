package station

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	commonBase "github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/sync/errgroup"
)

type impl struct {
	user domain.IUserAssetEventUseCase
}

func New(user domain.IUserAssetEventUseCase) station.PaymentService {
	return &impl{user: user}
}

func (i *impl) ChangeAsset(ctx context.Context, req *station.BizAssetChangesReq) (resp *station.BizAssetChangesResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	resp = &station.BizAssetChangesResp{
		BizAssetChangeResList: make([]*station.BizEventAssetChangerRes, len(req.BizAssetChanges)),
		BaseResp:              &commonBase.BaseResp{},
	}

	eg, ctx := errgroup.WithContext(ctx)
	for index, item := range req.BizAssetChanges {
		index := index
		item := item
		eg.Go(func() error {
			userAsset, txErr := i.user.UserAssetChangeTx(ctx, item, func(ctx context.Context) (incrAssetCn int64) {
				incrAssetCn = -int64(item.OpUserAssetChange.Incr)
				return
			})
			if txErr != nil {
				klog.CtxErrorf(ctx, "UserAssetChangeTx item:%+v err:%s", item, txErr.Error())
				resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
					EventId:     req.BizAssetChanges[index].EventId,
					ChangeRes:   false,
					FailMsg:     txErr.Error(),
					OpUserAsset: nil,
				}
				return txErr
			}
			resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
				EventId:     req.BizAssetChanges[index].EventId,
				ChangeRes:   true,
				FailMsg:     "",
				OpUserAsset: userAsset,
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		klog.CtxErrorf(ctx, "ChangeAssets err:%s", err.Error())
		resp.BaseResp.SetErrCode(int64(common.Err_PaymentStationInteralError))
		resp.BaseResp.SetErrMsg(common.Err_PaymentBadRequest.String())
		resp.BaseResp.SetExtra(map[string]string{"err": err.Error()})
		return
	}

	return
}
