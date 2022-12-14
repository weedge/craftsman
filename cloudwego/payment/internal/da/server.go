package da

import (
	"context"
	"net"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/da/paymentservice"
	commonConstants "github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

type Server struct {
	opts                *ServerOptions
	svc                 da.PaymentService
	kitexKVLogger       logutils.IKitexZapKVLogger
	rmqConsumerOpts     map[string]*subscriber.RmqPushConsumerOptions
	mapSubscribeHandler map[string]subscriber.IRocketMQConsumerSubscribeHandler
}

type ServerOptions struct {
	Addr                      string                 `mapstructure:"addr"`
	LogLevel                  logutils.Level         `mapstructure:"logLevel"`
	ProjectName               string                 `mapstructure:"projectName"`
	LogMeta                   map[string]interface{} `mapstructure:"logMeta"`
	OltpGrpcCollectorEndpoint string                 `mapstructure:"oltpCollectorGrpcEndpoint"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Addr:                      ":8003",
		LogLevel:                  logutils.LevelDebug,
		ProjectName:               constants.ProjectName,
		LogMeta:                   map[string]interface{}{},
		OltpGrpcCollectorEndpoint: ":4317",
	}
}

// Run reviews server
func (s *Server) Run(ctx context.Context) error {
	klog.SetLogger(s.kitexKVLogger)
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

	subscriber.InitPushConsumerSubscribes(s.rmqConsumerOpts, s.mapSubscribeHandler)
	defer subscriber.Close()

	tracingProvider := provider.NewOpenTelemetryProvider(
		provider.WithExportEndpoint(s.opts.OltpGrpcCollectorEndpoint),
		provider.WithEnableMetrics(true),
		provider.WithEnableTracing(true),
		provider.WithServiceName(strings.Join([]string{s.opts.ProjectName, commonConstants.MisDaServiceName}, ".")),
		provider.WithInsecure(),
	)
	defer func(ctx context.Context, p provider.OtelProvider) {
		_ = p.Shutdown(ctx)
	}(ctx, tracingProvider)

	// for service eg: k8s service
	addr, err := net.ResolveTCPAddr("tcp", s.opts.Addr)
	if err != nil {
		klog.Fatal(err)
	}
	svr := paymentservice.NewServer(
		s.svc,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
	)
	if err := svr.Run(); err != nil {
		klog.Fatal(err)
	}

	return nil
}
