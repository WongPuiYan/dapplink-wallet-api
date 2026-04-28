package tron

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultRequestTimeout = 10 * time.Second
	defaultRetryCount     = 3
)

// TronClient Define a Tron RPC client
type TronClient struct {
	rpc *resty.Client
}

// DialTronClient Initialize and return a TronClient instance
func DialTronClient(rpcURL, rpcUser, rpcPass string) *TronClient {
	client := resty.New()
	if rpcUser != "" && rpcPass != "" {
		client.SetHeader("TRON-PRO-API-KEY", rpcPass)
	}
	client.SetBaseURL(rpcURL)
	client.SetTimeout(defaultRequestTimeout)
	client.SetRetryCount(defaultRetryCount)

	return &TronClient{
		rpc: client,
	}
}

func (client *TronClient) JsonRpcBlock(params interface{}, result interface{}) error {
	var idOrNum string
	switch v := params.(type) {
	case int64:
		idOrNum = fmt.Sprintf("\"%d\"", v)
	case string:
		idOrNum = fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Errorf("unsupported params type: %T", params)
	}

	requestBody := map[string]interface{}{
		"id_or_num": json.RawMessage(idOrNum),
		"detail":    true,
	}

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/walletsolidity/getblock")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (client *TronClient) JsonRpcBlockHeader(params interface{}, result interface{}) error {
	var idOrNum string
	switch v := params.(type) {
	case int64:
		idOrNum = fmt.Sprintf("\"%d\"", v)
	case string:
		idOrNum = fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Errorf("unsupported params type: %T", params)
	}

	requestBody := map[string]interface{}{
		"id_or_num": json.RawMessage(idOrNum),
		"detail":    false,
	}

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/wallet/getblock")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// Solidity Call Solidity
func (client *TronClient) Solidity(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/walletsolidity/" + method)
	return err
}

// Wallet Call Wallet
func (client *TronClient) Wallet(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/wallet/" + method)
	return err
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockByNumber(blockNumber interface{}) (*BlockResponse, error) {

	var response BlockResponse
	err := client.JsonRpcBlock(blockNumber, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}
	return &response, nil
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockByHush(hush string) (*Block, error) {
	params := []interface{}{hush, true}
	var response Response[Block]
	err := client.JsonRpcBlock(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}
	return &response.Result, nil
}

// GetBlockHeaderByNumber 获取区块头信息
func (client *TronClient) GetBlockHeaderByNumber(blockNumber int64) (*BlockResponse, error) {
	var response BlockResponse
	err := client.JsonRpcBlockHeader(blockNumber, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: %v", err)
	}

	// 检查响应是否为空
	if response.BlockID == "" {
		return nil, fmt.Errorf("empty response received")
	}

	return &response, nil
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockHeaderByHash(blockHush string) (*BlockResponse, error) {
	var response BlockResponse
	err := client.JsonRpcBlockHeader(blockHush, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: %v", err)
	}

	// 检查响应是否为空
	if response.BlockID == "" {
		return nil, fmt.Errorf("empty response received")
	}

	return &response, nil
}

// GetBlockByHash Obtain block information based on block hash
func (client *TronClient) GetBlockByHash(blockHash string) (*Block, error) {
	params := []interface{}{blockHash, false}
	var response Response[Block]
	err := client.JsonRpcGetBlockByHash(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response.Result, nil
}
func (client *TronClient) JsonRpcGetBlockByHash(params interface{}, result interface{}) error {
	return nil
}
func (client *TronClient) JsonRpcGetBalance(params interface{}, result interface{}) error {
	requestBody := map[string]interface{}{
		"address": params,
		"visible": true,
	}

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/wallet/getaccount")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// GetAccount Get account information
func (client *TronClient) GetBalance(address string) (*Account, error) {
	params := []interface{}{address}
	var response Account
	err := client.JsonRpcGetBalance(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response, nil

}

// GetAccount Get account information
func (client *TronClient) GetTransactionByHash(hush string) (*Transaction, error) {
	params := []interface{}{hush}
	var response Transaction
	err := client.JsonRpcGetTransactionByHash(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response, nil

}

func (client *TronClient) JsonRpcGetTransactionByHash(params interface{}, result interface{}) error {
	requestBody := map[string]interface{}{
		"value": params,
	}

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/walletsolidity/gettransactionbyid")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
