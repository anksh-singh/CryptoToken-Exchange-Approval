package socket

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/onrik/ethrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type ISocket interface {
	GetChains() (*pb.BridgeChainResponse, error)
	GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error)
	GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error)
	GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error)
	GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error)
}

type Socket struct {
	env          *config.Config
	logger       *zap.SugaredLogger
	httpRequest  utils.IHttpRequest
	helper       *utils.Helpers
	includeDexes string
	tool         string
	logoURL      string
	rpc          map[string]*ethrpc.EthRPC
	coinGecko    coingecko.ICoinGecko
	defaultToken []string
}

func NewSocket(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest, helper *utils.Helpers, coinGecko coingecko.ICoinGecko) *Socket {
	rpc := make(map[string]*ethrpc.EthRPC)
	//utilsManager := utils.NewUtils(logger, env)
	//Do not continue if no EVM configurations are provided
	//logger.Info("Supported EVM chains:", len(env.EVM.Cfg.Wallets))
	//if len(env.EVM.Cfg.Wallets) < 1 {
	//	logger.Fatal("No EVM wallet configurations found")
	//}
	//Initialize EVM RPC configurations
	for i, w := range env.EVM.Cfg.Wallets {
		i++
		rpc[w.ChainName] = ethrpc.New(w.RPC)
	}
	return &Socket{
		env:          env,
		logger:       logger,
		helper:       helper,
		httpRequest:  httpRequest,
		includeDexes: "oneinch",
		tool:         "Socket",
		logoURL:      env.Socket.LogoUrl,
		rpc:          rpc,
		coinGecko:    coinGecko,
		defaultToken: []string{"0x4200000000000000000000000000000000000006", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "0x0000000000000000000000000000000000001010"},
	}
}

func (s *Socket) GetChains() (*pb.BridgeChainResponse, error) {
	url := s.env.Socket.EndPoint + "/supported/chains"
	body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var chainInfo ChainsSocket
	err = json.Unmarshal(body, &chainInfo)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	var chainResponse pb.BridgeChainResponse

	if chainInfo.Result != nil {
		if len(chainInfo.Result) > 0 {
			for _, item := range chainInfo.Result {
				chainResponse.Chains = append(chainResponse.Chains, &pb.BridgeChainInfo{
					Symbol:    s.GetChainSymbol(item.ChainID),
					ChainType: "EVM",
					Name:      item.Name,
					Coin:      item.Currency.Symbol,
					ChainId:   item.ChainID,
					LogoUrl:   item.Icon,
					MainNet:   true,
				})
			}
		} else {
			return nil, status.Errorf(codes.NotFound, "no chains found", "no chains found")
		}
	}
	return &chainResponse, err
}

func (s *Socket) GetChainTokens(request *pb.BridgeChainTokensRequest) (*pb.BridgeChainTokensResponse, error) {
	if request.FromChain != "" && request.ToChain != "" && request.FromToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "From token cannot be empty for this case")
	}
	var chain string
	_, isToChain := s.FindChainNameAndValidate(request.ToChain)
	chain = request.ToChain
	if !isToChain {
		_, isToChain = s.FindChainNameAndValidate(request.FromChain)
		chain = request.FromChain
	}
	if isToChain {
		url := fmt.Sprintf(s.env.Socket.EndPoint+"/token-lists/chain?chainId=%s", chain)
		body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
		if err != nil {
			s.logger.Error(err)
			var ErrorResponse ErrorSocket
			err = json.Unmarshal(body, &ErrorResponse)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
			}
			return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
		}
		var chainTokensInfo ChainTokensSocket
		err = json.Unmarshal(body, &chainTokensInfo)
		if err != nil {
			s.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		var chainTokensResponse pb.BridgeChainTokensResponse
		var nativeTokenInfo pb.BridgeTokens
		if chainTokensInfo.Result != nil {
			if len(chainTokensInfo.Result) > 0 {
				for _, item := range chainTokensInfo.Result {
					if contains(s.defaultToken, strings.ToLower(item.Address)) {
						nativeTokenInfo.TokenAddress = s.FindNativeTokenAddress(request.ToChain, item.Address)
						nativeTokenInfo.TokenDecimals = item.Decimals
						nativeTokenInfo.TokenSymbol = item.Symbol
						nativeTokenInfo.TokenName = item.Name
						nativeTokenInfo.TokenLogoUrl = item.LogoURI
					} else {
						chainTokensResponse.Tokens = append(chainTokensResponse.Tokens, &pb.BridgeTokens{
							TokenAddress:  item.Address,
							TokenDecimals: item.Decimals,
							TokenSymbol:   item.Symbol,
							TokenName:     item.Name,
							TokenLogoUrl:  item.LogoURI,
						})
					}

				}
				chainTokensResponse.Tokens = append([]*pb.BridgeTokens{
					&nativeTokenInfo,
				}, chainTokensResponse.Tokens...)
			} else {
				return nil, status.Errorf(codes.NotFound, "no token found", "no tokens found")
			}
		}
		return &chainTokensResponse, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "chain not served by socket", "Chain not served by Socket")
	}

}

