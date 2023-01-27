package dodo

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"github.com/onrik/ethrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type DODO interface {
	GetExchangeQuote(SellToken string, BuyToken string, SellAmount string, chain string, chainId string, srcTokenDecimals string, dstTokenDecimals string, slippage string, userAddr string, approveAddress string) (*pb.ExchangeQuoteResponse, error)
	GetExchangeSwap(SellToken string, BuyToken string, SellAmount string, chain string, chainId string, srcTokenDecimals string, dstTokenDecimals string, slippage string, userAddr string, chainConfig config.Wallets) (*pb.ExchangeSwapResponse, error)
}

type ServiceDodo struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	rpc          map[string]*ethrpc.EthRPC
	defaultToken []string
	utils        *utils.UtilConf
}

const (
	RouteAPIStatusKey   = "status"
	RouteAPIDataKey     = "data"
	SuccessResponseCode = 200
)

func NewServiceDodo(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *ServiceDodo {
	rpc := make(map[string]*ethrpc.EthRPC)
	utilsManager := utils.NewUtils(logger, env)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported EVM chains:", len(env.EVM.Cfg.Wallets))
	if len(env.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
	}
	//Initialize EVM RPC configurations
	for i, w := range env.EVM.Cfg.Wallets {
		i++
		rpc[w.ChainName] = ethrpc.New(w.RPC)
	}

	return &ServiceDodo{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		rpc:          rpc,
		utils:        utilsManager,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (s *ServiceDodo) CallDoDo(fromTokenAddress string, fromTokenDecimals string, toTokenAddress string, toTokenDecimals string, fromAmount string, dodoUrl string, slippage string, userAddr string, chainId string) (*ExchangeModel, error) {
	url := fmt.Sprintf(dodoUrl+"/dodoapi/getdodoroute?"+"fromTokenAddress=%s&fromTokenDecimals=%s&toTokenAddress=%s&toTokenDecimals=%s&fromAmount=%s&slippage=%s&userAddr=%s&chainId=%s", fromTokenAddress, fromTokenDecimals, toTokenAddress, toTokenDecimals, fromAmount, "1", userAddr, chainId)
	body, err := s.httpRequest.GetRequest(url)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
	}

	//Work around to deal with dodo's unconventional way to log error messages in data field
	//It is inefficient to unmarshal a response twice. TODO: Enhance processing unstructured data
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err == nil && result[RouteAPIStatusKey].(float64) != SuccessResponseCode {
		return nil, status.Errorf(codes.InvalidArgument, result[RouteAPIDataKey].(string), "json error")
	}

	var exchangeQuote ExchangeModel
	err = json.Unmarshal(body, &exchangeQuote)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshaling error")
	}
	return &exchangeQuote, err
}

func (s *ServiceDodo) GetExchangeQuote(SellToken string, BuyToken string, SellAmount string, chain string, chainId string, srcTokenDecimals string, dstTokenDecimals string, slippage string, userAddr string, dodoUrl string, approveAddress string) (*pb.ExchangeQuoteResponse, error) {
	var response pb.ExchangeQuoteResponse
	chainSupported := s.GetSupportedChains(chain)
	if chainSupported {
		sellAmountParam, err := s.helper.ConvertStringValueToFloatWei(SellAmount, srcTokenDecimals)
		sellAmountParam = math.Floor(sellAmountParam)
		_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)
		if err != nil {
			s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "format error")
		}
		if contains(s.defaultToken, strings.ToLower(SellToken)) {
			SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(s.defaultToken, strings.ToLower(BuyToken)) {
			BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		// Make API Call to Dodo swap
		responseDodo, err := s.CallDoDo(SellToken, srcTokenDecimals, BuyToken, dstTokenDecimals, _sellAmountParam, dodoUrl, slippage, userAddr, chainId)
		if err != nil {
			s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return &pb.ExchangeQuoteResponse{}, err
		}
		response.ResAmount = strconv.FormatFloat(responseDodo.Data.ResAmount, 'f', -1, 64)
		response.PriceImpact = strconv.FormatFloat(responseDodo.Data.PriceImpact*100, 'f', -1, 64)
		response.ResPricePerFromToken = strconv.FormatFloat(responseDodo.Data.ResPricePerFromToken, 'f', -1, 64)
		response.ResPricePerToToken = strconv.FormatFloat(responseDodo.Data.ResPricePerToToken, 'f', -1, 64)
		var toQuote float64
		var fromQuote float64
		if strings.ToLower(BuyToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			quoteRate, err := s.coinGecko.GetTokenExchange("usd", chain)
			if err != nil {
				s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				toQuote = 0
			} else {
				floatBuy, err := strconv.ParseFloat(response.ResAmount, 64)
				if err != nil {
					s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				}
				toQuote = quoteRate.Price * floatBuy
			}
		} else {
			quoteRate, err := s.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(BuyToken), "usd")
			if err != nil {
				s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				toQuote = 0
			} else {
				floatBuy, err := strconv.ParseFloat(response.ResAmount, 64)
				if err != nil {
					s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				}
				toQuote = quoteRate.Price * floatBuy
			}

		}
		if strings.ToLower(SellToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			quoteRate, err := s.coinGecko.GetTokenExchange("usd", chain)

			if err != nil {
				s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				fromQuote = 0
			} else {
				floatSell, err := strconv.ParseFloat(SellAmount, 64)
				if err != nil {
					s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				}
				fromQuote = quoteRate.Price * floatSell
			}

		} else {
			quoteRate, err := s.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(SellToken), "usd")
			if err != nil {
				s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				fromQuote = 0
			} else {
				floatSell, err := strconv.ParseFloat(SellAmount, 64)
				if err != nil {
					s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				}
				fromQuote = quoteRate.Price * floatSell
			}
		}

		resAmount, ok := new(big.Float).SetString("0")
		if response.ResAmount != "" {
			resAmount, ok = new(big.Float).SetString(response.ResAmount)
			if !ok {
				return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error resAmount "+response.ResAmount, "type error")
			}
		}

		_slippage, ok := new(big.Float).SetString("0.01")
		if !ok {
			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error resAmount "+response.ResAmount, "type error")
		}
		response.SellToken = SellToken
		response.BuyToken = BuyToken
		response.FromTokenPrice = strconv.FormatFloat(fromQuote, 'f', -1, 64)
		response.ToTokenPrice = strconv.FormatFloat(toQuote, 'f', -1, 64)
		response.ApproveAddress = approveAddress
		// Calculating Minimum Receive  resAmount *(1-priceImpact/100)
		minimumReceived, _ := new(big.Float).
			Mul(resAmount, new(big.Float).
				Sub(big.NewFloat(1), new(big.Float).
					Quo(_slippage, big.NewFloat(100)))).Float64()
		response.MinimumReceived = strconv.FormatFloat(minimumReceived, 'f', -1, 64)
		return &response, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "Chain not supported", "Chain not supported")
	}
}

