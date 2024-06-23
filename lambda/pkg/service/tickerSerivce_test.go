package service

import (
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gregdel/pushover"
	"github.com/stretchr/testify/mock"
	"lambda/mocks"
	"strings"
	"testing"
	"time"
)

var spyC, spyDP, spyD float32 = 544.51, -0.13, -0.73
var qqqC, qqqDP, qqqD float32 = 480.18, -0.27, -1.29
var amznC, amznDP, amznD float32 = 189.08, 1.60, 2.98
var mockNotificationApi *mocks.NotificationApi
var mockStockApi *mocks.StockApi
var spyQuoteHandler *mock.Call
var marketStatusHandler *mock.Call

func setup() {
	mockNotificationApi = &mocks.NotificationApi{}
	mockStockApi = &mocks.StockApi{}
	mockNotificationApi.On("Send", mock.Anything).Return(pushover.Response{}, nil)
	var no = false
	marketStatusHandler = mockStockApi.On("MarketStatus", "US").Return(finnhub.MarketStatus{
		IsOpen: &no,
	}, nil)
	spyQuoteHandler = mockStockApi.On("Quote", "SPY").Return(finnhub.Quote{
		C: &spyC, D: &spyD, Dp: &spyDP,
	}, nil)
	mockStockApi.On("Quote", "QQQ").Return(finnhub.Quote{
		C: &qqqC, D: &qqqD, Dp: &qqqDP,
	}, nil)
	mockStockApi.On("Quote", "AMZN").Return(finnhub.Quote{
		C: &amznC, D: &amznD, Dp: &amznDP,
	}, nil)
}

func Test_tickerService_Run_Market_Down_Closed(t *testing.T) {
	setup()
	want := "Market: ▼, closed\n" +
		"SPY: 544.51$, -0.13%, -0.73$\n" +
		"QQQ: 480.18$, -0.27%, -1.29$\n" +
		"AMZN: 189.08$, +1.60%, +2.98$\n"
	wantErr := false

	tickerService := &tickerService{
		notificationApi: mockNotificationApi,
		stockApi:        mockStockApi,
	}
	got, err := tickerService.Run()
	if (err != nil) != wantErr {
		t.Errorf("Run() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != want {
		t.Errorf("Run() got = %v, want %v", got, want)
	}
}

func Test_tickerService_Run_Market_Up_Closed(t *testing.T) {
	setup()
	spyQuoteHandler.Unset()
	var spyC, spyDP, spyD float32 = 544.51, 0.13, 0.73
	spyQuoteHandler.On("Quote", "SPY").Return(finnhub.Quote{
		C: &spyC, D: &spyD, Dp: &spyDP,
	}, nil)

	want := "Market: ▲, closed\n" +
		"SPY: 544.51$, +0.13%, +0.73$\n" +
		"QQQ: 480.18$, -0.27%, -1.29$\n" +
		"AMZN: 189.08$, +1.60%, +2.98$\n"
	wantErr := false

	tickerService := &tickerService{
		notificationApi: mockNotificationApi,
		stockApi:        mockStockApi,
	}
	got, err := tickerService.Run()
	if (err != nil) != wantErr {
		t.Errorf("Run() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != want {
		t.Errorf("Run() got = %v, want %v", got, want)
	}
}

func Test_tickerService_Run_Market_Down_Open(t *testing.T) {
	setup()
	marketStatusHandler.Unset()
	var yes = true
	marketStatusHandler = mockStockApi.On("MarketStatus", "US").Return(finnhub.MarketStatus{
		IsOpen: &yes,
	}, nil)

	easternTime, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(easternTime)
	marketClose := time.Date(now.Year(), now.Month(), now.Day(), 16, 00, 0, 0, easternTime)
	timeTillClose := marketClose.Sub(now).Truncate(time.Minute)
	var wantBuilder strings.Builder
	fmt.Fprintf(&wantBuilder, "Market: ▼, open for %v\n"+
		"SPY: 544.51$, -0.13%%, -0.73$\n"+
		"QQQ: 480.18$, -0.27%%, -1.29$\n"+
		"AMZN: 189.08$, +1.60%%, +2.98$\n", timeTillClose)
	want := wantBuilder.String()
	wantErr := false

	tickerService := &tickerService{
		notificationApi: mockNotificationApi,
		stockApi:        mockStockApi,
	}
	got, err := tickerService.Run()
	if (err != nil) != wantErr {
		t.Errorf("Run() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != want {
		t.Errorf("Run() got = %v, want %v", got, want)
	}
}
