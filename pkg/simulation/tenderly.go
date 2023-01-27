package simulation

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	web32 "github.com/chenzhijie/go-web3"
	_ "github.com/ethereum/go-ethereum/eth/tracers/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Tenderly struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
	helpers     utils.Helpers
	utils       utils.UtilConf
}

func NewTenderlyService(env *config.Config, logger *zap.SugaredLogger) *Tenderly {
	helpers := utils.Helpers{}
	httRequest := utils.NewHttpRequest(logger)
	utils := utils.NewUtils(logger, env)
	return &Tenderly{
		env:         env,
		logger:      logger,
		httpRequest: httRequest,
		helpers:     helpers,
		utils:       *utils,
	}
}

func (t *Tenderly) SimulateTx(request SimulateTxRequest) (*SimulateTxResponse, error) {
	simulateResUrl := fmt.Sprintf("https://api.tenderly.co/api/v1/account/%s/project/%s/simulate", t.env.Tenderly.UserName, t.env.Tenderly.Project)
	t.logger.Infof(simulateResUrl)
	if val, ok := TenderlyChainNetworkId[request.Chain]; ok {
		sellAmountParamBigInt, _ := new(big.Int).SetString(request.Value, 10)
		gas, _ := strconv.Atoi(request.Gas)
		if gas == 0 {
			gas = 8000000 //max limit
		}
		txBody := TenderlyTxBody{
			NetworkId:      val,
			From:           request.From,
			To:             request.To,
			Input:          strings.Join(request.InputData, ""),
			Gas:            gas,
			GasPrice:       request.GasPrice,
			Value:          sellAmountParamBigInt,
			SaveIfFails:    true,
			Save:           false,
			SimulationType: "quick",
		}
		marshalTxBody, err := json.Marshal(txBody)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Unmarshalling error")
		}
		txBodyString := string(marshalTxBody)
		t.logger.Infof(txBodyString)
		res, err := t.httpRequest.PostRequestWithHeaders(simulateResUrl, txBodyString, "X-Access-Key", t.env.Tenderly.AccessKey)
		if err != nil {
			var tenderlyError TenderlyTxErrorResponse
			err = json.Unmarshal(res, &tenderlyError)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
			}
			return nil, status.Errorf(codes.Internal, tenderlyError.Error.Message, "Transaction Simulation Failed")
		}
		var simulationTxRes TenderlyTxResponse
		err = json.Unmarshal(res, &simulationTxRes)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
		}

		var simulationResp SimulateTxResponse
		simulationResp.ProviderId = "3"
		if simulationTxRes.Transaction.ErrorMessage != "" {
			simulationResp.ChainId = simulationTxRes.Transaction.NetworkID
			simulationResp.SimulationResult = simulationTxRes.Transaction.Status
			if simulationTxRes.Transaction.Method == "" {
				simulationResp.Type = ""
			}
			simulationResp.Action = "NONE"
			simulationResp.ActionMessage = simulationTxRes.Transaction.ErrorMessage
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
			simulationResp.GasUsed = strconv.Itoa(simulationTxRes.Transaction.GasUsed)
			simulationResp.HumanReadableForm = "Transaction Simulation Failed"
		} else if simulationTxRes.Transaction.TransactionInfo.Logs == nil {
			chain, _ := strconv.Atoi(val)
			info := t.utils.GetWalletInfo(strconv.Itoa(chain))
			simulationResp.ChainId = simulationTxRes.Transaction.NetworkID
			simulationResp.SimulationResult = simulationTxRes.Transaction.Status
			if simulationTxRes.Transaction.Method == "" {
				simulationResp.Type = "Asset Transfer"
			}
			simulationResp.Action = "WARN"
			simulationResp.ActionMessage = "You are transferring ER20 tokens directly to their own token contract. In most cases this will lead to you losing them forever."
			simulationResp.ActionMessage = simulationTxRes.Transaction.ErrorMessage
			simulationResp.GasUsed = strconv.Itoa(simulationTxRes.Transaction.GasUsed)
			nativeTokenAddress := t.GetFrontierSpecificNativeToken(chain)
			amountValue := simulationTxRes.Transaction.Value
			assetsSentTokenName := info.NativeTokenInfo.Name
			assestsSentSymbol := info.NativeTokenInfo.Symbol
			assestsSentDecimals := info.NativeTokenInfo.Decimals
			amount := t.HexToDecimals(amountValue)
			decAmount, _ := t.GetDecimals(amount, assestsSentDecimals)
			humanReadableForm := fmt.Sprintf("Send %v %v", decAmount, assestsSentSymbol)
			simulationResp.HumanReadableForm = humanReadableForm
			simulationResp.SimulationData.AssetsSent = []AssetsSent{
				{
					ContractAddress: nativeTokenAddress,
					Name:            assetsSentTokenName,
					Symbol:          assestsSentSymbol,
					Decimals:        assestsSentDecimals,
				},
			}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		} else if simulationTxRes.Transaction.To == request.To && simulationTxRes.Transaction.TransactionInfo.Logs[0].Raw.Data == approveHexValue {
			toTokenAddress := simulationTxRes.Transaction.TransactionInfo.CallTrace.To
			assetsSentTokenName, assestsSentSymbol, assestsSentDecimals, _ := t.GetContractData(toTokenAddress, val)
			amount := t.HexToDecimals(simulationTxRes.Transaction.TransactionInfo.Logs[0].Raw.Data)
			decAmount, _ := t.GetDecimals(amount, assestsSentDecimals)
			humanReadableForm := fmt.Sprintf("Approve to transfer up to %v %v", decAmount, assestsSentSymbol)
			simulationResp.ChainId = simulationTxRes.Transaction.NetworkID
			simulationResp.SimulationResult = simulationTxRes.Transaction.Status
			if simulationTxRes.Transaction.Method == "" {
				simulationResp.Type = "Token Approval"
			}
			simulationResp.Action = "NONE"
			simulationResp.GasUsed = strconv.Itoa(simulationTxRes.Transaction.GasUsed)
			simulationResp.HumanReadableForm = humanReadableForm
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{
				{
					ContractAddress: toTokenAddress,
					Name:            assetsSentTokenName,
					Symbol:          assestsSentSymbol,
					Decimals:        assestsSentDecimals,
				},
			}
		} else if simulationTxRes.Transaction.TransactionInfo.Logs != nil {
			simulationResp.ChainId = simulationTxRes.Transaction.NetworkID
			simulationResp.SimulationResult = simulationTxRes.Transaction.Status
			if simulationTxRes.Transaction.Method == "" {
				simulationResp.Type = "Asset Transfer"
			}
			simulationResp.Action = "WARN"
			simulationResp.ActionMessage = "You are transferring ER20 tokens directly to their own token contract. In most cases this will lead to you losing them forever."
			simulationResp.ActionMessage = simulationTxRes.Transaction.ErrorMessage
			simulationResp.GasUsed = strconv.Itoa(simulationTxRes.Transaction.GasUsed)
			toTokenAddress := simulationTxRes.Transaction.TransactionInfo.Logs[0].Raw.Address
			amountValue := simulationTxRes.Transaction.TransactionInfo.Logs[0].Raw.Data
			assetsSentTokenName, assestsSentSymbol, assestsSentDecimals, _ := t.GetContractData(toTokenAddress, val)
			amount := t.HexToDecimals(amountValue)
			decAmount, _ := t.GetDecimals(amount, assestsSentDecimals)
			humanReadableForm := fmt.Sprintf("Send %v %v", decAmount, assestsSentSymbol)
			simulationResp.HumanReadableForm = humanReadableForm
			simulationResp.SimulationData.AssetsSent = []AssetsSent{
				{
					ContractAddress: toTokenAddress,
					Name:            assetsSentTokenName,
					Symbol:          assestsSentSymbol,
					Decimals:        assestsSentDecimals,
				},
			}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		} else {
			simulationResp.ChainId = simulationTxRes.Transaction.NetworkID
			simulationResp.SimulationResult = simulationTxRes.Transaction.Status
			if simulationTxRes.Transaction.Method == "" {
				simulationResp.Type = "Swap"
			}
			simulationResp.Action = "NONE"
			simulationResp.ActionMessage = simulationTxRes.Transaction.ErrorMessage
			simulationResp.GasUsed = strconv.Itoa(simulationTxRes.Transaction.GasUsed)
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		}
		return &simulationResp, nil
	} else {
		return nil, status.Errorf(codes.Unavailable, "Chain not supported for tenderly simulation")
	}
}