func (s *ServiceDodo) GetExchangeSwap(request *pb.ExchangeSwapRequest, chainId string, srcTokenDecimals string, dstTokenDecimals string, dodoUrl string, chainConfig config.Wallets) (*pb.ExchangeSwapResponse, error) {
	sellAmountParam, err := s.helper.ConvertStringValueToFloatWei(request.SellAmount, srcTokenDecimals)
	sellAmountParam = math.Floor(sellAmountParam)
	_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)

	var response pb.ExchangeSwapResponse
	// Calculating the SellAmount in native token := srcAmount/10^decimalValues
	if contains(s.defaultToken, strings.ToLower(request.SellToken)) {
		request.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	}
	if contains(s.defaultToken, strings.ToLower(request.BuyToken)) {
		request.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	}
	responseDodo, err := s.CallDoDo(request.SellToken, srcTokenDecimals, request.BuyToken, dstTokenDecimals, _sellAmountParam, dodoUrl, request.Slippage, request.TakerAddress, chainId)
	if err != nil {
		s.logger.Errorf("Error for Exchange Swap request  is : %v", err.Error())
		return &pb.ExchangeSwapResponse{}, err
	}
	//response.From = SellToken
	response.To = responseDodo.Data.To
	response.Data = responseDodo.Data.Data
	if strings.ToLower(request.SellToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		response.Value = _sellAmountParam
	} else {
		response.Value = "0"
	}

	// Calculate Gas Limit
	gasLimitData := responseDodo.Data.Data
	gaslimitTo := responseDodo.Data.To
	gaslimitFrom := request.TakerAddress

	// Something wrong with Params
	if request.Chain == "avalanche" && gaslimitTo == "0xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" {
		gaslimitTo = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	}

	if request.Chain == "harmony" {
		gaslimitTo = s.utils.ResolveBech32Address(gaslimitTo)
		gaslimitFrom = s.utils.ResolveBech32Address(gaslimitFrom)
	}

	if request.Chain == "xinfin" {
		gaslimitTo = s.utils.ResolveXDCAddress(gaslimitTo)
		gaslimitFrom = s.utils.ResolveXDCAddress(gaslimitFrom)
		if strings.ToLower(request.BuyToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			gaslimitTo = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
			gasLimitData = "0x"
		}
	}

	if request.Chain == "boba" && strings.ToLower(gaslimitTo) == "0x4200000000000000000000000000000000000006" {
		gaslimitTo = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		gasLimitData = "0x"
	}

	transaction := ethrpc.T{
		From: gaslimitFrom,
		To:   gaslimitTo,
		Data: gasLimitData,
	}
	gasLimit, err := s.rpc[request.Chain].EthEstimateGas(transaction)

	if err != nil {
		s.logger.Error("Error fetching gas estimate")
		gasLimit = 400000 //TODO:To be refactored
		err = nil
	}
	gasPrice, err := s.rpc[request.Chain].EthGasPrice()
	if err != nil {
		s.logger.Error("Error fetching gas price")
		return nil, status.Errorf(codes.Internal, "Error fetching gas price", "Internal Error")
	}

	_gaslimit := float64(gasLimit)
	gasLimitDodo, _ := new(big.Float).Mul(new(big.Float).SetFloat64(_gaslimit), new(big.Float).SetFloat64(chainConfig.GasLimitFactor.Dodo)).Float64()

	if chainConfig.GasLimitFactor.Dodo == 1 {
		gasLimitDodo, _ = new(big.Float).Mul(new(big.Float).SetFloat64(_gaslimit), new(big.Float).SetFloat64(1.1)).Float64()
	}
	response.GasLimit = strconv.FormatFloat(gasLimitDodo, 'f', -1, 64)
	response.Gas = fmt.Sprint(&gasPrice)
	response.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", response.To, response.Value, response.Data, response.GasLimit)
	return &response, err

}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func (s *ServiceDodo) GetSupportedChains(chain string) bool {
	switch chain {
	case "ethereum", "heco", "boba", "aurora", "arbitrum", "moonriver":
		return true
	default:
		return false
	}
}
