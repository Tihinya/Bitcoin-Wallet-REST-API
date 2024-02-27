package controllers

import (
	"bitcoin-wallet/api"
	"bitcoin-wallet/models"
	responseHandler "bitcoin-wallet/responseHandler"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	sql "bitcoin-wallet/db/database"
)

func GetListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := sql.GetAllTransactions()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting all transactions")
		return
	}
	w.WriteHeader(http.StatusOK)
	responseHandler.ReturnJsonMessage(w, "success", "", models.Transactions{
		Transactions: transactions,
	})
}

func GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	currentBalance, err := sql.GetCurrentBTCBalance()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting current balance")
		return
	}
	bitcoinPrice, err := api.GetCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting Bitcoin price")
		return
	}

	w.WriteHeader(http.StatusOK)
	responseHandler.ReturnJsonMessage(w, "success", "", models.CurrentBalance{
		BTC: currentBalance,
		EUR: math.Round(bitcoinPrice*currentBalance*100) / 100,
	})
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		responseHandler.ReturnJsonMessage(w, "error", "Bad request")
		return
	}
	bitcoinPrice, err := api.GetCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting Bitcoin price")
		return
	}
	a := fmt.Sprintf("%.8f", transaction.Amount/bitcoinPrice)
	transaction.Amount, _ = strconv.ParseFloat(a, 64)

	currentBalance, err := sql.GetCurrentBTCBalance()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting current balance")
		return
	}
	if transaction.Amount < 0.00001 {
		w.WriteHeader(http.StatusBadRequest)
		responseHandler.ReturnJsonMessage(w, "error", "Transfer amount should be greater than or equal to 0.00001 BTC")
		return
	}
	if transaction.Amount > currentBalance {
		w.WriteHeader(http.StatusBadRequest)
		responseHandler.ReturnJsonMessage(w, "error", "Not enough funds on balance. The transfer amount exceeds the available balance")
		return
	}
	allUnspentTransactions, err := sql.GetAllUnspentTransactions()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting all Untransactions")
		return
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
			responseHandler.ReturnJsonMessage(w, "error", "Data has not been written to the database")
			return
		}
	}

	for _, transactionId := range transactionIds {
		err := sql.MarkTransactionAsSpent(transactionId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			responseHandler.ReturnJsonMessage(w, "error", "Error mark transaction as spent")
			return
		}
	}

	if err := sql.CreateTransaction(transaction.Amount, true); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		responseHandler.ReturnJsonMessage(w, "error", "Data has not been written to the database")
		return
	}
	w.WriteHeader(http.StatusOK)
	responseHandler.ReturnJsonMessage(w, "success", "Transfer is well-processed")

}

func TopUp(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		responseHandler.ReturnJsonMessage(w, "error", "Bad request")
		return
	}
	bitcoinPrice, err := api.GetCurrentBTCPrice()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		responseHandler.ReturnJsonMessage(w, "error", "Error getting Bitcoin price")
		return
	}

	a := fmt.Sprintf("%.8f", transaction.Amount/bitcoinPrice)
	transaction.Amount, _ = strconv.ParseFloat(a, 64)

	if transaction.Amount < 0.00001 {
		price := math.Round(bitcoinPrice*0.00001*100) / 100

		w.WriteHeader(http.StatusBadRequest)
		responseHandler.ReturnJsonMessage(w, "error", "Transfer amount should be greater than or equal to "+fmt.Sprintf("%f", price))
		return
	}

	if err := sql.CreateTransaction(transaction.Amount, false); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		responseHandler.ReturnJsonMessage(w, "error", "Data has not been written to the database")
		return
	}
	w.WriteHeader(http.StatusOK)
	responseHandler.ReturnJsonMessage(w, "success", "Top-Up is well-processed")
}
