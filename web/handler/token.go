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

// // TokenDetailHandler  TokenDetail godoc
// // @Summary      get token details
// // @Description  get token details
// // @Tags         Token
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,near,aptos,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        contractAddress   query      string  true  "contractAddress"
// // @Router       /token/tokenDetail [get]
// func (h *handler) TokenDetailHandler(ctx *gin.Context) {
// 	contractAddress := ctx.Query("contractAddress")
// 	chain := ctx.Query("chain")
// 	chainGroup := chain
// 	if h.util.IsEVM(chain) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(contractAddress, chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	tokenDetailReq := &pb.TokenDetailRequest{
// 		ContractAddress: validAddress,
// 		Chain:           chain, //Not chain group , should be specific chain name
// 	}
// 	tokenDetailRes, err := h.coingecko.GetTokenDetail(tokenDetailReq)
// 	if err != nil {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	utils.APIResponse(ctx, "Token details fetched successfully", codes.OK, http.MethodGet, tokenDetailRes)

// }

// // TokenInfoHandler TokenInfo godoc
// // @Summary      get token info
// // @Description  get token info
// // @Tags         Token
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,near,aptos,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        token   query      string  true  "token"
// // @Param        range   query      string  true  "range"
// // @Router       /token/info [get]
// func (h *handler) TokenInfoHandler(ctx *gin.Context) {
// 	token := ctx.Query("token")
// 	tokenInfoRange := ctx.Query("range")
// 	chain := ctx.Query("chain")
// 	tokenInfoReq := &pb.TokenInfoRequest{
// 		Token: token,
// 		Range: tokenInfoRange,
// 		Chain: chain,
// 	}
// 	tokenInfoRes, err := h.coingecko.GetTokenInfo(tokenInfoReq)
// 	if err != nil {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	utils.APIResponse(ctx, "Token info fetched successfully", codes.OK, http.MethodGet, tokenInfoRes)
// }

// // GetTokenPrice  godoc
// // @Summary      get token price
// // @Tags         Token
// // @Accept       json
// // @Produce      json
// // @Description  `EVM supported chains list` :- arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,zksync
// // @Description  `Cosmos supported chains list` :- axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence
// // @Description                irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,cronos-cosmos,evmos-cosmos,terra-classic,terra2
// // @Description  `Other chains list` :- solana,near,aptos
// // @Param        tokenPrice   query   models.GetTokenPrice  true  "get token price"
// // @Router       /token/price [get]
// func (h *handler) GetTokenPrice(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	getTokenPriceReq := &pb.TokenPriceRequest{
// 		Currency: ctx.Query("vs_currencies"),
// 		Chain:    ctx.Query("ids"),
// 	}
// 	err := h.util.ValidateTokenPriceRequest(getTokenPriceReq)
// 	if err != nil {
// 		//Defaults to USD if currency isn't supplied
// 		getTokenPriceReq.Currency = "$"
// 	}
// 	//TODO:Need to be refactored
// 	chainGroup := ctx.Query("ids")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.TokenPrice(context2, getTokenPriceReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Token price is fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ApprovalHandler  godoc
// // @Summary      get token approve
// // @Tags         Token
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,near,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,evmos,fuse,cronos,astar,zksync)
// // @Param        target   query      string  true  "target"
// // @Param        token   query      string  true  "token"
// // @Router       /token/approve [get]
// func (h *handler) ApprovalHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	getApproveReq := &pb.ApprovalRequest{
// 		Target: ctx.Query("target"),
// 		Token:  ctx.Query("token"),
// 		Chain:  ctx.Query("chain"),
// 	}
// 	chainGroup := ctx.Query("chain")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.Approve(context2, getApproveReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Token Approve is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // UserHandler GetBalances godoc
// // @Summary      get balance of tokens by address
// // @Description  Fetches the balances of all tokens associated with public key
// // @Tags         User
// // @Accept       json
// // @Produce      json
// // @Param        chain	path	string  true  "chain"  Enums(ethereum,polygon,bsc)
// // @Param        address   query		string  true  "address"
// // @Param		 contractAddress	query	string true "contractAddress"
// // @Router       /{chain}/token/userdata [get]
// func (h *handler) UserHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	userDataRequest := &pb.UserDataRequest{
// 		Chain:    ctx.Param("chain"),
// 		Address:  ctx.Query("address"),
// 		Contract: ctx.Query("contractAddress"),
// 	}
// 	chainGroup := userDataRequest.Chain
// 	chain := userDataRequest.Chain
// 	if h.util.IsEVM(userDataRequest.Chain) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(userDataRequest.Address, chainGroup, chain)
// 	if err != nil || !valid {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	userDataRequest.Address = validAddress
// 	//This API is more of a utility. Hence, routing requests to evm adapter
// 	if grpcClient, ok := h.grpcClient[utils.EVM]; ok {
// 		resp, err := grpcClient.UserData(context2, userDataRequest)
// 		if err != nil {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		if resp == nil {
// 			resp = &pb.UserDataResponse{}
// 		}
// 		utils.APIResponse(ctx, "User Data fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // BulkApprovalHandler  godoc
// // @Summary      get token approve
// // @Tags         Token
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,evmos,fuse,cronos,astar)
// // @Param        target   query      string  true  "target"
// // @Param        token   query      string  true  "token"
// // @Router       /token/bulk-approve [get]
// func (h *handler) BulkApprovalHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	getApproveReq := &pb.ApprovalRequest{
// 		Target: ctx.Query("target"),
// 		Token:  ctx.Query("token"),
// 		Chain:  ctx.Query("chain"),
// 	}
// 	chainGroup := ctx.Query("chain")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.BulkApproval(context2, getApproveReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }
