package cowswap

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/lifi"
	"bridge-allowance/utils"
	"bridge-allowance/web/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type ICowSwap interface {
	GetExchangeTokens(info config.Wallets) (*pb.ExchangeTokenResponse, error)
	GetExchangeQuote(in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error)
	GetExchangeSwap(in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error)
}

type CowSwapService struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	util         *utils.UtilConf
	lifi         *lifi.LiFi
	defaultToken []string
}

func NewCowSwapService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko, lifi *lifi.LiFi) *CowSwapService {
	return &CowSwapService{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		lifi:         lifi,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

// CallCowExchange is used to call the cowswap exchange quote api
func (c *CowSwapService) CallCowExchange(inputRequest InputRequestData, info config.ChainData) (*QuoteExchange, error) {
	jsonReq, _ := json.Marshal(inputRequest)
	reqBody := bytes.NewBuffer(jsonReq)
	var chainName string
	switch info.ChainId {
	case 1:
		chainName = ETHEREUM
	case 100:
		chainName = GNOSIS
	}

	url := c.env.Swap.CowSwapUrl + chainName + "/api/v1/quote"
	body, err := c.httpRequest.PostRequest(url, reqBody)
	if err != nil {
		c.logger.Error("Error while calling the cowswap exchange quote api", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Error while calling the cowswap exchange quote api")
	}
	var quoteExchange QuoteExchange
	err = json.Unmarshal(body, &quoteExchange)
	if err != nil {
		c.logger.Error("Error while unmarshalling the cowswap exchange quote api response", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Error while unmarshalling the cowswap exchange quote api response")
	}
	return &quoteExchange, nil
}

func (c *CowSwapService) GetExchangeTokens(info config.ChainData, walletInfo config.Wallets) (*pb.ExchangeTokenResponse, error) {
	responseStruct := pb.ExchangeTokenResponse{}
	chainSupported := c.IsSupportedChain(info)
	if !chainSupported {
		return nil, status.Errorf(codes.Unavailable, "Chain not supported", "Chain not supported")
	}
	var exchangeTokens ExchangeCowResponse
	var tokens Tokens
	var tokenUrlList = []string{CowDaoList, UmaList, AaveList, SynthetixList, WrappedList, SetList, OpynList, RollList, CmcAllList, CmcStableCoin, GeminiList, BaList, UniList, OptimismList, ArbitrumList, CeloList, HoneSwapXDAi}
	// to store the token list into unbuffered channel
	chanResponse := make(chan Tokens)
	// wait group is used to wait for all the goroutines to finish
	wg := new(sync.WaitGroup)

	for _, url := range tokenUrlList {
		// increment the wait group counter
		wg.Add(1)
		// call the goroutine
		go c.FetchTokenList(url, exchangeTokens, info, wg, chanResponse)
		// Add the token list to the tokens struct
		tokens = append(tokens, <-chanResponse...)
	}
	wg.Wait()
	close(chanResponse)

	// Remove duplicates from the tokens
	tokens = tokens.RemoveDuplicates()

	// Assigning the tokens to the response struct
	var nativeTokenInfo pb.ExchangeTokenInfo
	for _, token := range tokens {
		exchangeTokenInfo := pb.ExchangeTokenInfo{
			TokenAddress:  token.Address,
			TokenDecimals: fmt.Sprint(token.Decimals),
			TokenSymbol:   token.Symbol,
			TokenName:     token.Name,
			TokenLogoUrl:  token.LogoURI,
			LogoUrl:       c.env.Swap.CowSwapLogoUrl,
		}
		responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
	}
	// Assigning the native tokens to the response struct for some chains
	if info.ChainId == 1 || info.ChainId == 100 {
		nativeTokenInfo.TokenAddress = walletInfo.NativeTokenInfo.Address
		nativeTokenInfo.TokenDecimals = walletInfo.NativeTokenInfo.Decimals
		nativeTokenInfo.TokenSymbol = walletInfo.NativeTokenInfo.Symbol
		nativeTokenInfo.TokenName = walletInfo.NativeTokenInfo.Name
		nativeTokenInfo.TokenLogoUrl = walletInfo.NativeTokenInfo.LogoURI
		nativeTokenInfo.LogoUrl = c.env.Swap.CowSwapLogoUrl

	}
	responseStruct.ExchangeTokens = append([]*pb.ExchangeTokenInfo{
		&nativeTokenInfo,
	}, responseStruct.ExchangeTokens...)

	if len(responseStruct.ExchangeTokens) == 0 {
		return nil, status.Errorf(codes.NotFound, "No tokens found for "+info.ChainName+" chain", "No tokens found")
	}

	return &responseStruct, nil
}

func (c *CowSwapService) GetExchangeQuote(request *pb.ExchangeQuoteRequest, info config.ChainData, srcTokenDecimals string, dstTokenDecimals string, walletInfo config.Wallets) (*pb.ExchangeQuoteResponse, error) {
	var response pb.ExchangeQuoteResponse
	chainSupported := c.IsSupportedChain(info)
	if !chainSupported {
		return nil, status.Errorf(codes.Unavailable, "Chain not supported")
	}
	nativeTokenAddress := c.GetNativeTokenAddress(walletInfo)
	wrappedTokenAddress := c.GetWNativeTokenAddress(info)
	if (strings.ToLower(nativeTokenAddress) == strings.ToLower(request.SellToken) && strings.ToLower(wrappedTokenAddress) == strings.ToLower(request.BuyToken)) ||
		(strings.ToLower(wrappedTokenAddress) == strings.ToLower(request.SellToken) && strings.ToLower(nativeTokenAddress) == strings.ToLower(request.BuyToken)) {
		return c.lifi.GetLiFiQuote(request.SellToken, request.BuyToken, request.SellAmount, info.ChainName, srcTokenDecimals, dstTokenDecimals, request.TakerAddress, request.Slippage)
	} else if strings.ToLower(request.SellToken) == nativeTokenAddress || strings.ToLower(request.BuyToken) == nativeTokenAddress {
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.InvalidArgument, "", "The chain's native token cannot be used as the sell/buy token")
	} else {
		sellAmountFee := ConvertToWei(request.SellAmount, srcTokenDecimals)
		inputRequest := InputRequestData{
			SellToken:           strings.ToLower(request.SellToken),
			BuyToken:            strings.ToLower(request.BuyToken),
			Receiver:            request.TakerAddress,
			AppData:             "0x0000000000000000000000000000000000000000000000000000000000000000",
			SellTokenBalance:    "erc20",
			BuyTokenBalance:     "erc20",
			PriceQuality:        "optimal",
			SigningScheme:       "eip712",
			Kind:                "sell",
			From:                request.TakerAddress,
			SellAmountBeforeFee: sellAmountFee.String(),
		}

		// Call Cow Exchange API
		quoteExchange, err := c.CallCowExchange(inputRequest, info)
		if err != nil {
			c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeQuoteResponse{}, err
		}

		buyAmount := ToDecimal(quoteExchange.Quote.BuyAmount, dstTokenDecimals).String()
		resAmount, ok := new(big.Float).SetString(buyAmount)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "Invalid sell amount")
		}
		response.ResAmount = resAmount.String()

		var fromQuote float64
		quoteRate, err := c.coinGecko.GetTokenExchangeForContract(request.Chain, strings.ToLower(request.SellToken), "usd")
		if err != nil {
			//o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}

		sellAmount, ok := new(big.Float).SetString(request.SellAmount)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "Invalid sell amount")
		}

		// Calculating token price to sell : quoteRate*sellAmount
		fromTokenPrice, _ := new(big.Float).Mul(big.NewFloat(fromQuote), sellAmount).Float64()
		response.FromTokenPrice = strconv.FormatFloat(fromTokenPrice, 'f', -1, 64)
		response.ToTokenPrice = resAmount.String()

		resPricePerToken, _ := new(big.Float).Quo(resAmount, sellAmount).Float64()
		response.ResPricePerFromToken = strconv.FormatFloat(resPricePerToken, 'f', -1, 64)

		// Calculating Price per Token for User = 1/resPricePerToken
		resPricePerToToken, _ := new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
		response.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)

		response.PriceImpact = c.GetPriceImpact(response)
		minimumReceived, err := c.GetMinimumReceived(response.ResAmount, request.Slippage)
		if err != nil {
			c.logger.Errorf("error while getting minimumReceived", err.Error())
			return &pb.ExchangeQuoteResponse{}, err
		}
		response.MinimumReceived = minimumReceived
		response.ApproveAddress = info.CowSwapConfig.ApproveAddress
		response.SellToken = request.SellToken
		response.BuyToken = request.BuyToken
		return &response, nil
	}
}

// GetSignatureData returns the signature data for the given order
func (c *CowSwapService) GetSignatureData(request *pb.ExchangeSignatureRequest, info config.ChainData, srcDecimal string, dstDecimal string, walletInfo config.Wallets) (*pb.ExchangeSignatureResponse, error) {
	var signatureResponse pb.ExchangeSignatureResponse
	chainSupported := c.IsSupportedChain(info)
	if !chainSupported {
		return nil, status.Errorf(codes.Unavailable, "Chain not supported")
	}
	nativeTokenAddress := c.GetNativeTokenAddress(walletInfo)

	if strings.ToLower(request.SellToken) == nativeTokenAddress || strings.ToLower(request.BuyToken) == nativeTokenAddress {
		return &pb.ExchangeSignatureResponse{}, status.Error(codes.InvalidArgument, "The chain's native token cannot be used as the sell/buy token")
	} else {
		sellAmountFee := ConvertToWei(request.SellAmount, srcDecimal)
		inputs := InputRequestData{
			SellToken:           strings.ToLower(request.SellToken),
			BuyToken:            strings.ToLower(request.BuyToken),
			Receiver:            request.TakerAddress,
			AppData:             "0x0000000000000000000000000000000000000000000000000000000000000000",
			SellTokenBalance:    "erc20",
			BuyTokenBalance:     "erc20",
			PriceQuality:        "optimal",
			SigningScheme:       "eip712",
			Kind:                "sell",
			From:                request.TakerAddress,
			SellAmountBeforeFee: sellAmountFee.String(),
		}

		// Call Cow Exchange API
		cowSwapResponse, err := c.CallCowExchange(inputs, info)
		if err != nil {
			c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeSignatureResponse{}, err
		}
		mappedResponse := &pb.CowSwapSignatureResponse{
			SellToken:         cowSwapResponse.Quote.SellToken,
			BuyToken:          cowSwapResponse.Quote.BuyToken,
			Receiver:          cowSwapResponse.Quote.Receiver,
			SellAmount:        cowSwapResponse.Quote.SellAmount,
			BuyAmount:         cowSwapResponse.Quote.BuyAmount,
			ValidTo:           cowSwapResponse.Quote.ValidTo,
			AppData:           cowSwapResponse.Quote.AppData,
			FeeAmount:         cowSwapResponse.Quote.FeeAmount,
			Kind:              cowSwapResponse.Quote.Kind,
			PartiallyFillable: cowSwapResponse.Quote.PartiallyFillable,
			SellTokenBalance:  cowSwapResponse.Quote.SellTokenBalance,
			BuyTokenBalance:   cowSwapResponse.Quote.BuyTokenBalance,
			From:              cowSwapResponse.From,
			QuoteID:           cowSwapResponse.ID,
		}
		signatureResponse = pb.ExchangeSignatureResponse{
			CowswapData: mappedResponse,
		}
		return &signatureResponse, nil
	}
}

