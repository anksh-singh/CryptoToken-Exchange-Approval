package bridge

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/bridge/application/core/bridge/v1Proxy"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBridge interface {
	GetChains(request *pb.BridgeChainRequest) (*pb.BridgeChainResponse, error)
	GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error)
	GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error)
	GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error)
	GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error)
}

type Bridge struct {
	config   *config.Config
	logger   *zap.SugaredLogger
	services common.Services
	util     *utils.UtilConf
	helper   *utils.Helpers
	IV1Proxy v1Proxy.IV1Proxy
}

func NewBridge(config *config.Config, logger *zap.SugaredLogger, services common.Services) *Bridge {
	logger.Info("Supported EVM chains:", len(config.EVM.Cfg.Wallets))
	helper := utils.Helpers{}
	httpRequest := utils.NewHttpRequest(logger)
	newUtils := utils.NewUtils(logger, config)
	v1proxy := v1Proxy.NewV1Proxy(httpRequest)
	return &Bridge{
		config:   config,
		logger:   logger,
		services: services,
		util:     newUtils,
		helper:   &helper,
		IV1Proxy: v1proxy,
	}
}

func (b *Bridge) GetChains(request *pb.BridgeChainRequest) (*pb.BridgeChainResponse, error) {
	//source := b.util.GetWalletSource(request.Chain)
	b.logger.Infof("Bridge Chains: %v", b.config.Bridge.BridgeCfg.BridgeExchangeChains)
	if request.BridgeProvider == "lifi" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "liFiBridge"
	} else if request.BridgeProvider == "xy" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "XYBridge"
	} else if request.BridgeProvider == "socket" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "SocketBridge"
	} else if request.BridgeProvider == "router" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "RouterBridge"
	} else if request.BridgeProvider == "rango" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "RangoBridge"
	} else if request.BridgeProvider == "debridge" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChains = "DeBridge"
	} else {
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported '%s' Bridge", request.BridgeProvider), "Unsupported")
	}

	switch b.config.Bridge.BridgeCfg.BridgeExchangeChains {
	case "liFiBridge":
		return b.services.LiFi.GetChains()
	case "XYBridge":
		return b.IV1Proxy.GetChains(request, b.config.PROXIES_ENDPOINT)
	case "RouterBridge":
		return b.IV1Proxy.GetChains(request, b.config.PROXIES_ENDPOINT)
	case "RangoBridge":
		return b.IV1Proxy.GetChains(request, b.config.PROXIES_ENDPOINT)
	case "SocketBridge":
		return b.services.Socket.GetChains()
	case "DeBridge":
		return b.services.DeBridge.GetChains()
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", b.config.Bridge.BridgeCfg.BridgeExchangeChains), "Unsupported operation")
	}
}

func (b *Bridge) GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	if request.BridgeProvider == "lifi" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "liFiBridge"
	} else if request.BridgeProvider == "xy" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "XYBridge"
	} else if request.BridgeProvider == "socket" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "SocketBridge"
	} else if request.BridgeProvider == "router" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "RouterBridge"
	} else if request.BridgeProvider == "rango" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "RangoBridge"
	} else if request.BridgeProvider == "debridge" {
		b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens = "DeBridge"
	} else {
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported '%s' Bridge", request.BridgeProvider), "Unsupported")
	}
	switch b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens {
	case "liFiBridge":
		return b.services.LiFi.GetChainTokens(request)
	case "SocketBridge":
		return b.services.Socket.GetChainTokens(request)
	case "XYBridge":
		return b.IV1Proxy.GetChainTokens(request, b.config.PROXIES_ENDPOINT)
	case "RouterBridge":
		return b.IV1Proxy.GetChainTokens(request, b.config.PROXIES_ENDPOINT)
	case "RangoBridge":
		return b.IV1Proxy.GetChainTokens(request, b.config.PROXIES_ENDPOINT)
	case "DeBridge":
		return b.services.DeBridge.GetChainTokens(request)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", b.config.Bridge.BridgeCfg.BridgeExchangeChainTokens), "Unsupported operation")
	}
}

