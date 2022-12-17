package station

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station/paymentservice"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

type Server struct {
	opts *ServerOptions
	svc  station.PaymentService
}

// ServerOptions server options
type ServerOptions struct {
	Addr     string         `mapstructure:"addr"`
	LogLevel logutils.Level `mapstructure:"logLevel"`
}

// DefaultServerOptions default opts
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Addr:     ":8002",
		LogLevel: logutils.LevelInfo,
	}
}

// Run kitex server
func (s *Server) Run(ctx context.Context) (err error) {
	klog.SetLogger(kitexzap.NewLogger())
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

	tracingProvider := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(constants.StationServiceName),
		provider.WithInsecure(),
	)
	defer func(p provider.OtelProvider, ctx context.Context) {
		_ = p.Shutdown(ctx)
	}(tracingProvider, ctx)

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

	return
}
