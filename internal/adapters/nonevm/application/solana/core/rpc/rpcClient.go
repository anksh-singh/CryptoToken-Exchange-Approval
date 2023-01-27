package rpc

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/jupiter"
	"bridge-allowance/utils"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
)

type SolanaManager struct {
	client      *rpc.Client
	env         *config.Config
	logger      *zap.SugaredLogger
	coingecko   *coingecko.CoinGecko
	solSwap     *jupiter.JupiterSwap
	helper      *utils.Helpers
	util        *utils.UtilConf
	httpRequest utils.IHttpRequest
}

func NewSolanaManager(config *config.Config, logger *zap.SugaredLogger, coingecko *coingecko.CoinGecko,
	solSwap *jupiter.JupiterSwap) *SolanaManager {
	endpoint := config.NONEVM_GRPC_ENDPOINT
	var solanaClusEndpoint string
	if endpoint == "dev" {
		solanaClusEndpoint = rpc.DevNet_RPC
	} else if endpoint == "test" {
		solanaClusEndpoint = rpc.TestNet_RPC
	} else {
		solanaClusEndpoint = rpc.MainNetBeta_RPC
	}
	utilsManager := utils.NewUtils(logger, config)
	logger.Infof(solanaClusEndpoint)

	solClient := rpc.New(solanaClusEndpoint)
	httpRequest := utils.NewHttpRequest(logger)
	utilsHelper := &utils.Helpers{}
	return &SolanaManager{
		client:      solClient,
		env:         config,
		logger:      logger,
		coingecko:   coingecko,
		solSwap:     solSwap,
		helper:      utilsHelper,
		httpRequest: httpRequest,
		util:        utilsManager,
	}
}
