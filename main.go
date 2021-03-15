package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance"
	"github.com/sunabozu/binance-price-change-go/utils"
)

func main() {
	parentPath, err := utils.GetParentPath()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	keys, err := utils.LoadKeys(parentPath + "/env.json")

	if err != nil {
		os.Exit(1)
	}

	client := binance.NewClient(keys.BinanceKey, keys.BinanceSecret)

	ticker := time.NewTicker(time.Second * 20)
	var drop20, drop15, drop10, drop5 bool = true, true, true, true

	for _ = range ticker.C {
		stats, err := client.NewListPriceChangeStatsService().Symbol("BTCBUSD").Do(context.Background())

		if err != nil || len(stats) < 1 {
			log.Println(err)
			continue
		}

		// log.Printf("%+v", stats[0])

		relativeChange, err := strconv.ParseFloat(stats[0].PriceChangePercent, 64)
		change, err := strconv.ParseFloat(stats[0].PriceChange, 64)
		lastPrice, err := strconv.ParseFloat(stats[0].LastPrice, 64)
		highPrice, err := strconv.ParseFloat(stats[0].HighPrice, 64)
		log.Print(relativeChange)

		if err != nil {
			continue
		}

		msg := " BTC dropped by "
		// log.Print(relativeChange, "%")
		if relativeChange <= -20.0 && drop20 {
			msg = "⚫️⚫️⚫️" + msg + "20%"
			go disableFor24h(&drop15)
		} else if relativeChange <= -15.0 && drop15 {
			msg = "🟣🟣🟣" + msg + "15%"
			go disableFor24h(&drop15)
		} else if relativeChange <= -10.0 && drop10 {
			msg = "🔴🔴🔴" + msg + "10%"
			go disableFor24h(&drop10)
		} else if relativeChange <= -5.0 && drop5 {
			msg = "🟡🟡🟡" + msg + "5%"
			go disableFor24h(&drop5)
		} else {
			continue
		}

		// msg += " ($" + stats[0].PriceChange + "), from $" + stats[0].HighPrice + " to $" + stats[0].LastPrice
		msg += fmt.Sprintf("($ %.0f), from $%.0f to $%.0f", change, highPrice, lastPrice)

		go utils.SendPushNotification(keys, msg)
	}
}

func disableFor24h(val *bool) {
	*val = false
	log.Print("disabling...", *val)
	time.Sleep(time.Hour * 24)
	*val = true
	log.Print("enabling...", *val)
}
