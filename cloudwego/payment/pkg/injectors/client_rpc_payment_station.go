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

type PaymentStationClientOptions struct {
	Endpoint  string   `mapstructure:"endpoint"`
	EnableXDS bool     `mapstructure:"enableXDS"`
	XDSAddr   string   `mapstructure:"xdsAddr"`
	HostPorts []string `mapstructure:"hostPorts"`
}

func DefaultPaymentStationClientOptions() *PaymentStationClientOptions {
	return &PaymentStationClientOptions{
		Endpoint:  "payment.station:8002",
		EnableXDS: false,
		XDSAddr:   "istiod.istio-system.svc:15010",
		HostPorts: []string{":8002"},
	}
}

func InitPaymentStationClient(opts *PaymentStationClientOptions) (paymentservice.Client, error) {
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
		constants.StationServiceName,
		client.WithHostPorts(opts.HostPorts...),
		client.WithSuite(tracing.NewClientSuite()),
	)
}
