package _x

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/stretchr/testify/mock"
)

type MockIOx struct {
	mock.Mock
}

func (m *MockIOx) GetExchangeQuote(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, coinGeckoChainId string, srcTokenDecimals string, dstTokenDecimals string, slippage string) (*pb.ExchangeQuoteResponse, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*pb.ExchangeQuoteResponse), args.Error(1)
}

func (m *MockIOx) GetExchangeSwap(SellToken string, BuyToken string, SellAmount string, zeroXChainUrl string, srcTokenDecimals string, slippage string) (*pb.ExchangeSwapResponse, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*pb.ExchangeSwapResponse), args.Error(1)
}
