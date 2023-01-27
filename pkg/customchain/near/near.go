package near

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/jsonrpc"
	"bridge-allowance/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/onrik/ethrpc"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"strconv"
	"strings"
)

type INear interface {
	GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error)
	//GetTokenPriceV2(request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error)
	GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error)
	GetNonce(request *pb.NonceRequest) (*pb.NonceResponse, error)
	//GasLimit(request *pb.GasLimitRequest) (*pb.GasLimitResponse, error)
	SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error)
	ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error)
	//GetUserData(request *pb.UserDataRequest) (*pb.UserDataResponse, error)
	//GetProcessingFee(request *pb.ProcessingFeeRequest) (*pb.ProcessingFeeResponse, error)
	//GetTxStatus(request *pb.TxStatusRequest) (*pb.TxStatusResponse, error)
	//GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error)
	//TokenApprove(request *pb.ApprovalRequest) (*pb.ApprovalResponse, error)
}

type ServiceNear struct {
	rpc    map[string]*ethrpc.EthRPC
	env    *config.Config
	logger *zap.SugaredLogger
	util   *utils.UtilConf

	ethgoRpc      map[string]*ethgoJsonRPC.Client
	httpRequest   utils.IHttpRequest
	rpcHandler    *jsonrpc.RPCHandler
	RestAPI       string
	coingecko     coingecko.ICoinGecko
	helper        utils.Helpers
	TokenListInfo map[string]TokenPriceItem
}

func NewServiceNear(config *config.Config, logger *zap.SugaredLogger, gecko coingecko.ICoinGecko) *ServiceNear {
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
	var tokeninfo map[string]TokenPriceItem
	url := fmt.Sprintf("https://indexer.ref.finance/list-token-price")

	body, err := httpRequest.GetRequest(url)
	if err != nil {
		//l.logger.Error(err)
		//return nil, status.Errorf(codes.Internal, string(body), "lifi third party error")
	}
	err = json.Unmarshal(body, &tokeninfo)
	if err != nil {
		//n.logger.Error(err)
		//return nil, status.Errorf(codes.Internal, err)
	}
	return &ServiceNear{rpc, config, logger, utilsManager, ethgoRpc,
		httpRequest, rpcHandler, restApi, gecko, utils.Helpers{}, tokeninfo}
}

func (n *ServiceNear) GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error) {

	requestObject := BalanceRequest{
		Jsonrpc: "2.0",
		Method:  "query",
		Params: BalanceParams{
			RequestType: "view_account",
			Finality:    "final",
			AccountId:   request.Address,
		},
		Id: 1}

	url := n.util.GetNonEVMWalletInfo(request.Chain).Near.Rpc
	buildData, err := json.Marshal(requestObject)
	if err != nil {
		n.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshaling error")
	}

	body, err := n.httpRequest.PostRequest(url, bytes.NewBuffer(buildData))
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "http error")
	}

	var balanceResponse BalanceResponse
	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
	}
	var assetResponse pb.BalanceResponse
	var quotePrice = 0.0
	var quotePriceChange24H = ""
	var quotePricePercentageChange24H = 0.0
	var quoteRate = 0.0
	tokenInfo, err := n.GetTokenPrice("usd", "")
	if err == nil {
		balance, _ := strconv.ParseFloat(balanceResponse.Result.Amount, 64)
		quotePrice = balance / math.Pow(10, float64(18)) * tokenInfo.Price
		quoteRate = tokenInfo.Price
		//quotePriceChange24H = strconv.FormatFloat(tokenInfo.Token.MarketData.PriceChange24H, 'f', 2, 64)
		//quotePricePercentageChange24H = tokenInfo.Token.MarketData.PriceChangePercentage24HInCurrency.Usd
	}
	assetResponse.Token = append(assetResponse.Token, &pb.TokenBalance{
		ContractName:         "Near",
		ContractTickerSymbol: "Near",
		ContractDecimals:     18,
		ContractAddress:      "",
		Coin:                 0,
		Balance:              balanceResponse.Result.Amount,
		Quote:                quotePrice,
		QuotePrice:           strconv.FormatFloat(quotePrice, 'f', -1, 64),
		QuoteRate:            quoteRate,
		LogoUrl:              "",
		QuoteRate_24H:        quotePriceChange24H,
		QuotePctChange_24H:   quotePricePercentageChange24H,
	})
	return &assetResponse, nil
}

func (n *ServiceNear) GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error) {
	return n.coingecko.GetTokenExchange(currency, "near")
}

