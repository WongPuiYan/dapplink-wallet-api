package tron

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dapplink-labs/dapplink-wallet-api/config"
	wallet_api "github.com/dapplink-labs/dapplink-wallet-api/protobuf/wallet-api"
)

func TestChainAdaptor_ConvertAddresses(t *testing.T) {
	// Create a mock server for TronData
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     "https://api.trongrid.io",
				RpcUser:    "",
				RpcPass:    "",
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	tests := []struct {
		name    string
		req     *wallet_api.ConvertAddressesRequest
		wantErr bool
	}{
		{
			name: "convert valid public key",
			req: &wallet_api.ConvertAddressesRequest{
				PublicKey: []*wallet_api.PublicKey{
					{
						PublicKey: "0x04a614f803b6fd780986a42c78ec9c7f77e6ded13c6e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e9e",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "convert invalid public key",
			req: &wallet_api.ConvertAddressesRequest{
				PublicKey: []*wallet_api.PublicKey{
					{
						PublicKey: "invalid",
					},
				},
			},
			wantErr: false, // Should return empty address, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adaptor.ConvertAddresses(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Error("ConvertAddresses() returned nil response")
			}
		})
	}
}

func TestChainAdaptor_ValidAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     "https://api.trongrid.io",
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	tests := []struct {
		name    string
		req     *wallet_api.ValidAddressesRequest
		wantErr bool
	}{
		{
			name: "validate valid address",
			req: &wallet_api.ValidAddressesRequest{
				Addresses: []*wallet_api.Addresses{
					{
						Address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "validate invalid address",
			req: &wallet_api.ValidAddressesRequest{
				Addresses: []*wallet_api.Addresses{
					{
						Address: "invalid",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adaptor.ValidAddresses(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Error("ValidAddresses() returned nil response")
			}
		})
	}
}

func TestChainAdaptor_GetLastestBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/walletsolidity/getblock" {
			response := BlockResponse{
				BlockID: "test-block-id",
				BlockHeader: BlockHeader{
					RawData: BlockHeaderRaw{
						Number:    12345,
						Timestamp: time.Now().Unix(),
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	req := &wallet_api.LastestBlockRequest{}

	resp, err := adaptor.GetLastestBlock(context.Background(), req)
	if err != nil {
		t.Errorf("GetLastestBlock() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("GetLastestBlock() returned nil response")
	}
}

func TestChainAdaptor_GetBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/walletsolidity/getblock" {
			response := BlockResponse{
				BlockID: "test-block-id",
				BlockHeader: BlockHeader{
					RawData: BlockHeaderRaw{
						Number:    12345,
						Timestamp: time.Now().Unix(),
					},
				},
				Transactions: []Transaction{
					{
						TxID: "test-tx-id",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	req := &wallet_api.BlockRequest{
		HashHeight: "12345",
	}

	resp, err := adaptor.GetBlock(context.Background(), req)
	if err != nil {
		t.Errorf("GetBlock() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("GetBlock() returned nil response")
	}
}

func TestChainAdaptor_GetTransactionByHash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/walletsolidity/gettransactionbyid" {
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
		}
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	req := &wallet_api.TransactionByHashRequest{
		Hash: "test-tx-id",
	}

	resp, err := adaptor.GetTransactionByHash(context.Background(), req)
	if err != nil {
		t.Errorf("GetTransactionByHash() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("GetTransactionByHash() returned nil response")
	}
}

func TestChainAdaptor_GetAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/wallet/getaccount" {
			response := Account{
				Address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
				Balance: 1000000000,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	tests := []struct {
		name    string
		req     *wallet_api.AccountBalanceRequest
		wantErr bool
	}{
		{
			name: "get native TRX balance",
			req: &wallet_api.AccountBalanceRequest{
				Address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			},
			wantErr: false,
		},
		{
			name: "get TRC20 token balance",
			req: &wallet_api.AccountBalanceRequest{
				Address:         "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
				ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := adaptor.GetAccountBalance(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Error("GetAccountBalance() returned nil response")
			}
		})
	}
}

func TestChainAdaptor_SendTransaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	req := &wallet_api.SendTransactionsRequest{
		RawTx: []*wallet_api.RawTransaction{
			{
				RawTx: "test-raw-tx",
			},
		},
	}

	resp, err := adaptor.SendTransaction(context.Background(), req)
	if err != nil {
		t.Errorf("SendTransaction() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("SendTransaction() returned nil response")
	}
}

func TestChainAdaptor_BuildTransactionSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := &config.Config{
		WalletNode: config.WalletNode{
			Tron: config.Node{
				RpcUrl:     server.URL,
				DataApiUrl: server.URL,
				DataApiKey: "test-key",
			},
		},
	}

	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		t.Fatalf("NewChainAdaptor() error = %v", err)
	}

	req := &wallet_api.TransactionSchemaRequest{}

	resp, err := adaptor.BuildTransactionSchema(context.Background(), req)
	if err != nil {
		t.Errorf("BuildTransactionSchema() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("BuildTransactionSchema() returned nil response")
	}
}

func TestHexToTronAddressInTron(t *testing.T) {
	tests := []struct {
		name    string
		hexAddr string
		want    string
	}{
		{
			name:    "valid hex address",
			hexAddr: "41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HexToTronAddress(tt.hexAddr)
			if got == "" {
				t.Error("HexToTronAddress() returned empty string")
			}
		})
	}
}
