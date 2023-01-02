package gw

import (
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
)

type Options struct {
	Server                  *ServerOptions                         `mapstructure:"server"`
	HttpThriftGenericClient *injectors.HttpThriftGenericClientOpts `mapstructure:"httpThriftGenericClient"`
}

// DefaultOptions default opts
func DefaultOptions() *Options {
	return &Options{
		Server:                  DefaultServerOptions(),
		HttpThriftGenericClient: injectors.DefaultHttpThriftGenericClientOpts(),
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