func (s *Socket) GetRouteTransactionData(routeData RoutePayload) (*RouteTransactionResponseSocket, error) {
	buildTxData, err := json.Marshal(routeData)
	txString := string(buildTxData)
	if err != nil {
		s.logger.Errorf("error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error(), "marshaling error")
	}
	url := fmt.Sprintf(s.env.Socket.EndPoint + "/build-tx")

	body, err := s.httpRequest.PostRequestWithHeaders(url, txString, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var transactionRoute RouteTransactionResponseSocket
	err = json.Unmarshal(body, &transactionRoute)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
	}
	return &transactionRoute, err
}

//Not Being Used

func (s *Socket) CheckAllowance(chainId string, owner string, allowanceTarget string, tokenAddress string) (*CheckAllowanceResponse, error) {
	url := fmt.Sprintf(s.env.Socket.EndPoint+"/check-allowance?chainID=%s&owner=%s&allowanceTarget=%s&tokenAddress=%s", chainId, owner, allowanceTarget, tokenAddress)
	body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var checkAllowance CheckAllowanceResponse
	err = json.Unmarshal(body, &checkAllowance)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	return &checkAllowance, err
}

//Not Being Used

func (s *Socket) GetApprovalTransaction(chainId string, owner string, allowanceTarget string, tokenAddress string, amount string) (*ApprovalTransactionResponse, error) {
	url := fmt.Sprintf(s.env.Socket.EndPoint+"/approval/build-tx?chainID=%s&owner=%s&allowanceTarget=%s&tokenAddress=%s&amount=%s", chainId, owner, allowanceTarget, tokenAddress, amount)
	body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var approvalTransaction ApprovalTransactionResponse
	err = json.Unmarshal(body, &approvalTransaction)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	return &approvalTransaction, err
}

func (s *Socket) GetBridgeStatus(transaction string, fromChainId string, toChainId string) (*TransactionStatusSocket, error) {
	url := fmt.Sprintf(s.env.Socket.EndPoint+"/bridge-status?transactionHash=%s&fromChainId=%s&toChainId=%s", transaction, fromChainId, toChainId)
	body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var statusTransaction TransactionStatusSocket
	err = json.Unmarshal(body, &statusTransaction)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	return &statusTransaction, err
}

func (s *Socket) GetQuoteSocket(request *pb.BridgeQuoteRequest) (*QuoteSocketResponse, error) {
	url := fmt.Sprintf(s.env.Socket.EndPoint+"/quote?fromChainId=%s&fromTokenAddress=%s&toChainId=%s&toTokenAddress=%s&fromAmount=%s&userAddress=%s&recipient=%s&uniqueRoutesPerBridge=%v&sort=%s&singleTxOnly=%v",
		request.FromChain, request.FromToken, request.ToChain, request.ToToken, request.FromAmount, request.FromAddress, request.ToAddress, true, "output", true)
	body, err := s.httpRequest.GetRequestWithHeaders(url, "API-KEY", s.env.Socket.APIKey)
	if err != nil {
		s.logger.Error(err)
		var ErrorResponse ErrorSocket
		err = json.Unmarshal(body, &ErrorResponse)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
		}
		return nil, status.Errorf(codes.Internal, string(body), ErrorResponse.Error.Message)
	}
	var quote QuoteSocketResponse
	err = json.Unmarshal(body, &quote)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "unmarshall error")
	}
	return &quote, err
}

