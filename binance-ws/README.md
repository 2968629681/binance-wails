# Websocket 

> websocket地址: addr/ws

addr可在config.yaml中修改: binancews.addr


## 消息订阅/取消

### 结构体定义

```go
type RequestMessage struct {
	Method string `json:"method"`
	Param  string `json:"param"`
}
```

### 请求

```json
{
	"method": "method",
	"param": "param"
}
```

### Method

* SUBSCRIBE: 订阅
* UNSUBSCRIBE: 取消订阅
* PING: 心跳

### Param

Param只有订阅的时候需要填写

#### Depth
> param: "DEPTH@symbol@level"

**Level**: 

* 5
* 10
* 20

```json
{
	"method": "SUBSCRIBE",
	"param": "DEPTH@BTCUSDT@20"
}
```

#### Depth100Ms
> param: "DEPTH100MS@symbol@level"

**Level:**
* 5
* 10
* 20

```json
{
	"method": "SUBSCRIBE",
	"param": "DEPTH100MS@BTCUSDT@20"
}
```

#### KLINE
> param: "KLINE@symbol@interval"

**Interval:**
* 1s 
* 1m
* 3m
* 5m
* 15m
* 30m
* 1h
* 2h
* 4h
* 6h
* 8h
* 12h
* 1d
* 3d
* 1w
* 1M

```json
{
	"method": "SUBSCRIBE",
	"param": "KLINE@BTCUSDT@1m"
}
```
