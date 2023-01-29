package application

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/bridge/application/core/bridge"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"context"
	"go.uber.org/zap"
)

type BridgeServerHandler struct {
	config *config.Config
	logger *zap.SugaredLogger
	bridge bridge.IBridge
	utils  *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewBridgeServerHandler(config *config.Config, logger *zap.SugaredLogger, bridge bridge.IBridge) *BridgeServerHandler {
	utilManager := utils.NewUtils(logger, config)
	return &BridgeServerHandler{
		config: config,
		logger: logger,
		utils:  utilManager,
		bridge: bridge,
	}
}

func (bridge *BridgeServerHandler) BridgeChain(ctx context.Context, request *pb.BridgeChainRequest) (*pb.BridgeChainResponse, error) {
	defer bridge.utils.CleanUp(bridge.logger)
	bridge.logger.Info("initiating request to get bridge chains")
	res, err := bridge.bridge.GetChains(request)
	if err != nil {
		bridge.logger.Errorf("Get bridge chains resulted in error: %v", err.Error())
		return nil, err
	}
	bridge.logger.Info("Return bridge chains")
	return res, err
}

func (bridge *BridgeServerHandler) BridgeChainTokens(ctx context.Context, request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	defer bridge.utils.CleanUp(bridge.logger)
	bridge.logger.Info("initiating request to get bridge chains tokens")
	res, err := bridge.bridge.GetChainTokens(request)
	if err != nil {
		bridge.logger.Errorf("Get bridge chains tokens resulted in error: %v", err.Error())
		return nil, err
	}
	bridge.logger.Info("Return bridge chains tokens")
	return res, err
}

func (bridge *BridgeServerHandler) BridgeQuote(ctx context.Context, request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error) {
	defer bridge.utils.CleanUp(bridge.logger)
	bridge.logger.Info("initiating request to get bridge quote")
	res, err := bridge.bridge.GetQuote(request)
	if err != nil {
		bridge.logger.Errorf("Get bridge quote resulted in error: %v", err.Error())
		return nil, err
	}
	bridge.logger.Info("Return bridge quote")
	return res, err
}

func (bridge *BridgeServerHandler) BridgeTransaction(ctx context.Context, request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error) {
	defer bridge.utils.CleanUp(bridge.logger)
	bridge.logger.Info("initiating request to get bridge transaction")
	res, err := bridge.bridge.GetTransaction(request)
	if err != nil {
		bridge.logger.Errorf("Get bridge transaction resulted in error: %v", err.Error())
		return nil, err
	}
	bridge.logger.Info("Return bridge transaction")
	return res, err
}

// func (bridge *BridgeServerHandler) BridgeTransactionStatus(ctx context.Context, request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error) {
// 	defer bridge.utils.CleanUp(bridge.logger)
// 	bridge.logger.Info("initiating request to get bridge transaction status")
// 	res, err := bridge.bridge.GetTransactionStatus(request)
// 	if err != nil {
// 		bridge.logger.Errorf("Get bridge transaction status resulted in error: %v", err.Error())
// 		return nil, err
// 	}
// 	bridge.logger.Info("Return bridge transaction status")
// 	return res, err
// }
