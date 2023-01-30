package handler

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) AllowanceHandler(ctx *gin.Context) {
	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)

	defer cancel()
	chainGroup := ctx.Query("chain")
	chain := ctx.Query("chain")
	if h.util.IsEVM(chainGroup) {
		chainGroup = utils.EVM
	} else if h.util.IsCosmos(chainGroup) {
		chainGroup = utils.COSMOS
	} else if h.util.IsNonEVM(chainGroup) {
		chainGroup = "nonevm"
	}
	valid, validOwnerAddress, err := h.util.ValidateAddress(ctx.Query("owner"), chainGroup, chain)
	if err != nil || valid == false {
		statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
		return
	}
	getTokenAllowanceReq := &pb.AllowanceRequest{
		Contract: ctx.Query("contract"),
		Owner:    validOwnerAddress,
		Spender:  ctx.Query("spender"),
		Chain:    ctx.Query("chain"),
	}

	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
		resp, err := grpcClient.Allowance(context2, getTokenAllowanceReq)
		if err != nil {
			s, _ := status.FromError(err)
			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
			return
		}
		utils.APIResponse(ctx, "Token allowance is success", codes.OK, http.MethodGet, resp)
		return
	} else {
		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
		return
	}
}
