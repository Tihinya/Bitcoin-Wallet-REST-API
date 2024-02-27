package api

import (
	"bitcoin-wallet/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetCurrentBTCPrice() (float64, error) {
	req, err := http.NewRequest("GET", "http://api-cryptopia.adca.sh/v1/prices/history", nil)
	if err != nil {
		return 0, err
	}
	q := req.URL.Query()

	q.Add("time", time.Now().Format("2006-01-02T15:04:05Z"))
	q.Add("symbol", "BTC/EUR")

	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	priceResourceCollection := &models.PriceResourceCollection{}
	err = json.NewDecoder(res.Body).Decode(priceResourceCollection)
	if err != nil {
		return 0, err
	}

	if len(priceResourceCollection.Data) < 1 || priceResourceCollection.Data[0].Value == "" {
		return 0, fmt.Errorf("invalid data: can't find price resource value")
	}

	bitcoinPrice, err := strconv.ParseFloat(priceResourceCollection.Data[0].Value, 64)
	if err != nil {
		return 0, err
	}

	return bitcoinPrice, nil
}
