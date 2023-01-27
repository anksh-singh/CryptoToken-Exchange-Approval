package handler

// import (
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"context"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// 	"time"
// )

// // BalanceHandler GetBalances godoc
// // @Summary      get balance of tokens by address
// // @Description  Fetches the balances of all tokens associated with public key
// // @Description  `EVM supported chains list` :- arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync
// // @Description  `Cosmos supported chains list` :- axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence
// // @Description                irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,cronos-cosmos,evmos-cosmos,terra-classic,terra2,bluzelle
// // @Description  `Other chains list` :- solana,near,aptos
// // @Tags         Balance
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"
// // @Param        address   path      string  true  "address"
// // @Router       /{chain}/address/{address}/balances [get]
// func (h *handler) BalanceHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	address := ctx.Param("address")
// 	chain := ctx.Param("chain")
// 	var validAddress string
// 	validAddress = address
// 	chainGroup := chain
// 	h.logger.Info("Before validation Validated Address: ", address)
// 	if h.util.IsEVM(chain) {
// 		chainGroup = utils.EVM
// 		isEVM, evmValidAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 		h.logger.Info("Validated Address: ", evmValidAddress)
// 		if err != nil || isEVM == false {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		validAddress = evmValidAddress
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 		isCOSMOS, cosmosValidAddress, err := h.util.ValidateCosmosAddress(address, chain)
// 		h.logger.Info("Validated Address: ", cosmosValidAddress)
// 		if err != nil || isCOSMOS == false {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		validAddress = cosmosValidAddress
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	h.logger.Info("Before validation Validated Address: ", address)
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	h.logger.Info("Validated Address: ", validAddress)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return

// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		getBalanceReq := &pb.BalanceRequest{
// 			Address: validAddress,
// 			Chain:   chain, //Not chain group , should be specific chain name
// 		}
// 		if chainGroup != utils.COSMOS {
// 			resp, err := grpcClient.Balance(context2, getBalanceReq)
// 			if err != nil {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			if resp != nil {
// 				resp.Address = address
// 				resp.QuoteCurrency = "USD" //as of now only USD
// 			}
// 			//Grpc status code 0
// 			if len(resp.Token) == 0 {
// 				resp.Token = []*pb.TokenBalance{}
// 			}
// 			utils.APIResponse(ctx, "Balance Fetched Successfully", codes.OK, http.MethodGet, resp)
// 			return
// 		} else {
// 			resp, err := grpcClient.CosmosAssets(context2, getBalanceReq)
// 			if err != nil {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			if resp != nil {
// 				resp.Address = address
// 				resp.QuoteCurrency = "USD" //as of now only USD
// 			}
// 			//Grpc status code 0
// 			if len(resp.Token) == 0 {
// 				resp.Token = []*pb.CosmosTokenBalance{}
// 			}
// 			utils.APIResponse(ctx, "Balance Fetched Successfully", codes.OK, http.MethodGet, resp)
// 			return
// 		}
// 	} else {
// 		//grpc unavailable 14
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return

// 	}
// }
