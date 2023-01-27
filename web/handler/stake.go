package handler

// import (
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"bridge-allowance/utils/models"
// 	"context"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// )

// // Opportunities @Summary     Get Opportunities
// // @Tags         Stake
// // @Accept       json
// // @Produce      json
// // @Description  `EVM supported chains list` :- ethereum
// // @Description  `Cosmos supported chains list` :- axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence
// // @Description                irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,cronos-cosmos,evmos-cosmos,terra-classic,terra2,bluzelle
// // @Param        current   query      string  true  "chain"
// // @Router       /stake/opportunity [get]
// func (h *handler) Opportunities(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	chain := ctx.Query("current")
// 	chainGroup := chain
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	}
// 	cosmosOpportunitiesRequest := &pb.GetOpportunitiesRequest{
// 		Chain: chain,
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.GetOpportunites(context2, cosmosOpportunitiesRequest)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		var opportunities models.Opportunities
// 		err = json.Unmarshal(resp.Opportunities, &opportunities)
// 		if err != nil {
// 			utils.APIResponse(ctx, "json unmarshalling error "+err.Error(), codes.Internal, http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, opportunities)
// 		return
// 	}
// 	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, nil)
// 	return
// }
