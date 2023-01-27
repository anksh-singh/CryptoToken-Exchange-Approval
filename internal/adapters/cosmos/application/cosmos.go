package application

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/cosmos/application/core"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"context"
	"go.uber.org/zap"
)

type CosmosServerHandler struct {
	config     *config.Config
	logger     *zap.SugaredLogger
	coingecko  *coingecko.CoinGecko
	cosmosCore *core.Handler
	utils      *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewCosmosServerHandler(config *config.Config, log *zap.SugaredLogger, coingecko *coingecko.CoinGecko) (*CosmosServerHandler, error) {
	httRequest := utils.NewHttpRequest(log)
	utilsManager := utils.NewUtils(log, config)
	cosmosHandler, err := core.NewCosmosHandler(config, log, httRequest, utilsManager, coingecko)
	if err != nil {
		log.Errorf("Error in CosmosHandler %v", err)
		return nil, err
	}
	return &CosmosServerHandler{
		config:     config,
		logger:     log,
		coingecko:  coingecko,
		cosmosCore: cosmosHandler,
		utils:      utilsManager,
	}, nil
}

func (c *CosmosServerHandler) TokenPrice(ctx context.Context, request *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.coingecko.GetTokenExchange(request.Currency, request.Chain)
}
func (c *CosmosServerHandler) ListTransaction(ctx context.Context, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.ListTransaction(request)
}
func (c *CosmosServerHandler) GetValidators(ctx context.Context, request *pb.CosmosValidatorsRequest) (*pb.CosmosValidatorsResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetValidators(request)
}
func (c *CosmosServerHandler) GetCosmosAprRates(ctx context.Context, request *pb.CosmosAprRatesRequest) (*pb.CosmosAprRatesResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetCosmosAprRates(request)
}
func (c *CosmosServerHandler) GetDelegations(ctx context.Context, request *pb.CosmosDelegationsRequest) (*pb.CosmosDelegationsResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetDelegations(request)
}

func (c *CosmosServerHandler) GetCosmosCDPParams(ctx context.Context, request *pb.CosmosCDPParametersRequest) (*pb.CosmosCDPParametersResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetCosmosCDPParams(request)
}
func (c *CosmosServerHandler) TokenPriceV2(ctx context.Context, request *pb.TokenPriceRequest) (*pb.TokenPriceResponseV2, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	res, err := c.coingecko.GetTokenExchangeV2(request.Currency, request.Chain)
	if err != nil {
		c.logger.Error("Handler: Error fetching Token Price V2", err)
		return nil, err
	}
	c.logger.Info("Handler: Fetching token price request V2")
	return res, nil
}

func (c *CosmosServerHandler) CosmosAssets(ctx context.Context, request *pb.BalanceRequest) (*pb.CosmosAssetResponse, error) {
	c.logger.Info(request)
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetBalance(request)
}

func (c *CosmosServerHandler) CosmosSendTx(ctx context.Context, request *pb.CosmosSendTxRequest) (*pb.SendTransactionResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.SendTx(request)
}

func (c *CosmosServerHandler) CosmosSimulateTX(ctx context.Context, request *pb.CosmosSimulateTxRequest) (*pb.CosmosSimulateTxResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.CosmosSimulateTx(request)
}

func (c *CosmosServerHandler) CosmosGetBlockHeight(ctx context.Context, request *pb.GetCosmosBlockHeightRequest) (*pb.GetCosmosBlockHeightResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetLatestBlockHeight(request)
}

func (c *CosmosServerHandler) GetOpportunites(ctx context.Context, request *pb.GetOpportunitiesRequest) (*pb.GetOpportunitesResponse, error) {
	defer c.utils.CleanUp(c.logger)
	request.Chain = c.utils.RenameEVMCosmosCompatibleChains(request.Chain)
	return c.cosmosCore.GetCosmosOpportunities(request)
}
