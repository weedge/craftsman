package router

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	commonConstants "github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/internal/gw/handler"
	"github.com/weedge/craftsman/cloudwego/payment/internal/gw/middleware"
)

func SetupRoutes(h *server.Hertz) {
	SetupProbeRoutes(h)

	SetupHttpThriftGenericRoutes(h)
}

func SetupProbeRoutes(h *server.Hertz) {
	// ready
	h.GET("/readiness", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, commonConstants.UIGateWayServiceName+" is readiness")
	})

	// liveness
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, commonConstants.UIGateWayServiceName+" is running")
	})
}

// SetupHttpThriftGenericRoutes
func SetupHttpThriftGenericRoutes(h *server.Hertz) {
	passGroup := h.Group("/openapi").Use(middleware.OpenApiAuth())
	passGroup.POST("/:svc", handler.OpenApiHandle)
}
