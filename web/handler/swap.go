package handler

// import (
// 	"bridge-allowance/pkg/debankAPI"
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"bridge-allowance/web/models"
// 	"context"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// 	"strings"
// 	"time"
// )

// // ExchangeTokens godoc
// // @Summary      get Token List
// // @Description  get tokens by exchange type
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,klaytn,moonbeam,heco,gnosis,fuse,evmos,astar,cronos,tomochain,bttc,,kava,osmosis,terra,band, iotex)
// // @Param        exchange_type   query  string  true "exchange_type" Enums(0x,dodo,lifi, 1inch,xy,dzap,jupiter,pancakeswap,zeroswap, sushiswap,netswap,elk,paraswap,cowswap,kyber,xswap,diffusion,luaswap,arthswap,vvsfinance,soyfinance,quackswap,meshswap,klayswap)
// // @Router       /exchange/{chain}/tokens [get]
// func (h *handler) ExchangeTokens(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOut25)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		ExchangeTokenReq := &pb.ExchangeTokenRequest{
// 			Chain:        ctx.Param("chain"),
// 			ExchangeType: strings.ToLower(ctx.Query("exchange_type")),
// 		}
// 		resp, err := grpcClient.ExchangeTokens(context2, ExchangeTokenReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		var msg = "Token list fetched successfully"
// 		var code = codes.OK
// 		var emptyToken = make([]*pb.ExchangeTokenInfo, 0)
// 		if resp.ExchangeTokens == nil {
// 			resp.ExchangeTokens = emptyToken
// 		}
// 		utils.APIResponse(ctx, msg, code, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ExchangeQuote godoc
// // @Summary      get Exchange quote
// // @Description  get exchange quote
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,,arbitrum,fantom,klaytn,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,heco,gnosis,fuse,evmos,astar,cronos,tomochain,bttc,,kava,osmosis,terra,band,iotex)
// // @Param        exchange_type   query  string  true "exchange_type" Enums(0x,dodo,lifi,1inch,xy,jupiter,pancakeswap,zeroswap,sushiswap,netswap,elk,paraswap,cowswap,kyber,xswap,diffusion,luaswap,arthswap,vvsfinance,soyfinance,quackswap,meshswap,klayswap)
// // @Param        taker_address   query      string  true  "taker_address"
// // @Param        sell_token   query      string  true  "sell_token"
// // @Param        buy_token   query      string  true  "buy_token"
// // @Param        sell_amount   query      string  true  "sell_amount"
// // @Param        slippage   query      string  true  "slippage"
// // @Router       /exchange/{chain}/quote [get]
// func (h *handler) ExchangeQuote(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOut25)
// 	defer cancel()
// 	address := ctx.Query("taker_address")
// 	chain := ctx.Param("chain")
// 	chainGroup := chain
// 	if h.util.IsEVM(chain) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || !valid {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressSellToken := ctx.Query("sell_token")
// 	validSellToken, validAddressSellToken, err := h.util.ValidateAddress(addressSellToken, chainGroup, chain)
// 	if err != nil || !validSellToken {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressBuyToken := ctx.Query("buy_token")
// 	validBuyToken, validAddressBuyToken, err := h.util.ValidateAddress(addressBuyToken, chainGroup, chain)
// 	if err != nil || !validBuyToken {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		ExchangeQuoteReq := &pb.ExchangeQuoteRequest{
// 			Chain:        ctx.Param("chain"),
// 			ExchangeType: ctx.Query("exchange_type"),
// 			TakerAddress: validAddress,
// 			SellToken:    validAddressSellToken,
// 			BuyToken:     validAddressBuyToken,
// 			SellAmount:   ctx.Query("sell_amount"),
// 			Slippage:     ctx.Query("slippage"),
// 		}
// 		resp, err := grpcClient.ExchangeQuote(context2, ExchangeQuoteReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Exchange quote is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ExchangeMultiQuote godoc
// // @Summary      get Exchange Multi quote
// // @Description  get exchange multi quote
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,,arbitrum,fantom,klaytn,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,heco,gnosis,fuse,evmos,astar,cronos,tomochain,bttc,,kava,osmosis,terra,band, iotex)
// // @Param        exchange_type  query string  true "exchange_type" Enums(0x,dodo,lifi,1inch,xy,jupiter,pancakeswap,zeroswap,sushiswap,netswap,elk,paraswap,cowswap,kyber,xswap,diffusion,luaswap,arthswap,vvsfinance,soyfinance,dzap,quackswap)
// // @Param        taker_address   query      string  true  "taker_address"
// // @param		 models.MultiSwapRequestBody body models.MultiSwapRequestBody true "msg"
// // @Router       /exchange/{chain}/multi_quote [post]
// func (h *handler) ExchangeMultiQuote(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOut25)
// 	defer cancel()
// 	chain := ctx.Param("chain")
// 	chainGroup := chain
// 	if h.util.IsEVM(chain) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(ctx.Query("taker_address"), chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	var requestBody models.MultiSwapRequestBody
// 	decoder := json.NewDecoder(ctx.Request.Body)
// 	err = decoder.Decode(&requestBody)
// 	if err != nil {
// 		h.logger.Errorf("error %s", err)
// 		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
// 		return
// 	}

// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		// Create Grpc Request Body
// 		var grpcReqObj pb.ExchangeMultiQuoteRequest
// 		for _, item := range requestBody.SwapParams {
// 			grpcReqObj.ExchangeType = ctx.Query("exchange_type")
// 			grpcReqObj.Chain = chain
// 			grpcReqObj.TakerAddress = validAddress
// 			var multichainRequest pb.MultiChainRequests
// 			validFromToken, validAddressFromToken, err := h.util.ValidateAddress(item.SellToken, chainGroup, chain)
// 			if err != nil || !validFromToken {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			validToToken, validAddresstoToken, err := h.util.ValidateAddress(item.BuyToken, chainGroup, chain)
// 			if err != nil || !validToToken {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			multichainRequest.SellToken = validAddressFromToken
// 			multichainRequest.BuyToken = validAddresstoToken
// 			multichainRequest.Slippage = item.Slippage
// 			multichainRequest.SellAmount = item.SellAmount
// 			grpcReqObj.MultiChainRequests = append(grpcReqObj.MultiChainRequests, &multichainRequest)
// 		}
// 		resp, err := grpcClient.ExchangeMultiQuote(context2, &grpcReqObj)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Exchange quote is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ExchangeMultiSwap godoc
// // @Summary      get Exchange Multi swap
// // @Description  get exchange multi swap
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,,arbitrum,fantom,klaytn,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,heco,gnosis,fuse,evmos,astar,cronos,tomochain,bttc,,kava,osmosis,terra,band, iotex)
// // @Param        exchange_type  query string  true "exchange_type" Enums(0x,dodo,lifi,1inch,xy,jupiter,pancakeswap,zeroswap,sushiswap,netswap,elk,paraswap,kyber,xswap,diffusion,luaswap,arthswap,vvsfinance,soyfinance,dzap,quackswap)
// // @Param        taker_address   query      string  true  "taker_address"
// // @param		 models.MultiSwapRequestBody body models.MultiSwapRequestBody true "msg"
// // @Router       /exchange/{chain}/multi_swap [post]
// func (h *handler) ExchangeMultiSwap(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOut25)
// 	defer cancel()
// 	//address := ctx.Query("taker_address")
// 	chain := ctx.Param("chain")
// 	chainGroup := chain
// 	if h.util.IsEVM(chain) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(ctx.Query("taker_address"), chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	var requestBody models.MultiSwapRequestBody
// 	decoder := json.NewDecoder(ctx.Request.Body)
// 	err = decoder.Decode(&requestBody)
// 	if err != nil {
// 		h.logger.Errorf("error %s", err)
// 		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
// 		return
// 	}

// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		// Create Grpc Request Body
// 		var grpcReqObj pb.ExchangeMultiSwapRequest
// 		for _, item := range requestBody.SwapParams {
// 			grpcReqObj.ExchangeType = ctx.Query("exchange_type")
// 			grpcReqObj.Chain = chain
// 			grpcReqObj.TakerAddress = validAddress
// 			var multichainRequest pb.MultiChainRequests
// 			validFromToken, validAddressFromToken, err := h.util.ValidateAddress(item.SellToken, chainGroup, chain)
// 			if err != nil || !validFromToken {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			validToToken, validAddresstoToken, err := h.util.ValidateAddress(item.BuyToken, chainGroup, chain)
// 			if err != nil || !validToToken {
// 				statusCode, _ := status.FromError(err)
// 				utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 				return
// 			}
// 			multichainRequest.SellToken = validAddressFromToken
// 			multichainRequest.BuyToken = validAddresstoToken
// 			multichainRequest.Slippage = item.Slippage
// 			multichainRequest.SellAmount = item.SellAmount
// 			grpcReqObj.MultiChainRequests = append(grpcReqObj.MultiChainRequests, &multichainRequest)
// 		}
// 		resp, err := grpcClient.ExchangeMultiSwap(context2, &grpcReqObj)
// 		if len(resp.Approval) == 0 {
// 			resp.Approval = []*pb.ApprovalResponse{}
// 		}
// 		if resp.Transaction.MultiRouteData == nil {
// 			resp.Transaction.MultiRouteData = &pb.MultiSwapTxs{
// 				Data: make([]string, 0),
// 			}
// 		}
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Exchange swap is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ExchangeSwap godoc
// // @Summary      get exchange swap
// // @Description  get exchange swap
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,,arbitrum,fantom,klaytn,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,moonriver,moonbeam,heco,gnosis,fuse,evmos,astar,cronos,tomochain,bttc,,kava,osmosis,terra,band, iotex)
// // @Param        exchange_type   query      string  true "exchange_type" Enums(0x,dodo,lifi,1inch,xy,jupiter,pancakeswap,zeroswap,sushiswap,netswap,elk,paraswap,kyber,xswap,diffusion,luaswap,arthswap,vvsfinance,soyfinance,quackswap,meshswap,klayswap)
// // @Param        taker_address   query      string  true  "taker_address"
// // @Param        sell_token   query      string  true  "sell_token"
// // @Param        buy_token   query      string  true  "buy_token"
// // @Param        sell_amount   query      string  true  "sell_amount"
// // @Param        slippage   query      string  true  "slippage"
// // @Router       /exchange/{chain}/swap [get]
// func (h *handler) ExchangeSwap(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOut25)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	chain := ctx.Param("chain")
// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(ctx.Query("taker_address"), chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressSellToken := ctx.Query("sell_token")
// 	validSellToken, validAddressSellToken, err := h.util.ValidateAddress(addressSellToken, chainGroup, chain)
// 	if err != nil || validSellToken == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressBuyToken := ctx.Query("buy_token")
// 	validBuyToken, validAddressBuyToken, err := h.util.ValidateAddress(addressBuyToken, chainGroup, chain)
// 	if err != nil || validBuyToken == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		ExchangeSwapReq := &pb.ExchangeSwapRequest{
// 			Chain:        ctx.Param("chain"),
// 			ExchangeType: ctx.Query("exchange_type"),
// 			TakerAddress: validAddress,
// 			SellToken:    validAddressSellToken,
// 			BuyToken:     validAddressBuyToken,
// 			SellAmount:   ctx.Query("sell_amount"),
// 			Slippage:     ctx.Query("slippage"),
// 		}
// 		resp, err := grpcClient.ExchangeSwap(context2, ExchangeSwapReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		if resp.MultiRouteData == nil {
// 			resp.MultiRouteData = &pb.MultiSwapTxs{
// 				Data: make([]string, 0),
// 			}
// 		}
// 		utils.APIResponse(ctx, "Exchange swap is success", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // ExchangeSignature godoc
// // @Summary      get exchange signature data
// // @Description  get exchange signature data
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(ethereum,fantom,polygon,celo,bsc, avalanche, gnosis, optimism)
// // @Param        exchange_type   query      string  true "exchange_type" Enums(zeroswap, cowswap)
// // @Param        taker_address   query      string  true  "taker_address"
// // @Param        sell_token   query      string  true  "sell_token"
// // @Param        buy_token   query      string  true  "buy_token"
// // @Param        sell_amount   query      string  true  "sell_amount"
// // @Param        slippage   query      string  true  "slippage"
// // @Router       /exchange/{chain}/signature [get]
// func (h *handler) ExchangeSignature(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	chain := ctx.Param("chain")
// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(ctx.Query("taker_address"), chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressSellToken := ctx.Query("sell_token")
// 	validSellToken, validAddressSellToken, err := h.util.ValidateAddress(addressSellToken, chainGroup, chain)
// 	if err != nil || validSellToken == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	addressBuyToken := ctx.Query("buy_token")
// 	validBuyToken, validAddressBuyToken, err := h.util.ValidateAddress(addressBuyToken, chainGroup, chain)
// 	if err != nil || validBuyToken == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		ExchangeSwapReq := &pb.ExchangeSignatureRequest{
// 			Chain:        ctx.Param("chain"),
// 			ExchangeType: ctx.Query("exchange_type"),
// 			TakerAddress: validAddress,
// 			SellToken:    validAddressSellToken,
// 			BuyToken:     validAddressBuyToken,
// 			SellAmount:   ctx.Query("sell_amount"),
// 			Slippage:     ctx.Query("slippage"),
// 		}
// 		resp, err := grpcClient.ExchangeSwapSignature(context2, ExchangeSwapReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}

