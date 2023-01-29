package application

import (
	"bridge-allowance/config"
	// "bridge-allowance/pkg/coingecko"
	// "bridge-allowance/pkg/customchain/aptos"
	// "bridge-allowance/pkg/customchain/near"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"

	"go.uber.org/zap"
)

type NonEVMServerHandler struct {
	config *config.Config
	logger *zap.SugaredLogger
	// solanaManager *rpc.SolanaManager
	// nearManager  near.INear
	// aptosManager aptos.IAptos
	// coingecko    *coingecko.CoinGecko
	utils        *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewEVMServerServerHandler(config *config.Config, log *zap.SugaredLogger) *NonEVMServerHandler {
	utilsManager := utils.NewUtils(log, config)
	s := &NonEVMServerHandler{
		config: config,
		logger: log,
		utils:        utilsManager,
	
	}
	return s
}
