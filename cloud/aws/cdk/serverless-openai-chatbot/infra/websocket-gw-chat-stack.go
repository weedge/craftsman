package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WsGwChatStackProps struct {
	awscdk.StackProps
}

func NewWsGwChatStack(scope constructs.Construct, id string, props *WsGwChatStackProps) (constructs.Construct, awssns.Topic, awslambda.Function) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	sendPromptNoticationTopic := awssns.NewTopic(stack, jsii.String("sendPromptNotication"), &awssns.TopicProps{
		DisplayName: jsii.String("sendPromptAlertNotication"),
	})

	chatHandler := awslambda.NewFunction(stack, jsii.String("chatHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("server/lambda/nodejs/chat"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("index.handler"),
		Environment: &map[string]*string{
			"SNS_TOPIC_ARN": sendPromptNoticationTopic.TopicArn(),
		},
	})

	sendPromptNoticationTopic.GrantPublish(chatHandler)

	return stack, sendPromptNoticationTopic, chatHandler
}
