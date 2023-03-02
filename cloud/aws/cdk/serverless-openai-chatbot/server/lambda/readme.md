### intro
this dir is aws lambda function, u can use sam manage lambda.

```shell
sam init
cd {sam-app}
sam build
sam validate

# test by event file
sam local invoke "HelloWorldFunction" -e events/event.json
# or
sam local start-api



# if local test ok, deploy 
sam deploy --guided
```

### reference
1. sam:  https://docs.aws.amazon.com/zh_cn/serverless-application-model/latest/developerguide/what-is-sam.html