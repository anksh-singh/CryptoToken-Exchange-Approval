package rpc

import (
	"bridge-allowance/internal/adapters/nonevm/application/solana/core"
	"bridge-allowance/pkg/grpc/proto/pb"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *SolanaManager) GetBalance(in *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	reqUrl := s.env.Unmarshall.EndPoint + "/solana/address/" + in.Address + "/assets"
	req, _ := http.NewRequest("GET", reqUrl, nil)
	query := req.URL.Query()
	query.Add("auth_key", s.env.Unmarshall.APIkey)
	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorf("Error:: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var unmarshalTokenBalRes []core.TokenBalance
	var assetResponse []*pb.TokenBalance
	err = json.Unmarshal(body, &unmarshalTokenBalRes)
	if err != nil {
		s.logger.Errorf("Error in Unmarshalling Json: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	for _, tokenBal := range unmarshalTokenBalRes {
		res := &pb.TokenBalance{
			ContractName:         tokenBal.ContractName,
			ContractTickerSymbol: tokenBal.ContractTickerSymbol,
			ContractDecimals:     tokenBal.ContractDecimals,
			ContractAddress:      tokenBal.ContractAddress,
			Coin:                 tokenBal.Coin,
			Balance:              tokenBal.Balance,
			Quote:                tokenBal.Quote,
			QuotePrice:           strconv.FormatFloat(tokenBal.Quote, 'f', -1, 64),
			QuoteRate:            tokenBal.QuoteRate,
			LogoUrl:              tokenBal.LogoUrl,
			QuoteRate_24H:        tokenBal.QuoteRate_24H,
			QuotePctChange_24H:   tokenBal.QuotePctChange_24H,
		}
		assetResponse = append(assetResponse, res)
	}
	return &pb.BalanceResponse{
		Token: assetResponse,
	}, nil
}
