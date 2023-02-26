/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/weedge/craftsman/doraemon/openai/internal/api"
)

// textCompletionCmd represents the textCompletion command
var textCompletionCmd = &cobra.Command{
	Use:   "textCompletion",
	Short: "text completion",
	Long:  `see more: https://platform.openai.com/docs/guides/completion/text-completion`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("textCompletion called")
		api.InitClient(os.Getenv("OPENAI_API_SK"))
		api.GetTextCompletionStream(context.TODO(), cmd.Flag("prompt").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(textCompletionCmd)
	textCompletionCmd.Flags().StringP("prompt", "p", "hi openai", "prompt mode")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// textCompletionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCompletionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
