package oneinch

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strconv"
	"strings"
)

type IOneInch interface {
	GetExchangeTokens(info config.Wallets, nativeTokenInfo *pb.ExchangeTokenInfo) (*pb.ExchangeTokenResponse, error)
	GetExchangeQuote(request *pb.ExchangeQuoteRequest, info config.ChainData, srcTokenDecimals string) (*pb.ExchangeQuoteResponse, error)
	GetExchangeSwap(request *pb.ExchangeSwapRequest, info config.Wallets, infoSwap config.ChainData, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error)
}

type OneInchService struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	defaultToken []string
}

func NewOneInchService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *OneInchService {
	return &OneInchService{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (o *OneInchService) CallOneInch(inputRequest InputRequest, info config.ChainData, query string) (*ExchangeOneInchModel, error) {
	var errorBody ErrorResponse
	chainID := fmt.Sprint(info.ChainId)
	url := fmt.Sprintf(o.env.Swap.OneInchEndpoint + chainID + "/" + query + "?fromTokenAddress=" + inputRequest.SellToken + "&toTokenAddress=" + inputRequest.BuyToken + "&amount=" + inputRequest.SellAmount + "&fromAddress=" + inputRequest.TakerAddress + "&slippage=" + inputRequest.Slippage)
	body, errorResponse, err := o.httpRequest.GetRequestWithErrorResponse(url)
	if err != nil {
		o.logger.Error(err)
		// if error is present in the response body then will unmarshal the error response
		err = json.Unmarshal(errorResponse, &errorBody)
		if strings.Contains(errorBody.Description, fmt.Sprintf("Not enough")) {
			return nil, status.Errorf(codes.NotFound, errorBody.Description, "Insufficient balance")
		}
		return nil, status.Errorf(codes.Internal, errorBody.Description, "Internal Error")
	}
	var exchangeQuote ExchangeOneInchModel
	err = json.Unmarshal(body, &exchangeQuote)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	return &exchangeQuote, err
}

func (o *OneInchService) GetExchangeTokens(info config.Wallets, nativeTokenInfo *pb.ExchangeTokenInfo) (*pb.ExchangeTokenResponse, error) {
	chainID := fmt.Sprint(info.ChainID)
	url := fmt.Sprintf(o.env.Swap.OneInchEndpoint + chainID + "/tokens")

	body, err := o.httpRequest.GetRequest(url)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	tokens := result["tokens"].(map[string]interface{})
	responseStruct := pb.ExchangeTokenResponse{}
	var nativeToken pb.ExchangeTokenInfo
	for _, item := range tokens {
		token := item.(map[string]interface{})
		if token["address"].(string) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" && info.NativeTokenInfo.Address != "" {
			nativeToken.TokenAddress = info.NativeTokenInfo.Address
			nativeToken.TokenDecimals = info.NativeTokenInfo.Decimals
			nativeToken.TokenSymbol = info.NativeTokenInfo.Symbol
			nativeToken.TokenName = info.NativeTokenInfo.Name
			nativeToken.TokenLogoUrl = o.CheckNilForString(token["logoURI"]).(string)
			nativeToken.LogoUrl = o.env.Swap.OneInchLogoUrl
		} else {
			exchangeTokenInfo := pb.ExchangeTokenInfo{
				TokenAddress:  token["address"].(string),
				TokenDecimals: fmt.Sprint(o.CheckNilForString(token["decimals"])),
				TokenSymbol:   o.CheckNilForString(token["symbol"]).(string),
				TokenName:     o.CheckNilForString(token["name"]).(string),
				TokenLogoUrl:  o.CheckNilForString(token["logoURI"]).(string),
				LogoUrl:       o.env.Swap.OneInchLogoUrl,
			}
			responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
		}
	}
	responseStruct.ExchangeTokens = append([]*pb.ExchangeTokenInfo{
		&nativeToken,
	}, responseStruct.ExchangeTokens...)
	return &responseStruct, err
}

func (o *OneInchService) GetExchangeQuote(request *pb.ExchangeQuoteRequest, info config.ChainData, srcTokenDecimals string) (*pb.ExchangeQuoteResponse, error) {
	var response pb.ExchangeQuoteResponse
	var inputRequest InputRequest

	// assigning input request to new struct
	inputRequest = InputRequest{
		Chain:        request.Chain,
		TakerAddress: request.TakerAddress,
		SellToken:    request.SellToken,
		BuyToken:     request.BuyToken,
		SellAmount:   request.SellAmount,
		Slippage:     request.Slippage,
		ExchangeType: request.ExchangeType,
	}
	sellAmountParams, err := o.helper.ConvertToWei(inputRequest.SellAmount, srcTokenDecimals)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	inputRequest.SellAmount = sellAmountParams.String()
	// in optimism chain 0x4200000000000000000000000000000000000006 is a valid contract
	if inputRequest.Chain != "optimism" {
		if contains(o.defaultToken, strings.ToLower(inputRequest.SellToken)) {
			inputRequest.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(o.defaultToken, strings.ToLower(inputRequest.BuyToken)) {
			inputRequest.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	// Make API Call to 1inch swap
	response1inch, err := o.CallOneInch(inputRequest, info, "quote")
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, err
	}
	toTokenDecimalString := strconv.Itoa(response1inch.ToToken.Decimals)
	resAmountWithDecimals, err := o.helper.ConvertStringValueToBigFloat(response1inch.ToTokenAmount, toTokenDecimalString)
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	response.ResAmount = strconv.FormatFloat(resAmountWithDecimals, 'f', -1, 64)

	resAmount, ok := new(big.Float).SetString(response.ResAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error resAmount "+response.ResAmount, "type error")
	}

	sellAmount, ok := new(big.Float).SetString(request.SellAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error sellAmount "+request.SellAmount, "type error")
	}
	// Calculate Res Price Per FromToken
	resPricePerToken, _ := new(big.Float).Quo(resAmount, sellAmount).Float64()
	response.ResPricePerFromToken = strconv.FormatFloat(resPricePerToken, 'f', -1, 64)

	// Calculate Res Price Per ToToken
	resPricePerToToken, _ := new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
	response.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)

	minimumReceived, err := o.GetMinimumReceived(response.ResAmount, inputRequest.Slippage)
	if err != nil {
		o.logger.Errorf("Error while getting minimumReceived", err.Error())
		return &pb.ExchangeQuoteResponse{}, err
	}

	var toQuote float64
	var fromQuote float64
	if strings.ToLower(inputRequest.BuyToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := o.coinGecko.GetTokenExchange("usd", inputRequest.Chain)
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}
	} else {
		quoteRate, err := o.coinGecko.GetTokenExchangeForContract(inputRequest.Chain, strings.ToLower(inputRequest.BuyToken), "usd")
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}

	}
	if strings.ToLower(inputRequest.SellToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := o.coinGecko.GetTokenExchange("usd", inputRequest.Chain)
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}

	} else {
		quoteRate, err := o.coinGecko.GetTokenExchangeForContract(inputRequest.Chain, strings.ToLower(inputRequest.SellToken), "usd")
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}
	}

	fromTokenPrice, _ := new(big.Float).Mul(big.NewFloat(fromQuote), sellAmount).Float64()
	toTokenPrice, _ := new(big.Float).Mul(big.NewFloat(toQuote), resAmount).Float64()

	response.FromTokenPrice = strconv.FormatFloat(fromTokenPrice, 'f', -1, 64)
	response.ToTokenPrice = strconv.FormatFloat(toTokenPrice, 'f', -1, 64)

	response.PriceImpact = o.GetPriceImpact(response)
	response.MinimumReceived = minimumReceived
	response.ApproveAddress = info.OneInchSwapConfig.ApproveAddress
	response.BuyToken = request.BuyToken
	response.SellToken = request.SellToken
	return &response, nil
}

