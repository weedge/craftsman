package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WsGwConnectStackProps struct {
	awscdk.StackProps
}

func NewWsGwConnectStack(scope constructs.Construct, id string, props *WsGwConnectStackProps) (constructs.Construct, awslambda.Function) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	jwt_secret := stack.Node().TryGetContext(jsii.String("jwt_secret")).(string)
	stage := stack.Node().TryGetContext(jsii.String("stage")).(string)

	connectHandler := awslambda.NewFunction(stack, jsii.String("connectHandler"), &awslambda.FunctionProps{
		Code: awslambda.Code_FromAsset(jsii.String("server/lambda/rust/connect/target/lambda/connect"), nil),
		// The runtime environment for the Lambda function that you are uploading.
		// For valid values, see the Runtime property in the AWS Lambda Developer Guide.
		// Use Runtime.FROM_IMAGE when defining a function from a Docker image.
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		//handler: — This name does not matter in case of custom runtime for rust lambda functions
		Handler: jsii.String("does_not_matter"),
		Environment: &map[string]*string{
			"TOKEN_KEY": jsii.String(jwt_secret),
		},
		FunctionName: jsii.String("chatbot-connect-" + stage),
		Description:  jsii.String("On connect event, check jwt authorization"),
	})

	if _, ok := StageAutoDeploy[stage]; !ok {
		return stack, connectHandler
	}

	fnUrl := connectHandler.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("connectHandlerUrl"), &awscdk.CfnOutputProps{
		Value: fnUrl.Url(),
	})

	return stack, connectHandler
}
