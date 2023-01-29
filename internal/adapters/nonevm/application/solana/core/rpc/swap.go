package rpc

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	// "bridge-allowance/pkg/jupiter"
	// "fmt"
	// "math/big"
	// "reflect"
	// "strconv"
	// "strings"

	_ "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SolanaManager) GetExchangeTokens(in *pb.ExchangeTokenRequest) (*pb.ExchangeTokenResponse, error) {
	var exchangeTokenData []*pb.ExchangeTokenInfo
	// var url string

	// chainInfo := s.util.GetNonEVMWalletInfo(in.Chain)
	if in.ExchangeType == "jupiter" {

		// url = chainInfo.Solana.JupiterApiTokenListUrl
	} else {
		return nil, status.Errorf(codes.Unavailable, "Exchange type "+in.ExchangeType+" not supported", "Unsupported type")
	}
	// tokenList, err := s.solSwap.GetSolanaExchangeTokenList(url, in.ExchangeType)
	// if err != nil {
	// 	s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
	// 	return &pb.ExchangeTokenResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	// }
	// if in.ExchangeType == "jupiter" {
	// 	// tokenInfo := tokenList.(*jupiter.TokenList)
	// 	for _, data := range *tokenInfo {
	// 		res := pb.ExchangeTokenInfo{
	// 			TokenAddress:  data.Address,
	// 			TokenDecimals: strconv.Itoa(data.Decimals),
	// 			TokenSymbol:   data.Symbol,
	// 			TokenName:     data.Name,
	// 			TokenLogoUrl:  data.LogoURI,
	// 			LogoUrl:       chainInfo.Solana.SolanaLogoUrl,
	// 		}
	// 		exchangeTokenData = append(exchangeTokenData, &res)
	// 	}
	// } else {
	// 	tokenInfo := tokenList.(*jupiter.SolanaTokenList)
	// 	for _, data := range tokenInfo.Tokens {
	// 		//Main Net Chain ID is 101
	// 		if data.ChainId == 101 {
	// 			res := pb.ExchangeTokenInfo{
	// 				TokenAddress:  data.Address,
	// 				TokenDecimals: strconv.Itoa(data.Decimals),
	// 				TokenSymbol:   data.Symbol,
	// 				TokenName:     data.Name,
	// 				TokenLogoUrl:  data.LogoURI,
	// 				LogoUrl:       tokenInfo.LogoURI,
	// 			}
	// 			exchangeTokenData = append(exchangeTokenData, &res)
	// 		}
	// 	}
	// }
	return &pb.ExchangeTokenResponse{
		ExchangeTokens: exchangeTokenData,
	}, nil
}

