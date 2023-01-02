package gw

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	commonConstants "github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/internal/gw/router"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

type Server struct {
	opts          *ServerOptions
	kitexKVLogger logutils.IKitexZapKVLogger
}

// ServerOptions server options
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
		Addr:                      ":8001",
		LogLevel:                  logutils.LevelDebug,
		ProjectName:               constants.ProjectName,
		LogMeta:                   map[string]interface{}{},
		OltpGrpcCollectorEndpoint: ":4317",
	}
}

// Run hertz server
func (s *Server) Run(ctx context.Context) error {
	klog.SetLogger(s.kitexKVLogger)
	klog.SetLevel(s.opts.LogLevel.KitexLogLevel())

	p := provider.NewOpenTelemetryProvider(
		provider.WithExportEndpoint(s.opts.OltpGrpcCollectorEndpoint),
		provider.WithEnableMetrics(true),
		provider.WithEnableTracing(true),
		provider.WithServiceName(strings.Join([]string{s.opts.ProjectName, commonConstants.UIGateWayServiceName}, ".")),
		provider.WithInsecure(),
	)

	tracer, cfg := tracing.NewServerTracer()
	h := server.Default(
		tracer,
		server.WithHostPorts(s.opts.Addr),
	)
	h.Use(tracing.ServerMiddleware(cfg))

	router.SetupRoutes(h)

	h.Spin()

	return p.Shutdown(ctx)
}
