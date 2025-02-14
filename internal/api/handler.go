package api

import (
	"encoding/json"
	"net/http"
	"github.com/gulovv/wallet-service/internal/service"
	"github.com/gorilla/mux"
)
func CreateOperationHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    walletID := vars["walletId"]  

    var op service.Operation
    if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    balance, err := service.ProcessOperation(walletID, op)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["walletId"]

	balance, err := service.GetBalance(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}