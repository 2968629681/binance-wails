package binancews

const (
	SUBSCRIBE   = "SUBSCRIBE"
	UNSUBSCRIBE = "UNSUBSCRIBE"
	PING        = "PING"
)

const (
	DEPTH      = "DEPTH"
	DEPTH100MS = "DEPTH100MS"
	KLINE      = "KLINE"
	AGGTRADE   = "AGGTRADE"
)

type RequestMessage struct {
	Method string `json:"method"`
	Param  string `json:"param"`
}
