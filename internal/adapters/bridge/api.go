package bridge

import (
	"bridge-allowance/config"
	app "bridge-allowance/internal/adapters/bridge/application"
	"bridge-allowance/internal/adapters/bridge/application/core/bridge"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"net"
)

type bridgeServer struct {
	config        *config.Config
	logger        *zap.SugaredLogger
	serverHandler *app.BridgeServerHandler
	services      common.Services
}

func NewServer(config *config.Config, logger *zap.SugaredLogger, services common.Services) *bridgeServer {
	bridgeService := bridge.NewBridge(config, logger, services)
	serverHandler := app.NewBridgeServerHandler(config, logger, bridgeService)
	return &bridgeServer{
		config:        config,
		logger:        logger,
		serverHandler: serverHandler,
		services:      services,
	}
}
func (bridge *bridgeServer) Start() {
	port := bridge.config.Bridge.BridgeCfg.ServerPort
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		bridge.logger.Errorf("failed to listen: %v", err)
	}
	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(bridge.config.BRIDGE_DATADOG_SERVICE))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(bridge.config.BRIDGE_DATADOG_SERVICE))

	bridge.logger.Infof("bridge gRPC Server Started at %v", port)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterUnifrontServer(grpcServer, bridge.serverHandler)
	if err = grpcServer.Serve(lis); err != nil {
		bridge.logger.Errorf("Failed to serve gRPC server over port : %v", err)
	}
}

var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Unifront Framework-bridge",
	Long:  `Unifront Framework:  Bridge Server`,
	// @BasePath  /v2
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig("", "")
		logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.Bridge.BridgeCfg.LogFile, conf.LOG_ENCODING_FORMAT)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		logger.Info(conf)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		httRequest := utils.NewHttpRequest(logger)
		helper := utils.Helpers{}
		services := common.Services{
			Http:      httRequest,
			Helper:    &helper,
		}
		bridgeServer := NewServer(conf, logger, services)
		bridgeServer.Start()
	},
}