// func (s *SolanaManager) GetExchangeQuote(in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error) {
// 	if in.ExchangeType != "jupiter" {
// 		s.logger.Errorf("Exchange type %v not supported", in.ExchangeType)
// 		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Unavailable, "Exchange Type Not Supported", "Unavailable")
// 	}
// 	// chainInfo := s.util.GetNonEVMWalletInfo(in.Chain)
// 	tokenList, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{
// 		Chain:        "solana",
// 		ExchangeType: "jupiter",
// 	})
// 	if err != nil {
// 		s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
// 		return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	isSrcTokenAddrExists, srcTokenDecimals := s.helper.CheckTokenListData(strings.ToLower(in.SellToken), "TokenDecimals", reflect.ValueOf(&tokenList.ExchangeTokens), reflect.TypeOf(tokenList.ExchangeTokens), "TokenAddress")
// 	isDstTokenAddrExists, dstTokenDecimals := s.helper.CheckTokenListData(strings.ToLower(in.BuyToken), "TokenDecimals", reflect.ValueOf(&tokenList.ExchangeTokens), reflect.TypeOf(tokenList.ExchangeTokens), "TokenAddress")
// 	fmt.Println(dstTokenDecimals)
// 	var response pb.ExchangeQuoteResponse
// 	if isSrcTokenAddrExists && isDstTokenAddrExists {
// 		sellAmountParam, err := s.helper.ConvertStringFloatToIntWithDecimals(in.SellAmount, srcTokenDecimals)
// 		fmt.Println(sellAmountParam)
// 		if err != nil {
// 			s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
// 			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Error in data type conversion ")
// 		}
// 		// quoteUrl := chainInfo.Solana.JupiterApi + fmt.Sprintf("/quote?inputMint=%s&outputMint=%s&amount=%d&slippage=%s&feeBps=0",
// 		// 	in.SellToken, in.BuyToken, sellAmountParam, in.Slippage)
// 		// quoteRes, err := s.solSwap.GetSolanaExchangeQuote(quoteUrl)
// 		if err != nil {
// 			s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
// 			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 		}
// 		bestQuote := quoteRes.Data[0]
// 		// resAmountWithDecimals, err := s.helper.ConvertStringValueToBigFloat(strconv.Itoa(bestQuote.OutAmountWithSlippage), dstTokenDecimals)
// 		// if err != nil {
// 		// 	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
// 		// 	return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, err.Error(), "Error in data type conversion")
// 		// }
// 		// response.ResAmount = fmt.Sprint(resAmountWithDecimals)
// 		// response.PriceImpact = fmt.Sprintf("%v", bestQuote.PriceImpactPct)
// 		resAmountBigFloat, ok := new(big.Float).SetString(response.ResAmount)
// 		if !ok {
// 			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error", "Internal Error")
// 		}
// 		sellAmountFloat, ok := new(big.Float).SetString(in.SellAmount)
// 		if !ok {
// 			return &pb.ExchangeQuoteResponse{}, status.Errorf(codes.Internal, "SetString: error", "Internal Error")
// 		}
// 		// Calculating Price per token := resAmount/sellAmount
// 		resPricePerToken, _ := new(big.Float).Quo(resAmountBigFloat, sellAmountFloat).Float64()
// 		response.ResPricePerFromToken = fmt.Sprint(resPricePerToken)
// 		// Calculating Price per Token for User = 1/resPricePerToken
// 		resPricePerToToken, _ := new(big.Float).Quo(big.NewFloat(1), big.NewFloat(resPricePerToken)).Float64()
// 		response.ResPricePerToToken = fmt.Sprint(resPricePerToToken)
// 		quoteRate, err := s.GetTokenPrice(&pb.TokenPriceRequest{
// 			Currency: "usd",
// 		})
// 		// Calculating token price to sell : quoteRate*sellAmount
// 		fromTokenPrice, _ := new(big.Float).Mul(big.NewFloat(quoteRate.Price), sellAmountFloat).Float64()
// 		// Calculate token price to buy : quoteRate*buyAmount
// 		toTokenPrice, _ := new(big.Float).Mul(big.NewFloat(quoteRate.Price), resAmountBigFloat).Float64()
// 		response.FromTokenPrice = fmt.Sprint(fromTokenPrice)
// 		response.ToTokenPrice = fmt.Sprint(toTokenPrice)
// 		// Calculating Minimum Receive  resAmount *(1-priceImpact/100)
// 		minimumReceived, _ := new(big.Float).
// 			Mul(resAmountBigFloat, new(big.Float).
// 				Sub(big.NewFloat(1), big.NewFloat(bestQuote.PriceImpactPct))).Float64()
// 		response.MinimumReceived = fmt.Sprint(minimumReceived)
// 		return &response, nil
// 	}
// 	return nil, status.Errorf(codes.Unavailable, "SellToken And/Or BuyToken Not Found in TokenList", "Internal Error")
// }

// func (s *SolanaManager) GetExchangeMultiQuote(request *pb.ExchangeMultiQuoteRequest) (*pb.ExchangeMultiQuoteResponse, error) {
// 	response, err := s.GetExchangeQuote(&pb.ExchangeQuoteRequest{
// 		Chain:        request.Chain,
// 		TakerAddress: request.TakerAddress,
// 		SellToken:    request.MultiChainRequests[0].SellToken,
// 		BuyToken:     request.MultiChainRequests[0].BuyToken,
// 		SellAmount:   request.MultiChainRequests[0].SellAmount,
// 		Slippage:     request.MultiChainRequests[0].Slippage,
// 		ExchangeType: request.ExchangeType,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s.ConvertSingleQuoteToMultiQuoteResponse(response, request.Chain), err

// }

func (s *SolanaManager) ConvertSingleQuoteToMultiQuoteResponse(returnData *pb.ExchangeQuoteResponse, chainName string) *pb.ExchangeMultiQuoteResponse {
	var multiQuoteResponse pb.ExchangeMultiQuoteResponse
	multiQuoteResponse.Chain = chainName
	multiQuoteResponse.MultiChainResponse = append(multiQuoteResponse.MultiChainResponse, &pb.MultiChainResponse{
		PriceImpact:          returnData.PriceImpact,
		ResAmount:            returnData.ResAmount,
		ResPricePerFromToken: returnData.ResPricePerFromToken,
		ResPricePerToToken:   returnData.ResPricePerToToken,
		FromTokenPrice:       returnData.FromTokenPrice,
		ToTokenPrice:         returnData.ToTokenPrice,
		MinimumReceived:      returnData.MinimumReceived,
		ApproveAddress:       returnData.ApproveAddress,
	})
	return &multiQuoteResponse
}

