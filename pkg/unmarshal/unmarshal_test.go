package unmarshal

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUnMarshalService(t *testing.T) {
	config := config.Config{}
	httpRequestMocker := new(utils.MockHttpRequest)

	t.Run("GetAssets", func(t *testing.T) {
		httpRequestMocker.On("GetRequest").Return([]byte{}, nil)
		unmarshall := UnmarshallService{
			env:         config,
			logger:      logger,
			httpRequest: httpRequestMocker,
		}
		// with invalid address
		publicAddress := "0x73D23cDaBBb25B0E039470ea940514Ca30744277"
		assets, _ := unmarshall.GetAssets(publicAddress, config.Fantom.FantomUnmarshallChainId)
		assert.Nil(t, assets, "assets should be nil")

	})

	t.Run("ListTransaction", func(t *testing.T) {
		httpRequestMocker.On("GetRequest").Return([]byte{}, nil)
		unmarshall := UnmarshallService{
			env:         config,
			logger:      logger,
			httpRequest: httpRequestMocker,
		}
		// With Valid Pub Address
		publicAddress := "0x73D23cDaBBb25B0E039470ea940514Ca30744277"
		listData, _ := unmarshall.ListTransaction(&pb.ListTransactionRequest{Address: publicAddress}, config.Fantom.FantomUnmarshallChainId)
		assert.Nil(t, listData, "list transaction data should be nil")
	})
}
