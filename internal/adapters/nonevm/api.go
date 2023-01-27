package nonevm

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/nonevm/application"
	"bridge-allowance/pkg/coingecko"
	aptos2 "bridge-allowance/pkg/customchain/aptos"
	near2 "bridge-allowance/pkg/customchain/near"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"net"
)

type NonEVMServer struct {
	config        *config.Config
	logger        *zap.SugaredLogger
	ServerHandler *application.NonEVMServerHandler
	coingecko     *coingecko.CoinGecko
}

func NewNonEVMServer(config *config.Config, log *zap.SugaredLogger, coinGecko *coingecko.CoinGecko) *NonEVMServer {
	near := near2.NewServiceNear(config, log, coinGecko)
	aptos := aptos2.NewServiceAptos(config, log, coinGecko)
	newServerHander := application.NewEVMServerServerHandler(config, log, coinGecko, near, aptos)
	return &NonEVMServer{
		config:        config,
		logger:        log,
		ServerHandler: newServerHander,
		coingecko:     coinGecko,
	}
}

func (s *NonEVMServer) Start() {
	port := s.config.NonEVMConfig.ServerPort
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		s.logger.Errorf("failed to listen:" + err.Error())
	}
	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(s.config.NONEVM_DATADOG_SERVICE))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(s.config.NONEVM_DATADOG_SERVICE))

	s.logger.Infof("Non EVM gRPC Server Started at %v", port)
	// Initialize the grpc server as normal, using the tracing interceptor.
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterUnifrontServer(grpcServer, s.ServerHandler)
	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Errorf("Failed to serve gRPC server over port : %v", err)
	}
}

var SolanaCmd = &cobra.Command{
	Use:   "nonevm",
	Short: "Unifront Framework-NonEVM",
	Long:  `Unifront Framework:  NonEVM Server`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig("", "")
		logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.NonEVMConfig.LogFile, conf.LOG_ENCODING_FORMAT)
		httpRequest := utils.NewHttpRequest(logger)
		coinGecko := coingecko.NewCoinGecko(conf, logger, httpRequest)
		solanaServer := NewNonEVMServer(conf, logger, coinGecko)
		solanaServer.Start()
	},
}
