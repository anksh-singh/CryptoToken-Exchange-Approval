package application

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/nonevm/application/solana/core/rpc"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/customchain/aptos"
	"bridge-allowance/pkg/customchain/near"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/pkg/jupiter"
	"bridge-allowance/utils"
	"context"
	r "github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NonEVMServerHandler struct {
	config        *config.Config
	logger        *zap.SugaredLogger
	solanaManager *rpc.SolanaManager
	nearManager   near.INear
	aptosManager  aptos.IAptos
	coingecko     *coingecko.CoinGecko
	utils         *utils.UtilConf
	pb.UnimplementedUnifrontServer
}

func NewEVMServerServerHandler(config *config.Config, log *zap.SugaredLogger, coingecko *coingecko.CoinGecko, nearManager near.INear, aptosManager aptos.IAptos) *NonEVMServerHandler {
	http := utils.NewHttpRequest(log)
	newJupiter := jupiter.NewJupiterSwap(config, log, http)
	newSolManager := rpc.NewSolanaManager(config, log, coingecko, newJupiter)
	utilsManager := utils.NewUtils(log, config)
	s := &NonEVMServerHandler{
		config:        config,
		logger:        log,
		solanaManager: newSolManager,
		coingecko:     coingecko,
		utils:         utilsManager,
		nearManager:   nearManager,
		aptosManager:  aptosManager,
	}
	return s
}

func (s *NonEVMServerHandler) TokenPrice(ctx context.Context, request *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		res, err := s.solanaManager.GetTokenPrice(request)
		return res, err
	} else if request.Chain == "near" {
		return s.nearManager.GetTokenPrice(request.Currency, "")
	} else if request.Chain == "aptos" {
		return s.aptosManager.GetTokenPrice(request.Currency, "")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) Balance(ctx context.Context, request *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		res, err := s.solanaManager.GetBalance(request)
		return res, err
	} else if request.Chain == "near" {
		return s.nearManager.GetAssets(request)
	} else if request.Chain == "aptos" {
		return s.aptosManager.GetAssets(request)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ProcessingFee(ctx context.Context, request *pb.ProcessingFeeRequest) (*pb.ProcessingFeeResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		if request.GetTransactionFee() == true {
			res, err := s.solanaManager.GetTransactionFee()
			return res, err
		}
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "Solana supports transaction fee only")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) Nonce(ctx context.Context, request *pb.NonceRequest) (*pb.NonceResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "Nonce is not supported in solana")
	} else if request.Chain == "near" {
		return s.nearManager.GetNonce(request)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}
}

func (s *NonEVMServerHandler) SendTransaction(ctx context.Context, request *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		res, err := s.solanaManager.SendTransaction(request)
		return res, err
	} else if request.Chain == "near" {
		return s.nearManager.SendTransaction(request)
	} else if request.Chain == "aptos" {
		return s.aptosManager.SendTransaction(request)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ListTransaction(ctx context.Context, request *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if request.Chain == "solana" {
		res, err := s.solanaManager.ListTransaction(request)
		return res, err
	} else if request.Chain == "near" {
		return s.nearManager.ListTransaction(request)
		//} else if request.Chain == "aptos" {
		//	return s.aptosManager.ListTransaction(request)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ExchangeTokens(ctx context.Context, in *pb.ExchangeTokenRequest) (*pb.ExchangeTokenResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("initiating request for Exchange Tokens")
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetExchangeTokens(in)
		if err != nil {
			s.logger.Errorf("Request Exchange Token with Error %v: ", err.Error())
		}
		s.logger.Infof("Outputing Exchange Token API")
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ExchangeQuote(ctx context.Context, in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("initiating request for Exchange Quote")
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetExchangeQuote(in)
		if err != nil {
			s.logger.Errorf("Request Exchange Quote with Error %v: ", err.Error())
		}
		s.logger.Infof("Outputing Exchange Quote API")
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ExchangeSwap(ctx context.Context, in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("initiating request for Exchange Swap API %v", in)
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetExchangeSwap(in)
		if err != nil {
			s.logger.Errorf("Request Exchange Swap with Error %v: ", err.Error())
		}
		s.logger.Infof("Outputing Exchange Swap API %v", res)
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) TxStatus(ctx context.Context, in *pb.TxStatusRequest) (*pb.TxStatusResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("TxStatus Request %v", in)
	if in.Chain == "solana" {
		res, err := s.solanaManager.TransactionStatus(in)
		if err != nil {
			s.logger.Errorf("Error %v: ", err.Error())
		}
		s.logger.Infof("TxStatus Response %v", res)
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else if in.Chain == "aptos" {
		return s.aptosManager.GetTxStatus(in)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) GetRecentBlockHash() (*r.GetRecentBlockhashResult, error) {
	defer s.utils.CleanUp(s.logger)
	return s.solanaManager.GetRecentBlockHash()
}

func (s *NonEVMServerHandler) ExchangeMultiQuote(ctx context.Context, in *pb.ExchangeMultiQuoteRequest) (*pb.ExchangeMultiQuoteResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("initiating request for Exchange Quote")
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetExchangeMultiQuote(in)
		if err != nil {
			s.logger.Errorf("Request Exchange Quote with Error %v: ", err.Error())
			return nil, err
		}
		s.logger.Infof("Outputing Exchange Quote API")
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) ExchangeMultiSwap(ctx context.Context, in *pb.ExchangeMultiSwapRequest) (*pb.ExchangeMultipleSwapResponse, error) {
	defer s.utils.CleanUp(s.logger)
	s.logger.Infof("initiating request for Exchange Swap API")
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetExchangeMultiSwap(in)
		if err != nil {
			s.logger.Errorf("Request Exchange Swap with Error %v: ", err.Error())
			return nil, err
		}
		s.logger.Infof("Outputing Exchange Swap API")
		return res, err
	} else if in.Chain == "near" {
		return nil, status.Errorf(codes.Unavailable, "Unavailable", "supported in near")
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}

}

func (s *NonEVMServerHandler) GetPositions(ctx context.Context, in *pb.PositionChainData) (*pb.GetPositionsResponse, error) {
	defer s.utils.CleanUp(s.logger)
	if in.Chain == "solana" {
		res, err := s.solanaManager.GetPositions(ctx, in)
		if err != nil {
			s.logger.Errorf("Request Get Positions with Error %v: ", err.Error())
			return nil, err
		}
		s.logger.Infof("Output : Get Positions API")
		return res, err
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Chain", "Invalid Chain")
	}
}
