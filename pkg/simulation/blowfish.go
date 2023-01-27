package simulation

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type BlowFish struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewBlowFish(env *config.Config, logger *zap.SugaredLogger) *BlowFish {
	httpRequest := utils.NewHttpRequest(logger)
	return &BlowFish{env, logger, httpRequest}
}

func (b *BlowFish) NewSimulate(request SimulateTxRequest) (*SimulateTxResponse, error) {
	if val, ok := BlowFishSupportedChains[request.Chain]; ok {
		if val == "solana" {
			tx, err := b.SolanaSimulateTx(request)
			return tx, err
		} else {
			tx, err := b.SimulateTx(request)
			return tx, err
		}
	} else {
		return nil, status.Errorf(codes.Unavailable, "Unsupported Chain: %s", request.Chain)
	}
}

func (b *BlowFish) SolanaSimulateTx(request SimulateTxRequest) (*SimulateTxResponse, error) {
	blowFishURL := fmt.Sprintf(b.env.BlowFish.EndPoint+"/%s/v0/mainnet/scan/transactions", request.Chain)
	b.logger.Infof(blowFishURL)
	txPayload := BlowFishSolanaTxBody{
		Transactions: request.InputData,
		Metadata: BlowFishMetadata{
			Origin: request.Website,
		},
		UserAccount: request.From,
	}
	marshalTxBody, err := json.Marshal(txPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Marshalling error")
	}
	txBodyString := string(marshalTxBody)
	body, err := b.httpRequest.PostRequestWithHeaders(blowFishURL, txBodyString, "X-API-KEY", b.env.BlowFish.AccessKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Simulation Failed")
	}

	var blofishSolanaResp BlowFishSolanaResponse
	var simulationResp SimulateTxResponse
	err = json.Unmarshal(body, &blofishSolanaResp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}

	simulationResp.ProviderId = "1"
	if blofishSolanaResp.Action != "NONE" {
		simulationResp.Action = blofishSolanaResp.Action
		simulationResp.ActionMessage = blofishSolanaResp.Warnings[0].Message
	} else {
		simulationResp.Action = blofishSolanaResp.Action
	}
	result := make([]string, len(blofishSolanaResp.SimulationResults.ExpectedStateChanges))
	simulationResp.GasUsed = strconv.Itoa(blofishSolanaResp.SimulationResults.Raw.UnitsConsumed)
	if len(result) == 2 {
		simulationResp.HumanReadableForm = blofishSolanaResp.SimulationResults.ExpectedStateChanges[1].HumanReadableDiff + " & " + blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].HumanReadableDiff
	} else if len(result) == 1 {
		simulationResp.HumanReadableForm = blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].HumanReadableDiff
	}
	if blofishSolanaResp.SimulationResults.Error.HumanReadableError != "" {
		simulationResp.ChainId = request.Chain
		simulationResp.SimulationResult = blofishSolanaResp.SimulationResults.IsRecentBlockhashExpired
		simulationResp.Type = ""
		simulationResp.HumanReadableForm = blofishSolanaResp.SimulationResults.Error.HumanReadableError
		simulationResp.SimulationData.AssetsSent = []AssetsSent{}
		simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
		simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
	} else if len(result) == 1 {
		simulationResp.ChainId = request.Chain
		simulationResp.SimulationResult = blofishSolanaResp.SimulationResults.IsRecentBlockhashExpired
		if blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Kind == "SOL_TRANSFER" && len(result) == 1 {
			simulationResp.Type = "Asset Transfer"
			simulationResp.SimulationData.AssetsSent = []AssetsSent{
				{
					ContractAddress: "",
					Name:            blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Name,
					Symbol:          blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Symbol,
					Decimals:        strconv.Itoa(blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Decimals),
				},
			}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		}
	} else if len(result) == 2 {
		simulationResp.ChainId = request.Chain
		simulationResp.Type = "Swap"
		simulationResp.SimulationData.AssetsSent = []AssetsSent{
			{
				ContractAddress: "",
				Name:            blofishSolanaResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Name,
				Symbol:          blofishSolanaResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Symbol,
				Decimals:        strconv.Itoa(blofishSolanaResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Decimals),
			},
		}
		simulationResp.SimulationData.AssetsReceived = []AssetsReceived{
			{
				ContractAddress: "",
				Name:            blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Name,
				Symbol:          blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Symbol,
				Decimals:        strconv.Itoa(blofishSolanaResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Decimals),
			},
		}
		simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
	}
	return &simulationResp, nil
}

