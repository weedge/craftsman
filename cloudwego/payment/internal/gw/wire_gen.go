// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package gw

import (
	"context"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
)

// Injectors from wire.go:

// NewServer build server with wire, dependency obj inject, so init random
func NewServer(ctx context.Context) (*Server, error) {
	provider := configparser.Default()
	options, err := Configure(provider)
	if err != nil {
		return nil, err
	}
	serverOptions := options.Server
	level := serverOptions.LogLevel
	v := serverOptions.LogMeta
	iKitexZapKVLogger := logutils.NewkitexZapKVLogger(level, v)
	httpThriftGenericClientOpts := options.HttpThriftGenericClient
	v2 := injectors.InitHttpThriftGenericClients(httpThriftGenericClientOpts)
	v3 := injectors.InitHttpThriftGenericEndpointsOpts(httpThriftGenericClientOpts)
	server := &Server{
		opts:                 serverOptions,
		kitexKVLogger:        iKitexZapKVLogger,
		mapCli:               v2,
		mapGenericClientOpts: v3,
	}
	return server, nil
}
