package cocoswap

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenExchangeStruct struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewTokenExchangeService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *TokenExchangeStruct {
	return &TokenExchangeStruct{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

func (t *TokenExchangeStruct) GetExchangeTokens(chainUrl string, exchangeLogoUrl string, chainId int) (*pb.ExchangeTokenResponse, error) {
	body, err := t.httpRequest.GetRequest(chainUrl)
	var tokenExchange []ExchangeToken
	err = json.Unmarshal(body, &tokenExchange)
	if err != nil {
		t.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	responseStruct := pb.ExchangeTokenResponse{}
	for _, item := range tokenExchange {
		if item.ChainId == chainId {
			exchangeTokenInfo := pb.ExchangeTokenInfo{
				TokenAddress:  item.Address,
				TokenDecimals: fmt.Sprint(item.Decimals),
				TokenLogoUrl:  item.LogoURI,
				TokenSymbol:   item.Symbol,
				TokenName:     item.Name,
				LogoUrl:       exchangeLogoUrl,
			}
			responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
		}

	}
	return &responseStruct, err
}
