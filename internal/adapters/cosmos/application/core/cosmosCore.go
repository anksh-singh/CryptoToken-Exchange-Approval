package core

import (
	"bridge-allowance/config"
// 	// "bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	// "bridge-allowance/utils/models"
	"bytes"
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	_ "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type Handler struct {
	env             *config.Config
	// coingecko       *coingecko.CoinGecko
	logger          *zap.SugaredLogger
	rpcClientMap    map[string]jsonrpc.RPCClient
	chainConfigData map[string]config.CosmosWallets
	httpRequest     utils.IHttpRequest
	utils           *utils.UtilConf
	helper          *utils.Helpers
	ibcTokenInfo    map[string]IBCData
	denomInfo       map[string]DenomInfo
}

type ConfMap struct {
	rpcClientMap    map[string]jsonrpc.RPCClient
	chainConfigData map[string]config.CosmosWallets
	ibcTokenInfo    map[string]IBCData
	denomInfo       map[string]DenomInfo
}

func NewCosmosHandler(config *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest,
	utilsConf *utils.UtilConf) (*Handler, error) {
	confMap, err := getRpcAndRestInfo(config.Cosmos, logger, httpRequest)
	helper := utils.Helpers{}
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &Handler{
		config,
		// coingecko,
		logger,
		confMap.rpcClientMap,
		confMap.chainConfigData,
		httpRequest,
		utilsConf,
		&helper,
		confMap.ibcTokenInfo,
		confMap.denomInfo,
	}, nil
}

// ListTransaction Lists the transactions activity for an address
// TODO: Add support for 3P services
func (cosmos *Handler) ListTransaction(request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	source := cosmos.utils.GetCosmosWalletSource(request.Chain)
	switch source.HistorySource {
	case "custom":
		return getCustomTxHistory(cosmos, request)
	default:
		return nil, status.Errorf(codes.Unavailable, "Unsupported Operation", "Unsupported Operation")
	}
}

func getCustomTxHistory(cosmos *Handler, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	var txResponse pb.ListTransactionResponse
	txResponse.Transactions = make([]*pb.TransactionData, 0)
	switch request.Chain {
	case "kava":
		return getCosmosTxHistory(cosmos, request)
	default:
		return nil, status.Errorf(codes.Unavailable, "Unsupported Chain", "Unsupported Chain")
	}
}
func getCosmosTxHistory(cosmos *Handler, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)

	var sentFieldsList []*pb.TransactionInfo
	var receivedFieldsList []*pb.TransactionInfo
	var othersList []*pb.TransactionInfo
	// var transactionData *pb.TransactionData
	var transactionType = ""
	var fromAddress = ""
	var toAddress = ""
	var statusTxn = ""
	var remarks = ""
	var value = "0"
	var pageNumber, _ = strconv.ParseInt(request.Page, 0, 64)
	var pageSize, _ = strconv.ParseInt(request.PageSize, 0, 64)

	var pageOffSet = (pageNumber - 1) * pageSize
	reqUrlReceipient := fmt.Sprintf(walletInfo.REST+"/cosmos/tx/v1beta1/txs?pagination.limit=%v&pagination.offset=%v&orderBy=ORDER_BY_DESC&events=transfer.recipient='%v'", request.PageSize, pageOffSet, request.Address)
	reqUrlSender := fmt.Sprintf(walletInfo.REST+"/cosmos/tx/v1beta1/txs?pagination.limit=%v&pagination.offset=%v&orderBy=ORDER_BY_DESC&events=transfer.sender='%v'", request.PageSize, pageOffSet, request.Address)

	bodyRec, errRec := cosmos.httpRequest.GetRequest(reqUrlReceipient)
	bodySend, errSend := cosmos.httpRequest.GetRequest(reqUrlSender)
	var cosmosScanRec *CosmosScan
	var cosmosScanSend *CosmosScan
	if errRec != nil {
		cosmos.logger.Error("Error fetching history. Err: ", errRec)
		// return nil, errRec
	} else {
		errRec = json.Unmarshal(bodyRec, &cosmosScanRec)
	}
	if errSend != nil {
		cosmos.logger.Error("Error fetching history. Err: ", errRec)
		// return nil, errRec
	} else {
		errSend = json.Unmarshal(bodySend, &cosmosScanSend)
	}
	var transactionResponse pb.ListTransactionResponse
	transactionResponse.Transactions = make([]*pb.TransactionData, 0)
	cosmos.logger.Info("before marshaling")

	cosmos.logger.Info("before marshaling")
	if errRec != nil {
		cosmos.logger.Error("Error fetching history. Err: ", errRec)
		// return nil, errRec
	}
	if errSend != nil {
		cosmos.logger.Error("Error fetching history. Err: ", errSend)
		// return nil, errSend
	}
	//Set pagination properties
	var totalTxnsRec, _ = strconv.ParseInt(cosmosScanRec.Pagination.Total, 0, 64)
	var totalTxnsSend, _ = strconv.ParseInt(cosmosScanSend.Pagination.Total, 0, 64)
	var countOfRecordsRec = len(cosmosScanRec.TxResponses)
	var countOfRecordsSend = len(cosmosScanSend.TxResponses)
	transactionResponse.TotalTxs = totalTxnsRec + totalTxnsSend
	transactionResponse.Page, _ = strconv.ParseInt(request.Page, 0, 64)
	transactionResponse.ItemsOnPage = int64(countOfRecordsRec + countOfRecordsSend)
	quotient, remainder := (totalTxnsRec+totalTxnsSend)/int64(countOfRecordsRec+countOfRecordsSend), (totalTxnsRec+totalTxnsSend)%int64(countOfRecordsRec+countOfRecordsSend)
	if remainder > 0 {
		quotient++
	}
	transactionResponse.TotalPages = quotient
	// Format response data as per *pb.AssetResponse contract
	for _, item := range cosmosScanRec.TxResponses {
		transactionType = ""
		fromAddress = ""
		toAddress = ""
		statusTxn = ""
		remarks = ""
		value = "0"
		sentFieldsList = []*pb.TransactionInfo{}
		receivedFieldsList = []*pb.TransactionInfo{}
		othersList = []*pb.TransactionInfo{}
		//TODO: Add description

		txReceipt := cosmosTxReceipt(cosmos, item.Txhash, walletInfo)
		if txReceipt.TxResponse.Code == 0 {
			statusTxn = "completed"
		} else {
			statusTxn = "error"
		}

		remarks = txReceipt.Tx.Body.Memo
		var counterMessages = -1
		for _, msg := range txReceipt.Tx.Body.Messages {
			counterMessages++
			if strings.Contains(msg.Type, "MsgDelegate") {
				var cosmosAmount = msg.Amount.(map[string]interface{})
				transactionType = "delegate"
				sentFields := pb.TransactionInfo{
					Name:      cosmosAmount["denom"].(string),
					Symbol:    cosmosAmount["denom"].(string),
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount["amount"].(string),
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				sentFieldsList = append(sentFieldsList, &sentFields)
				fromAddress = msg.DelegatorAddress
				toAddress = msg.ValidatorAddress
				value = cosmosAmount["amount"].(string)
			} else if strings.Contains(msg.Type, "MsgUndelegate") {
				var cosmosAmount = msg.Amount.(CosmosAmount)
				transactionType = "undelegate"
				receivedFields := pb.TransactionInfo{
					Name:      cosmosAmount.Denom,
					Symbol:    cosmosAmount.Denom,
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount.Amount,
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				receivedFieldsList = append(receivedFieldsList, &receivedFields)
				fromAddress = msg.ValidatorAddress
				toAddress = msg.DelegatorAddress
				value = cosmosAmount.Amount
			} else if strings.Contains(msg.Type, "MsgWithdrawDelegatorReward") {
				transactionType = "claimrewards"
				var denomLocal = ""
				var amountLocal = ""
				var log = txReceipt.TxResponse.Logs[counterMessages]
				for _, event := range log.Events {
					if event.Type == "coin_received" {
						for _, attr := range event.Attributes {
							if attr.Key == "amount" {
								amountLocal = string(attr.Value[0:strings.Index(attr.Value, "ukava")])
								denomLocal = "ukava"
								receivedFields := pb.TransactionInfo{
									Name:      denomLocal,
									Symbol:    denomLocal,
									TokenId:   "",
									Decimals:  walletInfo.Decimals,
									Value:     amountLocal,
									QuoteRate: 0,
									LogoUrl:   "",
									From:      msg.ValidatorAddress,
									To:        msg.DelegatorAddress,
								}
								receivedFieldsList = append(receivedFieldsList, &receivedFields)
								fromAddress = msg.ValidatorAddress
								toAddress = msg.DelegatorAddress
								value = amountLocal
							}
						}
					}
				}

			} else if strings.Contains(msg.Type, "MsgBeginRedelegate") {
				var cosmosAmount = msg.Amount.(CosmosAmount)
				transactionType = "delegate"
				sentFields := pb.TransactionInfo{
					Name:      cosmosAmount.Denom,
					Symbol:    cosmosAmount.Denom,
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount.Amount,
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				sentFieldsList = append(sentFieldsList, &sentFields)
				fromAddress = msg.DelegatorAddress
				toAddress = msg.ValidatorAddress
				value = cosmosAmount.Amount
			} else if strings.Contains(msg.Type, "MsgSend") {
				if request.Address == msg.FromAddress {
					var cosmosAmount = msg.Amount.([]CosmosAmount)
					transactionType = "send"
					var amountEnt = cosmosAmount[0]
					sentFields := pb.TransactionInfo{
						Name:      amountEnt.Denom,
						Symbol:    amountEnt.Denom,
						TokenId:   "",
						Decimals:  walletInfo.Decimals,
						Value:     amountEnt.Amount,
						QuoteRate: 0,
						LogoUrl:   "",
						From:      msg.DelegatorAddress,
						To:        msg.ValidatorAddress,
					}
					sentFieldsList = append(sentFieldsList, &sentFields)
					fromAddress = msg.DelegatorAddress
					toAddress = msg.ValidatorAddress
					value = amountEnt.Amount
				} else if request.Address == msg.FromAddress {
					var cosmosAmount = msg.Amount.([]CosmosAmount)
					transactionType = "send"
					var amountEnt = cosmosAmount[0]
					receivedFields := pb.TransactionInfo{
						Name:      amountEnt.Denom,
						Symbol:    amountEnt.Denom,
						TokenId:   "",
						Decimals:  walletInfo.Decimals,
						Value:     amountEnt.Amount,
						QuoteRate: 0,
						LogoUrl:   "",
						From:      msg.DelegatorAddress,
						To:        msg.ValidatorAddress,
					}
					receivedFieldsList = append(receivedFieldsList, &receivedFields)
					fromAddress = msg.DelegatorAddress
					toAddress = msg.ValidatorAddress
					value = amountEnt.Amount
				}
			} else if strings.Contains(msg.Type, "MsgRequestData") {
				transactionType = "request"
				fromAddress = request.Address
				toAddress = ""
				value = "0"
			} else if strings.Contains(msg.Type, "MsgExec") {
				transactionType = "report"
				fromAddress = request.Address
				toAddress = ""
				value = "0"
			}

		}
		var blockNumber, _ = strconv.ParseInt(item.Height, 0, 64)
		transactionData := pb.TransactionData{
			Id:                  Clean(item.Txhash).(string),
			From:                Clean(fromAddress).(string),
			To:                  Clean(toAddress).(string),
			Fee:                 "0",
			Date:                item.Timestamp.Unix(),
			Status:              statusTxn,
			Type:                transactionType,
			Block:               blockNumber,
			Value:               value,
			Nonce:               0,
			NativeTokenDecimals: walletInfo.Decimals,
			Description:         remarks,
			Sent:                sentFieldsList,
			Received:            receivedFieldsList,
			Others:              othersList,
		}
		transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)

	}

	for _, item := range cosmosScanSend.TxResponses {
		transactionType = ""
		fromAddress = ""
		toAddress = ""
		statusTxn = ""
		remarks = ""
		value = "0"
		sentFieldsList = []*pb.TransactionInfo{}
		receivedFieldsList = []*pb.TransactionInfo{}
		othersList = []*pb.TransactionInfo{}
		//TODO: Add description

		txReceipt := cosmosTxReceipt(cosmos, item.Txhash, walletInfo)
		if txReceipt.TxResponse.Code == 0 {
			statusTxn = "completed"
		} else {
			statusTxn = "error"
		}

		remarks = txReceipt.Tx.Body.Memo
		var counterMessages = -1
		for _, msg := range txReceipt.Tx.Body.Messages {
			counterMessages++
			if strings.Contains(msg.Type, "MsgDelegate") {
				var cosmosAmount = msg.Amount.(map[string]interface{})
				transactionType = "delegate"
				sentFields := pb.TransactionInfo{
					Name:      cosmosAmount["denom"].(string),
					Symbol:    cosmosAmount["denom"].(string),
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount["amount"].(string),
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				sentFieldsList = append(sentFieldsList, &sentFields)
				fromAddress = msg.DelegatorAddress
				toAddress = msg.ValidatorAddress
				value = cosmosAmount["amount"].(string)
			} else if strings.Contains(msg.Type, "MsgUndelegate") {
				var cosmosAmount = msg.Amount.(map[string]interface{})
				transactionType = "undelegate"
				receivedFields := pb.TransactionInfo{
					Name:      cosmosAmount["denom"].(string),
					Symbol:    cosmosAmount["denom"].(string),
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount["amount"].(string),
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				receivedFieldsList = append(receivedFieldsList, &receivedFields)
				fromAddress = msg.ValidatorAddress
				toAddress = msg.DelegatorAddress
				value = cosmosAmount["amount"].(string)
			} else if strings.Contains(msg.Type, "MsgWithdrawDelegatorReward") {
				transactionType = "claimrewards"
				var denomLocal = ""
				var amountLocal = ""
				var log = txReceipt.TxResponse.Logs[counterMessages]
				for _, event := range log.Events {
					if event.Type == "coin_received" {
						for _, attr := range event.Attributes {
							if attr.Key == "amount" {
								amountLocal = string(attr.Value[0:strings.Index(attr.Value, "ukava")])
								denomLocal = "ukava"
								receivedFields := pb.TransactionInfo{
									Name:      denomLocal,
									Symbol:    denomLocal,
									TokenId:   "",
									Decimals:  walletInfo.Decimals,
									Value:     amountLocal,
									QuoteRate: 0,
									LogoUrl:   "",
									From:      msg.ValidatorAddress,
									To:        msg.DelegatorAddress,
								}
								receivedFieldsList = append(receivedFieldsList, &receivedFields)
								fromAddress = msg.ValidatorAddress
								toAddress = msg.DelegatorAddress
								value = amountLocal
							}
						}
					}
				}

			} else if strings.Contains(msg.Type, "MsgBeginRedelegate") {
				var cosmosAmount = msg.Amount.(map[string]interface{})
				transactionType = "delegate"
				sentFields := pb.TransactionInfo{
					Name:      cosmosAmount["denom"].(string),
					Symbol:    cosmosAmount["denom"].(string),
					TokenId:   "",
					Decimals:  walletInfo.Decimals,
					Value:     cosmosAmount["amount"].(string),
					QuoteRate: 0,
					LogoUrl:   "",
					From:      msg.DelegatorAddress,
					To:        msg.ValidatorAddress,
				}
				sentFieldsList = append(sentFieldsList, &sentFields)
				fromAddress = msg.DelegatorAddress
				toAddress = msg.ValidatorAddress
				value = cosmosAmount["amount"].(string)
			} else if strings.Contains(msg.Type, "MsgSend") {
				if request.Address == msg.FromAddress {
					var cosmosAmount = msg.Amount.([]CosmosAmount)
					transactionType = "send"
					var amountEnt = cosmosAmount[0]
					sentFields := pb.TransactionInfo{
						Name:      amountEnt.Denom,
						Symbol:    amountEnt.Denom,
						TokenId:   "",
						Decimals:  walletInfo.Decimals,
						Value:     amountEnt.Amount,
						QuoteRate: 0,
						LogoUrl:   "",
						From:      msg.DelegatorAddress,
						To:        msg.ValidatorAddress,
					}
					sentFieldsList = append(sentFieldsList, &sentFields)
					fromAddress = msg.DelegatorAddress
					toAddress = msg.ValidatorAddress
					value = amountEnt.Amount
				} else if request.Address == msg.FromAddress {
					var cosmosAmount = msg.Amount.([]CosmosAmount)
					transactionType = "send"
					var amountEnt = cosmosAmount[0]
					receivedFields := pb.TransactionInfo{
						Name:      amountEnt.Denom,
						Symbol:    amountEnt.Denom,
						TokenId:   "",
						Decimals:  walletInfo.Decimals,
						Value:     amountEnt.Amount,
						QuoteRate: 0,
						LogoUrl:   "",
						From:      msg.DelegatorAddress,
						To:        msg.ValidatorAddress,
					}
					receivedFieldsList = append(receivedFieldsList, &receivedFields)
					fromAddress = msg.DelegatorAddress
					toAddress = msg.ValidatorAddress
					value = amountEnt.Amount
				}
			} else if strings.Contains(msg.Type, "MsgRequestData") {
				transactionType = "request"
				fromAddress = request.Address
				toAddress = ""
				value = "0"
			} else if strings.Contains(msg.Type, "MsgExec") {
				transactionType = "report"
				fromAddress = request.Address
				toAddress = ""
				value = "0"
			}

		}
		var blockNumber, _ = strconv.ParseInt(item.Height, 0, 64)
		transactionData := pb.TransactionData{
			Id:                  Clean(item.Txhash).(string),
			From:                Clean(fromAddress).(string),
			To:                  Clean(toAddress).(string),
			Fee:                 "0",
			Date:                item.Timestamp.Unix(),
			Status:              statusTxn,
			Type:                transactionType,
			Block:               blockNumber,
			Value:               value,
			Nonce:               0,
			NativeTokenDecimals: walletInfo.Decimals,
			Description:         remarks,
			Sent:                sentFieldsList,
			Received:            receivedFieldsList,
			Others:              othersList,
		}
		transactionResponse.Transactions = append(transactionResponse.Transactions, &transactionData)

	}
	return &transactionResponse, nil
}
func cosmosTxReceipt(cosmos *Handler, txHash string, walletInfo config.CosmosWallets) *CosmosTxReceipt {
	url := fmt.Sprintf(walletInfo.REST+"/cosmos/tx/v1beta1/txs/%v", txHash)
	body, err := cosmos.httpRequest.GetRequest(url)
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
			endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[walletInfo.ChainName]
			if endpoint != "" {
				url = fmt.Sprintf(endpoint+"cosmos/tx/v1beta1/txs/%v", txHash)
				body, err = cosmos.httpRequest.GetRequest(url)
				if err != nil {
					cosmos.logger.Error("Error reading transaction receipt for hash: %v", txHash)
				}
				err = nil
			}
		}
		if err != nil {
			cosmos.logger.Error("Error reading transaction receipt for hash: %v", txHash)
		}
	}
	var cosmosTxReceipt *CosmosTxReceipt
	err = json.Unmarshal(body, &cosmosTxReceipt)
	if err != nil {
		cosmos.logger.Error("Error marshaling tx receipt structure : %v", txHash)
	}
	return cosmosTxReceipt
}
func Clean(arg interface{}) interface{} {
	if arg == nil {
		return ""
	} else {
		return arg
	}
}
func getRpcAndRestInfo(cosmos config.CosmosConfig, log *zap.SugaredLogger, httpRequest utils.IHttpRequest) (*ConfMap, error) {
	chainConfigData := make(map[string]config.CosmosWallets)
	rpcClientMap := make(map[string]jsonrpc.RPCClient)
	ibcTokenMap := make(map[string]IBCData)
	denomInfo := make(map[string]DenomInfo)
	IBCbody, err := httpRequest.GetRequest(cosmos.Cfg.IBCInfo) //Fetches the IBC Token Data
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(IBCbody, &ibcTokenMap)
	if err != nil {
		log.Errorf("Error in Unmarshalling Json: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	denomBody, err := httpRequest.GetRequest(cosmos.Cfg.DenomInfo) //Fetches the IBC Token Data
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(denomBody, &denomInfo)
	if err != nil {
		log.Errorf("Error in Unmarshalling Json: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	for _, chain := range cosmos.Cfg.Wallets {
		var cosmostationAsset CosmostationAssetInfo
		if chain.CosmostationAssetsInfo != "" {
			cosmoStationAssetBody, err := httpRequest.GetRequest(chain.CosmostationAssetsInfo) //Fetches the IBC Token Data
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(cosmoStationAssetBody, &cosmostationAsset)
			if err != nil {
				log.Errorf("Error in Unmarshalling Json: %v", err.Error())
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
		for _, assetInfo := range cosmostationAsset {
			if !strings.Contains(assetInfo.Denom, "ibc") {
				denomInfoKey := assetInfo.Denom + "__" + chain.ChainName
				if assetInfo.Description == "" {
					assetInfo.Description = fmt.Sprintf("The %s token of %s", assetInfo.Type, chain.ChainName)
				}
				if _, ok := denomInfo[denomInfoKey]; !ok {
					denomInfo[denomInfoKey] = DenomInfo{
						Chain:       assetInfo.OriginChain,
						Name:        assetInfo.BaseDenom,
						Denom:       assetInfo.Denom,
						Symbol:      assetInfo.DpDenom,
						Decimals:    assetInfo.Decimal,
						Description: assetInfo.Description,
						CoingeckoID: assetInfo.CoinGeckoID,
						Logos: struct {
							Png string `json:"png"`
						}{
							Png: cosmos.Cfg.CosmostationImageUrl + assetInfo.Image,
						},
					}
				}
			}
		}
		chainConfigData[chain.ChainName] = chain
		rpcClient := jsonrpc.NewClient(chain.RPC)
		rpcClientMap[chain.ChainName] = rpcClient
	}

	return &ConfMap{
		rpcClientMap:    rpcClientMap,
		chainConfigData: chainConfigData,
		ibcTokenInfo:    ibcTokenMap,
		denomInfo:       denomInfo,
	}, nil
}

func (h *Handler) getCosmosAccountInfo(address string, chain string) (CosmosAccountInfo, error) {
	var accountInfo CosmosAccountInfo
	chainConf := h.chainConfigData[chain]
	accountInfoEndpoint := fmt.Sprintf(chainConf.REST+"/cosmos/auth/v1beta1/accounts/%s", address)
	body, err := h.httpRequest.GetRequest(accountInfoEndpoint)
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
			endpoint := h.env.Cosmos.Cfg.BackUpUrls[chain]
			if endpoint != "" {
				accountInfoEndpoint = fmt.Sprintf(endpoint+"cosmos/auth/v1beta1/accounts/%s", address)
				body, err = h.httpRequest.GetRequest(accountInfoEndpoint)
				if err != nil {
					h.logger.Errorf("Error calling %s cosmos account info end point %s", chain, err.Error())
					return accountInfo, err
				}
				err = nil
			}
		}
		if err != nil {
			h.logger.Errorf("Error calling %s cosmos account info end point %s", chain, err.Error())
			return accountInfo, err
		}
	}
	err = json.Unmarshal(body, &accountInfo)
	if err != nil {
		h.logger.Errorf("Error in Unmarshalling Json: %v", err.Error())
		return accountInfo, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	return accountInfo, nil
}

func reverse(less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		return !less(i, j)
	}
}

// func (h *Handler) getBlueZelleAssets(in *pb.BalanceRequest, chainConf config.CosmosWallets) (*pb.CosmosAssetResponse, error) {
// 	balanceEndpoint := fmt.Sprintf(chainConf.REST+"/bank/balances/%s", in.Address)
// 	body, err := h.httpRequest.GetRequest(balanceEndpoint)
// 	if err != nil {
// 		h.logger.Errorf("Error calling %s balance end point %s", in.Chain, err.Error())
// 		return nil, err
// 	}
// 	var tokenBalance []*pb.CosmosTokenBalance
// 	var balanceRes BluzelleBalanceRes
// 	var quote, quoteRate float64
// 	// var quoteMarketData coingecko.QuoteData
// 	var accountInfo BluzelleAccountInfo
// 	err = json.Unmarshal(body, &balanceRes)
// 	if err != nil {
// 		h.logger.Errorf("Error in Unmarshalling Json: %v", err.Error())
// 		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 	}
// 	for _, balanceData := range balanceRes.Result {
// 		tokenPrice, err := h.coingecko.GetTokenExchangeByCoingeckoId("usd", blzInfo.CoingeckoId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		quoteRate = tokenPrice
// 		amount, err := strconv.ParseFloat(balanceData.Amount, 64)
// 		if err != nil {
// 			h.logger.Errorf("Error in Float Conversion %v", amount)
// 		}
// 		amount = amount / math.Pow10(int(chainConf.Decimals))
// 		quote = amount * quoteRate
// 		quoteMarketData, err = h.coingecko.GetCosmosTokenMarketInfo(blzInfo.CoingeckoId)
// 		if err != nil {
// 			h.logger.Errorf("Error in Fetching Market Info %v", err)
// 		}
// 		if quoteMarketData.QuoteRateChange24h == "" {
// 			quoteMarketData.QuoteRateChange24h = "0" //"0" is default string fpr this key (client's expectation)
// 		}
// 		tokenData := pb.CosmosTokenBalance{
// 			ContractName:         blzInfo.Description,
// 			ContractDecimals:     int32(chainConf.Decimals),
// 			ContractTickerSymbol: blzInfo.Symbol,
// 			Balance:              balanceData.Amount,
// 			LogoUrl:              blzInfo.LogoUrl,
// 			QuotePrice:           strconv.FormatFloat(quote, 'f', -1, 64),
// 			QuoteRate:            quoteRate,
// 			QuoteRate_24H:        quoteMarketData.QuoteRateChange24h,
// 			QuotePctChange_24H:   quoteMarketData.QuoteRatePctChange24h,
// 			Denom:                balanceData.Denom,
// 		}
// 		tokenBalance = append(tokenBalance, &tokenData)

// 	}
// 	less := func(i, j int) bool {
// 		return tokenBalance[i].QuotePrice < tokenBalance[j].QuotePrice
// 	}
// 	sort.Slice(tokenBalance, reverse(less))

// 	if len(tokenBalance) != 0 {
// 		accountInfoEndpoint := fmt.Sprintf(chainConf.REST+"/auth/accounts//%s", in.Address)
// 		accountBody, accErr := h.httpRequest.GetRequest(accountInfoEndpoint)
// 		if accErr != nil {
// 			h.logger.Errorf("Error calling %s cosmos account info end point %s", in.Chain, err.Error())
// 			return nil, err
// 		}
// 		err = json.Unmarshal(accountBody, &accountInfo)
// 		if err != nil {
// 			h.logger.Errorf("Error in Unmarshalling Json: %v", err.Error())
// 			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 		}
// 	} else {
// 		accountInfo.Result.Value.AccountNumber = "0"
// 		accountInfo.Result.Value.Sequence = "0"
// 	}
// 	var accountNum, sequence string
// 	accountNum = accountInfo.Result.Value.AccountNumber
// 	sequence = accountInfo.Result.Value.Sequence
// 	accountNumberInt, _ := strconv.Atoi(accountNum)
// 	sequenceInt, _ := strconv.Atoi(sequence)
// 	return &pb.CosmosAssetResponse{
// 		AccountNumber: int64(accountNumberInt),
// 		ChainId:       chainConf.ChainID,
// 		Sequence:      int64(sequenceInt),
// 		Token:         tokenBalance,
// 	}, nil
// }

// func (h *Handler) GetBalance(in *pb.BalanceRequest) (*pb.CosmosAssetResponse, error) {
// 	chainConf := h.chainConfigData[in.Chain]
// 	if in.Chain == "bluzelle" {
// 		return h.getBlueZelleAssets(in, chainConf)
// 	}
// 	balanceEndpoint := fmt.Sprintf(chainConf.REST+"/cosmos/bank/v1beta1/balances/%s", in.Address)
// 	body, err := h.httpRequest.GetRequest(balanceEndpoint)
// 	if err != nil {
// 		h.logger.Errorf("Error calling %s balance end point %s", in.Chain, err.Error())
// 		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
// 			endpoint := h.env.Cosmos.Cfg.BackUpUrls[in.Chain]
// 			if endpoint != "" {
// 				balanceEndpoint = fmt.Sprintf(endpoint+"cosmos/bank/v1beta1/balances/%s", in.Address)
// 				body, err = h.httpRequest.GetRequest(balanceEndpoint)
// 				if err != nil {
// 					h.logger.Errorf("Error calling %s balance end point %s", in.Chain, err.Error())
// 					return nil, err
// 				}
// 				err = nil
// 			}
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	var balanceRes CosmosBalanceRes
// 	var tokenBalance []*pb.CosmosTokenBalance
// 	var accountInfo CosmosAccountInfo
// 	err = json.Unmarshal(body, &balanceRes)
// 	if err != nil {
// 		h.logger.Errorf("Error in Unmarshalling Json: %v", err.Error())
// 		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
// 	}
// 	for _, balanceData := range balanceRes.Balances {
// 		var balanceDenom = balanceData.Denom
// 		var assetKey, contractName string
// 		var decimals int32
// 		var quote, quoteRate float64
// 		// var quoteMarketData coingecko.QuoteData
// 		if in.Chain == "cryptoorg" && balanceData.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 			balanceDenom = balanceData.Denom + "__cryptoorgchain"
// 		} else {
// 			balanceDenom = balanceData.Denom + "__" + in.Chain //To match it with key name of ibc data
// 		}

// 		if strings.Contains(balanceData.Denom, "ibc") {
// 			ibcInfo := h.ibcTokenInfo[balanceDenom]
// 			if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 				assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 				contractName = "IBC Token"
// 			}
// 		} else {
// 			assetKey = balanceData.Denom + "__" + in.Chain //To match it with key name of ibc data
// 		}
// 		denomInfo := h.denomInfo[assetKey]
// 		if denomValue, ok := h.denomInfo[assetKey]; ok && denomValue.CoingeckoID != "" {
// 			tokenPrice, err := h.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 			if err != nil {
// 				return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 			}
// 			quoteRate = tokenPrice
// 			amount, err := strconv.ParseFloat(balanceData.Amount, 64)
// 			if err != nil {
// 				h.logger.Errorf("Error in Float Conversion %v", amount)
// 			}
// 			amount = amount / math.Pow10(denomValue.Decimals)
// 			quote = amount * quoteRate
// 			quoteMarketData, err = h.coingecko.GetCosmosTokenMarketInfo(denomValue.CoingeckoID)
// 			if err != nil {
// 				h.logger.Errorf("Error in Fetching Market Info %v", err)
// 			}
// 		}
// 		if quoteMarketData.QuoteRateChange24h == "" {
// 			quoteMarketData.QuoteRateChange24h = "0" //"0" is default string fpr this key (client's expectation)
// 		}
// 		if contractName == "" { //For the Identified IBC Tokens
// 			contractName = denomInfo.Description
// 		}
// 		if contractName == "" { //For the UnIdentified IBC Tokens
// 			contractName = "Unknown"
// 			denomInfo.Symbol = "Unknown"
// 		}
// 		if denomInfo.Decimals == 0 {
// 			decimals = 6
// 		} else {
// 			decimals = int32(denomInfo.Decimals)
// 		}
// 		tokenData := pb.CosmosTokenBalance{
// 			ContractName:         contractName,
// 			ContractDecimals:     decimals,
// 			ContractTickerSymbol: denomInfo.Symbol,
// 			Balance:              balanceData.Amount,
// 			LogoUrl:              denomInfo.Logos.Png,
// 			QuotePrice:           strconv.FormatFloat(quote, 'f', -1, 64),
// 			QuoteRate:            quoteRate,
// 			QuoteRate_24H:        quoteMarketData.QuoteRateChange24h,
// 			QuotePctChange_24H:   quoteMarketData.QuoteRatePctChange24h,
// 			Denom:                balanceData.Denom,
// 		}
// 		tokenBalance = append(tokenBalance, &tokenData)
// 	}

// 	less := func(i, j int) bool {
// 		return tokenBalance[i].QuotePrice < tokenBalance[j].QuotePrice
// 	}
// 	sort.Slice(tokenBalance, reverse(less))

// 	if len(tokenBalance) != 0 {
// 		accountInfo, err = h.getCosmosAccountInfo(in.Address, in.Chain)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		accountInfo.Account.AccountNumber = "0"
// 		accountInfo.Account.Sequence = "0"
// 	}
// 	var accountNum, sequence string
// 	accountNum = accountInfo.Account.AccountNumber
// 	sequence = accountInfo.Account.Sequence
// 	if accountNum == "" {
// 		accountNum = accountInfo.Account.BaseAccount.AccountNumber
// 	}
// 	if sequence == "" {
// 		sequence = accountInfo.Account.BaseAccount.Sequence
// 	}
// 	accountNumberInt, _ := strconv.Atoi(accountNum)
// 	sequenceInt, _ := strconv.Atoi(sequence)
// 	return &pb.CosmosAssetResponse{
// 		AccountNumber: int64(accountNumberInt),
// 		ChainId:       chainConf.ChainID,
// 		Sequence:      int64(sequenceInt),
// 		Token:         tokenBalance,
// 	}, nil
// }

func newCosmosTxConfig() (client.TxConfig, codec.Codec, error) {
	var TestConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: "runtime",
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: "clientTest",
				}),
			},
		},
	})

	var (
		pcdc codec.ProtoCodecMarshaler
		cdc  codec.Codec
	)
	err := depinject.Inject(TestConfig, &pcdc, &cdc)
	if err != nil {
		return nil, nil, err
	}
	return authtx.NewTxConfig(pcdc, authtx.DefaultSignModes), cdc, nil
}

