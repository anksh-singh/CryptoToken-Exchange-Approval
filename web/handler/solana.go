package handler

// import (
// 	"bridge-allowance/utils"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// )

// // GetRecentBlockHash godoc
// // @Summary      Get Recent Block hash
// // @Tags         Solana
// // @Accept       json
// // @Produce      json
// // @Router       /solana/getRecentBlockhash [get]
// func (h *handler) GetRecentBlockHash(ctx *gin.Context) {
// 	resp, err := h.nonEVMHandler.GetRecentBlockHash()
// 	if err != nil {
// 		s, _ := status.FromError(err)
// 		utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	utils.APIResponse(ctx, "Recent Blockhash has been fetched successfully", codes.OK, http.MethodGet, resp)
// 	return
// }
