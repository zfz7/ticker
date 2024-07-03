package apis

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

const CHAT_GPT_API_KEY = "CHAT_GPT_API_KEY"

type ChatGptApi interface {
	Call(request string) (string, error)
}

type chatGptApi struct {
	chatGptClient *openai.Client
}

func NewChatGptApi() ChatGptApi {
	client := openai.NewClient(os.Getenv(CHAT_GPT_API_KEY))

	return &chatGptApi{
		chatGptClient: client,
	}
}

func (api chatGptApi) Call(request string) (string, error) {
	resp, err := api.chatGptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: request,
				},
			},
		},
	)

	if err != nil {
		return fmt.Sprintf("ChatCompletion error: %v\n", err), err
	}
	return resp.Choices[0].Message.Content, nil
}
