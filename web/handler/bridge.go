package handler

// import (
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"context"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// 	"strconv"
// )

// // GetBridgeChains BridgeChains godoc
// // @Summary      get Bridge Chain  List
// // @Description  get Bridge chain list for bridge
// // @Tags         Bridge
// // @Accept       json
// // @Produce      json
// // @Param        bridge_provider   query  string  true "bridge_provider" Enums(lifi,socket,xy, router,debridge)
// // @Router       /bridge/chains [get]
// func (h *handler) GetBridgeChains(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()

// 	adapter := "bridge"
// 	bridgeProvider := ctx.Query("bridge_provider")
// 	if bridgeProvider == "" {
// 		// default bridgeProvider is lifi if no provider is specified
// 		bridgeProvider = "lifi"
// 	}
// 	if grpcClient, ok := h.grpcClient[adapter]; ok {
// 		resp, err := grpcClient.BridgeChain(context2, &pb.BridgeChainRequest{
// 			BridgeProvider: bridgeProvider,
// 		})
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Bridge chain list fetched successfully.", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "service not available", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetBridgeChainTokens BridgeChainTokens godoc
// // @Summary      get Bridge Chain token  List
// // @Description  get Bridge chain tokens list for bridge
// // @Tags         Bridge
// // @Accept       json
// // @Produce      json
// // @Param        from_chain   query    string  true  "from_chain"
// // @Param        to_chain   query    string  true  "to_chain"
// // @Param        from_token   query    string  true  "from_token"
// // @Param		 full_list	  query		boolean  false "full_list"
// // @Param        bridge_provider   query  string  true "bridge_provider" Enums(lifi,socket,xy, router,debridge)
// // @Router       /bridge/chains/tokens [get]
// func (h *handler) GetBridgeChainTokens(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	fromChain := ctx.Query("from_chain")
// 	toChain := ctx.Query("to_chain")
// 	fromToken := ctx.Query("from_token")
// 	fullList, err := strconv.ParseBool(ctx.Query("full_list"))
// 	bridgeProvider := ctx.Query("bridge_provider")
// 	if bridgeProvider == "" {
// 		// default bridgeProvider is lifi if no provider is specified
// 		bridgeProvider = "lifi"
// 	}
// 	if err != nil {
// 		fullList = false
// 	}
// 	reqChainTokens := &pb.BridgeChainTokensRequest{
// 		FromChain:      fromChain,
// 		FromToken:      fromToken,
// 		ToChain:        toChain,
// 		FullList:       fullList,
// 		BridgeProvider: bridgeProvider,
// 	}
// 	adapter := "bridge"
// 	if grpcClient, ok := h.grpcClient[adapter]; ok {
// 		resp, err := grpcClient.BridgeChainTokens(context2, reqChainTokens)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		if len(resp.Tokens) == 0 {
// 			resp.Tokens = []*pb.BridgeTokens{}
// 		}
// 		utils.APIResponse(ctx, "Bridge chain token list fetched successfully.", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "service not available", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetBridgeQuote BridgeQuote godoc
// // @Summary      get Bridge quote
// // @Description  get bridge quote value
// // @Tags         Bridge
// // @Accept       json
// // @Produce      json
// // @Param        from_chain   query    string  true  "from_chain"
// // @Param        to_chain   query    string  true  "to_chain"
// // @Param        from_token   query    string  true  "from_token"
// // @Param        to_token   query    string  true  "to_token"
// // @Param        from_amount   query    string  true  "from_amount"
// // @Param        from_address   query    string  true  "from_address"
// // @Param        to_address   query    string  true  "to_address"
// // @Param        bridge_provider   query  string  true "bridge_provider" Enums(lifi,socket,xy, router,debridge)
// // @Router       /bridge/quote [get]
// func (h *handler) GetBridgeQuote(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	fromChain := ctx.Query("from_chain")
// 	fromToken := ctx.Query("from_token")
// 	toToken := ctx.Query("to_token")
// 	fromAmount := ctx.Query("from_amount")
// 	fromAddress := ctx.Query("from_address")
// 	toAddress := ctx.Query("to_address")
// 	toChain := ctx.Query("to_chain")
// 	bridgeProvider := ctx.Query("bridge_provider")
// 	if bridgeProvider == "" {
// 		// default bridgeProvider is lifi if no provider is specified
// 		bridgeProvider = "lifi"
// 	}
// 	reqChainTokenReq := &pb.BridgeQuoteRequest{
// 		FromChain:      fromChain,
// 		FromToken:      fromToken,
// 		ToToken:        toToken,
// 		FromAmount:     fromAmount,
// 		FromAddress:    fromAddress,
// 		ToAddress:      toAddress,
// 		ToChain:        toChain,
// 		BridgeProvider: bridgeProvider,
// 	}

