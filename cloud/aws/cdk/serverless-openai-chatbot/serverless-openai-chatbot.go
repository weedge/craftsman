package main

import (
	"serverless-openai-chatbot/infra"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	infra.NewHttpLoginApiStack(app, "http-api-gateway-login", &infra.HttpLoginApiStackProps{
		StackProps: awscdk.StackProps{
			Env:         env(),
			StackName:   jsii.String("HttpLoginApi"),
			Description: jsii.String("http api gateway /login,"),
		},
	})

	wsConnectConstruct, connectHandler := infra.NewWsGwConnectStack(app, "websocket-api-gateway-connect", &infra.WsGwConnectStackProps{
		StackProps: awscdk.StackProps{
			Env:         env(),
			StackName:   jsii.String("WsApiGwConnectStack"),
			Description: jsii.String("websocket api gateway /$connect"),
		},
	})
	_, _ = wsConnectConstruct, connectHandler

	wsChatConstruct, chatMsgTopic, chatHandler := infra.NewWsGwChatStack(app, "websocket-api-gateway-chat", &infra.WsGwChatStackProps{
		StackProps: awscdk.StackProps{
			Env:         env(),
			StackName:   jsii.String("WsApiGwChatStack"),
			Description: jsii.String("websocket api gateway /$connect,"),
		},
	})
	_, _, _ = wsChatConstruct, chatMsgTopic, chatHandler

	wsPushConstruct, pushHandler := infra.NewWsGwPushStack(app, "async-ai-chat-push-ws-gw", &infra.WsGwPushStackProps{
		StackProps: awscdk.StackProps{
			Env:         env(),
			StackName:   jsii.String("PushAIContent2WsStack"),
			Description: jsii.String("Push AI content to websocket api gateway"),
		},
		Topic: chatMsgTopic,
	})
	_, _ = wsPushConstruct, pushHandler

	infra.NewWsGwStack(app, "websocket-api-gateway", &infra.WsGwStackProps{
		StackProps:     awscdk.StackProps{},
		ConnectHandler: connectHandler,
		ChatHandler:    chatHandler,
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
