package binancews

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func RunServer() {
	r := gin.Default()
	r.GET(viper.GetString("binancews.path"), binanceWs)
	log.Panic(r.Run(viper.GetString("binancews.addr")))
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func binanceWs(c *gin.Context) {
	log.Println("Get a request from: ", c.Request.RemoteAddr)

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		log.Println("Close connection: ", c.Request.RemoteAddr)
		ws.Close()
	}()

	NewBinanceWs(ws).Run()
}
