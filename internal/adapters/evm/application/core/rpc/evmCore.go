package rpc

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/evm/application/core"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/jsonrpc"
	"bridge-allowance/utils"
	// "bridge-allowance/utils/models"
	"encoding/json"
	"fmt"
	"github.com/umbracle/ethgo/builtin/erc20"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"

	// "bridge-allowance/pkg/unmarshal"
	"github.com/onrik/ethrpc"
	"github.com/umbracle/ethgo"
	_ "github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	_ "github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	_ "github.com/umbracle/ethgo/contract"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EvmCore interface {
	GetTokenPrice(request *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error)
	GetTokenPriceV2(request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error)
	GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error)
	
}

var TxStatus = func() map[string]string {
	return map[string]string{
		"0x1": "SUCCESS",
		"0x0": "FAILED",
	}
}

const (
	ABI_TRANSFER_FUNCTION string = "function transfer(address, uint256) view returns (bool)"
	LatestBlock           string = "latest"
)

type evmCore struct {
	rpc         map[string]*ethrpc.EthRPC
	env         *config.Config
	logger      *zap.SugaredLogger
	services    common.Services
	util        *utils.UtilConf
	ethgoRpc    map[string]*ethgoJsonRPC.Client
	httpRequest utils.IHttpRequest
	rpcHandler  *jsonrpc.RPCHandler
}

// NewEVMCore Manager to initialize EVM specific rpc node endpoints
func NewEVMCore(config *config.Config, logger *zap.SugaredLogger, services common.Services) *evmCore {
	rpc := make(map[string]*ethrpc.EthRPC)
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	//Do not continue if no EVM configurations are provided
	logger.Info("Supported EVM chains:", len(config.EVM.Cfg.Wallets))
	if len(config.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
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
	return &evmCore{rpc, config, logger, services, utilsManager, ethgoRpc,
		httpRequest, rpcHandler}
}

// GetTokenPrice Core logic to get token price
func (evm *evmCore) GetTokenPrice(request *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error) {
	source := evm.util.GetWalletSource(request.Chain)
	switch source.TokenPriceSource {
	// case "coingecko":
	// 	return evm.services.CoinGecko.GetTokenExchange(request.Currency, request.Chain)
	// case "custom":
	// 	if request.Chain == "zksync" {
	// 		return evm.services.Zksync.GetTokenPrice(request.Currency, "")
	// 	} else {
	// 		return nil, status.Error(codes.Unimplemented,
	// 			fmt.Sprintf("Unsupported operation: Source %v is unsupported", source.TokenPriceSource))
	// 	}
	default:
		return nil, status.Errorf(codes.Unavailable,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", source.TokenPriceSource), "Unsupported source")
	}
}

// GetTokenPriceV2 Returns the response from GetTokenPrice's price in string format
func (evm *evmCore) GetTokenPriceV2(request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error) {
	source := evm.util.GetWalletSource(request.Chain)
	switch source.TokenPriceSource {
	// case "coingecko":
	// 	return evm.services.CoinGecko.GetTokenExchangeV2(request.Currency, request.Chain)
	default:
		return nil, status.Errorf(codes.Unavailable,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", source.TokenPriceSource), "Unsupported source")
	}
}

// GetAssets Core logic to get asset balances across token holdings
func (evm *evmCore) GetAssets(request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	source := evm.util.GetWalletSource(request.Chain)
	evm.logger.Debugf("Sources := %s", source)
	switch source.HistorySource {
	// case "unmarshal":
	// 	return transformUnmarshalBalances(evm, request)
	case "custom":
		return getCustomBalances(evm, request)
	default:
		return nil, status.Errorf(codes.Unavailable,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", source.HistorySource), "Unsupported source")
	}
}

// getCustomBalances construct token balances using custom logic
func getCustomBalances(evm *evmCore, request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	var assetResponse pb.BalanceResponse
	assetResponse.Token = make([]*pb.TokenBalance, 0)
	switch request.Chain {
	// case "tomochain":
	// 	return getTomoAssets(evm, request)
	// case "zksync":
	// 	return evm.services.Zksync.GetAssets(request)
	default:
		return nil, status.Errorf(codes.Unavailable, "Unsupported balances source", "Unsupported source")
	}
}

// // transformUnmarshalBalances transform unmarshal balances response
// func transformUnmarshalBalances(evm *evmCore, request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
// 	var assetResponse pb.BalanceResponse
// 	assetResponse.Token = make([]*pb.TokenBalance, 0)
// 	// Format response data according to *pb.AssetResponse contract
// 	balances, err := evm.services.Unmarshall.GetAssets(request.GetAddress(), request.Chain)
// 	if err != nil {
// 		evm.logger.Error("Error fetching balances. Err: ", err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	for _, item := range balances {
// 		if strings.ToLower(request.Chain) == "xinfin" && item.ContractAddress == "xdceeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
// 			item.ContractAddress = "xdceeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
// 		} else if strings.ToLower(request.Chain) == "optimism" && strings.ToLower(item.ContractAddress) == "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0000" {
// 			item.ContractAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
// 		}
// 		asset := pb.TokenBalance{
// 			ContractName:         Clean(item.ContractName).(string),
// 			ContractTickerSymbol: Clean(item.ContractTickerSymbol).(string),
// 			ContractDecimals:     item.ContractDecimals,
// 			ContractAddress:      Clean(item.ContractAddress).(string),
// 			Coin:                 item.Coin,
// 			Balance:              Clean(item.Balance).(string),
// 			Quote:                Clean(item.Quote).(float64),
// 			QuotePrice:           strconv.FormatFloat(Clean(item.Quote).(float64), 'f', -1, 64),
// 			QuoteRate:            Clean(item.QuoteRate).(float64),
// 			LogoUrl:              Clean(item.LogoURL).(string),
// 			QuoteRate_24H:        Clean(fmt.Sprintf("%v", item.QuoteRate24H)).(string),
// 			QuotePctChange_24H:   item.QuotePctChange24H,
// 		}
// 		assetResponse.Token = append(assetResponse.Token, &asset)
// 	}
// 	evm.logger.Debug("Unmarshall balances response: ", &balances)
// 	return &assetResponse, nil
// }

// getTomoAssets retrieve tokens list from tomo wallet API
// func getTomoAssets(evm *evmCore, request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
// 	var assetResponse pb.BalanceResponse
// 	assetResponse.Token = make([]*pb.TokenBalance, 0)
// 	//Fetch native token balance
// 	info := evm.util.GetWalletInfo(request.Chain)
// 	rpc := info.RPC
// 	client, err := ethgoJsonRPC.NewClient(rpc)
// 	if err != nil {
// 		panic(err)
// 	}
// 	balance, err := client.Eth().GetBalance(ethgo.HexToAddress(request.Address), ethgo.Latest)
// 	nativeTokenBalance := fmt.Sprintf("%v", balance)
// 	newNativeTokenBalance, err := strconv.ParseFloat(nativeTokenBalance, 64)
// 	tokenInfoRequest := &pb.TokenInfoRequest{
// 		Token: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
// 		Range: "",
// 		Chain: "tomochain",
// 	}
// 	//Default values
// 	var quotePrice = 0.0
// 	var quotePriceChange24H = ""
// 	var quotePricePercentageChange24H = 0.0
// 	var quoteRate = 0.0
// 	tokenInfo, err := evm.services.CoinGecko.GetTokenInfo(tokenInfoRequest)
// 	if err == nil {
// 		quotePrice = newNativeTokenBalance / math.Pow(10, float64(18)) * tokenInfo.Token.MarketData.CurrentPrice.Usd
// 		quoteRate = tokenInfo.Token.MarketData.CurrentPrice.Usd
// 		quotePriceChange24H = strconv.FormatFloat(tokenInfo.Token.MarketData.PriceChange24H, 'f', 2, 64)
// 		quotePricePercentageChange24H = tokenInfo.Token.MarketData.PriceChangePercentage24HInCurrency.Usd
// 	}
// 	//Set native token balances
// 	asset := pb.TokenBalance{
// 		ContractName:         "tomochain",
// 		ContractTickerSymbol: "TOMO",
// 		ContractDecimals:     18,
// 		ContractAddress:      "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
// 		Coin:                 88,
// 		Balance:              nativeTokenBalance,
// 		Quote:                quotePrice,
// 		QuotePrice:           strconv.FormatFloat(quotePrice, 'f', -1, 64),
// 		QuoteRate:            quoteRate,
// 		LogoUrl:              "https://raw.githubusercontent.com/tomochain/tokens/master/tokens/0x0000000000000000000000000000000000000001.png",
// 		QuoteRate_24H:        quotePriceChange24H,
// 		QuotePctChange_24H:   quotePricePercentageChange24H,
// 	}
// 	assetResponse.Token = append(assetResponse.Token, &asset)

// 	//Fetch balances for non-native tokens
// 	walletUrl := fmt.Sprintf("https://wallet.tomochain.com/api/tokens/?holder=%v", request.Address)
// 	body, err := evm.httpRequest.GetRequest(walletUrl)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tomoWalletBalances := make(TomoWalletBalances, 0)
// 	err = json.Unmarshal(body, &tomoWalletBalances)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 	}
// 	for _, item := range tomoWalletBalances {
// 		var quotePrice = 0.0
// 		var quotePriceChange24H = ""
// 		var quotePricePercentageChange24H = 0.0
// 		var quoteRate = 0.0
// 		asset := pb.TokenBalance{
// 			//Balances currently defaults to price in USD
// 			ContractName:         Clean(item.Name).(string),
// 			ContractTickerSymbol: Clean(item.Symbol).(string),
// 			ContractDecimals:     int32(item.Decimals),
// 			ContractAddress:      Clean(item.TokenAddress).(string),
// 			Coin:                 88,
// 			Balance:              Clean(item.Balance).(string),
// 			Quote:                quotePrice,
// 			QuotePrice:           strconv.FormatFloat(quotePrice, 'f', -1, 64),
// 			QuoteRate:            quoteRate,
// 			LogoUrl:              Clean(item.Icon).(string),
// 			QuoteRate_24H:        quotePriceChange24H,
// 			QuotePctChange_24H:   quotePricePercentageChange24H,
// 		}
// 		assetResponse.Token = append(assetResponse.Token, &asset)
// 	}
// 	return &assetResponse, nil
// }

// Quotepctchange24h calculate 24hour quote percentage change
func Quotepctchange24h(quoteRate, quoteRate24H float64) float64 {
	if quoteRate24H == 0 {
		return 0.0
	} else {
		return (quoteRate - quoteRate24H) / quoteRate24H * 100
	}
}

// Clean the given argument to remove nil references
func Clean(arg interface{}) interface{} {
	if arg == nil {
		return ""
	} else {
		return arg
	}
}

// // ListTransaction Lists the transactions activity for an address
// // TODO: Add support for 3P services
// func (evm *evmCore) ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
// 	source := evm.util.GetWalletSource(request.Chain)
// 	switch source.HistorySource {
// 	case "unmarshal":
// 		return transformUnmarshalTxHistory(evm, request)
// 	case "custom":
// 		return getCustomTxHistory(evm, request)
// 	default:
// 		return nil, status.Errorf(codes.Unavailable, "Unsupported Operation", "Unsupported Operation")
// 	}
// }

// getCustomTxHistory construct transactions history using custom logic
func getCustomTxHistory(evm *evmCore, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	var txResponse pb.ListTransactionResponse
	txResponse.Transactions = make([]*pb.TransactionData, 0)
	switch request.Chain {
	case "tomochain":
		return getTomoTxHistory(evm, request)
	// case "zksync":
	// 	return evm.services.Zksync.ListTransaction(request)
	default:
		return nil, status.Errorf(codes.Unavailable, "Unsupported balances source", "Unsupported balances source")
	}
}

func getTomoTxHistory(evm *evmCore, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	reqUrl := fmt.Sprintf("https://scan.tomochain.com/api/txs/listByAccount/%v?page=%v&limit=%v", request.Address, request.Page, request.PageSize)
	evm.logger.Info("TOMO Scan endpoint:= ", reqUrl)
	body, err := evm.httpRequest.GetRequest(reqUrl)
	if err != nil {
		evm.logger.Error("Error fetching history. Err: ", err)
		return nil, err
	}
	var transactionResponse pb.ListTransactionResponse
	transactionResponse.Transactions = make([]*pb.TransactionData, 0)
	evm.logger.Info("before marshaling")
	var tomoScan *TomoScan
	err = json.Unmarshal(body, &tomoScan)
	evm.logger.Info("before marshaling")
	if err != nil {
		evm.logger.Error("Error fetching history. Err: ", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	//Set pagination properties
	transactionResponse.TotalTxs = int64(tomoScan.Total)
	transactionResponse.Page = int64(tomoScan.CurrentPage)
	transactionResponse.ItemsOnPage = int64(tomoScan.PerPage)
	transactionResponse.TotalPages = int64(tomoScan.Pages)
	const InputFunction = "Function: "
	// Format response data as per *pb.AssetResponse contract
	for _, item := range tomoScan.Items {
		var sent []*pb.TransactionInfo
		var received []*pb.TransactionInfo
		var others []*pb.TransactionInfo
		var transactionType = ""
		// Set the transactionType
		txReceipt := tomoTxReceipt(evm, item.Hash)
		if txReceipt != nil {
			if txReceipt.InputData != "" {
				if strings.Contains(txReceipt.InputData, InputFunction) {
					transactionType = txReceipt.InputData[strings.Index(txReceipt.InputData, InputFunction)+10 : strings.Index(txReceipt.InputData, "(")]
				} else {
					transactionType = "unknown"
				}
				if len(txReceipt.Trc20Txs) != 0 {
					for _, tx := range txReceipt.Trc20Txs {
						var sentTx pb.TransactionInfo
						var receiveTx pb.TransactionInfo
						if strings.ToLower(tx.From) == strings.ToLower(request.Address) {
							if transactionType == "transfer" {
								transactionType = "send"
							}
							sentTx.Name = tx.Symbol
							sentTx.Symbol = tx.Symbol
							sentTx.TokenId = tx.Address
							sentTx.Decimals = int64(tx.Decimals)
							sentTx.Value = tx.Value
							sentTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", tx.Address)
							sentTx.From = tx.From
							sentTx.To = tx.To
							sent = append(sent, &sentTx)
						} else if strings.ToLower(tx.To) == strings.ToLower(request.Address) {
							if transactionType == "transfer" {
								transactionType = "receive"
							}
							receiveTx.Name = tx.Symbol
							receiveTx.Symbol = tx.Symbol
							receiveTx.TokenId = tx.Address
							receiveTx.Decimals = int64(tx.Decimals)
							receiveTx.Value = tx.Value
							receiveTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", tx.Address)
							receiveTx.From = tx.From
							receiveTx.To = tx.To
							received = append(received, &receiveTx)
						}
					}
				}
				if len(txReceipt.Trc21Txs) != 0 {
					for _, tx := range txReceipt.Trc21Txs {
						var sentTx pb.TransactionInfo
						var receiveTx pb.TransactionInfo
						if strings.ToLower(tx.From) == strings.ToLower(request.Address) {
							if transactionType == "transfer" {
								transactionType = "send"
							}
							sentTx.Name = tx.Symbol
							sentTx.Symbol = tx.Symbol
							sentTx.TokenId = tx.Address
							sentTx.Decimals = int64(tx.Decimals)
							sentTx.Value = tx.Value
							sentTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", tx.Address)
							sentTx.From = tx.From
							sentTx.To = tx.To
							sent = append(sent, &sentTx)
						} else if strings.ToLower(tx.To) == strings.ToLower(request.Address) {
							if transactionType == "transfer" {
								transactionType = "receive"
							}
							receiveTx.Name = tx.Symbol
							receiveTx.Symbol = tx.Symbol
							receiveTx.TokenId = tx.Address
							receiveTx.Decimals = int64(tx.Decimals)
							receiveTx.Value = tx.Value
							receiveTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", tx.Address)
							receiveTx.From = tx.From
							receiveTx.To = tx.To
							received = append(received, &receiveTx)
						}
					}
				}
			} else {
				txReceiptValue, _ := strconv.Atoi(txReceipt.Value)
				var sentTx pb.TransactionInfo
				var receiveTx pb.TransactionInfo
				if txReceiptValue > 0 {
					if strings.ToLower(item.From) == strings.ToLower(request.Address) {
						transactionType = "send"
						sentTx.Name = "TOMOCHAIN"
						sentTx.Symbol = "TOMO"
						sentTx.TokenId = "0x7576BA850a5485B381A68beE77F1383414F1837f"
						sentTx.Decimals = 18
						sentTx.Value = item.Value
						sentTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", "0x7576BA850a5485B381A68beE77F1383414F1837f")
						sentTx.From = item.From
						sentTx.To = item.To
						sent = append(sent, &sentTx)
					} else if strings.ToLower(item.To) == strings.ToLower(request.Address) {
						transactionType = "receive"
						receiveTx.Name = "TOMOCHAIN"
						receiveTx.Symbol = "TOMO"
						receiveTx.TokenId = "0x7576BA850a5485B381A68beE77F1383414F1837f"
						receiveTx.Decimals = 18
						receiveTx.Value = item.Value
						receiveTx.LogoUrl = fmt.Sprintf("https://raw.githubusercontent.com/tomochain/luaswap-token-list/master/src/tokens/icons/tomochain/%s.png", "0x7576BA850a5485B381A68beE77F1383414F1837f")
						receiveTx.From = item.From
						receiveTx.To = item.To
						received = append(received, &receiveTx)
					} else {
						transactionType = "unknown"
					}
				} else {
					transactionType = "unknown"
				}
			}
			gasPrice, err := strconv.ParseInt(item.GasPrice, 10, 64)
			if err != nil {
				gasPrice = 0
			}
			fee := fmt.Sprintf("%d", int64(item.GasUsed)*int64(gasPrice))
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
			transactionData := pb.TransactionData{
				Id:                  Clean(item.Hash).(string),
				From:                Clean(item.From).(string),
				To:                  Clean(item.To).(string),
				Fee:                 fee,
				Date:                item.Timestamp.Unix(),
				Status:              getStatus(item.Status),
				Type:                transactionType,
				Block:               int64(item.BlockNumber),
				Value:               item.Value,
				Nonce:               int64(item.Nonce),
				NativeTokenDecimals: 18,
				Description:         description,
				Sent:                sent,
				Received:            received,
				Others:              others,
			}
			transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)
		}
	}
	return &transactionResponse, nil
}
func getStatus(status bool) string {
	if status {
		return "completed"
	}
	return "error"
}
func tomoTxReceipt(evm *evmCore, txHash string) *TomoTxReceipt {
	url := fmt.Sprintf("https://scan.tomochain.com/api/txs/%v", txHash)
	evm.logger.Info("url: ", url)
	body, err := evm.httpRequest.GetRequest(url)
	if err != nil {
		evm.logger.Infof("Error reading transaction receipt for hash: %v", txHash)
	}
	var tomoTxReceipt *TomoTxReceipt
	err = json.Unmarshal(body, &tomoTxReceipt)
	if err != nil {
		evm.logger.Infof("Error marshaling tx receipt structure : %v", txHash)
	}
	return tomoTxReceipt
}

/*
func parseTRCEvents(events []LogEvents, request *pb.ListTransactionRequest) ([]*pb.TransactionInfo, []*pb.TransactionInfo, []*pb.TransactionInfo) {
	var sentTransactionInfoList []*pb.TransactionInfo
	var receivedTransactionInfoList []*pb.TransactionInfo
	for count, event := range events {
		if count == 0 {
			continue
		}
		var toParam Params
		var valueParam Params
		//Check if the event decoded name is deposit or withdrawal
		if strings.ToLower(event.Decoded.Name) == "deposit" || strings.ToLower(event.Decoded.Name) == "withdrawal" {
			//If the event decoded param name is deposit and one of the param is "dst"
			//This qualifies for a "to" in sent field
			if strings.ToLower(event.Decoded.Name) == "deposit" {
				if toParam = findKey(event.Decoded.Params, "dst"); toParam != (Params{}) {
					valueParam = findKey(event.Decoded.Params, "wad")
				}
			}
			//If the event decoded param name is withdrawl and one of the param is "dst"
			//This qualifies for a "to" in sent field
			if strings.ToLower(event.Decoded.Name) == "withdrawal" {
				if toParam = findKey(event.Decoded.Params, "src"); toParam != (Params{}) {
					valueParam = findKey(event.Decoded.Params, "wad")
				}
			}
			sentInfo := &pb.TransactionInfo{
				Name:     event.SenderName,
				Symbol:   event.SenderContractTickerSymbol,
				TokenId:  event.SenderAddress,
				Decimals: int64(event.SenderContractDecimals),
				Value:    valueParam.Value,
				//QuoteRate: 0,
				LogUrl: event.SenderLogoURL,
				From:   request.Address,
				To:     toParam.Value,
			}
			sentTransactionInfoList = append(sentTransactionInfoList, sentInfo)
		}

		// Process sent events
		transferEventFrom := findKey(event.Decoded.Params, "from")
		if strings.ToLower(transferEventFrom.Name) == "from" && strings.ToLower(transferEventFrom.Value) == strings.ToLower(request.Address) {
			valueParam := findKey(event.Decoded.Params, "value")
			if valueParam == (Params{}) {
				valueParam = findKey(event.Decoded.Params, "amount")
			}
			fromParam := findKey(event.Decoded.Params, "from")
			toParam := findKey(event.Decoded.Params, "to")
			sentInfo := &pb.TransactionInfo{
				Name:     event.SenderName,
				Symbol:   event.SenderContractTickerSymbol,
				TokenId:  event.SenderAddress,
				Decimals: int64(event.SenderContractDecimals),
				Value:    valueParam.Value,
				//QuoteRate: 0,
				LogUrl: event.SenderLogoURL,
				From:   fromParam.Value,
				To:     toParam.Value,
			}
			sentTransactionInfoList = append(sentTransactionInfoList, sentInfo)
		}
		// Process received events
		receiveEventFrom := findKey(event.Decoded.Params, "from")
		if receiveEventFrom.Name == "from" {
			if receiveEventRequestAddr := findKey(event.Decoded.Params, "to"); receiveEventRequestAddr.Value == strings.ToLower(request.Address) {
				valueParam := findKey(event.Decoded.Params, "value")
				if valueParam == (Params{}) {
					valueParam = findKey(event.Decoded.Params, "amount")
				}
				fromParam := findKey(event.Decoded.Params, "from")
				toParam := findKey(event.Decoded.Params, "to")
				receivedInfo := &pb.TransactionInfo{
					Name:     event.SenderName,
					Symbol:   event.SenderContractTickerSymbol,
					TokenId:  event.SenderAddress,
					Decimals: int64(event.SenderContractDecimals),
					Value:    valueParam.Value,
					//QuoteRate: 0,
					LogUrl: event.SenderLogoURL,
					From:   fromParam.Value,
					To:     toParam.Value,
				}
				receivedTransactionInfoList = append(receivedTransactionInfoList, receivedInfo)
			}
		}
	}
	return sentTransactionInfoList, receivedTransactionInfoList, nil
}
*/

// func transformUnmarshalTxHistory(evm *evmCore, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
// 	transactions, err := evm.services.Unmarshall.ListTransaction(request, request.Chain)
// 	if err != nil {
// 		evm.logger.Error("Error fetching history. Err: ", err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	var transactionResponse pb.ListTransactionResponse
// 	transactionResponse.Transactions = make([]*pb.TransactionData, 0)
// 	transactionResponse.Page = transactions.Page
// 	transactionResponse.ItemsOnPage = transactions.ItemsOnPage
// 	transactionResponse.TotalPages = transactions.TotalPages
// 	transactionResponse.TotalTxs = transactions.TotalTxs

// 	// Format response data as per *pb.AssetResponse contract
// 	for _, item := range transactions.Transactions {
// 		var transactionData pb.TransactionData
// 		transactionData.Id = item.ID
// 		transactionData.Date = item.Date
// 		transactionData.Type = item.Type
// 		transactionData.To = item.To
// 		transactionData.Value = item.Value
// 		transactionData.Nonce = item.Nonce
// 		transactionData.NativeTokenDecimals = item.NativeTokenDecimals
// 		transactionData.From = item.From
// 		transactionData.Status = item.Status
// 		transactionData.Fee = item.Fee
// 		transactionData.Block = item.Block
// 		transactionData.Description = item.Description
// 		for _, sent := range item.Sent {
// 			var sentTx pb.TransactionInfo
// 			sentTx.Name = sent.Name
// 			sentTx.To = sent.To
// 			sentTx.Value = sent.Value
// 			sentTx.LogoUrl = sent.LogoURL
// 			sentTx.QuoteRate = sent.QuoteRate
// 			sentTx.Decimals = sent.Decimals
// 			sentTx.Symbol = sent.Symbol
// 			sentTx.TokenId = sent.TokenID
// 			sentTx.From = sent.From
// 			transactionData.Sent = append(transactionData.Sent, &sentTx)
// 		}
// 		for _, received := range item.Received {
// 			var receivedTx pb.TransactionInfo
// 			receivedTx.Name = received.Name
// 			receivedTx.To = received.To
// 			receivedTx.Value = received.Value
// 			receivedTx.LogoUrl = received.LogoURL
// 			receivedTx.QuoteRate = received.QuoteRate
// 			receivedTx.Decimals = received.Decimals
// 			receivedTx.Symbol = received.Symbol
// 			receivedTx.TokenId = received.TokenID
// 			receivedTx.From = received.From
// 			transactionData.Received = append(transactionData.Received, &receivedTx)
// 		}
// 		for _, others := range item.Others {
// 			var OtherTx pb.TransactionInfo
// 			OtherTx.Name = others.Name
// 			OtherTx.To = others.To
// 			OtherTx.Value = others.Value
// 			OtherTx.LogoUrl = others.LogoURL
// 			OtherTx.QuoteRate = others.QuoteRate
// 			OtherTx.Decimals = others.Decimals
// 			OtherTx.Symbol = others.Symbol
// 			OtherTx.TokenId = others.TokenID
// 			OtherTx.From = others.From
// 			transactionData.Others = append(transactionData.Others, &OtherTx)
// 		}
// 		transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)
// 	}
// 	return &transactionResponse, nil
// }

var TransactionTypes = func() map[string]string {
	return map[string]string{
		"swap":         "swap",
		"transfer":     "send",
		"mint":         "addLiquidity",
		"addLiquidity": "addLiquidity",
		"withdrawal":   "withdraw",
		"withdrawn":    "withdraw",
		"withdraw":     "withdraw",
		"approval":     "approve",
	}
}

// func (evm *evmCore) GetNonce(request *pb.NonceRequest) (*pb.NonceResponse, error) {
// 	source := evm.util.GetWalletSource(request.Chain)
// 	switch source.NonceSource {
// 	case "debank":
// 		return getNonce(evm, request)
// 	case "custom":
// 		return getCustomNonce(evm, request)
// 	default:
// 		return nil, status.Errorf(codes.Unavailable, "Unsupported Operation", "Unsupported Operation")
// 	}
// }

// getEthTransactionCount returns the latest transactions count for the given address
func getEthTransactionCount(evm *evmCore, chain string, address string) (int, error) {
	nonce, err := evm.rpc[chain].EthGetTransactionCount(address, "latest")
	if err != nil {
		evm.logger.Error("Error fetching transaction count. Err: ", err)
		return 0, err
	}
	return nonce, nil
}

// getEthGasPrice returns the gas price of a chain
func getEthGasPrice(evm *evmCore, chain string) (big.Int, error) {
	gasPrice, err := evm.rpc[chain].EthGasPrice()
	if err != nil {
		evm.logger.Error("Error fetching gas price. Err: ", err)
		return big.Int{}, nil
	}
	return gasPrice, nil
}

// getCustomNonce fetch nonce for an address based on a custom logic
// func getCustomNonce(evm *evmCore, request *pb.NonceRequest) (*pb.NonceResponse, error) {
// 	evm.logger.Info("Getting into custom nonce")
// 	nonce, err := getEthTransactionCount(evm, request.Chain, request.Address)
// 	if err != nil {
// 		evm.logger.Error("Error fetching nonce. Err: ", err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	tokenPrice, err := evm.services.CoinGecko.GetTokenExchange("usd", request.Chain)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Error fetching quote price")
// 	}
// 	gasPrice, err := getEthGasPrice(evm, request.Chain)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Error fetching gas price")
// 	}
// 	evm.logger.Info("gasPrice: ", gasPrice)
// 	gasPriceInfo := pb.GasPriceInfo{}
// 	var fastestFee = 0.0
// 	var slowAvgFee = 0.0
// 	var fastFee = 0.0
// 	acc := big.Exact
// 	switch request.Chain {
// 	case "tomochain", "xinfin":
// 		//V1 compatible custom gas price calculation
// 		//Convert gas price to Gwei gasprice / 10^9
// 		gasPriceGwei := new(big.Float).Quo(new(big.Float).SetInt(&gasPrice), new(big.Float).SetInt64(utils.GWei))
// 		//fastest fee = base gas price * 20
// 		fastestFee, acc = new(big.Float).Mul(gasPriceGwei, new(big.Float).SetInt64(20)).Float64()
// 		if acc != big.Exact {
// 			//TODO:Handle loss of precision
// 		}
// 		//slow & average fee = base gas price * 4
// 		slowAvgFee, acc = new(big.Float).Mul(gasPriceGwei, new(big.Float).SetInt64(4)).Float64()
// 		if acc != big.Exact {
// 			//TODO:Handle loss of precision
// 		}
// 		//fast fee = base gas price * 8
// 		fastFee, acc = new(big.Float).Mul(gasPriceGwei, new(big.Float).SetInt64(8)).Float64()
// 		if acc != big.Exact {
// 			//TODO:Handle loos of precision
// 		}
// 	//Klaytn requires gas fee to be exactly 250 at the time of writing - https://docs.klaytn.foundation/klaytn/design/transaction-fees
// 	case "klaytn":
// 		fastFee = 250.0
// 		slowAvgFee = 250.0
// 		fastestFee = 250.0
// 	case "zksync":
// 		fastFee = 250.0
// 		slowAvgFee = 250.0
// 		fastestFee = 250.0
// 		tokenPrice, err = evm.services.Zksync.GetTokenPrice("usd", "")
// 	default:
// 		return nil, status.Error(codes.Unavailable, "Unsupported nonce source")
// 	}
// 	gasPriceInfo = pb.GasPriceInfo{
// 		Fast:        fastFee,
// 		SafeLow:     slowAvgFee,
// 		Fastest:     fastestFee,
// 		Average:     slowAvgFee,
// 		SafeLowWait: 5,
// 		AvgWait:     2,
// 		FastWait:    1,
// 		FastestWait: 0.5,
// 	}

// 	return &pb.NonceResponse{
// 		Nonce:      int64(nonce),
// 		QuoteValue: tokenPrice.Price,
// 		GasPrice:   &gasPriceInfo,
// 		OpL1Fee:    0.0,
// 	}, nil
// }

// func getNonce(evm *evmCore, request *pb.NonceRequest) (*pb.NonceResponse, error) {
// 	nonce, err := getEthTransactionCount(evm, request.Chain, request.Address)
// 	if err != nil {
// 		evm.logger.Error("Error fetching nonce. Err: ", err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	tokenPrice, err := evm.services.CoinGecko.GetTokenExchange("usd", request.Chain)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Error fetching quote price")
// 	}
// 	evm.logger.Info("quote price := ", tokenPrice.Price)
// 	gasPrice, err := evm.services.Debank.GetGasPriceInfo(request.Chain)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	optL1Fee := evm.optContractABI(request.Chain)
// 	gasPriceInfo := pb.GasPriceInfo{
// 		Fast:        float64(gasPrice.Fast),
// 		SafeLow:     float64(gasPrice.Slow),
// 		Fastest:     float64(gasPrice.Fast + (gasPrice.Fast * .2)),
// 		Average:     float64(gasPrice.Normal),
// 		SafeLowWait: 10,
// 		AvgWait:     2,
// 		FastWait:    1,
// 		FastestWait: 0.5,
// 	}
// 	return &pb.NonceResponse{
// 		Nonce:      int64(nonce),
// 		QuoteValue: tokenPrice.Price,
// 		GasPrice:   &gasPriceInfo,
// 		OpL1Fee:    optL1Fee,
// 	}, nil
// }

func (evm *evmCore) SendTransaction(request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	txId, err := evm.rpc[request.Chain].EthSendRawTransaction(request.Msg)
	if err != nil {
		evm.logger.Error("Error while sending raw transaction: ", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	sendTransactionResponse := pb.SendTransactionResponse{
		TransactionId: txId,
	}
	return &sendTransactionResponse, nil
}

// GasLimit RPC call to estimate the required gas to complete a given transaction
func (evm *evmCore) GasLimit(request *pb.GasLimitRequest) (*pb.GasLimitResponse, error) {
	//Pre-checks and customization/hacks
	switch request.Chain {
	case "avalanche":
		if "0xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" == request.To {
			request.To = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	case "xinfin":
		request.From = evm.util.ConvertXdcAddressTo0x(request.From)
		request.To = evm.util.ConvertXdcAddressTo0x(request.To)
	}
	if len(request.Data) == 0 {
		data, _ := evm.createContractABI(request)
		request.Data = string(data)
	}
	transaction := ethrpc.T{
		From:  request.From,
		To:    request.To,
		Gas:   int(request.Gas),
		Value: big.NewInt(request.Value),
		Data:  request.Data,
		Nonce: 0, //Nonce is omitted for gas estimation
	}
	gasLimit, err := evm.rpc[request.Chain].EthEstimateGas(transaction)
	if err != nil {
		evm.logger.Error("Error fetching gas estimate")
		gasLimit = 50000 //TODO:To be refactored
		evm.logger.Error(err)
	}
	//TODO: Move gas limit calculation to ethogo library
	//toAddr := ethgo.HexToAddress(request.To)
	//msg := ethgo.CallMsg{
	//	From:     ethgo.HexToAddress(request.From),
	//	To:       &toAddr,
	//	Data:     []byte(request.Data),
	//	GasPrice: 21000,
	//	Gas:      big.NewInt(request.Gas),
	//	Value:    big.NewInt(request.Value),
	//}
	//gasLimitEthGoRPC, err := evm.ethgoRpc[request.Chain].Eth().EstimateGas(&msg)
	//if err != nil {
	//	evm.logger.Info("Error occurred while fetching gas estimate from ethgo rpc")
	//	evm.logger.Error(err.Error())
	//} else {
	//	evm.logger.Info("No error occured while fetching gas estimate from ethgo rpc")
	//}
	//evm.logger.Info("Gaslimit from ethgo rpc", gasLimitEthGoRPC)
	//evm.logger.Error("Gas limit before, old rpc: ", gasLimit)
	//TODO:Move gas delta calculation to config
	//Customized gas limit calculation
	switch request.Chain {
	case "polygon":
		gasLimit = int(math.Ceil(float64(gasLimit * 3)))
	//Keep both polygon & matic naming convention for backward compatibility
	case "matic":
		gasLimit = int(math.Ceil(float64(gasLimit * 3)))
	case "optimism":
		gasLimit = int(math.Ceil(float64(gasLimit * 2)))
	case "klaytn":
		gasLimit = int(math.Ceil(float64(gasLimit * 5)))
	default:
		gasLimit = int(math.Ceil(float64(gasLimit) * 1.5))
	}
	gasLimitResponse := pb.GasLimitResponse{
		GasLimit:  int64(gasLimit),
		InputData: request.Data,
	}
	evm.logger.Info("GasLimit = ", gasLimit)
	return &gasLimitResponse, nil
}

// getEthTokenDecimals retrieve token decimals from a contract address
func (evm *evmCore) getEthTokenDecimals(contractAddr string, chain string) (int, error) {
	erc20 := erc20.NewERC20(ethgo.HexToAddress(contractAddr), contract.WithJsonRPC(evm.ethgoRpc[chain].Eth()))
	decimals, err := erc20.Decimals()
	if err != nil {
		evm.logger.Info("Error while fetching token decimals", err.Error())
		decimals = 0
	}
	return int(decimals), err
}

func (evm *evmCore) createContractABI(request *pb.GasLimitRequest) ([]byte, error) {
	transferFunction := []string{core.AbiTransferFunction}
	abiContract, err := abi.NewABIFromList(transferFunction)
	if err != nil {
		evm.logger.Error("Method: transfer not found")
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	addr := ethgo.HexToAddress(request.To)
	contractInstance := contract.NewContract(addr, abiContract, contract.WithJsonRPC(evm.ethgoRpc[request.Chain].Eth()))
	method := contractInstance.GetABI().GetMethod("transfer")
	if method == nil {
		evm.logger.Error("Method: transfer not found")
	}
	data, err := method.Encode(request.Value)
	return data, nil
}

// optContractABI get opL1 fee from contact abi
func (evm *evmCore) optContractABI(chain string) float64 {
	if chain != "optimism" {
		//L1 fee not applicable to chains other than optimism
		return 0.0
	}
	abiInstance, err := abi.NewABI(core.OpGasPriceOracleABI)
	if err != nil {
		evm.logger.Errorf("Error generating ABI interface: %v", err)
		return core.DefaultOPL1Fee
	}
	contractInstance := contract.NewContract(ethgo.HexToAddress("0x420000000000000000000000000000000000000F"), abiInstance, contract.WithJsonRPC(evm.ethgoRpc[chain].Eth()))
	method := contractInstance.GetABI().GetMethod(core.OpL1FeeABIMethod)
	if method == nil {
		evm.logger.Error("Error fetching method: getL1Fee")
		return core.DefaultOPL1Fee
	}
	fee, err := contractInstance.Call("getL1Fee", ethgo.Latest, "0x")
	//ABI contract should return a map with only one entry
	if fee == nil || len(fee) != 1 {
		return core.DefaultOPL1Fee
	}
	var formattedFee = core.DefaultOPL1Fee
	for _, v := range fee {
		value := v.(*big.Int)
		//Convert to Gwei
		etherFee := new(big.Float).Quo(new(big.Float).SetInt(value), new(big.Float).SetFloat64(utils.Ether))
		//Multiply by a delta to handle higher fee
		etherFeeMulDelta := new(big.Float).Mul(etherFee, new(big.Float).SetFloat64(2.0))
		finalFee, _ := etherFeeMulDelta.Float64()
		finalFeeStr := strconv.FormatFloat(finalFee, 'f', 5, 64)
		formattedFee, _ = strconv.ParseFloat(finalFeeStr, 64)
	}
	if formattedFee > 0.0 {
		//Handle precision loss
		return formattedFee
	}
	return core.DefaultOPL1Fee
}

func (evm *evmCore) contractABI(request ContractABIRequest) (string, error) {
	abiContract, err := abi.NewABI(core.TokenABI)
	if err != nil {
		evm.logger.Errorf("Error in generating ABI Interface %v", err)
		return "", err
	}
	addr := ethgo.HexToAddress(request.Contract) //contract address
	c := contract.NewContract(addr, abiContract, contract.WithJsonRPC(evm.ethgoRpc[request.Chain].Eth()))
	switch request.Method {
	case "approve":
		//write call
		n := new(big.Int)
		value, ok := n.SetString(request.Data, 10)
		if !ok {
			return "", err
		}
		method := c.GetABI().GetMethod("approve")
		if method == nil {
			evm.logger.Error("Method: approve not found")
		}
		data, err := method.Encode(map[string]interface{}{
			"_spender": ethgo.HexToAddress(request.To),
			"_value":   value,
		})
		if err != nil {
			evm.logger.Error(err)
		}
		return fmt.Sprintf("0x%x", data), nil
	case "allowance":
		//Read call
		res, err := c.Call("allowance", ethgo.Latest, ethgo.HexToAddress(request.From), ethgo.HexToAddress(request.To))
		if err != nil {
			evm.logger.Error(err)
			return "", status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}
		return res["amount"].(*big.Int).String(), nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "", "Unsupported Method")
	}
}

// func (evm *evmCore) TokenApprove(request *pb.ApprovalRequest) (*pb.ApprovalResponse, error) {
// 	var res pb.ApprovalResponse
// 	var gasPrice string
// 	var blockNativeResponse *BlockNativeGasPrice
// 	if request.Chain == "xinfin" {
// 		request.Token = evm.util.ResolveXDCAddress(request.Token)
// 	}
// 	data, err := evm.contractABI(ContractABIRequest{
// 		Contract: request.Token,
// 		To:       request.Target,
// 		Data:     "79228162514264337593543950335", //Maximum Decimal Value
// 		Method:   "approve",
// 		Chain:    request.Chain,
// 	})
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	gasData := data
// 	transaction := ethrpc.T{
// 		From: request.Token,
// 		To:   request.Target,
// 		Data: gasData,
// 	}
// 	gasLimit, err := evm.rpc[request.Chain].EthEstimateGas(transaction) //TODO: EthEstimateGas need to be generic
// 	if err != nil {
// 		evm.logger.Error("Error fetching gas estimate")
// 		gasLimit = 60000 //Constant Gas Limit Value
// 		evm.logger.Error(err)
// 	}
// 	evm.logger.Error("Gas limit before: ", gasLimit)
// 	switch request.Chain {
// 	case "boba":
// 		gasLimit = int(math.Ceil(float64(gasLimit * 5)))
// 	case "aurora":
// 		gasLimit = int(math.Ceil(float64(gasLimit * 2)))
// 	case "arbitrum":
// 		gasLimit = int(math.Ceil(float64(gasLimit * 20)))
// 	default:
// 		gasLimit = int(math.Ceil(float64(gasLimit) * 1.2))
// 	}
// 	if request.Chain != "bsc" {
// 		body, err := evm.httpRequest.GetRequestWithHeaders(evm.env.BlockNative.EndPoint+"/gasprices/blockprices",
// 			"Authorization", evm.env.BlockNative.AuthHeader)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := json.Unmarshal(body, &blockNativeResponse); err != nil {
// 			evm.logger.Error(err)
// 			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 		}
// 		if len(blockNativeResponse.BlockPrices) != 0 {
// 			gasPrice = fmt.Sprintf("%v", blockNativeResponse.BlockPrices[0].EstimatedPrices[0].Price)
// 		}
// 	}
// 	if gasPrice != "" {
// 		res.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s&gasPrice=%v", request.Token, "0", gasData, strconv.Itoa(gasLimit), gasPrice)
// 	} else {
// 		res.TxLink = fmt.Sprintf("https://txlink.io/tx?to=%s&value=%s&data=%s&gaslimit=%s", request.Token, "0", gasData, strconv.Itoa(gasLimit))
// 	}
// 	res.To = request.Token
// 	res.Value = "0"
// 	res.Data = gasData
// 	res.GasLimit = strconv.Itoa(gasLimit)
// 	return &res, nil

// }
func (evm *evmCore) isChainNativeToken(tokenAddress string, chain string) bool {
	var isNativeChain bool
	for _, c := range evm.env.EVM.Cfg.Wallets {
		if chain == c.ChainName && tokenAddress == c.NativeTokenInfo.Address {
			isNativeChain = true
			return isNativeChain
		} else {
			isNativeChain = false
		}
	}
	return isNativeChain
}

func (evm *evmCore) GetTokenAllowance(request *pb.AllowanceRequest) (*pb.AllowanceResponse, error) {
	allowanceForNativeToken := 999999999999 //standard value for native token
	var allowance string
	if request.Chain == "xinfin" {
		request.Contract = evm.util.ResolveXDCAddress(request.Contract)
		request.Owner = evm.util.ResolveXDCAddress(request.Owner)
		if request.Contract == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
			request.Contract = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
	}
	isNativeChain := evm.isChainNativeToken(strings.ToLower(request.Contract), request.Chain)
	if isNativeChain {
		allowance = strconv.Itoa(allowanceForNativeToken)
	} else {
		data, err := evm.contractABI(ContractABIRequest{
			From:     request.Owner,
			To:       request.Spender,
			Contract: request.Contract,
			Chain:    request.Chain,
			Method:   "allowance",
		})
		if err != nil {
			return nil, err
		}
		allowance = data
	}
	return &pb.AllowanceResponse{
		Allowance: evm.util.ToDecimal(allowance, 18).String(),
	}, nil
}

// GetTxStatus gets transaction status by making an RPC call to transaction receipt
func (evm *evmCore) GetTxStatus(request *pb.TxStatusRequest) (*pb.TxStatusResponse, error) {
	txReceipt, err := evm.rpc[request.Chain].EthGetTransactionReceipt(request.TxHash)
	var logs = make([]*pb.Log, 0)
	pendingTxReceipt := &pb.TxStatusResponse{
		TransactionHash:   "",
		TransactionIndex:  0,
		BlockHash:         "",
		BlockNumber:       0,
		CumulativeGasUsed: 0,
		GasUsed:           0,
		ContractAddress:   "",
		Logs:              logs,
		LogsBloom:         "",
		Root:              "",
		Status:            "PENDING",
	}
	//TxStatus API pooling from clients require a 200 response for tx in mem pool
	if err != nil || txReceipt == nil {
		evm.logger.Info("Transaction pending/unknown")
		evm.logger.Error("Error: ", err.Error())
		return pendingTxReceipt, status.Errorf(codes.OK, err.Error(), "Transaction pending/unknown")
	}
	response := pb.TxStatusResponse{}
	response = pb.TxStatusResponse{
		TransactionHash:   Clean(txReceipt.TransactionHash).(string),
		TransactionIndex:  int64(txReceipt.TransactionIndex),
		BlockHash:         Clean(txReceipt.BlockHash).(string),
		BlockNumber:       int64(txReceipt.BlockNumber),
		CumulativeGasUsed: int64(txReceipt.CumulativeGasUsed),
		GasUsed:           int64(txReceipt.GasUsed),
		ContractAddress:   Clean(txReceipt.ContractAddress).(string),
		LogsBloom:         Clean(txReceipt.LogsBloom).(string),
		Root:              txReceipt.Root,
		Status:            TxStatus()[txReceipt.Status],
	}
	//Copy transaction logs
	if txReceipt.Logs != nil {
		for _, log2 := range txReceipt.Logs {
			var log = pb.Log{}
			log.Removed = log2.Removed
			log.LogIndex = int64(log2.LogIndex)
			log.TransactionIndex = int64(log2.TransactionIndex)
			log.TransactionHash = log2.TransactionHash
			log.BlockNumber = int64(log2.BlockNumber)
			log.BlockHash = log2.BlockHash
			log.Address = log2.Address
			log.Data = log2.Data
			log.Topics = log2.Topics
			logs = append(logs, &log)
		}
	}
	response.Logs = logs
	return &response, nil
}

func (evm *evmCore) GetProcessingFee(request *pb.ProcessingFeeRequest) (*pb.ProcessingFeeResponse, error) {
	return nil, nil
}

// // GetUserData Retrieve user data from on the provided source
// func (evm *evmCore) GetUserData(request *pb.UserDataRequest) (*pb.UserDataResponse, error) {
// 	source := evm.util.GetWalletSource(request.Chain)
// 	switch source.UserDataSource {
// 	case "unmarshal":
// 		return transformUnmarshalUserData(evm, request)
// 	//Hack considering unmarshal as the only data source.
// 	//TODO: common utilities doesn't need adapter routing.
// 	case "":
// 		return transformUnmarshalUserData(evm, request)
// 	default:
// 		return nil, status.Error(codes.Unavailable, "Unsupported Operation")
// 	}
// }

// // transformUnmarshalUserData Transforms unmarshal data model into user data response model
// func transformUnmarshalUserData(evm *evmCore, request *pb.UserDataRequest) (*pb.UserDataResponse, error) {
// 	userData, err := evm.services.Unmarshall.GetUserData(request, request.Chain)
// 	if err != nil {
// 		evm.logger.Error("Error fetching user data. Err: ", err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	userDataResponse := pb.UserDataResponse{
// 		QuoteRate:              userData.QuoteRate,
// 		TotalFeesPaid:          userData.TotalFeesPaid,
// 		TotalFeesPaidUsd:       userData.TotalFeesPaidUsd,
// 		AverageTokenPrice:      userData.AverageTokenPrice,
// 		OverallProfitLoss:      userData.OverallProfitLoss,
// 		CurrentHoldingQuantity: userData.CurrentHoldingQuantity,
// 		PercentageChange_24H:   userData.PercentageChange24H,
// 		PriceChange_24H:        userData.PriceChange24H,
// 	}
// 	return &userDataResponse, nil
// }

// func (evm *evmCore) GetNftCollections(request *pb.NftCollectionRequest) (*pb.ListNftCollectionResponse, error) {
// 	var endpoint = "https://api.opensea.io"
// 	url := fmt.Sprintf("%s%s%s%s%s%s%s", endpoint, "/api/v1/assets?owner=", request.Address, "&offset=", request.Page, "&limit=", request.PageSize)
// 	body, err := evm.httpRequest.GetRequestWithHeaders(url, "X-API-KEY", "8bb16c6088134b3fb77903d39aad0cce")
// 	if err != nil {
// 		evm.logger.Error("NftCollectionsHandler info Logging Error  is : %v", err.Error())
// 	}
// 	var jsonResponseStruct *unmarshal.NFTCollectionDataModel
// 	err = json.Unmarshal(body, &jsonResponseStruct)
// 	if err != nil {
// 		evm.logger.Error(" NftCollectionsHandler Logging Error  is : %v", err.Error())
// 		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 	}

// 	var nftCollections []*pb.NftCollectionResponse
// 	for _, asset := range jsonResponseStruct.Assets {
// 		collection := asset.Collection
// 		collection_item := pb.NftCollectionResponse{
// 			BannerImageUrl:          collection.BannerImageUrl,
// 			ChatUrl:                 collection.ChatUrl,
// 			CreatedDate:             collection.CreatedDate,
// 			DefaultToFiat:           collection.DefaultToFiat,
// 			Description:             collection.Description,
// 			DevBuyerFeeBasisPoints:  collection.DevBuyerFeeBasisPoints,
// 			DevSellerFeeBasisPoints: collection.DevSellerFeeBasisPoints,
// 			DiscordUrl:              collection.DiscordUrl,
// 			DisplayData: &pb.NFTDisplayData{
// 				CardDisplayStyle: collection.DisplayData.CardDisplayStyle,
// 			},
// 			ExternalUrl:                 collection.ExternalUrl,
// 			Featured:                    collection.Featured,
// 			FeaturedImageUrl:            collection.FeaturedImageUrl,
// 			Hidden:                      collection.Hidden,
// 			SafelistRequestStatus:       collection.SafelistRequestStatus,
// 			ImageUrl:                    collection.ImageUrl,
// 			IsSubjectToWhitelist:        collection.IsSubjectToWhitelist,
// 			LargeImageUrl:               collection.LargeImageUrl,
// 			MediumUsername:              collection.MediumUsername,
// 			Name:                        collection.Name,
// 			OnlyProxiedTransfers:        collection.OnlyProxiedTransfers,
// 			OpenseaBuyerFeeBasisPoints:  collection.OpenseaBuyerFeeBasisPoints,
// 			OpenseaSellerFeeBasisPoints: collection.OpenseaSellerFeeBasisPoints,
// 			PayoutAddress:               collection.PayoutAddress,
// 			RequireEmail:                collection.RequireEmail,
// 			ShortDescription:            collection.ShortDescription,
// 			Slug:                        collection.Slug,
// 			TelegramUrl:                 collection.TelegramUrl,
// 			TwitterUsername:             collection.TwitterUsername,
// 			InstagramUsername:           collection.InstagramUsername,
// 			WikiUrl:                     collection.WikiUrl,
// 			IsNsfw:                      collection.IsNsfw,
// 			NftData: []*pb.NftData{
// 				&pb.NftData{
// 					Id:                   asset.Id,
// 					NumSales:             asset.NumSales,
// 					BackgroundColor:      asset.BackgroundColor,
// 					ImageUrl:             asset.ImageUrl,
// 					ImagePreviewUrl:      asset.ImagePreviewUrl,
// 					ImageThumbnailUrl:    asset.ImageThumbnailUrl,
// 					ImageOriginalUrl:     asset.ImageOriginalUrl,
// 					AnimationUrl:         asset.AnimationUrl,
// 					AnimationOriginalUrl: asset.AnimationOriginalUrl,
// 					Name:                 asset.Name,
// 					Description:          asset.Description,
// 					ExternalLink:         asset.ExternalLink,
// 					AssetContract:        &pb.NftDataAssetContract{},
// 					Permalink:            asset.Permalink,
// 					Decimals:             asset.Decimals,
// 					TokenMetadata:        asset.TokenMetadata,
// 					IsNsfw:               asset.IsNsfw,
// 					Owner:                &pb.NftDataOwner{},
// 					SellOrders: []*pb.NftDataSellOrders{
// 						&pb.NftDataSellOrders{},
// 					},
// 					SeaportSellOrders: asset.SeaportSellOrders,
// 					Creator: &pb.NftDataCreator{
// 						User: &pb.NftDataUser{
// 							Username: asset.Creator.User.Username,
// 						},
// 						ProfileImgUrl: asset.Creator.ProfileImgUrl,
// 						Address:       asset.Creator.Address,
// 						Config:        asset.Creator.Config,
// 					},
// 					Traits:                  GetTraits(asset.Traits),
// 					LastSale:                nil,
// 					TopBid:                  asset.TopBid,
// 					ListingDate:             asset.ListingDate,
// 					IsPresale:               asset.IsPresale,
// 					TransferFeePaymentToken: asset.TransferFeePaymentToken,
// 					TransferFee:             asset.TransferFee,
// 					TokenId:                 asset.TokenId,
// 					CollectionName:          collection.Name,
// 					ContractAddress:         asset.AssetContract.Address,
// 				},
// 			},
// 		}
// 		nftCollections = append(nftCollections, &collection_item)
// 	}

// 	return &pb.ListNftCollectionResponse{
// 		Nft: nftCollections,
// 	}, nil
// }

// func (evm *evmCore) BulkApproval(request *pb.ApprovalRequest) (*pb.BulkApprovalResponse, error) {
// 	var bulkApprovalResponse pb.BulkApprovalResponse
// 	var tokens = strings.Split(request.Token, ",")
// 	chanResponse := make(chan *pb.ApprovalResponse)
// 	var listApprovalResponse []*pb.ApprovalResponse
// 	wg := new(sync.WaitGroup)
// 	for _, value := range tokens {
// 		wg.Add(1)
// 		req := &pb.ApprovalRequest{
// 			Target: request.Target,
// 			Chain:  request.Chain,
// 			Token:  strings.Trim(value, " "),
// 		}
// 		go evm.GetTokenApproval(req, wg, chanResponse)
// 		listApprovalResponse = append(listApprovalResponse, <-chanResponse)
// 	}
// 	bulkApprovalResponse.Response = listApprovalResponse
// 	return &bulkApprovalResponse, nil
// }

// func (evm *evmCore) BulkAllowance(request *pb.AllowanceRequest) (*pb.BulkAllowanceResponse, error) {
// 	var bulkAllowanceResponse pb.BulkAllowanceResponse
// 	var tokens = strings.Split(request.Contract, ",")
// 	chanResponse := make(chan *pb.AllowanceResponse)
// 	var listAllowanceResponse []*pb.AllowanceResponse
// 	wg := new(sync.WaitGroup)
// 	for _, value := range tokens {
// 		wg.Add(1)
// 		req := &pb.AllowanceRequest{
// 			Chain:    request.Chain,
// 			Contract: strings.Trim(value, " "),
// 			Owner:    request.Owner,
// 			Spender:  request.Spender,
// 		}
// 		go evm.GetBulkTokenAllowance(req, wg, chanResponse)
// 		listAllowanceResponse = append(listAllowanceResponse, <-chanResponse)
// 	}
// 	bulkAllowanceResponse.Response = listAllowanceResponse
// 	return &bulkAllowanceResponse, nil
// }

type Traits []struct {
	TraitType   string `json:"trait_type"`
	Value       string `json:"value"`
	DisplayType string `json:"display_type"`
	MaxValue    string `json:"max_value"`
	TraitCount  int64  `json:"trait_count"`
	Order       string `json:"order"`
}

func GetTraits(t Traits) []*pb.NftDataTraits {
	var trait []*pb.NftDataTraits
	for _, traititem := range t {
		trait = append(trait, &pb.NftDataTraits{
			TraitType:   traititem.TraitType,
			Value:       traititem.Value,
			DisplayType: traititem.DisplayType,
			MaxValue:    traititem.MaxValue,
			TraitCount:  traititem.TraitCount,
			Order:       traititem.Order,
		})
	}
	return trait
}

// func (evm *evmCore) GetTokenApproval(req *pb.ApprovalRequest, wg *sync.WaitGroup, chanResponse chan *pb.ApprovalResponse) {
// 	defer wg.Done()
// 	tokenApprovalResponse, _ := evm.TokenApprove(req)
// 	chanResponse <- tokenApprovalResponse
// 	return
// }

func (evm *evmCore) GetBulkTokenAllowance(req *pb.AllowanceRequest, wg *sync.WaitGroup, chanResponse chan *pb.AllowanceResponse) {
	defer wg.Done()
	tokenAllowanceResponse, _ := evm.GetTokenAllowance(req)
	chanResponse <- tokenAllowanceResponse
	return
}
// func (evm *evmCore) GetOpportunites(request *pb.GetOpportunitiesRequest) (*pb.GetOpportunitesResponse, error) {
// 	var opportunityResForEVM Opportunities
// 	var opportunities models.Opportunities
// 	reqUrl := fmt.Sprintf(evm.env.PROXIES_ENDPOINT+"/v1/stake/opportunity-new?current=%s", request.Chain)
// 	res, err := evm.httpRequest.GetRequest(reqUrl)
// 	if err != nil {
// 		evm.logger.Error(err)
// 		return nil, status.Errorf(codes.Internal, err.Error(), "error at fetching response from v1")
// 	}
// 	err = json.Unmarshal(res, &opportunityResForEVM)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, current := range opportunityResForEVM.Current {
// 		switch current.Apr.(type) {
// 		case float64:
// 			current.Apr = strconv.FormatFloat(current.Apr.(float64), 'f', -1, 64)
// 			opportunityData := models.OpportunityData{
// 				Apr:                         current.Apr.(string),
// 				Chain:                       current.Chain,
// 				Logo:                        current.Logo,
// 				StakeTokenName:              current.StakeTokenName,
// 				ReceiptTokenName:            current.ReceiptTokenName,
// 				ContractDecimals:            current.ContractDecimals,
// 				StakeTokenLogoUrl:           current.StakeTokenLogoUrl,
// 				StakeTokenContractAddress:   current.StakeTokenContractAddress,
// 				ReceiptTokenLogoUrl:         current.ReceiptTokenLogoUrl,
// 				ReceiptTokenContractAddres:  current.ReceiptTokenContractAddres,
// 				ReceiptTokenContractAddress: current.ReceiptTokenContractAddress,
// 				StakeToReceiptExchangeRate:  current.StakeToReceiptExchangeRate,
// 				ReceiptToStakeExchangeRate:  current.ReceiptToStakeExchangeRate,
// 				QuoteRate:                   current.QuoteRate,
// 				ReceiptQuoteRate:            current.ReceiptQuoteRate,
// 				StakingType:                 current.StakingType,
// 				ProtocolName:                current.ProtocolName,
// 				CoolDownPeriod:              current.CoolDownPeriod,
// 				MinLockup:                   current.MinLockup,
// 				RewardSchedule:              current.RewardSchedule,
// 				TokenName:                   current.TokenName,
// 			}
// 			opportunities.Current = append(opportunities.Current, opportunityData)
// 		default:
// 			opportunityData := models.OpportunityData{
// 				Apr:                         current.Apr.(string),
// 				Chain:                       current.Chain,
// 				Logo:                        current.Logo,
// 				StakeTokenName:              current.StakeTokenName,
// 				ReceiptTokenName:            current.ReceiptTokenName,
// 				ContractDecimals:            current.ContractDecimals,
// 				StakeTokenLogoUrl:           current.StakeTokenLogoUrl,
// 				StakeTokenContractAddress:   current.StakeTokenContractAddress,
// 				ReceiptTokenLogoUrl:         current.ReceiptTokenLogoUrl,
// 				ReceiptTokenContractAddres:  current.ReceiptTokenContractAddres,
// 				ReceiptTokenContractAddress: current.ReceiptTokenContractAddress,
// 				StakeToReceiptExchangeRate:  current.StakeToReceiptExchangeRate,
// 				ReceiptToStakeExchangeRate:  current.ReceiptToStakeExchangeRate,
// 				QuoteRate:                   current.QuoteRate,
// 				ReceiptQuoteRate:            current.ReceiptQuoteRate,
// 				StakingType:                 current.StakingType,
// 				ProtocolName:                current.ProtocolName,
// 				CoolDownPeriod:              current.CoolDownPeriod,
// 				MinLockup:                   current.MinLockup,
// 				RewardSchedule:              current.RewardSchedule,
// 				TokenName:                   current.TokenName,
// 			}
// 			opportunities.Current = append(opportunities.Current, opportunityData)
// 		}
// 	}
// 	for _, others := range opportunityResForEVM.Others {
// 		switch others.Apr.(type) {
// 		case float64:
// 			others.Apr = strconv.FormatFloat(others.Apr.(float64), 'f', -1, 64)
// 			opportunityData := models.OpportunityData{
// 				Apr:                         others.Apr.(string),
// 				Chain:                       others.Chain,
// 				Logo:                        others.Logo,
// 				StakeTokenName:              others.StakeTokenName,
// 				ReceiptTokenName:            others.ReceiptTokenName,
// 				ContractDecimals:            others.ContractDecimals,
// 				StakeTokenLogoUrl:           others.StakeTokenLogoUrl,
// 				StakeTokenContractAddress:   others.StakeTokenContractAddress,
// 				ReceiptTokenLogoUrl:         others.ReceiptTokenLogoUrl,
// 				ReceiptTokenContractAddres:  others.ReceiptTokenContractAddres,
// 				ReceiptTokenContractAddress: others.ReceiptTokenContractAddress,
// 				StakeToReceiptExchangeRate:  others.StakeToReceiptExchangeRate,
// 				ReceiptToStakeExchangeRate:  others.ReceiptToStakeExchangeRate,
// 				QuoteRate:                   others.QuoteRate,
// 				ReceiptQuoteRate:            others.ReceiptQuoteRate,
// 				StakingType:                 others.StakingType,
// 				ProtocolName:                others.ProtocolName,
// 				CoolDownPeriod:              others.CoolDownPeriod,
// 				MinLockup:                   others.MinLockup,
// 				RewardSchedule:              others.RewardSchedule,
// 				TokenName:                   others.TokenName,
// 			}
// 			opportunities.Others = append(opportunities.Others, opportunityData)
// 		default:
// 			opportunityData := models.OpportunityData{
// 				Apr:                         others.Apr.(string),
// 				Chain:                       others.Chain,
// 				Logo:                        others.Logo,
// 				StakeTokenName:              others.StakeTokenName,
// 				ReceiptTokenName:            others.ReceiptTokenName,
// 				ContractDecimals:            others.ContractDecimals,
// 				StakeTokenLogoUrl:           others.StakeTokenLogoUrl,
// 				StakeTokenContractAddress:   others.StakeTokenContractAddress,
// 				ReceiptTokenLogoUrl:         others.ReceiptTokenLogoUrl,
// 				ReceiptTokenContractAddres:  others.ReceiptTokenContractAddres,
// 				ReceiptTokenContractAddress: others.ReceiptTokenContractAddress,
// 				StakeToReceiptExchangeRate:  others.StakeToReceiptExchangeRate,
// 				ReceiptToStakeExchangeRate:  others.ReceiptToStakeExchangeRate,
// 				QuoteRate:                   others.QuoteRate,
// 				ReceiptQuoteRate:            others.ReceiptQuoteRate,
// 				StakingType:                 others.StakingType,
// 				ProtocolName:                others.ProtocolName,
// 				CoolDownPeriod:              others.CoolDownPeriod,
// 				MinLockup:                   others.MinLockup,
// 				RewardSchedule:              others.RewardSchedule,
// 				TokenName:                   others.TokenName,
// 			}
// 			opportunities.Others = append(opportunities.Others, opportunityData)
// 		}
// 	}

// 	marshalRes, err := json.Marshal(opportunities)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.GetOpportunitesResponse{
// 		Opportunities: marshalRes,
// 	}, nil
// }