func (n *ServiceNear) GetNonce(request *pb.NonceRequest) (*pb.NonceResponse, error) {
	requestObj := NonceRequest{Jsonrpc: "2.0",
		Method: "query",
		Params: Params{
			RequestType: "view_access_key",
			Finality:    "final",
			AccountId:   request.Address,
			PublicKey:   "ed25519:4oWce12aRTr9sdgyYm3wcp995Q2ugNapkcNsp1iWSQ32",
		},
		Id: 1,
	}
	url := n.util.GetNonEVMWalletInfo(request.Chain).Near.Rpc
	buildData, err := json.Marshal(requestObj)
	if err != nil {
		n.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshaling error")
	}

	body, err := n.httpRequest.PostRequest(url, bytes.NewBuffer(buildData))
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "http error")
	}

	var nonceResponse NonceResponse
	err = json.Unmarshal(body, &nonceResponse)
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
	}
	tokenPrice, err := n.coingecko.GetTokenExchange("usd", "near")

	var fastestFee = 0.0
	var slowAvgFee = 0.0
	var fastFee = 0.0
	gasPriceInfo := pb.GasPriceInfo{
		Fast:        fastFee,
		SafeLow:     slowAvgFee,
		Fastest:     fastestFee,
		Average:     slowAvgFee,
		SafeLowWait: 5,
		AvgWait:     2,
		FastWait:    1,
		FastestWait: 0.5,
	}
	return &pb.NonceResponse{
		Nonce:      int64(nonceResponse.Result.Nonce),
		QuoteValue: tokenPrice.Price,
		GasPrice:   &gasPriceInfo,
		OpL1Fee:    0.0,
	}, nil
}

func (n *ServiceNear) SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	requestObj := SendTransactionRequest{Jsonrpc: "2.0",
		Method: "broadcast_tx_commit",
		Params: []string{request.Msg},
		Id:     1,
	}
	url := n.util.GetNonEVMWalletInfo(request.Chain).Near.Rpc
	buildData, err := json.Marshal(requestObj)
	if err != nil {
		n.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshaling error")
	}

	body, err := n.httpRequest.PostRequest(url, bytes.NewBuffer(buildData))
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "http error")
	}

	var sendTxResponse SendTrandsactionResponse
	err = json.Unmarshal(body, &sendTxResponse)
	if sendTxResponse.Result == "" {
		var sendtTxError SendTxError
		err = json.Unmarshal(body, &sendtTxError)
		if err != nil {
			n.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), sendtTxError.Error.Cause.Name)
	}
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
	}
	sendTransactionResponse := pb.SendTransactionResponse{
		TransactionId: sendTxResponse.Result,
	}
	return &sendTransactionResponse, nil
}

