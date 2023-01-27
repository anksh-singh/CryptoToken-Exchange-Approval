package openApi

import (
	bridge_allowance "bridge-allowance"
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type IOpenAPI interface {
	GetGasPriceInfo(chainId string) (*GasPriceStruct, error)
	GetPosData(address string) (*PositionalDataModel, error)
}

type OpenAPI struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewOpenAPI(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *OpenAPI {
	return &OpenAPI{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

func (o *OpenAPI) GetGasPriceInfo(chainName string) (*GasPriceStruct, error) {
	chainId, err := bridge_allowance.GetDebankId(*o.env, chainName)
	if err != nil {
		o.logger.Error("Chain not supported")
		return nil, errors.New("Chain not supported")
	}
	url := fmt.Sprintf("%s%s%s", o.env.OpenAPI.EndPoint, "/wallet/gas_market?chain_id=", chainId)
	body, err := o.httpRequest.GetRequestWithHeaders(url, "AccessKey", o.env.OpenAPI.APIKey)
	if err != nil {
		o.logger.Errorf("Gas Price info Logging Error  is : %v", err.Error())
		return &GasPriceStruct{}, err
	}
	var jsonResponseStruct []*GasPriceOpenAPIStruct
	err = json.Unmarshal(body, &jsonResponseStruct)
	if err != nil {
		o.logger.Errorf(" Gas Price info Logging Error  is : %v", err.Error())
		return &GasPriceStruct{}, err
	}

	var response GasPriceStruct
	if len(jsonResponseStruct) > 0 {
		for _, item := range jsonResponseStruct {
			if item.Level == "slow" {
				response.Slow = item.Price / utils.GWei
			} else if item.Level == "normal" {
				response.Normal = item.Price / utils.GWei
			} else if item.Level == "fast" {
				response.Fast = item.Price / utils.GWei
			}
		}
	}
	return &response, err
}