func (o *OneInchService) GetExchangeSwap(request *pb.ExchangeSwapRequest, info config.Wallets, infoSwap config.ChainData, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error) {
	var response pb.ExchangeSwapResponse
	var inputRequest InputRequest

	// assigning input request to new struct
	inputRequest = InputRequest{
		Chain:        request.Chain,
		TakerAddress: request.TakerAddress,
		SellToken:    request.SellToken,
		BuyToken:     request.BuyToken,
		SellAmount:   request.SellAmount,
		Slippage:     request.Slippage,
		ExchangeType: request.ExchangeType,
	}
	sellAmountParams, err := o.helper.ConvertToWei(inputRequest.SellAmount, srcTokenDecimals)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	inputRequest.SellAmount = sellAmountParams.String()

	if inputRequest.Slippage != "" {
		slippagePercent, ok := new(big.Float).SetString(inputRequest.Slippage)
		if !ok {
			return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, "SetString: error for Slippage : "+inputRequest.Slippage, "type error")
		}
		_resSplippage, _ := new(big.Float).Quo(slippagePercent, big.NewFloat(100)).Float64()
		inputRequest.Slippage = strconv.FormatFloat(_resSplippage, 'f', -1, 64)

	}
	// in optimism chain 0x4200000000000000000000000000000000000006 is a valid contract
	if inputRequest.Chain != "optimism" {
		if contains(o.defaultToken, strings.ToLower(inputRequest.SellToken)) {
			inputRequest.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(o.defaultToken, strings.ToLower(inputRequest.BuyToken)) {
			inputRequest.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	response1inch, err := o.CallOneInch(inputRequest, infoSwap, "swap")
	if err != nil {
		o.logger.Errorf("Error for Exchange Swap request  is : %v", err.Error())
		return &pb.ExchangeSwapResponse{}, err
	}
	response.To = response1inch.Tx.To
	response.Data = response1inch.Tx.Data
	if contains(o.defaultToken, strings.ToLower(inputRequest.SellToken)) {
		response.Value = inputRequest.SellAmount
	} else {
		response.Value = "0"
	}
	// if the Gas Value is 0 from the 1inch API then we will calculate the gas value through ethrpc call
	if response1inch.Tx.Gas == 0 {
		value, err := o.helper.ConvertToWei(request.SellAmount, srcTokenDecimals)
		if err != nil {
			o.logger.Error("Error while converting to wei", err.Error())
			return &pb.ExchangeSwapResponse{}, err
		}
		bigIntValue, err := o.GetBigIntValue(request.SellToken, value, info.NativeTokenInfo.Address)
		if err != nil {
			o.logger.Error("Error for getting big int value", err.Error())
			return &pb.ExchangeSwapResponse{}, err
		}

		// Calling GetGasLimit for 0 gas value from 1inch API
		gasLimit, err := o.helper.GetGasLimit(info, infoSwap.OneInchSwapConfig.ApproveAddress, response1inch.Tx.Data, bigIntValue, request.TakerAddress)
		if err != nil {
			o.logger.Errorf("Error while calculating gasLimit : %v", err.Error())
			return &pb.ExchangeSwapResponse{}, err
		}
		response1inch.Tx.Gas = gasLimit
	}
	response.GasLimit = fmt.Sprintf("%v", response1inch.Tx.Gas)
	response.Gas = response1inch.Tx.GasPrice
	response.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", response.To, response.Value, response.Data, response.GasLimit)
	return &response, err
}

// GetMinimumReceived calculating minimum received for swap the tokens depending on slippage
func (o *OneInchService) GetMinimumReceived(resAmount string, slippage string) (string, error) {
	toTokenAmountFloat := o.helper.ConvertStringToFloat64(resAmount)
	slippageFloat := o.helper.ConvertStringToFloat64(slippage)
	minReceived := (toTokenAmountFloat * (100.0 - slippageFloat)) / 100.0
	minimumReceived := fmt.Sprintf("%f", minReceived)
	return minimumReceived, nil
}

// GetPriceImpact if there is no response from the APIs
func (o *OneInchService) GetPriceImpact(response pb.ExchangeQuoteResponse) string {
	var priceImpactString string
	toTokenPrice := o.helper.ConvertStringToFloat64(response.ToTokenPrice)
	fromTokenPrice := o.helper.ConvertStringToFloat64(response.FromTokenPrice)
	if response.ToTokenPrice == "0" || response.FromTokenPrice == "0" {
		priceImpactString = "1"
	} else {
		priceImpactFloat := (1.0 - (toTokenPrice / fromTokenPrice)) * 100
		priceImpactString = strconv.FormatFloat(priceImpactFloat, 'f', -1, 64)
	}
	if priceImpactString == "NaN" || priceImpactString < "0" {
		priceImpactString = "0.01"
	}
	return priceImpactString
}
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// CheckNilForString check if the string is nil or not
func (o *OneInchService) CheckNilForString(v interface{}) interface{} {
	if v == nil {
		return ""
	} else {
		return v
	}
}

// GetBigIntValue Updating bigInt value on the basis of native and non-native tokens for gas limit calculation
func (o *OneInchService) GetBigIntValue(sellToken string, value *big.Int, nativeTokenAddress string) (*big.Int, error) {
	if strings.ToLower(nativeTokenAddress) == strings.ToLower(sellToken) {
		return value, nil
	}
	return big.NewInt(0), nil
}
