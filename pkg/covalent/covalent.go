package covalent

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type CovalentService struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewCovalentService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *CovalentService {
	return &CovalentService{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

func (c *CovalentService) GetAssets(address string, chainid string) (*AssetsResponseForCovalent, error) {
	c.logger.Infof("initiating covalent request for Assets for public address : %v", address)
	url := fmt.Sprintf(c.env.Covalent.EndPoint+"/%s/address/"+
		"%s/balances_v2/?&key=%s", chainid, address, c.env.Covalent.APIKey)
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Errorf(" covalent request for Assets Logging Error  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var jsonResponseStruct *AssetsResponseForCovalent
	err = json.Unmarshal(body, &jsonResponseStruct)
	if err != nil {
		c.logger.Errorf(" covalent request for Assets Logging Error  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshal error")
	}
	return jsonResponseStruct, err
}

func (c *CovalentService) ListTransaction(in *pb.ListTransactionRequest, chainid string) (*ListTransactionListForCovalent, error) {
	c.logger.Infof("initiating covalent List transactions request for public address : %v", in.Address)
	pageInt, _ := strconv.Atoi(in.Page)
	pageCorrection := pageInt - 1 //Covalent pagination starts with 0(need to handle it by reducing 1)
	url := fmt.Sprintf(c.env.Covalent.EndPoint+"/%s/address/"+
		"%s/transactions_v2/?&page-number=%s&page-size=%s&key=%s", chainid, in.Address, string(pageCorrection), in.PageSize, c.env.Covalent.APIKey)
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Errorf(" List transactions request Logging Error is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var jsonResponseStruct *ListTransactionListForCovalent
	err = json.Unmarshal(body, &jsonResponseStruct)
	if err != nil {
		c.logger.Errorf(" List transactions request Logging Error is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshal error")
	}
	return jsonResponseStruct, err
}
