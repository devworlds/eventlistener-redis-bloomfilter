package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devworlds/eventlistener-redis-performance/internal/api"
)

func apiInit() {
	router := api.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
