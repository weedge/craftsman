// Code generated by Kitex v0.4.3. DO NOT EDIT.

package paymentservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	station "github.com/weedge/craftsman/cloudwego/common/kitex_gen/payment/station"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	ChangeAsset(ctx context.Context, req *station.BizAssetChangesReq, callOptions ...callopt.Option) (r *station.BizAssetChangesResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kPaymentServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kPaymentServiceClient struct {
	*kClient
}

func (p *kPaymentServiceClient) ChangeAsset(ctx context.Context, req *station.BizAssetChangesReq, callOptions ...callopt.Option) (r *station.BizAssetChangesResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ChangeAsset(ctx, req)
}
