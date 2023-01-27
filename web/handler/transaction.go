package handler

// import (
// 	"bridge-allowance/pkg/grpc/proto/pb"
// 	"bridge-allowance/utils"
// 	"bridge-allowance/web/models"
// 	"context"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"net/http"
// )

// // ListTransactionHandler  ListTransactions godoc
// // @Summary      List Transaction
// // @Description  list transaction by public key
// // @Tags         Transaction
// // @Accept       json
// // @Produce      json
// // @Param        chain   path      string  true  "chain"  Enums(solana,near,aptos,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        address   path      string  true  "address"
// // @Param        page   query      string  true  "page"
// // @Param        page-size   query      string  true  "page-size"
// // @Param        tokenContractAddress   query      string  false  "tokenContractAddress"
// // @Router       /{chain}/address/{address}/transactions [get]
// func (h *handler) ListTransactionHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TimeOutOneMinute)
// 	defer cancel()
// 	pageSize := ctx.Query("page-size")
// 	page := ctx.Query("page")
// 	address := ctx.Param("address")
// 	chain := ctx.Param("chain")
// 	h.logger.Info("Before validation Validated Address: ", address)
// 	chainGroup := chain
// 	var validAddress string
// 	var err error
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 		isEvm, evmValidAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 		h.logger.Info("Validated Address: ", evmValidAddress)
// 		if err != nil || isEvm == false {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		validAddress = evmValidAddress
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 		isCosmos, cosmosValidAddress, err := h.util.ValidateCosmosAddress(address, chain)
// 		h.logger.Info("Validated Address: ", cosmosValidAddress)
// 		if err != nil || isCosmos == false {
// 			statusCode, _ := status.FromError(err)
// 			utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		validAddress = cosmosValidAddress

// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	listTransactionRequest := &pb.ListTransactionRequest{
// 		Address:              validAddress,
// 		Testnet:              false,
// 		Page:                 page,
// 		PageSize:             pageSize,
// 		TokenContractAddress: ctx.Query("tokenContractAddress"),
// 		Chain:                chain,
// 	}
// 	err = h.util.ValidateListTransactionRequest(listTransactionRequest)
// 	if err != nil {
// 		s, _ := status.FromError(err)
// 		utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		resp, err := grpcClient.ListTransaction(context2, listTransactionRequest)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		if len(resp.Transactions) == 0 {
// 			resp.Transactions = []*pb.TransactionData{}
// 		}
// 		var t = make([]*pb.TransactionInfo, 0)
// 		for _, transRes := range resp.Transactions {
// 			if transRes.Sent == nil {
// 				transRes.Sent = t
// 			}
// 			if transRes.Received == nil {
// 				transRes.Received = t
// 			}
// 			if transRes.Others == nil {
// 				transRes.Others = t
// 			}
// 		}
// 		utils.APIResponse(ctx, "List Transaction is successful", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // SendTransactionHandler   SendTransaction godoc
// // @Summary      Send Transaction
// // @Tags         Transaction
// // @Accept       json
// // @Produce      json
// // @Description  `EVM supported chains list` :- arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync
// // @Description  `Cosmos supported chains list` :- axelar,akash,bandchain,cosmoshub,crescent,cryptoorgchain,injective,juno,kujira,kava,osmosis,secretnetwork,sifchain,umee,regen,stargaze,sentinel,persistence
// // @Description   irisnet,agoric,shentu,impacthub,emoney,sommelier,bostrom,gravitybridge,stride,assetmantle,cronos-cosmos,evmos-cosmos,terra-classic,terra2,bluzelle
// // @Description  `Other chains list` :- solana,near
// // @Param        chain   path      string  true  "chain"
// // @param		 models.SendTxBody body string true "msg"
// // @Router       /{chain}/transactions [post]
// func (h *handler) SendTransactionHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), TxTimeOut) //Increased the Request Time-out due to retry attempts
// 	defer cancel()
// 	chainGroup := ctx.Param("chain")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 		h.CosmosSendTx(ctx)
// 		return
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	var sendTxBody models.SendTxBody
// 	decoder := json.NewDecoder(ctx.Request.Body)
// 	err := decoder.Decode(&sendTxBody)
// 	if err != nil {
// 		h.logger.Errorf("error %s", err)
// 		utils.APIResponse(ctx, "Error decoding input message", codes.Unavailable, http.MethodPost, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		sendTxReq := &pb.SendTransactionRequest{
// 			Msg:   sendTxBody.Message,
// 			Chain: ctx.Param("chain"),
// 		}
// 		resp, err := grpcClient.SendTransaction(context2, sendTxReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodPost, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Send Transaction is successful", codes.OK, http.MethodPost, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodPost, nil)
// 		return
// 	}
// }

