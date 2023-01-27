package handler

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"bridge-allowance/web/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

// GetValidators godoc
// @Summary      Get GetValidators List
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   query      string  true  "chain"  Enums(axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence,irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,terra-classic,evmos-cosmos,bluzelle,terra2)
// @Param        address   query    string  true  "address"
// @Param        testnet   query    string  true  "testnet" Enums(false,true)
// @Router       /cosmos/validators [get]
func (h *handler) GetValidators(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	address := ctx.Query("address")
	chain := ctx.Query("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	h.logger.Info("Before validation Validated Address: ", address)
	valid, validAddress, err := h.util.ValidateCosmosAddress(address, chain)
	h.logger.Info("Validated Address: ", validAddress)
	if err != nil || valid == false {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return
	}
	getCosmosValidatorsRequest := &pb.CosmosValidatorsRequest{
		Address: address,
		Testnet: false,
		Chain:   chain,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.GetValidators(context2, getCosmosValidatorsRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
		return
	}
	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, nil)
	return
}

// GetCosmosAprRates godoc
// @Summary      Get Max APR Rate from all validator List
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   query      string  true  "chain"  Enums(axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence,irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,terra-classic,evmos-cosmos,bluzelle,terra2)
// @Param        testnet   query    string  true  "testnet" Enums(false,true)
// @Router       /cosmos/apr [get]
func (h *handler) GetCosmosAprRates(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	chain := ctx.Query("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	getCosmosAPRRequest := &pb.CosmosAprRatesRequest{
		Testnet: false,
		Chain:   chain,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.GetCosmosAprRates(context2, getCosmosAPRRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
		return
	}
	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, nil)
	return
}

// GetDelegations godoc
// @Summary      Get Delegations List
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   query      string  true  "chain"  Enums(axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence,irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,terra-classic,evmos-cosmos,bluzelle,terra2)
// @Param        address   query    string  true  "address"
// @Param        testnet   query    string  true  "testnet" Enums(false,true)
// @Router       /cosmos/delegations [get]
func (h *handler) GetDelegations(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	address := ctx.Query("address")
	chain := ctx.Query("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM

	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	h.logger.Info("Before validation Validated Address: ", address)
	valid, validAddress, err := h.util.ValidateCosmosAddress(address, chain)
	h.logger.Info("Validated Address: ", validAddress)
	if err != nil || valid == false {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return
	}
	getCosmosDelegationsRequest := &pb.CosmosDelegationsRequest{
		Address: address,
		Testnet: false,
		Chain:   chain,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.GetDelegations(context2, getCosmosDelegationsRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		if resp.Delegations == nil {
			resp.Delegations = []*pb.DelegationsInfo{}
		}
		for _, rewardsInfo := range resp.Delegations {
			if rewardsInfo.RewardsDetails.Reward == nil {
				rewardsInfo.RewardsDetails.Reward = []*pb.RewardsListInfo{}
			}
		}
		for _, indStakedVal := range resp.IndividualStakeValues {

			if indStakedVal.StakeBalance == "" {
				indStakedVal.StakeBalance = "0"
			}
			if indStakedVal.UnstakeBalance == "" {
				indStakedVal.UnstakeBalance = "0"
			}
			if indStakedVal.RewardsBalance == "" {
				indStakedVal.RewardsBalance = "0"
			}
			if indStakedVal.TotalStakedQuote == "" {
				indStakedVal.TotalStakedQuote = "0"
			}
			if indStakedVal.TotalUnstakedQuote == "" {
				indStakedVal.TotalUnstakedQuote = "0"
			}
			if indStakedVal.TotalRewardsQuote == "" {
				indStakedVal.TotalRewardsQuote = "0"
			}
		}
		if resp.IndividualStakeValues == nil {
			resp.IndividualStakeValues = []*pb.DenomsWiseValues{}
		}
		if resp.UnboundDelegations == nil {
			resp.UnboundDelegations = []*pb.UnDelegationsInfo{}
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
		return
	}
	utils.APIResponse(ctx, "Internal Error", codes.Unavailable, http.MethodGet, nil)
	return
}

func (h *handler) CosmosSendTx(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	chainGroup := ctx.Param("chain")
	if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	var payload models.CosmosSendTxRequest
	var bluzellePayload []byte
	if ctx.Param("chain") == "bluzelle" {
		var bluzelleRequest models.BluzelleSendTxRequest
		decoder := json.NewDecoder(ctx.Request.Body)
		err := decoder.Decode(&bluzelleRequest)
		if err != nil {
			h.logger.Errorf("error %s", err)
			utils.APIResponse(ctx, "Error decoding input message", codes.Internal, http.MethodPost, nil)
			return
		}
		bluzellePayload, err = json.Marshal(bluzelleRequest)
		if err != nil {
			h.logger.Errorf("error %s", err)
			utils.APIResponse(ctx, "Error decoding input message", codes.Internal, http.MethodPost, nil)
			return
		}
	} else {
		decoder := json.NewDecoder(ctx.Request.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			h.logger.Errorf("error %s", err)
			utils.APIResponse(ctx, "Error decoding input message", codes.Internal, http.MethodPost, nil)
			return
		}
	}

	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		sendTxReq := &pb.CosmosSendTxRequest{
			Chain:     ctx.Param("chain"),
			TxBytes:   payload.TxBytes,
			Mode:      payload.Mode,
			TxDetails: bluzellePayload,
		}
		resp, err := grpcClient.CosmosSendTx(context2, sendTxReq)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodPost, nil)
			return
		}
		utils.APIResponse(ctx, "Transaction submitted successfully", codes.OK, http.MethodPost, resp)
		return
	} else {
		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodPost, nil)
		return
	}
}

