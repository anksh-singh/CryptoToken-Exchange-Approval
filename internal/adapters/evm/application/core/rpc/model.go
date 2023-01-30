package rpc

import (
	// "time"
)
type BlockNativeGasPrice struct {
	System             string `json:"system"`
	Network            string `json:"network"`
	Unit               string `json:"unit"`
	MaxPrice           int    `json:"maxPrice"`
	CurrentBlockNumber int    `json:"currentBlockNumber"`
	MsSinceLastBlock   int    `json:"msSinceLastBlock"`
	BlockPrices        []struct {
		BlockNumber               int     `json:"blockNumber"`
		EstimatedTransactionCount int     `json:"estimatedTransactionCount"`
		BaseFeePerGas             float64 `json:"baseFeePerGas"`
		EstimatedPrices           []struct {
			Confidence           float64     `json:"confidence"`
			Price                interface{} `json:"price"`
			MaxPriorityFeePerGas float64     `json:"maxPriorityFeePerGas"`
			MaxFeePerGas         float64     `json:"maxFeePerGas"`
		} `json:"estimatedPrices"`
	} `json:"blockPrices"`
	EstimatedBaseFees []struct {
		Pending1 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+1,omitempty"`
		Pending2 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+2,omitempty"`
		Pending3 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+3,omitempty"`
		Pending4 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+4,omitempty"`
		Pending5 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+5,omitempty"`
	} `json:"estimatedBaseFees"`
}

type TokenAllowance struct {
	Contract string `json:"contract"`
	Owner    string `json:"owner"`
	Spender  string `json:"spender"`
}

type AllowanceRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type ContractABIRequest struct {
	Contract string `json:"contract"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    int64  `json:"value"`
	Data     string `json:"data"`
	Chain    string `json:"chain"`
	Method   string `json:"method"`
}