func (n *ServiceNear) ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	url := "https://helper.sender.org/account/" + request.Address + "/activity"
	body, err := n.httpRequest.GetRequest(url)
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var history []*HistoryResponseNew
	err = json.Unmarshal(body, &history)
	if err != nil {
		n.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var transactionResponse pb.ListTransactionResponse
	page, _ := strconv.Atoi(request.Page)
	transactionResponse.Transactions = make([]*pb.TransactionData, 0)
	transactionResponse.TotalTxs = int64(page)
	transactionResponse.Page = int64(page)
	transactionResponse.ItemsOnPage = int64(page)
	transactionResponse.TotalPages = int64(page)

	for _, item := range history {
		var sent []*pb.TransactionInfo
		var received []*pb.TransactionInfo
		var others []*pb.TransactionInfo
		var transactionData pb.TransactionData
		var transactionType = ""
		toAddress := ""
		value := ""
		if item.ActionKind == "TRANSFER" {
			if strings.ToLower(item.ReceiverId) == strings.ToLower(request.Address) {
				transactionType = "recieve"
				//tokenINfo:= n.TokenListInfo
				var tokenPrice float64
				tokenInfo, err := n.coingecko.GetTokenExchange("usd", "near")
				if err != nil {
					tokenPrice = 0
				} else {
					tokenPrice = tokenInfo.Price
				}
				qrate := tokenPrice
				toAddress = item.ReceiverId
				value = item.Args.Deposit
				received = append(received, &pb.TransactionInfo{
					Name:      "near",
					Decimals:  24,
					Value:     item.Args.Deposit,
					QuoteRate: float32(qrate),
					LogoUrl:   "",
					From:      item.SignerId,
					To:        item.ReceiverId,
					Symbol:    "near",
					TokenId:   "",
				})
			}
			if strings.ToLower(item.SignerId) == strings.ToLower(request.Address) {
				transactionType = "sent"
				//tokenINfo:= n.TokenListInfo
				var tokenPrice float64
				tokenInfo, err := n.coingecko.GetTokenExchange("usd", "near")
				if err != nil {
					tokenPrice = 0
				} else {
					tokenPrice = tokenInfo.Price
				}
				qrate := tokenPrice
				toAddress = item.ReceiverId
				value = item.Args.Deposit
				sent = append(sent, &pb.TransactionInfo{
					Name:      item.ReceiverId,
					Decimals:  24,
					Value:     item.Args.Deposit,
					QuoteRate: float32(qrate),
					LogoUrl:   "",
					From:      item.SignerId,
					To:        item.ReceiverId,
					Symbol:    "near",
					TokenId:   "",
				})
			}
		}
		if item.ActionKind == "FUNCTION_CALL" {
			if item.Args.MethodName == "ft_transfer" {
				if strings.ToLower(item.SignerId) == strings.ToLower(request.Address) {
					transactionType = "sent"
					tokenINfo := n.TokenListInfo
					toAddress = item.Args.ArgsJson.ReceiverId
					value = item.Args.ArgsJson.Amount
					quoteRate := n.helper.ConvertStringToFloat64(tokenINfo[item.ReceiverId].Price)
					sent = append(sent, &pb.TransactionInfo{
						Name:      item.ReceiverId,
						Decimals:  int64(tokenINfo[item.ReceiverId].Decimal),
						Value:     item.Args.ArgsJson.Amount,
						QuoteRate: float32(quoteRate),
						LogoUrl:   "",
						From:      item.SignerId,
						To:        item.Args.ArgsJson.ReceiverId,
						Symbol:    tokenINfo[item.ReceiverId].Symbol,
						TokenId:   item.ReceiverId,
					})
				}
				if strings.ToLower(item.Args.ArgsJson.ReceiverId) == strings.ToLower(request.Address) {
					transactionType = "recieve"
					tokenINfo := n.TokenListInfo
					quoteRate := n.helper.ConvertStringToFloat64(tokenINfo[item.ReceiverId].Price)
					toAddress = item.Args.ArgsJson.ReceiverId
					value = item.Args.ArgsJson.Amount
					received = append(received, &pb.TransactionInfo{
						Name:      item.ReceiverId,
						Decimals:  int64(tokenINfo[item.ReceiverId].Decimal),
						Value:     item.Args.ArgsJson.Amount,
						QuoteRate: float32(quoteRate),
						LogoUrl:   "",
						From:      item.SignerId,
						To:        item.Args.ArgsJson.ReceiverId,
						Symbol:    tokenINfo[item.ReceiverId].Symbol,
						TokenId:   item.ReceiverId,
					})
				}
			}
			if item.Args.MethodName == "deposit_and_stake" {
				transactionType = "deposit_and_stake"
				tokenINfo := n.TokenListInfo
				quoteRate := n.helper.ConvertStringToFloat64(tokenINfo[item.ReceiverId].Price)
				toAddress = item.ReceiverId
				value = item.Args.Deposit
				others = append(others, &pb.TransactionInfo{
					Name:      item.ReceiverId,
					Decimals:  int64(tokenINfo[item.ReceiverId].Decimal),
					Value:     item.Args.Deposit,
					QuoteRate: float32(quoteRate),
					LogoUrl:   "",
					From:      item.SignerId,
					To:        item.ReceiverId,
					Symbol:    tokenINfo[item.ReceiverId].Symbol,
					TokenId:   item.ReceiverId,
				})
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
		timestamp, err := strconv.Atoi(item.BlockTimestamp)
		if err != nil {
			timestamp = 0
		}
		transactionData.Id = item.Hash
		transactionData.Date = int64(timestamp)
		transactionData.To = toAddress
		transactionData.Value = value
		transactionData.Description = description
		transactionData.Status = ""
		transactionData.From = item.SignerId

		var _nonce int64
		nonce, err := n.GetNonce(&pb.NonceRequest{
			Address: item.SignerId,
			Chain:   request.Chain,
		})
		if err != nil {
			_nonce = 0
		} else {
			_nonce = nonce.Nonce
		}
		transactionData.Nonce = _nonce
		transactionData.Block = int64(0)
		transactionData.Fee = ""
		transactionData.Type = transactionType

		transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)
	}
	transactionResponse.TotalTxs = int64(len(transactionResponse.Transactions))
	transactionResponse.ItemsOnPage = int64(len(transactionResponse.Transactions))
	return &transactionResponse, nil
}
