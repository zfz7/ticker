package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"os"
	"strings"
)

const CHAT_GPT_API_KEY = "CHAT_GPT_API_KEY"

type ChatGptApi interface {
	Call(request string) (string, error)
	GenerateMessageToWife() (string, error)
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
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}
	//fmt.Printf(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}

func (api chatGptApi) GenerateMessageToWife() (string, error) {
	prompt := "I am trying to send my wife a nice message to boost her mood, it can be funny or sweet. " +
		"Can you generate a 15 prompts that would result in a that message? Please return the prompts in valid " +
		"json using the follow format. {\"prompts\":[\"prompt1\", \"prompt2\"]}"
	promptRes, err := api.Call(prompt)
	if err != nil {
		return "", err
	}
	startJSON := strings.Index(promptRes, "{")
	endJSON := strings.LastIndex(promptRes, "}") + 1
	jsonData := promptRes[startJSON:endJSON]
	jsonData = strings.TrimSpace(jsonData)

	var p Prompts
	err = json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n, err %v", jsonData, err)
		return "", nil
	}

	i := rand.Intn(len(p.Prompts))
	message, err := api.Call(p.Prompts[i] + " Please call her Babe, and return as the following json: {\"message\":\"message\"}")
	if err != nil {
		return "", nil
	}
	startJSON = strings.Index(message, "{")
	endJSON = strings.LastIndex(message, "}") + 1
	jsonData = message[startJSON:endJSON]
	jsonData = strings.TrimSpace(jsonData)

	var m Message
	err = json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n, err %v", jsonData, err)
		return "", nil
	}
	return m.Message, nil
}

type Prompts struct {
	Prompts []string `json:"prompts"`
}

type Message struct {
	Message string `json:"message"`
}
