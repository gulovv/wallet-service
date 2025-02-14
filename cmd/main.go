package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/gulovv/wallet-service/internal/api"
    "github.com/gulovv/wallet-service/internal/db"
    "github.com/gulovv/wallet-service/pkg/config"
)

func main() {
    cfg := config.LoadConfig()

    var err error
    maxRetries := 10
    for i := 0; i < maxRetries; i++ {
        err = db.InitDB(cfg.DBConnString)
        if err == nil {
            break
        }
        log.Printf("Failed to initialize database (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(5 * time.Second) 
    }

    if err != nil {
        log.Fatalf("Failed to initialize database after %d attempts: %v", maxRetries, err)
    }

    r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/{walletId}/operations", api.CreateOperationHandler).Methods("POST")  
	r.HandleFunc("/api/v1/wallets/{walletId}", api.GetBalanceHandler).Methods("GET")

    log.Printf("Server started at %s", cfg.ServerAddress)
    log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}