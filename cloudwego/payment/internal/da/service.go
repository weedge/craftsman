package da

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"go.opentelemetry.io/otel/baggage"
)

type impl struct {
}

func New() da.PaymentService {
	return &impl{}
}

func (i *impl) GetAsset(ctx context.Context, req *da.GetAssetReq) (resp *da.GetAssetResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	return
}
