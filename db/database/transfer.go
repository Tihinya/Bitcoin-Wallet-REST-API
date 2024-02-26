package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func CreateTransaction(transferAmount float64, spent bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlStmt, err := db.PrepareContext(ctx, `INSERT INTO transactions(amount, spent, created_at) VALUES($1, $2, LOCALTIMESTAMP)`)
	if err != nil {
		return err
	}

	_, err = sqlStmt.Exec(transferAmount, spent)
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentBTCBalance() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT SUM(amount) FROM transactions WHERE spent = false"

	var totalAmount float64
	err := db.QueryRowContext(ctx, query).Scan(&totalAmount)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return totalAmount, nil
}
func GetAllTransactions() ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT transaction_id, amount, spent, created_at FROM transactions"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.TransactionId, &transaction.Amount, &transaction.Spent, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func GetAllUnspentTransactions() ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT transaction_id, amount FROM transactions WHERE spent = false"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.TransactionId, &transaction.Amount)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func MarkTransactionAsSpent(transactionID string) error {
	query := "UPDATE transactions SET spent = true WHERE transaction_id = $1"
	_, err := db.Exec(query, transactionID)
	if err != nil {
		return fmt.Errorf("Error updating spent status: %v", err)
	}

	return nil
}
