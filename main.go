package main

import (
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

	errHandler := func(err error) {
		log.Printf("error in the handler: %+v", err)
	}

	var drop20, drop15, drop10, drop5 bool = true, true, true, true

	aggHandler := func(event *binance.WsMarketStatEvent) {
		change, err := strconv.ParseFloat(event.PriceChange, 64)
		lastPrice, err := strconv.ParseFloat(event.LastPrice, 64)

		relativeChange := change / lastPrice * 100

		log.Print("Relative change: ", relativeChange)

		if err != nil {
			return
		}

		msg := "BTC dropped by "
		// log.Print(relativeChange, "%")
		if relativeChange <= -20.0 && drop20 {
			msg += "20%! ⚫️⚫️⚫️"
			go disableFor24h(&drop15)
		} else if relativeChange <= -15.0 && drop15 {
			msg += "15%! 🟣🟣🟣"
			go disableFor24h(&drop15)
		} else if relativeChange <= -10.0 && drop10 {
			msg += "10%! 🔴🔴🔴"
			go disableFor24h(&drop10)
		} else if relativeChange <= -5.0 && drop5 {
			msg += "5%! 🟡🟡🟡"
			go disableFor24h(&drop5)
		} else {
			return
		}

		go utils.SendPushNotification(keys, msg)
	}

	doneC, _, err := binance.WsMarketStatServe("BTCBUSD", aggHandler, errHandler)
	log.Printf("a channel: %+v\n error: %+v\n", <-doneC, err)
}

func disableFor24h(val *bool) {
	*val = false
	log.Print("disabling...", *val)
	time.Sleep(time.Hour * 24)
	*val = true
	log.Print("enabling...", *val)
}