func (s *Socket) GetQuote(request *pb.BridgeQuoteRequest) (*pb.BridgeQuoteResponse, error) {
	fromChainName, isFromChain := s.FindChainNameAndValidate(request.FromChain)
	toChainName, isToChain := s.FindChainNameAndValidate(request.ToChain)
	if isFromChain && isToChain {
		if contains(s.defaultToken, strings.ToLower(request.FromToken)) {
			request.FromToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(s.defaultToken, strings.ToLower(request.ToToken)) {
			request.ToToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		quote, err := s.GetQuoteSocket(request)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		var quoteResponse pb.BridgeQuoteResponse
		var toolDetails pb.ToolDetails
		var estimate pb.Estimate
		//var bridgeFee pb.BridgeFee
		//fromQuoteRate, err := s.coinGecko.GetTokenExchangeForContract(fromChainName, request.FromToken, "usd")
		//if err != nil {
		//	fromQuoteRate.Price = 0
		//}
		//toQuoteRate, err := s.coinGecko.GetTokenExchangeForContract(toChainName, request.ToToken, "usd")
		//if err != nil {
		//	toQuoteRate.Price = 0
		//}

		var toQuote float64
		var fromQuote float64
		// Calculate FromToken Price
		if contains(s.defaultToken, strings.ToLower(request.FromToken)) {
			quoteRate, err := s.coinGecko.GetTokenExchange("usd", fromChainName)
			if err != nil {
				//o.logger.Errorf("Error for Exchange Quote Price for token  request  is : %v", err.Error())
				fromQuote = 0
			} else {
				fromQuote = quoteRate.Price
			}

		} else {
			quoteRate, err := s.coinGecko.GetTokenExchangeForContract(fromChainName, strings.ToLower(request.FromToken), "usd")
			if err != nil {
				//o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				fromQuote = 0
			} else {
				fromQuote = quoteRate.Price
			}
		}
		// Calculate ToToken Price
		if contains(s.defaultToken, strings.ToLower(request.ToToken)) {
			quoteRate, err := s.coinGecko.GetTokenExchange("usd", toChainName)

			if err != nil {
				//o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				toQuote = 0
			} else {
				toQuote = quoteRate.Price
			}

		} else {
			quoteRate, err := s.coinGecko.GetTokenExchangeForContract(toChainName, strings.ToLower(request.ToToken), "usd")
			if err != nil {
				//o.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
				toQuote = 0
			} else {
				toQuote = quoteRate.Price
			}
		}
		quoteResponse.Tool = s.tool
		toolDetails.Key = s.tool
		toolDetails.Name = s.tool
		toolDetails.LogoUrl = s.logoURL
		if quote.Result.Routes != nil {
			if len(quote.Result.Routes) > 0 {
				estimate.ToAmount = quote.Result.Routes[0].ToAmount
				estimate.FromAmount = quote.Result.Routes[0].FromAmount
				estimate.ToAmountUsd = s.CalculateQuotePrice(quote.Result.Routes[0].ToAmount, quote.Result.ToAsset.Decimals, toQuote)
				estimate.FromAmountUsd = s.CalculateQuotePrice(quote.Result.Routes[0].FromAmount, quote.Result.FromAsset.Decimals, fromQuote)
				estimate.ApproveAddress = quote.Result.Routes[0].UserTxs[0].ApprovalData.AllowanceTarget
				if estimate.ApproveAddress == "" {
					estimate.ApproveAddress = request.ToToken
				}
				estimate.ToTokenDecimals = quote.Result.ToAsset.Decimals
				estimate.FromTokenDecimals = quote.Result.FromAsset.Decimals
				if quote.Result.Routes[0].UserTxs != nil {
					estimate.ToAmount = quote.Result.Routes[0].ToAmount
					estimate.ToAmountMin = quote.Result.Routes[0].UserTxs[0].Steps[len(quote.Result.Routes[0].UserTxs[0].Steps)-1].MinAmountOut
					estimate.ExecutionDuration = float64(quote.Result.Routes[0].ServiceTime)
				}
				quoteResponse.Estimate = &estimate
				//quoteResponse.BridgeFee = &bridgeFee
			} else {
				return nil, status.Errorf(codes.NotFound, "no routes found", "No routes found")
			}
		}

		quoteResponse.ToolDetails = &toolDetails
		return &quoteResponse, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "chain not served by Socket", "Chain not served by Socket")
	}

}

func (s *Socket) GetTransaction(request *pb.BridgeTransactionRequest) (*pb.BridgeTransactionResponse, error) {
	// calculate Quote
	fromChainNameFrontier, isFromChain := s.FindChainNameAndValidate(request.FromChain)
	_, isToChain := s.FindChainNameAndValidate(request.ToChain)
	if isFromChain && isToChain {
		if contains(s.defaultToken, strings.ToLower(request.FromToken)) {
			request.FromToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		if contains(s.defaultToken, strings.ToLower(request.ToToken)) {
			request.ToToken = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		}
		quote, err := s.GetQuoteSocket(&pb.BridgeQuoteRequest{
			ToToken:     request.ToToken,
			FromAmount:  request.FromAmount,
			ToChain:     request.ToChain,
			FromChain:   request.FromChain,
			FromToken:   request.FromToken,
			FromAddress: request.FromAddress,
			ToAddress:   request.ToAddress,
			Chain:       request.Chain,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		// Calculate Route from Quote
		var route Route
		//var gasLimit int64
		if quote.Result.Routes != nil {
			if len(quote.Result.Routes) > 0 {
				route.RouteID = quote.Result.Routes[0].RouteID
				route.IsOnlySwapRoute = quote.Result.Routes[0].IsOnlySwapRoute
				route.FromAmount = quote.Result.Routes[0].FromAmount
				if quote.Result.Routes[0].UserTxs != nil {
					if len(quote.Result.Routes[0].UserTxs) > 0 {
						route.UserTxs = quote.Result.Routes[0].UserTxs[0:1]
						//gasLimit = quote.Result.Routes[0].UserTxs[0].GasFees.GasLimit
						//route.UserTxs[0].Steps = route.UserTxs[0].Steps[0:1]
					}
				}
				route.ToAmount = quote.Result.Routes[0].ToAmount
				route.MaxServiceTime = quote.Result.Routes[0].MaxServiceTime
				route.Recipient = quote.Result.Routes[0].Recipient
				route.Sender = quote.Result.Routes[0].Sender
				route.ServiceTime = quote.Result.Routes[0].ServiceTime
				route.TotalGasFeesInUsd = quote.Result.Routes[0].TotalGasFeesInUsd
				route.TotalUserTx = quote.Result.Routes[0].TotalUserTx
				route.UsedBridgeNames = quote.Result.Routes[0].UsedBridgeNames
			} else {
				return nil, status.Errorf(codes.NotFound, "no routes found", "No routes found")
			}

		}
		// Calculate TransactionData
		var routePayload RoutePayload
		routePayload.Route = route
		transactionData, err := s.GetRouteTransactionData(routePayload)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		var transactionResponse pb.BridgeTransactionResponse
		var toolDetails pb.ToolDetails
		var transactionRequest pb.TransactionRequest
		transactionResponse.Tool = s.tool
		toolDetails.Key = s.tool
		toolDetails.Name = s.tool
		toolDetails.LogoUrl = s.logoURL
		transactionRequest.Data = transactionData.Result.TxData
		transactionRequest.To = transactionData.Result.TxTarget
		transactionRequest.Value = strconv.FormatFloat(s.helper.ConvertHexToFloat64(transactionData.Result.Value), 'f', -1, 64)

		// Calculate Gas Limit
		gasLimitData := transactionRequest.Data
		gaslimitTo := request.ToAddress
		gaslimitFrom := request.FromAddress
		transaction := ethrpc.T{
			From: gaslimitFrom,
			To:   gaslimitTo,
			Data: gasLimitData,
		}
		gasLimit, err := s.rpc[fromChainNameFrontier].EthEstimateGas(transaction)
		if err != nil {
			s.logger.Error("Error fetching gas estimate")
			gasLimit = 8000000 //TODO:To be refactored
			err = nil
		}
		gasLimit = gasLimit * 31
		transactionRequest.GasLimit = strconv.Itoa(gasLimit)
		transactionResponse.TransactionRequest = &transactionRequest
		transactionResponse.ToolDetails = &toolDetails
		transactionResponse.ToolDetails.LogoUrl = s.logoURL
		return &transactionResponse, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "chain not served by Socket", "Chain not served by Socket")
	}

}

func (s *Socket) GetTransactionStatus(request *pb.BridgeTransactionStatusRequest) (*pb.BridgeTransactionStatusResponse, error) {
	_, isFromChain := s.FindChainNameAndValidate(request.FromChain)
	_, isToChain := s.FindChainNameAndValidate(request.ToChain)
	if isFromChain && isToChain {
		statusMsg, err := s.GetBridgeStatus(request.TxHash, request.FromChain, request.ToChain)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		var response pb.BridgeTransactionStatusResponse
		response.Status = statusMsg.Result.DestinationTxStatus
		response.TxHash = statusMsg.Result.DestinationTransactionHash
		response.Msg = "Swap request done"
		response.IsSuccess = true
		return &response, err
	} else {
		return nil, status.Errorf(codes.Unknown, errors.New("chain not served by Socket").Error(), "Chain not served by Socket")
	}

}

func (s *Socket) CalculateQuotePrice(amount string, tokenDecimal int64, quoteRate float64) string {
	amountRate := s.helper.CalculateRateWithDecimal(amount, tokenDecimal)
	return strconv.FormatFloat(amountRate*quoteRate, 'f', -1, 64)
}

func (s *Socket) FindChainNameAndValidate(fromChainId string) (string, bool) {
	isSocketChain := false
	for _, item := range s.env.EVM.Cfg.Wallets {
		if strings.ToLower(fromChainId) == strings.ToLower(item.ChainName) {
			if item.SocketSupport == true {
				isSocketChain = true
				return item.ChainName, isSocketChain
			}

		} else {
			chainId := strconv.FormatInt(int64(item.ChainID), 10)
			if fromChainId == chainId {
				if item.SocketSupport == true {
					isSocketChain = true
					return item.ChainName, isSocketChain
				}
			}
		}

	}
	return "", isSocketChain
}

func (s *Socket) GetChainSymbol(chainId int64) string {
	for _, item := range s.env.EVM.Cfg.Wallets {
		if chainId == int64(item.ChainID) {
			return item.Nomenclature.ChainNameShort
		}
	}
	return ""
}
func (s *Socket) FindNativeTokenAddress(toChainId string, tokenAddress string) string {
	for _, item := range s.env.EVM.Cfg.Wallets {
		if strings.ToLower(toChainId) == strings.ToLower(item.ChainName) {
			if item.NativeTokenInfo.Address != "" {
				return item.NativeTokenInfo.Address
			}
		} else {
			chainId := strconv.FormatInt(int64(item.ChainID), 10)
			if toChainId == chainId {
				if item.NativeTokenInfo.Address != "" {
					return item.NativeTokenInfo.Address
				}
			}
		}
	}
	return tokenAddress
}
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
