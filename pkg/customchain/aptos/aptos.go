package aptos

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
	"strconv"
	"strings"
)

type IAptos interface {
	GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error)
	GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error)
	//ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error)
	GetTxStatus(request *pb.TxStatusRequest) (*pb.TxStatusResponse, error)
	SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error)
}

type ServiceAptos struct {
	rpc         map[string]*ethrpc.EthRPC
	env         *config.Config
	logger      *zap.SugaredLogger
	util        *utils.UtilConf
	ethgoRpc    map[string]*ethgoJsonRPC.Client
	httpRequest utils.IHttpRequest
	rpcHandler  *jsonrpc.RPCHandler
	coingecko   coingecko.ICoinGecko
	helper      utils.Helpers
}

func NewServiceAptos(config *config.Config, logger *zap.SugaredLogger, coingecko coingecko.ICoinGecko) *ServiceAptos {
	rpc := make(map[string]*ethrpc.EthRPC)
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported Non-EVM chains:", len(config.EVM.Cfg.Wallets))
	if len(config.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No Non-EVM wallet configurations found")
	}
	//Initialize EVM RPC configurations
	for i, w := range config.EVM.Cfg.Wallets {
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

	return &ServiceAptos{rpc, config, logger, utilsManager, ethgoRpc,
		httpRequest, rpcHandler, coingecko, utils.Helpers{}}
}

func (s *ServiceAptos) GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	url := "https://mainnet-aptos-api.nodereal.io/api/account/" + request.Address + "/coins"
	body, err := s.httpRequest.GetRequest(url)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Http request error")
	}
	var aptosResponse BalanceResp
	err = json.Unmarshal(body, &aptosResponse)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}

	var quotePrice = 0.0
	var quotePriceChange24H = ""
	var quotePricePercentageChange24H = 0.0
	var quoteRate = 0.0

	var balanceResponse pb.BalanceResponse

	for _, j := range aptosResponse.Data {
		if j.Amount > 0 {
			if j.CoinType == "0x1::aptos_coin::AptosCoin" {
				amount := j.Amount
				var nativeQuotePrice float64
				var nativeQuoteRate float64
				nativeBalance := strconv.FormatFloat(amount, 'f', -1, 64)
				tokenInfo, err := s.coingecko.GetTokenExchange("usd", request.Chain)
				if err == nil {
					nativeQuotePrice = amount * tokenInfo.Price
					nativeQuoteRate = tokenInfo.Price
				}
				balanceResponse.Token = append(balanceResponse.Token, &pb.TokenBalance{
					ContractName:         j.Name,
					ContractTickerSymbol: j.Symbol,
					ContractDecimals:     int32(j.Decimals),
					ContractAddress:      "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
					Coin:                 0,
					Balance:              nativeBalance,
					Quote:                nativeQuotePrice,
					QuotePrice:           strconv.FormatFloat(nativeQuotePrice, 'f', -1, 64),
					QuoteRate:            nativeQuoteRate,
					LogoUrl:              s.util.GetNonEVMWalletInfo(request.Chain).Aptos.AptosLogoUrl,
					QuoteRate_24H:        quotePriceChange24H,
					QuotePctChange_24H:   quotePricePercentageChange24H,
				})
			} else {
				address := s.TrimExtraCharsInAddress(j.CoinType)
				balance := strconv.FormatFloat(j.Amount, 'f', -1, 64)
				balanceResponse.Token = append(balanceResponse.Token, &pb.TokenBalance{
					ContractName:         j.Name,
					ContractTickerSymbol: j.Symbol,
					ContractDecimals:     int32(j.Decimals),
					ContractAddress:      address,
					Coin:                 0,
					Balance:              balance,
					Quote:                quotePrice,
					QuotePrice:           strconv.FormatFloat(quotePrice, 'f', -1, 64),
					QuoteRate:            quoteRate,
					LogoUrl:              s.util.GetNonEVMWalletInfo(request.Chain).Aptos.AptosLogoUrl,
					QuoteRate_24H:        quotePriceChange24H,
					QuotePctChange_24H:   quotePricePercentageChange24H,
				})
			}
		}
	}
	return &balanceResponse, nil
}

func (s *ServiceAptos) GetTokenPrice(currency string, tokenAddress string) (*pb.TokenPriceResponse, error) {
	return s.coingecko.GetTokenExchange(currency, "aptos")
}

//commented out for now as this needs to be refactored and more enhanced for future use

