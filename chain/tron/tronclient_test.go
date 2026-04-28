package tron

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDialTronClient(t *testing.T) {
	tests := []struct {
		name    string
		rpcURL  string
		rpcUser string
		rpcPass string
	}{
		{
			name:    "create client with credentials",
			rpcURL:  "https://api.trongrid.io",
			rpcUser: "test-api-key",
			rpcPass: "test-pass",
		},
		{
			name:    "create client without credentials",
			rpcURL:  "https://api.trongrid.io",
			rpcUser: "",
			rpcPass: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := DialTronClient(tt.rpcURL, tt.rpcUser, tt.rpcPass)
			if client == nil {
				t.Error("DialTronClient() returned nil")
			}
			if client.rpc == nil {
				t.Error("DialTronClient() rpc client is nil")
			}
		})
	}
}

func TestTronClient_GetBlockByNumber(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/walletsolidity/getblock" {
			t.Errorf("Expected path '/walletsolidity/getblock', got %s", r.URL.Path)
		}

		response := BlockResponse{
			BlockID: "0000000002b6c78800e8c8e8c8e8c8e8c8e8c8e8c8e8c8e8c8e8c8e8c8e8c8e8",
			BlockHeader: BlockHeader{
				RawData: BlockHeaderRaw{
					Number:    45678472,
					Timestamp: 1234567890000,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	tests := []struct {
		name        string
		blockNumber interface{}
		wantErr     bool
	}{
		{
			name:        "get block by number (int64)",
			blockNumber: int64(45678472),
			wantErr:     false,
		},
		{
			name:        "get block by string 'latest'",
			blockNumber: "latest",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := client.GetBlockByNumber(tt.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && block == nil {
				t.Error("GetBlockByNumber() returned nil block")
			}
		})
	}
}

func TestTronClient_GetBalance(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wallet/getaccount" {
			t.Errorf("Expected path '/wallet/getaccount', got %s", r.URL.Path)
		}

		response := Account{
			Address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			Balance: 1000000000,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "get balance for valid address",
			address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := client.GetBalance(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && account == nil {
				t.Error("GetBalance() returned nil account")
			}
			if !tt.wantErr && account.Balance == 0 {
				t.Error("GetBalance() returned zero balance")
			}
		})
	}
}

func TestTronClient_GetTransactionByHash(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/walletsolidity/gettransactionbyid" {
			t.Errorf("Expected path '/walletsolidity/gettransactionbyid', got %s", r.URL.Path)
		}

		response := Transaction{
			TxID: "test-tx-id",
			RawData: TxRawData{
				Contract: []Contract{
					{
						Type: "TransferContract",
						Parameter: Parameter{
							Value: ContractValue{
								OwnerAddress: "41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
								ToAddress:    "41e552f6487585c2b58bc2c9bb4492bc1f17132cd0",
								Amount:       1000000,
							},
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	tests := []struct {
		name    string
		txHash  string
		wantErr bool
	}{
		{
			name:    "get transaction by valid hash",
			txHash:  "test-tx-id",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := client.GetTransactionByHash(tt.txHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tx == nil {
				t.Error("GetTransactionByHash() returned nil transaction")
			}
		})
	}
}

func TestTronClient_JsonRpcBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := BlockResponse{
			BlockID: "test-block-id",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	tests := []struct {
		name    string
		params  interface{}
		wantErr bool
	}{
		{
			name:    "valid int64 param",
			params:  int64(12345),
			wantErr: false,
		},
		{
			name:    "valid string param",
			params:  "latest",
			wantErr: false,
		},
		{
			name:    "invalid param type",
			params:  12.34,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result BlockResponse
			err := client.JsonRpcBlock(tt.params, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonRpcBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTronClient_GetBlockHeaderByNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := BlockResponse{
			BlockID: "test-block-id",
			BlockHeader: BlockHeader{
				RawData: BlockHeaderRaw{
					Number: 12345,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	tests := []struct {
		name        string
		blockNumber int64
		wantErr     bool
	}{
		{
			name:        "valid block number",
			blockNumber: 12345,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := client.GetBlockHeaderByNumber(tt.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockHeaderByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && block == nil {
				t.Error("GetBlockHeaderByNumber() returned nil")
			}
		})
	}
}

func TestTronClient_Solidity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	var result map[string]interface{}
	err := client.Solidity("testmethod", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Errorf("Solidity() error = %v", err)
	}
}

func TestTronClient_Wallet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	client := DialTronClient(server.URL, "", "")

	var result map[string]interface{}
	err := client.Wallet("testmethod", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Errorf("Wallet() error = %v", err)
	}
}
