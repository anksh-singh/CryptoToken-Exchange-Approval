package openApi

import (
	"github.com/stretchr/testify/mock"
)

type MockOpenAPI struct {
	mock.Mock
}

func (m *MockOpenAPI) GetGasPriceInfo(chainId string) (*GasPriceStruct, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*GasPriceStruct), args.Error(1)
}

func (m *MockOpenAPI) GetPosData(address string) (*PositionalDataModel, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*PositionalDataModel), args.Error(1)
}
