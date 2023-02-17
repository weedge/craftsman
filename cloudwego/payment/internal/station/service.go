package station

import (
	"context"
	"fmt"
	"time"

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

	if len(req.BizAssetChanges) > constants.StationChangeMaxAssetsCn {
		klog.CtxWarnf(ctx, "ChangeAsset len(BizAssetChanges) > %d ", constants.StationChangeMaxAssetsCn)
		resp.BaseResp.SetErrCode(int64(common.Err_PaymentBadRequest))
		resp.BaseResp.SetErrMsg(common.Err_PaymentBadRequest.String())
		return
	}

	resp = &station.BizAssetChangesResp{
		BizAssetChangeResList: make([]*station.BizEventAssetChangerRes, len(req.BizAssetChanges)),
		BaseResp:              &commonBase.BaseResp{},
	}

	ctx, cancel := context.WithTimeout(ctx, constants.StationChangeAssetExeTimeoutS*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	for index, item := range req.BizAssetChanges {
		index, item := index, item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				timeOutExeErr := fmt.Errorf("%ds timeOutExeErr", constants.StationChangeAssetExeTimeoutS)
				resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
					EventId:     req.BizAssetChanges[index].EventId,
					ChangeRes:   false,
					FailMsg:     timeOutExeErr.Error(),
					OpUserAsset: nil,
				}
				klog.CtxErrorf(ctx, "UserAssetChangeTx item:%+v err:%s", item, timeOutExeErr.Error())
				return nil
			default:
				userAsset, txErr := i.user.UserAssetChangeTx(ctx, constants.OpUserTypeActive, item, func(ctx context.Context) (incrAssetCn int64) {
					incrAssetCn = int64(item.OpUserAssetChange.Incr)
					return
				})
				if txErr != nil {
					resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
						EventId:     req.BizAssetChanges[index].EventId,
						ChangeRes:   false,
						FailMsg:     txErr.Error(),
						OpUserAsset: nil,
					}
					if txErr == domain.ErrorNoEnoughAsset {
						klog.CtxWarnf(ctx, "UserAssetChangeTx item:%+v err:%s", item, txErr.Error())
						return nil
					}
					klog.CtxErrorf(ctx, "UserAssetChangeTx item:%+v err:%s", item, txErr.Error())

					//return txErr

					// notice: don't need to return err, err save to BizEventAssetChangerRes, eventId transaction :)
					return nil
				}
				resp.BizAssetChangeResList[index] = &station.BizEventAssetChangerRes{
					EventId:     req.BizAssetChanges[index].EventId,
					ChangeRes:   true,
					FailMsg:     "",
					OpUserAsset: userAsset,
				}

				return nil
			}
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
