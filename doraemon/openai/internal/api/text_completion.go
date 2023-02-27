package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	MaxTokens   = 2000
	Temperature = 0.7
)

var client gogpt.Client
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

func StringModels() (modelsStr string) {
	for m := range textModels {
		modelsStr += m + " "
	}
	return
}

func InitClient(key string, model string) {
	if _, ok := textModels[model]; !ok {
		log.Fatalln("un support model", model)
		return
	}
	fmt.Println("use model:", model)
	client = *gogpt.NewClient(key)
}

func GetTextCompletion(ctx context.Context, prompt string, model string) {
	req := gogpt.CompletionRequest{
		Model:       model,
		MaxTokens:   MaxTokens,
		Prompt:      prompt,
		Temperature: Temperature,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(strings.TrimPrefix(resp.Choices[0].Text, "\n\n"))
}

func GetTextCompletionStream(ctx context.Context, prompt string, model string) (string, error) {
	req := gogpt.CompletionRequest{
		Model:       model,
		MaxTokens:   MaxTokens,
		Prompt:      prompt,
		Stream:      true,
		Temperature: Temperature,
	}
	stream, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		fmt.Println("Stream err", err.Error())
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
			fmt.Printf("\nStream error:%s\n", err.Error())
			return "", err
		}

		if index > 1 {
			word := response.Choices[0].Text
			res.WriteString(word)
			fmt.Print(word)
		}
		index++
	}
	fmt.Println()
	return res.String(), nil
}
