package tron

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	base582 "github.com/btcsuite/btcutil/base58"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/mr-tron/base58"
)

const (
	AddressPrefix = "41"
)

// Base58ToHex Convert TRON address from base58 to hexadecimal

// Base58ToHex Convert TRON address from base58 to hexadecimal
func Base58ToHex(base58Addr string) (string, error) {
	// Decode base58 address
	dec, err := base58.Decode(base58Addr)
	if err != nil {
		return "", fmt.Errorf("failed to decode base58 address: %w", err)
	}

	// Check if decoded length is 25 bytes
	if len(dec) != 25 {
		return "", fmt.Errorf("invalid address length: expected 25, got %d", len(dec))
	}

	// Extract initial address (first 21 bytes)
	initialAddress := dec[:21]

	// Calculate verification code
	expectedVerificationCode := make([]byte, 4)
	hash := sha256.Sum256(initialAddress)
	hash2 := sha256.Sum256(hash[:])
	copy(expectedVerificationCode, hash2[:4])

	// Verify verification code
	if !bytes.Equal(dec[21:], expectedVerificationCode) {
		return "", fmt.Errorf("invalid verification code")
	}

	// Convert initial address to hex string with "0x" prefix
	hexAddress := "0x" + hex.EncodeToString(initialAddress)
	return hexAddress, nil
}

// PadLeftZero Fill the left side of the hexadecimal string with zero to the specified length
func PadLeftZero(hexStr string, length int) string {
	return strings.Repeat("0", length-len(hexStr)) + hexStr
}

// ParseTRC20TransferData Extract the 'to' address and 'amount' from ABI encoded data
func ParseTRC20TransferData(data string) (string, *big.Int, error) {
	// Check minimum data length (method signature 8 chars + address 64 chars + amount 64 chars)
	if len(data) < 136 {
		return "", nil, fmt.Errorf("invalid data length: expected at least 136, got %d", len(data))
	}

	// Extract the receiving address (positions 32 to 72)
	toAddressHex := data[32:72]
	toAddress, err := address.HexToAddress(AddressPrefix + toAddressHex)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse address: %w", err)
	}

	// Get amount (positions 72 to 136)
	valueHex := data[72:136]
	value := new(big.Int)
	_, ok := value.SetString(valueHex, 16)
	if !ok {
		return "", nil, fmt.Errorf("failed to parse amount from hex: %s", valueHex)
	}

	return toAddress.String(), value, nil
}

// Helper functions
func HexToTronAddress(hexAddr string) string {
	hexAddr = strings.TrimPrefix(hexAddr, "0x")
	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return ""
	}
	return base582.CheckEncode(addrBytes[1:], addrBytes[0])
}

func TronAddressToHex(addr string) string {
	decoded, version, err := base582.CheckDecode(addr)
	if err != nil {
		return ""
	}
	return "0x" + hex.EncodeToString(append([]byte{version}, decoded...))
}

func FormatTronAddress(address string) string {
	if strings.HasPrefix(address, "T") {
		return "0x" + hex.EncodeToString(base582.Decode(address))
	}
	if !strings.HasPrefix(address, "0x") {
		return "0x" + address
	}
	return address
}
