package handler

// import (
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"context"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// )

// // NftCollectionsHandler GetNftCollections godoc
// // @Summary      get nft collections
// // @Tags         NFT
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(nonevm,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex)
// // @Param        address   path      string  true  "address"
// // @Param        page   query      string  true  "page"
// // @Param        page-size   query      string  true  "page-size"
// // @Router       /{chain}/address/{address}/nftcollections [get]
// func (h *handler) NftCollectionsHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	address := ctx.Param("address")
// 	chain := ctx.Param("chain")
// 	chainGroup := chain
// 	pageSize := ctx.Query("page-size")
// 	page := ctx.Query("page")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || valid == false {
// 		h.logger.Errorf(" NftCollectionsHandler Logging Error  is : %v", err.Error())
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}

// 	NftCollectionRequest := &pb.NftCollectionRequest{
// 		Address:  validAddress,
// 		Page:     page,
// 		PageSize: pageSize,
// 		Chain:    chain,
// 	}

// 	err = h.util.ValidateNftCollectionRequest(NftCollectionRequest)
// 	if err != nil {
// 		s, _ := status.FromError(err)
// 		utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 		return
// 	}

// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.GetNftCollections(context2, NftCollectionRequest)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}

// 		utils.APIResponse(ctx, "NFT Collection fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }
