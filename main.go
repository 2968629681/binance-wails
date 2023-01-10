package main

import (
	binancews "binance-wails/binance-ws"
	"binance-wails/conf"
	"embed"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	conf.Init()
}

func main() {
	go binancews.Run()

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "binance-wails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
