package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	awscdklambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WsGwPushStackProps struct {
	awscdk.StackProps
	Topic awssns.Topic
}

func NewWsGwPushStack(scope constructs.Construct, id string, props *WsGwPushStackProps) (constructs.Construct, awslambda.Function) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	openai_api_key := stack.Node().TryGetContext(jsii.String("openai_api_key")).(string)
	stage := stack.Node().TryGetContext(jsii.String("stage")).(string)

	pushHandler := awscdklambdago.NewGoFunction(stack, jsii.String("pushHandler"), &awscdklambdago.GoFunctionProps{
		FunctionName: jsii.String("chatbot-push-" + stage),
		Description:  jsii.String("sub SNS topic get prompt, request openai API resp stream, then send resp text by connectId"),
		Entry:        jsii.String("server/lambda/golang/push"),
		Environment: &map[string]*string{
			"OPENAI_API_KEY": jsii.String(openai_api_key),
		},
		Timeout: awscdk.Duration_Minutes(jsii.Number(5)),
	})

	props.Topic.AddSubscription(awssnssubscriptions.NewLambdaSubscription(pushHandler, &awssnssubscriptions.LambdaSubscriptionProps{}))

	pushHandler.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonAPIGatewayInvokeFullAccess")))

	fnUrl := pushHandler.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String("pushHandlerUrl"), &awscdk.CfnOutputProps{
		Value: fnUrl.Url(),
	})

	return stack, pushHandler
}
