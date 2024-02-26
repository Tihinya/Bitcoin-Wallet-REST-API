package controllers

import (
	errorhandler "bitcoin-wallet/errorHandler"
	"bitcoin-wallet/models"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	sql "bitcoin-wallet/db/database"
)

func GetListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := sql.GetAllTransactions()
	if err != nil {
		log.Println(err)
		errorhandler.ReturnJsonMessage(w, "Error getting all transactions", http.StatusInternalServerError, "error")
	}

	json.NewEncoder(w).Encode(models.Transactions{
		Transactions: transactions,
	})
}

func GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	currentBalance, err := sql.GetCurrentBTCBalance()
	if err != nil {
		log.Println(err)
	}
	bitcoinPrice, err := getCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		errorhandler.ReturnJsonMessage(w, "Error getting Bitcoin price", http.StatusInternalServerError, "error")
	}

	json.NewEncoder(w).Encode(models.CurrentBalance{
		BTC: currentBalance,
		EUR: math.Round(bitcoinPrice*currentBalance*100) / 100,
	})
}

func getCurrentBTCPrice() (float64, error) {
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
		return 0, fmt.Errorf("Invalid data: can't find price resource value")
	}

	bitcoinPrice, err := strconv.ParseFloat(priceResourceCollection.Data[0].Value, 64)
	if err != nil {
		return 0, err
	}

	return bitcoinPrice, nil
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		errorhandler.ReturnJsonMessage(w, "Bad request", http.StatusBadRequest, "error")
		return
	}
	bitcoinPrice, err := getCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		errorhandler.ReturnJsonMessage(w, "Error getting Bitcoin price", http.StatusInternalServerError, "error")
	}
	// transaction.Amount = toFixed(transaction.Amount/bitcoinPrice, 8)
	a := fmt.Sprintf("%.8f", transaction.Amount/bitcoinPrice)
	transaction.Amount, _ = strconv.ParseFloat(a, 64)

	currentBalance, err := sql.GetCurrentBTCBalance()
	if err != nil {
		errorhandler.ReturnJsonMessage(w, "Error getting current balance:", http.StatusInternalServerError, "error")
		return
	}
	if transaction.Amount < 0.00001 {
		errorhandler.ReturnJsonMessage(w, "Transfer amount should be greater than or equal to 0.00001 BTC.:", http.StatusBadRequest, "error")
		return
	}
	if transaction.Amount > currentBalance {
		errorhandler.ReturnJsonMessage(w, "Not enough funds on balance. The transfer amount exceeds the available balance.:", http.StatusInternalServerError, "error")
		return
	}
	allUnspentTransactions, err := sql.GetAllUnspentTransactions()
	if err != nil {
		log.Println(err)
		errorhandler.ReturnJsonMessage(w, "Error getting all Untransactions", http.StatusInternalServerError, "error")
	}
	sum := 0.0
	var transactionIds []string
	for _, transaction1 := range allUnspentTransactions {
		sum += transaction1.Amount
		transactionIds = append(transactionIds, transaction1.TransactionId)
		if sum > transaction.Amount {
			break
		}
	}
	if sum != transaction.Amount {
		err := sql.CreateTransaction(sum-transaction.Amount, false)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusConflict)
			errorhandler.ReturnJsonMessage(w, "Data has not been written to the database", http.StatusConflict, "error")
		}
	}

	for _, transactionId := range transactionIds {
		err := sql.MarkTransactionAsSpent(transactionId)
		if err != nil {
			log.Println(err)
			errorhandler.ReturnJsonMessage(w, "Error mark transaction as spent", http.StatusInternalServerError, "error")
		}
	}

	if err := sql.CreateTransaction(transaction.Amount, true); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		errorhandler.ReturnJsonMessage(w, "Data has not been written to the database", http.StatusConflict, "error")
		return
	}
	errorhandler.ReturnJsonMessage(w, "Transmission is well-processed", http.StatusOK, "success")

}

func TopUp(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		errorhandler.ReturnJsonMessage(w, "Bad request", http.StatusBadRequest, "error")
		return
	}
	bitcoinPrice, err := getCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		errorhandler.ReturnJsonMessage(w, "Error getting Bitcoin price", http.StatusInternalServerError, "error")
	}

	// transaction.Amount = toFixed(transaction.Amount/bitcoinPrice, 8)
	a := fmt.Sprintf("%.8f", transaction.Amount/bitcoinPrice)
	transaction.Amount, _ = strconv.ParseFloat(a, 64)

	if transaction.Amount < 0.00001 {
		price := math.Round(bitcoinPrice*0.00001*100) / 100
		errorhandler.ReturnJsonMessage(w, "Transfer amount should be greater than or equal to "+fmt.Sprintf("%f", price), http.StatusBadRequest, "error")
		return
	}

	if err := sql.CreateTransaction(transaction.Amount, false); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		errorhandler.ReturnJsonMessage(w, "Data has not been written to the database", http.StatusConflict, "error")
		return
	}
	errorhandler.ReturnJsonMessage(w, "Transmission is well-processed", http.StatusOK, "success")

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
