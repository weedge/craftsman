### Init FE client Code
from: https://github.com/aws-samples/aws-serverless-openai-chatbot-demo/tree/main/v2-websocket/client
```shell
git clone https://github.com/weedge/craftmaster.git
cd craftsman/cloud/aws/cdk/serverless-openai-chatbot/client
git init
git remote add origin https://github.com/aws-samples/aws-serverless-openai-chatbot-demo.git
git config core.sparsecheckout true
echo "v2-websocket/client/*" >> .git/info/sparse-checkout
git pull origin main
```

### Change https/wss GW API
 change [v2-websocket/client/src/commons/apigw.js](src/commons/apigw.js) for dev/test/pre/production. 
 ```js
// Change to your own WebSocket API Gateway endpoint
export const API_socket = 'wss://{apiid}.execute-api.{region}.amazonaws.com/{stage}';
// Change to your own HTTP API Gateway endpoint
export const API_http = 'https://{apiid}.execute-api.{region}.amazonaws.com/{stage}';
 ```

### Build
```shell
npm install
npm run build

# local to review
npm install -g serve
serve -s build
```

### Sync to S3
```shell
# use aws cli to upload all files in build folder to the bucket
aws s3 ./build/ s3://{bucket-name}/
```
