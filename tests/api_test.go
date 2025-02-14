package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateDepositOperation(t *testing.T) {
    op := map[string]interface{}{
        "operationType": "DEPOSIT",
        "amount": 500.00,
    }
    data, _ := json.Marshal(op)

    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/wallets/1/operations", bytes.NewBuffer(data))
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusOK, res.StatusCode)

    var response map[string]float64
    json.NewDecoder(res.Body).Decode(&response)
    assert.Greater(t, response["balance"], 0.0, "Balance should be updated after deposit")
}

func TestCreateWithdrawOperation(t *testing.T) {
    op := map[string]interface{}{
        "operationType": "WITHDRAW",
        "amount": 200.00,
    }
    data, _ := json.Marshal(op)

    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/wallets/2/operations", bytes.NewBuffer(data))
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusOK, res.StatusCode)

    var response map[string]float64
    json.NewDecoder(res.Body).Decode(&response)
    assert.Less(t, response["balance"], 500.0, "Balance should decrease after withdrawal")
}

func TestCreateWithdrawWithInsufficientFunds(t *testing.T) {
    op := map[string]interface{}{
        "operationType": "WITHDRAW",
        "amount": 100.00,  // Exceeding balance
    }
    data, _ := json.Marshal(op)

    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/wallets/1/operations", bytes.NewBuffer(data))
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusInternalServerError, res.StatusCode)  // Expecting error due to insufficient funds
}

func TestGetBalance(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/wallets/3", nil)
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusOK, res.StatusCode)

    var response map[string]float64
    json.NewDecoder(res.Body).Decode(&response)
    assert.Equal(t, 1601.5, response["balance"], "Balance should be correct")
}

func TestInvalidOperationType(t *testing.T) {
    op := map[string]interface{}{
        "operationType": "INVALID",
        "amount": 100.00,
    }
    data, _ := json.Marshal(op)

    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/wallets/1/operations", bytes.NewBuffer(data))
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusBadRequest, res.StatusCode)  // Expecting error due to invalid operation type
}

func TestNonExistingWallet(t *testing.T) {
    op := map[string]interface{}{
        "operationType": "DEPOSIT",
        "amount": 100.00,
    }
    data, _ := json.Marshal(op)

    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/wallets/non-existing-wallet/operations", bytes.NewBuffer(data))
    res, _ := http.DefaultClient.Do(req)

    assert.Equal(t, http.StatusNotFound, res.StatusCode)  // Expecting not found error for non-existing wallet
}