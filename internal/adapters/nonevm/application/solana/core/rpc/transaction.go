package rpc

import (
	"bridge-allowance/internal/adapters/nonevm/application/solana/core"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
)

func (s *SolanaManager) SendTransaction(in *pb.SendTransactionRequest) (*pb.SendTransactionResponse, error) {
	var sendTxRes core.SendTransactionResponse
	sendTxMsg := core.SendTransaction{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "sendTransaction",
		Params:  []string{in.Msg},
	}
	sendTxRawMsg, err := json.Marshal(sendTxMsg)
	if err != nil {
		s.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Marshaling the SendTx Request")
	}
	res, err := s.httpRequest.PostRequest(rpc.MainNetBeta_RPC, bytes.NewBuffer(sendTxRawMsg)) //Mainnet only Enabled
	if err != nil {
		s.logger.Errorf(" Error in Sending Transaction: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Sending Transaction")
	}
	err = json.Unmarshal(res, &sendTxRes)
	if err != nil {
		s.logger.Errorf("Error in Unmarshaling SendTx Response: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Unmarshaling SendTx Response")
	}
	if sendTxRes.Error.Code < 0 {
		s.logger.Errorf("Error in  SendTx Response: %v", string(res))
		return nil, status.Errorf(codes.FailedPrecondition, sendTxRes.Error.Message, "Error in  SendTx Response")
	}
	return &pb.SendTransactionResponse{
		TransactionId: sendTxRes.Result,
	}, nil
}

func (s *SolanaManager) ListTransaction(in *pb.ListTransactionRequest) (*pb.ListTransactionResponse, error) {
	reqUrl := s.env.Unmarshall.EndPoint + "/solana/address/" + in.Address + "/transactions"
	req, _ := http.NewRequest("GET", reqUrl, nil)
	query := req.URL.Query()
	query.Add("page", in.Page)
	query.Add("pageSize", in.PageSize)
	query.Add("contract", " ")
	query.Add("auth_key", s.env.Unmarshall.APIkey)
	req.URL.RawQuery = query.Encode()
	res, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorf("List Transaction  error : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Fetching Solana List Transaction")
	}
	var transResponse *pb.ListTransactionResponse
	err = json.Unmarshal(body, &transResponse)
	if err != nil {
		s.logger.Errorf("Error in Unmarshalling Json Response: %v", err.Error())
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("Unmarshal Response %s causes json unmarshal error", string(body)),
			"Error in Fetching Solana List Transaction")
	}
	return transResponse, nil
}

func (s *SolanaManager) TransactionStatus(in *pb.TxStatusRequest) (*pb.TxStatusResponse, error) {
	solanaSig, err := solana.SignatureFromBase58(in.TxHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Solana Signature")
	}
	res, err := s.client.GetSignatureStatuses(
		context.TODO(),
		true,
		solanaSig,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Solana Get Signature Status API")
	}
	txStatusInfo := res.Value[0]
	if txStatusInfo == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Empty Response Value in Signature Status API", "Invalid TxHash")
	}
	if txStatusInfo.Err == nil {
		return &pb.TxStatusResponse{
			TransactionHash: in.TxHash,
			Status:          fmt.Sprintf("%v", txStatusInfo.ConfirmationStatus),
		}, nil
	} else {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("%v", txStatusInfo.Err), "Failed in Fetching Transaction Status")
	}
}
func (s *SolanaManager) GetRecentBlockHash() (*rpc.GetRecentBlockhashResult, error) {
	res, err := s.client.GetRecentBlockhash(
		context.TODO(), rpc.CommitmentConfirmed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Failed in Fetching Recent Blockhash")
	}
	return res, nil
}
