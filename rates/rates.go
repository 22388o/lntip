package rates

import (
	"net/http"
	"time"

	coingecko "github.com/superoo7/go-gecko/v3"
	"go.uber.org/zap"
)

var client = coingecko.NewClient(http.DefaultClient)

var Price float32

func Run() {
	for {
		resp, err := client.SimpleSinglePrice("bitcoin", "eur")
		if err != nil {
			zap.S().Error(err)
			time.Sleep(time.Second * 10)
			continue
		}

		Price = resp.MarketPrice
		time.Sleep(time.Second * 60)
	}
}
