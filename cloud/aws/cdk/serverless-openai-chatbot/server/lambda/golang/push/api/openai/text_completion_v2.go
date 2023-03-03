package openai

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func GetChatCompletion(ctx context.Context, client *gogpt.Client, req gogpt.ChatCompletionRequest) (res gogpt.ChatCompletionMessage, err error) {
	if req.Stream {
		return
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Println("CreateCompletion err", err.Error())
		return
	}

	return resp.Choices[0].Message, nil
}

// https://platform.openai.com/docs/api-reference/chat/create
// stream events: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format
// u should use channel hash connectId to concurrently send msg to websocket
func GetChatCompletionStream(ctx context.Context, client *gogpt.Client, req gogpt.ChatCompletionRequest) (gogpt.ChatCompletionMessage, error) {
	if !req.Stream {
		return gogpt.ChatCompletionMessage{}, nil
	}

	log.Printf("CreateChatCompletionStream req %+v \n", req)
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("CreateChatCompletionStream req %+v err:%s\n", req, err.Error())
		return gogpt.ChatCompletionMessage{}, err
	}
	defer stream.Close()

	index := 0
	res := strings.Builder{}
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("Stream error:%s\n", err.Error())
			return gogpt.ChatCompletionMessage{}, err
		}

		if index > 1 {
			word := response.Choices[0].Delta.Content
			res.WriteString(word)
		}
		index++
	}
	return gogpt.ChatCompletionMessage{Content: res.String()}, err
}
