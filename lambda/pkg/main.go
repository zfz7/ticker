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
	marketDown := true
	for idx, ticker := range tickers {
		quote, _, err := finnhubClient.Quote(context.Background()).Symbol(ticker).Execute()
		if err != nil {
			return "", err
		}
		if idx == 0 && *quote.Dp > 0 {
			marketDown = false
			fmt.Fprintf(&title, "Market: ▲,")
		} else if idx == 0 && *quote.Dp < 0 {
			marketDown = true
			fmt.Fprintf(&title, "Market: ▼,")
		}
		fmt.Fprintf(&message, "%v: %v$, %+2.2f%%, %+.2f$\n", ticker, *quote.C, *quote.Dp, *quote.D)
	}
	easternTime, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(easternTime)
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 30, 0, 0, easternTime)
	duration := marketClose.Sub(now).Truncate(time.Minute)

	marketStatus, _, _ := finnhubClient.MarketStatus(context.Background()).Exchange("US").Execute()
	if *marketStatus.IsOpen {
		fmt.Fprintf(&title, " open for %v", duration)
	} else {
		fmt.Fprintf(&title, " closed")
	}

	recipient := pushover.NewRecipient(os.Getenv(PUSHOVER_RECIPIENT))
	messageRequest := &pushover.Message{
		Message:  message.String(),
		Title:    title.String(),
		Priority: pushover.PriorityNormal,
	}
	if marketDown && *marketStatus.IsOpen && duration <= time.Hour {
		messageRequest.Priority = pushover.PriorityEmergency
		messageRequest.Retry = 5 * time.Minute
		messageRequest.Expire = time.Hour
		messageRequest.Sound = pushover.SoundMechanical
	}
	_, err := pushoverClient.SendMessage(messageRequest, recipient)
	if err != nil {
		return "", err
	}
	return title.String() + "|" + message.String(), nil
}
