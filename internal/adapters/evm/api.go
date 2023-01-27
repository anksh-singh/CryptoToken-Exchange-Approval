package evm

import (
	"bridge-allowance/config"
	app "bridge-allowance/internal/adapters/evm/application"
	"bridge-allowance/internal/adapters/evm/application/core/rpc"
	"bridge-allowance/internal/adapters/evm/application/core/rpc/swap"
	"bridge-allowance/internal/common"
	_x "bridge-allowance/pkg/0x"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/covalent"
	"bridge-allowance/pkg/cowswap"
	"bridge-allowance/pkg/customchain/zksync"
	"bridge-allowance/pkg/dodo"
	dzap2 "bridge-allowance/pkg/dzap"
	"bridge-allowance/pkg/grpc/proto/pb"
	lifi2 "bridge-allowance/pkg/lifi"
	"bridge-allowance/pkg/oneinch"
	"bridge-allowance/pkg/openApi"
	"bridge-allowance/pkg/tokenExchange/cocoswap"
	"bridge-allowance/pkg/tokenExchange/dodoEth"
	uniswap2 "bridge-allowance/pkg/tokenExchange/ubeswap"
	"bridge-allowance/pkg/trustwallet"
	"bridge-allowance/pkg/unmarshal"
	"bridge-allowance/pkg/zeroswap"
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
	newSwap := swap.NewSwap(config, logger, services)
	serverHandler := app.NewEVMServerHandler(*config, logger, core, newSwap)
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
	Short: "Unifront Framework-evm",
	Long:  `Unifront Framework:  EVM Server`,
	// @BasePath  /v2
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig("", "")
		logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.EVM.Cfg.LogFile, conf.LOG_ENCODING_FORMAT)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		logger.Info(conf)
		logger.Info(":::::::::::::::::::::::::::::::: Configuration ::::::::::::::::::::::::::::::::")
		httRequest := utils.NewHttpRequest(logger)
		coinGecko := coingecko.NewCoinGecko(conf, logger, httRequest)
		covalentService := covalent.NewCovalentService(conf, logger, httRequest)
		unmarshall := unmarshal.NewUnMarshalService(conf, logger, httRequest)
		debank := openApi.NewOpenAPI(conf, logger, httRequest)
		helper := utils.Helpers{}
		cocoSwap := cocoswap.NewTokenExchangeService(conf, logger, httRequest)
		uniswap := uniswap2.NewTokenExchangeService(conf, logger, httRequest)
		dodoex := dodoEth.NewTokenExchangeService(conf, logger, httRequest)
		dodoExCache := dodoEth.NewTokenExchangeServiceCache(conf, logger, httRequest)
		zeroX := _x.NewOXService(conf, logger, httRequest, &helper, coinGecko)
		dodoExsWap := dodo.NewServiceDodo(conf, logger, httRequest, &helper, coinGecko)
		trustWallet := trustwallet.NewTrustWallet(conf, logger, httRequest)
		utils := utils.NewUtils(logger, conf)
		lifi := lifi2.NewLiFiService(conf, logger, httRequest, &helper, utils.GetEVMBridgeData("lifi"), utils)
		oneInch := oneinch.NewOneInchService(conf, logger, httRequest, &helper, coinGecko)
		zeroSwap := zeroswap.NewZeroSwapService(conf, logger, httRequest, &helper, coinGecko)
		cowSwap := cowswap.NewCowSwapService(conf, logger, httRequest, &helper, coinGecko, lifi)
		zsynk := zksync.NewServiceZksync(conf, logger, coinGecko)
		dzap := dzap2.NewServiceDZap(conf, logger, httRequest, &helper, coinGecko)
		services := common.Services{
			Http:       httRequest,
			CoinGecko:  coinGecko,
			Covalent:   covalentService,
			Helper:     &helper,
			Unmarshall: unmarshall,
			//TODO: zerox swap swap
			//zeroX:      _x.OXService{},
			CocoSwapTokenExchange:    cocoSwap,
			UniSwapTokenExchange:     uniswap,
			ZeroX:                    zeroX,
			Debank:                   debank,
			DoDoExTokenExchange:      dodoex,
			DoDoExTokenExchangeCache: dodoExCache,
			DodoSwap:                 dodoExsWap,
			TrustWallet:              trustWallet,
			LiFi:                     lifi,
			OneInch:                  oneInch,
			ZeroSwap:                 zeroSwap,
			CowSwap:                  cowSwap,
			DZap:                     dzap,
			Zksync:                   zsynk,
		}
		evmServer := NewServer(conf, logger, services)
		evmServer.Start()
	},
}
