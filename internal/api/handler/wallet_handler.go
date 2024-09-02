package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devworlds/eventlistener-redis-performance/internal/db/service"
)

type AddWalletRequest struct {
	Address string `json:"address"`
}

func AddWalletHandler(w http.ResponseWriter, r *http.Request) {
	var req AddWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := service.CreateWallet(req.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
