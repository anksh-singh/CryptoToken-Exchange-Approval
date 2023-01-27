package _x

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
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

type I0x interface {
	GetExchangeQuote(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, chain string, srcTokenDecimals string, dstTokenDecimals string, slippage string) (*pb.ExchangeQuoteResponse, error)
	GetExchangeSwap(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, srcTokenDecimals string, slippage string, chainConfig config.Wallets) (*pb.ExchangeSwapResponse, error)
}

type OXService struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	coinGecko    coingecko.ICoinGecko
	defaultToken []string
}

func NewOXService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *OXService {
	return &OXService{
		env:          env,
		logger:       logger,
		httpRequest:  httpRequest,
		helper:       helper,
		coinGecko:    coinGecko,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (o *OXService) CallZeroX(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, slippage string) (*Exchange0xModel, error) {
	url := fmt.Sprintf(zeroXChainUrl+"/swap/v1/quote?sellToken=%s&buyToken=%s&sellAmount=%s&affiliateAddress=%s", SellToken, BuyToken, SellAmount, o.env.EVM.Cfg.ZeroxAffilateAddress)
	body, err := o.httpRequest.GetRequest(url)
	if err != nil {
		o.logger.Error(err)
		var ErrorResponse ErrorResponse
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		errorString := string(body)
		return nil, status.Errorf(codes.Internal, errorString, ErrorResponse.Reason)
	}
	var exchangeQuote Exchange0xModel
	err = json.Unmarshal(body, &exchangeQuote)
	if err != nil {
		o.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	return &exchangeQuote, err
}

func (o *OXService) GetExchangeQuote(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, chain string, srcTokenDecimals string, dstTokenDecimals string, slippage string) (*pb.ExchangeQuoteResponse, error) {
	var response pb.ExchangeQuoteResponse
	sellAmountParam, err := o.helper.ConvertStringValueToFloatWei(SellAmount, srcTokenDecimals)
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	sellAmountParam = math.Floor(sellAmountParam)
	_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)
	// in optimism chain 0x4200000000000000000000000000000000000006 is a valid contract
	if chain != "optimism" {
		if contains(o.defaultToken, strings.ToLower(SellToken)) {
			SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(o.defaultToken, strings.ToLower(BuyToken)) {
			BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	// Make API Call to 0x swap
	response0x, err := o.CallZeroX(SellToken, BuyToken, _sellAmountParam, zeroXChainUrl, slippage)
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, err
	}
	resAmountWithDecimals, err := o.helper.ConvertStringValueToBigFloat(response0x.BuyAmount, dstTokenDecimals)
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	o.logger.Infof("***** %v ****", resAmountWithDecimals)
	response.ResAmount = strconv.FormatFloat(resAmountWithDecimals, 'f', -1, 64)

	if response0x.EstimatedPriceImpact == "" || response0x.EstimatedPriceImpact == "0" {
		response0x.EstimatedPriceImpact = "0.01"
	}
	priceImpact, ok := new(big.Float).SetString(response0x.EstimatedPriceImpact)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, errors.New("SetString: error for Guaranteed Price " + response0x.GuaranteedPrice)
	}
	response.PriceImpact = response0x.EstimatedPriceImpact
	resAmount, ok := new(big.Float).SetString(response.ResAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, errors.New("SetString: error resAmount " + response.ResAmount)
	}
	sellAmount, ok := new(big.Float).SetString(SellAmount)
	if !ok {
		return &pb.ExchangeQuoteResponse{}, errors.New("SetString: error sellAmount " + SellAmount)
	}
	// Calculating Price per token := resAmount/sellAmount
	resPricePerToken, _ := new(big.Float).Quo(resAmount, sellAmount).Float64()

	response.ResPricePerFromToken = strconv.FormatFloat(resPricePerToken, 'f', -1, 64)

	// Calculating Price per Token for User = 1/resPricePerToken
	resPricePerToToken, _ := new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
	response.ResPricePerToToken = strconv.FormatFloat(resPricePerToToken, 'f', -1, 64)

	var toQuote float64
	var fromQuote float64
	if strings.ToLower(BuyToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := o.coinGecko.GetTokenExchange("usd", chain)
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote Price for token  request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}

	} else {
		quoteRate, err := o.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(BuyToken), "usd")
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			toQuote = 0
		} else {
			toQuote = quoteRate.Price
		}
	}
	if strings.ToLower(SellToken) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		quoteRate, err := o.coinGecko.GetTokenExchange("usd", chain)

		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			fromQuote = 0
		} else {
			fromQuote = quoteRate.Price
		}

	} else {
		quoteRate, err := o.coinGecko.GetTokenExchangeForContract(chain, strings.ToLower(SellToken), "usd")
		if err != nil {
			o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
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
	response.ApproveAddress = response0x.To

	// Calculating Minimum Receive  resAmount *(1-priceImpact/100)

	minimumReceived, _ := new(big.Float).
		Mul(resAmount, new(big.Float).
			Sub(big.NewFloat(1), new(big.Float).
				Quo(priceImpact, big.NewFloat(100)))).Float64()
	response.MinimumReceived = strconv.FormatFloat(minimumReceived, 'f', -1, 64)

	response.SellToken = SellToken
	response.BuyToken = BuyToken

	return &response, err
}

func (o *OXService) GetExchangeSwap(request *pb.ExchangeSwapRequest, zeroXChainUrl string, srcTokenDecimals string, chainConfig config.Wallets) (*pb.ExchangeSwapResponse, error) {
	sellAmountParam, err := o.helper.ConvertStringValueToFloatWei(request.SellAmount, srcTokenDecimals)
	if err != nil {
		o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, err.Error(), "conversion error")
	}
	sellAmountParam = math.Floor(sellAmountParam)
	_sellAmountParam := strconv.FormatFloat(sellAmountParam, 'f', -1, 64)

	var response pb.ExchangeSwapResponse
	// Calculating the SellAmount in native token := srcAmount/10^decimalValues
	resSplippage := ""
	if request.Slippage != "" {
		slippagePercent, ok := new(big.Float).SetString(request.Slippage)
		if !ok {
			return &pb.ExchangeSwapResponse{}, errors.New("SetString: error for Slippage : " + request.Slippage)
		}
		_resSplippage, _ := new(big.Float).Quo(slippagePercent, big.NewFloat(100)).Float64()
		resSplippage = strconv.FormatFloat(_resSplippage, 'f', -1, 64)
	}

	// in optimism chain 0x4200000000000000000000000000000000000006 is a valid contract
	if request.Chain != "optimism" {
		if contains(o.defaultToken, strings.ToLower(request.SellToken)) {
			request.SellToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(o.defaultToken, strings.ToLower(request.BuyToken)) {
			request.BuyToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	response0x, err := o.CallZeroX(request.SellToken, request.BuyToken, _sellAmountParam, zeroXChainUrl, fmt.Sprint(resSplippage))
	if err != nil {
		o.logger.Errorf("Error for Exchange Swap request  is : %v", err.Error())
		return &pb.ExchangeSwapResponse{}, err
	}
	response.To = response0x.To
	response.Data = response0x.Data
	if contains(o.defaultToken, strings.ToLower(request.SellToken)) {
		response.Value = _sellAmountParam
	} else {
		response.Value = "0"
	}
	gasLimitFromOx := ""
	if response0x.EstimatedGas != "" {
		gasLimit, ok := new(big.Float).SetString(response0x.EstimatedGas)
		if !ok {
			return &pb.ExchangeSwapResponse{}, errors.New("SetString: error for EstimatedGas " + response0x.EstimatedGas)
		}
		// Calculating gasLimit :=( gasLimit + gasLimit*0.75)
		var gasLimit0x float64
		if chainConfig.GasLimitFactor.Zerox != 0 {
			gasLimit0x, _ = new(big.Float).Mul(gasLimit, new(big.Float).SetFloat64(chainConfig.GasLimitFactor.Zerox)).Float64()
			//gasLimit0x, _ = new(big.Float).Add(gasLimit, new(big.Float).Mul(gasLimit, big.NewFloat(0.75))).Float64()
		} else {
			gasLimit0x, _ = new(big.Float).Add(gasLimit, new(big.Float).Mul(gasLimit, big.NewFloat(0.8))).Float64()
		}
		gasLimitFromOx = fmt.Sprint(int(gasLimit0x))
	} else {
		gasLimitFromOx = "600000"
	}
	response.GasLimit = gasLimitFromOx
	response.Gas = response0x.GasPrice
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
