//go:build wireinject
// +build wireinject

package station

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/consumer"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/domain"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/repository"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/usecase"
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
			"AssetChangeEventProducer",
			"RedisClusterClient",
			"PaymentDaClient",
			"UserAssetTxMethod",
			"RmqConsumers",
		),
		wire.FieldsOf(new(*ServerOptions), "LogLevel", "LogMeta"),

		logutils.NewkitexZapKVLogger,

		wire.Bind(new(primitive.TransactionListener), new(domain.IUserAssetEventMsgListener)),
		injectors.InitRmqTransactionProducerClient,

		injectors.InitRedisClusterDefaultClient,
		ProvideUniversalClients,
		injectors.InitRedsync,
		injectors.InitPaymentDaClient,

		repository.ProviderSet,
		usecase.ProviderSet,
		NewService,
		consumer.RegisterUserAssetEvent,

		wire.Struct(new(Server), "*"),
	))
}

func ProvideUniversalClients(
	c redis.UniversalClient,
) []redis.UniversalClient {
	return []redis.UniversalClient{
		c,
	}
}
