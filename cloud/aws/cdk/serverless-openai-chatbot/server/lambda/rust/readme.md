### Develop (cargo new)
```shell
cd craftsman/cloud/aws/cdk/serverless-openai-chatbot/server/lambda/rust
# change {my_lambda_function} to your handler
cargo new {my_lambda_function} --bin
```
### Use Aws Rust Runtime On MacOs (cargo lambda)
```shell
# see more: https://github.com/awslabs/aws-lambda-rust-runtime
brew tap cargo-lambda/cargo-lambda
brew install cargo-lambda
# eg: connect use websocket gw api
cargo lambda new YOUR_FUNCTION_NAME
cd YOUR_FUNCTION_NAME
cargo add serde
cargo add jsonwebtoken
cargo lambda watch
```
tips: u should use this `cargo lambda`

### Install
```shell
git clone github.com/weedge/craftsman 
cd craftsman/cloud/aws/cdk/serverless-openai-chatbot/server/lambda/rust/connect
```


### Deploy
```shell
cargo lambda build --release --target x86_64-unknown-linux-musl
cargo lambda deploy \
  --iam-role arn:aws:iam::XXXXXXXXXXXXX:role/your_lambda_execution_role \
  YOUR_FUNCTION_NAME
# aws-cli deploy see https://github.com/awslabs/aws-lambda-rust-runtime#22-deploying-with-the-aws-cli
# asm deploy see https://github.com/awslabs/aws-lambda-rust-runtime#23-aws-serverless-application-model-sam
# serverless deploy see https://github.com/awslabs/aws-lambda-rust-runtime#24-serverless-framework
# docker deploy use lambci/lambda:provided container https://github.com/awslabs/aws-lambda-rust-runtime#24-serverless-framework
cargo lambda build --release --target x86_64-unknown-linux-musl --output-format zip
# or u can use CDK deploy release bootstrap zip like aws-cli 
```

### Reference
1. https://course.rs/usecases/aws-rust.html
2. https://aws.amazon.com/cn/blogs/opensource/sustainability-with-rust/ 
3. https://aws.amazon.com/cn/developer/language/rust/
4. https://aws.amazon.com/cn/blogs/china/rust-runtime-for-aws-lambda/
5. https://github.com/awslabs/aws-lambda-rust-runtime 
6. https://blog.logrocket.com/deploy-lambda-functions-rust/
7. https://medium.com/techhappily/rust-based-aws-lambda-with-aws-cdk-deployment-14a9a8652d62


**<u>rust so cool~ learn by good example, axiba, rust lifetime<u>**