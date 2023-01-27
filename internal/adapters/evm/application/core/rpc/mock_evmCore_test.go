package rpc

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"fmt"
	gomock "github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestGetTokenPrice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.TokenPriceRequest
		want  pb.TokenPriceResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.TokenPriceRequest{
				Chain:    "arbitrum",
				Currency: "usd",
			},
			want: pb.TokenPriceResponse{
				Price:          1234,
				CurrencyCode:   "USD",
				CurrencySymbol: "$",
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.TokenPriceRequest{
				Chain:    "arbitrum",
				Currency: "usd",
			},
			want: pb.TokenPriceResponse{
				Price:          12345,
				CurrencyCode:   "USD",
				CurrencySymbol: "$",
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.TokenPriceRequest
			mockCtrl.RecordCall(mockevm, "GetTokenPrice", &arg.input).DoAndReturn(
				func(arg *pb.TokenPriceRequest) *pb.TokenPriceResponse {
					doCalled = true
					argument = arg
					var r = pb.TokenPriceResponse{
						Price:          1234,
						CurrencyCode:   "USD",
						CurrencySymbol: "$",
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetTokenPrice", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.TokenPriceResponse)
			if arg.match {
				if ret.Price != arg.want.Price {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret.Price, arg.want.Price)
				}
			}
			if !arg.match {
				if ret.Price == arg.want.Price {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret.Price, arg.want.Price)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGasLimit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)
	type testdata struct {
		name  string
		input pb.GasLimitRequest
		want  pb.GasLimitResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.GasLimitRequest{
				Chain: "etherium",
				From:  "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				To:    "0x07a565b7ed7d7a678680a4c162885bedbb695fe0",
				Value: 1,
				Data:  "test data",
			},
			want: pb.GasLimitResponse{
				GasLimit:  75000,
				InputData: "test data",
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.GasLimitRequest{
				Chain: "etherium",
				From:  "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				To:    "0x07a565b7ed7d7a678680a4c162885bedbb695fe0",
				Value: 1,
				Data:  "test data",
			},
			want: pb.GasLimitResponse{
				GasLimit:  7500,
				InputData: "test data",
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.GasLimitRequest
			mockCtrl.RecordCall(mockevm, "GasLimit", &arg.input).DoAndReturn(
				func(arg *pb.GasLimitRequest) *pb.GasLimitResponse {
					doCalled = true
					argument = arg
					var r = pb.GasLimitResponse{
						GasLimit:  75000,
						InputData: "",
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GasLimit", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.GasLimitResponse)
			if arg.match {
				if ret.GasLimit != arg.want.GasLimit {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret.GasLimit, arg.want.GasLimit)
				}
			}
			if !arg.match {
				if ret.GasLimit == arg.want.GasLimit {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret.GasLimit, arg.want.GasLimit)
				}
			}
		})
	}
	mockCtrl.Finish()
}

func TestGetAssets(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)
	type testdata struct {
		name  string
		input pb.BalanceRequest
		want  pb.BalanceResponse
		match bool
	}

	asset := pb.TokenBalance{
		ContractName:         Clean("Ethereum").(string),
		ContractTickerSymbol: Clean("ETH").(string),
		ContractDecimals:     18,
		ContractAddress:      Clean("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee").(string),
		Coin:                 60,
		Type:                 Clean("ERC20").(string),
		Balance:              Clean("49976423895523587").(string),
		Quote:                Clean(77.8502745590129).(float64),
		QuoteRate:            Clean(1557.74).(float64),
		LogoUrl:              Clean("https://assets.unmarshal.io/tokens/ethereum_0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png").(string),
		QuoteRate_24H:        Clean(fmt.Sprintf("%v", "4.65")).(string),
		QuotePctChange_24H:   0.29962,
	}
	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.BalanceRequest{
				Chain:   "etherium",
				Address: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			want: pb.BalanceResponse{
				Address:       "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				QuoteCurrency: "USD",
				ChainId:       1,
				Token: []*pb.TokenBalance{
					&asset,
				},
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.BalanceRequest{
				Chain:   "etherium",
				Address: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			want: pb.BalanceResponse{
				Address:       "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				QuoteCurrency: "INR",
				ChainId:       12,
				Token: []*pb.TokenBalance{
					&asset,
				},
			},
			match: true,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.BalanceRequest
			mockCtrl.RecordCall(mockevm, "GetAssets", &arg.input).DoAndReturn(
				func(arg *pb.BalanceRequest) *pb.BalanceResponse {
					doCalled = true
					argument = arg
					var r = pb.BalanceResponse{
						Address:       "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
						QuoteCurrency: "USD",
						ChainId:       1,
						Token: []*pb.TokenBalance{
							&asset,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetAssets", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.BalanceResponse)
			if arg.match {
				if !reflect.DeepEqual(ret.Token, arg.want.Token) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, arg.want.Token)
				}
			}
			if !arg.match {
				//  if ret.GasLimit == arg.want.GasLimit
				if reflect.DeepEqual(ret.Token, arg.want.Token) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, arg.want.Token)
				}
			}
		})
	}
	mockCtrl.Finish()
}

func TestGetNonce(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.NonceRequest
		want  pb.NonceResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.NonceRequest{
				Chain:   "arbitrum",
				Address: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			want: pb.NonceResponse{
				Nonce:      0,
				QuoteValue: 1569.31,
				GasPrice: &pb.GasPriceInfo{
					Fast:        0.5,
					SafeLow:     0.30000001192092896,
					Fastest:     0.6000000238418579,
					Average:     0.4000000059604645,
					SafeLowWait: 10,
					AvgWait:     2,
					FastWait:    1,
					FastestWait: 0.5,
				},
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.NonceRequest{
				Chain:   "arbitrum",
				Address: "usd",
			},
			want: pb.NonceResponse{
				Nonce:      0,
				QuoteValue: 1569.31,
				GasPrice: &pb.GasPriceInfo{
					Fast:        0.51,
					SafeLow:     0.300000011920928961,
					Fastest:     0.60000002384185791,
					Average:     0.40000000596046451,
					SafeLowWait: 10,
					AvgWait:     2,
					FastWait:    1,
					FastestWait: 0.5,
				},
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.NonceRequest
			mockCtrl.RecordCall(mockevm, "GetNonce", &arg.input).DoAndReturn(
				func(arg *pb.NonceRequest) *pb.NonceResponse {
					doCalled = true
					argument = arg
					var r = pb.NonceResponse{
						Nonce:      0,
						QuoteValue: 1569.31,
						GasPrice: &pb.GasPriceInfo{
							Fast:        0.5,
							SafeLow:     0.30000001192092896,
							Fastest:     0.6000000238418579,
							Average:     0.4000000059604645,
							SafeLowWait: 10,
							AvgWait:     2,
							FastWait:    1,
							FastestWait: 0.5,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetNonce", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.NonceResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestSendTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.SendTransactionRequest
		want  pb.SendTransactionResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.SendTransactionRequest{
				Chain: "arbitrum",
				Msg:   "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			want: pb.SendTransactionResponse{
				TransactionId: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.SendTransactionRequest{
				Chain: "arbitrum",
				Msg:   "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
			},
			want: pb.SendTransactionResponse{
				TransactionId: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db1111",
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.SendTransactionRequest
			mockCtrl.RecordCall(mockevm, "SendTransaction", &arg.input).DoAndReturn(
				func(arg *pb.SendTransactionRequest) *pb.SendTransactionResponse {
					doCalled = true
					argument = arg
					var r = pb.SendTransactionResponse{
						TransactionId: "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "SendTransaction", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.SendTransactionResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestListTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)
	var tran = pb.TransactionData{
		Id:                  "0x753c53c04dda66ee62933db72074ab3ecbf1e671262fbc4cd419479ffd2d5af4",
		From:                "0x13756e7adf2ff5991593fab310f8906b768e5d89",
		To:                  "0x473037de59cf9484632f4a27b509cfe8d4a31404",
		Fee:                 "611025138041715",
		Date:                1658481272,
		Type:                "send",
		Block:               15191475,
		Value:               "0",
		Nonce:               27,
		NativeTokenDecimals: 18,
		Description:         "Sent 249.9194 GST",
		Sent:                nil,
		Received:            nil,
		Others:              nil,
	}

	type testdata struct {
		name  string
		input pb.ListTransactionRequest
		want  pb.ListTransactionResponse
		match bool
	}
	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.ListTransactionRequest{
				Chain:                "arbitrum",
				Address:              "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				Testnet:              true,
				Page:                 "1",
				PageSize:             "1",
				TokenContractAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			},
			want: pb.ListTransactionResponse{
				Page:        1,
				TotalPages:  1,
				ItemsOnPage: 1,
				TotalTxs:    1,
				Transactions: []*pb.TransactionData{
					&tran,
				},
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.ListTransactionRequest{
				Chain:                "arbitrum",
				Address:              "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				Testnet:              true,
				Page:                 "1",
				PageSize:             "1",
				TokenContractAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			},
			want: pb.ListTransactionResponse{
				Page:        1,
				TotalPages:  10,
				ItemsOnPage: 10,
				TotalTxs:    1,
				Transactions: []*pb.TransactionData{
					&tran,
				},
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ListTransactionRequest
			mockCtrl.RecordCall(mockevm, "ListTransaction", &arg.input).DoAndReturn(
				func(arg *pb.ListTransactionRequest) *pb.ListTransactionResponse {
					doCalled = true
					argument = arg
					var r = pb.ListTransactionResponse{
						Page:        1,
						TotalPages:  1,
						ItemsOnPage: 1,
						TotalTxs:    1,
						Transactions: []*pb.TransactionData{
							&tran,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "ListTransaction", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.ListTransactionResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGetUserData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.UserDataRequest
		want  pb.UserDataResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.UserDataRequest{
				Chain:    "arbitrum",
				Address:  "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				Contract: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			},
			want: pb.UserDataResponse{
				QuoteRate:              1591.28,
				TotalFeesPaid:          0,
				TotalFeesPaidUsd:       0,
				AverageTokenPrice:      0,
				OverallProfitLoss:      0,
				CurrentHoldingQuantity: 0,
				PercentageChange_24H:   6.50931,
				PriceChange_24H:        0,
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.UserDataRequest{
				Chain:    "arbitrum",
				Address:  "0x1923f626bb8dc025849e00f99c25fe2b2f7fb0db",
				Contract: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			},
			want: pb.UserDataResponse{
				QuoteRate:              1591.28,
				TotalFeesPaid:          10,
				TotalFeesPaidUsd:       10,
				AverageTokenPrice:      10,
				OverallProfitLoss:      10,
				CurrentHoldingQuantity: 0,
				PercentageChange_24H:   6.50931,
				PriceChange_24H:        0,
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.UserDataRequest
			mockCtrl.RecordCall(mockevm, "GetUserData", &arg.input).DoAndReturn(
				func(arg *pb.UserDataRequest) *pb.UserDataResponse {
					doCalled = true
					argument = arg
					var r = pb.UserDataResponse{
						QuoteRate:              1591.28,
						TotalFeesPaid:          0,
						TotalFeesPaidUsd:       0,
						AverageTokenPrice:      0,
						OverallProfitLoss:      0,
						CurrentHoldingQuantity: 0,
						PercentageChange_24H:   6.50931,
						PriceChange_24H:        0,
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetUserData", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.UserDataResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGetProcessingFee(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.ProcessingFeeRequest
		want  pb.ProcessingFeeResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.ProcessingFeeRequest{
				Chain:          "arbitrum",
				GasPrice:       true,
				TransactionFee: true,
			},
			want: pb.ProcessingFeeResponse{
				Value: 33.5,
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.ProcessingFeeRequest{
				Chain:          "arbitrum",
				GasPrice:       true,
				TransactionFee: true,
			},
			want: pb.ProcessingFeeResponse{
				Value: 30,
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ProcessingFeeRequest
			mockCtrl.RecordCall(mockevm, "GetProcessingFee", &arg.input).DoAndReturn(
				func(arg *pb.ProcessingFeeRequest) *pb.ProcessingFeeResponse {
					doCalled = true
					argument = arg
					var r = pb.ProcessingFeeResponse{
						Value: 33.5,
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetProcessingFee", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.ProcessingFeeResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGetTxStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.TxStatusRequest
		want  pb.TxStatusResponse
		match bool
	}

	var logs = pb.Log{
		Removed:          false,
		LogIndex:         4,
		TransactionIndex: 0,
		TransactionHash:  "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
		BlockNumber:      15186345,
		BlockHash:        "0x5e6010b29c1425a6d48c597bf0464b9003cbb7b69a7dc3f24e3a5135e11bd19c",
		Address:          "0x9ea468e0ffc6a6dbfc6c81e81efe15719fd88b67",
		Data:             "0x000000000000000000000000000000000000000000090ea91a1c4a1bd77d52ba00000000000000000000000000000000000000000000000c8a68adae4fa598bf",
		Topics:           []string{"0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"},
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.TxStatusRequest{
				Chain:  "arbitrum",
				TxHash: "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
			},
			want: pb.TxStatusResponse{
				TransactionHash:   "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
				TransactionIndex:  0,
				BlockHash:         "0x5e6010b29c1425a6d48c597bf0464b9003cbb7b69a7dc3f24e3a5135e11bd19c",
				BlockNumber:       15186345,
				CumulativeGasUsed: 104709,
				GasUsed:           104709,
				ContractAddress:   "",
				Logs:              []*pb.Log{&logs},
				LogsBloom:         "0x00200000000000000000000080000000002000000000000000000000020000000000000000000004000000000000020002000000080001000000000000200000000000000000000000000008000000200000000000000000000000008000200000000000000000000000000000000000000200000000000000000010000000000000000000000000100000000000000000000001000000080000004000000000004800000000000000000000000000000000000000000000040000000000000000000002000000000000000000000000000000000000001000000000000000000000200000000000000005000000000000000000000000400000000000000000",
				Root:              "",
				Status:            "SUCCESS",
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.TxStatusRequest{
				Chain:  "arbitrum",
				TxHash: "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
			},
			want: pb.TxStatusResponse{
				TransactionHash:   "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
				TransactionIndex:  0,
				BlockHash:         "0x5e6010b29c1425a6d48c597bf0464b9003cbb7b69a7dc3f24e3a5135e11bd19c",
				BlockNumber:       151863451,
				CumulativeGasUsed: 1047091,
				GasUsed:           1047091,
				ContractAddress:   "",
				Logs:              []*pb.Log{&logs},
				LogsBloom:         "0x00200000000000000000000080000000002000000000000000000000020000000000000000000004000000000000020002000000080001000000000000200000000000000000000000000008000000200000000000000000000000008000200000000000000000000000000000000000000200000000000000000010000000000000000000000000100000000000000000000001000000080000004000000000004800000000000000000000000000000000000000000000040000000000000000000002000000000000000000000000000000000000001000000000000000000000200000000000000005000000000000000000000000400000000000000000",
				Root:              "",
				Status:            "SUCCESS",
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.TxStatusRequest
			mockCtrl.RecordCall(mockevm, "GetTxStatus", &arg.input).DoAndReturn(
				func(arg *pb.TxStatusRequest) *pb.TxStatusResponse {
					doCalled = true
					argument = arg
					var r = pb.TxStatusResponse{
						TransactionHash:   "0x68231b646b4d9635dba01f50143cd34705042f167d603ae34726ad1c8aae0c6f",
						TransactionIndex:  0,
						BlockHash:         "0x5e6010b29c1425a6d48c597bf0464b9003cbb7b69a7dc3f24e3a5135e11bd19c",
						BlockNumber:       15186345,
						CumulativeGasUsed: 104709,
						GasUsed:           104709,
						ContractAddress:   "",
						Logs:              []*pb.Log{&logs},
						LogsBloom:         "0x00200000000000000000000080000000002000000000000000000000020000000000000000000004000000000000020002000000080001000000000000200000000000000000000000000008000000200000000000000000000000008000200000000000000000000000000000000000000200000000000000000010000000000000000000000000100000000000000000000001000000080000004000000000004800000000000000000000000000000000000000000000040000000000000000000002000000000000000000000000000000000000001000000000000000000000200000000000000005000000000000000000000000400000000000000000",
						Root:              "",
						Status:            "SUCCESS",
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetTxStatus", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.TxStatusResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGetTokenAllowance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.AllowanceRequest
		want  pb.AllowanceResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.AllowanceRequest{
				Chain:    "etherium",
				Contract: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
				Owner:    "shi",
				Spender:  "24",
			},
			want: pb.AllowanceResponse{
				Allowance: "0.000001",
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.AllowanceRequest{
				Chain:    "etherium",
				Contract: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
				Owner:    "shi",
				Spender:  "24",
			},
			want: pb.AllowanceResponse{
				Allowance: "0.000002",
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.AllowanceRequest
			mockCtrl.RecordCall(mockevm, "GetTokenAllowance", &arg.input).DoAndReturn(
				func(arg *pb.AllowanceRequest) *pb.AllowanceResponse {
					doCalled = true
					argument = arg
					var r = pb.AllowanceResponse{
						Allowance: "0.000001",
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetTokenAllowance", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.AllowanceResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestTokenApprove(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.ApprovalRequest
		want  pb.ApprovalResponse
		match bool
	}

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.ApprovalRequest{
				Chain:  "etherium",
				Target: "addr1",
				Token:  "dum token",
			},
			want: pb.ApprovalResponse{
				To:        "dum token",
				Value:     "0",
				Data:      "0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff",
				GasLimit:  "72000",
				TxLink:    "https://txlink.io/tx?to=wddmwldm&value=0&data=0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff&gaslimit=72000&gasPrice=8",
				Allowance: nil,
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.ApprovalRequest{
				Chain:  "etherium",
				Target: "addr1",
				Token:  "dum token",
			},
			want: pb.ApprovalResponse{
				To:        "dum token",
				Value:     "10",
				Data:      "0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff",
				GasLimit:  "2000",
				TxLink:    "https://txlink.io/tx?to=wddmwldm&value=0&data=0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff&gaslimit=72000&gasPrice=8",
				Allowance: nil,
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ApprovalRequest
			mockCtrl.RecordCall(mockevm, "TokenApprove", &arg.input).DoAndReturn(
				func(arg *pb.ApprovalRequest) *pb.ApprovalResponse {
					doCalled = true
					argument = arg
					var r = pb.ApprovalResponse{
						To:        "dum token",
						Value:     "0",
						Data:      "0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff",
						GasLimit:  "72000",
						TxLink:    "https://txlink.io/tx?to=wddmwldm&value=0&data=0x095ea7b300000000000000000000000000000000000000000000000000000000000aaaaa0000000000000000000000000000000000000000ffffffffffffffffffffffff&gaslimit=72000&gasPrice=8",
						Allowance: nil,
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "TokenApprove", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.ApprovalResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

func TestGetNftCollections(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testdata struct {
		name  string
		input pb.NftCollectionRequest
		want  pb.ListNftCollectionResponse
		match bool
	}

	var resp = pb.NftCollectionResponse{
		BannerImageUrl:          "https://openseauserdata.com/files/935f943ff19a4857de0f2fc0157051b0.png",
		ChatUrl:                 "",
		CreatedDate:             "2022-07-15T23:43:29.587300",
		DefaultToFiat:           false,
		Description:             "🇸🇦 5,555 MoonShabaabs have taken over the moon 🇸🇦 We got the oil, now we buy the moon ☽\n2 Free per wallet, rest 0.0069✨",
		DevBuyerFeeBasisPoints:  "0",
		DevSellerFeeBasisPoints: "1000",
		DiscordUrl:              "",
		DisplayData: &pb.NFTDisplayData{
			CardDisplayStyle: "contain",
		},
		ExternalUrl:                 "https://etherscan.io/address/0xd4c7d17565de4108dd1d01877d9424f899e058d8",
		Featured:                    false,
		FeaturedImageUrl:            "https://openseauserdata.com/files/5d6eda55a1c3a433948e1f4f830553d6.png",
		Hidden:                      false,
		SafelistRequestStatus:       "not_requested",
		ImageUrl:                    "https://openseauserdata.com/files/25cdc59132225af8ce5b4e0525f3b97e.jpg",
		IsSubjectToWhitelist:        false,
		LargeImageUrl:               "https://openseauserdata.com/files/5d6eda55a1c3a433948e1f4f830553d6.png",
		MediumUsername:              "",
		Name:                        "MoonShabaabs Official",
		OnlyProxiedTransfers:        false,
		OpenseaBuyerFeeBasisPoints:  "0",
		OpenseaSellerFeeBasisPoints: "250",
		PayoutAddress:               "0x3d42aaad435ce4a226714f6b882bcdff1264d2c4",
		RequireEmail:                false,
		ShortDescription:            "",
		Slug:                        "moonshabaabs",
		TelegramUrl:                 "",
		TwitterUsername:             "MoonShabaabs",
		InstagramUsername:           "",
		WikiUrl:                     "",
		IsNsfw:                      false,
		NftData: []*pb.NftData{
			&pb.NftData{
				Id:                   542154241,
				NumSales:             0,
				BackgroundColor:      "",
				ImageUrl:             "https://lh3.googleusercontent.com/IUM2BEwJHRMrHzgiqTwZ7m0c6TWICCz9C1dTTYljRVKOKTs9oqnW4X-CRjBlXfXCzuklJgFN1MJclFmQOQ00LKQqvgGXBf0OkqY1Sw",
				ImagePreviewUrl:      "https://lh3.googleusercontent.com/IUM2BEwJHRMrHzgiqTwZ7m0c6TWICCz9C1dTTYljRVKOKTs9oqnW4X-CRjBlXfXCzuklJgFN1MJclFmQOQ00LKQqvgGXBf0OkqY1Sw=s250",
				ImageThumbnailUrl:    "https://lh3.googleusercontent.com/IUM2BEwJHRMrHzgiqTwZ7m0c6TWICCz9C1dTTYljRVKOKTs9oqnW4X-CRjBlXfXCzuklJgFN1MJclFmQOQ00LKQqvgGXBf0OkqY1Sw=s128",
				ImageOriginalUrl:     "https://opensea.mypinata.cloud/ipfs/QmQrZXqMhreCRVrtE4JfnvshJxnDgntyvbX3HhK6F18tQA/4128.png",
				AnimationUrl:         "",
				AnimationOriginalUrl: "",
				Name:                 "",
				Description:          "We got the oil, now we buy the Moon!",
				ExternalLink:         "",
				AssetContract: &pb.NftDataAssetContract{

					Address: "",
				},
				Permalink:         "https://opensea.io/assets/ethereum/0xd4c7d17565de4108dd1d01877d9424f899e058d8/4128",
				Decimals:          "",
				TokenMetadata:     "https://opensea.mypinata.cloud/ipfs/QmXEcsvvF473NokwnWZDTq8Yi2SARv8vKwoo6sZbh6UBXk/4128.json",
				IsNsfw:            false,
				Owner:             &pb.NftDataOwner{},
				SellOrders:        nil,
				SeaportSellOrders: "",
				Creator: &pb.NftDataCreator{
					User: &pb.NftDataUser{
						Username: "",
					},
					ProfileImgUrl: "https://storage.googleapis.com/opensea-static/opensea-profile/2.png",
					Address:       "0xe1c4be771735f270cc678246e238c57949b46ddb",
					Config:        "",
				},
				Traits: []*pb.NftDataTraits{
					&pb.NftDataTraits{
						TraitType:   "Head",
						Value:       "Dark",
						DisplayType: "",
						MaxValue:    "",
						TraitCount:  0,
						Order:       "",
					},
				},
				LastSale:                nil,
				TopBid:                  "",
				ListingDate:             "",
				IsPresale:               false,
				TransferFeePaymentToken: "",
				TransferFee:             "",
				TokenId:                 "4128",
				CollectionName:          "MoonShabaabs Official",
				ContractAddress:         "0xd4c7d17565de4108dd1d01877d9424f899e058d8",
			},
		},
	}
	var invalidresp pb.NftCollectionResponse = resp
	invalidresp.BannerImageUrl = "dummy"
	invalidresp.Name = "invalid name"

	var testdatainput = []testdata{
		testdata{
			name: "validreturn",
			input: pb.NftCollectionRequest{
				Address:  "0xbaf28d992ad354e56a4d6da997900ab3738d11ab",
				Page:     "1",
				PageSize: "10",
				Chain:    "etherium",
			},
			want: pb.ListNftCollectionResponse{
				Nft: []*pb.NftCollectionResponse{
					&resp,
				},
			},
			match: true,
		},
		testdata{
			name: "invalidreturn",
			input: pb.NftCollectionRequest{
				Address:  "0xbaf28d992ad354e56a4d6da997900ab3738d11ab",
				Page:     "1",
				PageSize: "10",
				Chain:    "etherium",
			},
			want: pb.ListNftCollectionResponse{
				Nft: []*pb.NftCollectionResponse{
					&invalidresp,
				},
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.NftCollectionRequest
			mockCtrl.RecordCall(mockevm, "GetNftCollections", &arg.input).DoAndReturn(
				func(arg *pb.NftCollectionRequest) *pb.ListNftCollectionResponse {
					doCalled = true
					argument = arg
					var r = pb.ListNftCollectionResponse{
						Nft: []*pb.NftCollectionResponse{
							&resp,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			rets := mockCtrl.Call(mockevm, "GetNftCollections", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(rets) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(rets))
			}
			ret := rets[0].(*pb.ListNftCollectionResponse)
			if arg.match {
				if !reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(ret, &arg.want) {
					t.Errorf("DoAndReturn return value: got %v, want %v", ret, &arg.want)
				}
			}
		})
	}

	mockCtrl.Finish()
}

// TestCases for GetProtocols and GetStakeProtocols
func TestGetProtocols(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testinputs struct {
		name  string
		input pb.ProtocolRequest
		want  pb.ProtocolResponse
		match bool
	}

	// inuput request test data
	inputRequest := pb.ProtocolRequest{
		Chain: "bsc",
	}

	successResponse := pb.ProtocolInfo{
		ProtocolName: "pancakeswap",
		Logo:         "https://icons.llama.fi/pancakeswap.jpg",
		Description:  "The #1 AMM and yield farm on Binance Smart Chain",
		EarningModes: []string{"staking", "farming"},
		Tvl:          "https://storage.googleapis.com/opensea-static/opensea-tvl/bsc.png",
	}

	failResponse := pb.ProtocolInfo{
		ProtocolName: "madhuswap",
		Logo:         "https://icons.llama.fi/pancakeswap.jpg",
		Description:  "Bsc is a decentralized protocol for the blockchain",
		EarningModes: []string{"staking", "farming"},
		Tvl:          "3067468159.2999897",
	}

	var testdatainput = []testinputs{
		{
			name:  "success",
			input: inputRequest,
			want: pb.ProtocolResponse{
				Protocols: []*pb.ProtocolInfo{
					&successResponse,
				},
			},
			match: true,
		},
		{
			name:  "failure",
			input: inputRequest,
			want: pb.ProtocolResponse{
				Protocols: []*pb.ProtocolInfo{
					&failResponse,
				},
			},
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ProtocolRequest
			mockCtrl.RecordCall(mockevm, "GetProtocols", &arg.input).DoAndReturn(
				func(arg *pb.ProtocolRequest) *pb.ProtocolResponse {
					doCalled = true
					argument = arg
					var r = pb.ProtocolResponse{
						Protocols: []*pb.ProtocolInfo{
							&successResponse,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			response := mockCtrl.Call(mockevm, "GetProtocols", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}

			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(response) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(response))
			}
			results := response[0].(*pb.ProtocolResponse)
			if arg.match {
				if !reflect.DeepEqual(results.Protocols, arg.want.Protocols) {
					t.Errorf("DoAndReturn return value: got %v, want %v", results.Protocols, arg.want.Protocols)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(results.Protocols, arg.want.Protocols) {
					t.Errorf("DoAndReturn return value: got %v, want %v", results.Protocols, arg.want.Protocols)
				}
			}
		})
	}
	mockCtrl.Finish()
}

func TestGetStakeProtocols(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockevm := NewMockEvmCore(mockCtrl)

	type testinputs struct {
		name  string
		input pb.ProtocolRequest
		want  pb.StakeProtocolResponse
		match bool
	}

	// inuput request test data
	inputRequest := pb.ProtocolRequest{
		Chain: "bsc",
	}

	successResponse := pb.StakeProtocolResponse{
		StakeProtocols: []string{"pancakeswap", "apeswap"},
	}

	failResponse := pb.StakeProtocolResponse{
		StakeProtocols: []string{"rangeswap", "apeswap"},
	}

	var testdatainput = []testinputs{
		{
			name:  "Success",
			input: inputRequest,
			want:  successResponse,
			match: true,
		},
		{
			name:  "Failure",
			input: inputRequest,
			want:  failResponse,
			match: false,
		},
	}

	for _, arg := range testdatainput {
		t.Run(arg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ProtocolRequest
			mockCtrl.RecordCall(mockevm, "GetStakeProtocols", &arg.input).DoAndReturn(
				func(arg *pb.ProtocolRequest) *pb.StakeProtocolResponse {
					doCalled = true
					argument = arg
					var r = successResponse
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			response := mockCtrl.Call(mockevm, "GetStakeProtocols", &arg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}

			if &arg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(response) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(response))
			}
			results := response[0].(*pb.StakeProtocolResponse)
			if arg.match {
				if !reflect.DeepEqual(results.StakeProtocols, arg.want.StakeProtocols) {
					t.Errorf("DoAndReturn return value: got %v, want %v", results.StakeProtocols, arg.want.StakeProtocols)
				}
			}
			if !arg.match {
				if reflect.DeepEqual(results.StakeProtocols, arg.want.StakeProtocols) {
					t.Errorf("DoAndReturn return value: got %v, want %v", results.StakeProtocols, arg.want.StakeProtocols)
				}
			}
		})
	}
	mockCtrl.Finish()
}
