package debridge

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type IdeBridge interface {
	GetChains() (*pb.BridgeChainResponse, error)
	GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error)
	GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error)
	GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error)
	//GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error)
}

type DeBridge struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
	helper      *utils.Helpers
	tool        string
	coinGecko   coingecko.ICoinGecko
}

func NewDeBridgeService(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *DeBridge {
	return &DeBridge{
		env:         env,
		logger:      logger,
		httpRequest: httpRequest,
		helper:      helper,
		tool:        "deBridge",
		coinGecko:   coinGecko,
	}
}

func (d *DeBridge) GetChains() (*pb.BridgeChainResponse, error) {
	var chainResponse pb.BridgeChainResponse
	var err error
	for _, item := range d.env.EVM.Cfg.Wallets {
		if item.DebridgeSupport == true {
			chainResponse.Chains = append(chainResponse.Chains, &pb.BridgeChainInfo{
				Symbol:    item.ChainInfo.ID,
				ChainType: "EVM",
				Name:      item.ChainInfo.Name,
				Coin:      item.ChainInfo.NativeTokenID,
				ChainId:   int64(item.ChainID),
				LogoUrl:   item.ChainInfo.LogoUrl,
				MainNet:   true,
			})
		}
	}
	return &chainResponse, err
}

func (d *DeBridge) GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	if request.FromChain != "" && request.ToChain != "" && request.FromToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "From token cannot be empty for this case")
	}
	var chain string
	toChain := d.FindChainKey(request.ToChain)
	_, isChain := d.FindChainNameAndValidate(toChain)
	chain = toChain
	if !isChain {
		fromChain := d.FindChainKey(request.FromChain)
		_, isChain = d.FindChainNameAndValidate(fromChain)
		chain = fromChain
	}
	var chainTokensResponse pb.BridgeChainTokensResponse
	var tokenAddress string
	if isChain {
		url := fmt.Sprintf(d.env.Swap.OneInchEndpoint+"%s/tokens", chain)
		body, err := d.httpRequest.GetRequest(url)
		if err != nil {
			d.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			d.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
		}
		tokens := result["tokens"].(map[string]interface{})
		var nativeTokenInfo pb.BridgeTokens
		for _, token := range tokens {
			token := token.(map[string]interface{})
			if token["logoURI"] == nil {
				if token["address"].(string) == "0x3642c0680329ae3e103e2b5ab29ddfed4d43cbe5" {
					token["logoURI"] = "https://plenny.link/wp-content/uploads/2020/11/plenny_logo.png"
				} else if token["address"].(string) == "0xde903e2712288a1da82942dddf2c20529565ac30" {
					token["logoURI"] = "https://swapr.eth.link/static/media/swapr_white_no_badge.e8d9cb38a2a453eafc00c6cbe419d942.svg"
				} else {
					token["logoURI"] = ""
				}
			}
			tokenAddress = token["address"].(string)
			tags := token["tags"].([]interface{})
			for _, item := range tags {
				if item.(string) == "native" {
					nativeTokenInfo.TokenAddress = d.GetFrontierNativeToken(toChain, tokenAddress)
					nativeTokenInfo.TokenDecimals = int64(token["decimals"].(float64))
					nativeTokenInfo.TokenSymbol = token["symbol"].(string)
					nativeTokenInfo.TokenName = token["name"].(string)
					nativeTokenInfo.TokenLogoUrl = token["logoURI"].(string)
				}
			}
			chainTokensResponse.Tokens = append(chainTokensResponse.Tokens, &pb.BridgeTokens{
				TokenAddress:  tokenAddress,
				TokenDecimals: int64(token["decimals"].(float64)),
				TokenSymbol:   token["symbol"].(string),
				TokenName:     token["name"].(string),
				TokenLogoUrl:  token["logoURI"].(string),
			})
		}
		chainTokensResponse.Tokens = append([]*pb.BridgeTokens{
			&nativeTokenInfo,
		}, chainTokensResponse.Tokens...)
		return &chainTokensResponse, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "Chain not served by DeBridge", "Chain not served by DeBridge")
	}
}

