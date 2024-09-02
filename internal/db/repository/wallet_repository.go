package repository

import (
	"fmt"

	"github.com/devworlds/eventlistener-redis-performance/internal/db"
)

type Wallet struct {
	Address string
}

func GetWallets() ([]Wallet, error) {
	rows, err := db.DB.Query("SELECT address FROM wallets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var wallet Wallet
		if err := rows.Scan(&wallet.Address); err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	db.DbCall++
	fmt.Print("DBCall: ", db.DbCall, "\n")
	return wallets, nil
}

func AddWallet(address string) error {
	_, err := db.DB.Exec(
		"INSERT INTO wallets (address) VALUES ($1)",
		address,
	)
	if err != nil {
		return err
	}
	return nil
}
