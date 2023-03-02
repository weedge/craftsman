package openai

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	MaxTokens   = 2000
	Temperature = 0.7
)

var textModels = map[string]struct{}{
	gogpt.GPT3TextDavinci003:      {},
	gogpt.GPT3TextDavinci002:      {},
	gogpt.GPT3TextCurie001:        {},
	gogpt.GPT3TextBabbage001:      {},
	gogpt.GPT3TextAda001:          {},
	gogpt.GPT3TextDavinci001:      {},
	gogpt.GPT3DavinciInstructBeta: {},
	gogpt.GPT3Davinci:             {},
	gogpt.GPT3CurieInstructBeta:   {},
	gogpt.GPT3Curie:               {},
	gogpt.GPT3Ada:                 {},
	gogpt.GPT3Babbage:             {},
}

func GetTextCompletion(ctx context.Context, client *gogpt.Client, req gogpt.CompletionRequest) (res string, err error) {
	if req.Stream {
		return
	}

	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		log.Println("CreateCompletion err", err.Error())
		return
	}

	return strings.TrimPrefix(resp.Choices[0].Text, "\n\n"), nil
}

func GetTextCompletionStream(ctx context.Context, client *gogpt.Client, req gogpt.CompletionRequest) (string, error) {
	if !req.Stream {
		return "", nil
	}

	stream, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		log.Println("Stream err", err.Error())
		return "", err
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
			return "", err
		}

		if index > 1 {
			word := response.Choices[0].Text
			res.WriteString(word)
		}
		index++
	}
	return res.String(), nil
}