// GetCosmosCDPParams fetches the Kava CDP params
// TODO: Add support for 3P services
func (cosmos *Handler) GetCosmosCDPParams(request *pb.CosmosCDPParametersRequest) (*pb.CosmosCDPParametersResponse, error) {
	if strings.ToLower(request.Chain) == "kava" {
		walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)
		var respo pb.CosmosCDPParametersResponse
		var cosmosCDPParameters CosmosCDPParameterResp
		reqUrlCDPParameters := fmt.Sprintf(walletInfo.REST + "/cdp/parameters")
		bodyCDPParameters, errCDPParameters := cosmos.httpRequest.GetRequest(reqUrlCDPParameters)
		if errCDPParameters != nil {
			if strings.Contains(errCDPParameters.Error(), "429") || strings.Contains(errCDPParameters.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlCDPParameters = fmt.Sprintf(endpoint + "/cdp/parameters")
					bodyCDPParameters, errCDPParameters = cosmos.httpRequest.GetRequest(reqUrlCDPParameters)
					if errCDPParameters != nil {
						cosmos.logger.Error("Error fetching CDP Parameters. Err: ", errCDPParameters)
						return nil, errCDPParameters
					}
					errCDPParameters = nil
				}
			}
			if errCDPParameters != nil {
				cosmos.logger.Error("Error fetching CDP Parameters. Err: ", errCDPParameters)
				return nil, errCDPParameters
			}
		} else {
			errCDPParameters = json.Unmarshal(bodyCDPParameters, &cosmosCDPParameters)
			if errCDPParameters != nil {
				cosmos.logger.Error("Error Unmarshalling CDP Parameters. Err: ", errCDPParameters)
				return nil, status.Errorf(codes.Internal, errCDPParameters.Error(), "json unmarshalling error")
			}
		}
		for _, colParam := range cosmosCDPParameters.Result.CollateralParams {
			var colParamInfo pb.CollateralParams
			var debtLimit pb.DebtLimit

			colParamInfo.Denom = colParam.Denom
			colParamInfo.Type = colParam.Type
			colParamInfo.LiquidationRatio = colParam.LiquidationRatio
			colParamInfo.StabilityFee = colParam.StabilityFee
			colParamInfo.AuctionSize = colParam.AuctionSize
			colParamInfo.LiquidationPenalty = colParam.LiquidationPenalty
			colParamInfo.SpotMarketId = colParam.SpotMarketID
			colParamInfo.LiquidationMarketId = colParam.LiquidationMarketID
			colParamInfo.KeeperRewardPercentage = colParam.KeeperRewardPercentage
			colParamInfo.CheckCollateralizationIndexCount = colParam.CheckCollateralizationIndexCount
			colParamInfo.ConversionFactor = colParam.ConversionFactor
			debtLimit.Denom = colParam.DebtLimit.Denom
			debtLimit.Amount = colParam.DebtLimit.Amount
			colParamInfo.DebtLimit = &debtLimit
			respo.CollateralParams = append(respo.CollateralParams, &colParamInfo)
		}
		var debtParam pb.DebtParam
		debtParam.Denom = cosmosCDPParameters.Result.DebtParam.Denom
		debtParam.ReferenceAsset = cosmosCDPParameters.Result.DebtParam.ReferenceAsset
		debtParam.ConversionFactor = cosmosCDPParameters.Result.DebtParam.ConversionFactor
		debtParam.DebtFloor = cosmosCDPParameters.Result.DebtParam.DebtFloor
		respo.DebtParam = &debtParam

		var globalDebtLimit pb.GlobalDebtLimit
		globalDebtLimit.Denom = cosmosCDPParameters.Result.GlobalDebtLimit.Denom
		globalDebtLimit.Amount = cosmosCDPParameters.Result.GlobalDebtLimit.Amount
		respo.GlobalDebtLimit = &globalDebtLimit

		respo.SurplusAuctionThreshold = cosmosCDPParameters.Result.SurplusAuctionThreshold
		respo.SurplusAuctionLot = cosmosCDPParameters.Result.SurplusAuctionLot
		respo.DebtAuctionThreshold = cosmosCDPParameters.Result.DebtAuctionThreshold
		respo.DebtAuctionLot = cosmosCDPParameters.Result.DebtAuctionLot

		return &respo, nil
	} else {
		return nil, status.Errorf(codes.Internal, "Chain is not supported", "Chain is not supported")
	}
}

