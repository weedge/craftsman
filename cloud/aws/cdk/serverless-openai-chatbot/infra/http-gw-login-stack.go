package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	awscdklambdapython "github.com/aws/aws-cdk-go/awscdklambdapythonalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HttpLoginApiStackProps struct {
	awscdk.StackProps
}

var StageAutoDeploy map[string]bool = map[string]bool{
	"dev":        true,
	"test":       true,
	"pre":        false,
	"production": false,
}

func NewHttpLoginApiStack(scope constructs.Construct, id string, props *HttpLoginApiStackProps) constructs.Construct {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	jwt_secret := stack.Node().TryGetContext(jsii.String("jwt_secret")).(string)
	login_dynamodb_table := stack.Node().TryGetContext(jsii.String("login_dynamodb_table")).(string)
	stage := stack.Node().TryGetContext(jsii.String("stage")).(string)

	loginTable := awsdynamodb.NewTable(stack, jsii.String(login_dynamodb_table), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("user_name"), Type: awsdynamodb.AttributeType_STRING},
		//SortKey:       &awsdynamodb.Attribute{Name: jsii.String("created_at"), Type: awsdynamodb.AttributeType_STRING},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Encryption:    awsdynamodb.TableEncryption_AWS_MANAGED,
		TableName:     jsii.String(login_dynamodb_table + "_" + stage),
		//ReadCapacity:  jsii.Number(7),
	})

	loginHandler := awscdklambdapython.NewPythonFunction(stack, jsii.String("loginHandler"), &awscdklambdapython.PythonFunctionProps{
		Entry:   jsii.String("server/lambda/python/login"),
		Index:   jsii.String("login.py"),
		Handler: jsii.String("handler"),
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Environment: &map[string]*string{
			"TOKEN_KEY":       jsii.String(jwt_secret),
			"USER_TABLE_NAME": loginTable.TableName(),
		},
		FunctionName: jsii.String("chatbot-login-" + stage),
		Description:  jsii.String("user login and authentication user information from dynamodb, ok return jwt for authorization"),
	})

	/*
		loginHandler := awslambda.NewFunction(stack, jsii.String("loginHandler"), &awslambda.FunctionProps{
			Handler: jsii.String("login.handler"),
			Code: awslambda.Code_FromAsset(jsii.String("server/lambda/python/login"), &awss3assets.AssetOptions{
				Bundling: &awscdk.BundlingOptions{
					Image: awslambda.Runtime_PYTHON_3_9().BundlingImage(),
					Command: &[]*string{
						jsii.String("bash"),
						jsii.String("-c"),
						jsii.String("pip install -r requirements.txt -t /asset-output && cp -au . /asset-output"),
					},
				},
			}),
		})
	*/

	loginTable.GrantReadWriteData(loginHandler)

	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("http-gw-chatbot"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("http-gateway-chatbot-" + stage),
		CorsPreflight: &awscdkapigatewayv2alpha.CorsPreflightOptions{
			AllowCredentials: jsii.Bool(false),
			AllowHeaders: &[]*string{
				jsii.String("*"),
			},
			AllowMethods: &[]awscdkapigatewayv2alpha.CorsHttpMethod{
				awscdkapigatewayv2alpha.CorsHttpMethod_ANY,
				//awscdkapigatewayv2alpha.CorsHttpMethod_GET,
				//awscdkapigatewayv2alpha.CorsHttpMethod_POST,
			},
			AllowOrigins: &[]*string{
				jsii.String("*"),
			},
		},
		CreateDefaultStage: jsii.Bool(false),
		Description:        jsii.String("htttp api gateway"),
	})
	httpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(
			jsii.String("login-route"), loginHandler,
			&awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
		Path:    jsii.String("/login"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
	})

	if _, ok := StageAutoDeploy[stage]; !ok {
		return stack
	}
	httpApi.AddStage(jsii.String(stage), &awscdkapigatewayv2alpha.HttpStageOptions{
		AutoDeploy: jsii.Bool(StageAutoDeploy[stage]),
		StageName:  jsii.String(stage),
	})

	return stack
}
