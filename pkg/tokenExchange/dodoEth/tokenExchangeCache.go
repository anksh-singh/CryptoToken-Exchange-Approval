package dodoEth

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"sync"
)

//TODO: Limitation Caching will be up until the adapter is restarted ,To be moved to different caching layer

type TokenExchangeStructCache struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
	helper      *utils.Helpers
}

func NewTokenExchangeServiceCache(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *TokenExchangeStructCache {
	helper := utils.Helpers{}
	return &TokenExchangeStructCache{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
		helper:      &helper,
	}
}

var once sync.Once

func (t *TokenExchangeStructCache) GetExchangeTokens(chainUrl string, exchangeLogoUrl string) (*pb.ExchangeTokenResponse, error) {
	var methodErr error
	var responseStruct pb.ExchangeTokenResponse
	body, err := t.httpRequest.GetRequest(chainUrl)
	var tokenExchange ExchangeToken
	err = json.Unmarshal(body, &tokenExchange)
	if err != nil {
		t.logger.Error(err)
		methodErr = err
	}
	reflectTokenExchange := reflect.ValueOf(tokenExchange)
	for i := 0; i < reflectTokenExchange.NumField(); i++ {
		var tokenValue pb.ExchangeTokenInfo
		if reflectTokenExchange.Field(i).FieldByName("Name").String() != "" {
			tokenValue.TokenSymbol = reflectTokenExchange.Field(i).FieldByName("Symbol").String()
			tokenValue.TokenName = reflectTokenExchange.Field(i).FieldByName("Name").String()
			tokenValue.TokenLogoUrl = "https://assets.unmarshal.io/tokens/" + t.helper.CheckSumAddress(reflectTokenExchange.Field(i).FieldByName("Address").String()) + ".png"
			tokenValue.TokenDecimals = fmt.Sprint(reflectTokenExchange.Field(i).FieldByName("Decimals").Int())
			tokenValue.TokenAddress = reflectTokenExchange.Field(i).FieldByName("Address").String()
			tokenValue.LogoUrl = exchangeLogoUrl
			responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &tokenValue)
		}
	}
	if methodErr != nil {
		t.logger.Error(methodErr)
		return nil, status.Errorf(codes.Internal, methodErr.Error())
	}
	return &responseStruct, methodErr
}