// // TxStatusHandler   TransactionStatus godoc
// // @Summary      Transaction status
// // @Tags         Transaction
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,near,aptos,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        txHash   query      string  true  "txHash"
// // @Router       /utils/txStatus [get]
// func (h *handler) TxStatusHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	var logs = make([]*pb.Log, 0)
// 	chainGroup := ctx.Query("chain")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		TxStatusReq := &pb.TxStatusRequest{
// 			Chain:  ctx.Query("chain"),
// 			TxHash: ctx.Query("txHash"),
// 		}
// 		resp, err := grpcClient.TxStatus(context2, TxStatusReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			if resp != nil {
// 				//Handle pending/unknown case
// 				if len(resp.Logs) == 0 {
// 					resp.Logs = logs
// 				}
// 				utils.APIResponse(ctx, s.Message(), codes.OK, http.MethodGet, resp)
// 				return
// 			} else {
// 				//Handle grpc errors
// 				utils.APIResponse(ctx, defaultSuccessMsg, codes.OK, http.MethodGet, nil)
// 				return
// 			}
// 		}
// 		if len(resp.Logs) == 0 {
// 			resp.Logs = logs
// 		}
// 		utils.APIResponse(ctx, "Transaction status is fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // GasLimitHandler  GasLimit godoc
// // @Tags         Transaction
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        from   query      string  false  "from"
// // @Param		 to query  string true "to"
// // @Param		 data query  string false "data"
// // @Param        value query string true "value"
// // @Router       /utils/getGasLimit [get]
// func (h *handler) GasLimitHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	chainGroup := ctx.Query("chain")
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	flag, fromAddress, err := h.util.ValidateAddress(ctx.Query("from"), chainGroup, ctx.Query("chain"))
// 	if !flag || err != nil {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	flag, toAddress, err := h.util.ValidateAddress(ctx.Query("to"), chainGroup, ctx.Query("chain"))
// 	if !flag || err != nil {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		gasLimitReq := &pb.GasLimitRequest{
// 			Chain:    ctx.Query("chain"),
// 			From:     fromAddress,
// 			To:       toAddress,
// 			Gas:      0,
// 			GasPrice: 0,
// 			Value:    ctx.GetInt64("value"),
// 			Data:     ctx.Query("data"),
// 		}
// 		resp, err := grpcClient.GasLimit(context2, gasLimitReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Gas Limit data fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }

// // NonceHandler Nonce godoc
// // @Tags         Nonce
// // @Accept       json
// // @Produce      json
// // @Param        chain   query      string  true  "chain"  Enums(solana,near,,arbitrum,fantom,ethereum,harmony,polygon,celo,optimism,xinfin,metis,avalanche,aurora,bsc,boba,heco,gnosis,moonriver,moonbeam,klaytn,bttc,iotex,tomochain,evmos,fuse,cronos,astar,zksync)
// // @Param        address   query      string  true  "address"
// // @Router       /utils/nonce [get]
// func (h *handler) NonceHandler(ctx *gin.Context) {
// 	context2, cancel := context.WithTimeout(context.Background(), defaultTimeOutMills)
// 	defer cancel()
// 	address := ctx.Query("address")
// 	chain := ctx.Query("chain")
// 	chainGroup := chain
// 	if h.util.IsEVM(chainGroup) {
// 		chainGroup = utils.EVM
// 	} else if h.util.IsCosmos(chainGroup) {
// 		chainGroup = utils.COSMOS
// 	} else if h.util.IsNonEVM(chainGroup) {
// 		chainGroup = "nonevm"
// 	}
// 	valid, validAddress, err := h.util.ValidateAddress(address, chainGroup, chain)
// 	if err != nil || valid == false {
// 		statusCode, _ := status.FromError(err)
// 		utils.APIResponse(ctx, statusCode.Message(), statusCode.Code(), http.MethodGet, nil)
// 		return
// 	}
// 	if grpcClient, ok := h.grpcClient[chainGroup]; ok {
// 		nonceReq := &pb.NonceRequest{
// 			Address: validAddress,
// 			Chain:   chain,
// 		}
// 		resp, err := grpcClient.Nonce(context2, nonceReq)
// 		if err != nil {
// 			s, _ := status.FromError(err)
// 			utils.APIResponse(ctx, s.Message(), s.Code(), http.MethodGet, nil)
// 			return
// 		}
// 		utils.APIResponse(ctx, "Nonce data fetched successfully", codes.OK, http.MethodGet, resp)
// 		return
// 	} else {
// 		utils.APIResponse(ctx, "Chain is not supported", codes.Unavailable, http.MethodGet, nil)
// 		return
// 	}
// }
