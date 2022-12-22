package da

import (
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
)

// Options payment da opts
type Options struct {
	Server        *ServerOptions                  `mapstructure:"server"`
	MysqlDBClient *injectors.MysqlDBClientOptions `mapstructure:"mysqlDBClient"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:        DefaultServerOptions(),
		MysqlDBClient: injectors.DefaultMysqlDBClientOptions(),
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
