package lifi

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type ILiFi interface {
	GetChains() (*pb.BridgeChainResponse, error)
	GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error)
	GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error)
	GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error)
	// GetExchangeTokens Swap function
	GetExchangeTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error)
	// GetLiFiQuote Swap function
	GetLiFiQuote(SellToken string, BuyToken string, SellAmount string, chain string, srcTokenDecimals string, dstTokenDecimals string, takerAddress string, slippage string) (*pb.ExchangeQuoteResponse, error)
	// GetLiFiSwap Swap function
	GetLiFiSwap(request *pb.ExchangeSwapRequest, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error)
}
type LiFi struct {
	env            *config.Config
	logger         *zap.SugaredLogger
	httpRequest    utils.IHttpRequest
	helper         *utils.Helpers
	allowedBridges string
	lifiChainList  []*config.BridgeNomenclature
	utils          *utils.UtilConf
}

func NewLiFiService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, lifiChainList []*config.BridgeNomenclature, utils *utils.UtilConf) *LiFi {
	return &LiFi{
		env:            env,
		logger:         logger,
		helper:         helper,
		httpRequest:    httpRequest,
		allowedBridges: "connext",
		lifiChainList:  lifiChainList,
		utils:          utils,
	}
}
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
func (l *LiFi) GetChains() (*pb.BridgeChainResponse, error) {
	url := l.env.Swap.LIFIEndpoint + "/chains"
	body, err := l.httpRequest.GetRequest(url)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
	}
	var chainInfo ChainsInfo
	err = json.Unmarshal(body, &chainInfo)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var chainResponse pb.BridgeChainResponse

	for _, item := range chainInfo.Chains {
		info := l.utils.GetWalletInfo(strings.ToLower(item.Name))
		if info.ChainID == item.ID {
			if item.Name != "" {
				chainResponse.Chains = append(chainResponse.Chains, &pb.BridgeChainInfo{
					Symbol:    item.Key,
					ChainType: item.ChainType,
					Name:      item.Name,
					Coin:      item.Coin,
					ChainId:   int64(item.ID),
					LogoUrl:   item.LogoURI,
					MainNet:   item.Mainnet,
				})
			}
		}
	}
	return &chainResponse, err
}

func (l *LiFi) GetExchangeTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	// Validate Chain
	chainKey := l.FindChainKey(request.Chain)
	url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/tokens?chains=%s", chainKey)

	body, err := l.httpRequest.GetRequest(url)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	tokenLifi := result["tokens"].(map[string]interface{})
	var chainTokensResponse pb.BridgeChainTokensResponse
	if len(tokenLifi) > 0 {
		tokens := result["tokens"].(map[string]interface{})[chainKey].([]interface{})
		var nativeTokenInfo pb.BridgeTokens
		if len(tokens) > 0 {
			//tokens = tokens[0:5]
			for _, item := range tokens {
				inItem := item.(map[string]interface{})
				tokenAddress := l.GetFrontierSpecificNativeToken(request.Chain, inItem["address"].(string))
				if inItem["address"].(string) == tokenAddress {
					chainTokensResponse.Tokens = append(chainTokensResponse.Tokens, &pb.BridgeTokens{
						TokenAddress:  tokenAddress,
						TokenDecimals: int64(inItem["decimals"].(float64)),
						TokenSymbol:   inItem["symbol"].(string),
						TokenName:     inItem["name"].(string),
						TokenLogoUrl:  l.CheckNilForString(inItem["logoURI"]).(string),
					})
				} else {
					nativeTokenInfo = pb.BridgeTokens{
						TokenAddress:  tokenAddress,
						TokenDecimals: int64(inItem["decimals"].(float64)),
						TokenSymbol:   inItem["symbol"].(string),
						TokenName:     inItem["name"].(string),
						TokenLogoUrl:  l.CheckNilForString(inItem["logoURI"]).(string),
					}
				}

			}
			chainTokensResponse.Tokens = append([]*pb.BridgeTokens{
				&nativeTokenInfo,
			}, chainTokensResponse.Tokens...)
		}
	} else {
		return nil, status.Errorf(codes.NotFound, errors.New("no tokens for the chain").Error(), "no tokens for the chain "+request.Chain)
	}

	return &chainTokensResponse, err
}

func (l *LiFi) CheckNilForString(v interface{}) interface{} {
	if v == nil {
		return ""
	} else {
		return v
	}
}

func (l *LiFi) GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	// Validate Chain
	if request.FromChain != "" && request.ToChain != "" && request.FromToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "From token cannot be empty for this case")
	}
	IsValidateChain, fromChain, toChain := l.ValidateChain(request.FromChain, request.ToChain)
	if IsValidateChain {
		request.FromToken = l.GetLiFiSpecificNativeToken(request.FromChain, request.FromToken)

		url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/connections?fromChain=%s&fromToken=%s&toChain=%s", fromChain, request.FromToken, toChain)
		body, err := l.httpRequest.GetRequest(url)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
		}
		var connections ChainsToken
		err = json.Unmarshal(body, &connections)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		var chainTokensResponse pb.BridgeChainTokensResponse
		chainTokensResponse.Tokens = make([]*pb.BridgeTokens, 0)
		var nativeTokenInfo *pb.BridgeTokens
		if len(connections.Connections) > 0 {
			tokens := connections.Connections[0].ToTokens
			for _, item := range tokens {
				tokenAddress := l.GetFrontierSpecificNativeToken(toChain, item.Address)
				if item.Address == tokenAddress {
					chainTokensResponse.Tokens = append(chainTokensResponse.Tokens, &pb.BridgeTokens{
						TokenAddress:  tokenAddress,
						TokenDecimals: int64(item.Decimals),
						TokenSymbol:   item.Symbol,
						TokenName:     item.Name,
						TokenLogoUrl:  item.LogoURI,
					})
				} else {
					nativeTokenInfo = &pb.BridgeTokens{
						TokenAddress:  tokenAddress,
						TokenDecimals: int64(item.Decimals),
						TokenSymbol:   item.Symbol,
						TokenName:     item.Name,
						TokenLogoUrl:  item.LogoURI,
					}
				}

			}
			if nativeTokenInfo != nil {
				chainTokensResponse.Tokens = append([]*pb.BridgeTokens{
					nativeTokenInfo,
				}, chainTokensResponse.Tokens...)
			}

		}
		if len(chainTokensResponse.Tokens) > 500 && !request.FullList {
			chainTokensResponse.Tokens = chainTokensResponse.Tokens[0:500]
			l.logger.Infof("GetChainTokens: Return full list? %v", request.FullList)
		}
		l.logger.Infof("GetChainTokens: Returning full list? %v", request.FullList)
		return &chainTokensResponse, err
	} else {
		return nil, status.Errorf(codes.Unknown, errors.New("chain not supported").Error(), "Chain not served by Frontier")
	}

}

func (l *LiFi) GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error) {
	IsValidateChain, fromChain, toChain := l.ValidateChain(request.FromChain, request.ToChain)
	if IsValidateChain {
		request.FromToken = l.GetLiFiSpecificNativeToken(request.FromChain, request.FromToken)
		request.ToToken = l.GetLiFiSpecificNativeToken(request.ToChain, request.ToToken)
		url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/quote?fromChain=%s&toChain=%s&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&fromAmount=%s",
			fromChain, toChain, request.FromToken, request.ToToken, request.FromAddress, request.ToAddress, request.FromAmount)
		body, err := l.httpRequest.GetRequest(url)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
		}

		var quote Quote
		err = json.Unmarshal(body, &quote)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "unmarshall error")
		}
		var quoteResponse pb.BridgeQuoteResponse
		var toolDetails pb.ToolDetails
		var estimate pb.Estimate
		quoteResponse.Tool = quote.Tool
		toolDetails.Key = quote.ToolDetails.Key
		toolDetails.Name = quote.ToolDetails.Name
		toolDetails.LogoUrl = quote.ToolDetails.LogoURI
		estimate.FromAmount = quote.Estimate.FromAmount
		estimate.FromTokenDecimals = int64(quote.Action.FromToken.Decimals)
		estimate.ToAmount = quote.Estimate.ToAmount
		estimate.ToAmountMin = quote.Estimate.ToAmountMin
		estimate.ToTokenDecimals = int64(quote.Action.ToToken.Decimals)
		estimate.ExecutionDuration = quote.Estimate.ExecutionDuration
		estimate.FromAmountUsd = quote.Estimate.FromAmountUSD
		estimate.ToAmountUsd = quote.Estimate.ToAmountUSD
		estimate.ApproveAddress = quote.Estimate.ApprovalAddress
		quoteResponse.ToolDetails = &toolDetails
		quoteResponse.Estimate = &estimate
		return &quoteResponse, err
	} else {
		return nil, status.Errorf(codes.Unknown, errors.New("chain not supported ").Error(), "Chain not served by Frontier")
	}

}
func (l *LiFi) GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error) {
	IsValidateChain, fromChain, toChain := l.ValidateChain(request.FromChain, request.ToChain)
	if IsValidateChain {
		request.FromToken = l.GetLiFiSpecificNativeToken(request.FromChain, request.FromToken)
		request.ToToken = l.GetLiFiSpecificNativeToken(request.ToChain, request.ToToken)
		url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/quote?fromChain=%s&toChain=%s&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&fromAmount=%s",
			fromChain, toChain, request.FromToken, request.ToToken, request.FromAddress, request.ToAddress, request.FromAmount)
		body, err := l.httpRequest.GetRequest(url)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
		}
		var quote Quote
		err = json.Unmarshal(body, &quote)
		if err != nil {
			l.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		var transactionResponse pb.BridgeTransactionResponse
		var toolDetails pb.ToolDetails
		var transactionRequest pb.TransactionRequest
		transactionResponse.Tool = quote.Tool
		toolDetails.Key = quote.ToolDetails.Key
		toolDetails.Name = quote.ToolDetails.Name
		toolDetails.LogoUrl = quote.ToolDetails.LogoURI
		transactionRequest.Data = quote.TransactionRequest.Data
		transactionRequest.To = quote.TransactionRequest.To
		transactionRequest.Value = l.helper.ConvertHexFloatString(quote.TransactionRequest.Value)
		transactionRequest.GasLimit = l.helper.ConvertHexFloatString(quote.TransactionRequest.GasLimit)
		transactionResponse.ToolDetails = &toolDetails
		transactionResponse.TransactionRequest = &transactionRequest
		return &transactionResponse, err
	} else {
		return nil, status.Errorf(codes.Unknown, errors.New("chain not supported").Error(), "Chain not served by Frontier")
	}
}