func (cosmos *Handler) BluzelleSendTx(request *pb.CosmosSendTxRequest, walletInfo config.CosmosWallets) (*pb.SendTransactionResponse, error) {
	sendTxUrl := walletInfo.REST + "/txs"
	var sendTxRes BluzelleSendTxResponse
	res, err := cosmos.httpRequest.PostRequest(sendTxUrl, bytes.NewBuffer(request.TxDetails))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &sendTxRes)
	if err != nil {
		return nil, err
	}
	if sendTxRes.Height == "0" && sendTxRes.RawLog != "" {
		return nil, status.Errorf(codes.InvalidArgument, sendTxRes.RawLog)
	} else if sendTxRes.Height == "0" {
		cosmos.logger.Errorf("error: %v", sendTxRes)
		return nil, status.Errorf(codes.InvalidArgument, string(res), "value of height is 0")
	}
	return &pb.SendTransactionResponse{
		TransactionId: sendTxRes.Txhash,
	}, nil
}

func (cosmos *Handler) SendTx(request *pb.CosmosSendTxRequest) (*pb.SendTransactionResponse, error) {
	walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)
	if request.Chain == "bluzelle" {
		return cosmos.BluzelleSendTx(request, walletInfo)
	}
	sendTxUrl := walletInfo.REST + "/cosmos/tx/v1beta1/txs"
	var sendTxRes CosmosSendTxRes
	sendTxBody := CosmosSendTxRequest{
		TxBytes: request.TxBytes,
		Mode:    request.Mode,
	}
	sendTxRawMsg, err := json.Marshal(sendTxBody)
	if err != nil {
		cosmos.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json marshalling error")
	}
	res, err := cosmos.httpRequest.PostRequest(sendTxUrl, bytes.NewBuffer(sendTxRawMsg))
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
			endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
			if endpoint != "" {
				sendTxUrl = endpoint + "/cosmos/tx/v1beta1/txs"
				res, err = cosmos.httpRequest.PostRequest(sendTxUrl, bytes.NewBuffer(sendTxRawMsg))
				if err != nil {
					return nil, err
				}
				err = nil
			}
		}
		if err != nil {
			return nil, err
		}
	}
	err = json.Unmarshal(res, &sendTxRes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	if sendTxRes.TxResponse.Height == "0" && sendTxRes.TxResponse.RawLog != "" {
		return nil, status.Errorf(codes.InvalidArgument, sendTxRes.TxResponse.RawLog, "Something went wrong")
	} else if sendTxRes.TxResponse.Height == "0" {
		cosmos.logger.Errorf("error: %v", sendTxRes.TxResponse)
		return nil, status.Errorf(codes.InvalidArgument, string(res), "value of tx response height is 0")
	}
	return &pb.SendTransactionResponse{
		TransactionId: sendTxRes.TxResponse.Txhash,
	}, nil
}

// GetDelegations Lists the Delegations for any wallet address
// TODO: Add support for 3P services
// func (cosmos *Handler) GetDelegations(request *pb.CosmosDelegationsRequest) (*pb.CosmosDelegationsResponse, error) {
// 	walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)
// 	// , _ = strconv.ParseInt(request.PageSize, 0, 64)
// 	var respo pb.CosmosDelegationsResponse
// 	var reqUrlValidators, reqUrlDelegations, reqUrlRewards, reqUrlUnboundDelegations, reqUrlValidatorsForDelegator, reqUrlPools, reqUrlInflation, reqUrlAnnualProv, reqUrlDistrParams string
// 	if request.Chain == "bluzelle" {
// 		var pageOffSet = 1
// 		var pageSizeForValidators = int64(1000)
// 		reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/staking/validators?status=bonded&page=%v&limit=%v", pageOffSet, pageSizeForValidators)
// 		reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/staking/delegators/%v/delegations", request.Address)
// 		reqUrlRewards = fmt.Sprintf(walletInfo.REST+"/distribution/delegators/%v/rewards", request.Address)
// 		reqUrlUnboundDelegations = fmt.Sprintf(walletInfo.REST+"/staking/delegators/%v/unbonding_delegations", request.Address)
// 		reqUrlValidatorsForDelegator = fmt.Sprintf(walletInfo.REST+"/staking/delegators/%v/validators", request.Address)
// 		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/minting/inflation")
// 		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/staking/pool")
// 		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/minting/annual-provisions")
// 		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/distribution/parameters")

// 		bodyValidators, errValidators := cosmos.httpRequest.GetRequest(reqUrlValidators)
// 		bodyDelegations, errDelegations := cosmos.httpRequest.GetRequest(reqUrlDelegations)
// 		bodyRewards, errRewards := cosmos.httpRequest.GetRequest(reqUrlRewards)
// 		bodyUnboundDelegations, errUnboundDelegations := cosmos.httpRequest.GetRequest(reqUrlUnboundDelegations)
// 		bodyValidatorsForDelegator, errDValidatorsForDelegator := cosmos.httpRequest.GetRequest(reqUrlValidatorsForDelegator)
// 		bodyInflation, errInflation := cosmos.httpRequest.GetRequest(reqUrlInflation)
// 		bodyPools, errPools := cosmos.httpRequest.GetRequest(reqUrlPools)
// 		bodyAnnualProv, errAnnualProv := cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
// 		bodyDistrParams, errDistrParams := cosmos.httpRequest.GetRequest(reqUrlDistrParams)

// 		var cosmosValidators CosmosValidatorsLaunchPad
// 		var cosmosDelegations CosmosDelegationsLaunchPad
// 		var cosmosRewards CosmosRewardsLaunchPad
// 		var cosmosUnboundDelegations AutoGeneratedLaunchPad
// 		var cosmosValidatorsForDelegator CosmosValidatorsForDelegatorLaunchPad
// 		var cosmosInflation CosmosInflationLaunchPad
// 		var cosmosPool CosmosPoolLaunchPad
// 		var cosmosAnnualProvision CosmosAnnualProvisionLaunchPad
// 		var cosmosDistributionParams CosmosDistributionParamsLaunchPad

// 		if errValidators != nil {
// 			cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
// 		} else {
// 			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
// 			if errValidators != nil {
// 				cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
// 				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errDelegations != nil {
// 			cosmos.logger.Error("Error fetching Delegations. Err: ", errDelegations)
// 		} else {
// 			errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
// 			if errDelegations != nil {
// 				cosmos.logger.Error("Error unmarshalling Delegations List. Err: ", errDelegations)
// 				return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errRewards != nil {
// 			cosmos.logger.Error("Error fetching Rewards. Err: ", errRewards)
// 		} else {
// 			errRewards = json.Unmarshal(bodyRewards, &cosmosRewards)
// 			if errRewards != nil {
// 				cosmos.logger.Error("Error unmarshalling Rewards List. Err: ", errRewards)
// 				return nil, status.Errorf(codes.Internal, errRewards.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errUnboundDelegations != nil {
// 			cosmos.logger.Error("Error fetching Inflation. Err: ", errUnboundDelegations)
// 		} else {
// 			errUnboundDelegations = json.Unmarshal(bodyUnboundDelegations, &cosmosUnboundDelegations)
// 			if errUnboundDelegations != nil {
// 				cosmos.logger.Error("Error unmarshalling Unbonding delegations List. Err: ", errUnboundDelegations)
// 				return nil, status.Errorf(codes.Internal, errUnboundDelegations.Error(), "json unmarshalling error")
// 			}
// 		}
// 		if errDValidatorsForDelegator != nil {
// 			cosmos.logger.Error("Error fetching Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 		} else {
// 			errDValidatorsForDelegator = json.Unmarshal(bodyValidatorsForDelegator, &cosmosValidatorsForDelegator)
// 			if errDValidatorsForDelegator != nil {
// 				cosmos.logger.Error("Error unmarshalling Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 				return nil, status.Errorf(codes.Internal, errDValidatorsForDelegator.Error(), "json unmarshalling error")
// 			}
// 		}
// 		if errInflation != nil {
// 			cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
// 			cosmosInflation.Inflation = "0"
// 		} else {
// 			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
// 			if cosmosInflation.Inflation == "" {
// 				cosmosInflation.Inflation = "0"
// 			}
// 			if errInflation != nil {
// 				cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
// 				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
// 			}
// 		}
// 		if errPools != nil {
// 			cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
// 		} else {
// 			errPools = json.Unmarshal(bodyPools, &cosmosPool)
// 			if errPools != nil {
// 				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errInflation)
// 				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errAnnualProv != nil {
// 			cosmos.logger.Error("Error fetching Annual Provision data. Err: ", errAnnualProv)
// 		} else {
// 			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
// 			if errAnnualProv != nil {
// 				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
// 				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errDistrParams != nil {
// 			cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams)
// 		} else {
// 			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
// 			if errDistrParams != nil {
// 				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
// 				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
// 			}
// 		}

// 		var totalStakedValue = 0.0
// 		var totalRewardsValue = 0.0
// 		var totalUnStakedValue = 0.0
// 		var totalStakedQuote = 0.0
// 		var totalRewardsQuote = 0.0
// 		var totalUnStakedQuote = 0.0
// 		respo.Delegations = make([]*pb.DelegationsInfo, 0)
// 		respo.UnboundDelegations = make([]*pb.UnDelegationsInfo, 0)

// 		var bondedToken = cosmosPool.Result.BondedTokens
// 		var quoteRate string
// 		var denomWiseData = make(map[string]pb.DenomsWiseValues)
// 		for _, delegation := range cosmosDelegations.Result {
// 			if delegation.Balance.Amount == "0" {
// 				continue
// 			}
// 			var delegationInfo pb.DelegationsInfo
// 			var rewardsInfo pb.RewardsInfo
// 			var delegatorInfo pb.DelegationDetail
// 			var balancerInfo pb.BalanceDetail
// 			delegatorInfo.DelegatorAddress = delegation.DelegatorAddress
// 			delegatorInfo.ValidatorAddress = delegation.ValidatorAddress
// 			delegatorInfo.Shares = delegation.Shares
// 			delegationInfo.Delegation = &delegatorInfo

// 			balancerInfo.Denom = delegation.Balance.Denom
// 			var balanceDenom = ""
// 			var assetKey string
// 			if request.Chain == "cryptoorg" && delegation.Balance.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 				balanceDenom = delegation.Balance.Denom + "__cryptoorgchain"
// 			} else {
// 				var chainName string
// 				if delegation.Balance.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				balanceDenom = delegation.Balance.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if strings.Contains(delegation.Balance.Denom, "ibc") {
// 				ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 				if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 					assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 				}
// 			} else {
// 				var chainName string
// 				if delegation.Balance.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				assetKey = delegation.Balance.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			var tokenPrice float64
// 			var tokenDecimals int
// 			if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 				if denomValue.CoingeckoID != "" {
// 					tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 				}
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.TickerSymbol = denomValue.Symbol
// 				balancerInfo.LogoUrl = denomValue.Logos.Png
// 				if denomValue.CoingeckoID != "" {
// 					balancerInfo.Decimals = int64(denomValue.Decimals)
// 					tokenDecimals = denomValue.Decimals
// 				} else {
// 					balancerInfo.Decimals = walletInfo.Decimals
// 					tokenDecimals = int(walletInfo.Decimals)
// 				}
// 			} else {
// 				tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.Decimals = walletInfo.Decimals
// 				tokenDecimals = int(walletInfo.Decimals)
// 			}
// 			totalStakedValue = totalStakedValue + cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)
// 			totalStakedQuote = totalStakedQuote + (cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 			balancerInfo.Amount = delegation.Balance.Amount
// 			balancerInfo.QuoteRate = quoteRate
// 			balancerInfo.Quote = strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 			delegationInfo.Balance = &balancerInfo
// 			if denomWiseDataExists, ok := denomWiseData[delegation.Balance.Denom]; ok {
// 				denomWiseData[delegation.Balance.Denom] = pb.DenomsWiseValues{
// 					StakeBalance:       strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.StakeBalance)+(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 					TotalStakedQuote:   strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalStakedQuote)+(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					UnstakeBalance:     denomWiseDataExists.UnstakeBalance,
// 					TotalUnstakedQuote: denomWiseDataExists.TotalUnstakedQuote,
// 					RewardsBalance:     denomWiseDataExists.RewardsBalance,
// 					TotalRewardsQuote:  denomWiseDataExists.TotalRewardsQuote,
// 					Denom:              delegation.Balance.Denom,
// 				}
// 			} else {
// 				denomWiseData[delegation.Balance.Denom] = pb.DenomsWiseValues{

// 					StakeBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals), 'f', -1, 64),
// 					TotalStakedQuote: strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:            delegation.Balance.Denom,
// 				}
// 			}

// 			for _, validator := range cosmosValidators.Result {
// 				if delegation.ValidatorAddress == validator.OperatorAddress {
// 					var info pb.ValidatorForDelegatorInfo
// 					var description pb.ValidatorDescription
// 					var validatorCommission pb.ValidatorCommission
// 					var validatorCommissionRates pb.ValidatorCommissionRates

// 					info.OperatorAddress = validator.OperatorAddress
// 					if len(validator.ConsensusPubkey) == 0 {
// 						info.ConsensusPubkey = ""
// 					} else {
// 						info.ConsensusPubkey = validator.ConsensusPubkey
// 					}
// 					info.Jailed = validator.Jailed
// 					info.Status = string(validator.Status)
// 					info.Tokens = validator.Tokens

// 					var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
// 					if !errParseTokens {
// 						cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
// 					}
// 					var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
// 					if !errParseBondedToken {
// 						cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
// 					}
// 					divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
// 					divResult = divResult.Mul(divResult, big.NewFloat(100))
// 					info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
// 					s := fmt.Sprintf("%.2f", divResult)
// 					info.VotingPower = s
// 					// Description Details
// 					description.Moniker = validator.Description.Moniker
// 					description.Identity = validator.Description.Identity
// 					description.Website = validator.Description.Website
// 					description.SecurityContact = validator.Description.SecurityContact
// 					description.Details = validator.Description.Details
// 					info.Description = &description

// 					info.UnbondingHeight = validator.UnbondingHeight
// 					info.UnbondingTime = validator.UnbondingTime.String()
// 					// Commission Details
// 					var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
// 					var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
// 					var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
// 					validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
// 					validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
// 					validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
// 					validatorCommission.CommissionRates = &validatorCommissionRates
// 					info.Commission = &validatorCommission
// 					info.Commission.UpdateTime = validator.Commission.UpdateTime.String()

// 					info.MinSelfDelegation = validator.MinSelfDelegation
// 					info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
// 					if request.Chain == "bluzelle" {
// 						info.Apr = 20
// 					} else if request.Chain == "terra" {
// 						commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 						if commissionRate == 0 {
// 							info.Apr = 12.9
// 						} else {
// 							commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 							commissionConsume := commissionRate * 12.9
// 							validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
// 							info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
// 						}
// 					} else {
// 						if validator.Status == 2 {
// 							if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
// 								var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Result.CommunityTax, cosmosPool.Result.BondedTokens, validator.Commission.CommissionRates.Rate)
// 								if apr < 0 {
// 									info.Apr = 0
// 								} else {
// 									info.Apr = Clean(apr * 100).(float64)
// 								}
// 							} else {
// 								info.Apr = 0
// 							}
// 						} else {
// 							info.Apr = 0
// 						}
// 					}
// 					info.Apr = math.Round(info.Apr*100) / 100
// 					delegationInfo.ValidatorDetails = &info
// 					break
// 				}
// 			}
// 			for _, reward := range cosmosRewards.Result.Rewards {
// 				if reward.ValidatorAddress == delegation.ValidatorAddress {
// 					rewardsInfo.ValidatorAddress = reward.ValidatorAddress
// 					for _, rewardEnt := range reward.Reward {
// 						var rewardLocal pb.RewardsListInfo
// 						var tokenPrice float64
// 						var assetKey, quoteRate string
// 						var balanceDenom = ""
// 						var tokenDecimals int
// 						if request.Chain == "cryptoorg" && rewardEnt.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 							balanceDenom = rewardEnt.Denom + "__cryptoorgchain"
// 						} else {
// 							var chainName string
// 							if rewardEnt.Denom != "uluna" && request.Chain == "terra2" {
// 								chainName = "terra"
// 							} else {
// 								chainName = request.Chain

