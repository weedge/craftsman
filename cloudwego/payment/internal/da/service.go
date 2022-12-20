package da

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"go.opentelemetry.io/otel/baggage"
	"gorm.io/gorm"
)

type impl struct {
	dbClient *gorm.DB
}

func New(dbClient *gorm.DB) da.PaymentService {
	return &impl{dbClient: dbClient}
}

func (i *impl) GetAsset(ctx context.Context, req *da.GetAssetReq) (resp *da.GetAssetResp, err error) {
	klog.CtxInfof(ctx, "req: %+v", req)
	klog.CtxDebugf(ctx, "otel tracing baggage: %s", baggage.FromContext(ctx).String())

	return
}
