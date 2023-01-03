package da

import (
	"fmt"

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
		RmqConsumers:  map[string]*subscriber.RmqPushConsumerOptions{},
	}
}

func (m *Options) String() (str string) {
	str += fmt.Sprintf(" Server:%+v ", m.Server)
	str += fmt.Sprintf(" MysqlDBClient:%+v ", m.MysqlDBClient)
	str += "RmqConsumer:"
	for name, item := range m.RmqConsumers {
		str += fmt.Sprintf("%s:%+v ", name, item)
	}
	return
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

	klog.Infof("config:%v", opt)

	return opt, nil
}