func (t *Tenderly) GetRPC(chain string) (string, error) {
	var rpc string
	var err error
	for _, item := range t.env.EVM.Cfg.Wallets {
		if chain == item.ChainName || chain == item.NativeTokenInfo.ChainId {
			rpc = item.RPC
			return rpc, err
		}
	}
	return rpc, err
}

func (t *Tenderly) GetContractData(address string, chain string) (string, string, string, error) {
	rpc, err := t.GetRPC(chain)
	web3, err := web32.NewWeb3(rpc)
	client, err := web3.Eth.NewContract(tokenABI, address)

	tName, err := client.Call(name)
	if err != nil {
		t.logger.Error(err)
		return "", "", "", err
	}

	tSymbol, err := client.Call(symbol)
	if err != nil {
		t.logger.Error(err)
		return "", "", "", err
	}

	tDecimals, err := client.Call(decimals)
	if err != nil {
		t.logger.Error(err)
		return "", "", "", err
	}

	tokenName := fmt.Sprintf("%v", tName)
	tokenSymbol := fmt.Sprintf("%v", tSymbol)
	tokenDecimals := fmt.Sprintf("%v", tDecimals)

	return tokenName, tokenSymbol, tokenDecimals, nil
}

func (t *Tenderly) WeiToEther(value string, decimals int64) (string, error) {
	newValue, _ := strconv.ParseInt(value, 10, 64)
	toPow := math.Pow(float64(10), float64(decimals))
	toEther := float64(newValue) / toPow
	return fmt.Sprintf("%v", toEther), nil
}

func (t *Tenderly) StringToInt64(value string) int64 {
	intValue, _ := strconv.ParseInt(value, 10, 64)
	return intValue
}

func (t *Tenderly) HexToDecimals(value string) string {
	valueBigint := new(big.Int)
	valueBigint.SetString(value[2:], 16)
	return valueBigint.String()
}

func (t *Tenderly) GetDecimals(value string, decimals string) (string, error) {
	floatValue, _ := strconv.ParseFloat(value, 64)
	floatDecimals, _ := strconv.ParseFloat(decimals, 64)
	powValue := math.Pow(10, floatDecimals)

	res := floatValue / powValue
	stringRes := fmt.Sprintf("%f", res)
	return stringRes, nil
}

func (t *Tenderly) GetFrontierSpecificNativeToken(chain int) string {
	var tokenAddress string
	info := t.utils.GetWalletInfo(strconv.Itoa(chain))
	if info.ChainID == chain {
		tokenAddress = info.NativeTokenInfo.Address
	}
	return tokenAddress
}
