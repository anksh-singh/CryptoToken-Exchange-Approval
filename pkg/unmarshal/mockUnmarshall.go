package unmarshal

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/stretchr/testify/mock"
)

type MockUnmarshal struct {
	mock.Mock
}

func (m *MockUnmarshal) GetAssets(address string, chainId string) ([]*UnmarshallAssetModel, error) {
	args := m.Called()
	result := args.Get(0)
	return result.([]*UnmarshallAssetModel), args.Error(1)
}

func (m *MockUnmarshal) ListTransaction(in *pb.ListTransactionRequest, chainId string) (*UnmarshallTransactionModel, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*UnmarshallTransactionModel), args.Error(1)
}
