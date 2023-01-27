package rpc

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SolanaManager) GetTransactionFee() (*pb.ProcessingFeeResponse, error) {
	resp, err := s.client.GetFees(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		s.logger.Errorf("Logging Error for Getting Fee Response is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	spew.Dump(resp)
	s.logger.Infof("Logging Output for Getting Fee Response with response : %v", resp.Value)
	return &pb.ProcessingFeeResponse{
		Value: float64(resp.Value.FeeCalculator.LamportsPerSignature),
	}, nil
}
func (s *SolanaManager) GetTokenPrice(in *pb.TokenPriceRequest) (*pb.TokenPriceResponse, error) {
	return s.coingecko.GetTokenExchange(in.Currency, "solana")
}
