package ubeswap

// ubeswap and quickswap are same
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
	var tokenExchange ExchangeToken
	err = json.Unmarshal(body, &tokenExchange)
	if err != nil {
		t.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}

	var tokenExchangeAvalanche ExchangeToken
	// TODO:: Change the Logic for Avalanche specific TokenExchange two url concatenation
	if chainId == 43114 {
		bodyA, err := t.httpRequest.GetRequest("https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json")
		err = json.Unmarshal(bodyA, &tokenExchangeAvalanche)
		if err != nil {
			t.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
		}
	}
	responseStruct := pb.ExchangeTokenResponse{}
	if len(tokenExchange.Tokens) > 0 {
		tokenExchange.Tokens = append(tokenExchange.Tokens, tokenExchangeAvalanche.Tokens...)
		tokenExchange.Tokens = RemoveDuplicate(tokenExchange.Tokens)
		for _, item := range tokenExchange.Tokens {
			if item.ChainID == chainId {
				exchangeTokenInfo := pb.ExchangeTokenInfo{
					TokenAddress:  item.Address,
					TokenDecimals: fmt.Sprint(item.Decimals),
					TokenLogoUrl:  item.LogoURI,
					TokenSymbol:   item.Symbol,
					TokenName:     item.Name,
					LogoUrl:       exchangeLogoUrl,
				}
				//if exchangeTokenInfo.LogoUrl == "" {
				//	exchangeTokenInfo.LogoUrl = exchangeLogoUrl
				//}
				responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
			}

		}
	}

	// TODO: Need to find Better Approach
	// Hard Coded Front Token if Chain is Matic
	if chainId == 43114 {
		exchangeTokenInfo := pb.ExchangeTokenInfo{
			TokenAddress:  "0x2f86508f41310d8d974b76deb3d246c0caa71cf5",
			TokenDecimals: fmt.Sprint(18),
			TokenLogoUrl:  "https://assets.coingecko.com/coins/images/15706/large/Hotcross.png?1632197570",
			TokenSymbol:   "HOTCROSS",
			TokenName:     "Hot Cross",
			LogoUrl:       exchangeLogoUrl,
		}
		responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
	}
	if chainId == 137 {
		exchangeTokenInfo := pb.ExchangeTokenInfo{
			TokenAddress:  "0xa3ed22eee92a3872709823a6970069e12a4540eb",
			TokenDecimals: fmt.Sprint(18),
			TokenLogoUrl:  "https://assets.coingecko.com/coins/images/12479/large/frontier_logo.png?1600145472",
			TokenSymbol:   "FRONT",
			TokenName:     "Front Token",
			LogoUrl:       exchangeLogoUrl,
		}
		responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
	}
	return &responseStruct, err
}

func RemoveDuplicate(slice []*Tokens) []*Tokens {
	uniqMap := make(map[string]*Tokens)
	for _, v := range slice {
		uniqMap[v.Address] = v
	}
	uniqSlice := make([]*Tokens, 0, len(uniqMap))
	for _, v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
