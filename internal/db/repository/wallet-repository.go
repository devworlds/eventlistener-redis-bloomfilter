package repository

import (
	"database/sql"
	"fmt"

	"github.com/devworlds/eventlistener-redis-performance/internal/db"
)

type Wallet struct {
	Address string
}

func GetWallets(DB *sql.DB) ([]Wallet, error) {
	rows, err := DB.Query("SELECT address FROM wallets")
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
	fmt.Print("Databaase Call:", db.DbCall, "\n")
	return wallets, nil
}
