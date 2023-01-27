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
)

type TokenExchangeStruct struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
	helper      *utils.Helpers
}

func NewTokenExchangeService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *TokenExchangeStruct {
	helper := utils.Helpers{}
	return &TokenExchangeStruct{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
		helper:      &helper,
	}
}

func (t *TokenExchangeStruct) GetExchangeTokens(chainUrl string, exchangeLogoUrl string) (*pb.ExchangeTokenResponse, error) {
	body, err := t.httpRequest.GetRequest(chainUrl)
	var tokenExchange ExchangeToken
	err = json.Unmarshal(body, &tokenExchange)
	if err != nil {
		t.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	responseStruct := pb.ExchangeTokenResponse{}
	reflectTokenExchange := reflect.ValueOf(tokenExchange)
	for i := 0; i < reflectTokenExchange.NumField(); i++ {
		var tokenValue pb.ExchangeTokenInfo
		tokenValue.TokenSymbol = reflectTokenExchange.Field(i).FieldByName("Symbol").String()
		tokenValue.TokenName = reflectTokenExchange.Field(i).FieldByName("Name").String()
		tokenValue.TokenLogoUrl = "https://assets.unmarshal.io/tokens/" + t.helper.CheckSumAddress(reflectTokenExchange.Field(i).FieldByName("Address").String()) + ".png"
		tokenValue.TokenDecimals = fmt.Sprint(reflectTokenExchange.Field(i).FieldByName("Decimals").Int())
		tokenValue.TokenAddress = reflectTokenExchange.Field(i).FieldByName("Address").String()
		tokenValue.LogoUrl = exchangeLogoUrl
		responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &tokenValue)
	}
	return &responseStruct, err
}