// GetCosmosCDPParams godoc
// @Summary      Get Cosmos CDP Params
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   query      string  true  "chain"  Enums(kava)
// @Param        testnet   query    string  true  "testnet" Enums(false,true)
// @Router       /cosmos/cdpParams [get]
func (h *handler) GetCosmosCDPParams(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	chain := ctx.Query("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	getCosmosCDPParametersRequest := &pb.CosmosCDPParametersRequest{
		Testnet: false,
		Chain:   chain,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.GetCosmosCDPParams(context2, getCosmosCDPParametersRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
		return
	}
	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, nil)
	return
}

// CosmosSimulateTx @Summary      Simulation of transaction to estimate the gas usage
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   path      string  true  "chain"  Enums(axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence,irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,terra-classic,evmos-cosmos,bluzelle,terra2)
// @param		 SimulateTxRequest body models.SimulateTxBody true "msg"
// @Router       /cosmos/{chain}/simulateTx [post]
func (h *handler) CosmosSimulateTx(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	chain := ctx.Param("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	var sendTxBody models.SimulateTxBody
	decoder := json.NewDecoder(ctx.Request.Body)
	err := decoder.Decode(&sendTxBody)
	if err != nil {
		h.logger.Errorf("error %s", err)
		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
		return
	}
	cosmosSimulateTxRequest := &pb.CosmosSimulateTxRequest{
		Chain:   chain,
		TxBytes: sendTxBody.TxBytes,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.CosmosSimulateTX(context2, cosmosSimulateTxRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodPost, nil)
			return
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodPost, resp)
		return
	}
	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodPost, nil)
	return
}

// CosmosGetBlockHeight @Summary Fetches the latest block height
// @Tags         Cosmos
// @Accept       json
// @Produce      json
// @Param        chain   path      string  true  "chain"  Enums(axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence,irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,terra-classic,evmos-cosmos,bluzelle,terra2)
// @Router       /cosmos/{chain}/getLatestBlockHeight [get]
func (h *handler) CosmosGetBlockHeight(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
	defer cancel()
	chain := ctx.Param("chain")
	chainGroup := chain
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	}
	cosmosSimulateTxRequest := &pb.GetCosmosBlockHeightRequest{
		Chain: chain,
	}
	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.CosmosGetBlockHeight(context2, cosmosSimulateTxRequest)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		h.logger.Infof("Response ", resp)
		utils.APIResponse(ctx, "Latest Block Height Fetched Successfully", codes.OK, http.MethodGet, resp)
		return
	}
	utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
	return
}
