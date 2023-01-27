package zksync

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/jsonrpc"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"github.com/onrik/ethrpc"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var currencySymbol = make(map[string]string)
var currencyCode = make(map[string]string)

type IZksync interface {
	GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error)
	//GetTokenPriceV2(request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error)
	GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error)
	//GetNonce(request *pb.NonceRequest) (*pb.NonceResponse, error)
	//GasLimit(request *pb.GasLimitRequest) (*pb.GasLimitResponse, error)
	//SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error)
	ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error)
	//GetUserData(request *pb.UserDataRequest) (*pb.UserDataResponse, error)
	//GetProcessingFee(request *pb.ProcessingFeeRequest) (*pb.ProcessingFeeResponse, error)
	//GetTxStatus(request *pb.TxStatusRequest) (*pb.TxStatusResponse, error)
	//GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error)
	//TokenApprove(request *pb.ApprovalRequest) (*pb.ApprovalResponse, error)
}

type ServiceZksync struct {
	rpc    map[string]*ethrpc.EthRPC
	env    *config.Config
	logger *zap.SugaredLogger
	util   *utils.UtilConf

	ethgoRpc    map[string]*ethgoJsonRPC.Client
	httpRequest utils.IHttpRequest
	rpcHandler  *jsonrpc.RPCHandler
	RestAPI     string
	coingecko   coingecko.ICoinGecko
	helper      utils.Helpers
}

func NewServiceZksync(config *config.Config, logger *zap.SugaredLogger, gecko coingecko.ICoinGecko) *ServiceZksync {
	rpc := make(map[string]*ethrpc.EthRPC)
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported EVM chains:", len(config.EVM.Cfg.Wallets))
	if len(config.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
	}
	restApi := ""
	//Initialize EVM RPC configurations
	for i, w := range config.EVM.Cfg.Wallets {
		if w.ChainName == "zksync" {
			restApi = w.RestAPI
		}
		i++
		rpc[w.ChainName] = ethrpc.New(w.RPC)
		var err error
		ethgoRpc[w.ChainName], err = ethgoJsonRPC.NewClient(w.RPC)
		if err != nil {
			logger.Errorf(err.Error())
			logger.Fatalf("Error initializing go RPC client for `%s` chain", w.ChainName)
		}
		logger.Infof("%v. %v EVM chain initialized with configurations %v", i, w.ChainName, w)
	}
	utilsManager := utils.NewUtils(logger, config)
	httpRequest := utils.NewHttpRequest(logger)
	rpcHandler := jsonrpc.NewJsonRPCHandler(config, logger, httpRequest)
	return &ServiceZksync{rpc, config, logger, utilsManager, ethgoRpc,
		httpRequest, rpcHandler, restApi, gecko, utils.Helpers{}}
}

