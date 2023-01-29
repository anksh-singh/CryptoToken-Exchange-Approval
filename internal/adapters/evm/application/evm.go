package application

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/evm/application/core/rpc"

	// "bridge-allowance/internal/adapters/evm/application/core/rpc/swap"
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
	// evmSwap swap.IExchangeSwap
	utils *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewEVMServerHandler(config config.Config, log *zap.SugaredLogger, core rpc.EvmCore) *evmServerHandler {
	utilManager := utils.NewUtils(log, &config)
	handler := &evmServerHandler{
		config:  &config,
		logger:  log,
		evmCore: core,
		// evmSwap: evmSwap,
		utils: utilManager,
	}
	return handler
}


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
