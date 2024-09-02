package api

import (
	"github.com/devworlds/eventlistener-redis-performance/internal/api/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/wallets", handler.AddWalletHandler).Methods("POST")
	return r
}
