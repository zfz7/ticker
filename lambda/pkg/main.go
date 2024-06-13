package main

import (
	"context"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gregdel/pushover"
	"os"
	"strings"
	"time"
)

const FINNHUB_API_KEY = "FINNHUB_API_KEY"
const PUSHOVER_APP_KEY = "PUSHOVER_APP_KEY"
const PUSHOVER_RECIPIENT = "PUSHOVER_RECIPIENT"

func main() {
	lambda.Start(router)
}

func router(ctx context.Context, req map[string]string) (string, error) {
	finnHubCfg := finnhub.NewConfiguration()
	finnHubCfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv(FINNHUB_API_KEY))
	finnhubClient := finnhub.NewAPIClient(finnHubCfg).DefaultApi

	pushoverClient := pushover.New(os.Getenv(PUSHOVER_APP_KEY))

	output, err := runTicker(finnhubClient, pushoverClient)
	if err != nil {
		return err.Error(), err
	}

	return output, nil
}

func runTicker(finnhubClient *finnhub.DefaultApiService, pushoverClient *pushover.Pushover) (string, error) {
	tickers := [3]string{"SPY", "QQQ", "AMZN"}

	var title strings.Builder
	var message strings.Builder

	for idx, ticker := range tickers {
		quote, _, err := finnhubClient.Quote(context.Background()).Symbol(ticker).Execute()
		if err != nil {
			return "", err
		}
		if idx == 0 && *quote.Dp > 0 {
			fmt.Fprintf(&title, "Market: UP")
		} else if idx == 0 && *quote.Dp < 0 {
			fmt.Fprintf(&title, "Market: DOWN")
		}
		fmt.Fprintf(&message, "%v: %v$, %%\u0394: %2.2f, $\u0394: %v\n", ticker, *quote.C, *quote.Dp, *quote.D)
	}

	easternTime, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(easternTime)
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 30, 0, 0, easternTime)

	duration := marketClose.Sub(now).Truncate(time.Minute)
	fmt.Fprintf(&title, ", cls in %v", duration)

	recipient := pushover.NewRecipient(os.Getenv(PUSHOVER_RECIPIENT))
	messageRequest := pushover.NewMessageWithTitle(message.String(), title.String())
	_, err := pushoverClient.SendMessage(messageRequest, recipient)
	if err != nil {
		return "", err
	}
	return title.String() + "|" + message.String(), nil
}
