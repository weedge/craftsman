package injectors

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/xds"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	xdsmanager "github.com/kitex-contrib/xds"
	"github.com/kitex-contrib/xds/xdssuite"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da/paymentservice"
	"github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/common/pkg/metadata"
)

type PaymentDaClientOptions struct {
	Endpoint  string   `mapstructure:"endpoint"`
	EnableXDS bool     `mapstructure:"enableXDS"`
	XDSAddr   string   `mapstructure:"xdsAddr"`
	HostPorts []string `mapstructure:"hostPorts"`
}

func DefaultPaymentDaClientOptions() *PaymentDaClientOptions {
	return &PaymentDaClientOptions{
		Endpoint:  "payment.da:8003",
		EnableXDS: false,
		XDSAddr:   "istiod.istio-system.svc:15010",
		HostPorts: []string{":8002"},
	}
}

func InitPaymentDaClient(opts *PaymentDaClientOptions) (paymentservice.Client, error) {
	if opts.EnableXDS {
		err := xdsmanager.Init(xdsmanager.WithXDSServerAddress(opts.XDSAddr))
		if err != nil {
			klog.Fatal(err.Error())
		}
		return paymentservice.NewClient(
			opts.Endpoint,
			client.WithSuite(tracing.NewClientSuite()),
			client.WithXDSSuite(xds.ClientSuite{
				RouterMiddleware: xdssuite.NewXDSRouterMiddleware(
					xdssuite.WithRouterMetaExtractor(metadata.ExtractFromPropagator),
				),
				Resolver: xdssuite.NewXDSResolver(),
			}),
		)
	}

	return paymentservice.NewClient(
		constants.MisDaServiceName,
		client.WithHostPorts(opts.HostPorts...),
		client.WithSuite(tracing.NewClientSuite()),
	)
}
