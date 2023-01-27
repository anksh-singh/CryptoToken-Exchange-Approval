package rpc

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"context"
	"encoding/json"
)

type IPosition interface {
	GetPositions(ctx context.Context, in *pb.PositionChainData) *pb.GetPositionsResponse
}

func (s *SolanaManager) GetPositions(ctx context.Context, in *pb.PositionChainData) (*pb.GetPositionsResponse, error) {
	var response pb.GetPositionsResponse
	err := json.Unmarshal([]byte(SONAR_WATCH_DUMMY_RESPONSE), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
