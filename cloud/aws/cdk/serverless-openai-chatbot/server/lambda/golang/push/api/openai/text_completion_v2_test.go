package openai

import (
	"context"
	"os"
	"testing"

	gogpt "github.com/sashabaranov/go-gpt3"
)

var gogptClient *gogpt.Client

func TestMain(m *testing.M) {
	gogptClient = gogpt.NewClient(os.Getenv("OPENAI_API_KEY"))
	m.Run()
}

func TestGetChatCompletionStream(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *gogpt.Client
		req    gogpt.ChatCompletionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test context ok",
			args: args{
				ctx:    context.TODO(),
				client: gogptClient,
				req: gogpt.ChatCompletionRequest{
					Model: gogpt.GPT3Dot5Turbo,
					Messages: []gogpt.ChatCompletionMessage{
						{
							Role:    "user",
							Content: "hello",
						},
					},
					//MaxTokens:        4096,
					//Temperature:      0,
					//TopP:             1,
					//N:                0,
					Stream: true,
					//Stop:             []string{},
					//PresencePenalty:  0,
					//FrequencyPenalty: 0,
					//LogitBias:        map[string]int{},
					//User:             "",
				},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println("adfadf")
			got, err := GetChatCompletionStream(tt.args.ctx, tt.args.client, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChatCompletionStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Content == tt.want {
				t.Errorf("GetChatCompletionStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
