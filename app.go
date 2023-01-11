package main

import (
	binancehttp "binance-wails/binance-http"
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) KlineService(symbol string, interval string, limit int) (res []*binance.Kline, err error) {
	return binancehttp.KlinesService(symbol, interval, limit)
}

func (a *App) KlineServiceWithStartAndEnd(symbol string, start, end int64, interval string, limit int) (res []*binance.Kline, err error) {
	return binancehttp.KlinesServiceWithStartAndEnd(symbol, start, end, interval, limit)
}
