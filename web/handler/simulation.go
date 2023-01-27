package handler

// import (
// 	"bridge-allowance/pkg/simulation"
// 	"bridge-allowance/utils"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// )

// // SimulateTx godoc
// // @Summary      Simulate Transaction
// // @Tags         Simulation
// // @Accept       json
// // @Produce      json
// // @param		  SimulateTxRequest body  simulation.SimulateTxRequest true "msg"
// // @Router       /simulateTx [post]
// func (h *handler) SimulateTx(ctx *gin.Context) {
// 	var payload simulation.SimulateTxRequest
// 	var err error
// 	decoder := json.NewDecoder(ctx.Request.Body)
// 	err = decoder.Decode(&payload)
// 	if err != nil {
// 		h.logger.Errorf("Error %s", err)
// 		utils.APIResponse(ctx, "Error decoding input ", codes.Internal, http.MethodPost, nil)
// 		return
// 	}
// 	tx, err := h.simulation.SimulateTx(payload)
// 	if err != nil {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodPost, nil)
// 		return
// 	}
// 	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, tx)
// 	return
// }