func (z *ServiceZksync) GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error) {
	if tokenAddress == "" {
		tokenAddress = "0x0000000000000000000000000000000000000000"
	}
	url := z.RestAPI + "/token/" + tokenAddress
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var tokenResponse TokenQuoteResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]`).MatchString
	var symbol string
	const DefaultSymbol = "$"
	const DefaultCurrency = "USD"
	if !IsLetter(currency) { //Check if currency is symbol
		if symbolCode, ok := currencySymbol[currency]; ok {
			symbol = currency
			currency = symbolCode
		} else {
			//Defaults to USD currency & $symbol
			symbol = DefaultSymbol
			currency = DefaultCurrency
		}
	} else {
		if curSymbol, ok := currencyCode[strings.ToLower(currency)]; ok {
			symbol = curSymbol
		} else {
			//Defaults to $symbol
			symbol = DefaultSymbol
			currency = DefaultCurrency
		}
	}
	currency = strings.ToLower(currency)
	var price float64

	quote := tokenResponse.UsdPrice
	if quote != "" {
		price, err = strconv.ParseFloat(quote, 8)
		if err != nil {
			z.logger.Errorf("quote parsing error")
			price = 0
		}
	}
	return &pb.TokenPriceResponse{
		Price:          price,
		CurrencyCode:   strings.ToUpper(currency),
		CurrencySymbol: symbol,
	}, nil

}

//func (z *ServiceZksync) GetPriceUsd() (map[string]interface{}, error) {
//	url := fmt.Sprintf("https://api.zksync.io/api/v0.2/tokens/ETH/priceIn/usd")
//	body, err := z.httpRequest.GetRequest(url)
//	if err != nil {
//		z.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//	var response TokenPriceUSD
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		z.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//	if len(response.Result.)
//}

func (z *ServiceZksync) GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	url := fmt.Sprintf("https://zksync2-testnet.zkscan.io/api?module=%v&action=%v&address=%v", "account", "tokenlist", request.Address)
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var balanceResponse BalanceResponse
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	//info := z.util.GetWalletInfo(request.Chain)
	//rpc := info.RPC
	//client, err := ethgoJsonRPC.NewClient(rpc)
	if err != nil {
		panic(err)
	}
	//balance, err := client.Eth().GetBalance(ethgo.HexToAddress(request.Address), ethgo.Latest)
	//nativeTokenBalance := fmt.Sprintf("%v", balance)
	//newNativeTokenBalance, err := strconv.ParseFloat(nativeTokenBalance, 64)
	//var quotePrice = 0.0
	//var quotePriceChange24H = ""
	//var quotePricePercentageChange24H = 0.0
	//var quoteRate = 0.0
	//tokenInfo, err := z.GetTokenPrice("usd", item.ContractAddress)
	//if err == nil {
	//	quotePrice = newNativeTokenBalance / math.Pow(10, float64(18)) * tokenInfo.Token.MarketData.CurrentPrice.Usd
	//	quoteRate = tokenInfo.Token.MarketData.CurrentPrice.Usd
	//	quotePriceChange24H = strconv.FormatFloat(tokenInfo.Token.MarketData.PriceChange24H, 'f', 2, 64)
	//	quotePricePercentageChange24H = tokenInfo.Token.MarketData.PriceChangePercentage24HInCurrency.Usd
	//}

	var assetResponse pb.BalanceResponse
	assetResponse.Token = make([]*pb.TokenBalance, 0)
	if len(balanceResponse.Result) > 0 {
		for _, item := range balanceResponse.Result {
			if item.ContractAddress == "0x000000000000000000000000000000000000800a" {
				var quotePrice = 0.0
				var quotePriceChange24H = ""
				var quotePricePercentageChange24H = 0.0
				var quoteRate = 0.0
				tokenInfo, err := z.GetTokenPrice("usd", "")
				if err == nil {
					balance, _ := strconv.ParseFloat(item.Balance, 64)
					quotePrice = balance / math.Pow(10, float64(18)) * tokenInfo.Price
					quoteRate = tokenInfo.Price
					//quotePriceChange24H = strconv.FormatFloat(tokenInfo.Token.MarketData.PriceChange24H, 'f', 2, 64)
					//quotePricePercentageChange24H = tokenInfo.Token.MarketData.PriceChangePercentage24HInCurrency.Usd
				}
				assetResponse.Token = append(assetResponse.Token, &pb.TokenBalance{
					ContractName:         "Ether",
					ContractTickerSymbol: "ETH",
					ContractDecimals:     18,
					ContractAddress:      item.ContractAddress,
					Coin:                 280,
					Balance:              item.Balance,
					Quote:                quotePrice,
					QuotePrice:           strconv.FormatFloat(quotePrice, 'f', -1, 64),
					QuoteRate:            quoteRate,
					LogoUrl:              "https://assets.unmarshal.io/tokens/ethereum_0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png",
					QuoteRate_24H:        quotePriceChange24H,
					QuotePctChange_24H:   quotePricePercentageChange24H,
				})
			} else {
				var quotePrice = 0.0
				var quotePriceChange24H = ""
				var quotePricePercentageChange24H = 0.0
				var quoteRate = 0.0
				var asset pb.TokenBalance
				asset.LogoUrl = "" // To be decided later
				asset.Coin = 280
				tokenInfo, err := z.GetTokenPrice("usd", item.ContractAddress)
				if err == nil {
					balance, _ := strconv.ParseFloat(item.Balance, 64)
					quotePrice = balance / math.Pow(10, float64(18)) * tokenInfo.Price
					quoteRate = tokenInfo.Price
					//quotePriceChange24H = strconv.FormatFloat(tokenInfo.Token.MarketData.PriceChange24H, 'f', 2, 64)
					//quotePricePercentageChange24H = tokenInfo.Token.MarketData.PriceChangePercentage24HInCurrency.Usd
				}
				asset.ContractAddress = item.ContractAddress
				asset.Quote = quotePrice
				asset.QuoteRate = quoteRate
				asset.QuotePctChange_24H = quotePricePercentageChange24H
				asset.QuoteRate_24H = quotePriceChange24H
				contractDecimal, err := strconv.Atoi(item.Decimals)
				if err != nil {
					contractDecimal = 0
				}
				asset.ContractDecimals = int32(contractDecimal)
				asset.Balance = item.Balance
				asset.QuotePrice = strconv.FormatFloat(quotePrice, 'f', -1, 64)
				asset.ContractTickerSymbol = item.Symbol
				asset.ContractName = item.Name
				assetResponse.Token = append(assetResponse.Token, &asset)
			}

		}
	}
	return &assetResponse, nil
}

func (z *ServiceZksync) ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	url := fmt.Sprintf("https://zksync2-testnet-explorer.zksync.dev/transactions?limit=%v&direction=older&accountAddress=%v", request.Page, request.Address)
	body, err := z.httpRequest.GetRequest(url)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var zsynclistTransactionResponse ListTransactionResponse
	err = json.Unmarshal(body, &zsynclistTransactionResponse)
	if err != nil {
		z.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// Data Logic
	page, _ := strconv.Atoi(request.Page)
	//pageSize, _ := strconv.Atoi(request.PageSize)
	var transactionResponse pb.ListTransactionResponse
	transactionResponse.Transactions = make([]*pb.TransactionData, 0)
	transactionResponse.TotalTxs = int64(zsynclistTransactionResponse.Total)
	transactionResponse.Page = int64(page)
	transactionResponse.ItemsOnPage = int64(page)
	transactionResponse.TotalPages = int64(page)

	for _, item := range zsynclistTransactionResponse.List {
		var sent []*pb.TransactionInfo
		var received []*pb.TransactionInfo
		var others []*pb.TransactionInfo
		var transactionData pb.TransactionData

		var transactionType = ""
		balanceChanges := item.BalanceChanges
		for _, initem := range balanceChanges {
			transactionType = initem.Type
			if initem.Type == "transfer" {
				if strings.ToLower(initem.From) == strings.ToLower(request.Address) {
					transactionType = "send"
					qrate := z.helper.ConvertStringToFloat64(initem.TokenInfo.UsdPrice)
					sent = append(sent, &pb.TransactionInfo{
						Name:      initem.TokenInfo.Name,
						Decimals:  int64(initem.TokenInfo.Decimals),
						Value:     z.helper.ConvertHexFloatString(initem.Amount),
						QuoteRate: float32(qrate),
						LogoUrl:   "",
						From:      initem.From,
						To:        initem.To,
						Symbol:    initem.TokenInfo.Symbol,
						TokenId:   initem.TokenInfo.Address,
					})
				}
				if strings.ToLower(initem.To) == strings.ToLower(request.Address) {
					qrate := z.helper.ConvertStringToFloat64(initem.TokenInfo.UsdPrice)
					transactionType = "receive"
					received = append(received, &pb.TransactionInfo{
						Name:      initem.TokenInfo.Name,
						Decimals:  int64(initem.TokenInfo.Decimals),
						Value:     z.helper.ConvertHexFloatString(initem.Amount),
						QuoteRate: float32(qrate),
						LogoUrl:   "",
						From:      initem.From,
						To:        initem.To,
						Symbol:    initem.TokenInfo.Symbol,
						TokenId:   initem.TokenInfo.Address,
					})
				}
			}

		}
		description := ""
		if transactionType == "deposit" || transactionType == "send" || transactionType == "transfer" ||
			transactionType == "addliquidity" || transactionType == "approve" {
			if len(sent) != 0 {
				for _, sendTxObject := range sent {
					sendTxValue, _ := strconv.ParseFloat(sendTxObject.Value, 64)
					sendTxValue = sendTxValue / math.Pow(10, float64(sendTxObject.Decimals))
					value := strconv.FormatFloat(sendTxValue, 'f', -1, 64)
					description = description + value + " " + sendTxObject.Symbol + " "
				}
				description = transactionType + " " + description
			}
		} else if transactionType == "receive" || transactionType == "withdraw" {
			if len(received) != 0 {
				for _, receiveTxObject := range received {
					receiveTxValue, _ := strconv.ParseFloat(receiveTxObject.Value, 64)
					receiveTxValue = receiveTxValue / math.Pow(10, float64(receiveTxObject.Decimals))
					value := strconv.FormatFloat(receiveTxValue, 'f', -1, 64)
					description = description + value + " " + receiveTxObject.Symbol + " "
				}
				description = transactionType + " " + description
			}
		} else if transactionType == "swap" {
			var descriptionSent string
			var descriptionReceive string
			if len(sent) != 0 {
				for _, sendTxObject := range sent {
					sendTxValue, _ := strconv.ParseFloat(sendTxObject.Value, 64)
					sendTxValue = sendTxValue / math.Pow(10, float64(sendTxObject.Decimals))
					value := strconv.FormatFloat(sendTxValue, 'f', -1, 64)
					descriptionSent = descriptionSent + value + " " + sendTxObject.Symbol + " "
				}
			}
			if len(received) != 0 {
				for _, receiveTxObject := range received {
					receiveTxValue, _ := strconv.ParseFloat(receiveTxObject.Value, 64)
					receiveTxValue = receiveTxValue / math.Pow(10, float64(receiveTxObject.Decimals))
					value := strconv.FormatFloat(receiveTxValue, 'f', -1, 64)
					descriptionReceive = descriptionReceive + value + " " + receiveTxObject.Symbol + " "
				}
			}
			description = transactionType + " " + descriptionSent + " for " + descriptionReceive

		}
		transactionData.Sent = append(transactionData.Sent, sent...)
		transactionData.Received = append(transactionData.Received, received...)
		transactionData.Others = append(transactionData.Others, others...)

		transactionData.Id = item.TransactionHash
		transactionData.Date = item.ReceivedAt.Unix()
		transactionData.To = item.InitiatorAddress
		transactionData.Value = item.Data.Value
		transactionData.Description = description
		transactionData.Status = item.Status
		transactionData.From = item.Data.ContractAddress
		transactionData.Nonce = int64(item.Nonce)
		transactionData.Block = int64(item.BlockNumber)
		transactionData.Fee = z.helper.ConvertHexFloatString(item.Fee)
		transactionData.Type = transactionType

		transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)

	}
	return &transactionResponse, nil
}
