package models

import (
	sql "bitcoin-wallet/db/database"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Object  []interface{}  `json:"object,omitempty"`
}

type Transaction struct {
	Amount float64 `json:"amount,omitempty"`
}

type PriceResource struct {
	Symbol     string `json:"symbol"`
	Value      string `json:"value"`
	Source     int    `json:"source"`
	Updated_at string `json:"updated_at"`
}

type PriceResourceCollection struct {
	Data []PriceResource `json:"data"`
}

type CurrentBalance struct {
	BTC float64 `json:"btc"`
	EUR float64 `json:"eur"`
}

type Transactions struct {
	Transactions []sql.Transaction
}
