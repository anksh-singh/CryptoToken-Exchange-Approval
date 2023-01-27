package application

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/evm/application/core/rpc"
	"bridge-allowance/internal/adapters/evm/application/core/rpc/swap"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type EVMServerHandler interface {
	pb.UnifrontServer
}

type evmServerHandler struct {
	config  *config.Config
	logger  *zap.SugaredLogger
	evmCore rpc.EvmCore
	evmSwap swap.IExchangeSwap
	utils   *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewEVMServerHandler(config config.Config, log *zap.SugaredLogger, core rpc.EvmCore, evmSwap swap.IExchangeSwap) *evmServerHandler {
	utilManager := utils.NewUtils(log, &config)
	handler := &evmServerHandler{
		config:  &config,
		logger:  log,
		evmCore: core,
		evmSwap: evmSwap,
		utils:   utilManager,
	}
	return handler
}

// func (evm *evmServerHandler) TokenPrice(ctx context.Context, request *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetTokenPrice(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error fetching Token Price", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching token price request")
// 	return res, nil
// }

// func (evm *evmServerHandler) TokenPriceV2(ctx context.Context, request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetTokenPriceV2(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error fetching Token Price V2", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching token price request V2")
// 	return res, nil
// }

// func (evm *evmServerHandler) Balance(ctx context.Context, request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetAssets(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error fetching assets", err)
// 		return nil, err
// 	}
// 	evm.logger.Debug("Handler: Fetching balances request")
// 	return res, nil
// }

// func (evm *evmServerHandler) Nonce(ctx context.Context, request *pb.NonceRequest) (*pb.NonceResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetNonce(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error fetching Nonce", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching Nonce request")
// 	return res, nil
// }

// func (evm *evmServerHandler) SendTransaction(ctx context.Context, request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.SendTransaction(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error while initiating send transaction")
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching SendTransaction request")
// 	return res, nil
// }

// func (evm *evmServerHandler) ListTransaction(ctx context.Context, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.ListTransaction(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error listing transactions history", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching transactions list")
// 	return res, nil
// }

// func (evm *evmServerHandler) GasLimit(ctx context.Context, request *pb.GasLimitRequest) (*pb.GasLimitResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GasLimit(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error listing transactions history", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching gas limit")
// 	return res, nil
// }

// func (evm *evmServerHandler) ProcessingFee(ctx context.Context, request *pb.ProcessingFeeRequest) (*pb.ProcessingFeeResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Processing Fee API")
// 	res, err := evm.evmCore.GetProcessingFee(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error requesting processingFee/gasPrice", err)
// 		return nil, err
// 	}
// 	evm.logger.Infof("Handler: Fetching ProcessingFee API")
// 	return res, err
// }

// func (evm *evmServerHandler) ExchangeTokens(ctx context.Context, in *pb.ExchangeTokenRequest) (*pb.ExchangeTokenResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Exchange Tokens")
// 	res, err := evm.evmSwap.GetExchangeTokens(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Token with Error %v: ", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Infof("Outputing Exchange Token API")
// 	return res, err
// }

// func (evm *evmServerHandler) ExchangeQuote(ctx context.Context, in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Exchange Quote")
// 	res, err := evm.evmSwap.GetExchangeQuote(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Quote with Error %v: ", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Infof("Outputing Exchange Quote API")
// 	return res, err
// }

// func (evm *evmServerHandler) ExchangeMultiQuote(ctx context.Context, in *pb.ExchangeMultiQuoteRequest) (*pb.ExchangeMultiQuoteResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Exchange Quote")
// 	res, err := evm.evmSwap.GetExchangeMultiQuote(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Quote with Error %v: ", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Infof("Outputing Exchange Quote API")
// 	return res, nil
// }

// func (evm *evmServerHandler) ExchangeSwap(ctx context.Context, in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Exchange Swap API")
// 	res, err := evm.evmSwap.GetExchangeSwap(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Swap with Error %v: ", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Infof("Outputing Exchange Swap API")
// 	return res, nil
// }

// func (evm *evmServerHandler) ExchangeMultiSwap(ctx context.Context, in *pb.ExchangeMultiSwapRequest) (*pb.ExchangeMultipleSwapResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Exchange Swap API")
// 	res, err := evm.evmSwap.GetExchangeMultiSwap(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Swap with Error %v: ", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Infof("Outputing Exchange Swap API")
// 	return res, nil
// }

// func (evm *evmServerHandler) FreeTradeCount(ctx context.Context, in *pb.FreeTradeCountRequest) (*pb.FreeTradeCountResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Free Trade count API")
// 	res, err := evm.evmSwap.GetFreeTradeCount(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request to get free trade count with Error %v: ", err.Error())
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) ExchangeTokenApprove(ctx context.Context, in *pb.TokenApprovalRequest) (*pb.TokenApprovalResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiating request for Token approval API")
// 	res, err := evm.evmSwap.GetTokenApproval(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request to get token approval with Error %v: ", err.Error())
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) ExchangeSwapSignature(ctx context.Context, in *pb.ExchangeSignatureRequest) (*pb.ExchangeSignatureResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiate swap signature request")
// 	res, err := evm.evmSwap.GetExchangeSignature(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Signature with Error %v: ", err.Error())
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) ExchangeSwapExecute(ctx context.Context, in *pb.ExchangeSwapExecuteRequest) (*pb.ExchangeSwapExecuteResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Infof("initiate " + in.ExchangeType + " request")
// 	res, err := evm.evmSwap.GetExchangeSwapExecute(in)
// 	if err != nil {
// 		evm.logger.Errorf("Request Exchange Swap with Error %v: ", err.Error())
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) TxStatus(ctx context.Context, request *pb.TxStatusRequest) (*pb.TxStatusResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Info("initiating request to get tx status")
// 	res, err := evm.evmCore.GetTxStatus(request)
// 	if err != nil {
// 		evm.logger.Errorf("Get transaction status resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	evm.logger.Info("Return tx status")
// 	return res, nil
// }

func (evm *evmServerHandler) Allowance(ctx context.Context, request *pb.AllowanceRequest) (*pb.AllowanceResponse, error) {
	defer evm.utils.CleanUp(evm.logger)
	fmt.Println("inside evmshandlerss")
	evm.logger.Info("initiating request to get token allowance")
	res, err := evm.evmCore.GetTokenAllowance(request)
	if err != nil {
		evm.logger.Errorf("Get Token Allowance resulted in error: %v", err.Error())
		return nil, err
	}
	return res, nil
}

// func (evm *evmServerHandler) Approve(ctx context.Context, request *pb.ApprovalRequest) (*pb.ApprovalResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Info("initiating request to  token approval")
// 	res, err := evm.evmCore.TokenApprove(request)
// 	if err != nil {
// 		evm.logger.Errorf("Token Approval resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) UserData(ctx context.Context, request *pb.UserDataRequest) (*pb.UserDataResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Info("initiating request to fetch user data")
// 	res, err := evm.evmCore.GetUserData(request)
// 	if err != nil {
// 		evm.logger.Errorf("Token Approval resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (evm *evmServerHandler) GetNftCollections(ctx context.Context, request *pb.NftCollectionRequest) (*pb.ListNftCollectionResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetNftCollections(request)
// 	if err != nil {
// 		evm.logger.Error("Handler: Error getting GetNftCollections", err)
// 		return nil, err
// 	}
// 	evm.logger.Info("Handler: Fetching GetNftCollections")
// 	return res, nil
// }

// func (evm *evmServerHandler) BulkApproval(ctx context.Context, request *pb.ApprovalRequest) (*pb.BulkApprovalResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Info("initiating request to  bulk token approval")
// 	res, err := evm.evmCore.BulkApproval(request)
// 	if err != nil {
// 		evm.logger.Errorf("Token Approval resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	return res, err
// }

// func (evm *evmServerHandler) BulkAllowance(ctx context.Context, request *pb.AllowanceRequest) (*pb.BulkAllowanceResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	evm.logger.Info("initiating request to bulk token allowance")
// 	res, err := evm.evmCore.BulkAllowance(request)
// 	if err != nil {
// 		evm.logger.Errorf("Get Token Allowance resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	return res, err
// }

// func (evm *evmServerHandler) GetOpportunites(ctx context.Context, request *pb.GetOpportunitiesRequest) (*pb.GetOpportunitesResponse, error) {
// 	defer evm.utils.CleanUp(evm.logger)
// 	res, err := evm.evmCore.GetOpportunites(request)
// 	if err != nil {
// 		evm.logger.Errorf("Get Opportunities resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	return res, err
// }
