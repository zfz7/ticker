package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gregdel/pushover"
	"lambda/pkg/apis"
	"math/rand"
	"strings"
)

type AiService interface {
	GenerateMessageToWife() (string, error)
}

type aiService struct {
	notificationApis []apis.NotificationApi
	chatGptApi       apis.ChatGptApi
}

func NewAiService(notificationApis []apis.NotificationApi, chatGptApi apis.ChatGptApi) AiService {
	return &aiService{
		notificationApis: notificationApis,
		chatGptApi:       chatGptApi,
	}
}

var InitialPrompt = "I am trying to send my wife a nice message to boost her mood, it can be funny or sweet. " +
	"Can you generate a 15 prompts that would result in a that message? Please return the prompts in valid " +
	"json using the follow format. {\"prompts\":[\"prompt1\", \"prompt2\"]}"
var FinalPrompt = " Please call her Babe, and return as the following json: {\"message\":\"message\"}"

func (aiService aiService) GenerateMessageToWife() (string, error) {
	promptRes, err := aiService.chatGptApi.Call(InitialPrompt)
	if err != nil {
		return "Failed first prompt", err
	}

	prompts, err := parseJson[Prompts](promptRes)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %v\n, err %v", promptRes, err), err
	}

	i := rand.Intn(len(prompts.Prompts))

	message, err := aiService.chatGptApi.Call(prompts.Prompts[i] + FinalPrompt)
	if err != nil {
		return "Failed second prompt: " + prompts.Prompts[i], nil
	}

	m, err := parseJson[Message](message)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %v\n, err %v", promptRes, err), err
	}

	messageRequest := pushover.Message{
		Message:  m.Message,
		Title:    "I lovee you",
		Priority: pushover.PriorityNormal,
	}

	for _, notification := range aiService.notificationApis {
		_, err := notification.Send(&messageRequest)
		if err != nil {
			return "Did not send message: " + m.Message, err
		}
	}
	return m.Message, nil
}

func parseJson[V Prompts | Message](s string) (V, error) {
	var res V
	startJSON := strings.Index(s, "{")
	endJSON := strings.LastIndex(s, "}") + 1

	if startJSON < 0 || endJSON < 1 {
		return res, errors.New("no json found")
	}
	jsonData := s[startJSON:endJSON]
	jsonData = strings.TrimSpace(jsonData)

	err := json.Unmarshal([]byte(jsonData), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

type Prompts struct {
	Prompts []string `json:"prompts"`
}

type Message struct {
	Message string `json:"message"`
}