// 							}
// 							balanceDenom = rewardEnt.Denom + "__" + chainName //To match it with key name of ibc data
// 						}
// 						if strings.Contains(rewardEnt.Denom, "ibc") {
// 							ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 							if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 								assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 							}
// 						} else {
// 							var chainName string
// 							if rewardEnt.Denom != "uluna" && request.Chain == "terra2" {
// 								chainName = "terra"
// 							} else {
// 								chainName = request.Chain

// 							}
// 							assetKey = rewardEnt.Denom + "__" + chainName //To match it with key name of ibc data
// 						}
// 						if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 							if denomValue.CoingeckoID != "" {
// 								tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 							}
// 							quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 							rewardLocal.TickerSymbol = denomValue.Symbol
// 							rewardLocal.LogoUrl = denomValue.Logos.Png
// 							if denomValue.CoingeckoID != "" {
// 								rewardLocal.Decimals = int64(denomValue.Decimals)
// 								tokenDecimals = denomValue.Decimals
// 							} else {
// 								rewardLocal.Decimals = walletInfo.Decimals
// 								tokenDecimals = int(walletInfo.Decimals)
// 							}
// 						} else {
// 							tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 							quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 							rewardLocal.Decimals = walletInfo.Decimals
// 							tokenDecimals = int(walletInfo.Decimals)
// 						}
// 						totalRewardsValue = totalRewardsValue + cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)
// 						totalRewardsQuote = totalRewardsQuote + (cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 						rewardLocal.Denom = rewardEnt.Denom
// 						rewardLocal.Amount = rewardEnt.Amount
// 						rewardLocal.Quote = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals)*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 						rewardLocal.QuoteRate = quoteRate
// 						rewardsInfo.Reward = append(rewardsInfo.Reward, &rewardLocal)
// 						if denomWiseDataExists, ok := denomWiseData[rewardEnt.Denom]; ok {
// 							denomWiseData[rewardEnt.Denom] = pb.DenomsWiseValues{
// 								UnstakeBalance:     denomWiseDataExists.UnstakeBalance,
// 								TotalUnstakedQuote: denomWiseDataExists.TotalUnstakedQuote,
// 								StakeBalance:       denomWiseDataExists.StakeBalance,
// 								TotalStakedQuote:   denomWiseDataExists.TotalStakedQuote,
// 								RewardsBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.RewardsBalance)+(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 								TotalRewardsQuote:  strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalRewardsQuote)+(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 								Denom:              rewardEnt.Denom,
// 							}
// 						} else {

// 							denomWiseData[rewardEnt.Denom] = pb.DenomsWiseValues{
// 								RewardsBalance:    strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals), 'f', -1, 64),
// 								TotalRewardsQuote: strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 								Denom:             rewardEnt.Denom,
// 							}
// 						}
// 					}
// 					break
// 				}
// 			}
// 			delegationInfo.RewardsDetails = &rewardsInfo
// 			respo.Delegations = append(respo.Delegations, &delegationInfo)
// 		}
// 		for _, unbondDelegation := range cosmosUnboundDelegations.Result {
// 			var totalBalance = 0.0
// 			var unbondDelegationInfo pb.UnDelegationsInfo
// 			unbondDelegationInfo.DelegatorAddress = unbondDelegation.DelegatorAddress
// 			unbondDelegationInfo.ValidatorAddress = unbondDelegation.ValidatorAddress
// 			var quoteRate string
// 			var entries pb.Entries
// 			totalBalance = totalBalance + cosmos.helper.ConvertStringToFloat64(unbondDelegation.Balance)
// 			entries.CreationHeight = ""
// 			entries.CompletionTime = ""
// 			entries.InitialBalance = unbondDelegation.InitialBalance
// 			entries.Balance = unbondDelegation.Balance
// 			unbondDelegationInfo.Entries = append(unbondDelegationInfo.Entries, &entries)

// 			var balancerInfo pb.BalanceDetail
// 			balancerInfo.Denom = walletInfo.Denom
// 			balancerInfo.Amount = strconv.FormatFloat(totalBalance, 'f', -1, 64)
// 			balancerInfo.Decimals = walletInfo.Decimals
// 			var tokenPrice float64
// 			var assetKey string
// 			var balanceDenom = ""
// 			var tokenDecimals int
// 			// Currently we are using the denom as the native token denom in unbonded delegations, once cosmos undelegate api response provides denom similer to delegations
// 			// we can replace the below walletInfo.Denom to the one recieved in the response. with that the non native token support can also be added. (from line 1101 to 1112)
// 			if request.Chain == "cryptoorg" && walletInfo.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 				balanceDenom = walletInfo.Denom + "__cryptoorgchain"
// 			} else {
// 				var chainName string
// 				if walletInfo.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				balanceDenom = walletInfo.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if strings.Contains(walletInfo.Denom, "ibc") {
// 				ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 				if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 					assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 				}
// 			} else {
// 				var chainName string
// 				if walletInfo.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				assetKey = walletInfo.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 				if denomValue.CoingeckoID != "" {
// 					tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 				}
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.TickerSymbol = denomValue.Symbol
// 				balancerInfo.LogoUrl = denomValue.Logos.Png
// 				if denomValue.CoingeckoID != "" {
// 					balancerInfo.Decimals = int64(denomValue.Decimals)
// 					tokenDecimals = denomValue.Decimals
// 				} else {
// 					balancerInfo.Decimals = walletInfo.Decimals
// 					tokenDecimals = int(walletInfo.Decimals)
// 				}
// 			} else {
// 				tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.Decimals = walletInfo.Decimals
// 				tokenDecimals = int(walletInfo.Decimals)
// 			}
// 			totalUnStakedValue = totalUnStakedValue + totalBalance
// 			totalUnStakedQuote = totalUnStakedQuote + (totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 			balancerInfo.QuoteRate = quoteRate
// 			balancerInfo.Quote = strconv.FormatFloat((totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 			if denomWiseDataExists, ok := denomWiseData[walletInfo.Denom]; ok {
// 				denomWiseData[walletInfo.Denom] = pb.DenomsWiseValues{
// 					RewardsBalance:     denomWiseDataExists.RewardsBalance,
// 					TotalRewardsQuote:  denomWiseDataExists.TotalRewardsQuote,
// 					StakeBalance:       denomWiseDataExists.StakeBalance,
// 					TotalStakedQuote:   denomWiseDataExists.TotalStakedQuote,
// 					UnstakeBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.UnstakeBalance)+(totalBalance/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 					TotalUnstakedQuote: strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalUnstakedQuote)+(totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:              walletInfo.Denom,
// 				}
// 			} else {
// 				denomWiseData[walletInfo.Denom] = pb.DenomsWiseValues{
// 					UnstakeBalance:     strconv.FormatFloat(totalBalance/math.Pow10(tokenDecimals), 'f', -1, 64),
// 					TotalUnstakedQuote: strconv.FormatFloat((totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:              walletInfo.Denom,
// 				}
// 			}
// 			unbondDelegationInfo.Balance = &balancerInfo
// 			// unbondDelegationInfo.QuoteRate = cosmos.getTokenQuoteRate(request.Chain)
// 			for _, validator := range cosmosValidators.Result {
// 				if unbondDelegation.ValidatorAddress == validator.OperatorAddress {
// 					var info pb.ValidatorForDelegatorInfo
// 					var description pb.ValidatorDescription
// 					var validatorCommission pb.ValidatorCommission
// 					var validatorCommissionRates pb.ValidatorCommissionRates

// 					info.OperatorAddress = validator.OperatorAddress
// 					if len(validator.ConsensusPubkey) == 0 {
// 						info.ConsensusPubkey = ""
// 					} else {
// 						info.ConsensusPubkey = validator.ConsensusPubkey
// 					}
// 					info.Jailed = validator.Jailed
// 					info.Status = string(validator.Status)
// 					info.Tokens = validator.Tokens

// 					var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
// 					if !errParseTokens {
// 						cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
// 					}
// 					var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
// 					if !errParseBondedToken {
// 						cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
// 					}
// 					divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
// 					divResult = divResult.Mul(divResult, big.NewFloat(100))
// 					info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
// 					s := fmt.Sprintf("%.2f", divResult)
// 					info.VotingPower = s
// 					// Description Details
// 					description.Moniker = validator.Description.Moniker
// 					description.Identity = validator.Description.Identity
// 					description.Website = validator.Description.Website
// 					description.SecurityContact = validator.Description.SecurityContact
// 					description.Details = validator.Description.Details
// 					info.Description = &description

// 					info.UnbondingHeight = validator.UnbondingHeight
// 					info.UnbondingTime = validator.UnbondingTime.String()
// 					// Commission Details
// 					var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
// 					var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
// 					var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
// 					validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
// 					validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
// 					validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
// 					validatorCommission.CommissionRates = &validatorCommissionRates
// 					info.Commission = &validatorCommission
// 					info.Commission.UpdateTime = validator.Commission.UpdateTime.String()
// 					info.MinSelfDelegation = validator.MinSelfDelegation
// 					info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
// 					if request.Chain == "bluzelle" {
// 						info.Apr = 20
// 					} else if request.Chain == "terra" {
// 						commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 						if commissionRate == 0 {
// 							info.Apr = 12.9
// 						} else {
// 							commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 							commissionConsume := commissionRate * 12.9
// 							validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
// 							info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
// 						}
// 					} else {
// 						if validator.Status == 2 {
// 							if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
// 								var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Result.CommunityTax, cosmosPool.Result.BondedTokens, validator.Commission.CommissionRates.Rate)
// 								if apr < 0 {
// 									info.Apr = 0
// 								} else {
// 									info.Apr = Clean(apr * 100).(float64)
// 								}
// 							} else {
// 								info.Apr = 0
// 							}
// 						} else {
// 							info.Apr = 0
// 						}
// 					}
// 					info.Apr = math.Round(info.Apr*100) / 100
// 					unbondDelegationInfo.ValidatorDetails = &info
// 					break
// 				}
// 			}
// 			respo.UnboundDelegations = append(respo.UnboundDelegations, &unbondDelegationInfo)
// 		}
// 		getCosmosAPRRequest := &pb.CosmosAprRatesRequest{
// 			Testnet: false,
// 			Chain:   request.Chain,
// 		}
// 		var cosmosAPRRates, _ = cosmos.GetCosmosAprRates(getCosmosAPRRequest)
// 		if cosmosAPRRates != nil {
// 			respo.Apr = cosmosAPRRates.Apr
// 		}
// 		var overallDenomsValuesInfo pb.OverallDenomsValues
// 		overallDenomsValuesInfo.TotalStakedQuote = strconv.FormatFloat(totalStakedQuote, 'f', -1, 64)
// 		overallDenomsValuesInfo.TotalUnstakedQuote = strconv.FormatFloat(totalUnStakedQuote, 'f', -1, 64)
// 		overallDenomsValuesInfo.TotalRewardsQuote = strconv.FormatFloat(totalRewardsQuote, 'f', -1, 64)
// 		respo.NetStakeValues = &overallDenomsValuesInfo

// 		var denomStakeInfo []*pb.DenomsWiseValues
// 		for key, value := range denomWiseData {
// 			data := pb.DenomsWiseValues{
// 				StakeBalance:       value.StakeBalance,
// 				UnstakeBalance:     value.UnstakeBalance,
// 				RewardsBalance:     value.RewardsBalance,
// 				TotalStakedQuote:   value.TotalStakedQuote,
// 				TotalUnstakedQuote: value.TotalUnstakedQuote,
// 				TotalRewardsQuote:  value.TotalRewardsQuote,
// 				Denom:              key,
// 			}
// 			denomStakeInfo = append(denomStakeInfo, &data)
// 		}
// 		respo.IndividualStakeValues = denomStakeInfo
// 	} else {
// 		var pageSize = int64(100)
// 		var pageOffSet = 0
// 		var pageSizeForValidators = int64(1000)
// 		if request.Chain == "terra2" {
// 			reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegations/%v?pagination.limit=%v&pagination.count_total=true", request.Address, pageSize)
// 			reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?pagination.limit=%v&pagination.count_total=true", pageSizeForValidators)
// 			reqUrlUnboundDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations?pagination.limit=%v&pagination.count_total=true", request.Address, pageSize)

// 		} else {
// 			reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegations/%v?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 			reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSizeForValidators)
// 			reqUrlUnboundDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 		}
// 		reqUrlRewards = fmt.Sprintf(walletInfo.REST+"/cosmos/distribution/v1beta1/delegators/%v/rewards", request.Address)
// 		reqUrlValidatorsForDelegator = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegators/%v/validators?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/inflation")
// 		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/cosmos/staking/v1beta1/pool")
// 		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/annual_provisions")
// 		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/cosmos/distribution/v1beta1/params")

// 		bodyValidators, errValidators := cosmos.httpRequest.GetRequest(reqUrlValidators)
// 		bodyDelegations, errDelegations := cosmos.httpRequest.GetRequest(reqUrlDelegations)
// 		bodyRewards, errRewards := cosmos.httpRequest.GetRequest(reqUrlRewards)
// 		bodyUnboundDelegations, errUnboundDelegations := cosmos.httpRequest.GetRequest(reqUrlUnboundDelegations)
// 		bodyValidatorsForDelegator, errDValidatorsForDelegator := cosmos.httpRequest.GetRequest(reqUrlValidatorsForDelegator)
// 		bodyInflation, errInflation := cosmos.httpRequest.GetRequest(reqUrlInflation)
// 		bodyPools, errPools := cosmos.httpRequest.GetRequest(reqUrlPools)
// 		bodyAnnualProv, errAnnualProv := cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
// 		bodyDistrParams, errDistrParams := cosmos.httpRequest.GetRequest(reqUrlDistrParams)

// 		var cosmosValidators CosmosValidators
// 		var cosmosDelegations CosmosDelegations
// 		var cosmosRewards CosmosRewards
// 		var cosmosUnboundDelegations CosmosUnboundDelegations
// 		var cosmosValidatorsForDelegator CosmosValidatorsForDelegator
// 		var cosmosInflation CosmosInflation
// 		var cosmosPool CosmosPool
// 		var cosmosAnnualProvision CosmosAnnualProvision
// 		var cosmosDistributionParams CosmosDistributionParams

