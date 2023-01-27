package v1Proxy

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IV1Proxy interface {
	GetChains(request *pb.BridgeChainRequest, proxyUrl string) (*pb.BridgeChainResponse, error)
	GetChainTokens(request *pb.BridgeChainTokensRequest, proxyUrl string) (*pb.BridgeChainTokensResponse, error)
	GetQuote(request *pb.BridgeQuoteRequest, proxyUrl string) (*pb.BridgeQuoteResponse, error)
	GetTransaction(request *pb.BridgeTransactionRequest, proxyUrl string) (*pb.BridgeTransactionResponse, error)
	GetTransactionStatus(request *pb.BridgeTransactionStatusRequest, proxyUrl string) (*pb.BridgeTransactionStatusResponse, error)
}
type V1Proxy struct {
	httpRequest utils.IHttpRequest
}

func NewV1Proxy(httpRequest utils.IHttpRequest) *V1Proxy {
	return &V1Proxy{
		httpRequest: httpRequest,
	}
}

func (v *V1Proxy) GetChains(request *pb.BridgeChainRequest, proxyUrl string) (*pb.BridgeChainResponse, error) {
	url := fmt.Sprintf("%s/v1/bridge/chains?bridgeProvider=%s", proxyUrl, request.BridgeProvider)
	body, err := v.httpRequest.GetRequest(url)
	if err != nil {
		var v1ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &v1ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), v1ErrorResponse.Message)
	}
	var v1Response BridgeChains
	err = json.Unmarshal(body, &v1Response)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}

	var response pb.BridgeChainResponse
	if v1Response.Chains != nil {
		if len(v1Response.Chains) > 0 {
			for _, item := range v1Response.Chains {
				response.Chains = append(response.Chains, &pb.BridgeChainInfo{
					Symbol:    item.Symbol,
					ChainType: item.ChainType,
					Name:      item.Name,
					Coin:      item.Coin,
					ChainId:   item.ChainID,
					LogoUrl:   item.LogoURL,
					MainNet:   item.Mainnet,
				})
			}
		}
	}
	return &response, err
}

func (v *V1Proxy) GetChainTokens(request *pb.BridgeChainTokensRequest, proxyUrl string) (*pb.BridgeChainTokensResponse, error) {
	if request.FromChain == "" {
		request.FromChain = "0"
	}
	if request.ToChain == "" {
		request.ToChain = "0"
	}
	url := fmt.Sprintf("%s/v1/bridge/chains/tokens?fromChain=%s&toChain=%s&fromToken=%s&bridgeProvider=%s", proxyUrl, request.FromChain, request.ToChain, request.FromToken, request.BridgeProvider)
	body, err := v.httpRequest.GetRequest(url)
	if err != nil {
		var v1ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &v1ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), v1ErrorResponse.Message)
	}
	var v1Response BridgeExchangeTokens
	err = json.Unmarshal(body, &v1Response)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}

	var response pb.BridgeChainTokensResponse
	if v1Response.Tokens != nil {
		if len(v1Response.Tokens) > 0 {
			for _, item := range v1Response.Tokens {
				response.Tokens = append(response.Tokens, &pb.BridgeTokens{
					TokenAddress:  item.TokenAddress,
					TokenName:     item.TokenName,
					TokenSymbol:   item.TokenSymbol,
					TokenDecimals: item.TokenDecimals,
					TokenLogoUrl:  item.TokenLogoURL,
				})
			}
		}
	}
	return &response, err
}