func (s *SolanaManager) ConvertSingleQuoteToMultiSwapResponse(returnData *pb.ExchangeSwapResponse) *pb.ExchangeMultipleSwapResponse {
	//var multiSwapResponse pb.ExchangeMultiSwapResponse
	multiSwapResponse := &pb.ExchangeMultipleSwapResponse{
		Transaction: returnData,
	}
	return multiSwapResponse
}

// func (s *SolanaManager) GetExchangeSwap(in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error) {
// 	var values []string
// 	chainInfo := s.util.GetNonEVMWalletInfo(in.Chain)

// 	if in.ExchangeType != "jupiter" {
// 		s.logger.Errorf("Exchange type %v not supported", in.ExchangeType)
// 		return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Unavailable, "Exchange Type Not Supported", "Unavailable")
// 	}
// 	tokenList, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{
// 		Chain:        "solana",
// 		ExchangeType: "jupiter",
// 	})
// 	if err != nil {
// 		s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
// 		return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 	}
// 	//To Fetch Sell Token and But Token decimal values
// 	isSrcTokenAddrExists, srcTokenDecimals := s.helper.CheckTokenListData(strings.ToLower(in.SellToken), "TokenDecimals", reflect.ValueOf(&tokenList.ExchangeTokens), reflect.TypeOf(tokenList.ExchangeTokens), "TokenAddress")
// 	isDstTokenAddrExists, _ := s.helper.CheckTokenListData(strings.ToLower(in.BuyToken), "TokenDecimals", reflect.ValueOf(&tokenList.ExchangeTokens), reflect.TypeOf(tokenList.ExchangeTokens), "TokenAddress")
// 	if isSrcTokenAddrExists && isDstTokenAddrExists {
// 		sellAmountParam, err := s.helper.ConvertStringFloatToIntWithDecimals(in.SellAmount, srcTokenDecimals)
// 		if err != nil {
// 			s.logger.Errorf("Error for Exchange Swap request  is : %v", err.Error())
// 			return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, err.Error(), "Error in data type conversion")
// 		}
// 		quoteUrl := chainInfo.Solana.JupiterApi + fmt.Sprintf("/quote?inputMint=%s&outputMint=%s&amount=%d&slippage=%s&feeBps=0",
// 			in.SellToken, in.BuyToken, sellAmountParam, in.Slippage)
// 		quoteRes, err := s.solSwap.GetSolanaExchangeQuote(quoteUrl)
// 		if err != nil {
// 			s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
// 			return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 		}
// 		//Taking the best quote route
// 		bestQuote := quoteRes.Data[0]
// 		s.logger.Infof("Quote Route %v", bestQuote)
// 		swapUrl := chainInfo.Solana.JupiterApi + "/swap"
// 		swapRes, err := s.solSwap.GetSolanaExchangeSwap(swapUrl, bestQuote, in.TakerAddress)
// 		if err != nil {
// 			s.logger.Errorf("Error Exchange Tokens  for Solana  is : %v", err.Error())
// 			return &pb.ExchangeSwapResponse{}, status.Errorf(codes.Internal, err.Error(), "Internal Error")
// 		}
// 		v := reflect.ValueOf(swapRes)
// 		if v.Kind() == reflect.Ptr {
// 			v = v.Elem()
// 		}
// 		for i := 0; i < v.NumField(); i++ {
// 			if v.Field(i).String() != "" {
// 				values = append(values, v.Field(i).String())
// 			}
// 		}
// 		if values == nil {
// 			return nil, status.Errorf(codes.Internal, "Internal Error : Empty Data")
// 		}
// 		return &pb.ExchangeSwapResponse{
// 			To:       in.BuyToken,
// 			Value:    fmt.Sprint(bestQuote.OutAmountWithSlippage),
// 			GasLimit: "0",
// 			Gas:      "0",
// 			MultiRouteData: &pb.MultiSwapTxs{
// 				Data: values,
// 			},
// 		}, nil
// 	}
// 	return nil, status.Errorf(codes.Unavailable, "SellToken And/Or BuyToken Not Found in TokenList", "Internal Error")
// }

// func (s *SolanaManager) GetExchangeMultiSwap(request *pb.ExchangeMultiSwapRequest) (*pb.ExchangeMultipleSwapResponse, error) {
// 	response, err := s.GetExchangeSwap(&pb.ExchangeSwapRequest{
// 		Chain:        request.Chain,
// 		TakerAddress: request.TakerAddress,
// 		SellToken:    request.MultiChainRequests[0].SellToken,
// 		BuyToken:     request.MultiChainRequests[0].BuyToken,
// 		SellAmount:   request.MultiChainRequests[0].SellAmount,
// 		Slippage:     request.MultiChainRequests[0].Slippage,
// 		ExchangeType: request.ExchangeType,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s.ConvertSingleQuoteToMultiSwapResponse(response), err
// }
