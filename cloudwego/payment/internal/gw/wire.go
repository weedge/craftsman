package gw

import (
	"context"

	"github.com/google/wire"
	"github.com/weedge/craftsman/cloudwego/payment/internal/gw/handler"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

// NewServer build server with wire, dependency obj inject, so init random
func NewServer(ctx context.Context) (*Server, error) {
	panic(wire.Build(
		configparser.Default,
		Configure,
		wire.FieldsOf(new(*Options),
			"Server",
			"HttpThriftGenericClient",
		),

		wire.FieldsOf(new(*ServerOptions), "LogLevel", "LogMeta"),
		logutils.NewkitexZapKVLogger,

		injectors.InitHttpThriftGenericClients,
		handler.InitSvcGenericClientMap,

		wire.Struct(new(Server), "*"),
	))
}