// 	adapter := "bridge"
// 	if grpcClient, ok := h.grpcClient[adapter]; ok {
// 		resp, err := grpcClient.BridgeQuote(context2, reqChainTokenReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Bridge quote is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "service not available", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetBridgeTransaction BridgeTransaction godoc
// // @Summary      get Bridge Transaction
// // @Description  get bridge Transaction value
// // @Tags         Bridge
// // @Accept       json
// // @Produce      json
// // @Param        from_chain   query    string  true  "from_chain"
// // @Param        to_chain   query    string  true  "to_chain"
// // @Param        from_token   query    string  true  "from_token"
// // @Param        to_token   query    string  true  "to_token"
// // @Param        from_amount   query    string  true  "from_amount"
// // @Param        from_address   query    string  true  "from_address"
// // @Param        to_address   query    string  true  "to_address"
// // @Param        bridge_provider   query  string  true "bridge_provider" Enums(lifi,socket,xy, router,debridge)
// // @Router       /bridge/transaction [get]
// func (h *handler) GetBridgeTransaction(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	fromChain := ctx.Query("from_chain")
// 	fromToken := ctx.Query("from_token")
// 	toToken := ctx.Query("to_token")
// 	fromAmount := ctx.Query("from_amount")
// 	fromAddress := ctx.Query("from_address")
// 	toAddress := ctx.Query("to_address")
// 	toChain := ctx.Query("to_chain")
// 	bridgeProvider := ctx.Query("bridge_provider")
// 	if bridgeProvider == "" {
// 		// default bridgeProvider is lifi if no provider is specified
// 		bridgeProvider = "lifi"
// 	}
// 	reqChainTokenReq := &pb.BridgeTransactionRequest{
// 		FromChain:      fromChain,
// 		FromToken:      fromToken,
// 		ToToken:        toToken,
// 		FromAmount:     fromAmount,
// 		FromAddress:    fromAddress,
// 		ToAddress:      toAddress,
// 		ToChain:        toChain,
// 		BridgeProvider: bridgeProvider,
// 	}

// 	adapter := "bridge"
// 	if grpcClient, ok := h.grpcClient[adapter]; ok {
// 		resp, err := grpcClient.BridgeTransaction(context2, reqChainTokenReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Bridge transaction is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "service not available", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetBridgeTransactionStatus BridgeTransactionStatus godoc
// // @Summary      get Bridge Transaction Status
// // @Description  get bridge Transaction Status
// // @Tags         Bridge
// // @Accept       json
// // @Produce      json
// // @Param        bridge   query    string  true  "bridge"
// // @Param        from_chain   query    string  true  "from_chain"
// // @Param        to_chain   query    string  true  "to_chain"
// // @Param        tx_hash   query    string  true  "tx_hash"
// // @Param        bridge_provider   query  string  true "bridge_provider" Enums(lifi,socket,xy, router)
// // @Router       /bridge/status [get]
// func (h *handler) GetBridgeTransactionStatus(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	bridge := ctx.Query("bridge")
// 	fromChain := ctx.Query("from_chain")
// 	toChain := ctx.Query("to_chain")
// 	txHash := ctx.Query("tx_hash")
// 	bridgeProvider := ctx.Query("bridge_provider")
// 	if bridgeProvider == "" {
// 		// default bridgeProvider is lifi if no provider is specified
// 		bridgeProvider = "lifi"
// 	}
// 	reqChainTokenReq := &pb.BridgeTransactionStatusRequest{
// 		FromChain:      fromChain,
// 		Bridge:         bridge,
// 		TxHash:         txHash,
// 		ToChain:        toChain,
// 		BridgeProvider: bridgeProvider,
// 	}

// 	adapter := "bridge"
// 	if grpcClient, ok := h.grpcClient[adapter]; ok {
// 		resp, err := grpcClient.BridgeTransactionStatus(context2, reqChainTokenReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Transaction status fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "service not available", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }
