package zeroswap

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"bridge-allowance/web/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/chenzhijie/go-web3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IZeroSwap interface {
	GetExchangeTokens(info config.Wallets) (*pb.ExchangeTokenResponse, error)
	GetExchangeQuote(in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error)
	GetExchangeSwap(in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error)
}

type ZeroSwapService struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	util         *utils.UtilConf
	defaultToken []string
}

func NewZeroSwapService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *ZeroSwapService {
	return &ZeroSwapService{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (z *ZeroSwapService) IsChainSupported(info config.ChainData) bool {
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains")
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		return false
	}
	var result []ChainInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}
	for _, search := range result {
		if search.ChainId == info.ChainId {
			return true
		}
	}
	return false
}

func (z *ZeroSwapService) CallZeroSwap(inputRequest InputRequest, info config.ChainData, query string) (*ZeroSwapModel, error) {
	chainID := fmt.Sprint(info.ChainId)
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains/" + chainID + "/" + query + "?sellTokenAddress=" + strings.ToLower(inputRequest.SellToken) + "&buyTokenAddress=" + inputRequest.BuyToken + "&sellAmount=" + inputRequest.SellAmount + "&fromAddress=" + inputRequest.TakerAddress + "&slippagePercentage=" + inputRequest.Slippage)
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var to string
	if result["to"] != nil {
		to = result["to"].(string)
	}
	var data string
	if result["data"] != nil {
		data = result["data"].(string)
	}
	if result["estimatedPriceImpact"] == nil {
		result["estimatedPriceImpact"] = ""
	}
	// var exchangeQuote ZeroSwapModel
	exchangeQuote := ZeroSwapModel{
		EstimatedGas:         result["estimatedGas"].(string),
		GasPrice:             result["gasPrice"].(string),
		BuyAmount:            result["buyAmount"].(string),
		SellAmount:           result["sellAmount"].(string),
		Price:                result["price"].(string),
		To:                   to,
		Data:                 data,
		EstimatedPriceImpact: result["estimatedPriceImpact"].(string),
	}

	return &exchangeQuote, err
}

func (z *ZeroSwapService) GetExchangeTokens(info config.ChainData) (*pb.ExchangeTokenResponse, error) {
	chainID := fmt.Sprint(info.ChainId)
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.ExchangeTokenResponse{}, fmt.Errorf("chain " + info.ChainName + " not supported")
	}
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains/" + chainID + "/tokens")
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	tokens := result["tokens"].([]interface{})
	responseStruct := pb.ExchangeTokenResponse{}
	for _, item := range tokens {
		token := item.(map[string]interface{})
		var tokenName string
		if token["name"] != nil {
			tokenName = token["name"].(string)
		}
		var tokenLogo string
		if token["logoURI"] != nil {
			tokenLogo = token["logoURI"].(string)
		}
		tokenAddress := z.GetFrontierSpecificNativeToken(token["address"].(string), token["symbol"].(string))
		exchangeTokenInfo := pb.ExchangeTokenInfo{
			TokenAddress:  tokenAddress,
			TokenDecimals: fmt.Sprint(token["decimals"]),
			TokenSymbol:   token["symbol"].(string),
			TokenName:     tokenName,
			TokenLogoUrl:  tokenLogo,
			LogoUrl:       z.env.Swap.ZeroSwapLogoUrl,
		}
		responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
	}
	return &responseStruct, err
}

func (z *ZeroSwapService) GetFreeTradeCount(request *pb.FreeTradeCountRequest, info config.ChainData) (*pb.FreeTradeCountResponse, error) {
	chainID := fmt.Sprint(info.ChainId)
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.FreeTradeCountResponse{}, status.Errorf(codes.Unavailable, "chain "+info.ChainName+" not supported", "Unsupported")
	}
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains/" + chainID + "/trade-count?account=" + request.Account)
	body, err := z.httpRequest.GetRequestWithHeaders(url, "apikey", z.env.Swap.ZeroSwapApiKey)
	if err != nil {
		z.logger.Error(err)
		return &pb.FreeTradeCountResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var response pb.FreeTradeCountResponse
	var tradeCountInfo TradeCountInfo
	err = json.Unmarshal(body, &tradeCountInfo)
	if err != nil {
		z.logger.Error(err)
		return &pb.FreeTradeCountResponse{}, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	response.Account = tradeCountInfo.Account
	response.FreeTradeCount = tradeCountInfo.FreeTradeCount
	response.ChainId = fmt.Sprint(tradeCountInfo.ChainId)
	response.ChainName = tradeCountInfo.ChainName
	return &response, err
}

func (z *ZeroSwapService) GetExchangeQuote(request *pb.ExchangeQuoteRequest, info config.ChainData, walletInfo config.Wallets) (*pb.ExchangeQuoteResponse, error) {
	var response pb.ExchangeQuoteResponse
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "chain "+info.ChainName+" not supported", "Unsupported")
	}
	// assigning input request to new struct
	inputRequest := InputRequest{
		Chain:        request.Chain,
		TakerAddress: request.TakerAddress,
		SellToken:    request.SellToken,
		BuyToken:     request.BuyToken,
		SellAmount:   request.SellAmount,
		Slippage:     request.Slippage,
		ExchangeType: request.ExchangeType,
	}
	// Make API Call to zeroswap
	resp, err := z.CallZeroSwap(inputRequest, info, "quote")
	if err != nil {
		z.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	response.ResAmount = resp.BuyAmount
	response.PriceImpact = resp.EstimatedPriceImpact
	resAmount, ok := new(big.Float).SetString(response.ResAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error resAmount "+response.ResAmount, "Internal Error")
	}

	sellAmount, ok := new(big.Float).SetString(request.SellAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error sellAmount "+request.SellAmount, "Internal Error")
	}

	//Calculating Price per token := resAmount/sellAmount
	resPricePerToken, _ := new(big.Float).Quo(resAmount, sellAmount).Float64()
	response.ResPricePerFromToken = strconv.FormatFloat(resPricePerToken, 'f', -1, 64)

	// Calculating Price per Token for User = 1/resPricePerToken
	resPricePerToToken, _ := new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
	response.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)

	minimumReceived, err := z.GetMinimumReceived(response.ResPricePerFromToken, inputRequest.Slippage)
	if err != nil {
		z.logger.Errorf("Error while getting minimumReceived", err.Error())
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var toQuote float64
	var fromQuote float64
	if strings.ToLower(request.BuyToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := z.coinGecko.GetTokenExchange("usd", request.Chain)
		if err != nil {
			z.logger.Errorf("Error for Exchange Quote Price for token  request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}

	} else {

		quoteRate, err := z.coinGecko.GetTokenExchangeForContract(request.Chain, strings.ToLower(request.BuyToken), "usd")
		if err != nil {
			z.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}
	}
	if strings.ToLower(request.SellToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := z.coinGecko.GetTokenExchange("usd", request.Chain)
		if err != nil {
			z.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}

	} else {
		quoteRate, err := z.coinGecko.GetTokenExchangeForContract(request.Chain, strings.ToLower(request.SellToken), "usd")
		if err != nil {
			z.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}
	}

	// Calculating token price to sell : quoteRate*sellAmount
	fromTokenPrice, _ := new(big.Float).Mul(big.NewFloat(fromQuote), sellAmount).Float64()
	// Calculate token price to buy : quoteRate*buyAmount
	toTokenPrice, _ := new(big.Float).Mul(big.NewFloat(toQuote), resAmount).Float64()

	response.FromTokenPrice = strconv.FormatFloat(fromTokenPrice, 'f', -1, 64)
	response.ToTokenPrice = strconv.FormatFloat(toTokenPrice, 'f', -1, 64)
	response.SellToken = request.SellToken
	response.BuyToken = request.BuyToken
	response.MinimumReceived = minimumReceived
	isNativeToken := z.IsNativeToken(request.SellToken, walletInfo)
	if isNativeToken {
		response.ApproveAddress = info.ZeroswapSwapConfig.ApproveAddress
	} else {
		response.ApproveAddress = info.ZeroswapSwapConfig.GasLessApproveAddress
	}
	return &response, nil
}

func (z *ZeroSwapService) GetSignatureData(request *pb.ExchangeSignatureRequest, info config.ChainData) (*pb.ExchangeSignatureResponse, error) {
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.ExchangeSignatureResponse{}, status.Errorf(codes.Unavailable, "chain "+info.ChainName+" not supported", "Unsupported")
	}
	signatureRequest := SignatureRequest{
		Account:           request.TakerAddress,
		SellTokenAddress:  strings.ToLower(request.SellToken),
		BuyTokenAddress:   strings.ToLower(request.BuyToken),
		SlippageTolerance: request.Slippage,
		SellAmount:        request.SellAmount,
	}
	chainID := fmt.Sprint(info.ChainId)
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains/" + chainID + "/signature?account=" + signatureRequest.Account + "&sellTokenAddress=" + signatureRequest.SellTokenAddress + "&buyTokenAddress=" + signatureRequest.BuyTokenAddress + "&slippagePercentage=" + signatureRequest.SlippageTolerance + "&sellAmount=" + signatureRequest.SellAmount)
	body, err := z.httpRequest.GetRequestWithHeaders(url, "apikey", z.env.Swap.ZeroSwapApiKey)
	if err != nil {
		z.logger.Errorf("signature fetch error %v", err)
		return nil, err
	}
	var response pb.ZeroswapSignatureResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		z.logger.Errorf("signature unmarshal error %v", err)
		return nil, err
	}
	signatureResponse := &pb.ExchangeSignatureResponse{
		ZeroswapData: &response,
	}
	return signatureResponse, nil
}

func (z *ZeroSwapService) GetExchangeSwap(request *pb.ExchangeSwapRequest, info config.ChainData, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error) {
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Unavailable, "chain "+info.ChainName+" not supported", "Unsupported")
	}
	// assigning input request to new struct
	inputRequest := InputRequest{
		Chain:        request.Chain,
		TakerAddress: request.TakerAddress,
		SellToken:    request.SellToken,
		BuyToken:     request.BuyToken,
		SellAmount:   request.SellAmount,
		Slippage:     request.Slippage,
		ExchangeType: request.ExchangeType,
	}
	var response pb.ExchangeSwapResponse

	resp, err := z.CallZeroSwap(inputRequest, info, "swap")
	if err != nil {
		z.logger.Errorf("Error for Exchange Swap request  is : %v", err.Error())
		return &pb.ExchangeSwapResponse{}, err
	}
	sellAmountParam, err := z.helper.ConvertStringValueToFloatWei(inputRequest.SellAmount, srcTokenDecimals)
	sellAmountParam = math.Floor(sellAmountParam)

	if contains(z.defaultToken, strings.ToLower(inputRequest.SellToken)) {
		response.Value = strconv.FormatFloat(sellAmountParam, 'f', -1, 64)
	} else {
		response.Value = "0"
	}
	response.To = resp.To
	response.Data = resp.Data
	response.GasLimit = resp.EstimatedGas
	response.Gas = resp.GasPrice
	response.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", response.To, response.Value, response.Data, response.GasLimit)
	return &response, err
}

func (z *ZeroSwapService) ExecuteZeroSwap(request *pb.ExchangeSwapExecuteRequest, info config.ChainData, srcTokenDecimals string, walletInfo config.Wallets) (*pb.ExchangeSwapExecuteResponse, error) {
	var response pb.ExchangeSwapExecuteResponse
	chainID := fmt.Sprint(info.ChainId)
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.ExchangeSwapExecuteResponse{}, status.Errorf(codes.Unavailable, "chain "+info.ChainName+" not supported", "Unsupported")
	}
	buyTokenSymbol := z.GetTokenSymbol(walletInfo, strings.ToLower(request.ZeroSwapPayload.BuyToken), chainID)
	sellTokenSymbol := z.GetTokenSymbol(walletInfo, strings.ToLower(request.ZeroSwapPayload.SellToken), chainID)
	buyToken := z.GetZeroswapSpecicToken(strings.ToLower(request.ZeroSwapPayload.BuyToken))
	sellToken := z.GetZeroswapSpecicToken(strings.ToLower(request.ZeroSwapPayload.SellToken))
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "swap/execute")
	zeroSwapPayload := ZeroSwapPayload{
		ChainId:             chainID,
		Signature:           request.ZeroSwapPayload.Signature,
		BuyToken:            buyTokenSymbol,
		SellToken:           sellTokenSymbol,
		BuyTokenAddress:     buyToken,
		SellTokenAddress:    sellToken,
		SellAmount:          request.ZeroSwapPayload.SellAmount,
		Recipient:           request.ZeroSwapPayload.TakerAddress,
		SlippageTolerance:   request.ZeroSwapPayload.Slippage,
		TransactionDeadline: request.ZeroSwapPayload.TransactionalDeadline,
		AffiliateAddress:    FrontierAffiliateAddress,
	}
	payload, err := json.Marshal(zeroSwapPayload)
	body, err := z.httpRequest.PostRequestWithHeaders(url, string(payload), "apikey", z.env.Swap.ZeroSwapApiKey)
	if err != nil {
		var respError RespError
		err = json.Unmarshal(body, &respError)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
		}
		return &response, status.Errorf(codes.Internal, respError.Message, "Internal Error")
	}
	var exchangeSwapTransactionResponse *pb.ExchangeSwapTransactionResponse
	err = json.Unmarshal(body, &exchangeSwapTransactionResponse)
	if err != nil {
		return &response, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	var executeResponse pb.ExecuteResponse
	executeResponse.TxHash = exchangeSwapTransactionResponse.TransactionHash
	response.ExecuteResponse = &executeResponse
	return &response, err
}

func (z *ZeroSwapService) GetTokenApproval(request *pb.TokenApprovalRequest, info config.ChainData) (*pb.TokenApprovalResponse, error) {
	chainID := fmt.Sprint(info.ChainId)
	chainSupported := z.IsChainSupported(info)
	if !chainSupported {
		return &pb.TokenApprovalResponse{}, fmt.Errorf("chain " + info.ChainName + " not supported")
	}
	url := fmt.Sprintf(z.env.Swap.ZeroSwapEndpoint + "chains/" + chainID + "/tokens/" + request.Token + "/approve?gasless=" + request.Gasless)
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return &pb.TokenApprovalResponse{}, err
	}
	var response pb.TokenApprovalResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		z.logger.Error(err)
		return &pb.TokenApprovalResponse{}, err
	}
	return &response, err
}

func (z *ZeroSwapService) GetMinimumReceived(resPricePerFromToken string, slippage string) (string, error) {

	toTokenAmountFloat, err := strconv.ParseFloat(resPricePerFromToken, 64)
	if err != nil {
		toTokenAmountFloat = 0
	}
	slippageFloat, err := strconv.ParseFloat(slippage, 64)
	if err != nil {
		slippageFloat = 0
	}
	minReceived := (toTokenAmountFloat * (100.0 - slippageFloat)) / 100.0
	minimumReceived := fmt.Sprintf("%f", minReceived)
	return minimumReceived, nil
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func (z *ZeroSwapService) GetZeroswapExecuteRequest(ctx *gin.Context, chainGroup string, chain string) (*pb.ZeroSwapExecuteRequest, error) {
	var payload models.GasLessSwapBody
	decoder := json.NewDecoder(ctx.Request.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		z.logger.Errorf("error %s", err)
		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
		return nil, err
	}

	valid, validAddress, err := z.util.ValidateAddress(payload.TakerAddress, chainGroup, chain)
	if err != nil || !valid {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	validSellToken, validAddressSellToken, err := z.util.ValidateAddress(payload.SellToken, chainGroup, chain)
	if err != nil || !validSellToken {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	validBuyToken, validAddressBuyToken, err := z.util.ValidateAddress(payload.BuyToken, chainGroup, chain)
	if err != nil || !validBuyToken {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	zeroswapPayload := &pb.ZeroSwapExecuteRequest{
		TakerAddress:          validAddress,
		SellToken:             validAddressSellToken,
		BuyToken:              validAddressBuyToken,
		SellAmount:            payload.SellAmount,
		Signature:             payload.Signature,
		TransactionalDeadline: payload.TransactionalDeadline,
		Slippage:              payload.Slippage,
	}
	return zeroswapPayload, nil
}
func (z *ZeroSwapService) GetZeroswapSpecicToken(tokenAddress string) string {
	if tokenAddress == "0x0000000000000000000000000000000000001010" {
		return "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	} else {
		return tokenAddress
	}
}

func (z *ZeroSwapService) GetTokenSymbol(info config.Wallets, address string, chain string) string {
	var symbol string
	if chain == strconv.Itoa(info.ChainID) && address == info.NativeTokenInfo.Address {
		return info.NativeTokenInfo.Symbol
	} else {
		web3Client := z.GetClient(info)
		tokenContract, err := web3Client.Eth.NewContract(TokenABI, address)
		if err == nil {
			buyTokenSymbol, err := tokenContract.Call(Symbol)
			if err == nil {
				symbol = fmt.Sprint(buyTokenSymbol)
			}
		}
	}
	return symbol
}

func (z *ZeroSwapService) GetClient(info config.Wallets) *web3.Web3 {
	web3Client, err := web3.NewWeb3(info.RPC)
	if err != nil {
		z.logger.Error("Error initializing web3", err)
		return web3Client
	}
	return web3Client
}

func (z *ZeroSwapService) IsNativeToken(token string, info config.Wallets) bool {
	if strings.ToLower(token) == info.NativeTokenInfo.Address {
		return true
	} else {
		return false
	}
}

func (z *ZeroSwapService) GetFrontierSpecificNativeToken(tokenAddress string, symbol string) string {
	if symbol == "MATIC" {
		tokenAddress = "0x0000000000000000000000000000000000001010"
	}
	return tokenAddress
}