// ExecuteCowSwapOrder executes the order on CowSwap
func (c *CowSwapService) ExecuteCowSwapOrder(request *pb.ExchangeSwapExecuteRequest, info config.ChainData, srcTokenDecimals string, dstTokenDecimals string, walletInfo config.Wallets) (*pb.ExchangeSwapExecuteResponse, error) {
	var response pb.ExchangeSwapExecuteResponse
	chainSupported := c.IsSupportedChain(info)
	if !chainSupported {
		return nil, status.Errorf(codes.Unavailable, "", "Chain not supported")
	}
	nativeTokenAddress := c.GetNativeTokenAddress(walletInfo)
	wrappedTokenAddress := c.GetWNativeTokenAddress(info)

	if (strings.ToLower(nativeTokenAddress) == strings.ToLower(request.CowSwapPayload.SellToken) && strings.ToLower(wrappedTokenAddress) == strings.ToLower(request.CowSwapPayload.BuyToken)) ||
		(strings.ToLower(wrappedTokenAddress) == strings.ToLower(request.CowSwapPayload.SellToken) && strings.ToLower(nativeTokenAddress) == strings.ToLower(request.CowSwapPayload.BuyToken)) {
		swapRequestData := &pb.ExchangeSwapRequest{
			Chain:        info.ChainName,
			TakerAddress: request.CowSwapPayload.From,
			SellToken:    request.CowSwapPayload.SellToken,
			BuyToken:     request.CowSwapPayload.BuyToken,
			SellAmount:   request.CowSwapPayload.SellAmount,
			Slippage:     request.CowSwapPayload.Slippage,
		}
		swapResponseLifi, err := c.lifi.GetLiFiSwap(swapRequestData, srcTokenDecimals)
		if err != nil {
			c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeSwapExecuteResponse{}, err
		}
		response.ExchangeSwapResponse = swapResponseLifi
		return &response, nil
	} else if strings.ToLower(request.CowSwapPayload.SellToken) == nativeTokenAddress || strings.ToLower(request.CowSwapPayload.BuyToken) == nativeTokenAddress {
		return &pb.ExchangeSwapExecuteResponse{}, status.Error(codes.InvalidArgument, "The chain's native token cannot be used as the sell/buy token")
	} else {
		// call POST Order execution api
		orderResponse, err := c.CallPostOrderAPI(request, info)
		if err != nil {
			c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeSwapExecuteResponse{}, err
		}
		// call Get Order execution api to get the status of the order
		//orderStatus, err := c.CallGetOrderAPI(orderResponse.UID, info)@
		//if err != nil {
		//	c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		//	return &pb.ExchangeSwapExecuteResponse{}, err
		//}
		//if orderStatus.Status == "fulfilled" {
		//	cowTradeResponse, err := c.CallCowTradeAPI(orderResponse.UID, info)
		//	if err != nil {
		//		c.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		//		return &pb.ExchangeSwapExecuteResponse{}, err
		//	}
		//	response.ExecuteResponse.TxHash = cowTradeResponse[0].TxHash
		//}
		// As of now we are providing UID of swap order to the user
		var executeResponse pb.ExecuteResponse
		executeResponse.TxHash = orderResponse.UID
		response.ExecuteResponse = &executeResponse
		return &response, nil
	}
}

