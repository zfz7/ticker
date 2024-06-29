package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gregdel/pushover"
	"lambda/pkg/apis"
	"lambda/pkg/service"
	"os"
)

const PUSHOVER_RECIPIENT = "PUSHOVER_RECIPIENT"
const WIFE_RECIPIENT = "WIFE_RECIPIENT"

func main() {
	lambda.Start(router)
}

func router(ctx context.Context, req map[string]string) (string, error) {
	stockApi := apis.NewStockApi()
	notificationApi := apis.NewNotificationApi(os.Getenv(PUSHOVER_RECIPIENT))
	tickerService := service.NewTickerService(stockApi, notificationApi)

	output, err := tickerService.Run()
	if err != nil {
		return err.Error(), err
	}

	gptApi := apis.NewChatGptApi()
	gptRes, err := gptApi.GenerateMessageToWife()
	if err != nil {
		return err.Error(), err
	}
	messageRequest := pushover.Message{
		Message:  gptRes,
		Title:    "I lovee you",
		Priority: pushover.PriorityNormal,
	}
	wifeNotificationApi := apis.NewNotificationApi(os.Getenv(WIFE_RECIPIENT))
	_, err = wifeNotificationApi.Send(&messageRequest)
	if err != nil {
		return err.Error(), err
	}
	_, err = notificationApi.Send(&messageRequest)
	if err != nil {
		return err.Error(), err
	}

	return output + gptRes, nil
}
