package service

import (
	"github.com/gulovv/wallet-service/internal/db"
	"fmt"
)

type Operation struct {
	WalletID     string  `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount       float64 `json:"amount"`
}

func ProcessOperation(WalletID string, op Operation) (float64, error) {
	if op.OperationType == "DEPOSIT" {
		return db.UpdateBalance(WalletID, op.Amount)
	} else if op.OperationType == "WITHDRAW" {
		return db.UpdateBalance(WalletID, -op.Amount)
	}
	return 0, fmt.Errorf("invalid operation type")
}

func GetBalance(walletID string) (float64, error) {
	return db.GetBalance(walletID)
}

