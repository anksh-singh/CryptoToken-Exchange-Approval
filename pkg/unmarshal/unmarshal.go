package unmarshal

import (
	bridge_allowance "bridge-allowance"
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type IUnmarshall interface {
	GetAssets(address string, chainId string) ([]*UnmarshallAssetModel, error)
	ListTransaction(in *pb.ListTransactionRequest, chainId string) (*UnmarshallTransactionModel, error)
}

type UnmarshallService struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewUnMarshalService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *UnmarshallService {

	return &UnmarshallService{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
	}
}

// GetAssets Fetch wallet balances for an address
func (u *UnmarshallService) GetAssets(address string, chainName string) ([]*UnmarshallAssetModel, error) {
	u.logger.Infof("Initiating Unmarshal assets request for public address : %v", address)
	chainId, err := bridge_allowance.GetUnmarshalId(*u.env, chainName)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(u.env.Unmarshall.EndPoint+"/%s/address/"+
		"%s/assets?&auth_key=%s", chainId, address, u.env.Unmarshall.APIkey)
	body, err := u.httpRequest.GetRequest(url)
	u.logger.Infof("Url: %v", url)
	if err != nil {
		u.logger.Errorf(" Unmarshal request for Assets Logging Error  is : %v", err.Error())
		return nil, err
	}
	jsonResponseStruct, err := AssetsResponseTransformer(body)
	if err != nil {
		u.logger.Errorf(" Unmarshal request for Assets Logging Error  is : %v", err.Error())
		return nil, err
	}
	return jsonResponseStruct, err
}

// ListTransaction List of transactions for an address
func (u *UnmarshallService) ListTransaction(in *pb.ListTransactionRequest, chainName string) (*UnmarshallTransactionModel, error) {
	chainId, err := bridge_allowance.GetUnmarshalId(*u.env, chainName)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/address/%s/transactions?auth_key=%s&pageSize=%s&page=%s",
		u.env.Unmarshall.EndPointV2, chainId, in.Address, u.env.Unmarshall.APIkey, in.PageSize, in.Page)
	if len(in.TokenContractAddress) > 0 {
		url = fmt.Sprintf("%s&contract=%s", url, in.TokenContractAddress)
	}
	body, err := u.httpRequest.GetRequest(url)
	if err != nil {
		u.logger.Errorf("List Transaction  error : %v", err.Error())
		return nil, err
	}
	//var jsonResponseStruct UnmarshallTransactionModel
	//err = json.Unmarshal(body, &jsonResponseStruct)
	jsonResponseStruct, err := TransactionsResponseTransformer(body)
	if err != nil {
		u.logger.Errorf("List Transaction  error : %v", err.Error())
		return nil, err
	}
	return jsonResponseStruct, err
}

func (u *UnmarshallService) GetUserData(request *pb.UserDataRequest, chainName string) (*UserDataModel, error) {
	chainId, err := bridge_allowance.GetUnmarshalId(*u.env, chainName)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/address/%s/userData?contract=%s&auth_key=%s",
		u.env.Unmarshall.EndPointV2, chainId, request.Address, request.Contract, u.env.Unmarshall.APIkey)
	body, err := u.httpRequest.GetRequest(url)
	if err != nil {
		u.logger.Errorf("Error fetching url: %v  Err: %v", url, err.Error())
		return nil, err
	}
	response, err := TransformUserData(body)
	if err != nil {
		u.logger.Errorf("User Data transformation error : %v", err.Error())
	}
	return response, nil
}

func TransformUserData(data []byte) (*UserDataModel, error) {
	var responseObject UserDataModel
	err := json.Unmarshal(data, &responseObject)
	if err != nil {
		return nil, err
	}
	return &responseObject, nil
}

func TransactionsResponseTransformer(data []byte) (*UnmarshallTransactionModel, error) {
	var responseObject UnmarshallTransactionModel
	err := json.Unmarshal(data, &responseObject)
	if err != nil {
		return nil, err
	}
	return &responseObject, nil
}

func AssetsResponseTransformer(data []byte) ([]*UnmarshallAssetModel, error) {
	var responseObject []*UnmarshallAssetModel
	err := json.Unmarshal(data, &responseObject)
	if err != nil {
		return nil, err
	}
	return responseObject, nil

}