//func (s *ServiceAptos) ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
//	url := s.util.GetNonEVMWalletInfo(request.Chain).Aptos.Rpc + "accounts/" + request.Address + "/transactions"
//	body, err := s.httpRequest.GetRequest(url)
//	if err != nil {
//		s.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error(), "Http request error")
//	}
//
//	var aptosTxList ListTransactions
//	err = json.Unmarshal(body, &aptosTxList)
//	if err != nil {
//		s.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
//	}
//	var nativeTokenPrice float64
//	tokenInfo, err := s.coingecko.GetTokenExchange("usd", "aptos")
//	if err != nil {
//		nativeTokenPrice = 0
//	} else {
//		nativeTokenPrice = tokenInfo.Price
//	}
//
//	var txListResp pb.ListTransactionResponse
//	page, _ := strconv.Atoi(request.Page)
//	txListResp.Transactions = make([]*pb.TransactionData, 0)
//	txListResp.TotalTxs = int64(page)
//	txListResp.Page = int64(page)
//	txListResp.ItemsOnPage = int64(page)
//	txListResp.TotalPages = int64(page)
//
//	for _, item := range aptosTxList {
//		var sent []*pb.TransactionInfo
//		var received []*pb.TransactionInfo
//		var transactionData pb.TransactionData
//		var toAddress string
//		if len(item.Events) == 2 {
//			toAddress = item.Events[1].GUID.AccountAddress
//		}
//		fromAddress := item.Events[0].GUID.AccountAddress
//		if len(item.Events) == 2 {
//			if item.Type == "user_transaction" {
//				txType := s.TrimExtraCharsInType(item.Events[0].Type)
//				txType1 := s.TrimExtraCharsInType(item.Events[1].Type)
//				if txType == "WithdrawEvent" {
//					address := s.TrimExtraCharsInAddress(item.Changes[0].Data.Type)
//					var tokenPrice float64
//					var tokenName string
//					var tokenDecimal int64
//					var tokenSymbol string
//					var tokenLogo string
//					var tokenID string
//					if address == "0x1" {
//						tokenPrice = nativeTokenPrice
//						address = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
//						tokenName = "Aptos Coin"
//						tokenDecimal = 8
//						tokenSymbol = "APT"
//						tokenLogo = s.util.GetNonEVMWalletInfo(request.Chain).Aptos.AptosLogoUrl
//						tokenID = address
//					} else {
//						tokenPrice = 0
//						tokenName = "Ocean Park Coin"
//						tokenDecimal = 6
//						tokenSymbol = "OPC"
//						tokenLogo = ""
//						tokenID = address
//					}
//					qrate := tokenPrice
//					sent = append(sent, &pb.TransactionInfo{
//						Name:      tokenName,
//						Decimals:  tokenDecimal,
//						Value:     item.Events[0].Data.Amount,
//						QuoteRate: float32(qrate),
//						LogoUrl:   tokenLogo,
//						From:      fromAddress,
//						To:        toAddress,
//						Symbol:    tokenSymbol,
//						TokenId:   tokenID,
//					})
//				}
//				if txType1 == "DepositEvent" {
//					address := s.TrimExtraCharsInAddress(item.Changes[1].Data.Type)
//					var tokenPrice float64
//					var tokenName string
//					var tokenDecimal int64
//					var tokenSymbol string
//					var tokenLogo string
//					var tokenID string
//					if address == "0x1" {
//						address = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
//						tokenPrice = nativeTokenPrice
//						tokenName = "Aptos Coin"
//						tokenDecimal = 8
//						tokenSymbol = "APT"
//						tokenLogo = s.util.GetNonEVMWalletInfo(request.Chain).Aptos.AptosLogoUrl
//						tokenID = address
//					} else {
//						tokenPrice = 0
//						tokenName = "Ocean Park Coin"
//						tokenDecimal = 6
//						tokenSymbol = "OPC"
//						tokenLogo = ""
//						tokenID = address
//					}
//					qrate := tokenPrice
//					received = append(received, &pb.TransactionInfo{
//						Name:      tokenName,
//						Decimals:  tokenDecimal,
//						Value:     item.Events[0].Data.Amount,
//						QuoteRate: float32(qrate),
//						LogoUrl:   tokenLogo,
//						From:      fromAddress,
//						To:        toAddress,
//						Symbol:    tokenSymbol,
//						TokenId:   tokenID,
//					})
//				}
//			}
//		} else if len(item.Events) == 1 {
//			if item.Type == "user_transaction" {
//				txType := s.TrimExtraCharsInType(item.Events[0].Type)
//				//txType1 := s.TrimExtraCharsInType(item.Events[1].Type)
//				if txType == "WithdrawEvent" {
//					address := s.TrimExtraCharsInAddress(item.Changes[0].Data.Type)
//					var tokenPrice float64
//					var tokenName string
//					var tokenDecimal int64
//					var tokenSymbol string
//					var tokenLogo string
//					var tokenID string
//					if address == "0x1" {
//						tokenPrice = nativeTokenPrice
//						address = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
//						tokenName = "Aptos Coin"
//						tokenDecimal = 8
//						tokenSymbol = "APT"
//						tokenLogo = s.util.GetNonEVMWalletInfo(request.Chain).Aptos.AptosLogoUrl
//						tokenID = address
//					} else {
//						tokenPrice = 0
//						tokenName = "Ocean Park Coin"
//						tokenDecimal = 6
//						tokenSymbol = "OPC"
//						tokenLogo = ""
//						tokenID = address
//					}
//					qrate := tokenPrice
//					sent = append(sent, &pb.TransactionInfo{
//						Name:      tokenName,
//						Decimals:  tokenDecimal,
//						Value:     item.Events[0].Data.Amount,
//						QuoteRate: float32(qrate),
//						LogoUrl:   tokenLogo,
//						From:      fromAddress,
//						To:        toAddress,
//						Symbol:    tokenSymbol,
//						TokenId:   tokenID,
//					})
//				}
//			}
//		}
//		transactionData.Sent = append(transactionData.Sent, sent...)
//		transactionData.Received = append(transactionData.Received, received...)
//		transactionData.Id = item.Hash
//		timestamp, err := strconv.Atoi(item.Timestamp)
//		if err != nil {
//			timestamp = 0
//		}
//		transactionData.Date = int64(timestamp)
//		transactionData.To = toAddress
//		transactionData.Value = item.Events[0].Data.Amount
//		transactionData.Description = item.Type
//		transactionData.Status = fmt.Sprintf("%v", item.Success)
//		transactionData.From = fromAddress
//		transactionData.Nonce = 0
//		transactionData.NativeTokenDecimals = 8
//		block, _, _ := s.GetBlocks(request.Chain, item.Version)
//		transactionData.Block = block
//		transactionData.Fee = ""
//		transactionData.Type = item.Type
//
//		txListResp.Transactions = append(txListResp.Transactions, &transactionData)
//	}
//	return &txListResp, nil
//}

