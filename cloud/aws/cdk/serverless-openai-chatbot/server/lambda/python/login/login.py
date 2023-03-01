"""
人生苦短

help doc:
https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/dynamodb.html
"""

import logging
import boto3
import json
import jwt
import os
import base64
from datetime import datetime, timedelta


logger = logging.getLogger()
logger.setLevel(logging.INFO)

dynamodbClient = boto3.client("dynamodb", region_name="us-east-1")
TABLE_NAME = "chat_user_info"


def create_token(user_name):
    return jwt.encode(payload={
        "aud": "chatbot-user",
        "sub": "chatbot-openai",
        "username": user_name,
        "exp": datetime.utcnow() + timedelta(days=1)
    }, key=os.getenv("TOKEN_KEY"))


def format_response(code, err_msg, token):
    return {
        "isAuthorized": code == 200,
        "body": {
            "message": err_msg,
            "token": token,
        }
    }


def query_dynamodb(key):
    try:
        res = dynamodbClient.get_item(
            TableName=os.getenv("USER_TABLE_NAME"),
            Key={'username': {'S': key}}
        )
        logger.info("res", res)
        if res.Item:
            return res.Item.password.S
    except Exception as e:
        logger.error("except %s", e)
        return None

    return None


def handler(event, context):
    logger.info("Request: %s", event)
    bodyStr = base64.b64decode(event.get('body')).decode('utf-8')
    body = json.loads(bodyStr)
    logger.info("body: %s", body)
    pwd = query_dynamodb(body.username)
    if pwd is None:
        return format_response(403, "user not found", "")

    if pwd != body.password:
        return format_response(403, "Invalid password, please check it", "")

    token = create_token(body.username)
    logger.info("create token: %s", token)
    return format_response(200, "success", token)
