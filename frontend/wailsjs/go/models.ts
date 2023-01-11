export namespace binance {
	
	export class Kline {
	    openTime: number;
	    open: string;
	    high: string;
	    low: string;
	    close: string;
	    volume: string;
	    closeTime: number;
	    quoteAssetVolume: string;
	    tradeNum: number;
	    takerBuyBaseAssetVolume: string;
	    takerBuyQuoteAssetVolume: string;
	
	    static createFrom(source: any = {}) {
	        return new Kline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.openTime = source["openTime"];
	        this.open = source["open"];
	        this.high = source["high"];
	        this.low = source["low"];
	        this.close = source["close"];
	        this.volume = source["volume"];
	        this.closeTime = source["closeTime"];
	        this.quoteAssetVolume = source["quoteAssetVolume"];
	        this.tradeNum = source["tradeNum"];
	        this.takerBuyBaseAssetVolume = source["takerBuyBaseAssetVolume"];
	        this.takerBuyQuoteAssetVolume = source["takerBuyQuoteAssetVolume"];
	    }
	}

}

