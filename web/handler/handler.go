package handler

import (
	"bridge-allowance/config"
	// cosmos "bridge-allowance/internal/adapters/cosmos/application"
	// "bridge-allowance/internal/adapters/nonevm/application"
	grpcClient "bridge-allowance/pkg/grpc/client"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"go.uber.org/zap"
	"time"
)

type handler struct {
	grpcClient    map[string]pb.UnifrontClient
	config        *config.Config
	logger        *zap.SugaredLogger
	util          *utils.UtilConf
	// nonEVMHandler *application.NonEVMServerHandler
	// cosmosHandler *cosmos.CosmosServerHandler
}

const (
	//defaultTimeOutMills API time out
	defaultTimeOutMills = 10 * time.Second
	// TimeOut25 TimeOut25Seconds Special case for grpc nonevm sendTransaction retry attempts
	TimeOut25 = 25 * time.Second
	//TimeOutOneMinute Special case for bridge APIs
	TimeOutOneMinute = 60 * time.Second
	//defaultSuccessMsg A default API success response
	defaultSuccessMsg = "OK"

	TxTimeOut = 40 * time.Second
)

func NewHandler(config *config.Config, logger *zap.SugaredLogger, gc *grpcClient.GrpcClientManager,
	) *handler {
	utilConf := utils.NewUtils(logger, config)
	// nonEVMHandler := application.NewEVMServerServerHandler(config, logger)
	// cosmosHandler, _ := cosmos.NewCosmosServerHandler(config, logger)

	return &handler{
		grpcClient:    gc.MapGrpcClient(),
		config:        config,
		logger:        logger,
		util:          utilConf,
		// nonEVMHandler: nonEVMHandler,
		// cosmosHandler: cosmosHandler,
	
	}
}