func (b *BlowFish) SimulateTx(request SimulateTxRequest) (*SimulateTxResponse, error) {
	chain := BlowFishSupportedChains[request.Chain]
	blowFishURL := fmt.Sprintf(b.env.BlowFish.EndPoint+"/%s/v0/mainnet/scan/transaction", chain)
	b.logger.Infof(blowFishURL)
	txPayload := BlowFishEVMTxBody{
		TxObject: BlowFishTxObject{
			From:  request.From,
			To:    request.To,
			Data:  strings.Join(request.InputData, " "),
			Value: request.Value,
		},
		Metadata: BlowFishMetadata{
			Origin: request.Website,
		},
		UserAccount: request.From,
	}
	marshalTxBody, err := json.Marshal(txPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Marshalling error")
	}
	txBodyString := string(marshalTxBody)
	body, err := b.httpRequest.PostRequestWithHeaders(blowFishURL, txBodyString, "X-Api-Key", b.env.BlowFish.AccessKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Simulation Failed")
	}
	var blowFishResp BlowFishResponse
	var simulationResp SimulateTxResponse
	err = json.Unmarshal(body, &blowFishResp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Unmarshalling error")
	}
	var simulationResult bool
	if blowFishResp.SimulationResults.Error.Kind == "" {
		simulationResult = true
	} else {
		simulationResult = false
	}
	simulationResp.ProviderId = "1"
	if blowFishResp.Action != "NONE" {
		simulationResp.Action = blowFishResp.Action
		simulationResp.ActionMessage = blowFishResp.Warnings[0].Message
	} else {
		simulationResp.Action = blowFishResp.Action
	}
	result := make([]string, len(blowFishResp.SimulationResults.ExpectedStateChanges))
	simulationResp.GasUsed = blowFishResp.SimulationResults.Gas.GasLimit
	if len(result) == 2 {
		simulationResp.HumanReadableForm = blowFishResp.SimulationResults.ExpectedStateChanges[1].HumanReadableDiff + " & " + blowFishResp.SimulationResults.ExpectedStateChanges[0].HumanReadableDiff
	} else if len(result) == 1 {
		simulationResp.HumanReadableForm = blowFishResp.SimulationResults.ExpectedStateChanges[0].HumanReadableDiff
	}

	if blowFishResp.SimulationResults.Error.Kind != "SIMULATION_FAILED" {
		simulationResp.ChainId = ChainIds[chain]
		simulationResp.SimulationResult = simulationResult
		if blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Kind == "ERC20_APPROVAL" {
			simulationResp.Type = "Token Approval"
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{
				{
					ContractAddress: blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Contract.Address,
					Name:            blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Name,
					Symbol:          blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Symbol,
					Decimals:        strconv.Itoa(blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Decimals),
				},
			}
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
		} else if blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Kind == "NATIVE_ASSET_TRANSFER" || (len(result) == 1) || blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Kind == "ERC20_TRANSFER" {
			simulationResp.ActionMessage = "You are transferring ER20 tokens directly to their own token contract. In most cases this will lead to you losing them forever."
			simulationResp.Type = "Asset Transfer"
			simulationResp.SimulationData.AssetsSent = []AssetsSent{
				{
					ContractAddress: blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Contract.Address,
					Name:            blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Name,
					Symbol:          blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Symbol,
					Decimals:        strconv.Itoa(blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Decimals),
				},
			}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		} else if len(result) == 2 {
			simulationResp.Type = "Swap"
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{
				{
					ContractAddress: blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Contract.Address,
					Name:            blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Name,
					Symbol:          blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Symbol,
					Decimals:        strconv.Itoa(blowFishResp.SimulationResults.ExpectedStateChanges[0].RawInfo.Data.Decimals),
				},
			}
			simulationResp.SimulationData.AssetsSent = []AssetsSent{
				{
					ContractAddress: blowFishResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Contract.Address,
					Name:            blowFishResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Name,
					Symbol:          blowFishResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Symbol,
					Decimals:        strconv.Itoa(blowFishResp.SimulationResults.ExpectedStateChanges[1].RawInfo.Data.Decimals),
				},
			}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		} else if len(blowFishResp.SimulationResults.ExpectedStateChanges) == 0 {
			simulationResp.SimulationData.AssetsSent = []AssetsSent{}
			simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
			simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
		}
	} else {
		simulationResp.ChainId = ChainIds[request.Chain]
		simulationResp.SimulationResult = simulationResult
		simulationResp.Type = ""
		simulationResp.Action = blowFishResp.Action
		simulationResp.ActionMessage = blowFishResp.SimulationResults.Error.ParsedErrorMessage
		simulationResp.GasUsed = blowFishResp.SimulationResults.Gas.GasLimit
		simulationResp.HumanReadableForm = blowFishResp.SimulationResults.Error.HumanReadableError
		simulationResp.SimulationData.AssetsSent = []AssetsSent{}
		simulationResp.SimulationData.AssetsReceived = []AssetsReceived{}
		simulationResp.SimulationData.ApprovalAssets = []ApprovalAssets{}
	}
	return &simulationResp, nil
}

func (b *BlowFish) GetFrontierSpecificNativeToken(chain string, address string) string {
	var tokenAddress string
	if chain == "polygon" || chain == "137" && address == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		tokenAddress = "0x0000000000000000000000000000000000001010"
	}
	return tokenAddress
}
