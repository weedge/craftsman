package station

import (
	"context"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	commonBase "github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/sync/errgroup"
)

type impl struct {
	user domain.IUserAssetEventUseCase
}

func NewService(user domain.IUserAssetEventUseCase) station.PaymentService {
	return &impl{user: user}
}

func (i *impl) ChangeAsset(ctx context.Context, req *station.BizAssetChangesReq) (resp *station.BizAssetChangesResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	resp = &station.BizAssetChangesResp{
		BizAssetChangeResList: make([]*station.BizEventAssetChangerRes, len(req.BizAssetChanges)),
		BaseResp:              &commonBase.BaseResp{},
	}

	var mutex sync.RWMutex
	eg, ctx := errgroup.WithContext(ctx)
	for index, item := range req.BizAssetChanges {
		index, item := index, item
		eg.Go(func() error {
			userAsset, txErr := i.user.UserAssetChangeTx(ctx, constants.OpUserTypeActive, item, func(ctx context.Context) (incrAssetCn int64) {
				incrAssetCn = int64(item.OpUserAssetChange.Incr)
				return
			})
			if txErr != nil {
				mutex.Lock()
				resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
					EventId:     req.BizAssetChanges[index].EventId,
					ChangeRes:   false,
					FailMsg:     txErr.Error(),
					OpUserAsset: nil,
				}
				mutex.Unlock()
				if txErr == domain.ErrorNoEnoughAsset {
					klog.CtxWarnf(ctx, "UserAssetChangeTx item:%+v err:%s", item, txErr.Error())
					return nil
				}
				klog.CtxErrorf(ctx, "UserAssetChangeTx item:%+v err:%s", item, txErr.Error())

				return txErr
			}
			mutex.Lock()
			resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
				EventId:     req.BizAssetChanges[index].EventId,
				ChangeRes:   true,
				FailMsg:     "",
				OpUserAsset: userAsset,
			}
			mutex.Unlock()

			return nil
		})
	}

	gErr := eg.Wait()
	if gErr != nil {
		//klog.CtxErrorf(ctx, "ChangeAssets err:%s", gErr.Error())
		resp.BaseResp.SetErrCode(int64(common.Err_PaymentStationInteralError))
		resp.BaseResp.SetErrMsg(common.Err_PaymentStationInteralError.String())
		resp.BaseResp.SetExtra(map[string]string{"err": gErr.Error()})
		return
	}

	return
}