func (l *LiFi) GetFrontierSpecificNativeToken(chain string, tokenAddress string) string {
	info := l.utils.GetWalletInfo(chain)
	if tokenAddress != "0x0000000000000000000000000000000000000000" {
		return tokenAddress
	} else if strings.ToLower(tokenAddress) == "0x471ece3750da237f93b8e339c536989b8978a438" {
		return tokenAddress
	} else {
		if strings.ToLower(chain) == strings.ToLower(info.ChainName) || strings.ToLower(chain) == strings.ToLower(info.Bridge.ChainKey) || strings.ToLower(chain) == info.Bridge.ChainId {
			return info.NativeTokenInfo.Address
		}
	}
	return tokenAddress
}

func (l *LiFi) GetLiFiSpecificNativeToken(chain string, tokenAddress string) string {
	info := l.utils.GetWalletInfo(chain)
	if tokenAddress == info.NativeTokenInfo.Address && (info.ChainName == chain || info.NativeTokenInfo.ChainId == chain) {
		return "0x0000000000000000000000000000000000000000"
	} else {
		return tokenAddress
	}
}

func (l *LiFi) GetLiFiQuote(SellToken string, BuyToken string, SellAmount string, chain string, srcTokenDecimals string, dstTokenDecimals string, takerAddress string, slippage string) (*pb.ExchangeQuoteResponse, error) {
	sellAmountParam, err := l.helper.ConvertStringValueToFloatWei(SellAmount, srcTokenDecimals)
	sellAmountParam = math.Floor(sellAmountParam)
	_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)
	chainKey := l.FindChainKey(chain)
	SellToken = l.GetLiFiSpecificNativeToken(chain, SellToken)
	BuyToken = l.GetLiFiSpecificNativeToken(chain, BuyToken)
	url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/quote?fromChain=%s&toChain=%s&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&fromAmount=%s&allowBridges=%s",
		chainKey, chainKey, SellToken, BuyToken, takerAddress, takerAddress, _sellAmountParam, l.allowedBridges)
	body, err := l.httpRequest.GetRequest(url)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
	}
	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var response pb.ExchangeQuoteResponse

	resAmount, ok := new(big.Float).SetString("0")
	if quote.Estimate.ToAmount != "" {
		resAmount, ok = new(big.Float).SetString(quote.Estimate.ToAmount)
		if !ok {
			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "set string error")
		}
	}
	decimalAmount, ok := new(big.Float).SetString("0")
	if dstTokenDecimals != "" {
		decimalAmount, ok = new(big.Float).SetString(dstTokenDecimals)
		if !ok {
			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "set string error")
		}
	}

	_decimalAmount, _ := decimalAmount.Float64()
	var resPricePerFromToken float64
	var resPricePerToToken float64

	wei := math.Pow(10, _decimalAmount)
	_resAmount, _ := new(big.Float).Quo(resAmount, big.NewFloat(wei)).Float64()
	response.ResAmount = strconv.FormatFloat(_resAmount, 'f', -1, 64)
	fromTokenRate, ok := new(big.Float).SetString(quote.Action.FromToken.PriceUSD)
	if !ok {
		fromTokenRate, _ = new(big.Float).SetString("0")
	}
	toTokenRate, ok := new(big.Float).SetString(quote.Action.ToToken.PriceUSD)
	if !ok {
		toTokenRate, _ = new(big.Float).SetString("0")
	}
	if quote.Action.ToToken.PriceUSD == "0" {
		resPricePerFromToken = 0
	} else {
		resPricePerFromToken, _ = new(big.Float).Quo(fromTokenRate, toTokenRate).Float64()
	}
	if quote.Action.FromToken.PriceUSD == "0" {
		resPricePerToToken = 0
	} else {
		resPricePerToToken, _ = new(big.Float).Quo(toTokenRate, fromTokenRate).Float64()
	}

	//priceimpact
	if quote.Action.FromToken.PriceUSD == "" {
		response.ResPricePerFromToken = "0"
	} else {
		response.ResPricePerFromToken = strconv.FormatFloat(resPricePerFromToken, 'f', -1, 64)
	}
	if quote.Action.ToToken.PriceUSD == "" {
		response.ResPricePerToToken = "0"
	} else {
		response.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)
	}
	response.FromTokenPrice = quote.Estimate.FromAmountUSD
	response.ToTokenPrice = quote.Estimate.ToAmountUSD
	response.PriceImpact = strconv.FormatFloat(quote.Action.Slippage, 'f', -1, 64)
	toAmountIn, ok := new(big.Float).SetString("0")
	if quote.Estimate.ToAmountMin != "" {
		toAmountIn, ok = new(big.Float).SetString(quote.Estimate.ToAmountMin)
		if !ok {
			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "set string error")
		}
	}
	response.ApproveAddress = quote.Estimate.ApprovalAddress
	minimumReceived, _ := new(big.Float).Quo(toAmountIn, big.NewFloat(wei)).Float64()
	response.MinimumReceived = strconv.FormatFloat(minimumReceived, 'f', -1, 64)
	response.SellToken = SellToken
	response.BuyToken = BuyToken
	return &response, err
}