func (s *ServiceAptos) GetTxStatus(request *pb.TxStatusRequest) (*pb.TxStatusResponse, error) {
	pendingTxReceipt := &pb.TxStatusResponse{
		TransactionHash:   "",
		TransactionIndex:  0,
		BlockHash:         "",
		BlockNumber:       0,
		CumulativeGasUsed: 0,
		GasUsed:           0,
		ContractAddress:   "",
		LogsBloom:         "",
		Root:              "",
		Status:            "PENDING",
	}
	url := s.util.GetNonEVMWalletInfo(request.Chain).Aptos.Rpc + "transactions/by_hash/" + request.TxHash
	body, err := s.httpRequest.GetRequest(url)
	if err != nil || body == nil {
		s.logger.Info("Transaction pending/unknown")
		s.logger.Error("Error: ", err.Error())
		return pendingTxReceipt, status.Errorf(codes.OK, err.Error(), "Transaction pending/unknown")
	}
	var TxStatus string
	var aptosTxStatus TxByHash
	err = json.Unmarshal(body, &aptosTxStatus)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}
	response := pb.TxStatusResponse{}
	blockNumber, blockHash, err := s.GetBlocks(request.Chain, aptosTxStatus.Version)
	gasUsed, _ := strconv.Atoi(aptosTxStatus.GasUsed)
	cumulativeGasUsed, _ := strconv.Atoi(aptosTxStatus.MaxGasAmount)
	if aptosTxStatus.Success == true {
		TxStatus = "SUCCESS"
	} else if aptosTxStatus.Success == false {
		TxStatus = "FAILURE"
	}
	response = pb.TxStatusResponse{
		TransactionHash:   aptosTxStatus.Hash,
		TransactionIndex:  0,
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		CumulativeGasUsed: int64(cumulativeGasUsed),
		GasUsed:           int64(gasUsed),
		ContractAddress:   "",
		LogsBloom:         "",
		Root:              "",
		Status:            TxStatus,
	}
	return &response, nil
}

func (s *ServiceAptos) SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	url := s.util.GetNonEVMWalletInfo(request.Chain).Aptos.Rpc + "/transactions"
	body, err := s.httpRequest.PostRequestWithHeaders(url, request.Msg, "Content-Type", "application/x.aptos.signed_transaction+bcs")
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Http request error")
	}
	var sendTxResponse SendTxResponse
	err = json.Unmarshal(body, &sendTxResponse)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}
	sendTransactionResponse := pb.SendTransactionResponse{
		TransactionId: sendTxResponse.Hash,
	}
	return &sendTransactionResponse, nil
}

func (s *ServiceAptos) TrimExtraCharsInAddress(address string) string {
	//address1 := strings.SplitAfter(address, "<")
	trimmedAddress := strings.Split(address, "::")
	return fmt.Sprintf("%v", trimmedAddress[0])
}

func (s *ServiceAptos) TrimExtraCharsInType(data string) string {
	newData := strings.Split(data, "::")
	return fmt.Sprintf("%v", newData[2])
}

func (s *ServiceAptos) GetBlocks(chain string, version string) (int64, string, error) {
	url := s.util.GetNonEVMWalletInfo(chain).Aptos.Rpc + "blocks/by_version/" + version
	body, err := s.httpRequest.GetRequest(url)
	if err != nil {
		s.logger.Error(err)
		return 0, "", status.Errorf(codes.Internal, err.Error(), "Http request error")
	}

	var blocks Blocks
	err = json.Unmarshal(body, &blocks)
	if err != nil {
		s.logger.Error(err)
		return 0, "", status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}
	block := blocks.BlockHeight
	intBlock, _ := strconv.Atoi(block)

	blockHash := blocks.BlockHash
	return int64(intBlock), blockHash, nil
}
