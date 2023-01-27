package covalent

//
import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCovalentService(t *testing.T) {
	logger := utils.InitLogger()
	config := config.Config{}
	httpRequestMocker := new(utils.MockHttpRequest)

	t.Run("GetAssets", func(t *testing.T) {
		httpRequestMocker.On("GetRequest").Return([]byte{}, nil)
		covalentService := CovalentService{
			env:         config,
			logger:      logger,
			httpRequest: httpRequestMocker,
		}
		publicAddress := "0xaFf5D64aD24eb307e5EbeF3fCb789B5C21146f71"
		assets, _ := covalentService.GetAssets(publicAddress, config.Fantom.FantomUnmarshallChainId)
		assert.Nil(t, assets, "assets should be nil")
	})

	t.Run("List Transaction", func(t *testing.T) {
		httpRequestMocker.On("GetRequest").Return([]byte{}, nil)
		covalentService := CovalentService{
			env:         config,
			logger:      logger,
			httpRequest: httpRequestMocker,
		}
		publicAddress := "0xaFf5D64aD24eb307e5EbeF3fCb789B5C21146f71"
		listTransactionData, _ := covalentService.ListTransaction(&pb.ListTransactionRequest{Address: publicAddress}, config.Fantom.FantomUnmarshallChainId)
		assert.Nil(t, listTransactionData, "list transaction data should be nil")
	})
}
