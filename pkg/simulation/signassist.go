package simulation

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"strconv"
	"strings"
)

type SignAssist struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
	utils       *utils.UtilConf
}

func NewSignAssist(env *config.Config, logger *zap.SugaredLogger) *SignAssist {
	httpRequest := utils.NewHttpRequest(logger)
	utils := utils.NewUtils(logger, env)
	return &SignAssist{env: env, logger: logger, httpRequest: httpRequest, utils: utils}
}

func (s *SignAssist) SimulateTx(request SimulateTxRequest) (*SimulateTxResponse, error) {
	signAssistURL := fmt.Sprintf(s.env.SignAssist.Endpoint)
	s.logger.Infof(signAssistURL)
	if val, ok := SignAssistSupportedChains[request.Chain]; ok {
		intVal, err := strconv.Atoi(val)
		txPayload := SignAssistTxBody{
			NetworkId: intVal,
			TransactionParams: []TransactionParams{
				{
					From:     request.From,
					To:       request.To,
					Value:    request.Value,
					Gas:      request.Gas,
					GasPrice: request.GasPrice,
					Data:     strings.Join(request.InputData, ""),
				},
			},
		}
		marshalTxBody, err := json.Marshal(txPayload)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "Marshalling error")
		}
		txBodyString := string(marshalTxBody)
		body, err := s.httpRequest.PostRequestWithHeaders(signAssistURL, txBodyString, "X-API-KEY", s.env.SignAssist.AccessKey)
		if err != nil {
			var signassistError SignAssistError
			err = json.Unmarshal(body, &signassistError)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
			}
			return nil, status.Errorf(codes.Internal, signassistError.Error, "Transaction Simulation Failed")
		}
		var signAssistResp SignAssistResponse
		var simulationResp SimulateTxResponse
		err = json.Unmarshal(body, &signAssistResp)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
		}
		if signAssistResp.Simulation.TransactionInfo.SimulationResult.Status == true {
			intGasUsed := strconv.Itoa(signAssistResp.Simulation.TransactionInfo.GasUsed)
			simulationResp.ChainId = request.Chain
			simulationResp.ProviderId = "2"
			simulationResp.SimulationResult = signAssistResp.Simulation.TransactionInfo.SimulationResult.Status
			simulationResp.Type = signAssistResp.Simulation.TransactionInfo.Method
			simulationResp.Action = "NONE"
			simulationResp.ActionMessage = "Successful"
			simulationResp.GasUsed = intGasUsed
			if signAssistResp.Simulation.TransactionInfo.Method == "fallback" && len(signAssistResp.Simulation.AssetsSent) != 0 && len(signAssistResp.Simulation.AssetsReceived) != 0 && len(signAssistResp.Simulation.AssetsApprovals) != 0 {
				sentAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsSent[0].Symbol, signAssistResp.Simulation.AssetsSent[0].ContractAddress)
				receivedAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsReceived[0].Symbol, signAssistResp.Simulation.AssetsReceived[0].ContractAddress)
				approvedAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsApprovals[0].Symbol, signAssistResp.Simulation.AssetsApprovals[0].ContractAddress)
				sentAmount := signAssistResp.Simulation.AssetsSent[0].Amount
				receivedAmount := signAssistResp.Simulation.AssetsReceived[0].Amount
				approvedAmount := signAssistResp.Simulation.AssetsApprovals[0].Amount
				decimalSentAmount, _ := s.GetDecimals(sentAmount, strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals))
				decimalReceivedAmount, _ := s.GetDecimals(receivedAmount, strconv.Itoa(signAssistResp.Simulation.AssetsReceived[0].Decimals))
				decimalApprovedAmount, _ := s.GetDecimals(approvedAmount, strconv.Itoa(signAssistResp.Simulation.AssetsApprovals[0].Decimals))
				simulationResp.HumanReadableForm = fmt.Sprintf("Send %v %v & Receive %v %v by approving %v %v", decimalSentAmount, signAssistResp.Simulation.AssetsSent[0].Symbol, decimalReceivedAmount, signAssistResp.Simulation.AssetsReceived[0].Symbol, decimalApprovedAmount, signAssistResp.Simulation.AssetsApprovals[0].Symbol)
				simulationResp.Type = "Swap"
				simulationResp.SimulationData.AssetsSent = []AssetsSent{
					{
						ContractAddress: sentAddress,
						Name:            signAssistResp.Simulation.AssetsSent[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsSent[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals),
					},
				}
				simulationResp.SimulationData.AssetsReceived = []AssetsReceived{
					{
						ContractAddress: receivedAddress,
						Name:            signAssistResp.Simulation.AssetsReceived[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsReceived[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsReceived[0].Decimals),
					},
				}
				simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{
					{
						ContractAddress: approvedAddress,
						Name:            signAssistResp.Simulation.AssetsApprovals[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsApprovals[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsApprovals[0].Decimals),
					},
				}
			} else if signAssistResp.Simulation.TransactionInfo.Method == "approve" {
				approvedAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsApprovals[0].Symbol, signAssistResp.Simulation.AssetsApprovals[0].ContractAddress)
				approvedAmount := signAssistResp.Simulation.AssetsApprovals[0].Amount
				decimalApprovedAmount, _ := s.GetDecimals(approvedAmount, strconv.Itoa(signAssistResp.Simulation.AssetsApprovals[0].Decimals))
				simulationResp.HumanReadableForm = fmt.Sprintf("Approve to transfer up to %v %v", decimalApprovedAmount, signAssistResp.Simulation.AssetsApprovals[0].Symbol)
				simulationResp.Type = "Token Approval"
				simulationResp.SimulationData.AssetsSent = []AssetsSent{}
				simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
				simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{
					{
						ContractAddress: approvedAddress,
						Name:            signAssistResp.Simulation.AssetsApprovals[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsApprovals[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsApprovals[0].Decimals),
					},
				}
			} else if signAssistResp.Simulation.TransactionInfo.Method == "fallback" && len(signAssistResp.Simulation.AssetsSent) != 0 && len(signAssistResp.Simulation.AssetsReceived) != 0 && len(signAssistResp.Simulation.AssetsApprovals) == 0 {
				sentAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsSent[0].Symbol, signAssistResp.Simulation.AssetsSent[0].ContractAddress)
				receivedAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsReceived[0].Symbol, signAssistResp.Simulation.AssetsReceived[0].ContractAddress)
				sentAmount := signAssistResp.Simulation.AssetsSent[0].Amount
				receivedAmount := signAssistResp.Simulation.AssetsReceived[0].Amount
				decimalSentAmount, _ := s.GetDecimals(sentAmount, strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals))
				decimalReceivedAmount, _ := s.GetDecimals(receivedAmount, strconv.Itoa(signAssistResp.Simulation.AssetsReceived[0].Decimals))
				simulationResp.HumanReadableForm = fmt.Sprintf("Send %v %v & Receive %v %v", decimalSentAmount, signAssistResp.Simulation.AssetsSent[0].Symbol, decimalReceivedAmount, signAssistResp.Simulation.AssetsReceived[0].Symbol)
				simulationResp.Type = "Swap"
				simulationResp.SimulationData.AssetsSent = []AssetsSent{
					{
						ContractAddress: sentAddress,
						Name:            signAssistResp.Simulation.AssetsSent[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsSent[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals),
					},
				}
				simulationResp.SimulationData.AssetsReceived = []AssetsReceived{
					{
						ContractAddress: receivedAddress,
						Name:            signAssistResp.Simulation.AssetsReceived[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsReceived[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsReceived[0].Decimals),
					},
				}
				simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
			} else if len(signAssistResp.Simulation.AssetsSent) == 0 && len(signAssistResp.Simulation.AssetsReceived) == 0 && len(signAssistResp.Simulation.AssetsApprovals) == 0 {
				simulationResp.Type = "Asset Transfer"
				simulationResp.SimulationData.AssetsSent = []AssetsSent{}
				simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
				simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
			} else {
				sentAddress := s.GetFrontierSpecificNativeToken(intVal, signAssistResp.Simulation.AssetsSent[0].Symbol, signAssistResp.Simulation.AssetsSent[0].ContractAddress)
				sentAmount := signAssistResp.Simulation.AssetsSent[0].Amount
				decimalSentAmount, _ := s.GetDecimals(sentAmount, strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals))
				simulationResp.HumanReadableForm = fmt.Sprintf("Send %v %v", decimalSentAmount, signAssistResp.Simulation.AssetsSent[0].Symbol)
				simulationResp.Type = "Asset Transfer"
				simulationResp.SimulationData.AssetsSent = []AssetsSent{
					{
						ContractAddress: sentAddress,
						Name:            signAssistResp.Simulation.AssetsSent[0].Name,
						Symbol:          signAssistResp.Simulation.AssetsSent[0].Symbol,
						Decimals:        strconv.Itoa(signAssistResp.Simulation.AssetsSent[0].Decimals),
					},
				}
				simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
				simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
			}
		} else {
			intGasUsed := strconv.Itoa(signAssistResp.Simulation.TransactionInfo.GasUsed)
			simulationResp.ChainId = request.Chain
			simulationResp.ProviderId = "2"
			simulationResp.SimulationResult = signAssistResp.Simulation.TransactionInfo.SimulationResult.Status
			simulationResp.Type = ""
			simulationResp.Action = "NONE"
			simulationResp.GasUsed = intGasUsed
			simulationResp.HumanReadableForm = "Simulation Failed"
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		}
		return &simulationResp, nil
	} else {
		return nil, status.Errorf(codes.Unavailable, "Unsupported Chain: %s", request.Chain)
	}
}

func (s *SignAssist) GetDecimals(value string, decimals string) (string, error) {
	floatValue, _ := strconv.ParseFloat(value, 64)
	floatDecimals, _ := strconv.ParseFloat(decimals, 64)
	powValue := math.Pow(10, floatDecimals)

	res := floatValue / powValue
	stringRes := fmt.Sprintf("%f", res)
	return stringRes, nil
}

func (s *SignAssist) GetFrontierSpecificNativeToken(chain int, symbol string, address string) string {
	var tokenAddress string
	info := s.utils.GetWalletInfo(strconv.Itoa(chain))
	if info.CurrencySymbol == symbol && info.ChainID == chain {
		tokenAddress = info.NativeTokenInfo.Address
	} else {
		tokenAddress = address
	}
	return tokenAddress
}
