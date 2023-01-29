package cosmos

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/cosmos/application"
	// "bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"net"
)

type cosmos struct {
	config              *config.Config
	logger              *zap.SugaredLogger
	// coingecko           *coingecko.CoinGecko
	cosmosServerHandler *application.CosmosServerHandler
}

func NewServer(config *config.Config, log *zap.SugaredLogger) (*cosmos, error) {
	newCosServerHander, err := application.NewCosmosServerHandler(config, log)
	if err != nil {
		return nil, err
	}
	return &cosmos{
		config:              config,
		logger:              log,
		cosmosServerHandler: newCosServerHander,
		// coingecko:           coingecko,
	}, nil
}

func (c *cosmos) Start() {
	port := c.config.Cosmos.Cfg.ServerPort
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		c.logger.Errorf("failed to listen:" + err.Error())
	}
	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(c.config.COSMOS_DATADOG_SERVICE))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(c.config.COSMOS_DATADOG_SERVICE))

	c.logger.Infof("Cosmos gRPC Server Started at %v", port)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterUnifrontServer(grpcServer, c.cosmosServerHandler)
	if err = grpcServer.Serve(lis); err != nil {
		c.logger.Errorf("Failed to serve gRPC server over port : %v", err)
	}
}

var CosmosCmd = &cobra.Command{
	Use:   "cosmos",
	Short: "Unifront Framework-cosmos",
	Long:  `Unifront Framework:  Cosmos Server`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig("", "")
		logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.Cosmos.Cfg.LogFile, conf.LOG_ENCODING_FORMAT)
		// httpRequest := utils.NewHttpRequest(logger)
		// coinGecko := coingecko.NewCoinGecko(conf, logger, httpRequest)
		cosmosServer, err := NewServer(conf, logger)
		if err != nil {
			logger.Error("Error in Starting Cosmos Server")
			return
		}
		cosmosServer.Start()
	},
}
