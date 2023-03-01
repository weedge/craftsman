package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WsGwStackProps struct {
	awscdk.StackProps
	connectHandler awslambda.Function
	chatHandler    awslambda.Function
	pushHandler    awslambda.Function
}

func NewWsGwStack(scope constructs.Construct, id string, props *WsGwStackProps) constructs.Construct {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	wsApi := awscdkapigatewayv2alpha.NewWebSocketApi(stack, jsii.String("ws-gw-chatbot"), &awscdkapigatewayv2alpha.WebSocketApiProps{
		ApiName: jsii.String("web-socket-gateway-chatbot"),
		ConnectRouteOptions: &awscdkapigatewayv2alpha.WebSocketRouteOptions{
			Integration:    awscdkapigatewayv2integrationsalpha.NewWebSocketLambdaIntegration(jsii.String("ws-gw-chatbot-connect"), props.connectHandler),
			Authorizer:     nil,
			ReturnResponse: jsii.Bool(true),
		},
		Description: jsii.String("websocket gateway chatbot"),
	})
	wsApi.AddRoute(jsii.String("sendprompt"), &awscdkapigatewayv2alpha.WebSocketRouteOptions{
		Integration:    awscdkapigatewayv2integrationsalpha.NewWebSocketLambdaIntegration(jsii.String("ws-gw-chatbot-sendprompt"), props.chatHandler),
		ReturnResponse: jsii.Bool(false),
	})

	return stack
}
