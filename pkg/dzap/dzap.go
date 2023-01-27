package dzap

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strconv"
	"strings"
)

type IDZap interface {
	GetExchangeTokens(info config.Wallets) (*pb.ExchangeTokenResponse, error)
	GetExchangeQuote(in *pb.ExchangeMultiQuoteRequest, info config.Wallets, decimalMapper map[string]string) (*pb.ExchangeMultiQuoteResponse, error)
	GetExchangeSwap(in *pb.ExchangeMultiSwapRequest, info config.Wallets, decimalMapper map[string]string) (*pb.ExchangeMultiSwapResponse, error)
}

type ServiceDZap struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	ethgoRpc     map[string]*ethgoJsonRPC.Client
	defaultToken []string
	utils        *utils.UtilConf
}

func NewServiceDZap(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *ServiceDZap {
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	utilsManager := utils.NewUtils(logger, env)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported EVM chains:", len(env.EVM.Cfg.Wallets))
	if len(env.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
	}
	//Initialize EVM RPC configurations
	for i, w := range env.EVM.Cfg.Wallets {
		i++
		var err error
		ethgoRpc[w.ChainName], err = ethgoJsonRPC.NewClient(w.RPC)
		if err != nil {
			logger.Errorf(err.Error())
			//logger.Fatalf("Error initializing go RPC client for `%s` chain", w.ChainName)
		}
	}

	return &ServiceDZap{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		ethgoRpc:     ethgoRpc,
		utils:        utilsManager,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (d *ServiceDZap) GetExchangeTokens(info config.Wallets) (*pb.ExchangeTokenResponse, error) {
	url := fmt.Sprintf(d.env.Swap.DZapUrl+"/"+"token/get-all?chainId=%v", info.ChainID)
	body, err := d.httpRequest.GetRequest(url)
	if err != nil {
		d.logger.Error(err)
		var ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		errorString := string(body)
		return nil, status.Errorf(codes.Internal, errorString, "internal error")
	}

	var responseDZap []*ExchangeTokenResponse
	var response pb.ExchangeTokenResponse
	err = json.Unmarshal(body, &responseDZap)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var nativeTokenInfo pb.ExchangeTokenInfo
	for _, item := range responseDZap {
		address := item.Contract
		if contains(d.defaultToken, strings.ToLower(item.Contract)) {
			address = d.FindNativeTokenAddress(info.ChainName, item.Contract)
			nativeTokenInfo.TokenName = item.Name
			nativeTokenInfo.TokenSymbol = item.Symbol
			nativeTokenInfo.TokenAddress = address
			nativeTokenInfo.TokenDecimals = fmt.Sprint(item.Decimals)
			nativeTokenInfo.TokenLogoUrl = item.Logo
			nativeTokenInfo.LogoUrl = d.env.Swap.DZapLogoUrl
		}
		if !contains(d.defaultToken, strings.ToLower(item.Contract)) {
			response.ExchangeTokens = append(response.ExchangeTokens, &pb.ExchangeTokenInfo{
				TokenAddress:  address,
				TokenDecimals: fmt.Sprint(item.Decimals),
				TokenSymbol:   item.Symbol,
				TokenName:     item.Name,
				TokenLogoUrl:  item.Logo,
				LogoUrl:       d.env.Swap.DZapLogoUrl,
			})
		}

	}
	//prepending native token
	response.ExchangeTokens = append([]*pb.ExchangeTokenInfo{
		&nativeTokenInfo,
	}, response.ExchangeTokens...)

	return &response, err
}

func (d *ServiceDZap) GetExchangeQuote(in *pb.ExchangeMultiQuoteRequest, info config.Wallets, decimalMapper map[string]string) (*pb.ExchangeMultiQuoteResponse, error) {
	// create pay load
	var getPathPayload ExchangePathRequest
	var getParamPayload ExchangeParamsRequest
	getParamPayload.ChainId = info.ChainID
	getPathPayload.ChainId = info.ChainID
	for _, item := range in.MultiChainRequests {
		amount, err := d.helper.ConvertToWei(item.SellAmount, decimalMapper[item.SellToken])
		if err != nil {
			d.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "internal error")
		}
		//inputRequest.SellAmount = amount.String()
		if contains(d.defaultToken, strings.ToLower(item.BuyToken)) {
			item.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(d.defaultToken, strings.ToLower(item.SellToken)) {
			item.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		getPathPayload.Requests = append(getPathPayload.Requests, &Requests{
			Amount:           amount.String(),
			FromTokenAddress: item.SellToken,
			ToTokenAddress:   item.BuyToken,
			Slippage:         item.Slippage,
		})
		getParamPayload.SwapParams = append(getParamPayload.SwapParams, &SwapParams{
			Amount:           amount.String(),
			FromTokenAddress: item.SellToken,
			ToTokenAddress:   item.BuyToken,
			Slippage:         item.Slippage,
		})
	}

	// Make API call to DZap server

	buildPathRequestData, err := json.Marshal(getPathPayload)
	if err != nil {
		d.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshall error")
	}
	urlPathRequest := fmt.Sprintf(d.env.Swap.DZapUrl + "/swap/get-path")
	buildParamRequestData, err := json.Marshal(getParamPayload)
	if err != nil {
		d.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshall error")
	}
	urlParamRequest := fmt.Sprintf(d.env.Swap.DZapUrl + "/swap/get-params")
	// Make Put Request getPAth
	bodyPath, err := d.httpRequest.PutRequest(urlPathRequest, string(buildPathRequestData))
	if err != nil {
		d.logger.Error(err)
		var ErrorResponse ErrorResponse
		err = json.Unmarshal(bodyPath, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Unavailable, string(bodyPath), ErrorResponse.Name)
	}
	var successErrorResponse []ErrorResponseV1
	err = json.Unmarshal(bodyPath, &successErrorResponse)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	if successErrorResponse[0].Status == "Error" {
		d.logger.Error(err)
		statusCode := codes.Internal
		if successErrorResponse[0].Data.Errors.StatusCode == 400 {
			statusCode = codes.InvalidArgument
		}
		return nil, status.Errorf(statusCode, successErrorResponse[0].Data.Errors.Description, "invalid request")
	}
	var pathResponse []*MultiExchangePath
	err = json.Unmarshal(bodyPath, &pathResponse)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	// Make Put Request getParam
	bodyParam, err := d.httpRequest.PutRequest(urlParamRequest, string(buildParamRequestData))
	if err != nil {
		d.logger.Error(err)
		var ErrorResponse ErrorResponse
		err = json.Unmarshal(bodyParam, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(bodyParam), ErrorResponse.Name)
	}
	var paramResponse ExchangeParamsResponse
	err = json.Unmarshal(bodyParam, &paramResponse)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// Calculate Quote Data
	var quoteResponse pb.ExchangeMultiQuoteResponse
	quoteResponse.Chain = info.ChainName
	for _, item := range pathResponse {
		var multiChainResponse pb.MultiChainResponse
		//multiChainResponse.ResAmount = item.Data.ToTokenAmount

		actualResAmount := d.helper.CalculateRateWithDecimal(item.Data.ToTokenAmount, int64(item.Data.ToToken.Decimals))
		multiChainResponse.ResAmount = strconv.FormatFloat(actualResAmount, 'f', -1, 64)
		resAmount, ok := new(big.Float).SetString(multiChainResponse.ResAmount)
		if !ok {
			return &pb.ExchangeMultiQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error resAmount "+multiChainResponse.ResAmount, "type error")
		}
		actualSellAmount := d.helper.CalculateRateWithDecimal(item.Data.FromTokenAmount, int64(item.Data.FromToken.Decimals))
		sellAmount := big.NewFloat(actualSellAmount)
		multiChainResponse.ApproveAddress = d.GetSwapDetail(info.ChainName).ApproveAddress
		fromQuote, toQuote := d.CalculateQuotePrice(item.Data.FromToken.Address, item.Data.ToToken.Address, info.ChainName)
		fromTokenPrice, _ := new(big.Float).Mul(big.NewFloat(fromQuote), sellAmount).Float64()
		toTokenPrice, _ := new(big.Float).Mul(big.NewFloat(toQuote), resAmount).Float64()
		multiChainResponse.ToTokenPrice = strconv.FormatFloat(toTokenPrice, 'f', -1, 64)
		multiChainResponse.FromTokenPrice = strconv.FormatFloat(fromTokenPrice, 'f', -1, 64)
		// Calculating Price per token := resAmount/sellAmount
		var resPricePerToken float64
		if actualSellAmount == 0 {
			resPricePerToken = 0
		} else {
			resPricePerToken, _ = new(big.Float).Quo(resAmount, sellAmount).Float64()
		}
		multiChainResponse.ResPricePerFromToken = strconv.FormatFloat(resPricePerToken, 'f', -1, 64)
		// Calculating Price per Token for User = 1/resPricePerToken
		var resPricePerToToken float64
		if resPricePerToken == 0 {
			resPricePerToToken = 0
		} else {
			resPricePerToToken, _ = new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
		}
		multiChainResponse.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)
		fromToken := item.Data.FromToken.Address
		toToken := item.Data.ToToken.Address
		if contains(d.defaultToken, strings.ToLower(fromToken)) {
			fromToken = d.FindNativeTokenAddress(info.ChainName, fromToken)
		}
		if contains(d.defaultToken, strings.ToLower(toToken)) {
			fromToken = d.FindNativeTokenAddress(info.ChainName, toToken)
		}
		multiChainResponse.SellToken = fromToken
		multiChainResponse.BuyToken = toToken
		// Calculate minimum received

		for _, minItem := range paramResponse.ErcSwapDetails {
			if strings.ToLower(minItem.Desc.DstToken) == item.Data.ToToken.Address && strings.ToLower(minItem.Desc.SrcToken) == item.Data.FromToken.Address {
				minRecieved := d.helper.CalculateRateWithDecimal(minItem.Desc.MinReturnAmount, int64(item.Data.ToToken.Decimals))
				multiChainResponse.MinimumReceived = strconv.FormatFloat(minRecieved, 'f', -1, 64)
			}
		}
		multiChainResponse.PriceImpact = d.GetPriceImpact(&multiChainResponse)
		quoteResponse.MultiChainResponse = append(quoteResponse.MultiChainResponse, &multiChainResponse)
	}
	return &quoteResponse, err
}

func (d *ServiceDZap) GetExchangeSwap(in *pb.ExchangeMultiSwapRequest, info config.Wallets, decimalMapper map[string]string) (*pb.ExchangeMultiSwapResponse, error) {
	var getParamPayload ExchangeParamsRequest
	getParamPayload.ChainId = info.ChainID
	for _, item := range in.MultiChainRequests {
		if contains(d.defaultToken, strings.ToLower(item.BuyToken)) {
			item.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(d.defaultToken, strings.ToLower(item.SellToken)) {
			item.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		amount, err := d.helper.ConvertToWei(item.SellAmount, decimalMapper[item.SellToken])
		if err != nil {
			d.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "internal error")
		}
		getParamPayload.SwapParams = append(getParamPayload.SwapParams, &SwapParams{
			Amount:           amount.String(),
			FromTokenAddress: item.SellToken,
			ToTokenAddress:   item.BuyToken,
			Slippage:         item.Slippage,
		})
	}
	// Make API call to DZap server
	buildParamRequestData, err := json.Marshal(getParamPayload)
	if err != nil {
		d.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshall error")
	}
	urlParamRequest := fmt.Sprintf(d.env.Swap.DZapUrl + "/swap/get-params")
	// Make Put Request getPAth
	bodyParam, err := d.httpRequest.PutRequest(urlParamRequest, string(buildParamRequestData))
	if err != nil {
		d.logger.Error(err)
		var ErrorResponse ErrorResponse
		err = json.Unmarshal(bodyParam, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(bodyParam), ErrorResponse.Name)
	}
	var paramResponse ExchangeParamsResponse
	err = json.Unmarshal(bodyParam, &paramResponse)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var swapResponse pb.ExchangeMultiSwapResponse
	swapResponse.To = d.GetSwapDetail(info.ChainName).ApproveAddress
	swapResponse.Value = paramResponse.Value
	//Calculated by calling the contract

	abiContract, err := abi.NewABI(CONTRACT_ABI)
	addr := ethgo.HexToAddress(swapResponse.To) //contract address
	c := contract.NewContract(addr, abiContract, contract.WithJsonRPC(d.ethgoRpc[info.ChainName].Eth()))
	method := c.GetABI().GetMethod("swapTokensToTokens")
	if method == nil {
		d.logger.Error("Method: swapTokensToTokens not found")
	}

	// Convert to ERCSWAPDETAILSPAYLOAD

	var ercListMap []map[string]interface{}
	for _, item := range paramResponse.ErcSwapDetails {
		n := new(big.Int)
		ercSwapDetails := make(map[string]interface{})
		ercSwapDetails["executor"] = ethgo.HexToAddress(item.Executor)
		ercSwapDetails["minReturnAmount"], _ = n.SetString(item.MinReturnAmount, 10)
		ercSwapDetails["routeData"], _ = json.Marshal(item.RouteData)
		ercSwapDetails["permit"], _ = json.Marshal(item.Permit)
		ercSwapDetails["desc"] = make(map[string]interface{})
		ercSwapDetails["desc"].(map[string]interface{})["srcToken"] = ethgo.HexToAddress(item.Desc.SrcToken)
		ercSwapDetails["desc"].(map[string]interface{})["dstToken"] = ethgo.HexToAddress(item.Desc.DstToken)
		ercSwapDetails["desc"].(map[string]interface{})["permit"], _ = json.Marshal(item.Desc.Permit)
		ercSwapDetails["desc"].(map[string]interface{})["minReturnAmount"], _ = n.SetString(item.Desc.MinReturnAmount, 10)
		ercSwapDetails["desc"].(map[string]interface{})["amount"], _ = n.SetString(item.Desc.Amount, 10)
		ercSwapDetails["desc"].(map[string]interface{})["flags"], _ = n.SetString(item.Desc.Flags, 10)
		ercSwapDetails["desc"].(map[string]interface{})["dstReceiver"] = ethgo.HexToAddress(item.Desc.DstReceiver)
		ercSwapDetails["desc"].(map[string]interface{})["srcReceiver"] = ethgo.HexToAddress(item.Desc.SrcReceiver)
		ercListMap = append(ercListMap, ercSwapDetails)

	}
	data, err := method.Encode(map[string]interface{}{
		"data_":      ercListMap,
		"recipient_": ethgo.HexToAddress(in.TakerAddress),
		"nftId_":     big.NewInt(int64(0)),
	})
	if err != nil {
		d.logger.Error(err)
	}
	// Calculate Gas limit
	client, err := ethgoJsonRPC.NewClient(info.RPC)
	if err != nil {
		d.logger.Error("Error creating  rpc client")
	}
	to := ethgo.HexToAddress(swapResponse.To)
	gasLimit, err := client.Eth().EstimateGas(&ethgo.CallMsg{
		To:   &to,
		From: ethgo.HexToAddress(in.MultiChainRequests[0].SellToken),
		Data: data,
	})
	if err != nil {
		d.logger.Error("Error fetching gas estimate")
		gasLimit = 8000000 //TODO:To be refactored
		err = nil
	}
	gasLimit = gasLimit * 31
	swapResponse.GasLimit = fmt.Sprint(gasLimit)
	swapResponse.Data = fmt.Sprintf("0x%x", data)
	swapResponse.Gas = fmt.Sprint(gasLimit)
	swapResponse.Value = paramResponse.Value
	swapResponse.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", swapResponse.To, swapResponse.Value, swapResponse.Data, swapResponse.GasLimit)
	return &swapResponse, err
}

func (d *ServiceDZap) CalculateQuotePrice(fromTokenAddress string, toTokenAddress string, chain string) (float64, float64) {
	var toQuote float64
	var fromQuote float64
	if contains(d.defaultToken, strings.ToLower(fromTokenAddress)) {
		quoteRate, err := d.coinGecko.GetTokenExchange("usd", chain)
		if err != nil {
			d.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}
	} else {
		quoteRate, err := d.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(fromTokenAddress), "usd")
		if err != nil {
			d.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}

	}
	if contains(d.defaultToken, strings.ToLower(toTokenAddress)) {
		quoteRate, err := d.coinGecko.GetTokenExchange("usd", chain)
		if err != nil {
			d.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}

	} else {
		quoteRate, err := d.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(toTokenAddress), "usd")
		if err != nil {
			d.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}
	}
	return fromQuote, toQuote
}

func (d *ServiceDZap) GetSwapDetail(chain string) *config.DZapSwapConfig {
	for _, item := range d.env.Swap.ChainData {
		if item.ChainName == chain {
			return &item.DZapSwapConfig
		}
	}
	return nil
}

func (d *ServiceDZap) GetPriceImpact(response *pb.MultiChainResponse) string {
	var priceImpactString string
	toTokenPrice := d.helper.ConvertStringToFloat64(response.ToTokenPrice)
	fromTokenPrice := d.helper.ConvertStringToFloat64(response.FromTokenPrice)
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

func (d *ServiceDZap) FindNativeTokenAddress(toChainId string, tokenAddress string) string {
	for _, item := range d.env.EVM.Cfg.Wallets {
		if strings.ToLower(toChainId) == strings.ToLower(item.ChainName) {
			if item.NativeTokenInfo.Address != "" {
				return item.NativeTokenInfo.Address
			}
		} else {
			chainId := strconv.FormatInt(int64(item.ChainID), 10)
			if toChainId == chainId {
				if item.NativeTokenInfo.Address != "" {
					return item.NativeTokenInfo.Address
				}
			}
		}
	}
	return tokenAddress
}