// 		if errValidators != nil {
// 			if strings.Contains(errValidators.Error(), "429") || strings.Contains(errValidators.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSizeForValidators)
// 					if request.Chain == "terra2" {
// 						reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?pagination.limit=%v&pagination.count_total=true", pageSizeForValidators)
// 					}
// 					bodyValidators, errValidators = cosmos.httpRequest.GetRequest(reqUrlValidators)
// 					if errValidators != nil {
// 						cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
// 					}
// 					errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
// 					if errValidators != nil {
// 						cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
// 						return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
// 					}
// 					errValidators = nil
// 				}
// 			}
// 			if errValidators != nil {
// 				cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
// 			}
// 		} else {
// 			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
// 			if errValidators != nil {
// 				cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
// 				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errDelegations != nil {
// 			if strings.Contains(errDelegations.Error(), "429") || strings.Contains(errDelegations.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlDelegations = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegations/%v?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 					if request.Chain == "terra2" {
// 						reqUrlDelegations = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegations/%v?pagination.limit=%v&pagination.count_total=true", request.Address, pageSize)
// 					}
// 					bodyDelegations, errDelegations = cosmos.httpRequest.GetRequest(reqUrlDelegations)
// 					if errDelegations != nil {
// 						cosmos.logger.Error("Error fetching Validators. Err: ", errDelegations)
// 					}
// 					errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
// 					if errDelegations != nil {
// 						cosmos.logger.Error("Error unmarshalling Delegations List. Err: ", errDelegations)
// 						return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
// 					}
// 					errDelegations = nil
// 				}
// 			}
// 			if errDelegations != nil {
// 				cosmos.logger.Error("Error fetching Delegations. Err: ", errDelegations)
// 			}
// 		} else {
// 			errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
// 			if errDelegations != nil {
// 				cosmos.logger.Error("Error unmarshalling Delegations List. Err: ", errDelegations)
// 				return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errRewards != nil {
// 			if strings.Contains(errRewards.Error(), "429") || strings.Contains(errRewards.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlRewards = fmt.Sprintf(endpoint+"/cosmos/distribution/v1beta1/delegators/%v/rewards", request.Address)
// 					bodyRewards, errRewards = cosmos.httpRequest.GetRequest(reqUrlRewards)
// 					if errRewards != nil {
// 						cosmos.logger.Error("Error fetching Rewards. Err: ", errRewards)
// 					}
// 					errRewards = json.Unmarshal(bodyRewards, &cosmosRewards)
// 					if errRewards != nil {
// 						cosmos.logger.Error("Error unmarshalling Rewards List. Err: ", errRewards)
// 						return nil, status.Errorf(codes.Internal, errRewards.Error(), "json unmarshalling error")
// 					}
// 					errRewards = nil
// 				}
// 			}
// 			if errRewards != nil {
// 				cosmos.logger.Error("Error fetching Rewards. Err: ", errRewards)
// 			}
// 		} else {
// 			errRewards = json.Unmarshal(bodyRewards, &cosmosRewards)
// 			if errRewards != nil {
// 				cosmos.logger.Error("Error unmarshalling Rewards List. Err: ", errRewards)
// 				return nil, status.Errorf(codes.Internal, errRewards.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errUnboundDelegations != nil {
// 			if strings.Contains(errUnboundDelegations.Error(), "429") || strings.Contains(errUnboundDelegations.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlUnboundDelegations = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 					if request.Chain == "terra2" {
// 						reqUrlUnboundDelegations = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations?pagination.limit=%v&pagination.count_total=true", request.Address, pageSize)
// 					}
// 					bodyUnboundDelegations, errUnboundDelegations = cosmos.httpRequest.GetRequest(reqUrlUnboundDelegations)
// 					if errUnboundDelegations != nil {
// 						cosmos.logger.Error("Error fetching Inflation. Err: ", errUnboundDelegations)
// 					}
// 					errUnboundDelegations = json.Unmarshal(bodyUnboundDelegations, &cosmosUnboundDelegations)
// 					if errUnboundDelegations != nil {
// 						cosmos.logger.Error("Error unmarshalling Unbonding delegations List. Err: ", errUnboundDelegations)
// 						return nil, status.Errorf(codes.Internal, errUnboundDelegations.Error(), "json unmarshalling error")
// 					}
// 					errUnboundDelegations = nil
// 				}
// 			}
// 			if errUnboundDelegations != nil {
// 				cosmos.logger.Error("Error fetching Inflation. Err: ", errUnboundDelegations)
// 			}
// 		} else {
// 			errUnboundDelegations = json.Unmarshal(bodyUnboundDelegations, &cosmosUnboundDelegations)
// 			if errUnboundDelegations != nil {
// 				cosmos.logger.Error("Error unmarshalling Unbonding delegations List. Err: ", errUnboundDelegations)
// 				return nil, status.Errorf(codes.Internal, errUnboundDelegations.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errDValidatorsForDelegator != nil {
// 			if strings.Contains(errDValidatorsForDelegator.Error(), "429") || strings.Contains(errDValidatorsForDelegator.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlValidatorsForDelegator = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegators/%v/validators?pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", request.Address, pageOffSet, pageSize)
// 					bodyValidatorsForDelegator, errDValidatorsForDelegator = cosmos.httpRequest.GetRequest(reqUrlValidatorsForDelegator)
// 					if errDValidatorsForDelegator != nil {
// 						cosmos.logger.Error("Error fetching Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 					}
// 					errDValidatorsForDelegator = json.Unmarshal(bodyValidatorsForDelegator, &cosmosValidatorsForDelegator)
// 					if errDValidatorsForDelegator != nil {
// 						cosmos.logger.Error("Error unmarshalling Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 						return nil, status.Errorf(codes.Internal, errDValidatorsForDelegator.Error(), "json unmarshalling error")
// 					}
// 					errDValidatorsForDelegator = nil
// 				}
// 			}
// 			if errDValidatorsForDelegator != nil {
// 				cosmos.logger.Error("Error fetching Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 			}
// 		} else {
// 			errDValidatorsForDelegator = json.Unmarshal(bodyValidatorsForDelegator, &cosmosValidatorsForDelegator)
// 			if errDValidatorsForDelegator != nil {
// 				cosmos.logger.Error("Error unmarshalling Validator for the delegation. Err: ", errDValidatorsForDelegator)
// 				return nil, status.Errorf(codes.Internal, errDValidatorsForDelegator.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errInflation != nil {
// 			if strings.Contains(errInflation.Error(), "429") || strings.Contains(errInflation.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlInflation = fmt.Sprintf(endpoint + "/cosmos/mint/v1beta1/inflation")
// 					bodyInflation, errInflation = cosmos.httpRequest.GetRequest(reqUrlInflation)
// 					if errInflation != nil {
// 						cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
// 						cosmosInflation.Inflation = "0"
// 					}
// 					errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
// 					if cosmosInflation.Inflation == "" {
// 						cosmosInflation.Inflation = "0"
// 					}
// 					if errInflation != nil {
// 						cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
// 						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
// 					}
// 					errInflation = nil
// 				}
// 			}
// 			if errInflation != nil {
// 				cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
// 				cosmosInflation.Inflation = "0"
// 			}
// 		} else {
// 			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
// 			if cosmosInflation.Inflation == "" {
// 				cosmosInflation.Inflation = "0"
// 			}
// 			if errInflation != nil {
// 				cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
// 				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errPools != nil {
// 			if strings.Contains(errPools.Error(), "429") || strings.Contains(errPools.Error(), "No servers available") {
// 				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 				if endpoint != "" {
// 					reqUrlPools = fmt.Sprintf(endpoint + "/cosmos/staking/v1beta1/pool")
// 					bodyPools, errPools = cosmos.httpRequest.GetRequest(reqUrlPools)
// 					if errPools != nil {
// 						cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
// 					}
// 					errPools = json.Unmarshal(bodyPools, &cosmosPool)
// 					if errPools != nil {
// 						cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errInflation)
// 						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
// 					}
// 					errPools = nil
// 				}
// 			}
// 			if errPools != nil {
// 				cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
// 			}
// 		} else {
// 			errPools = json.Unmarshal(bodyPools, &cosmosPool)
// 			if errPools != nil {
// 				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errPools)
// 				return nil, status.Errorf(codes.Internal, errPools.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errAnnualProv != nil {
// 			cosmos.logger.Error("Error fetching Annual Provision data. Err: ", errAnnualProv)
// 		} else {
// 			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
// 			if errAnnualProv != nil {
// 				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
// 				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
// 			}
// 		}

// 		if errDistrParams != nil {
// 			cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams)
// 		} else {
// 			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
// 			if errDistrParams != nil {
// 				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
// 				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
// 			}
// 		}

// 		var totalStakedValue = 0.0
// 		var totalRewardsValue = 0.0
// 		var totalUnStakedValue = 0.0
// 		var totalStakedQuote = 0.0
// 		var totalRewardsQuote = 0.0
// 		var totalUnStakedQuote = 0.0
// 		respo.Delegations = make([]*pb.DelegationsInfo, 0)
// 		respo.UnboundDelegations = make([]*pb.UnDelegationsInfo, 0)

// 		var bondedToken = cosmosPool.Pool.BondedTokens
// 		var quoteRate string
// 		var denomWiseData = make(map[string]pb.DenomsWiseValues)
// 		for _, delegation := range cosmosDelegations.DelegationResponses {
// 			if delegation.Balance.Amount == "0" {
// 				continue
// 			}
// 			var delegationInfo pb.DelegationsInfo
// 			var rewardsInfo pb.RewardsInfo
// 			var delegatorInfo pb.DelegationDetail
// 			var balancerInfo pb.BalanceDetail
// 			delegatorInfo.DelegatorAddress = delegation.Delegation.DelegatorAddress
// 			delegatorInfo.ValidatorAddress = delegation.Delegation.ValidatorAddress
// 			delegatorInfo.Shares = delegation.Delegation.Shares
// 			delegationInfo.Delegation = &delegatorInfo

// 			balancerInfo.Denom = delegation.Balance.Denom
// 			var balanceDenom = ""
// 			var assetKey string
// 			if request.Chain == "cryptoorg" && delegation.Balance.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 				balanceDenom = delegation.Balance.Denom + "__cryptoorgchain"
// 			} else {
// 				var chainName string
// 				if delegation.Balance.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain
// 				}
// 				balanceDenom = delegation.Balance.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if strings.Contains(delegation.Balance.Denom, "ibc") {
// 				ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 				if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 					assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 				}
// 			} else {
// 				var chainName string
// 				if delegation.Balance.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				assetKey = delegation.Balance.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			var tokenPrice float64
// 			var tokenDecimals int
// 			if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 				if denomValue.CoingeckoID != "" {
// 					tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 				}
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.TickerSymbol = denomValue.Symbol
// 				balancerInfo.LogoUrl = denomValue.Logos.Png
// 				if denomValue.CoingeckoID != "" {
// 					balancerInfo.Decimals = int64(denomValue.Decimals)
// 					tokenDecimals = denomValue.Decimals
// 				} else {
// 					balancerInfo.Decimals = walletInfo.Decimals
// 					tokenDecimals = int(walletInfo.Decimals)
// 				}
// 			} else {
// 				tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.Decimals = walletInfo.Decimals
// 				tokenDecimals = int(walletInfo.Decimals)
// 			}
// 			totalStakedValue = totalStakedValue + cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)
// 			totalStakedQuote = totalStakedQuote + (cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 			balancerInfo.Amount = delegation.Balance.Amount
// 			balancerInfo.QuoteRate = quoteRate
// 			balancerInfo.Quote = strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 			delegationInfo.Balance = &balancerInfo
// 			if denomWiseDataExists, ok := denomWiseData[delegation.Balance.Denom]; ok {
// 				denomWiseData[delegation.Balance.Denom] = pb.DenomsWiseValues{
// 					StakeBalance:       strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.StakeBalance)+(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 					TotalStakedQuote:   strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalStakedQuote)+(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					UnstakeBalance:     denomWiseDataExists.UnstakeBalance,
// 					TotalUnstakedQuote: denomWiseDataExists.TotalUnstakedQuote,
// 					RewardsBalance:     denomWiseDataExists.RewardsBalance,
// 					TotalRewardsQuote:  denomWiseDataExists.TotalRewardsQuote,
// 					Denom:              delegation.Balance.Denom,
// 				}
// 			} else {
// 				denomWiseData[delegation.Balance.Denom] = pb.DenomsWiseValues{
// 					StakeBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals), 'f', -1, 64),
// 					TotalStakedQuote: strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(delegation.Balance.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:            delegation.Balance.Denom,
// 				}
// 			}

// 			for _, validator := range cosmosValidators.Validators {
// 				if delegation.Delegation.ValidatorAddress == validator.OperatorAddress {
// 					var info pb.ValidatorForDelegatorInfo
// 					var description pb.ValidatorDescription
// 					var validatorCommission pb.ValidatorCommission
// 					var validatorCommissionRates pb.ValidatorCommissionRates

// 					info.OperatorAddress = validator.OperatorAddress
// 					if len(validator.ConsensusPubkey.Key) == 0 {
// 						info.ConsensusPubkey = ""
// 					} else {
// 						info.ConsensusPubkey = validator.ConsensusPubkey.Key
// 					}
// 					info.Jailed = validator.Jailed
// 					info.Status = validator.Status
// 					info.Tokens = validator.Tokens

// 					var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
// 					if !errParseTokens {
// 						cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
// 					}
// 					var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
// 					if !errParseBondedToken {
// 						cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
// 					}
// 					divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
// 					divResult = divResult.Mul(divResult, big.NewFloat(100))
// 					info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
// 					s := fmt.Sprintf("%.2f", divResult)
// 					info.VotingPower = s
// 					// Description Details
// 					description.Moniker = validator.Description.Moniker
// 					description.Identity = validator.Description.Identity
// 					description.Website = validator.Description.Website
// 					description.SecurityContact = validator.Description.SecurityContact
// 					description.Details = validator.Description.Details
// 					info.Description = &description

// 					info.UnbondingHeight = validator.UnbondingHeight
// 					info.UnbondingTime = validator.UnbondingTime
// 					// Commission Details
// 					var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
// 					var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
// 					var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
// 					validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
// 					validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
// 					validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
// 					validatorCommission.CommissionRates = &validatorCommissionRates
// 					info.Commission = &validatorCommission
// 					info.Commission.UpdateTime = validator.Commission.UpdateTime

// 					info.MinSelfDelegation = validator.MinSelfDelegation
// 					info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
// 					if request.Chain == "bluzelle" {
// 						info.Apr = 20
// 					} else if request.Chain == "terra" {
// 						commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 						if commissionRate == 0 {
// 							info.Apr = 12.9
// 						} else {
// 							commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 							commissionConsume := commissionRate * 12.9
// 							validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
// 							info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
// 						}
// 					} else {
// 						if validator.Status == "BOND_STATUS_BONDED" {
// 							if len(cosmosAnnualProvision.AnnualProvisions) != 0 {

// 								var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Params.CommunityTax, cosmosPool.Pool.BondedTokens, validator.Commission.CommissionRates.Rate)
// 								if apr < 0 {
// 									info.Apr = 0
// 								} else {
// 									info.Apr = Clean(apr * 100).(float64)
// 								}
// 							} else {
// 								info.Apr = 0
// 							}
// 						} else {
// 							info.Apr = 0
// 						}
// 					}
// 					info.Apr = math.Round(info.Apr*100) / 100
// 					delegationInfo.ValidatorDetails = &info
// 					break
// 				}
// 			}
// 			for _, reward := range cosmosRewards.Rewards {
// 				if reward.ValidatorAddress == delegation.Delegation.ValidatorAddress {
// 					rewardsInfo.ValidatorAddress = reward.ValidatorAddress
// 					for _, rewardEnt := range reward.Reward {
// 						var rewardLocal pb.RewardsListInfo
// 						var tokenPrice float64
// 						var assetKey, quoteRate string
// 						var balanceDenom = ""
// 						var tokenDecimals int
// 						if request.Chain == "cryptoorg" && rewardEnt.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 							balanceDenom = rewardEnt.Denom + "__cryptoorgchain"
// 						} else {
// 							var chainName string
// 							if rewardEnt.Denom != "uluna" && request.Chain == "terra2" {
// 								chainName = "terra"
// 							} else {
// 								chainName = request.Chain

// 							}
// 							balanceDenom = rewardEnt.Denom + "__" + chainName //To match it with key name of ibc data
// 						}
// 						if strings.Contains(rewardEnt.Denom, "ibc") {
// 							ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 							if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 								assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 							}
// 						} else {
// 							var chainName string
// 							if rewardEnt.Denom != "uluna" && request.Chain == "terra2" {
// 								chainName = "terra"
// 							} else {
// 								chainName = request.Chain

// 							}
// 							assetKey = rewardEnt.Denom + "__" + chainName //To match it with key name of ibc data
// 						}
// 						if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 							if denomValue.CoingeckoID != "" {
// 								tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 							}
// 							quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 							rewardLocal.TickerSymbol = denomValue.Symbol
// 							rewardLocal.LogoUrl = denomValue.Logos.Png
// 							if denomValue.CoingeckoID != "" {
// 								rewardLocal.Decimals = int64(denomValue.Decimals)
// 								tokenDecimals = denomValue.Decimals
// 							} else {
// 								rewardLocal.Decimals = walletInfo.Decimals
// 								tokenDecimals = int(walletInfo.Decimals)
// 							}
// 						} else {
// 							tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 							quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 							rewardLocal.Decimals = walletInfo.Decimals
// 							tokenDecimals = int(walletInfo.Decimals)
// 						}
// 						totalRewardsValue = totalRewardsValue + cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)
// 						totalRewardsQuote = totalRewardsQuote + (cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 						rewardLocal.Denom = rewardEnt.Denom
// 						rewardLocal.Amount = rewardEnt.Amount
// 						rewardLocal.Quote = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals)*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 						rewardLocal.QuoteRate = quoteRate
// 						rewardsInfo.Reward = append(rewardsInfo.Reward, &rewardLocal)
// 						if denomWiseDataExists, ok := denomWiseData[rewardEnt.Denom]; ok {
// 							denomWiseData[rewardEnt.Denom] = pb.DenomsWiseValues{
// 								UnstakeBalance:     denomWiseDataExists.UnstakeBalance,
// 								TotalUnstakedQuote: denomWiseDataExists.TotalUnstakedQuote,
// 								StakeBalance:       denomWiseDataExists.StakeBalance,
// 								TotalStakedQuote:   denomWiseDataExists.TotalStakedQuote,
// 								RewardsBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.RewardsBalance)+(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 								TotalRewardsQuote:  strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalRewardsQuote)+(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 								Denom:              rewardEnt.Denom,
// 							}
// 						} else {

// 							denomWiseData[rewardEnt.Denom] = pb.DenomsWiseValues{
// 								RewardsBalance:    strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals), 'f', -1, 64),
// 								TotalRewardsQuote: strconv.FormatFloat((cosmos.helper.ConvertStringToFloat64(rewardEnt.Amount)/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 								Denom:             rewardEnt.Denom,
// 							}
// 						}
// 					}
// 					break
// 				}
// 			}
// 			delegationInfo.RewardsDetails = &rewardsInfo
// 			respo.Delegations = append(respo.Delegations, &delegationInfo)
// 		}
// 		for _, unbondDelegation := range cosmosUnboundDelegations.UnbondingResponses {
// 			var totalBalance = 0.0
// 			var delegatorInfo pb.DelegationDetail
// 			delegatorInfo.DelegatorAddress = unbondDelegation.DelegatorAddress
// 			delegatorInfo.ValidatorAddress = unbondDelegation.ValidatorAddress
// 			delegatorInfo.Shares = ""
// 			var unbondDelegationInfo pb.UnDelegationsInfo
// 			unbondDelegationInfo.DelegatorAddress = unbondDelegation.DelegatorAddress
// 			unbondDelegationInfo.ValidatorAddress = unbondDelegation.ValidatorAddress
// 			unbondDelegationInfo.Delegation = &delegatorInfo
// 			var quoteRate string
// 			for _, entry := range unbondDelegation.Entries {
// 				var entries pb.Entries
// 				totalBalance = totalBalance + cosmos.helper.ConvertStringToFloat64(entry.Balance)
// 				entries.CreationHeight = entry.CreationHeight
// 				entries.CompletionTime = entry.CompletionTime
// 				entries.InitialBalance = entry.InitialBalance
// 				entries.Balance = entry.Balance
// 				unbondDelegationInfo.Entries = append(unbondDelegationInfo.Entries, &entries)
// 			}
// 			var balancerInfo pb.BalanceDetail
// 			balancerInfo.Denom = walletInfo.Denom
// 			balancerInfo.Amount = strconv.FormatFloat(totalBalance, 'f', -1, 64)
// 			balancerInfo.Decimals = walletInfo.Decimals
// 			var tokenPrice float64
// 			var assetKey string
// 			var balanceDenom = ""
// 			var tokenDecimals int
// 			// Currently we are using the denom as the native token denom in unbonded delegations, once cosmos undelegate api response provides denom similer to delegations
// 			// we can replace the below walletInfo.Denom to the one recieved in the response. with that the non native token support can also be added. (from line 1101 to 1112)
// 			if request.Chain == "cryptoorg" && walletInfo.Denom == "basecro" { //basecro and uluna present in two chain networks, basecro handled,uluna need to check
// 				balanceDenom = walletInfo.Denom + "__cryptoorgchain"
// 			} else {
// 				var chainName string
// 				if walletInfo.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				balanceDenom = walletInfo.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if strings.Contains(walletInfo.Denom, "ibc") {
// 				ibcInfo := cosmos.ibcTokenInfo[balanceDenom]
// 				if ibcInfo.Origin.Chain != nil && ibcInfo.Origin.Denom != nil {
// 					assetKey = ibcInfo.Origin.Denom.(string) + "__" + ibcInfo.Origin.Chain.(string)
// 				}
// 			} else {
// 				var chainName string
// 				if walletInfo.Denom != "uluna" && request.Chain == "terra2" {
// 					chainName = "terra"
// 				} else {
// 					chainName = request.Chain