func (v *V1Proxy) GetQuote(request *pb.BridgeQuoteRequest, proxyUrl string) (*pb.BridgeQuoteResponse, error) {
	url := fmt.Sprintf("%s/v1/bridge/quote?fromChain=%s&fromToken=%s&toChain=%s&toToken=%s&fromAmount=%s&fromAddress=%s&toAddress=%s&bridgeProvider=%s", proxyUrl, request.FromChain, request.FromToken, request.ToChain, request.ToToken, request.FromAmount, request.FromAddress, request.ToAddress, request.BridgeProvider)
	body, err := v.httpRequest.GetRequest(url)
	if err != nil {
		var v1ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &v1ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), v1ErrorResponse.Message)
	}
	var successError SuccessErrorV1
	err = json.Unmarshal(body, &successError)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	if successError.Error.Message != "" {
		return nil, status.Errorf(codes.NotFound, string(body), successError.Error.Message)
	}
	var v1Response BridgeQuote
	jsonErr := json.Unmarshal(body, &v1Response)
	if jsonErr != nil {
		return nil, status.Errorf(codes.Internal, jsonErr.Error(), "Json Unmarshal Error")
	}
	var quoteResponse pb.BridgeQuoteResponse
	var toolDetails pb.ToolDetails
	var estimate pb.Estimate
	var bridgeFee pb.BridgeFee
	quoteResponse.Tool = v1Response.Tool
	toolDetails.Key = v1Response.ToolDetails.Key
	toolDetails.Name = v1Response.ToolDetails.Name
	toolDetails.LogoUrl = v1Response.ToolDetails.LogoURL
	estimate.FromAmount = v1Response.Estimate.FromAmount
	estimate.FromTokenDecimals = v1Response.Estimate.FromTokenDecimals
	estimate.ToAmount = v1Response.Estimate.ToAmount
	estimate.ToAmountMin = v1Response.Estimate.ToAmountMin
	estimate.ToTokenDecimals = v1Response.Estimate.ToTokenDecimals
	estimate.ExecutionDuration = v1Response.Estimate.ExecutionDuration
	estimate.FromAmountUsd = v1Response.Estimate.FromAmountUsd
	estimate.ToAmountUsd = v1Response.Estimate.ToAmountUsd
	estimate.ApproveAddress = v1Response.Estimate.ApproveAddress
	if v1Response.BridgeFee.ContractAddress != "" {
		bridgeFee.Symbol = v1Response.BridgeFee.Symbol
		bridgeFee.TokenDecimals = v1Response.BridgeFee.TokenDecimals
		bridgeFee.ContractAddress = v1Response.BridgeFee.ContractAddress
		bridgeFee.Amount = v1Response.BridgeFee.Amount
		bridgeFee.AmountUsd = v1Response.BridgeFee.AmountUsd
		quoteResponse.BridgeFee = &bridgeFee
	}
	quoteResponse.ToolDetails = &toolDetails
	quoteResponse.Estimate = &estimate
	return &quoteResponse, err

}

func (v *V1Proxy) GetTransaction(request *pb.BridgeTransactionRequest, proxyUrl string) (*pb.BridgeTransactionResponse, error) {
	url := fmt.Sprintf("%s/v1/bridge/transaction?fromChain=%s&fromToken=%s&toChain=%s&toToken=%s&fromAmount=%s&fromAddress=%s&toAddress=%s&bridgeProvider=%s", proxyUrl, request.FromChain, request.FromToken, request.ToChain, request.ToToken, request.FromAmount, request.FromAddress, request.ToAddress, request.BridgeProvider)
	body, err := v.httpRequest.GetRequest(url)
	if err != nil {
		var v1ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &v1ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), v1ErrorResponse.Message)
	}
	var v1Response BridgeTransaction
	err = json.Unmarshal(body, &v1Response)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}

	var transactionResponse pb.BridgeTransactionResponse
	var toolDetails pb.ToolDetails
	var transactionRequest pb.TransactionRequest
	transactionResponse.Tool = v1Response.Tool
	toolDetails.Key = v1Response.ToolDetails.Key
	toolDetails.Name = v1Response.ToolDetails.Name
	toolDetails.LogoUrl = v1Response.ToolDetails.LogoURL
	transactionRequest.Data = v1Response.TransactionRequest.Data
	transactionRequest.To = v1Response.TransactionRequest.To
	transactionRequest.Value = v1Response.TransactionRequest.Value
	transactionRequest.GasLimit = v1Response.TransactionRequest.GasLimit
	transactionResponse.ToolDetails = &toolDetails
	transactionResponse.TransactionRequest = &transactionRequest
	return &transactionResponse, err
}

func (v *V1Proxy) GetTransactionStatus(request *pb.BridgeTransactionStatusRequest, proxyUrl string) (*pb.BridgeTransactionStatusResponse, error) {
	url := fmt.Sprintf("%s/v1/bridge/status?bridge=%s&fromChain=%s&toChain=%s&txHash=%s&bridgeProvider=%s", proxyUrl, request.Bridge, request.FromChain, request.ToChain, request.TxHash, request.BridgeProvider)
	body, err := v.httpRequest.GetRequest(url)
	if err != nil {
		var v1ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &v1ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), v1ErrorResponse.Message)
	}
	var v1Response BridgeTransactionStatus
	err = json.Unmarshal(body, &v1Response)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var response pb.BridgeTransactionStatusResponse
	response.Status = v1Response.Status
	response.TxHash = v1Response.TxHash
	response.Msg = v1Response.Msg
	response.IsSuccess = v1Response.IsSuccess
	return &response, err
}