func (d *DeBridge) GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error) {
	fromChain := d.FindChainKey(request.FromChain)
	toChain := d.FindChainKey(request.ToChain)
	fromToken := d.GetDeBridgeNativeToken(fromChain, request.FromToken)
	toToken := d.GetDeBridgeNativeToken(toChain, request.ToToken)
	url := fmt.Sprintf(d.env.Debridge.EndPoint+"/estimation?srcChainId=%s&srcChainTokenIn=%s&srcChainTokenInAmount=%s&slippage=0&dstChainId=%s&dstChainTokenOut=%s&executionFeeAmount=auto&affiliatePercent=0",
		fromChain, fromToken, request.FromAmount, toChain, toToken)

	body, err := d.httpRequest.GetRequest(url)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	var quoteResponse pb.BridgeQuoteResponse
	var toolDetails pb.ToolDetails
	var estimate pb.Estimate
	var bridgeFee pb.BridgeFee
	var executionDuration float64
	var fromAmountInUSD float64
	var toAmountInUSD float64

	if request.FromChain == "1" {
		executionDuration = float64(ethEstimatedDuration)
	} else if request.FromChain == "56" {
		executionDuration = float64(bscEstimatedDuration)
	} else if request.FromChain == "137" {
		executionDuration = float64(maticEstimatedDuration)
	} else if request.FromChain == "42161" {
		executionDuration = float64(arbitrumEstimatedDuration)
	} else if request.FromChain == "43114" {
		executionDuration = float64(avalancheEstimatedDuration)
	}

	fromAmountInUSD = d.GetTokenAmountInUSD(fromChain, request.FromToken)
	toAmountInUSD = d.GetTokenAmountInUSD(toChain, request.ToToken)
	tokenData := d.GetNativeTokenInfo(fromChain)
	bridgeFees, err := d.GetProtocolFee(fromChain)
	bridgeFeeInDecimal := d.GetAmountUSDInDecimal(bridgeFees, tokenData.Decimals)
	bridgeFeeInUSD, err := d.GetProtocolFeeInUSD(request)
	bridgeFeeInFloat, err := strconv.ParseFloat(bridgeFeeInDecimal, 64)
	bridgeFeeUSDInFloat, err := strconv.ParseFloat(bridgeFeeInUSD, 64)
	bridgeFeeUSD := bridgeFeeInFloat * bridgeFeeUSDInFloat
	inAmount, err := strconv.ParseFloat(quote.Estimation.SrcChainTokenIn.Amount, 64)
	outAmount, err := strconv.ParseFloat(quote.Estimation.DstChainTokenOut.Amount, 64)
	inAmountUSD := inAmount * fromAmountInUSD
	outAmountUSD := outAmount * toAmountInUSD
	fromAmountInDecimal := d.GetAmountUSDInDecimal(strconv.FormatFloat(inAmountUSD, 'f', -1, 64), int64(quote.Estimation.SrcChainTokenIn.Decimals))
	toAmountInDecimal := d.GetAmountUSDInDecimal(strconv.FormatFloat(outAmountUSD, 'f', -1, 64), int64(quote.Estimation.DstChainTokenOut.Decimals))
	tokenAddress := d.GetFrontierNativeToken(fromChain, nativeTokenAddress)
	quoteResponse.Tool = d.tool
	toolDetails.Key = d.tool
	toolDetails.Name = d.tool
	toolDetails.LogoUrl = d.env.Debridge.LogoUrl
	estimate.FromAmount = quote.Estimation.SrcChainTokenIn.Amount
	estimate.FromTokenDecimals = int64(quote.Estimation.SrcChainTokenIn.Decimals)
	estimate.ToAmount = quote.Estimation.DstChainTokenOut.Amount
	estimate.ToAmountMin = quote.Estimation.DstChainTokenOut.MinAmount
	estimate.ToTokenDecimals = int64(quote.Estimation.DstChainTokenOut.Decimals)
	estimate.ExecutionDuration = executionDuration
	estimate.FromAmountUsd = fromAmountInDecimal
	estimate.ToAmountUsd = toAmountInDecimal
	estimate.ApproveAddress = quote.Tx.AllowanceTarget
	bridgeFee.ContractAddress = tokenAddress
	bridgeFee.Symbol = tokenData.Symbol
	bridgeFee.Amount = bridgeFees
	bridgeFee.TokenDecimals = tokenData.Decimals
	bridgeFee.AmountUsd = strconv.FormatFloat(bridgeFeeUSD, 'f', -1, 64)
	quoteResponse.ToolDetails = &toolDetails
	quoteResponse.Estimate = &estimate
	quoteResponse.BridgeFee = &bridgeFee
	return &quoteResponse, err
}

