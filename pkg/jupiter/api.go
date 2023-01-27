package jupiter

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JupiterSwap struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewJupiterSwap(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *JupiterSwap {
	return &JupiterSwap{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

func (j *JupiterSwap) GetSolanaExchangeTokenList(url string, exchangeType string) (interface{}, error) {
	body, err := j.httpRequest.GetRequest(url)
	if err != nil {
		j.logger.Errorf(err.Error())
	}
	if exchangeType == "jupiter" {
		var tokenList TokenList
		err = json.Unmarshal(body, &tokenList)
		if err != nil {
			j.logger.Error(err)
			return nil, err
		}
		return &tokenList, err
	} else {
		var tokenList SolanaTokenList
		err = json.Unmarshal(body, &tokenList)
		if err != nil {
			j.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return &tokenList, err
	}
}

func (j *JupiterSwap) GetSolanaExchangeQuote(url string) (*ExchangeQuoteRes, error) {
	body, err := j.httpRequest.GetRequest(url)
	if err != nil {
		j.logger.Errorf(err.Error())
	}
	var quoteRes ExchangeQuoteRes
	err = json.Unmarshal(body, &quoteRes)
	if err != nil {
		j.logger.Error(err)
		return nil, err
	}
	if len(quoteRes.Data) == 0 {
		return nil, errors.New("No Routes available for this exchange quote")
	}
	return &quoteRes, err
}

func (j *JupiterSwap) GetSolanaExchangeSwap(url string, route ExchangeQuoteData, publicKey string) (*SwapResponse, error) {
	req := SwapRequest{
		Route:         route,
		UserPublicKey: publicKey,
		WrapUnwrapSOL: true,
	}
	jsonReq, _ := json.Marshal(req)
	reqBody := bytes.NewBuffer(jsonReq)
	res, err := j.httpRequest.PostRequest(url, reqBody)
	if err != nil {
		j.logger.Errorf(err.Error())
		return nil, err
	}
	var swapRes SwapResponse
	err = json.Unmarshal(res, &swapRes)
	if err != nil {
		j.logger.Error(err)
		return nil, err
	}
	return &swapRes, err
}
