package evm

import (
	"bridge-allowance/config"
	app "bridge-allowance/internal/adapters/evm/application"
	"bridge-allowance/internal/adapters/evm/application/core/rpc"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"net"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type evmServer struct {
	config        config.Config
	logger        *zap.SugaredLogger
	serverHandler app.EVMServerHandler
	services      common.Services
}

func NewServer(config *config.Config, logger *zap.SugaredLogger, services common.Services) *evmServer {
	core := rpc.NewEVMCore(config, logger, services)
	serverHandler := app.NewEVMServerHandler(*config, logger, core)
	return &evmServer{
		config:        *config,
		logger:        logger,
		serverHandler: serverHandler,
		services:      services,
	}
}

func (evm *evmServer) Start() {
	port := evm.config.EVM.Cfg.ServerPort
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		evm.logger.Errorf("failed to listen: %v", err)
	}
	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(evm.config.EVM_DATADOG_SERVICE))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(evm.config.EVM_DATADOG_SERVICE))

	evm.logger.Infof("EVM gRPC Server Started at %v", port)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterUnifrontServer(grpcServer, evm.serverHandler)
	if err = grpcServer.Serve(lis); err != nil {
		evm.logger.Errorf("Failed to serve gRPC server over port : %v", err)
	}
}

var EvmCmd = &cobra.Command{
	Use:   "evm",
	Short: "Bridge Allowance-evm",
	Long:  `Bridge Allowance:  EVM Server`,
	// @BasePath  /v2
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig("", "")
		logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.EVM.Cfg.LogFile, conf.LOG_ENCODING_FORMAT)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		logger.Info(conf)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		httRequest := utils.NewHttpRequest(logger)
		helper := utils.Helpers{}
		services := common.Services{
			Http:       httRequest,
			Helper:     &helper,
			
		}
		evmServer := NewServer(conf, logger, services)
		evmServer.Start()
	},
}
