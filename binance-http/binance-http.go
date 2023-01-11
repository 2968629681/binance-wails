package binancehttp

import (
	"context"

	"github.com/adshao/go-binance/v2"
	"github.com/spf13/viper"
)

func getClient() *binance.Client {
	return binance.NewClient(
		viper.GetString("binance.apiKey"),
		viper.GetString("binance.secretKey"),
	)
}

func KlinesService(symbol string, interval string, limit int) (res []*binance.Kline, err error) {
	return getClient().NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
}

func KlinesServiceWithStartAndEnd(symbol string, start, end int64, interval string, limit int) (res []*binance.Kline, err error) {
	return getClient().NewKlinesService().
		Symbol(symbol).
		StartTime(start).
		EndTime(end).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
}
