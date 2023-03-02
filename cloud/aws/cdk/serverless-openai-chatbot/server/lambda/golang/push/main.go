package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"
	"unsafe"

	"push/api/openai"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	apigw "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var gogptClient *gogpt.Client
var apiWsGwClient *apigw.Client

func Init() {
	log.Println("openai_api_key", os.Getenv("OPENAI_API_KEY"))
	gogptClient = gogpt.NewClient(os.Getenv("OPENAI_API_KEY"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	apiWsGwClient = apigw.NewFromConfig(cfg)
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

type WsEventMsg struct {
	Action  string  `json:"action"`
	Payload Payload `json:"payload"`
}
type Payload struct {
	Msgid  string `json:"msgid"`
	Prompt string `json:"prompt"`
	Params Params `json:"params"`
}
type Params struct {
	FrequencyPenalty float32 `json:"frequency_penalty"`
	MaxTokens        int     `json:"max_tokens"`
	PresencePenalty  float32 `json:"presence_penalty"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"top_p"`
	ModelName        string  `json:"model_name"`
}
type Message struct {
	RequestContext events.APIGatewayWebsocketProxyRequestContext
	Payload        Payload
}

func Handler(ctx context.Context, snsEvent events.SNSEvent) (err error) {
	log.Printf("Event = %s \n", snsEvent)
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		log.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		msg := &Message{}
		err = json.Unmarshal(Str2Bytes(snsRecord.Message), msg)
		if err != nil {
			log.Printf("error:%s\n", err.Error())
			return
		}
		prompt := strings.Trim(msg.Payload.Prompt, " ")
		if prompt == "" {
			log.Println("empty prompt")
			return
		}

		params := msg.Payload.Params
		res := ""
		res, err = openai.GetTextCompletionStream(ctx, gogptClient, gogpt.CompletionRequest{
			Model:            params.ModelName,
			Prompt:           prompt,
			MaxTokens:        params.MaxTokens,
			Temperature:      float32(params.Temperature),
			TopP:             float32(params.TopP),
			Stream:           true,
			PresencePenalty:  float32(params.PresencePenalty),
			FrequencyPenalty: params.FrequencyPenalty,
		})
		if err != nil {
			log.Printf("openai.GetTextCompletionStream error: %v", err)
			res = err.Error()
		}

		reqCtx := msg.RequestContext
		postData, _ := json.Marshal(map[string]any{"msgid": reqCtx.MessageID, "text": res})
		log.Printf("postData: %s\n", postData)
		apiGwUrl := &url.URL{
			Scheme: "https",
			Host:   reqCtx.DomainName,
			Path:   reqCtx.Stage,
		}
		postRes, err1 := apiWsGwClient.PostToConnection(ctx, &apigw.PostToConnectionInput{
			ConnectionId: aws.String(reqCtx.ConnectionID),
			Data:         postData,
		}, func(o *apigw.Options) {
			log.Printf("apiGwUrl:%s\n", apiGwUrl.String())
			o.EndpointResolver = apigw.EndpointResolverFromURL(apiGwUrl.String())
		})
		if err1 != nil {
			log.Printf("post err:%s\n", err1.Error())
			err = err1
			return
		}
		log.Printf("post res:%+v\n", *postRes)
	}

	return
}

func main() {
	Init()
	lambda.Start(Handler)
}
