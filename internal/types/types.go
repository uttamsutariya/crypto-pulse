package types

type KafkaMessage struct {
	Bid                float64 `json:"bid"`
	BidQty             int64   `json:"bidQty"`
	Ask                float64 `json:"ask"`
	AskQty             int64   `json:"askQty"`
	Vol                int64   `json:"vol"`
	LTP                float64 `json:"ltp"`
	OI                 int64   `json:"oi"`
	LastTradedQuantity int64   `json:"last_traded_quantity"`
	Symbol             string  `json:"symbol"`
	LastTradeTime      string  `json:"last_trade_time"`
}
