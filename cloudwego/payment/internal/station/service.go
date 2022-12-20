package station

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"go.opentelemetry.io/otel/baggage"
)

type impl struct {
}

func New() station.PaymentService {
	return &impl{}
}

func (i *impl) ChangeAsset(ctx context.Context, req *station.BizAssetChangesReq) (resp *station.BizAssetChangesResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	return
}
