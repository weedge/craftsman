//go:build wireinject
// +build wireinject

package da

import (
	"context"

	"github.com/google/wire"
	"github.com/weedge/craftsman/cloudwego/kitex-contrib/gorm"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

// NewServer build server with wire
func NewServer(ctx context.Context) (*Server, error) {
	panic(wire.Build(
		configparser.Default,
		Configure,

		wire.FieldsOf(new(*Options), "Server", "MysqlDBClient"),
		wire.FieldsOf(new(*ServerOptions), "LogLevel", "LogMeta"),

		logutils.NewkitexZapKVLogger,
		wire.Bind(new(gorm.IkvLogger), new(logutils.IKitexZapKVLogger)),

		injectors.InitMysqlDBClient,
		NewSvc,

		wire.Struct(new(Server), "*"),
	))
}
