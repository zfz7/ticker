package apis

import (
	"context"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"os"
)

const FINNHUB_API_KEY = "FINNHUB_API_KEY"

type StockApi interface {
	Quote(ticker string) (finnhub.Quote, error)
	MarketStatus(exchange string) (finnhub.MarketStatus, error)
}

type stockApi struct {
	finnhubClient *finnhub.DefaultApiService
}

func NewStockApi() StockApi {
	finnHubCfg := finnhub.NewConfiguration()
	finnHubCfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv(FINNHUB_API_KEY))
	finnhubClient := finnhub.NewAPIClient(finnHubCfg).DefaultApi

	return &stockApi{
		finnhubClient: finnhubClient,
	}
}

func (stockService *stockApi) Quote(ticker string) (finnhub.Quote, error) {
	quote, _, err := stockService.finnhubClient.Quote(context.Background()).Symbol(ticker).Execute()
	if err != nil {
		return finnhub.Quote{}, err
	}
	return quote, nil
}

func (stockService *stockApi) MarketStatus(exchange string) (finnhub.MarketStatus, error) {
	marketStatus, _, err := stockService.finnhubClient.MarketStatus(context.Background()).Exchange(exchange).Execute()

	if err != nil {
		return finnhub.MarketStatus{}, err
	}

	return marketStatus, nil
}