// 				}
// 				assetKey = walletInfo.Denom + "__" + chainName //To match it with key name of ibc data
// 			}
// 			if denomValue, ok := cosmos.denomInfo[assetKey]; ok {
// 				if denomValue.CoingeckoID != "" {
// 					tokenPrice, _ = cosmos.coingecko.GetTokenExchangeByCoingeckoId("usd", denomValue.CoingeckoID)
// 				}
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.TickerSymbol = denomValue.Symbol
// 				balancerInfo.LogoUrl = denomValue.Logos.Png
// 				if denomValue.CoingeckoID != "" {
// 					balancerInfo.Decimals = int64(denomValue.Decimals)
// 					tokenDecimals = denomValue.Decimals
// 				} else {
// 					balancerInfo.Decimals = walletInfo.Decimals
// 					tokenDecimals = int(walletInfo.Decimals)
// 				}
// 			} else {
// 				tokenPrice = cosmos.getTokenQuoteRateAsFloatVal(request.Chain)
// 				quoteRate = strconv.FormatFloat(tokenPrice, 'f', -1, 64)
// 				balancerInfo.Decimals = walletInfo.Decimals
// 				tokenDecimals = int(walletInfo.Decimals)
// 			}
// 			totalUnStakedValue = totalUnStakedValue + totalBalance
// 			totalUnStakedQuote = totalUnStakedQuote + (totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate)
// 			balancerInfo.QuoteRate = quoteRate
// 			balancerInfo.Quote = strconv.FormatFloat((totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64)
// 			if denomWiseDataExists, ok := denomWiseData[walletInfo.Denom]; ok {
// 				denomWiseData[walletInfo.Denom] = pb.DenomsWiseValues{
// 					RewardsBalance:     denomWiseDataExists.RewardsBalance,
// 					TotalRewardsQuote:  denomWiseDataExists.TotalRewardsQuote,
// 					StakeBalance:       denomWiseDataExists.StakeBalance,
// 					TotalStakedQuote:   denomWiseDataExists.TotalStakedQuote,
// 					UnstakeBalance:     strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.UnstakeBalance)+(totalBalance/math.Pow10(tokenDecimals)), 'f', -1, 64),
// 					TotalUnstakedQuote: strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(denomWiseDataExists.TotalUnstakedQuote)+(totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:              walletInfo.Denom,
// 				}
// 			} else {
// 				denomWiseData[walletInfo.Denom] = pb.DenomsWiseValues{
// 					UnstakeBalance:     strconv.FormatFloat(totalBalance/math.Pow10(tokenDecimals), 'f', -1, 64),
// 					TotalUnstakedQuote: strconv.FormatFloat((totalBalance/math.Pow10(tokenDecimals))*cosmos.helper.ConvertStringToFloat64(quoteRate), 'f', -1, 64),
// 					Denom:              walletInfo.Denom,
// 				}
// 			}
// 			unbondDelegationInfo.Balance = &balancerInfo
// 			// unbondDelegationInfo.QuoteRate = cosmos.getTokenQuoteRate(request.Chain)
// 			for _, validator := range cosmosValidators.Validators {
// 				if unbondDelegation.ValidatorAddress == validator.OperatorAddress {
// 					var info pb.ValidatorForDelegatorInfo
// 					var description pb.ValidatorDescription
// 					var validatorCommission pb.ValidatorCommission
// 					var validatorCommissionRates pb.ValidatorCommissionRates

// 					info.OperatorAddress = validator.OperatorAddress
// 					if len(validator.ConsensusPubkey.Key) == 0 {
// 						info.ConsensusPubkey = ""
// 					} else {
// 						info.ConsensusPubkey = validator.ConsensusPubkey.Key
// 					}
// 					info.Jailed = validator.Jailed
// 					info.Status = validator.Status
// 					info.Tokens = validator.Tokens

// 					var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
// 					if !errParseTokens {
// 						cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
// 					}
// 					var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
// 					if !errParseBondedToken {
// 						cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
// 					}
// 					divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
// 					divResult = divResult.Mul(divResult, big.NewFloat(100))
// 					info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
// 					s := fmt.Sprintf("%.2f", divResult)
// 					info.VotingPower = s
// 					// Description Details
// 					description.Moniker = validator.Description.Moniker
// 					description.Identity = validator.Description.Identity
// 					description.Website = validator.Description.Website
// 					description.SecurityContact = validator.Description.SecurityContact
// 					description.Details = validator.Description.Details
// 					info.Description = &description

// 					info.UnbondingHeight = validator.UnbondingHeight
// 					info.UnbondingTime = validator.UnbondingTime
// 					// Commission Details
// 					var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
// 					var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
// 					var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
// 					validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
// 					validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
// 					validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
// 					validatorCommission.CommissionRates = &validatorCommissionRates
// 					info.Commission = &validatorCommission
// 					info.Commission.UpdateTime = validator.Commission.UpdateTime
// 					info.MinSelfDelegation = validator.MinSelfDelegation
// 					info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
// 					if request.Chain == "bluzelle" {
// 						info.Apr = 20
// 					} else if request.Chain == "terra" {
// 						commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 						if commissionRate == 0 {
// 							info.Apr = 12.9
// 						} else {
// 							commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
// 							commissionConsume := commissionRate * 12.9
// 							validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
// 							info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
// 						}
// 					} else {
// 						if validator.Status == "BOND_STATUS_BONDED" {
// 							if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
// 								var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Params.CommunityTax, cosmosPool.Pool.BondedTokens, validator.Commission.CommissionRates.Rate)
// 								if apr < 0 {
// 									info.Apr = 0
// 								} else {
// 									info.Apr = Clean(apr * 100).(float64)
// 								}
// 							} else {
// 								info.Apr = 0
// 							}
// 						} else {
// 							info.Apr = 0
// 						}
// 					}
// 					info.Apr = math.Round(info.Apr*100) / 100
// 					unbondDelegationInfo.ValidatorDetails = &info
// 					break
// 				}
// 			}
// 			respo.UnboundDelegations = append(respo.UnboundDelegations, &unbondDelegationInfo)
// 		}
// 		getCosmosAPRRequest := &pb.CosmosAprRatesRequest{
// 			Testnet: false,
// 			Chain:   request.Chain,
// 		}
// 		var cosmosAPRRates, _ = cosmos.GetCosmosAprRates(getCosmosAPRRequest)
// 		if cosmosAPRRates != nil {
// 			respo.Apr = cosmosAPRRates.Apr
// 		}
// 		var overallDenomsValuesInfo pb.OverallDenomsValues
// 		overallDenomsValuesInfo.TotalStakedQuote = strconv.FormatFloat(totalStakedQuote, 'f', -1, 64)
// 		overallDenomsValuesInfo.TotalUnstakedQuote = strconv.FormatFloat(totalUnStakedQuote, 'f', -1, 64)
// 		overallDenomsValuesInfo.TotalRewardsQuote = strconv.FormatFloat(totalRewardsQuote, 'f', -1, 64)
// 		respo.NetStakeValues = &overallDenomsValuesInfo

// 		var denomStakeInfo []*pb.DenomsWiseValues
// 		for key, value := range denomWiseData {
// 			data := pb.DenomsWiseValues{
// 				StakeBalance:       value.StakeBalance,
// 				UnstakeBalance:     value.UnstakeBalance,
// 				RewardsBalance:     value.RewardsBalance,
// 				TotalStakedQuote:   value.TotalStakedQuote,
// 				TotalUnstakedQuote: value.TotalUnstakedQuote,
// 				TotalRewardsQuote:  value.TotalRewardsQuote,
// 				Denom:              key,
// 			}
// 			denomStakeInfo = append(denomStakeInfo, &data)
// 		}
// 		respo.IndividualStakeValues = denomStakeInfo
// 	}
// 	return &respo, nil
// }

// GetValidators Lists the Validators
// TODO: Add support for 3P services
func (cosmos *Handler) GetValidators(request *pb.CosmosValidatorsRequest) (*pb.CosmosValidatorsResponse, error) {
	walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)
	var reqUrlValidators, reqUrlPools, reqUrlInflation, reqUrlDelegations, reqUrlAnnualProv, reqUrlDistrParams string
	if request.Chain == "bluzelle" {
		var pageSize = int64(400)
		var pageOffSet = 1
		reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/staking/validators?status=bonded&page=%v&limit=%v", pageOffSet, pageSize)
		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/staking/pool")
		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/minting/inflation")
		reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/staking/delegators/%v/delegations", request.Address)
		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/minting/annual-provisions")
		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/distribution/parameters")
	} else if request.Chain == "terra2" {
		// this block is specifically maintained for terra2 since cosmos validators api doesnt support page offset for terra2 rest endpoint
		var pageSize = int64(400)
		reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.limit=%v&pagination.count_total=true", pageSize)
		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/cosmos/staking/v1beta1/pool")
		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/inflation")
		reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegations/%v", request.Address)
		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/annual_provisions")
		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/cosmos/distribution/v1beta1/params")
	} else {
		var pageSize = int64(400)
		var pageOffSet = 0
		reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSize)
		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/cosmos/staking/v1beta1/pool")
		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/inflation")
		reqUrlDelegations = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/delegations/%v", request.Address)
		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/annual_provisions")
		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/cosmos/distribution/v1beta1/params")
	}

	bodyValidators, errValidators := cosmos.httpRequest.GetRequest(reqUrlValidators)
	bodyPools, errPools := cosmos.httpRequest.GetRequest(reqUrlPools)
	bodyInflation, errInflation := cosmos.httpRequest.GetRequest(reqUrlInflation)
	bodyDelegations, errDelegations := cosmos.httpRequest.GetRequest(reqUrlDelegations)
	bodyAnnualProv, errAnnualProv := cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
	bodyDistrParams, errDistrParams := cosmos.httpRequest.GetRequest(reqUrlDistrParams)

	var respo pb.CosmosValidatorsResponse

	if request.Chain == "bluzelle" {
		var cosmosValidators CosmosValidatorsLaunchPad
		var cosmosPool CosmosPoolLaunchPad
		var cosmosInflation CosmosInflationLaunchPad
		var cosmosDelegations CosmosDelegationsLaunchPad
		var cosmosAnnualProvision CosmosAnnualProvisionLaunchPad
		var cosmosDistributionParams CosmosDistributionParamsLaunchPad

		if errValidators != nil {
			cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
		} else {
			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
			if errValidators != nil {
				cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
			}
		}
		if errPools != nil {
			cosmos.logger.Error("Error fetching Pool. Err: ", errPools)
		} else {
			errPools = json.Unmarshal(bodyPools, &cosmosPool)
			if errPools != nil {
				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errPools)
				return nil, status.Errorf(codes.Internal, errPools.Error(), "json unmarshalling error")
			}
		}
		if errInflation != nil {
			cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
			cosmosInflation.Inflation = "0"
		} else {
			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
			if cosmosInflation.Inflation == "" {
				cosmosInflation.Inflation = "0"
			}
			if errInflation != nil {
				cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}
		if errDelegations != nil {
			cosmos.logger.Error("Error fetching Delegations. Err: ", errDelegations)
		} else {
			errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
			if errDelegations != nil {
				cosmos.logger.Error("Error unmarshalling Delegations data. Err: ", errDelegations)
				return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
			}
		}

		if errAnnualProv != nil {
			cosmos.logger.Error("Error fetching Annual Provision data. Err: ", errAnnualProv)
		} else {
			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
			if errAnnualProv != nil {
				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
			}
		}

		if errDistrParams != nil {
			cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams)
		} else {
			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
			if errDistrParams != nil {
				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
			}
		}

		var validatorAddress []string
		if len(cosmosDelegations.Result) > 0 {
			for _, delegation := range cosmosDelegations.Result {
				validatorAddress = append(validatorAddress, delegation.ValidatorAddress)
			}
		}
		var bondedToken = cosmosPool.Result.BondedTokens

		respo.Height = "0"
		for _, validator := range cosmosValidators.Result {
			var info pb.ValidatorInfo
			var description pb.ValidatorDescription
			var validatorCommission pb.ValidatorCommission
			var validatorCommissionRates pb.ValidatorCommissionRates

			info.OperatorAddress = validator.OperatorAddress
			if len(validator.ConsensusPubkey) == 0 {
				info.ConsensusPubkey = ""
			} else {
				info.ConsensusPubkey = validator.ConsensusPubkey
			}
			info.Jailed = validator.Jailed
			info.Status = "BOND_STATUS_BONDED"
			info.Tokens = validator.Tokens

			var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
			if !errParseTokens {
				cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
			}
			var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
			if !errParseBondedToken {
				cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
			}
			divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
			divResult = divResult.Mul(divResult, big.NewFloat(100))
			info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
			s := fmt.Sprintf("%.2f", divResult)
			info.VotingPower = s
			// Description Details
			description.Moniker = validator.Description.Moniker
			description.Identity = validator.Description.Identity
			description.Website = validator.Description.Website
			description.SecurityContact = validator.Description.SecurityContact
			description.Details = validator.Description.Details
			info.Description = &description

			info.UnbondingHeight = validator.UnbondingHeight
			info.UnbondingTime = validator.UnbondingTime.String()
			if cosmos.deriveContainsString(validatorAddress, validator.OperatorAddress) {
				info.ActiveStake = true
			} else {
				info.ActiveStake = false
			}
			// Commission Details
			var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
			var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
			var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
			validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
			validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
			validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
			validatorCommission.CommissionRates = &validatorCommissionRates
			info.Commission = &validatorCommission
			info.Commission.UpdateTime = validator.Commission.UpdateTime.String()

			info.MinSelfDelegation = validator.MinSelfDelegation
			info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
			// info.QuoteRate = cosmos.getTokenQuoteRate(request.Chain)
			if request.Chain == "bluzelle" {
				info.Apr = 20
			} else if request.Chain == "terra" {
				commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
				if commissionRate == 0 {
					info.Apr = 12.9
				} else {
					commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
					commissionConsume := commissionRate * 12.9
					validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
					info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
				}
			} else {
				if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
					var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Result.CommunityTax, cosmosPool.Result.BondedTokens, validator.Commission.CommissionRates.Rate)
					if apr < 0 {
						info.Apr = 0
					} else {
					}
				} else {
					info.Apr = 0
				}
			}
			info.Apr = math.Round(info.Apr*100) / 100
			respo.Validators = append(respo.Validators, &info)
		}
	} else {
		var cosmosValidators CosmosValidators
		var cosmosPool CosmosPool
		var cosmosInflation CosmosInflation
		var cosmosDelegations CosmosDelegations
		var cosmosAnnualProvision CosmosAnnualProvision
		var cosmosDistributionParams CosmosDistributionParams

		if errValidators != nil {
			if strings.Contains(errValidators.Error(), "429") || strings.Contains(errValidators.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					var pageSize = int64(400)
					var pageOffSet = 0
					reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSize)
					if request.Chain == "terra2" {
						reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.limit=%v&pagination.count_total=true", pageSize)
					}
					bodyValidators, errValidators = cosmos.httpRequest.GetRequest(reqUrlValidators)
					if errValidators != nil {
						cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
					}
					errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
					if errValidators != nil {
						cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
						return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
					}
					errValidators = nil
				}
			}
			if errValidators != nil {
				cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
			}
		} else {
			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
			if errValidators != nil {
				cosmos.logger.Error("Error unmarshalling Validators List. Err: ", errValidators)
				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
			}
		}

		if errPools != nil {
			if strings.Contains(errPools.Error(), "429") || strings.Contains(errPools.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlPools = fmt.Sprintf(endpoint + "/cosmos/staking/v1beta1/pool")
					bodyPools, errPools = cosmos.httpRequest.GetRequest(reqUrlPools)
					if errPools != nil {
						cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
					}
					errPools = json.Unmarshal(bodyPools, &cosmosPool)
					if errPools != nil {
						cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errPools)
						return nil, status.Errorf(codes.Internal, errPools.Error(), "json unmarshalling error")
					}
					errPools = nil
				}
			}
			if errPools != nil {
				cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
			}
		} else {
			errPools = json.Unmarshal(bodyPools, &cosmosPool)
			if errPools != nil {
				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errPools)
				return nil, status.Errorf(codes.Internal, errPools.Error(), "json unmarshalling error")
			}
		}

		if errInflation != nil {
			if strings.Contains(errInflation.Error(), "429") || strings.Contains(errInflation.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlInflation = fmt.Sprintf(endpoint + "/cosmos/mint/v1beta1/inflation")
					bodyInflation, errInflation = cosmos.httpRequest.GetRequest(reqUrlInflation)
					if errInflation != nil {
						cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
						cosmosInflation.Inflation = "0"
					}
					errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
					if cosmosInflation.Inflation == "" {
						cosmosInflation.Inflation = "0"
					}
					if errInflation != nil {
						cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
					}
					errInflation = nil
				}
			}
			if errInflation != nil {
				cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
				cosmosInflation.Inflation = "0"
			}
		} else {
			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
			if cosmosInflation.Inflation == "" {
				cosmosInflation.Inflation = "0"
			}
			if errInflation != nil {
				cosmos.logger.Error("Error unmarshalling Inflation data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}

		if errDelegations != nil {
			if strings.Contains(errDelegations.Error(), "429") || strings.Contains(errDelegations.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlDelegations = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/delegations/%v", request.Address)
					bodyDelegations, errDelegations = cosmos.httpRequest.GetRequest(reqUrlDelegations)
					if errDelegations != nil {
						cosmos.logger.Error("Error fetching Validators. Err: ", errDelegations)
					}
					errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
					if errDelegations != nil {
						cosmos.logger.Error("Error unmarshalling Delegations List. Err: ", errDelegations)
						return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
					}
					errDelegations = nil
				}
			}
			if errDelegations != nil {
				cosmos.logger.Error("Error fetching Delegations. Err: ", errDelegations)
			}
		} else {
			errDelegations = json.Unmarshal(bodyDelegations, &cosmosDelegations)
			if errDelegations != nil {
				cosmos.logger.Error("Error unmarshalling Delegations data. Err: ", errDelegations)
				return nil, status.Errorf(codes.Internal, errDelegations.Error(), "json unmarshalling error")
			}
		}

		if errAnnualProv != nil {
			cosmos.logger.Error("Error fetching Annual Provision data. Err: ", errAnnualProv)
		} else {
			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
			if errAnnualProv != nil {
				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
			}
		}

		if errDistrParams != nil {
			cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams)
		} else {
			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
			if errDistrParams != nil {
				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
			}
		}

		var validatorAddress []string
		if len(cosmosDelegations.DelegationResponses) > 0 {
			for _, delegation := range cosmosDelegations.DelegationResponses {
				validatorAddress = append(validatorAddress, delegation.Delegation.ValidatorAddress)
			}
		}
		var bondedToken = cosmosPool.Pool.BondedTokens

		respo.Height = "0"
		for _, validator := range cosmosValidators.Validators {
			var info pb.ValidatorInfo
			var description pb.ValidatorDescription
			var validatorCommission pb.ValidatorCommission
			var validatorCommissionRates pb.ValidatorCommissionRates

			info.OperatorAddress = validator.OperatorAddress
			if len(validator.ConsensusPubkey.Key) == 0 {
				info.ConsensusPubkey = ""
			} else {
				info.ConsensusPubkey = validator.ConsensusPubkey.Key
			}
			info.Jailed = validator.Jailed
			info.Status = validator.Status
			info.Tokens = validator.Tokens

			var bignumTokens, errParseTokens = new(big.Float).SetString(validator.Tokens)
			if !errParseTokens {
				cosmos.logger.Error("Error converting Validator Tokens. Err: ", errInflation)
			}
			var bignumBondedToken, errParseBondedToken = new(big.Float).SetString(bondedToken)
			if !errParseBondedToken {
				cosmos.logger.Error("Error converting Bonded Token. Err: ", errInflation)
			}
			divResult := new(big.Float).Quo(bignumTokens, bignumBondedToken)
			divResult = divResult.Mul(divResult, big.NewFloat(100))
			info.DelegatorShares = strconv.FormatFloat(cosmos.helper.ConvertStringToFloat64(validator.DelegatorShares)/math.Pow10(int(walletInfo.Decimals)), 'f', -1, 64)
			s := fmt.Sprintf("%.2f", divResult)
			info.VotingPower = s
			// Description Details
			description.Moniker = validator.Description.Moniker
			description.Identity = validator.Description.Identity
			description.Website = validator.Description.Website
			description.SecurityContact = validator.Description.SecurityContact
			description.Details = validator.Description.Details
			info.Description = &description

			info.UnbondingHeight = validator.UnbondingHeight
			info.UnbondingTime = validator.UnbondingTime
			if cosmos.deriveContainsString(validatorAddress, validator.OperatorAddress) {
				info.ActiveStake = true
			} else {
				info.ActiveStake = false
			}
			// Commission Details
			var commRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.Rate, 64)
			var commMaxRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxRate, 64)
			var maxChangeRate, _ = strconv.ParseFloat(validator.Commission.CommissionRates.MaxChangeRate, 64)
			validatorCommissionRates.Rate = fmt.Sprintf("%.2f", commRate*100)
			validatorCommissionRates.MaxRate = fmt.Sprintf("%.2f", commMaxRate*100)
			validatorCommissionRates.MaxChangeRate = fmt.Sprintf("%.2f", maxChangeRate*100)
			validatorCommission.CommissionRates = &validatorCommissionRates
			info.Commission = &validatorCommission
			info.Commission.UpdateTime = validator.Commission.UpdateTime

			info.MinSelfDelegation = validator.MinSelfDelegation
			info.ImageUrl = cosmos.getImageURL(validator.Description.Identity)
			// info.QuoteRate = cosmos.getTokenQuoteRate(request.Chain)
			if request.Chain == "bluzelle" {
				info.Apr = 20
			} else if request.Chain == "terra" {
				commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
				if commissionRate == 0 {
					info.Apr = 12.9
				} else {
					commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
					commissionConsume := commissionRate * 12.9
					validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
					info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
				}
			} else {
				if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
					var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Params.CommunityTax, cosmosPool.Pool.BondedTokens, validator.Commission.CommissionRates.Rate)
					if apr < 0 {
						info.Apr = 0
					} else {
						info.Apr = Clean(apr * 100).(float64)
					}
				} else {
					info.Apr = 0
				}

			}
			info.Apr = math.Round(info.Apr*100) / 100
			respo.Validators = append(respo.Validators, &info)
		}
	}

	return &respo, nil

}

