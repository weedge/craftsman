package station

import (
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
)

type Options struct {
	Server                   *ServerOptions                       `mapstructure:"server"`
	RedisClusterClient       *injectors.RedisClusterClientOptions `mapstructure:"redisClusterClient"`
	AssetChangeEventProducer *injectors.RmqProducerClientOptions  `mapstructure:"assetChangeEventProducer"`
	PaymentDaClient          *injectors.PaymentDaClientOptions    `mapstruct:"paymentDaClient"`

	UserAssetTxMethod string `mapstructure:"userAssetTxMethod"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:                   DefaultServerOptions(),
		RedisClusterClient:       injectors.DefaultRedisClusterClientOptions(),
		AssetChangeEventProducer: injectors.DefaultRmqProducerClientOptions(),
		PaymentDaClient:          injectors.DefaultPaymentDaClientOptions(),

		UserAssetTxMethod: constants.UserAssetTxMethodLua,
	}
}

// Configure inject config
func Configure(configProvider configparser.Provider) (*Options, error) {
	opt := DefaultOptions()

	cp, err := configProvider.Get()
	if err != nil {
		return nil, err
	}

	if err = cp.UnmarshalExact(opt); err != nil {
		return nil, err
	}

	return opt, nil
}