func (l *LiFi) GetLiFiSwap(request *pb.ExchangeSwapRequest, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error) {
	sellAmountParam, err := l.helper.ConvertStringValueToFloatWei(request.SellAmount, srcTokenDecimals)
	sellAmountParam = math.Floor(sellAmountParam)
	_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)
	chainKey := l.FindChainKey(request.Chain)
	request.SellToken = l.GetLiFiSpecificNativeToken(request.Chain, request.SellToken)
	request.BuyToken = l.GetLiFiSpecificNativeToken(request.Chain, request.BuyToken)
	url := fmt.Sprintf(l.env.Swap.LIFIEndpoint+"/quote?fromChain=%s&toChain=%s&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&fromAmount=%s&allowBridges=%s",
		chainKey, chainKey, request.SellToken, request.BuyToken, request.TakerAddress, request.TakerAddress, _sellAmountParam, l.allowedBridges)
	body, err := l.httpRequest.GetRequest(url)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
	}
	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		l.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var response pb.ExchangeSwapResponse
	response.To = quote.Estimate.ApprovalAddress
	response.Value = l.helper.ConvertHexFloatString(quote.TransactionRequest.Value)
	response.Data = quote.TransactionRequest.Data
	response.GasLimit = l.helper.ConvertHexFloatString(quote.TransactionRequest.GasLimit)
	response.Gas = response.GasLimit
	response.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", response.To, response.Value, response.Data, response.GasLimit)
	return &response, err
}

func (l *LiFi) ValidateChain(fromChain string, toChain string) (bool bool, from_chain string, to_chain string) {
	boolValFrom := false
	boolValTo := false
	if fromChain != "" && toChain == "" {
		toChain = fromChain
	} else if toChain != "" && fromChain == "" {
		fromChain = toChain
	}
	for _, item := range l.lifiChainList {
		if strings.ToLower(fromChain) == strings.ToLower(item.ChainKey) || fromChain == item.ChainId {
			boolValFrom = true
		}
		if strings.ToLower(toChain) == strings.ToLower(item.ChainKey) || toChain == item.ChainId {
			boolValTo = true
		}
	}
	if boolValTo && boolValFrom {
		return true, fromChain, toChain
	} else {
		return false, "", ""
	}
}

func (l *LiFi) FindChainKey(fromChain string) string {
	for _, item := range l.env.EVM.Cfg.Wallets {
		if strings.ToLower(fromChain) == strings.ToLower(item.ChainName) {
			return item.Bridge.ChainId
		}
	}
	return ""
}
