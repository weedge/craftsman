use jsonwebtoken::{decode, Algorithm, DecodingKey, Validation};
use lambda_http::{run, service_fn, Body, Error, Request, RequestExt, Response};
use serde::{Deserialize, Serialize};
use std::env;
use tracing::info;

fn get_env_token_key() -> String {
    let u = match env::var_os("TOKEN_KEY") {
        Some(v) => {
            let u = v.into_string().unwrap();
            info!("Got TOKEN_KEY: {}", u);
            u
        }
        None => {
            info!("TOKEN_KEY is not set");
            format!("")
        }
    };
    u
}
// bind encoder/decoder to payload
#[derive(Debug, Serialize, Deserialize)]
struct Payload {
    aud: String,
    sub: String,
    exp: usize,
    username: String,
}

/// This is the main body for the function.
/// Write your code inside it.
/// There are some code example in the following URLs:
/// - https://github.com/awslabs/aws-lambda-rust-runtime/tree/main/examples
async fn function_handler(event: Request) -> Result<Response<Body>, Error> {
    // Extract some useful information from the request
    Ok(match event.query_string_parameters().first("token") {
        Some(token) => {
            info!("request token: {:?}", token);
            let token_key = get_env_token_key();
            let mut status = 200;
            let mut err_msg = "success";
            if token_key == "" {
                status = 400;
                err_msg = "env arg TOKEN_KEY is empty";
            } else {
                let mut validation = Validation::new(Algorithm::HS256);
                validation.sub = Some("chatbot-openai".to_string());
                validation.set_audience(&["chatbot-user"]);
                let token_data = match decode::<Payload>(
                    &token,
                    &DecodingKey::from_secret(token_key.as_bytes()),
                    &validation,
                ) {
                    Ok(c) => {
                        info!("success decode token_data: {:?}", c);
                        c
                    }
                    Err(err) => {
                        info!("Error decoding {:?}", err);
                        return Err(Box::new(err));
                    }
                };
                info!(
                    "token_data.claims {:?} token_data.header {:?} ",
                    token_data.claims, token_data.header
                );
            }
            println!("{} {}", status, err_msg);
            Response::builder()
                .status(status)
                .body(err_msg.into())
                .expect("failed to render response")
        }
        _ => Response::builder()
            .status(400)
            .body("Empty Token".into())
            .expect("failed to render response"),
    })
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        // disable printing the name of the module in every log line.
        .with_target(false)
        // disabling time is handy because CloudWatch will add the ingestion time.
        .without_time()
        .init();

    run(service_fn(function_handler)).await
}