func (b *Bridge) GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error) {
	if request.BridgeProvider == "lifi" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "liFiBridge"
	} else if request.BridgeProvider == "xy" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "XYBridge"
	} else if request.BridgeProvider == "socket" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "SocketBridge"
	} else if request.BridgeProvider == "router" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "RouterBridge"
	} else if request.BridgeProvider == "rango" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "RangoBridge"
	} else if request.BridgeProvider == "debridge" {
		b.config.Bridge.BridgeCfg.BridgeExchangeQuote = "DeBridge"
	} else {
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported '%s' Bridge", request.BridgeProvider), "Unsupported")
	}
	switch b.config.Bridge.BridgeCfg.BridgeExchangeQuote {
	case "liFiBridge":
		return b.services.LiFi.GetQuote(request)
	case "SocketBridge":
		return b.services.Socket.GetQuote(request)
	case "XYBridge":
		return b.IV1Proxy.GetQuote(request, b.config.PROXIES_ENDPOINT)
	case "RangoBridge":
		return b.IV1Proxy.GetQuote(request, b.config.PROXIES_ENDPOINT)
	case "RouterBridge":
		return b.IV1Proxy.GetQuote(request, b.config.PROXIES_ENDPOINT)
	case "DeBridge":
		return b.services.DeBridge.GetQuote(request)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", b.config.Bridge.BridgeCfg.BridgeExchangeQuote), "Unsupported operation")
	}
}

func (b *Bridge) GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error) {
	if request.BridgeProvider == "lifi" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "liFiBridge"
	} else if request.BridgeProvider == "xy" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "XYBridge"
	} else if request.BridgeProvider == "socket" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "SocketBridge"
	} else if request.BridgeProvider == "router" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "RouterBridge"
	} else if request.BridgeProvider == "rango" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "RangoBridge"
	} else if request.BridgeProvider == "debridge" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransaction = "DeBridge"
	} else {
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported '%s' Bridge", request.BridgeProvider), "Unsupported")
	}
	switch b.config.Bridge.BridgeCfg.BridgeExchangeTransaction {
	case "liFiBridge":
		return b.services.LiFi.GetTransaction(request)
	case "XYBridge":
		return b.IV1Proxy.GetTransaction(request, b.config.PROXIES_ENDPOINT)
	case "SocketBridge":
		return b.services.Socket.GetTransaction(request)
	case "RangoBridge":
		return b.IV1Proxy.GetTransaction(request, b.config.PROXIES_ENDPOINT)
	case "RouterBridge":
		return b.IV1Proxy.GetTransaction(request, b.config.PROXIES_ENDPOINT)
	case "DeBridge":
		return b.services.DeBridge.GetTransaction(request)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", b.config.Bridge.BridgeCfg.BridgeExchangeTransaction), "Unsupported operation")
	}
}

func (b *Bridge) GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error) {
	if request.BridgeProvider == "lifi" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus = "liFiBridge"
	} else if request.BridgeProvider == "xy" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus = "XYBridge"
	} else if request.BridgeProvider == "router" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus = "RouterBridge"
	} else if request.BridgeProvider == "socket" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus = "SocketBridge"
	} else if request.BridgeProvider == "rango" {
		b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus = "RangoBridge"
	} else {
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported '%s' Bridge", request.BridgeProvider), "Unsupported")
	}
	switch b.config.Bridge.BridgeCfg.BridgeExchangeTransactionStatus {
	case "XYBridge":
		return b.IV1Proxy.GetTransactionStatus(request, b.config.PROXIES_ENDPOINT)
	case "RouterBridge":
		return b.IV1Proxy.GetTransactionStatus(request, b.config.PROXIES_ENDPOINT)
	case "SocketBridge":
		return b.services.Socket.GetTransactionStatus(request)
	case "RangoBridge":
		return b.IV1Proxy.GetTransactionStatus(request, b.config.PROXIES_ENDPOINT)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", b.config.Bridge.BridgeCfg.BridgeExchangeTransaction), "Unsupported operation")
	}
}
