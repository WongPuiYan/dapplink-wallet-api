package tron

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTronDataClient(t *testing.T) {
	tests := []struct {
		name    string
		baseUrl string
		apiKey  string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "create client with valid params",
			baseUrl: "https://www.oklink.com",
			apiKey:  "test-api-key",
			timeout: 15 * time.Second,
			wantErr: false,
		},
		{
			name:    "create client with empty api key",
			baseUrl: "https://www.oklink.com",
			apiKey:  "",
			timeout: 15 * time.Second,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewTronDataClient(tt.baseUrl, tt.apiKey, tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTronDataClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewTronDataClient() returned nil client")
			}
		})
	}
}

func TestTronData_GetTransactionsByAddress(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": [{
				"page": "1",
				"limit": "10",
				"totalPage": "1",
				"transactionList": [
					{
						"txId": "test-tx-id-1",
						"blockHeight": "12345",
						"transactionTime": "1234567890000"
					}
				]
			}]
		}`))
	}))
	defer server.Close()

	client, err := NewTronDataClient(server.URL, "test-key", 15*time.Second)
	if err != nil {
		t.Fatalf("NewTronDataClient() error = %v", err)
	}

	tests := []struct {
		name     string
		address  string
		page     int
		pageSize int
		wantErr  bool
	}{
		{
			name:     "get transactions for valid address",
			address:  "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			page:     1,
			pageSize: 10,
			wantErr:  false,
		},
		{
			name:     "get transactions with zero page size",
			address:  "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			page:     1,
			pageSize: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txs, err := client.GetTransactionsByAddress(tt.address, tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && txs == nil {
				t.Error("GetTransactionsByAddress() returned nil")
			}
		})
	}
}

func TestTronData_GetEstimateGasFee(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"code": "0",
			"msg": "success",
			"data": [{
				"chainFullName": "Tron",
				"chainShortName": "TRON",
				"symbol": "TRX",
				"bestTransactionFee": "1000000",
				"recommendedGasPrice": "420",
				"rapidGasPrice": "500",
				"standardGasPrice": "420",
				"slowGasPrice": "300"
			}]
		}`))
	}))
	defer server.Close()

	client, err := NewTronDataClient(server.URL, "test-key", 15*time.Second)
	if err != nil {
		t.Fatalf("NewTronDataClient() error = %v", err)
	}

	gasFee, err := client.GetEstimateGasFee()
	if err != nil {
		t.Errorf("GetEstimateGasFee() error = %v", err)
		return
	}
	if gasFee == nil {
		t.Error("GetEstimateGasFee() returned nil")
	}
}
