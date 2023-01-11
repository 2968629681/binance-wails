package binancews

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	binance "github.com/adshao/go-binance/v2"
	"github.com/gorilla/websocket"
)

var (
	heartbeatInterval = 10 * time.Minute
)

type WebsocketFunc func() (chan struct{}, chan struct{}, error)
type BinanceWs struct {
	ws *websocket.Conn
	f  WebsocketFunc
}

func NewBinanceWs(ws *websocket.Conn) *BinanceWs {
	return &BinanceWs{
		ws: ws,
	}
}

func (b *BinanceWs) Run() {
	var (
		err error

		doneC chan struct{}
		stopC chan struct{}

		msgCh     = make(chan RequestMessage)
		heartbeat = time.Now().Add(heartbeatInterval)

		wsRunning atomic.Bool
	)
	wsRunning.Store(false)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go b.listenReadMessage(msgCh, ctx)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-doneC:
			wsRunning.Store(false)
			return
		case msg := <-msgCh:
			heartbeat = time.Now().Add(heartbeatInterval)

			switch msg.Method {
			case SUBSCRIBE:
				if wsRunning.Load() {
					b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Running...")))
					continue
				}

				err = b.Parse(msg)
				if err != nil {
					b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
					continue
				}

				doneC, stopC, err = b.f()
				if err != nil {
					b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
					return
				}
				wsRunning.Store(true)
			case UNSUBSCRIBE:
				if wsRunning.Load() {
					stopC <- struct{}{}
					wsRunning.Store(false)
				}
			case PING:
				b.ws.WriteMessage(websocket.TextMessage, []byte("PONG"))
			default:
				continue
			}
		case <-ticker.C:
			if time.Now().Sub(heartbeat).Seconds() >= 0 {
				b.ws.WriteMessage(websocket.TextMessage, []byte("Connect expired"))
				if wsRunning.Load() {
					stopC <- struct{}{}
					wsRunning.Store(false)
				}
				return
			}
		}
	}
}

var (
	ParamsError = errors.New("Parmas error")
)

func (b *BinanceWs) Parse(r RequestMessage) error {
	strs := strings.Split(r.Param, "@")
	if len(strs) < 1 {
		return ParamsError
	}

	switch strings.ToUpper(strs[0]) {
	// DEPTH@symbel@level
	case DEPTH:
		if len(strs) < 3 {
			return ParamsError
		}
		b.partialDepthServeFunc(strs[1], strs[2])

	// DEPTH100MS@symbel@level
	case DEPTH100MS:
		if len(strs) < 3 {
			return ParamsError
		}
		b.partialDepthServe100MsFunc(strs[1], strs[2])

	// KLINE@symbol@interval
	case KLINE:
		if len(strs) < 3 {
			return ParamsError
		}
		b.kLineServeFunc(strs[1], strs[2])
	default:
		return ParamsError
	}

	return nil
}

func (b *BinanceWs) listenReadMessage(ch chan RequestMessage, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			var message = RequestMessage{}
			err := b.ws.ReadJSON(&message)
			if err != nil {
				b.ws.WriteMessage(websocket.TextMessage, []byte("Read message failed"))
				return
			}

			log.Printf("Recive message: [method: %v params: %v]\n", message.Method, message.Param)
			ch <- message
		}
	}

}

func (b *BinanceWs) partialDepthServe100MsFunc(symbol, level string) {
	wsDepthHandler := func(event *binance.WsPartialDepthEvent) {
		b.ws.WriteJSON(event)
	}

	errHandler := func(err error) {
		b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
	}

	b.f = func() (chan struct{}, chan struct{}, error) {
		return binance.WsPartialDepthServe100Ms(symbol, level, wsDepthHandler, errHandler)
	}
}

func (b *BinanceWs) partialDepthServeFunc(symbol, level string) {
	wsDepthHandler := func(event *binance.WsPartialDepthEvent) {
		b.ws.WriteJSON(event)
	}

	errHandler := func(err error) {
		b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
	}

	b.f = func() (chan struct{}, chan struct{}, error) {
		return binance.WsPartialDepthServe(symbol, level, wsDepthHandler, errHandler)
	}
}

func (b *BinanceWs) kLineServeFunc(symbol, interval string) {
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		b.ws.WriteJSON(event)
	}

	errHandler := func(err error) {
		b.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
	}

	b.f = func() (chan struct{}, chan struct{}, error) {
		return binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	}
}
