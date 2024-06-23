package service

import (
	"fmt"
	"github.com/gregdel/pushover"
	"lambda/pkg/apis"
	"strings"
	"time"
)

type TickerService interface {
	Run() (string, error)
}

type tickerService struct {
	notificationApi apis.NotificationApi
	stockApi        apis.StockApi
}

func NewTickerService(stockApi apis.StockApi, notificationAPi apis.NotificationApi) TickerService {
	return &tickerService{
		notificationApi: notificationAPi,
		stockApi:        stockApi,
	}
}

func (tickerService *tickerService) Run() (string, error) {
	tickers := [3]string{"SPY", "QQQ", "AMZN"}

	var title strings.Builder
	var message strings.Builder
	marketDown := true

	for idx, ticker := range tickers {
		quote, err := tickerService.stockApi.Quote(ticker)
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
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 00, 0, 0, easternTime)
	timeTillClose := marketClose.Sub(now).Truncate(time.Minute)

	marketStatus, _ := tickerService.stockApi.MarketStatus("US")
	if *marketStatus.IsOpen {
		fmt.Fprintf(&title, " open for %v", timeTillClose)
	} else {
		fmt.Fprintf(&title, " closed")
	}

	messageRequest := pushover.Message{
		Message:  message.String(),
		Title:    title.String(),
		Priority: pushover.PriorityNormal,
	}
	if marketDown && *marketStatus.IsOpen && timeTillClose <= time.Hour {
		messageRequest.Priority = pushover.PriorityEmergency
		messageRequest.Retry = 5 * time.Minute
		messageRequest.Expire = time.Hour
		messageRequest.Sound = pushover.SoundMechanical
	}
	_, err := tickerService.notificationApi.Send(&messageRequest)
	if err != nil {
		return "", err
	}
	return title.String() + "\n" + message.String(), nil
}
