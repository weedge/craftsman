## System architecture by using AWS infrastructure
![architecture](https://raw.githubusercontent.com/weedge/mypic/master/doraemon/aws-serverless-openai-chatbot.drawio.png)

### Welcome to your CDK Go project!

This is a blank project for Go development with CDK.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

### Useful commands

 * `cdk list`        list stack, u can check build before deployment
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

### Update aws-cdk
This CDK CLI is not compatible with the CDK library used by your application. Please upgrade the CLI to the latest version.
please do this 
```shell
npm uninstall -g aws-cdk
npm install -g aws-cdk
```

### Configure
config [cdk.context.json](./cdk.context.json) like this:
```json
{
  "login_dynamodb_table": "chatbot_user_info",
  "jwt_secret": "123",
  "stage": "dev",//dev,test,pre,production
  "openai_api_key": "",
  "sns_chat_openai_topic": "chatbot_openai_msg"
}
```
### Deploy
```shell
cdk list
# http-api-gateway-login
# websocket-api-gateway-chat
# websocket-api-gateway-connect
# async-ai-chat-push-ws-gw
# websocket-api-gateway

# choose one to deploy or cdk deploy all
cdk deploy --all
```

### Test
add test data use by aws cli, change {default} `~/.aws/credentials` , {login_dynamodb_table} config in [cdk.context.json](./cdk.context.json)
```shell
# https://docs.aws.amazon.com/cli/latest/reference/dynamodb/index.html#cli-aws-dynamodb
# https://docs.aws.amazon.com/cli/latest/reference/dynamodb/create-table.html
# aws --profile {default} dynamodb create-table \
#     --table-name {login_dynamodb_table} \
#     --attribute-definitions \
#         AttributeName=username,AttributeType=S \
#     --key-schema \
#         AttributeName=username,KeyType=HASH \
#     --provisioned-throughput \
#         ReadCapacityUnits=5,WriteCapacityUnits=5 \
#     --table-class STANDARD

# https://docs.aws.amazon.com/cli/latest/reference/dynamodb/put-item.html
# --transact-items
aws --profile {default} dynamodb put-item \
    --table-name {login_dynamodb_table} \
    --item '{"user_name":{"S":"root"},"password":{"S":"12345678"}}'

# https://docs.aws.amazon.com/cli/latest/reference/dynamodb/get-item.html
aws --profile {default} dynamodb get-item --table-name {login_dynamodb_table}  --key '{"user_name": {"S": "root"}}'

# change you deploy {region} and api gw id
curl -XPOST "https://{id}.execute-api.{region}.amazonaws.com/dev/login" -d '{"username":"root","password":"12345678"}'
```

### Notice
this dynamodb table is just a demo, if want have unique key, u can see this: 
[simulating-amazon-dynamodb-unique-constraints-using-transactions](https://aws.amazon.com/fr/blogs/database/simulating-amazon-dynamodb-unique-constraints-using-transactions/) 
[general-nosql-design-use-dynamodb](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-general-nosql-design.html)

u need change database model for your biz application

### Project specifics
1. Lambda function name: {OrgName}-{projectName}-{modelName/actionName}-{stageName}; eg: niubi-chatbot-login-dev
2. Database name: {OrgName}_{projectName}_{modelName/actionName}_{stageName}; eg: niubi_chatbot_user_info_dev
3. SNS TOPIC name: {OrgName}_{projectName}_{modelName/actionName}-{stageName}; eg: niubi_chatbot_msg-dev

deploy env, for production, please use different admin/root role to deploy to the user region

happy coding~ :)

### Reference:
1. https://github.com/aws-samples/aws-serverless-openai-chatbot-demo 




