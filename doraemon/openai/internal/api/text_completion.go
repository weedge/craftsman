package api

import (
	"context"
	"errors"
	"fmt"
	"io"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	MaxTokens = 1000
)

var client gogpt.Client

func InitClient(key string) {
	client = *gogpt.NewClient(key)
}

func GetTextCompletion(prompt string) {
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Ada,
		MaxTokens: MaxTokens,
		Prompt:    prompt,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return
	}
	fmt.Println(resp.Choices[0].Text)
}

func GetTextCompletionStream(ctx context.Context, prompt string) error {
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Ada,
		MaxTokens: MaxTokens,
		Prompt:    prompt,
		Stream:    true,
	}
	fmt.Printf("gogpt.CompletionRequest: %v\n", req)
	stream, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished", err.Error())
			return err
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}

		fmt.Printf("Stream response: %v\n", response)
	}
}
