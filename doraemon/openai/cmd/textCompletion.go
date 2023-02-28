/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/weedge/craftsman/doraemon/openai/internal/api"
)

// textCompletionCmd represents the textCompletion command
var textCompletionCmd = &cobra.Command{
	Use:   "textCompletion",
	Short: "text completion",
	Long:  `see more: https://platform.openai.com/docs/guides/completion/text-completion`,
	Run: func(cmd *cobra.Command, args []string) {
		model := cmd.Flag("model").Value.String()
		stream := cmd.Flag("stream").Value.String()
		chat := cmd.Flag("chat").Value.String()
		if os.Getenv("OPENAI_API_SK") == "" {
			println("please export OPENAI_API_SK")
			return
		}
		api.InitClient(os.Getenv("OPENAI_API_SK"), model)
		if chat == "open" {
			CmdChat(context.TODO(), stream, model)
			return
		}

		if stream == "true" {
			api.GetTextCompletionStream(context.TODO(), cmd.Flag("prompt").Value.String(), model)
		} else {
			api.GetTextCompletion(context.TODO(), cmd.Flag("prompt").Value.String(), model)
		}
	},
}

func CmdChat(ctx context.Context, stream, model string) {
	println("welcome! if u want exit chat, please input 'quit' ")
loop:
	for {
		q := ReadAskQuestionPrompt()
		switch q {
		case "quit":
			break loop
		}
		fmt.Print("A: ")
		if stream == "true" {
			api.GetTextCompletionStream(ctx, q, model)
		} else {
			api.GetTextCompletion(ctx, q, model)
		}
		println()
	}
}

func ReadAskQuestionPrompt() string {
	fmt.Print("Q: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return "quit"
}

func init() {
	rootCmd.AddCommand(textCompletionCmd)
	textCompletionCmd.Flags().StringP("prompt", "p", "hi openai", "prompt mode")
	textCompletionCmd.Flags().StringP("stream", "S", "true", "stream response")
	textCompletionCmd.Flags().StringP("model", "M", gogpt.GPT3TextDavinci003, "openai text models: "+api.StringModels())
	textCompletionCmd.Flags().StringP("chat", "c", "open", "open cmd mode interactive")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// textCompletionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCompletionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
