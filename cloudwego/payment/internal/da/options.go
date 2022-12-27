package da

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/subscriber"
)

// Options payment da opts
type Options struct {
	Server        *ServerOptions                                `mapstructure:"server"`
	MysqlDBClient *injectors.MysqlDBClientOptions               `mapstructure:"mysqlDBClient"`
	RmqConsumers  map[string]*subscriber.RmqPushConsumerOptions `mapstructure:"rmqConsumers"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:        DefaultServerOptions(),
		MysqlDBClient: injectors.DefaultMysqlDBClientOptions(),
		//RmqConsumers:  []*subscriber.RmqPushConsumerOptions{},
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

	klog.Infof("server: %+v, mysqlDBClient: %+v", opt.Server, opt.MysqlDBClient)

	return opt, nil
}