// GetCowSwapExecuteRequest returns the request for CowSwap to swap handler
func (c *CowSwapService) GetCowSwapExecuteRequest(ctx *gin.Context, chainGroup string, chain string) (*pb.CowSwapExecuteRequest, error) {
	var payload models.GasLessSwapBody
	decoder := json.NewDecoder(ctx.Request.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		c.logger.Errorf("error %s", err)
		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
		return nil, err
	}

	valid, validAddress, err := c.util.ValidateAddress(payload.Receiver, chainGroup, chain)
	if err != nil || !valid {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	validSellToken, validAddressSellToken, err := c.util.ValidateAddress(payload.SellToken, chainGroup, chain)
	if err != nil || !validSellToken {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	validBuyToken, validAddressBuyToken, err := c.util.ValidateAddress(payload.BuyToken, chainGroup, chain)
	if err != nil || !validBuyToken {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return nil, err
	}
	cowswapRequest := &pb.CowSwapExecuteRequest{
		SellToken:         validAddressSellToken,
		BuyToken:          validAddressBuyToken,
		Receiver:          validAddress,
		SellAmount:        payload.SellAmount,
		BuyAmount:         payload.BuyAmount,
		ValidTo:           int32(payload.ValidTo),
		AppData:           payload.AppData,
		FeeAmount:         payload.FeeAmount,
		Kind:              payload.Kind,
		PartiallyFillable: payload.PartiallyFillable,
		SellTokenBalance:  payload.SellTokenBalance,
		BuyTokenBalance:   payload.BuyTokenBalance,
		SigningScheme:     payload.SigningScheme,
		Signature:         payload.Signature,
		From:              payload.From,
		QuoteId:           int32(payload.QuoteId),
		Slippage:          payload.Slippage,
	}
	return cowswapRequest, nil
}

// GetMinimumReceived calculating minimum received for swap the tokens depending on slippage
func (c *CowSwapService) GetMinimumReceived(resAmount string, slippage string) (string, error) {
	toTokenAmountFloat := c.helper.ConvertStringToFloat64(resAmount)
	slippageFloat := c.helper.ConvertStringToFloat64(slippage)
	minReceived := (toTokenAmountFloat * (100.0 - slippageFloat)) / 100.0
	minimumReceived := fmt.Sprintf("%f", minReceived)
	return minimumReceived, nil
}

// GetPriceImpact if there is no response from the APIs
func (c *CowSwapService) GetPriceImpact(response pb.ExchangeQuoteResponse) string {
	toTokenPrice := c.helper.ConvertStringToFloat64(response.ToTokenPrice)
	fromTokenPrice := c.helper.ConvertStringToFloat64(response.FromTokenPrice)
	priceImpactFloat := (1.0 - (toTokenPrice / fromTokenPrice)) * 100
	priceImpactString := strconv.FormatFloat(priceImpactFloat, 'f', -1, 64)
	if priceImpactString == "NaN" || priceImpactString < "0" {
		priceImpactString = "0.01"
	}
	return priceImpactString
}

// CallPostOrderAPI CallOrderAPI is used to call the cowswap exchange order api
func (c *CowSwapService) CallPostOrderAPI(request *pb.ExchangeSwapExecuteRequest, info config.ChainData) (OrderResponse, error) {
	inputRequestPayload := &pb.CowSwapExecuteRequest{
		SellToken:         request.CowSwapPayload.SellToken,
		BuyToken:          request.CowSwapPayload.BuyToken,
		Receiver:          request.CowSwapPayload.Receiver,
		SellAmount:        request.CowSwapPayload.SellAmount,
		BuyAmount:         request.CowSwapPayload.BuyAmount,
		ValidTo:           request.CowSwapPayload.ValidTo,
		AppData:           request.CowSwapPayload.AppData,
		FeeAmount:         request.CowSwapPayload.FeeAmount,
		Kind:              request.CowSwapPayload.Kind,
		PartiallyFillable: request.CowSwapPayload.PartiallyFillable,
		SellTokenBalance:  request.CowSwapPayload.SellTokenBalance,
		BuyTokenBalance:   request.CowSwapPayload.BuyTokenBalance,
		SigningScheme:     request.CowSwapPayload.SigningScheme,
		Signature:         request.CowSwapPayload.Signature,
		From:              request.CowSwapPayload.From,
		QuoteId:           request.CowSwapPayload.QuoteId,
	}
	jsonReq, _ := json.Marshal(inputRequestPayload)
	reqBody := bytes.NewBuffer(jsonReq)
	var chainName string
	switch info.ChainId {
	case 1:
		chainName = ETHEREUM
	case 100:
		chainName = GNOSIS
	}
	url := c.env.Swap.CowSwapUrl + chainName + "/api/v1/orders"
	body, err := c.httpRequest.PostRequest(url, reqBody)
	if err != nil {
		c.logger.Errorf("Error while calling cowswap exchange POST Order api: %v", err)
		return OrderResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var orderResponse OrderResponse
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		c.logger.Error("Error while unmarshalling the response from cowswap", err)
		return OrderResponse{}, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	return orderResponse, nil
}

// CallGetOrderAPI is used to call the cowswap exchange order api
func (c *CowSwapService) CallGetOrderAPI(orderUID string, info config.ChainData) (OrderStatusResponse, error) {
	var chainName string
	switch info.ChainId {
	case 1:
		chainName = ETHEREUM
	case 100:
		chainName = GNOSIS
	}
	url := c.env.Swap.CowSwapUrl + chainName + "/api/v1/orders/" + orderUID
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Errorf("Error while calling cowswap exchange GET Order api: %v", err)
		return OrderStatusResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var orderStatusResponse OrderStatusResponse
	err = json.Unmarshal(body, &orderStatusResponse)
	if err != nil {
		c.logger.Error("Error while unmarshalling the response from cowswap", err)
		return OrderStatusResponse{}, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	return orderStatusResponse, nil
}

// CallCowTradeAPI is used to call the cowswap Trade api
func (c *CowSwapService) CallCowTradeAPI(orderUID string, info config.ChainData) (CowTradeResponse, error) {
	var chainName string
	switch info.ChainId {
	case 1:
		chainName = ETHEREUM
	case 100:
		chainName = GNOSIS
	}
	url := c.env.Swap.CowSwapUrl + chainName + "/api/v1/trades?orderUid=" + orderUID
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Errorf("Error while calling cowswap exchange GET Trade api: %v", err)
		return CowTradeResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var cowTradeResponse CowTradeResponse
	err = json.Unmarshal(body, &cowTradeResponse)
	if err != nil {
		c.logger.Error("Error while unmarshalling the response from cowswap", err)
		return CowTradeResponse{}, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	return cowTradeResponse, nil
}

// RemoveDuplicates Removes Duplicates from a slice
func (t Tokens) RemoveDuplicates() Tokens {
	encountered := map[string]bool{}
	result := Tokens{}

	for _, token := range t {
		if encountered[token.Address] {
			continue
		}
		encountered[token.Address] = true
		result = append(result, token)
	}
	return result
}

// ToDecimal converts any value to decimal
func ToDecimal(amount interface{}, decimals string) decimal.Decimal {
	decimalInt, err := strconv.Atoi(decimals)
	if err != nil {
		return decimal.Zero
	}
	value := new(big.Int)
	switch v := amount.(type) {
	case string:
		value.SetString(v, 10)
	case float64:
		value = big.NewInt(int64(v))
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimalInt)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ConvertToWei converts any value to wei as big int
func ConvertToWei(value interface{}, decimals string) *big.Int {
	decimalInt, err := strconv.Atoi(decimals)
	if err != nil {
		return big.NewInt(0)
	}
	amount := decimal.NewFromFloat(0)
	switch v := value.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimalInt)))
	result := amount.Mul(mul)

	weiAmount := new(big.Int)
	weiAmount.SetString(result.String(), 10)

	return weiAmount
}

// FetchTokenList fetches the token list from the given url
func (c *CowSwapService) FetchTokenList(url string, exchangeTokens ExchangeCowResponse, info config.ChainData, wg *sync.WaitGroup, chanResponse chan Tokens) {
	defer wg.Done()
	var tokenList Tokens
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Error(err)
		return
	}

	err = json.Unmarshal(body, &exchangeTokens)
	for i := range exchangeTokens.Tokens {
		if exchangeTokens.Tokens[i].ChainID == info.ChainId {
			tokenList = exchangeTokens.Tokens
		}
	}
	if err != nil {
		c.logger.Error(err)
	}
	chanResponse <- tokenList
	return
}

// GetNativeTokenAddress returns the native token address
func (c *CowSwapService) GetNativeTokenAddress(info config.Wallets) string {
	return info.NativeTokenInfo.Address
}

// GetWNativeTokenAddress returns the native token address
func (c *CowSwapService) GetWNativeTokenAddress(info config.ChainData) string {
	return info.CowSwapConfig.WrappedTokenID
}

// IsSupportedChain checks if the chain is supported by cowswap
func (c *CowSwapService) IsSupportedChain(info config.ChainData) bool {
	switch info.ChainId {
	case 1, 100:
		return true
	default:
		return false
	}
}
