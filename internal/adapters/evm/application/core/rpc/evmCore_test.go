package rpc

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/onrik/ethrpc"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clean(tt.args.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEVMCore(t *testing.T) {
	type args struct {
		config   config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	tests := []struct {
		name string
		args args
		want *evmCore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEVMCore(&tt.args.config, tt.args.logger, tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEVMCore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuotepctchange24h(t *testing.T) {
	type args struct {
		quoteRate    float64
		quoteRate24H float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quotepctchange24h(tt.args.quoteRate, tt.args.quoteRate24H); got != tt.want {
				t.Errorf("Quotepctchange24h() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_GasLimit(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.GasLimitRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GasLimitResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.GasLimit(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GasLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GasLimit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_GetAssets(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.BalanceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.BalanceResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.GetAssets(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAssets() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_GetNonce(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.NonceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.NonceResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.GetNonce(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNonce() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_GetProcessingFee(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.ProcessingFeeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ProcessingFeeResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.GetProcessingFee(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProcessingFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProcessingFee() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_GetTokenPrice(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.TokenPriceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TokenPriceResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.GetTokenPrice(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokenPrice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_ListTransaction(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.ListTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListTransactionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.ListTransaction(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evmCore_SendTransaction(t *testing.T) {
	type fields struct {
		rpc      map[string]*ethrpc.EthRPC
		env      config.Config
		logger   *zap.SugaredLogger
		services common.Services
	}
	type args struct {
		request *pb.SendTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.SendTransactionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evm := &evmCore{
				rpc:      tt.fields.rpc,
				env:      &tt.fields.env,
				logger:   tt.fields.logger,
				services: tt.fields.services,
			}
			got, err := evm.SendTransaction(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