func (d *DeBridge) GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error) {
	fromChain := d.FindChainKey(request.FromChain)
	toChain := d.FindChainKey(request.ToChain)
	fromToken := d.GetDeBridgeNativeToken(fromChain, request.FromToken)
	toToken := d.GetDeBridgeNativeToken(toChain, request.ToToken)
	url := fmt.Sprintf(d.env.Debridge.EndPoint+"/transaction?srcChainId=%s&srcChainTokenIn=%s&srcChainTokenInAmount=%s&slippage=1&dstChainId=%s&dstChainTokenOut=%s&dstChainTokenOutRecipient=%s&dstChainFallbackAddress=%s&referralCode=4467",
		fromChain, fromToken, request.FromAmount, toChain, toToken, request.ToAddress, request.FromAddress)
	body, err := d.httpRequest.GetRequest(url)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var transaction Transaction
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		d.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	gasLimit, err := d.GetEstimatedGas(request, transaction.Tx.Data)
	var transactionResponse pb.BridgeTransactionResponse
	var toolDetails pb.ToolDetails
	var transactionRequest pb.TransactionRequest

	transactionResponse.Tool = d.tool
	toolDetails.Key = d.tool
	toolDetails.Name = d.tool
	toolDetails.LogoUrl = d.env.Debridge.LogoUrl
	transactionRequest.Data = transaction.Tx.Data
	transactionRequest.To = transaction.Tx.To
	transactionRequest.Value = transaction.Tx.Value
	transactionRequest.GasLimit = gasLimit
	transactionResponse.ToolDetails = &toolDetails
	transactionResponse.TransactionRequest = &transactionRequest
	return &transactionResponse, err
}

//func (d *DeBridge) GetDebridgeStatus(transaction string) (*TransactionStatus, error) {
//	url := fmt.Sprintf(d.env.Debridge.TransactionEndPoint+"/Transactions/GetFullSubmissionInfo?filter=%s&filterType=1", transaction)
//	body, err := d.httpRequest.GetRequest(url)
//	if err != nil {
//		d.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	var transactionStatus TransactionStatus
//	err = json.Unmarshal(body, &transactionStatus)
//	if err != nil {
//		d.logger.Error(err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//	return &transactionStatus, err
//}
//
//func (d *DeBridge) GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error) {
//	toChain := d.FindChainKey(request.ToChain)
//	fromChain := d.FindChainKey(request.FromChain)
//	_, isFromChain := d.FindChainNameAndValidate(fromChain)
//	_, isToChain := d.FindChainNameAndValidate(toChain)
//	if isFromChain && isToChain {
//		statusMsg, err := d.GetDebridgeStatus(request.TxHash)
//		if err != nil {
//			d.logger.Error(err)
//			return nil, status.Errorf(codes.Internal, "Error getting transaction status from DeBridge")
//		}
//		var response pb.BridgeTransactionStatusResponse
//		fromChainKey := d.FindChainKey(fromChain)
//		toChainKey := d.FindChainKey(toChain)
//		fromChain, err := strconv.Atoi(fromChainKey)
//		toChain, err := strconv.Atoi(toChainKey)
//		if statusMsg.Send.IsExecuted == true && statusMsg.Send.EventOriginChainID == fromChain && statusMsg.Send.ChainToID == toChain {
//			response.Status = "Success"
//			response.TxHash = statusMsg.Send.TransactionHash
//			response.Msg = "Swap request done"
//			response.IsSuccess = statusMsg.Send.IsExecuted
//		} else if statusMsg.Send.IsExecuted == false && statusMsg.Send.EventOriginChainID == fromChain && statusMsg.Send.ChainToID == toChain {
//			response.Status = "Processing request"
//			response.TxHash = request.TxHash
//			response.Msg = "Swap request not yet completed"
//			response.IsSuccess = statusMsg.Send.IsExecuted
//		} else {
//			response.Status = "Failed"
//			response.TxHash = request.TxHash
//			response.Msg = "Swap request failed"
//			response.IsSuccess = statusMsg.Send.IsExecuted
//		}
//		return &response, err
//	} else {
//		return nil, status.Errorf(codes.Unknown, "Chain not served by DeBridge")
//	}
//}
