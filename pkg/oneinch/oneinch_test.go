package oneinch

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

// Test Cases for GetExchangeTokens
func TestGetExchangeTokens(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockExchange := NewMockIOneInch(mockCtrl)

	type testdata struct {
		name  string
		input InputRequest
		want  pb.ExchangeTokenResponse
		match bool
	}

	// input request test data
	ExchangeTokenReqInfo := InputRequest{
		Chain:        "bsc",
		ExchangeType: "1inch",
	}

	validResponse := pb.ExchangeTokenInfo{
		TokenAddress:  "0xb59490ab09a0f526cc7305822ac65f2ab12f9723",
		TokenDecimals: "18",
		TokenSymbol:   "LIT",
		TokenName:     "Litentry",
		TokenLogoUrl:  "https://tokens.1inch.io/0xb59490ab09a0f526cc7305822ac65f2ab12f9723.png",
		LogoUrl:       "https://assets.coingecko.com/coins/images/13469/large/1inch-token.png?1608803028",
	}

	invalidResponse := pb.ExchangeTokenInfo{
		TokenAddress:  "18",
		TokenDecimals: "0xb59490ab09a0f526cc7305822ac65f2ab12f9723",
		TokenSymbol:   "Litentry",
		TokenName:     "LIT",
		TokenLogoUrl:  "https://tokens.1inch.io/0xb59490ab09a0f526cc7305822ac65f2ab12f9723.png",
		LogoUrl:       "https://assets.coingecko.com/coins/images/13469/large/1inch-token.png?1608803028",
	}

	var testDataInput = []testdata{
		{
			name:  "Success Response",
			input: ExchangeTokenReqInfo,
			want: pb.ExchangeTokenResponse{
				ExchangeTokens: []*pb.ExchangeTokenInfo{
					&validResponse,
				},
			},
			match: true,
		},
		{
			name:  "Failure Response",
			input: ExchangeTokenReqInfo,
			want: pb.ExchangeTokenResponse{
				ExchangeTokens: []*pb.ExchangeTokenInfo{
					&invalidResponse,
				},
			},
			match: false,
		},
	}
	for _, inputArg := range testDataInput {
		t.Run(inputArg.name, func(t *testing.T) {
			doCalled := false
			var inputReq InputRequest
			mockCtrl.RecordCall(mockExchange, "GetExchangeTokens", inputArg.input).DoAndReturn(
				func(arg InputRequest) *pb.ExchangeTokenResponse {
					doCalled = true
					inputReq = arg
					var r = pb.ExchangeTokenResponse{
						ExchangeTokens: []*pb.ExchangeTokenInfo{
							&validResponse,
						},
					}
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			responseResult := mockCtrl.Call(mockExchange, "GetExchangeTokens", inputArg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if inputArg.input != inputReq {
				t.Error("Do callback received wrong argument.")
			}
			if len(responseResult) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(responseResult))
			}
			response := responseResult[0].(*pb.ExchangeTokenResponse)
			if inputArg.match {
				if !reflect.DeepEqual(response.ExchangeTokens, inputArg.want.ExchangeTokens) {
					t.Errorf("Return values from Call: got %v, want %v", response.ExchangeTokens, inputArg.want.ExchangeTokens)
				}
			}
			if !inputArg.match {
				if reflect.DeepEqual(response.ExchangeTokens, inputArg.want.ExchangeTokens) {
					t.Errorf("Return values from Call: got %v, want %v", response.ExchangeTokens, inputArg.want.ExchangeTokens)
				}
			}
		})
	}

	mockCtrl.Finish()
}

// Test Cases for GetExchangeQuote
func TestGetExchangeQuote(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockQuote := NewMockIOneInch(mockCtrl)

	type testdata struct {
		name  string
		input pb.ExchangeQuoteRequest
		want  pb.ExchangeQuoteResponse
		match bool
	}

	// input request test data
	requestInfo := pb.ExchangeQuoteRequest{
		Chain:        "bsc",
		ExchangeType: "1inch",
		TakerAddress: "0x3aF210Eec25E95f470fDA218C9DF9291adb7dcFD",
		SellToken:    "0x603c7f932ED1fc6575303D8Fb018fDCBb0f39a95",
		BuyToken:     "0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82",
		SellAmount:   "0.01",
		Slippage:     "0.1",
	}

	validResponse := pb.ExchangeQuoteResponse{
		ResAmount:            "0.000248159061439969",
		PriceImpact:          "0",
		ResPricePerFromToken: "0.0248159061439969",
		ResPricePerToToken:   "40.29673525509788",
		FromTokenPrice:       "0.00103584",
		ToTokenPrice:         "0.0010422680580478697",
		MinimumReceived:      "0.000248",
	}

	invalidResponse := pb.ExchangeQuoteResponse{
		ResAmount:            "0.000248159061439969",
		PriceImpact:          "0",
		ResPricePerFromToken: "0.0248159061439969",
		ResPricePerToToken:   "40.29673525509788",
		FromTokenPrice:       "0.00103584",
		ToTokenPrice:         "0.0010422680580478697",
		MinimumReceived:      "1",
	}

	var testDataInput = []testdata{
		{
			name:  "Success Response",
			input: requestInfo,
			want:  validResponse,
			match: true,
		},
		{
			name:  "Failure Response",
			input: requestInfo,
			want:  invalidResponse,
			match: false,
		},
	}

	for _, inputArg := range testDataInput {
		t.Run(inputArg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ExchangeQuoteRequest
			mockCtrl.RecordCall(mockQuote, "GetExchangeQuote", &inputArg.input).DoAndReturn(
				func(arg *pb.ExchangeQuoteRequest) *pb.ExchangeQuoteResponse {
					doCalled = true
					argument = arg
					var r = validResponse
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			responses := mockCtrl.Call(mockQuote, "GetExchangeQuote", &inputArg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &inputArg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(responses) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(responses))
			}
			resp := responses[0].(*pb.ExchangeQuoteResponse)
			if inputArg.match {
				if resp.MinimumReceived != inputArg.want.MinimumReceived {
					t.Errorf("DoAndReturn return value: got %v, want %v", resp.MinimumReceived, inputArg.want.MinimumReceived)
				}
			}
			if !inputArg.match {
				if resp.MinimumReceived == inputArg.want.MinimumReceived {
					t.Errorf("DoAndReturn return value: got %v, want %v", resp.MinimumReceived, inputArg.want.MinimumReceived)
				}
			}
		})
	}

	mockCtrl.Finish()
}

// Test Cases for GetExchangeSwap
func TestGetExchangeSwap(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockQuote := NewMockIOneInch(mockCtrl)

	type testdata struct {
		name  string
		input pb.ExchangeSwapRequest
		want  pb.ExchangeSwapResponse
		match bool
	}

	// input request test data
	requestInfo := pb.ExchangeSwapRequest{
		Chain:        "bsc",
		ExchangeType: "1inch",
		TakerAddress: "0x3aF210Eec25E95f470fDA218C9DF9291adb7dcFD",
		SellToken:    "0x603c7f932ED1fc6575303D8Fb018fDCBb0f39a95",
		BuyToken:     "0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82",
		SellAmount:   "0.01",
		Slippage:     "0.1",
	}

	validResponse := pb.ExchangeSwapResponse{
		To:       "0x1111111254fb6c44bac0bed2854e76f90643097d",
		Data:     "0x7c02520000000000000000000000000005ad60d9a2f1aa30ba0cdbaf1e0a0a145fbea16f00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000180000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e8930000000000000000000000003af210eec25e95f470fda218c9df9291adb7dcfd000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000e1b274ca648a00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a4b757fed60000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e893000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000001e84801111111254fb6c44bac0bed2854e76f90643097d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000cfee7c08",
		Value:    "0",
		GasLimit: "196150",
		Gas:      "15000000000",
		TxLink:   "https://txlink.io/tx?to=0x1111111254fb6c44bac0bed2854e76f90643097d&value=0&data=0x7c02520000000000000000000000000005ad60d9a2f1aa30ba0cdbaf1e0a0a145fbea16f00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000180000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e8930000000000000000000000003af210eec25e95f470fda218c9df9291adb7dcfd000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000e1b274ca648a00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a4b757fed60000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e893000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000001e84801111111254fb6c44bac0bed2854e76f90643097d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000cfee7c08&gaslimit=196150",
	}

	invalidResponse := pb.ExchangeSwapResponse{
		To:       "",
		Data:     "0x7c02520000000000000000000000000005ad60d9a2f1aa30ba0cdbaf1e0a0a145fbea16f00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000180000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e8930000000000000000000000003af210eec25e95f470fda218c9df9291adb7dcfd000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000e1b274ca648a00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a4b757fed60000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e893000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000001e84801111111254fb6c44bac0bed2854e76f90643097d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000cfee7c08",
		Value:    "1234",
		GasLimit: "196150",
		Gas:      "15000000000",
		TxLink:   "https://txlink.io/tx?to=0x1111111254fb6c44bac0bed2854e76f90643097d&value=0&data=0x7c02520000000000000000000000000005ad60d9a2f1aa30ba0cdbaf1e0a0a145fbea16f00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000180000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e8930000000000000000000000003af210eec25e95f470fda218c9df9291adb7dcfd000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000e1b274ca648a00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a4b757fed60000000000000000000000009949e1db416a8a05a0cac0ba6ea152ba8729e893000000000000000000000000603c7f932ed1fc6575303d8fb018fdcbb0f39a950000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce820000000000000000001e84801111111254fb6c44bac0bed2854e76f90643097d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000cfee7c08&gaslimit=196150",
	}

	var testDataInput = []testdata{
		{
			name:  "Success Response",
			input: requestInfo,
			want:  validResponse,
			match: true,
		},
		{
			name:  "Failure Response",
			input: requestInfo,
			want:  invalidResponse,
			match: false,
		},
	}

	for _, inputArg := range testDataInput {
		t.Run(inputArg.name, func(t *testing.T) {
			doCalled := false
			var argument *pb.ExchangeSwapRequest
			mockCtrl.RecordCall(mockQuote, "GetExchangeSwap", &inputArg.input).DoAndReturn(
				func(arg *pb.ExchangeSwapRequest) *pb.ExchangeSwapResponse {
					doCalled = true
					argument = arg
					var r = validResponse
					return &r
				})
			if doCalled {
				t.Error("Do() callback called too early.")
			}

			responses := mockCtrl.Call(mockQuote, "GetExchangeSwap", &inputArg.input)

			if !doCalled {
				t.Error("Do() callback not called.")
			}
			if &inputArg.input != argument {
				t.Error("Do callback received wrong argument.")
			}
			if len(responses) != 1 {
				t.Fatalf("Return values from Call: got %d, want 1", len(responses))
			}
			resp := responses[0].(*pb.ExchangeSwapResponse)
			if inputArg.match {
				if resp.Value != inputArg.want.Value {
					t.Errorf("DoAndReturn return value: got %v, want %v", resp.Value, inputArg.want.Value)
				}
			}
			if !inputArg.match {
				if resp.Value == inputArg.want.Value {
					t.Errorf("DoAndReturn return value: got %v, want %v", resp.Value, inputArg.want.Value)
				}
			}
		})
	}

	mockCtrl.Finish()
}
