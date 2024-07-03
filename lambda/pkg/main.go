package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
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
	wifeNotificationApi := apis.NewNotificationApi(os.Getenv(WIFE_RECIPIENT))
	notificationApis := []apis.NotificationApi{notificationApi, wifeNotificationApi}
	aiService := service.NewAiService(notificationApis, gptApi)

	gptRes, err := aiService.GenerateMessageToWife()
	if err != nil {
		return err.Error(), err
	}

	return output + gptRes, nil
}
