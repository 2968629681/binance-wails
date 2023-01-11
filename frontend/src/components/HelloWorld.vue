<script >
  import { dispose, init } from "klinecharts";

  export default{
    data(){
    return{
          websock: null, //建立的连接
          lockReconnect: false, //是否真正建立连接
          timeout: 60 * 1000 * 8, //8分组一次心跳
          timeoutObj: null, //心跳心跳倒计时
          serverTimeoutObj: null, //心跳倒计时
          timeoutnum: null, //断开 重连倒计时
          mapdata:[]
    }
    },
    mounted(){
        this.iniWebSocket();
         // 一般在mountd周期初始化,在destroyed周期清除
         this.kLineChart = init("init-kline", {
               // 蜡烛图
          candle: {
                // 蜡烛图上下间距，大于1为绝对值，大于0小余1则为比例
                margin: {
                    top: 0.2,
                    bottom: 0.1
            },
                          
                // 蜡烛图类型 'candle_solid'|'candle_stroke'|'candle_up_stroke'|'candle_down_stroke'|'ohlc'|'area'
                type: 'candle_solid',
                          
                // 蜡烛柱提颜色
                bar: {
                    upColor: '#26A69A',
                    downColor: '#EF5350',
                    noChangeColor: '#888888'
            },
                          
                // 面积图，只有当type为'area'时，当有效
                area: {
               lineSize: 2,
               lineColor: '#2196F3',
               value: 'close',
               backgroundColor: [{ 
                    offset: 0, //下半区颜色
                    color: 'rgba(33, 150, 243, 0.01)'
            }, {
                    offset: 1, //上半区颜色
                    color: 'rgba(33, 150, 243, 0.2)'
            }]
       },
                    
          // 价格标记
          priceMark: {
               show: true, //是否展示
               high: { // 最高价标记
                    show: true,
                    color: '#D9D9D9',
                    textMargin: 5,
                    textSize: 10,
                    textFamily: 'Helvetica Neue',
                    textWeight: 'normal'
            },
               low: { // 最低价标记
                    show: true,
                    color: '#D9D9D9',
                    textMargin: 5,
                    textSize: 10,
                    textFamily: 'Helvetica Neue',
                    textWeight: 'normal',
            },
               last: { // 最新价标记
                    show: true,
                    upColor: '#26A69A',
                    downColor: '#EF5350',
                    noChangeColor: '#888888',
                    line: {
                         show: true,
                         style: 'dash', // 'solid'|'dash'
                         dashValue: [4, 4], // 虚线时的紧密程度
                         size: 1 
                 },
                    text: { //最新价的字体样式
                         show: true,
                         size: 12,
                         paddingLeft: 2,
                         paddingTop: 2,
                         paddingRight: 2,
                         paddingBottom: 2,
                         color: '#FFFFFF',
                         family: 'Helvetica Neue',
                         weight: 'normal',
                         borderRadius: 2
                 }
            }
       },
                    
          // 提示 (顶部文案)
          tooltip: {
               // follow_cross是指图表失焦时不展示顶部文案
               showRule: 'always', // 'always' | 'follow_cross' | 'none'     
               // rect是指顶部文案垂直排列
               showType: 'standard', // 'standard' | 'rect'
               labels: ['时间', '开', '收', '高', '低', '成交量'],
               values: null,
               defaultValue: 'n/a',
               rect: { //垂直时排列样式
                    paddingLeft: 0,
                    paddingRight: 0,
                    paddingTop: 0,
                    paddingBottom: 6,
                    offsetLeft: 8,
                    offsetTop: 8,
                    offsetRight: 8,
                    borderRadius: 4,
                    borderSize: 1,
                    borderColor: '#3f4254',
                    backgroundColor: 'rgba(17, 17, 17, .3)'
            },
               text: { // 字体样式
                    size: 12,
                    family: 'Helvetica Neue',
                    weight: 'normal',
                    color: '#D9D9D9',
                    marginLeft: 8,
                    marginTop: 6,
                    marginRight: 8,
                    marginBottom: 0
                    }
              }
          },     
            // y轴
          yAxis: {
            show: true,
            size: 'auto',
            // 'left' | 'right'
            position: 'right',
            // 'normal' | 'percentage' | 'log'
            type: 'normal',
            inside: false,
            reverse: false,
            // y轴线
            axisLine: {
              show: true,
              color: '#888888',
              size: 1
            },
            // x轴分割文字
            tickText: {
              show: true,
              color: '#D9D9D9',
              family: 'Helvetica Neue',
              weight: 'normal',
              size: 12,
              marginStrat: 4,
              marginBottom: 4
            },
            // x轴分割线
            tickLine: {
              show: true,
              size: 1,
              length: 3,
              color: '#888888'
            }
          }, 
      });
    },
    methods: {
        iniWebSocket() {
          //初始化websocket
          const wsuri ="ws://127.0.0.1:2303/ws";
          this.websock = new WebSocket(wsuri);
          this.websock.onmessage = this.websocketonmessage;
          this.websock.onopen = this.websocketonopen;
          this.websock.onerror = this.websocketonerror;
          this.websock.onclose = this.websocketclose;
        },
        reconnect() {
          //重新连接
          var that = this;
          if (that.lockReconnect) {
            return;
          }
          that.lockReconnect = true;
          //没连接上会一直重连，设置延迟避免请求过多
          that.timeoutnum && clearTimeout(that.timeoutnum);
          that.timeoutnum = setTimeout(function () {
            //新连接
            that.iniWebSocket();
            that.lockReconnect = false;
          }, 1000);
        },
        reset() {
          //重置心跳
          var that = this;
          //清除时间
          clearTimeout(that.timeoutObj);
          clearTimeout(that.serverTimeoutObj);
          //重启心跳
          that.start();
        },
        start() {
          //开启心跳
          var self = this;
          self.timeoutObj && clearTimeout(self.timeoutObj);
          self.serverTimeoutObj && clearTimeout(self.serverTimeoutObj);
          self.timeoutObj = setTimeout(function () {
            //这里发送一个心跳，后端收到后，返回一个心跳消息
            if (self.websock.readyState == 'PONG') {
              //如果连接正常
              let temp = {"method": "PING",}
              self.websock.send(JSON.stringify(temp));
            } else {
              //否则重连
              self.reconnect();
            }
            self.serverTimeoutObj = setTimeout(function () {
              //超时关闭
              self.websock.close();
            }, self.timeout);
          }, self.timeout);
        },
        websocketonopen() {
          //连接建立之后执行send方法发送数据
          let actions = { 
            "method": "SUBSCRIBE",
            "param": "KLINE@BTCUSDT@1m"
          };
          this.websocketsend(JSON.stringify(actions));
          console.log("连接");
          //开启心跳
          this.start();
        },
        websocketonerror() {
          //连接建立失败重连
          this.iniWebSocket();
          this.reconnect();
          console.log("重新连接");
        },
        websocketonmessage(e) {
          //数据接收
          const redata = JSON.parse(e.data);
          let temp={
            open: null, // 开盘价，必要字段
            close: null, // 收盘价，必要字段
            high: null, // 最高价，必要字段
            low: null,  // 最低价，必要字段
            volume: null, // 成交量，非必须字段
            turnover: null, // 成交额，非必须字段，如果需要展示技术指标'EMV'和'AVP'，则需要为该字段填充数据
            timestamp: null, // 时间戳，毫秒级别，必要字段
            time: null // 时间戳，一分钟更新一次
          }
          temp.open=parseFloat(redata.k.o)
          temp.close=parseFloat(redata.k.c)
          temp.high=parseFloat(redata.k.h)
          temp.low=parseFloat(redata.k.l)
          temp.volume=parseFloat(redata.k.v)
          temp.turnover=parseFloat(redata.k.q)
          temp.timestamp=redata.k.T
          let len=this.mapdata.length
          if(len == 0) this.mapdata.push(temp)
          else if(temp.timestamp != this.mapdata[len-1].timestamp) this.mapdata.push(temp)
          else this.mapdata[len-1]=temp
          this.kLineChart.applyNewData(this.mapdata)
          this.reset();
        },
        websocketsend(Data) {
          //数据发送
          this.websock.send(Data);
        },
        websocketclose(e) {
          //关闭
          console.log("断开连接", e);
          //重连
          this.reconnect();
        },
    },
    destroyed () {
      dispose("update-k-line");
    }
  }
</script>

<template>
  <div id="init-kline" style="width: 80%;height: 80%;"/>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
