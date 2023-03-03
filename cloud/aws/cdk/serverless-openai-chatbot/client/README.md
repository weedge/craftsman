### FE
from: https://github.com/aws-samples/aws-serverless-openai-chatbot-demo/tree/main/v2-websocket/client

### Change https/wss GW API
 change [src/commons/apigw.js](src/commons/apigw.js) for dev/test/pre/production. 
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