// 		switch ctx.Query("exchange_type") {
// 		case "zeroswap":
// 			utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp.ZeroswapData)
// 			return
// 		case "cowswap":
// 			utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp.CowswapData)
// 			return
// 		default:
// 			utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
// 			return
// 		}

// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetFreeTradeCount FreeTradeCount godoc
// // @Summary      Get free trade count
// // @Description  Get free trade count
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(fantom,polygon,celo,bsc, avalanche,optimism)
// // @Param        exchange_type   query  string  true "exchange_type" Enums(zeroswap)
// // @Param        account   query      string  true  "account_address"
// // @Router       /exchange/{chain}/freeTradeCount [get]
// func (h *handler) FreeTradeCount(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	chain := ctx.Param("chain")
// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	address := ctx.Query("account")
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || !valid {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		inputRequest := &pb.FreeTradeCountRequest{
// 			Account:      validAddress,
// 			Chain:        ctx.Param("chain"),
// 			ExchangeType: ctx.Query("exchange_type"),
// 		}

// 		resp, err := grpcClient.FreeTradeCount(context2, inputRequest)
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

// // ExchangeSwapExecute  ExchangeZeroSwap godoc
// // @Summary      execute zero swap
// // @Description  execute zero swap
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(ethereum, fantom,polygon,celo,bsc, avalanche, gnosis,optimism)
// // @Param        exchange_type  query string  true "exchange_type" Enums(zeroswap, cowswap)
// // @Param        payload   body models.GasLessSwapBody false "payload"
// // @Router       /exchange/{chain}/gasless/swap [post]
// func (h *handler) ExchangeSwapExecute(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	chain := ctx.Param("chain")

// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	var zeroswapRequest *pb.ZeroSwapExecuteRequest
// 	var cowswapRequest *pb.CowSwapExecuteRequest
// 	var err error
// 	switch ctx.Query("exchange_type") {
// 	case "zeroswap":
// 		zeroswapRequest, err = h.zeroSwap.GetZeroswapExecuteRequest(ctx, chainGroup, chain)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodPost, nil)
// 			return
// 		}
// 	case "cowswap":
// 		cowswapRequest, err = h.cowSwap.GetCowSwapExecuteRequest(ctx, chainGroup, chain)
// 		if err != nil {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		ExchangeSwapExecuteReq := &pb.ExchangeSwapExecuteRequest{
// 			Chain:           ctx.Param("chain"),
// 			ExchangeType:    ctx.Query("exchange_type"),
// 			ZeroSwapPayload: zeroswapRequest,
// 			CowSwapPayload:  cowswapRequest,
// 		}

// 		resp, err := grpcClient.ExchangeSwapExecute(context2, ExchangeSwapExecuteReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		switch ctx.Query("exchange_type") {
// 		case "zeroswap":
// 			utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp.ExecuteResponse)
// 			return
// 		case "cowswap":
// 			info := h.util.GetSwapDataInfo(ctx.Param("chain"))
// 			walletInfo := h.util.GetWalletInfo(ctx.Param("chain"))
// 			nativeTokenAddress := h.cowSwap.GetNativeTokenAddress(walletInfo)
// 			wrappedTokenAddress := h.cowSwap.GetWNativeTokenAddress(info)

// 			if (strings.ToLower(nativeTokenAddress) == strings.ToLower(ExchangeSwapExecuteReq.CowSwapPayload.SellToken) && strings.ToLower(wrappedTokenAddress) == strings.ToLower(ExchangeSwapExecuteReq.CowSwapPayload.BuyToken)) ||
// 				(strings.ToLower(wrappedTokenAddress) == strings.ToLower(ExchangeSwapExecuteReq.CowSwapPayload.SellToken) && strings.ToLower(nativeTokenAddress) == strings.ToLower(ExchangeSwapExecuteReq.CowSwapPayload.BuyToken)) {
// 				utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp.ExchangeSwapResponse)
// 				return
// 			} else {
// 				utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp.ExecuteResponse)
// 				return
// 			}
// 		default:
// 			utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, resp)
// 			return
// 		}

// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GetTokenApproval TokenApproval godoc
// // @Summary      Get token approval
// // @Description  Get token approval
// // @Tags         Swap
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(fantom,polygon,celo,bsc, avalanche,optimism)
// // @Param        exchange_type   query  string  true "exchange_type" Enums(zeroswap)
// // @Param        token   path      string  true  "Token Address"
// // @Param		 gasless query 	   string  true  "gasless" Enums(true, false)
// // @Router       /exchange/{chain}/token/{token}/approve [get]
// func (h *handler) ExchangeTokenApprove(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	chain := ctx.Param("chain")
// 	if h.util.IsEVM(ctx.Param("chain")) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	address := ctx.Param("token")
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || !valid {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		inputRequest := &pb.TokenApprovalRequest{
// 			Token:        validAddress,
// 			Chain:        ctx.Param("chain"),
// 			Gasless:      ctx.Query("gasless"),
// 			ExchangeType: ctx.Query("exchange_type"),
// 		}

// 		resp, err := grpcClient.ExchangeTokenApprove(context2, inputRequest)
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

// // PositionHandler godoc
// // @Summary      Get Positions
// // @Tags         Positions
// // @Accept       json
// // @Produce      json
// // @Description  `Supported Chain List` :- "ethereum","bsc","gnosis","polygon","fantom","heco","avalanche","optimism","arbitrum","celo","moonriver","cronos","boba","metis","bttc","aurora","moonbeam","fuse","harmony","klaytn","astar","iotex","evmos"
// // @Param        payload   body models.GetPositionRequest false "payload"
// // @Router       /positions [post]
// func (h *handler) PositionHandler(ctx *gin.Context) {
// 	var payload models.GetPositionRequest
// 	decoder := json.NewDecoder(ctx.Request.Body)
// 	err := decoder.Decode(&payload)
// 	if err != nil {
// 		h.logger.Errorf("Error while decoding request body: %v", err)
// 		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
// 		return
// 	}

// 	var debankResponse debankAPI.PositionResponse
// 	var sonarWatchResponse *pb.GetPositionsResponse
// 	for _, item := range payload {
// 		var chainGroup string
// 		if h.util.IsEVM(item.Chain) {
// 			chainGroup = utils.EVM
// 		}
// 		if chainGroup == "evm" {
// 			var chainIds []string
// 			if item.Chain != "evm" {
// 				chainIds = []string{item.Chain}
// 			}
// 			positionPayload := &pb.PositionPayload{
// 				ChainIds:    chainIds,
// 				ProtocolIds: item.ProtocolIds,
// 			}
// 			positionRequest := &pb.PositionRequest{
// 				Address:         item.Address,
// 				PositionPayload: positionPayload,
// 			}
// 			resp, _ := h.debankAPI.GetPosition(positionRequest)
// 			debankResponse.Positions = append(debankResponse.Positions, resp.Positions...)
// 		} else if item.Chain == "solana" {
// 			context2, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 			defer cancel()
// 			if grpcClient, ok := h.grpcClient["nonevm"]; ok {
// 				inputRequest := &pb.PositionChainData{
// 					Chain:       item.Chain,
// 					Address:     item.Address,
// 					ProtocolIds: item.ProtocolIds,
// 				}
// 				resp, _ := grpcClient.GetPositions(context2, inputRequest)
// 				sonarWatchResponse = resp
// 			}
// 		}
// 	}

// 	// Merging response from debank and nonevm
// 	var debankStruct debankAPI.PositionResponse
// 	out, _ := json.Marshal(sonarWatchResponse)
// 	err = json.Unmarshal(out, &debankStruct)
// 	for _, item := range debankStruct.Positions {
// 		debankResponse.Positions = append(debankResponse.Positions, item)
// 	}
// 	utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, debankResponse)
// 	return
// }