// GetCosmosAprRates fetches the Maximim apr amongst all validators
// TODO: Add support for 3P services
func (cosmos *Handler) GetCosmosAprRates(request *pb.CosmosAprRatesRequest) (*pb.CosmosAprRatesResponse, error) {
	walletInfo := cosmos.utils.GetCosmosWalletInfo(request.Chain)
	var respoAPR pb.CosmosAprRatesResponse
	var reqUrlValidators, reqUrlInflation, reqUrlPools, reqUrlAnnualProv, reqUrlDistrParams string
	if request.Chain == "bluzelle" {
		var pageSize = int64(400)
		var pageOffSet = 1
		reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/staking/validators?status=bonded&page=%v&limit=%v", pageOffSet, pageSize)
		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/staking/pool")
		reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/minting/annual-provisions")
		reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/distribution/parameters")

		var cosmosValidators *CosmosValidatorsLaunchPad
		bodyValidators, errValidators := cosmos.httpRequest.GetRequest(reqUrlValidators)
		if errValidators != nil {
			cosmos.logger.Error("Error fetching Validators. Err: ", errValidators)
		} else {
			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
			if errValidators != nil {
				cosmos.logger.Error("Error unmarshalling Validaotrs List. Err: ", errValidators)
				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
			}
		}
		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/inflation")
		bodyInflation, errInflation := cosmos.httpRequest.GetRequest(reqUrlInflation)
		bodyPools, errPools := cosmos.httpRequest.GetRequest(reqUrlPools)
		bodyAnnualProv, errAnnualProv := cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
		bodyDistrParams, errDistrParams := cosmos.httpRequest.GetRequest(reqUrlDistrParams)

		var cosmosInflation CosmosInflationLaunchPad
		var cosmosPool CosmosPoolLaunchPad
		var cosmosAnnualProvision CosmosAnnualProvisionLaunchPad
		var cosmosDistributionParams CosmosDistributionParamsLaunchPad

		if errInflation != nil {
			if strings.Contains(errInflation.Error(), "429") || strings.Contains(errInflation.Error(), "No servers available") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlInflation = fmt.Sprintf(endpoint + "/cosmos/mint/v1beta1/inflation")
					bodyInflation, errInflation = cosmos.httpRequest.GetRequest(reqUrlInflation)
					if errInflation != nil {
						cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
						cosmosInflation.Inflation = "0"
					}
					errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
					if cosmosInflation.Inflation == "" {
						cosmosInflation.Inflation = "0"
					}
					if errInflation != nil {
						cosmos.logger.Error("Error unmarshalling Inflation Data. Err: ", errInflation)
						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
					}
					errInflation = nil
				}
			}
			if errInflation != nil {
				cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation)
				cosmosInflation.Inflation = "0"
			}
		} else {
			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
			if cosmosInflation.Inflation == "" {
				cosmosInflation.Inflation = "0"
			}
			if errInflation != nil {
				cosmos.logger.Error("Error unmarshalling Inflation Data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}

		if errPools != nil {
			cosmos.logger.Error("Error fetching Pool data. Err: ", errPools)
		} else {
			errPools = json.Unmarshal(bodyPools, &cosmosPool)
			if errPools != nil {
				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}

		if errAnnualProv != nil {
			cosmos.logger.Error("Error fetching Annual Provision data. Err: ", errAnnualProv)
		} else {
			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
			if errAnnualProv != nil {
				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
			}
		}

		if errDistrParams != nil {
			cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams)
		} else {
			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
			if errDistrParams != nil {
				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
			}
		}

		var respo pb.CosmosValidatorsResponse
		for _, validator := range cosmosValidators.Result {
			var info pb.ValidatorInfo
			if request.Chain == "bluzelle" {
				info.Apr = 20
			} else if request.Chain == "terra" {
				commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
				if commissionRate == 0 {
					info.Apr = 12.9
				} else {
					commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
					commissionConsume := commissionRate * 12.9
					validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
					info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
				}
			} else {
				if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
					var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Result.CommunityTax, cosmosPool.Result.BondedTokens, validator.Commission.CommissionRates.Rate)
					if apr < 0 {
						info.Apr = 0
					} else {
						info.Apr = Clean(apr * 100).(float64)
					}
				} else {
					info.Apr = 0
				}
			}
			respo.Validators = append(respo.Validators, &info)
		}
		sort.Slice(respo.Validators, func(i, j int) bool {
			return respo.Validators[i].Apr < respo.Validators[j].Apr
		})
		var maxAprRate float64
		if respo.Validators != nil {
			maxAprRate = respo.Validators[len(respo.Validators)-1].Apr
		}
		respoAPR.Apr = maxAprRate
		respoAPR.Apr = math.Round(respoAPR.Apr*100) / 100
	} else {
		var pageSize = int64(400)
		var pageOffSet = 0
		if request.Chain == "terra2" {
			reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.limit=%v&pagination.count_total=true", pageSize)
			reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/annual_provisions")
			reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/cosmos/distribution/v1beta1/params")
		} else {
			reqUrlValidators = fmt.Sprintf(walletInfo.REST+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSize)
			reqUrlAnnualProv = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/annual_provisions")
			reqUrlDistrParams = fmt.Sprintf(walletInfo.REST + "/cosmos/distribution/v1beta1/params")
		}
		reqUrlPools = fmt.Sprintf(walletInfo.REST + "/cosmos/staking/v1beta1/pool")

		var cosmosValidators CosmosValidators
		var cosmosPool CosmosPool
		var cosmosAnnualProvision CosmosAnnualProvision
		var cosmosDistributionParams CosmosDistributionParams

		bodyValidators, errValidators := cosmos.httpRequest.GetRequest(reqUrlValidators)
		bodyPools, errPools := cosmos.httpRequest.GetRequest(reqUrlPools)
		bodyAnnualProv, errAnnualProv := cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
		bodyDistrParams, errDistrParams := cosmos.httpRequest.GetRequest(reqUrlDistrParams)

		if errValidators != nil {
			if strings.Contains(errValidators.Error(), "Too Many Requests") || strings.Contains(errValidators.Error(), "No servers available") ||
				strings.Contains(errValidators.Error(), "rate limit exceeded") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.offset=%v&pagination.limit=%v&pagination.count_total=true", pageOffSet, pageSize)
					if request.Chain == "terra2" {
						reqUrlValidators = fmt.Sprintf(endpoint+"/cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED&pagination.limit=%v&pagination.count_total=true", pageSize)
					}
					bodyValidators, errValidators = cosmos.httpRequest.GetRequest(reqUrlValidators)
					if errValidators != nil {
						cosmos.logger.Error("Error fetching Validators by using backup url Err: ", errValidators, reqUrlValidators)
					}
					errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
					if errValidators != nil {
						cosmos.logger.Error("Error unmarshalling Validaotrs List.by using backup url  Err: ", errValidators)
						return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
					}
					errValidators = nil
				}
			}
			if errValidators != nil {
				cosmos.logger.Error("Error fetching Validators. Err: ", errValidators, reqUrlValidators)
			}
		} else {
			errValidators = json.Unmarshal(bodyValidators, &cosmosValidators)
			if errValidators != nil {
				cosmos.logger.Error("Error unmarshalling Validaotrs List. Err: ", errValidators)
				return nil, status.Errorf(codes.Internal, errValidators.Error(), "json unmarshalling error")
			}
		}
		reqUrlInflation = fmt.Sprintf(walletInfo.REST + "/cosmos/mint/v1beta1/inflation")
		bodyInflation, errInflation := cosmos.httpRequest.GetRequest(reqUrlInflation)

		var cosmosInflation CosmosInflation
		if errInflation != nil {
			if strings.Contains(errInflation.Error(), "Too Many Requests") || strings.Contains(errInflation.Error(), "No servers available") ||
				strings.Contains(errInflation.Error(), "rate limit exceeded") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlInflation = fmt.Sprintf(endpoint + "/cosmos/mint/v1beta1/inflation")
					bodyInflation, errInflation = cosmos.httpRequest.GetRequest(reqUrlInflation)
					if errInflation != nil {
						cosmos.logger.Error("Error fetching Inflation by using backup url Err: ", errInflation, reqUrlInflation)
						cosmosInflation.Inflation = "0"
					}
					errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
					if cosmosInflation.Inflation == "" {
						cosmosInflation.Inflation = "0"
					}
					if errInflation != nil {
						cosmos.logger.Error("Error unmarshalling Inflation databy using backup url. Err: ", errInflation)
						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
					}
					errInflation = nil
				}
			}
			if errInflation != nil {
				cosmos.logger.Error("Error fetching Inflation. Err: ", errInflation, reqUrlInflation)
				cosmosInflation.Inflation = "0"
			}
		} else {
			errInflation = json.Unmarshal(bodyInflation, &cosmosInflation)
			if cosmosInflation.Inflation == "" {
				cosmosInflation.Inflation = "0"
			}
			if errInflation != nil {
				cosmos.logger.Error("Error unmarshalling Inflation Data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}

		if errPools != nil {
			if strings.Contains(errPools.Error(), "Too Many Requests") || strings.Contains(errPools.Error(), "No servers available") ||
				strings.Contains(errPools.Error(), "rate limit exceeded") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlPools = fmt.Sprintf(endpoint + "/cosmos/staking/v1beta1/pool")
					bodyPools, errPools = cosmos.httpRequest.GetRequest(reqUrlPools)
					if errPools != nil {
						cosmos.logger.Error("Error fetching Pool data by using backup url. Err: ", errPools, reqUrlPools)
					}
					errPools = json.Unmarshal(bodyPools, &cosmosPool)
					if errPools != nil {
						cosmos.logger.Error("Error unmarshalling Pool data by using backup url. Err: ", errInflation)
						return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
					}
					errPools = nil
				}
			}
			if errPools != nil {
				cosmos.logger.Error("Error fetching Pool data. Err: ", errPools, reqUrlPools)
			}
		} else {
			errPools = json.Unmarshal(bodyPools, &cosmosPool)
			if errPools != nil {
				cosmos.logger.Error("Error unmarshalling Pool data. Err: ", errInflation)
				return nil, status.Errorf(codes.Internal, errInflation.Error(), "json unmarshalling error")
			}
		}

		if errAnnualProv != nil {
			if strings.Contains(errAnnualProv.Error(), "Too Many Requests") || strings.Contains(errAnnualProv.Error(), "No servers available") ||
				strings.Contains(errAnnualProv.Error(), "rate limit exceeded") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlAnnualProv = fmt.Sprintf(endpoint + "/cosmos/mint/v1beta1/annual_provisions")
					bodyAnnualProv, errAnnualProv = cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
					if errAnnualProv != nil {
						cosmos.logger.Error("Error fetching Annual Provision data by using backup url. Err: ", errPools, reqUrlAnnualProv)
					}
					errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
					if errAnnualProv != nil {
						cosmos.logger.Error("Error unmarshalling  Annual Provision data by using backup url. Err: ", errAnnualProv)
						return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error"+string(bodyAnnualProv))
					}
					errAnnualProv = nil
				}
			}
			if errAnnualProv != nil {
				cosmos.logger.Error("Error fetching Annual Provision data . Err: ", errAnnualProv)
			}
		} else {
			errAnnualProv = json.Unmarshal(bodyAnnualProv, &cosmosAnnualProvision)
			if errAnnualProv != nil {
				cosmos.logger.Error("Error unmarshalling Annual Provision data. Err: ", errAnnualProv)
				return nil, status.Errorf(codes.Internal, errAnnualProv.Error(), "json unmarshalling error")
			}
		}

		if errDistrParams != nil {
			if strings.Contains(errDistrParams.Error(), "Too Many Requests") || strings.Contains(errDistrParams.Error(), "No servers available") ||
				strings.Contains(errDistrParams.Error(), "rate limit exceeded") {
				endpoint := cosmos.env.Cosmos.Cfg.BackUpUrls[request.Chain]
				if endpoint != "" {
					reqUrlDistrParams = fmt.Sprintf(endpoint + "/cosmos/distribution/v1beta1/params")
					bodyDistrParams, errDistrParams = cosmos.httpRequest.GetRequest(reqUrlAnnualProv)
					if errDistrParams != nil {
						cosmos.logger.Error("Error fetching Distribution Params data by using backup url. Err: ", errDistrParams, reqUrlDistrParams)
					}
					errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
					if errDistrParams != nil {
						cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
						return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
					}
					errDistrParams = nil
				}
			}
			if errDistrParams != nil {
				cosmos.logger.Error("Error fetching Distribution Params data. Err: ", errDistrParams, reqUrlDistrParams)
			}
		} else {
			errDistrParams = json.Unmarshal(bodyDistrParams, &cosmosDistributionParams)
			if errDistrParams != nil {
				cosmos.logger.Error("Error unmarshalling Distribution Params data. Err: ", errDistrParams)
				return nil, status.Errorf(codes.Internal, errDistrParams.Error(), "json unmarshalling error")
			}
		}

		var respo pb.CosmosValidatorsResponse
		for _, validator := range cosmosValidators.Validators {
			var info pb.ValidatorInfo
			if request.Chain == "bluzelle" {
				info.Apr = 20
			} else if request.Chain == "terra" {
				commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
				if commissionRate == 0 {
					info.Apr = 12.9
				} else {
					commissionRate := cosmos.helper.ConvertStringToFloat64(validator.Commission.CommissionRates.Rate)
					commissionConsume := commissionRate * 12.9
					validatorApr := strconv.FormatFloat(12.9-commissionConsume, 'f', -1, 64)
					info.Apr = math.Round(cosmos.helper.ConvertStringToFloat64(validatorApr)*100) / 100
				}
			} else {
				if len(cosmosAnnualProvision.AnnualProvisions) != 0 {
					var apr = cosmos.calculateRealAPR(cosmosAnnualProvision.AnnualProvisions, cosmosDistributionParams.Params.CommunityTax, cosmosPool.Pool.BondedTokens, validator.Commission.CommissionRates.Rate)
					if apr < 0 {
						info.Apr = 0
					} else {

						info.Apr = Clean(apr * 100).(float64)
					}
				} else {
					info.Apr = 0
				}
			}
			respo.Validators = append(respo.Validators, &info)
		}
		sort.Slice(respo.Validators, func(i, j int) bool {
			return respo.Validators[i].Apr < respo.Validators[j].Apr
		})
		var maxAprRate float64
		if respo.Validators != nil {
			maxAprRate = respo.Validators[len(respo.Validators)-1].Apr
		}
		respoAPR.Apr = maxAprRate
		respoAPR.Apr = math.Round(respoAPR.Apr*100) / 100
	}

	return &respoAPR, nil
}

func (h *Handler) deriveContainsString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func (h *Handler) getImageURL(identity string) string {
	if identity != "" {
		if strings.HasSuffix(identity, "\n") {
			identity = strings.TrimRight(identity, "\n")
		}
		identityURL := fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?key_suffix=%v&fields=pictures", identity)
		bodyValidators, errImageURL := h.httpRequest.GetRequest(identityURL)
		var validatorImageURL *ValidatorImageURL
		if errImageURL != nil {
			h.logger.Error("Error fetching Image URL api response. Err: ", errImageURL)
			return ""
		} else {
			errImageURL = json.Unmarshal(bodyValidators, &validatorImageURL)
			if errImageURL != nil {
				h.logger.Error("Error Unmarshalling Image URL api response. Err: ", errImageURL)
				return ""
			}
		}
		// if validatorImageURL != nil && validatorImageURL.Them != nil && len(validatorImageURL.Them) > 0 && validatorImageURL.Them[0].Pictures.Primary != nil && validatorImageURL.Them[0].Pictures.Primary.URL != "" {
		if validatorImageURL != nil && validatorImageURL.Them != nil && len(validatorImageURL.Them) > 0 && validatorImageURL.Them[0].Pictures.Primary.URL != "" {
			return validatorImageURL.Them[0].Pictures.Primary.URL
		} else {
			return ""
		}

	} else {
		return ""
	}
}

// func (h *Handler) getTokenQuoteRate(chain string) string {
// 	var quotePrice string
// 	quoteRate, err := h.coingecko.GetTokenExchange("usd", chain)
// 	if err != nil {
// 		h.logger.Errorf("Error for Exchange Quote Price for token  request  is : %v", err.Error())
// 		quotePrice = "0"
// 	} else {
// 		quotePrice = strconv.FormatFloat(quoteRate.Price, 'f', -1, 64)
// 	}
// 	return quotePrice
// }

/*
	 func (h *Handler) calculateAPR(inflation string, commision string, bondedTokenRatio string) float64 {
		// aprRate := (h.helper.ConvertStringToFloat64(inflation) * (1 - h.helper.ConvertStringToFloat64(commision))) / h.helper.ConvertStringToFloat64(bondedTokenRatio)
		var aprRate = h.helper.ConvertStringToFloat64(inflation) -
			(h.helper.ConvertStringToFloat64(inflation) *
				h.helper.ConvertStringToFloat64(
					commision))
		return aprRate
	}
*/
func (h *Handler) calculateRealAPR(annual_provision string, community_tax string, bonded_tokens string, val_commission string) float64 {
	aprRate := (h.helper.ConvertStringToFloat64(annual_provision) * (1 - h.helper.ConvertStringToFloat64(community_tax))) / h.helper.ConvertStringToFloat64(bonded_tokens)
	aprRate = aprRate - (aprRate * h.helper.ConvertStringToFloat64(val_commission))
	return aprRate
}
func (h *Handler) getCDPTokenAmountDetails(inflation string, commision string, bondedTokenRatio string) float64 {
	aprRate := (h.helper.ConvertStringToFloat64(inflation) * (1 - h.helper.ConvertStringToFloat64(commision))) / h.helper.ConvertStringToFloat64(bondedTokenRatio)
	return aprRate
}

// func (h *Handler) getTokenQuoteRateAsFloatVal(chain string) float64 {
// 	var quotePrice float64
// 	quoteRate, err := h.coingecko.GetTokenExchange("usd", chain)
// 	if err != nil {
// 		h.logger.Errorf("Error for Exchange Quote Price for token  request  is : %v", err.Error())
// 		quotePrice = 0.0
// 	} else {
// 		quotePrice = quoteRate.Price
// 	}
// 	return quotePrice
// }

// func (h *Handler) CosmosSimulateTx(request *pb.CosmosSimulateTxRequest) (*pb.CosmosSimulateTxResponse, error) {
// 	walletInfo := h.utils.GetCosmosWalletInfo(request.Chain)
// 	var baseFee float32
// 	if request.Chain == "bluzelle" {
// 		baseFee = 0.001
// 		return &pb.CosmosSimulateTxResponse{
// 			SimulationResult: true,
// 			GasLimit:         750000, //static
// 			GasPrice: &pb.CosmosSimulateGasPrice{
// 				Fast:        1.3 * baseFee,
// 				SafeLow:     1 * baseFee,
// 				Fastest:     1.6 * baseFee,
// 				Average:     1.3 * baseFee,
// 				SafeLowWait: 5,
// 				AvgWait:     2,
// 				FastWait:    1,
// 				FastestWait: 0.5,
// 			},
// 			Fee:     1500, //static
// 			Message: "Simulation executed successfully",
// 		}, nil
// 	}
// 	reqUrl := fmt.Sprintf(walletInfo.REST + "/cosmos/tx/v1beta1/simulate")
// 	var txBody SimulateTxBody
// 	var simulationStatus bool
// 	var message string
// 	txBody.TxBytes = request.TxBytes
// 	jsonReq, _ := json.Marshal(txBody)
// 	reqBody := bytes.NewBuffer(jsonReq)
// 	res, err := h.httpRequest.PostRequest(reqUrl, reqBody)
// 	if err != nil {
// 		var fee int64
// 		switch request.Chain {
// 		case "emoney":
// 			fee = 75000
// 		case "irisnet":
// 			fee = 150000
// 		case "osmosis":
// 			fee = 20000
// 		case "kava":
// 			fee = 2000
// 		case "akash":
// 			fee = 20000
// 		case "cryptoorgchain":
// 			fee = 18750
// 		default:
// 			fee = 1000
// 		}
// 		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
// 			endpoint := h.env.Cosmos.Cfg.BackUpUrls[request.Chain]
// 			if endpoint != "" {
// 				reqUrl = endpoint + "cosmos/tx/v1beta1/simulate"
// 				res, err = h.httpRequest.PostRequest(reqUrl, reqBody)
// 				if err != nil {
// 					h.logger.Error(err)
// 					baseFee = walletInfo.BaseFee
// 					quoteRate, _ := h.coingecko.GetTokenExchange("usd", request.Chain)
// 					return &pb.CosmosSimulateTxResponse{
// 						SimulationResult: false,
// 						GasLimit:         750000, //static
// 						GasPrice: &pb.CosmosSimulateGasPrice{
// 							Fast:        1.3 * baseFee,
// 							SafeLow:     1 * baseFee,
// 							Fastest:     1.6 * baseFee,
// 							Average:     1.3 * baseFee,
// 							SafeLowWait: 5,
// 							AvgWait:     2,
// 							FastWait:    1,
// 							FastestWait: 0.5,
// 						},
// 						Fee:       int64(fee),
// 						QuoteRate: strconv.FormatFloat(quoteRate.Price, 'f', -1, 64),
// 						Message:   err.Error(),
// 					}, nil
// 				}
// 				err = nil
// 			}
// 		}
// 		h.logger.Error(err)
// 		baseFee = walletInfo.BaseFee
// 		quoteRate, _ := h.coingecko.GetTokenExchange("usd", request.Chain)
// 		return &pb.CosmosSimulateTxResponse{
// 			SimulationResult: false,
// 			GasLimit:         750000, //static
// 			GasPrice: &pb.CosmosSimulateGasPrice{
// 				Fast:        1.3 * baseFee,
// 				SafeLow:     1 * baseFee,
// 				Fastest:     1.6 * baseFee,
// 				Average:     1.3 * baseFee,
// 				SafeLowWait: 5,
// 				AvgWait:     2,
// 				FastWait:    1,
// 				FastestWait: 0.5,
// 			},
// 			Fee:       fee,
// 			QuoteRate: strconv.FormatFloat(quoteRate.Price, 'f', -1, 64),
// 			Message:   err.Error(),
// 		}, nil
// 	}
// 	var simulateTxResponse SimulateTxResponse
// 	err = json.Unmarshal(res, &simulateTxResponse)
// 	if err != nil {
// 		var fee int64
// 		switch request.Chain {
// 		case "emoney":
// 			fee = 75000
// 		case "irisnet":
// 			fee = 150000
// 		case "osmosis":
// 			fee = 20000
// 		case "kava":
// 			fee = 2000
// 		case "akash":
// 			fee = 20000
// 		case "cryptoorgchain":
// 			fee = 18750
// 		default:
// 			fee = 1000
// 		}
// 		quoteRate, _ := h.coingecko.GetTokenExchange("usd", request.Chain)
// 		return &pb.CosmosSimulateTxResponse{
// 			SimulationResult: false,
// 			GasLimit:         750000, //static
// 			GasPrice: &pb.CosmosSimulateGasPrice{
// 				Fast:        1.3 * baseFee,
// 				SafeLow:     1 * baseFee,
// 				Fastest:     1.6 * baseFee,
// 				Average:     1.3 * baseFee,
// 				SafeLowWait: 5,
// 				AvgWait:     2,
// 				FastWait:    1,
// 				FastestWait: 0.5,
// 			},
// 			Fee:       fee,
// 			QuoteRate: strconv.FormatFloat(quoteRate.Price, 'f', -1, 64),
// 			Message:   err.Error(),
// 		}, nil
// 	}
// 	gasUsed, _ := strconv.ParseFloat(simulateTxResponse.GasInfo.GasUsed, 64)
// 	if gasUsed > 0 {
// 		simulationStatus = true
// 		message = "Simulation executed successfully"
// 	}
// 	baseFee = walletInfo.BaseFee
// 	gasUsed = math.Round(1.5 * gasUsed)
// 	gasPrice := 1.6 * baseFee //used the fastest gas_price factor
// 	fee := float64(gasPrice) * gasUsed
// 	quoteRate, err := h.coingecko.GetTokenExchange("usd", request.Chain)
// 	if err != nil {
// 		message = "Simulation executed successfully but quote rate failed due to" + err.Error()
// 	}
// 	return &pb.CosmosSimulateTxResponse{
// 		SimulationResult: simulationStatus,
// 		GasLimit:         int64(gasUsed),
// 		GasPrice: &pb.CosmosSimulateGasPrice{
// 			Fast:        1.3 * baseFee, // base value 0.025  multiplied by frequency factor
// 			SafeLow:     1 * baseFee,
// 			Fastest:     1.6 * baseFee,
// 			Average:     1.3 * baseFee,
// 			SafeLowWait: 5,
// 			AvgWait:     2,
// 			FastWait:    1,
// 			FastestWait: 0.5,
// 		},
// 		Fee:       int64(fee),
// 		QuoteRate: strconv.FormatFloat(quoteRate.Price, 'f', -1, 64),
// 		Message:   message,
// 	}, nil
// }

func (h *Handler) GetLatestBlockHeight(request *pb.GetCosmosBlockHeightRequest) (*pb.GetCosmosBlockHeightResponse, error) {
	walletInfo := h.utils.GetCosmosWalletInfo(request.Chain)
	reqUrl := fmt.Sprintf(walletInfo.REST + "/cosmos/base/tendermint/v1beta1/blocks/latest")
	res, err := h.httpRequest.GetRequest(reqUrl)
	if err != nil {
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "No servers available") {
			endpoint := h.env.Cosmos.Cfg.BackUpUrls[request.Chain]
			if endpoint != "" {
				reqUrl = fmt.Sprintf(endpoint + "/cosmos/base/tendermint/v1beta1/blocks/latest")
				res, err = h.httpRequest.GetRequest(reqUrl)
				if err != nil {
					h.logger.Error("Error in LatestBlock http request:", err)
					return nil, err
				}
				err = nil
			}
		}
		if err != nil {
			h.logger.Error("Error in LatestBlock http request:", err)
			return nil, err
		}
	}
	var blockRes BlockHeight
	err = json.Unmarshal(res, &blockRes)
	if err != nil {
		h.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	blockResInInt, _ := strconv.Atoi(blockRes.Block.Header.Height)
	return &pb.GetCosmosBlockHeightResponse{
		BlockHeight: int64(blockResInInt),
	}, nil
}

// type OpportunityChannelResponse struct {
// 	Opportunity models.OpportunityData
// 	Err         error
// }

// func (h *Handler) GetCosmosOpportunities(request *pb.GetOpportunitiesRequest) (*pb.GetOpportunitesResponse, error) {
// 	var res models.Opportunities
// 	channelOut := make(chan OpportunityChannelResponse)
// 	for _, chainInfo := range h.env.Cosmos.Cfg.Wallets {
// 		denomKey := chainInfo.Denom + "__" + chainInfo.ChainName
// 		chainInfo := chainInfo
// 		go func() {
// 			channelOut <- h.fetchOpportunities(denomKey, chainInfo)
// 		}()
// 	}
// 	for {
// 		select {
// 		case channelInfo := <-channelOut:
// 			if channelInfo.Err != nil {
// 				return nil, status.Errorf(codes.Internal, channelInfo.Err.Error(), "Error in Fetching Cosmos Opportunities")
// 			}
// 			if channelInfo.Opportunity.Chain == request.Chain {
// 				res.Current = append(res.Current, channelInfo.Opportunity)
// 			} else {
// 				res.Others = append(res.Others, channelInfo.Opportunity)
// 			}
// 			if len(h.env.Cosmos.Cfg.Wallets) == len(res.Others)+len(res.Current) {
// 				marshal, err := json.Marshal(res)
// 				if err != nil {
// 					return nil, err
// 				}
// 				return &pb.GetOpportunitesResponse{
// 					Opportunities: marshal,
// 				}, nil
// 			}
// 		}
// 	}
// }

// func (h *Handler) fetchOpportunities(denomKey string, chainInfo config.CosmosWallets) OpportunityChannelResponse {
// 	var res OpportunityChannelResponse
// 	maxAprRequestInfo := &pb.CosmosAprRatesRequest{
// 		Chain: chainInfo.ChainName,
// 	}
// 	maxAprRate, err := h.GetCosmosAprRates(maxAprRequestInfo)
// 	if err != nil {
// 		res.Err = err
// 		return res
// 	}
// 	if chainInfo.ChainName == "bluzelle" {
// 		res.Opportunity.Apr = fmt.Sprintf("%.2f", maxAprRate.Apr)
// 		res.Opportunity.Apr = fmt.Sprintf("%.2f", maxAprRate.Apr)
// 		res.Opportunity.Chain = chainInfo.ChainName
// 		res.Opportunity.Logo = "https://api.frontierwallet.com/images/chain/bluzelle.svg"
// 		res.Opportunity.StakingType = "Network"
// 		res.Opportunity.TokenName = "BLZ"
// 		res.Opportunity.ProtocolName = "Bluzelle"
// 		res.Opportunity.CoolDownPeriod = "0"
// 		res.Opportunity.MinLockup = "0"
// 		res.Opportunity.RewardSchedule = "1x"
// 	} else {
// 		denomInfo := h.denomInfo[denomKey]
// 		res.Opportunity.Apr = fmt.Sprintf("%.2f", maxAprRate.Apr)
// 		res.Opportunity.Apr = fmt.Sprintf("%.2f", maxAprRate.Apr)
// 		res.Opportunity.Chain = denomInfo.Chain
// 		res.Opportunity.Logo = denomInfo.Logos.Png
// 		res.Opportunity.StakingType = "Network"
// 		res.Opportunity.TokenName = denomInfo.Symbol
// 		res.Opportunity.ProtocolName = denomInfo.Name
// 		res.Opportunity.CoolDownPeriod = "0"
// 		res.Opportunity.MinLockup = "0"
// 		res.Opportunity.RewardSchedule = "1x"
// 	}
// 	return res
// }
