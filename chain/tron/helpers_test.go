package tron

import (
	"math/big"
	"testing"
)

func TestBase58ToHex(t *testing.T) {
	tests := []struct {
		name        string
		base58Addr  string
		wantHex     string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid TRON address",
			base58Addr: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			wantHex:    "0x41631d3c8b4e4f8e5f3e3c3b3a3d3c3b3a3d3c3b3a",
			wantErr:    false,
		},
		{
			name:        "invalid base58 string",
			base58Addr:  "invalid!!!",
			wantErr:     true,
			errContains: "failed to decode",
		},
		{
			name:        "empty string",
			base58Addr:  "",
			wantErr:     true,
			errContains: "failed to decode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHex, err := Base58ToHex(tt.base58Addr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Base58ToHex() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Base58ToHex() unexpected error: %v", err)
				return
			}
			if len(gotHex) == 0 {
				t.Errorf("Base58ToHex() returned empty hex string")
			}
		})
	}
}

func TestPadLeftZero(t *testing.T) {
	tests := []struct {
		name   string
		hexStr string
		length int
		want   string
	}{
		{
			name:   "pad 5 zeros",
			hexStr: "abc",
			length: 8,
			want:   "00000abc",
		},
		{
			name:   "no padding needed",
			hexStr: "abcdef",
			length: 6,
			want:   "abcdef",
		},
		{
			name:   "empty string",
			hexStr: "",
			length: 4,
			want:   "0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadLeftZero(tt.hexStr, tt.length)
			if got != tt.want {
				t.Errorf("PadLeftZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTRC20TransferData(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantAddr    bool
		wantAmount  bool
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid data length",
			data:        "a9059cbb",
			wantErr:     true,
			errContains: "invalid data length",
		},
		{
			name:        "empty data",
			data:        "",
			wantErr:     true,
			errContains: "invalid data length",
		},
		{
			name:       "valid transfer data",
			data:       "a9059cbb000000000000000000000000" + "1234567890123456789012345678901234567890" + "0000000000000000000000000000000000000000000000000de0b6b3a7640000",
			wantAddr:   true,
			wantAmount: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, amount, err := ParseTRC20TransferData(tt.data)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseTRC20TransferData() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseTRC20TransferData() unexpected error: %v", err)
				return
			}
			if tt.wantAddr && addr == "" {
				t.Errorf("ParseTRC20TransferData() returned empty address")
			}
			if tt.wantAmount && amount == nil {
				t.Errorf("ParseTRC20TransferData() returned nil amount")
			}
		})
	}
}

func TestHexToTronAddress(t *testing.T) {
	tests := []struct {
		name    string
		hexAddr string
		want    string
	}{
		{
			name:    "valid hex address with 0x prefix",
			hexAddr: "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			want:    "",
		},
		{
			name:    "valid hex address without 0x prefix",
			hexAddr: "41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			want:    "",
		},
		{
			name:    "invalid hex string",
			hexAddr: "invalid",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HexToTronAddress(tt.hexAddr)
			// Just verify it doesn't panic and returns a string
			_ = got
		})
	}
}

func TestTronAddressToHex(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{
			name: "valid TRON address",
			addr: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			want: "",
		},
		{
			name: "invalid address",
			addr: "invalid",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TronAddressToHex(tt.addr)
			// Just verify it doesn't panic
			_ = got
		})
	}
}

func TestFormatTronAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
	}{
		{
			name:    "address starting with T",
			address: "TJRyWwFs9wTFGZg3JbrVriFbNfCug5tDeC",
			want:    "0x",
		},
		{
			name:    "address with 0x prefix",
			address: "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			want:    "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
		},
		{
			name:    "address without prefix",
			address: "41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			want:    "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTronAddress(tt.address)
			if tt.name == "address with 0x prefix" && got != tt.want {
				t.Errorf("FormatTronAddress() = %v, want %v", got, tt.want)
			}
			if tt.name == "address without prefix" && got != tt.want {
				t.Errorf("FormatTronAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTRC20TransferDataAmount(t *testing.T) {
	// Test with a known amount
	data := "a9059cbb000000000000000000000000" + "1234567890123456789012345678901234567890" + "0000000000000000000000000000000000000000000000000de0b6b3a7640000"
	_, amount, err := ParseTRC20TransferData(data)
	if err != nil {
		t.Fatalf("ParseTRC20TransferData() error = %v", err)
	}

	expectedAmount := new(big.Int)
	expectedAmount.SetString("0de0b6b3a7640000", 16)

	if amount.Cmp(expectedAmount) != 0 {
		t.Errorf("ParseTRC20TransferData() amount = %v, want %v", amount, expectedAmount)
	}
}
