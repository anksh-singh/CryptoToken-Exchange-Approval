package tokenExchange

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/stretchr/testify/mock"
)

type MockTokenExchange struct {
	mock.Mock
}

func (m *MockTokenExchange) GetExchangeTokens(chainUrl string) (*pb.ExchangeTokenResponse, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*pb.ExchangeTokenResponse), args.Error(1)
}
