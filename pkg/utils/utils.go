package utils

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/uttamsutariya/crypto-pulse/internal/types"
)

func FormatExchangeMessageToKafka(jsonBatch []map[string]interface{}, isSpotSymbol bool) []types.KafkaMessage {
	kafkaBatchMessages := lo.Map(jsonBatch, func(item map[string]interface{}, index int) types.KafkaMessage {
		emptyMessage := types.KafkaMessage{}

		strTimestamp, ok := lo.ValueOr(item, "E", 0).(string)
		if !ok {
			currentTime := time.Now()
			x := currentTime.UnixNano() / int64(time.Millisecond)
			strTimestamp = strconv.Itoa(int(x))
		}

		timestamp, err := strconv.Atoi(strTimestamp)
		if err != nil {
			log.Println("Error converting timestamp")
			return emptyMessage
		}

		ltt := time.Unix(int64(timestamp)/1000, 0).UTC().Format(time.RFC3339)

		// each property received from binance is in string format
		bidStr, ok := lo.ValueOr(item, "b", "0").(string)
		if !ok {
			bidStr = "0"
		}

		bid, err := strconv.ParseFloat(bidStr, 64)
		if err != nil {
			log.Println("Error converting bid")
			return emptyMessage
		}

		bidQtyStr, ok := lo.ValueOr(item, "B", "0").(string)
		if !ok {
			bidQtyStr = "0"
		}

		bidQty, err := strconv.ParseFloat(bidQtyStr, 64)
		if err != nil {
			log.Println("Error converting bidQty", err)
			return emptyMessage
		}

		askStr, ok := lo.ValueOr(item, "a", "0").(string)
		if !ok {
			askStr = "0"
		}

		ask, err := strconv.ParseFloat(askStr, 64)
		if err != nil {
			log.Println("Error converting ask")
			return emptyMessage
		}

		askQtyStr, ok := lo.ValueOr(item, "A", "0").(string)
		if !ok {
			askQtyStr = "0"
		}

		askQty, err := strconv.ParseFloat(askQtyStr, 64)
		if err != nil {
			log.Println("Error converting askQty")
			return emptyMessage
		}

		ltpStr, ok := lo.ValueOr(item, "c", "0").(string)
		if !ok {
			ltpStr = "0"
		}

		ltp, err := strconv.ParseFloat(ltpStr, 64)
		if err != nil {
			log.Println("Error converting ltp")
			return emptyMessage
		}

		ltqStr, ok := lo.ValueOr(item, "Q", "0").(string)
		if !ok {
			ltqStr = "0"
		}

		ltq, err := strconv.ParseFloat(ltqStr, 64)
		if err != nil {
			log.Println("Error converting ltq", err)
			return emptyMessage
		}

		symbol, ok := lo.ValueOr(item, "s", "").(string)
		if !ok {
			symbol = ""
		}

		if isSpotSymbol {
			symbol = fmt.Sprintf("%s SPOT", symbol)
		} else {
			symbol = fmt.Sprintf("%s PERP", symbol)
		}

		return types.KafkaMessage{
			Bid:                bid,
			BidQty:             int64(bidQty),
			Ask:                ask,
			AskQty:             int64(askQty),
			LTP:                ltp,
			LastTradedQuantity: int64(ltq),
			Symbol:             symbol,
			LastTradeTime:      ltt,
			Vol:                0, // skip vol
			OI:                 0, // skip oi
		}
	})

	return kafkaBatchMessages
}
