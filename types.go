package binance

const SymbolSOLBUSD = "SOLBUSD"
const SymbolSOLUSDT = "SOLUSDT"

type FutureOrder struct {
	Symbol            string `json:"symbol"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
	Price             string `json:"price"`
	OrigQuantity      string `json:"origQty"`
	ExecutedQuantity  string `json:"executedQty"`
	CumQuote          string `json:"cumQuote"`
	ReduceOnly        bool   `json:"reduceOnly"`
	StopPrice         string `json:"stopPrice"`
	UpdateTime        int64  `json:"updateTime"`
	ActivatePrice     string `json:"activatePrice"`
	PriceRate         string `json:"priceRate"`
	AvgPrice          string `json:"avgPrice"`
	ClosePosition     bool   `json:"closePosition"`
	PriceProtect      bool   `json:"priceProtect"`
	RateLimitOrder10s string `json:"rateLimitOrder10s,omitempty"`
	RateLimitOrder1m  string `json:"rateLimitOrder1m,omitempty"`
}

type FuturePosition struct {
	EntryPrice       string `json:"entryPrice"`
	MarginType       string `json:"marginType"`
	IsAutoAddMargin  string `json:"isAutoAddMargin"`
	IsolatedMargin   string `json:"isolatedMargin"`
	Leverage         string `json:"leverage"`
	LiquidationPrice string `json:"liquidationPrice"`
	MarkPrice        string `json:"markPrice"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	PositionAmt      string `json:"positionAmt"`
	Symbol           string `json:"symbol"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	PositionSide     string `json:"positionSide"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
}
