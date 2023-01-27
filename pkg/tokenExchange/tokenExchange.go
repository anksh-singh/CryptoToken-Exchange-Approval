package tokenExchange

import (
	"bridge-allowance/pkg/grpc/proto/pb"
)

type ITokenExchangeForChain interface {
	GetExchangeTokens(chainUrl string) (*pb.ExchangeTokenResponse, error)
}
