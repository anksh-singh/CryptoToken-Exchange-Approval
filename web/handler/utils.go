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

// AllowanceHandler  godoc
// @Summary      get token allowance
// @Tags         Utils
// @Accept       json
// @Produce      json
// @Param        chain   query      string  true  "chain"  Enums(solana,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,evmos,fuse,cronos,astar,tomochain,zksync,,kava,osmosis,terra,band)
// @Param        contract   query      string  true  "contract"
// @Param        owner   query      string  true  "owner"
// @Param        spender   query      string  true  "spender"
// @Router       /utils/getAllowance [get]
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

// EnsHandler  godoc
// @Summary      ENS domain resolution
// @Tags         Utils
// @Accept       json
// @Produce      json
// @Param        domain   query      string  true  "domain"
// @Router       /utils/ens [get]
func (h *handler) EnsHandler(ctx *gin.Context) {
	ensRequest := &pb.ENSRequest{
		Domain: ctx.Query("domain"),
	}
	// Domain resolver
	address := ""
	var err error
	//domain := strings.Split(ensRequest.Domain, ".")[1]

	address, err = h.util.ResolveENSAddress(ensRequest.Domain)
	if err != nil {
		err = nil
		address, err = h.util.ResolveZNSAddress(ensRequest.Domain)
		if err != nil || address == "" {
			address, err = h.util.ResolveUNSAddress(ensRequest.Domain)
		}
	}
	response := pb.ENSResponse{Address: address}
	if err != nil || address == "" {
		//statusCode, _ := status.FromError(err)
		utils.APIResponse(ctx, "cannot resolve domain", codes.Internal, http.MethodGet, nil)
		return
	}
	utils.APIResponse(ctx, "Domain resolution is successful", codes.OK, http.MethodGet, response)
	return
}
