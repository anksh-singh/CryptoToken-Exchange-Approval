package swap

import (
	"bridge-allowance/config"
	"bridge-allowance/internal/adapters/evm/application/core/rpc"
	"bridge-allowance/internal/common"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/builtin/erc20"
	"github.com/umbracle/ethgo/contract"
	ethgoJsonRPC "github.com/umbracle/ethgo/jsonrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IExchangeSwap interface {
	GetExchangeTokens(in *pb.ExchangeTokenRequest) (*pb.ExchangeTokenResponse, error)
	GetExchangeQuote(in *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error)
	GetExchangeMultiQuote(request *pb.ExchangeMultiQuoteRequest) (*pb.ExchangeMultiQuoteResponse, error)
	GetExchangeSwap(in *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error)
	GetExchangeMultiSwap(in *pb.ExchangeMultiSwapRequest) (*pb.ExchangeMultipleSwapResponse, error)
	GetFreeTradeCount(in *pb.FreeTradeCountRequest) (*pb.FreeTradeCountResponse, error)
	GetExchangeSignature(in *pb.ExchangeSignatureRequest) (*pb.ExchangeSignatureResponse, error)
	GetExchangeSwapExecute(in *pb.ExchangeSwapExecuteRequest) (*pb.ExchangeSwapExecuteResponse, error)
	GetTokenApproval(in *pb.TokenApprovalRequest) (*pb.TokenApprovalResponse, error)
}

type Swap struct {
	config   *config.Config
	logger   *zap.SugaredLogger
	services common.Services
	util     *utils.UtilConf
	helper   *utils.Helpers
	evmCore  rpc.EvmCore
	ethgoRpc map[string]*ethgoJsonRPC.Client
}

func NewSwap(config *config.Config, logger *zap.SugaredLogger, services common.Services) *Swap {
	ethgoRpc := make(map[string]*ethgoJsonRPC.Client)
	logger.Info("Supported EVM chains:", len(config.EVM.Cfg.Wallets))
	helper := utils.Helpers{}
	utilsManager := utils.NewUtils(logger, config)
	evmCoreManager := rpc.NewEVMCore(config, logger, services)
	if len(config.EVM.Cfg.Wallets) < 1 {
		logger.Fatal("No EVM wallet configurations found")
	}
	//Initialize EVM RPC configurations
	for i, w := range config.EVM.Cfg.Wallets {
		i++
		var err error
		//This is a duplicate initialization RPC client
		//TODO:Move RPC calls a new Utility
		ethgoRpc[w.ChainName], err = ethgoJsonRPC.NewClient(w.RPC)
		if err != nil {
			logger.Errorf(err.Error())
			logger.Fatalf("Error initializing go RPC client for `%s` chain", w.ChainName)
		}
	}
	return &Swap{
		config:   config,
		logger:   logger,
		services: services,
		util:     utilsManager,
		helper:   &helper,
		ethgoRpc: ethgoRpc,
		evmCore:  evmCoreManager,
	}
}

// getEthTokenDecimals retrieve token decimals from a contract address
func (s *Swap) getEthTokenDecimals(contractAddr string, chain string) (int, error) {
	//TODO: [Enhancement] Cache token decimals for faster retrieval
	data := erc20.NewERC20(ethgo.HexToAddress(contractAddr), contract.WithJsonRPC(s.ethgoRpc[chain].Eth()))
	decimals, err := data.Decimals()
	if err != nil {
		s.logger.Info("Error while fetching token decimals", err.Error())
		//Defaults to 18 decimals
		decimals = 18
	}
	return int(decimals), nil
}

func (s *Swap) GetExchangeTokens(request *pb.ExchangeTokenRequest) (*pb.ExchangeTokenResponse, error) {
	walletInfo := s.util.GetWalletInfo(request.Chain)
	source := s.GetSwapSource(request.Chain)
	exchangeLogoUrl := ""
	tokenSource := ""
	exchangeTokenUrl := ""
	nativeTokenInfo := &pb.ExchangeTokenInfo{
		TokenAddress:  walletInfo.NativeTokenInfo.Address,
		TokenDecimals: walletInfo.NativeTokenInfo.Decimals,
		TokenSymbol:   walletInfo.NativeTokenInfo.Symbol,
		TokenName:     walletInfo.NativeTokenInfo.Name,
		TokenLogoUrl:  walletInfo.NativeTokenInfo.LogoURI,
		LogoUrl:       exchangeLogoUrl,
	}
	switch request.ExchangeType {
	case "0x":
		if source.ChainName == request.Chain && source.ZeroxSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.ZeroxLogoUrl
			tokenSource = source.ZeroxSwapConfig.TokenSource
			exchangeTokenUrl = source.ZeroxSwapConfig.ExchangeTokenUrl
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "dodo":
		if source.ChainName == request.Chain && source.DodoSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.DodoLogoUrl
			exchangeTokenUrl = source.DodoSwapConfig.ExchangeTokenUrl
			tokenSource = source.DodoSwapConfig.TokenSource
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "lifi":
		if source.ChainName == request.Chain && source.LiFiSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.LIFILogoUrl
			res, err := s.services.LiFi.GetExchangeTokens(&pb.BridgeChainTokensRequest{
				Chain:     request.Chain,
				FromChain: "",
				ToChain:   "",
			})
			if err != nil {
				return nil, err
			}
			responseStruct := pb.ExchangeTokenResponse{}
			for _, item := range res.Tokens {
				exchangeTokenInfo := pb.ExchangeTokenInfo{
					TokenAddress:  item.TokenAddress,
					TokenDecimals: fmt.Sprint(item.TokenDecimals),
					TokenLogoUrl:  item.TokenLogoUrl,
					TokenSymbol:   item.TokenSymbol,
					TokenName:     item.TokenName,
					LogoUrl:       exchangeLogoUrl,
				}
				responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, &exchangeTokenInfo)
			}
			if nativeTokenInfo.TokenAddress != "" {
				addressExists, _ := s.helper.CheckTokenListData(nativeTokenInfo.TokenAddress, "TokenAddress", reflect.ValueOf(&responseStruct.ExchangeTokens), reflect.TypeOf(responseStruct.ExchangeTokens), "TokenAddress")
				if !addressExists {
					responseStruct.ExchangeTokens = append(responseStruct.ExchangeTokens, nativeTokenInfo)
				}
			}
			return &responseStruct, nil
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "1inch":
		if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.OneInchLogoUrl
			return s.services.OneInch.GetExchangeTokens(walletInfo, nativeTokenInfo)
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "dzap":
		if source.ChainName == request.Chain && source.DZapSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.DZapLogoUrl
			return s.services.DZap.GetExchangeTokens(walletInfo)
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "zeroswap":
		if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
			exchangeLogoUrl = s.config.Swap.ZeroSwapLogoUrl
			return s.services.ZeroSwap.GetExchangeTokens(source)
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "cowswap":
		if source.ChainName == request.Chain && source.CowSwapConfig.IsSupported == true {
			return s.services.CowSwap.GetExchangeTokens(source, walletInfo)
		} else {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	default:
		for _, proxy := range s.config.Proxies.ExchangeTypes {
			if request.ExchangeType == proxy {
				return s.exchangeTokensV1Proxy(request, s.config.PROXIES_ENDPOINT)
			}
		}
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType))
	}
	switch tokenSource {
	case "cocoswap":
		res, err := s.services.CocoSwapTokenExchange.GetExchangeTokens(exchangeTokenUrl, exchangeLogoUrl, source.ChainId)
		if err != nil {
			return nil, err
		}
		if nativeTokenInfo.TokenAddress != "" {
			// Check that native exists in the list
			addressExists, _ := s.helper.CheckTokenListData(nativeTokenInfo.TokenAddress, "TokenAddress", reflect.ValueOf(&res.ExchangeTokens), reflect.TypeOf(res.ExchangeTokens), "TokenAddress")

			if !addressExists {
				// apply prepend logic
				nativeTokenInfo.LogoUrl = exchangeLogoUrl
				res.ExchangeTokens = append([]*pb.ExchangeTokenInfo{
					nativeTokenInfo,
				}, res.ExchangeTokens...)
				//res.ExchangeTokens = append(res.ExchangeTokens, nativeTokenInfo)
			}
		}
		return res, nil
	case "ubeswap":
		res, err := s.services.UniSwapTokenExchange.GetExchangeTokens(exchangeTokenUrl, exchangeLogoUrl, source.ChainId)
		if err != nil {
			return nil, err
		}
		if nativeTokenInfo.TokenAddress != "" {
			addressExists, _ := s.helper.CheckTokenListData(nativeTokenInfo.TokenAddress, "TokenAddress", reflect.ValueOf(&res.ExchangeTokens), reflect.TypeOf(res.ExchangeTokens), "TokenAddress")
			if !addressExists {
				nativeTokenInfo.LogoUrl = exchangeLogoUrl
				res.ExchangeTokens = append([]*pb.ExchangeTokenInfo{
					nativeTokenInfo,
				}, res.ExchangeTokens...)
				///res.ExchangeTokens = append(res.ExchangeTokens, nativeTokenInfo)
			}
		}
		return res, nil
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", tokenSource), "Chain not supported")
	}
}

func (s *Swap) GetExchangeQuote(request *pb.ExchangeQuoteRequest) (*pb.ExchangeQuoteResponse, error) {
	chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeQuoteResponse{}, err
	}
	source := s.GetSwapSource(request.Chain)
	walletInfo := s.util.GetWalletInfo(request.Chain)
	_, srcTokenDecimals := CheckExistsWithValueData(strings.ToLower(request.SellToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")
	_, dstTokenDecimals := CheckExistsWithValueData(strings.ToLower(request.BuyToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")
	if srcTokenDecimals == "" {
		decimals, _ := s.getEthTokenDecimals(request.SellToken, request.Chain)
		srcTokenDecimals = strconv.Itoa(decimals)
	}
	if dstTokenDecimals == "" {
		decimals, err := s.getEthTokenDecimals(request.BuyToken, request.Chain)
		if err != nil {
			s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}
		dstTokenDecimals = strconv.Itoa(decimals)
	}
	if srcTokenDecimals != "" && dstTokenDecimals != "" {
		switch request.ExchangeType {
		case "0x":
			if source.ChainName == request.Chain && source.ZeroxSwapConfig.IsSupported == true {
				return s.services.ZeroX.GetExchangeQuote(request.SellToken,
					request.BuyToken, request.SellAmount, source.ZeroxSwapConfig.EndPoint, request.Chain, srcTokenDecimals, dstTokenDecimals, request.Slippage)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		case "dodo":
			if source.ChainName == request.Chain && source.DodoSwapConfig.IsSupported == true {
				return s.services.DodoSwap.GetExchangeQuote(request.SellToken,
					request.BuyToken, request.SellAmount, request.Chain, fmt.Sprint(source.ChainId), srcTokenDecimals, dstTokenDecimals, request.Slippage, request.TakerAddress, s.config.Swap.DodoEndpoint, source.DodoSwapConfig.ApproveAddress)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		case "lifi":
			if source.ChainName == request.Chain && source.LiFiSwapConfig.IsSupported == true {
				return s.services.LiFi.GetLiFiQuote(request.SellToken, request.BuyToken, request.SellAmount, request.Chain,
					srcTokenDecimals, dstTokenDecimals, request.TakerAddress, request.Slippage)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		case "1inch":
			if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
				return s.services.OneInch.GetExchangeQuote(request, source, srcTokenDecimals)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		case "zeroswap":
			if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
				return s.services.ZeroSwap.GetExchangeQuote(request, source, walletInfo)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		case "cowswap":
			if source.ChainName == request.Chain && source.CowSwapConfig.IsSupported == true {
				return s.services.CowSwap.GetExchangeQuote(request, source, srcTokenDecimals, dstTokenDecimals, walletInfo)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
			}
		default:
			for _, proxy := range s.config.Proxies.ExchangeTypes {
				if request.ExchangeType == proxy {
					return s.exchangeQuoteV1Proxy(request, s.config.PROXIES_ENDPOINT)
				}
			}
		}
	}
	return nil, status.Errorf(codes.Unimplemented,
		fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "Chain not supported")
}

func (s *Swap) GetExchangeMultiQuote(request *pb.ExchangeMultiQuoteRequest) (*pb.ExchangeMultiQuoteResponse, error) {
	chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return &pb.ExchangeMultiQuoteResponse{}, err
	}
	source := s.GetSwapSource(request.Chain)
	walletInfo := s.util.GetWalletInfo(request.Chain)
	decimalMapperSrc, decimalMapperDst := s.GetDecimalMapperQuote(request, chainTokens)
	switch request.ExchangeType {
	case "dzap":
		if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
			return s.services.DZap.GetExchangeQuote(request, walletInfo, decimalMapperSrc)
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "0x":
		if source.ChainName == request.Chain && source.ZeroxSwapConfig.IsSupported == true {
			returnData, err := s.services.ZeroX.GetExchangeQuote(request.MultiChainRequests[0].SellToken,
				request.MultiChainRequests[0].BuyToken, request.MultiChainRequests[0].SellAmount, source.ZeroxSwapConfig.EndPoint, request.Chain, decimalMapperSrc[request.MultiChainRequests[0].SellToken], decimalMapperDst[request.MultiChainRequests[0].BuyToken], request.MultiChainRequests[0].Slippage)
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "dodo":
		if source.ChainName == request.Chain && source.DodoSwapConfig.IsSupported == true {
			returnData, err := s.services.DodoSwap.GetExchangeQuote(request.MultiChainRequests[0].SellToken,
				request.MultiChainRequests[0].BuyToken, request.MultiChainRequests[0].SellAmount, request.Chain,
				fmt.Sprint(source.ChainId), decimalMapperSrc[request.MultiChainRequests[0].SellToken],
				decimalMapperDst[request.MultiChainRequests[0].BuyToken],
				request.MultiChainRequests[0].Slippage, request.TakerAddress, s.config.Swap.DodoEndpoint, source.DodoSwapConfig.ApproveAddress)
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "lifi":
		if source.ChainName == request.Chain && source.LiFiSwapConfig.IsSupported == true {
			returnData, err := s.services.LiFi.GetLiFiQuote(request.MultiChainRequests[0].SellToken,
				request.MultiChainRequests[0].BuyToken, request.MultiChainRequests[0].SellAmount, request.Chain,
				decimalMapperSrc[request.MultiChainRequests[0].SellToken],
				decimalMapperDst[request.MultiChainRequests[0].BuyToken], request.TakerAddress, request.MultiChainRequests[0].Slippage)
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "1inch":
		if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
			returnData, err := s.services.OneInch.GetExchangeQuote(&pb.ExchangeQuoteRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			}, source, decimalMapperSrc[request.MultiChainRequests[0].SellToken])
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "zeroswap":
		if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
			returnData, err := s.services.ZeroSwap.GetExchangeQuote(&pb.ExchangeQuoteRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			}, source, walletInfo)
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "cowswap":
		if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
			returnData, err := s.services.CowSwap.GetExchangeQuote(&pb.ExchangeQuoteRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			}, source, decimalMapperSrc[request.MultiChainRequests[0].SellToken],
				decimalMapperDst[request.MultiChainRequests[0].BuyToken], walletInfo)
			if err != nil {
				return nil, err
			}
			// convert to Multiresponse
			return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	default:
		for _, proxy := range s.config.Proxies.ExchangeTypes {
			if request.ExchangeType == proxy {
				returnData, err := s.exchangeQuoteV1Proxy(&pb.ExchangeQuoteRequest{
					Chain:        request.Chain,
					TakerAddress: request.TakerAddress,
					SellToken:    request.MultiChainRequests[0].SellToken,
					BuyToken:     request.MultiChainRequests[0].BuyToken,
					SellAmount:   request.MultiChainRequests[0].SellAmount,
					Slippage:     request.MultiChainRequests[0].Slippage,
					ExchangeType: request.ExchangeType,
				}, s.config.PROXIES_ENDPOINT)
				if err != nil {
					return nil, err
				}
				// convert to Multiresponse
				return s.ConvertSingleQuoteToMultiQuoteResponse(returnData, request.Chain), err
			}
		}
	}
	return nil, status.Errorf(codes.Unimplemented,
		fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "Unsupported source")
}

func (s *Swap) ConvertSingleQuoteToMultiQuoteResponse(returnData *pb.ExchangeQuoteResponse, chainName string) *pb.ExchangeMultiQuoteResponse {
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
		SellToken:            returnData.SellToken,
		BuyToken:             returnData.BuyToken,
	})
	return &multiQuoteResponse
}

func (s *Swap) GetExchangeSwap(request *pb.ExchangeSwapRequest) (*pb.ExchangeSwapResponse, error) {
	chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return nil, err
	}
	walletInfo := s.util.GetWalletInfo(request.Chain)
	source := s.GetSwapSource(request.Chain)

	_, srcTokenDecimals := CheckExistsWithValueData(strings.ToLower(request.SellToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")
	_, dstTokenDecimals := CheckExistsWithValueData(strings.ToLower(request.BuyToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")

	if srcTokenDecimals == "" {
		decimals, _ := s.getEthTokenDecimals(request.SellToken, request.Chain)

		srcTokenDecimals = strconv.Itoa(decimals)
	}
	if dstTokenDecimals == "" {
		decimals, _ := s.getEthTokenDecimals(request.BuyToken, request.Chain)
		//if err != nil {
		//	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		//	return nil, status.Errorf(codes.Internal, err.Error(), "Chain not supported")
		//}
		dstTokenDecimals = strconv.Itoa(decimals)
	}
	if srcTokenDecimals != "" && dstTokenDecimals != "" {
		switch request.ExchangeType {
		case "0x":
			if source.ChainName == request.Chain && source.ZeroxSwapConfig.IsSupported == true {
				return s.services.ZeroX.GetExchangeSwap(request, source.ZeroxSwapConfig.EndPoint, srcTokenDecimals, walletInfo)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
			}
		case "dodo":
			if source.ChainName == request.Chain && source.DodoSwapConfig.IsSupported == true {
				//Temporarily use V1 to route dodo swap requests
				if request.Chain == "heco" || request.Chain == "boba" || request.Chain == "aurora" || request.Chain == "arbitrum" || request.Chain == "moonriver" || request.Chain == "ethereum" {
					return s.services.DodoSwap.GetExchangeSwap(request, fmt.Sprint(source.ChainId), srcTokenDecimals, dstTokenDecimals, s.config.Swap.DodoEndpoint, walletInfo)
				} else {
					return s.exchangeSwapV1Proxy(request, s.config.PROXIES_ENDPOINT, srcTokenDecimals)
				}
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain))
			}
		case "lifi":
			if source.ChainName == request.Chain && source.LiFiSwapConfig.IsSupported == true {
				return s.services.LiFi.GetLiFiSwap(request, srcTokenDecimals)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
			}
		case "1inch":
			if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
				return s.services.OneInch.GetExchangeSwap(request, walletInfo, source, srcTokenDecimals)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
			}
		case "zeroswap":
			if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
				return s.services.ZeroSwap.GetExchangeSwap(request, source, srcTokenDecimals)
			} else {
				return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
			}
		default:
			for _, proxy := range s.config.Proxies.ExchangeTypes {
				if request.ExchangeType == proxy {
					return s.exchangeSwapV1Proxy(request, s.config.PROXIES_ENDPOINT, srcTokenDecimals)
				}
			}
		}
	}
	return nil, status.Errorf(codes.Unimplemented,
		fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "chain not supported")
}

func (s *Swap) GetExchangeMultiSwap(request *pb.ExchangeMultiSwapRequest) (*pb.ExchangeMultipleSwapResponse, error) {
	chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return nil, err
	}
	source := s.GetSwapSource(request.Chain)
	walletInfo := s.util.GetWalletInfo(request.Chain)
	decimalMapperSrc, decimalMapperDst := s.GetDecimalMapperSwap(request, chainTokens)
	var multipleSwapResponse pb.ExchangeMultipleSwapResponse
	switch request.ExchangeType {
	case "dzap":
		if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
			res, err := s.services.DZap.GetExchangeSwap(request, walletInfo, decimalMapperSrc)
			if err != nil {
				return nil, err
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, res.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = &pb.ExchangeSwapResponse{
				To:       res.To,
				Data:     res.Data,
				Value:    res.Value,
				GasLimit: res.GasLimit,
				Gas:      res.Gas,
				TxLink:   res.TxLink,
			}
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "Chain not supported")
		}
	case "0x":
		if source.ChainName == request.Chain && source.ZeroxSwapConfig.IsSupported == true {
			returnData, err := s.services.ZeroX.GetExchangeSwap(&pb.ExchangeSwapRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			},
				source.ZeroxSwapConfig.EndPoint,
				decimalMapperSrc[request.MultiChainRequests[0].SellToken], walletInfo)
			if err != nil {
				return nil, err
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, returnData.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = returnData
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
		}
	case "dodo":
		if source.ChainName == request.Chain && source.DodoSwapConfig.IsSupported == true {
			//Temporarily use V1 to route dodo swap requests
			var res *pb.ExchangeSwapResponse
			if request.Chain == "heco" {
				returnData, err := s.services.DodoSwap.GetExchangeSwap(&pb.ExchangeSwapRequest{
					Chain:        request.Chain,
					TakerAddress: request.TakerAddress,
					SellToken:    request.MultiChainRequests[0].SellToken,
					BuyToken:     request.MultiChainRequests[0].BuyToken,
					SellAmount:   request.MultiChainRequests[0].SellAmount,
					Slippage:     request.MultiChainRequests[0].Slippage,
					ExchangeType: request.ExchangeType,
				},
					fmt.Sprint(source.ChainId),
					decimalMapperSrc[request.MultiChainRequests[0].SellToken],
					decimalMapperDst[request.MultiChainRequests[0].BuyToken], s.config.Swap.DodoEndpoint, walletInfo)
				if err != nil {
					return nil, err
				}
				res = returnData
			} else {
				returnData, err := s.exchangeSwapV1Proxy(&pb.ExchangeSwapRequest{
					Chain:        request.Chain,
					TakerAddress: request.TakerAddress,
					SellToken:    request.MultiChainRequests[0].SellToken,
					BuyToken:     request.MultiChainRequests[0].BuyToken,
					SellAmount:   request.MultiChainRequests[0].SellAmount,
					Slippage:     request.MultiChainRequests[0].Slippage,
					ExchangeType: request.ExchangeType,
				}, s.config.PROXIES_ENDPOINT, decimalMapperSrc[request.MultiChainRequests[0].SellToken])
				if err != nil {
					return nil, err
				}
				res = returnData
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, res.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = res
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
		}
	case "lifi":
		if source.ChainName == request.Chain && source.LiFiSwapConfig.IsSupported == true {
			returnData, err := s.services.LiFi.GetLiFiSwap(&pb.ExchangeSwapRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			},
				decimalMapperSrc[request.MultiChainRequests[0].SellToken])
			if err != nil {
				return nil, err
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, returnData.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = returnData
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
		}
	case "1inch":
		if source.ChainName == request.Chain && source.OneInchSwapConfig.IsSupported == true {
			returnData, err := s.services.OneInch.GetExchangeSwap(&pb.ExchangeSwapRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			}, walletInfo, source, decimalMapperSrc[request.MultiChainRequests[0].SellToken])
			if err != nil {
				return nil, err
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, returnData.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = returnData
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
		}
	case "zeroswap":
		if source.ChainName == request.Chain && source.ZeroswapSwapConfig.IsSupported == true {
			returnData, err := s.services.ZeroSwap.GetExchangeSwap(&pb.ExchangeSwapRequest{
				Chain:        request.Chain,
				TakerAddress: request.TakerAddress,
				SellToken:    request.MultiChainRequests[0].SellToken,
				BuyToken:     request.MultiChainRequests[0].BuyToken,
				SellAmount:   request.MultiChainRequests[0].SellAmount,
				Slippage:     request.MultiChainRequests[0].Slippage,
				ExchangeType: request.ExchangeType,
			}, source, decimalMapperSrc[request.MultiChainRequests[0].SellToken])
			if err != nil {
				return nil, err
			}
			dataResponse, err2 := s.GetMultipleSwapData(request, returnData.To)
			if err2 != nil {
				s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
			}
			multipleSwapResponse = dataResponse
			multipleSwapResponse.Transaction = returnData
			return &multipleSwapResponse, nil
		} else {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Exchange type %s not supports for %s", request.ExchangeType, request.Chain), "chain not supported")
		}
	default:
		for _, proxy := range s.config.Proxies.ExchangeTypes {
			if request.ExchangeType == proxy {
				returnData, err := s.exchangeSwapV1Proxy(&pb.ExchangeSwapRequest{
					Chain:        request.Chain,
					TakerAddress: request.TakerAddress,
					SellToken:    request.MultiChainRequests[0].SellToken,
					BuyToken:     request.MultiChainRequests[0].BuyToken,
					SellAmount:   request.MultiChainRequests[0].SellAmount,
					Slippage:     request.MultiChainRequests[0].Slippage,
					ExchangeType: request.ExchangeType,
				}, s.config.PROXIES_ENDPOINT, "")
				if err != nil {
					return nil, err
				}
				dataResponse, err2 := s.GetMultipleSwapData(request, returnData.To)
				if err2 != nil {
					s.logger.Errorf("Error for getting multiple swap data : %v", err2.Error())
				}
				multipleSwapResponse = dataResponse
				multipleSwapResponse.Transaction = returnData
				return &multipleSwapResponse, nil
			}
		}
	}
	return nil, status.Errorf(codes.Unimplemented,
		fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "Unsupported source")
}

func (s *Swap) GetFreeTradeCount(request *pb.FreeTradeCountRequest) (*pb.FreeTradeCountResponse, error) {
	source := s.GetSwapSource(request.Chain)
	swapExchange := request.ExchangeType
	switch swapExchange {
	case "zeroswap":
		return s.services.ZeroSwap.GetFreeTradeCount(request, source)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "chain not supported")
	}
}

func (s *Swap) GetTokenApproval(request *pb.TokenApprovalRequest) (*pb.TokenApprovalResponse, error) {
	source := s.GetSwapSource(request.Chain)
	swapExchange := request.ExchangeType
	switch swapExchange {
	case "zeroswap":
		return s.services.ZeroSwap.GetTokenApproval(request, source)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "chain not supported")
	}
}

func (s *Swap) GetExchangeSignature(request *pb.ExchangeSignatureRequest) (*pb.ExchangeSignatureResponse, error) {
	source := s.GetSwapSource(request.Chain)
	var srcTokenDecimals, dstTokenDecimals string
	if request.ExchangeType == "cowswap" {
		chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
		if err != nil {
			s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			return nil, err
		}
		_, srcTokenDecimals = CheckExistsWithValueData(strings.ToLower(request.SellToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")
		_, dstTokenDecimals = CheckExistsWithValueData(strings.ToLower(request.BuyToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")

		if srcTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(request.SellToken, request.Chain)
			//if err != nil {
			//	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			//	return &pb.ExchangeSignatureResponse{}, status.Errorf(codes.Internal, err.Error())
			//}
			srcTokenDecimals = strconv.Itoa(decimals)
		}
		if dstTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(request.BuyToken, request.Chain)
			//if err != nil {
			//	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
			//	return &pb.ExchangeSignatureResponse{}, status.Errorf(codes.Internal, err.Error())
			//}
			dstTokenDecimals = strconv.Itoa(decimals)
		}
	}
	walletInfo := s.util.GetWalletInfo(request.Chain)
	swapExchange := request.ExchangeType
	switch swapExchange {
	case "zeroswap":
		return s.services.ZeroSwap.GetSignatureData(request, source)
	case "cowswap":
		return s.services.CowSwap.GetSignatureData(request, source, srcTokenDecimals, dstTokenDecimals, walletInfo)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "chain not supported")
	}
}

func (s *Swap) GetExchangeSwapExecute(request *pb.ExchangeSwapExecuteRequest) (*pb.ExchangeSwapExecuteResponse, error) {
	source := s.GetSwapSource(request.Chain)
	chainTokens, err := s.GetExchangeTokens(&pb.ExchangeTokenRequest{Chain: request.Chain, ExchangeType: request.ExchangeType})
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		return nil, err
	}
	walletInfo := s.util.GetWalletInfo(request.Chain)

	var sellToken, buyToken string
	if request.ExchangeType == "zeroswap" {
		sellToken = strings.ToLower(request.ZeroSwapPayload.SellToken)
		buyToken = strings.ToLower(request.ZeroSwapPayload.BuyToken)
	} else if request.ExchangeType == "cowswap" {
		sellToken = strings.ToLower(request.CowSwapPayload.SellToken)
		sellToken = strings.ToLower(request.CowSwapPayload.SellToken)
	}
	_, srcTokenDecimals := CheckExistsWithValueData(strings.ToLower(sellToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")
	_, dstTokenDecimals := CheckExistsWithValueData(strings.ToLower(buyToken), "TokenDecimals", reflect.ValueOf(&chainTokens.ExchangeTokens), reflect.TypeOf(chainTokens.ExchangeTokens), "TokenAddress")

	if srcTokenDecimals == "" {
		decimals, _ := s.getEthTokenDecimals(sellToken, request.Chain)
		//if err != nil {
		//	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		//	return &pb.ExchangeSwapExecuteResponse{}, status.Errorf(codes.Internal, err.Error())
		//}
		srcTokenDecimals = strconv.Itoa(decimals)
	}

	if dstTokenDecimals == "" {
		decimals, _ := s.getEthTokenDecimals(buyToken, request.Chain)
		//if err != nil {
		//	s.logger.Errorf("Error for Exchange Quote request  is : %v", err.Error())
		//	return nil, status.Errorf(codes.Internal, err.Error())
		//}
		dstTokenDecimals = strconv.Itoa(decimals)
	}

	swapExchange := request.ExchangeType
	switch swapExchange {
	case "zeroswap":
		return s.services.ZeroSwap.ExecuteZeroSwap(request, source, srcTokenDecimals, walletInfo)
	case "cowswap":
		return s.services.CowSwap.ExecuteCowSwapOrder(request, source, srcTokenDecimals, dstTokenDecimals, walletInfo)
	default:
		return nil, status.Errorf(codes.Unimplemented,
			fmt.Sprintf("Unsupported operation: Source %v is unsupported", request.ExchangeType), "chain not supported")
	}

}

// exchangeSwapV1Proxy exchange swap routed via V1 endpoint
func (s *Swap) exchangeSwapV1Proxy(request *pb.ExchangeSwapRequest, proxyUrl string, srcTokenDecimals string) (*pb.ExchangeSwapResponse, error) {
	swapUrlV1 := fmt.Sprintf(proxyUrl+"/v1/exchange/swap?takerAddress=%v&sellToken=%v&buyToken=%v&sellAmount=%v&slippage=%v&exchangeType=%v&chain=%v", request.TakerAddress, request.SellToken, request.BuyToken, request.SellAmount, request.Slippage, request.ExchangeType, request.Chain)
	body, err := s.services.Http.GetRequest(swapUrlV1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var swapResponse *rpc.ExchangeSwapResponseV1
	err = json.Unmarshal(body, &swapResponse)
	if err != nil {
		s.logger.Errorf("Error for Exchange Swap Proxy request  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}

	var successError rpc.V1ErrorResponse
	err = json.Unmarshal(body, &successError)
	if err != nil {
		s.logger.Errorf("Error for Exchange Swap Proxy request  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	if successError.Error.Message != "" {
		return nil, status.Errorf(codes.NotFound, successError.Error.Message, successError.Error.Message)
	}

	response := &pb.ExchangeSwapResponse{
		To:             swapResponse.To,
		Data:           swapResponse.Data,
		Value:          swapResponse.Value,
		GasLimit:       swapResponse.GasLimit,
		Gas:            swapResponse.Gas,
		TxLink:         swapResponse.Txlink,
		MultiRouteData: nil,
	}
	if response.Data == "" {
		return nil, status.Errorf(codes.Unavailable, "", "No Route Found")
	}
	return response, nil
}

// exchangeQuoteV1Proxy exchange quote routed via V1 endpoint
func (s *Swap) exchangeQuoteV1Proxy(request *pb.ExchangeQuoteRequest, proxyUrl string) (*pb.ExchangeQuoteResponse, error) {
	swapUrlV1 := fmt.Sprintf(proxyUrl+"/v1/exchange/quote?takerAddress=%v&sellToken=%v&buyToken=%v&sellAmount=%v&slippage=%v&exchangeType=%v&chain=%v", request.TakerAddress, request.SellToken, request.BuyToken, request.SellAmount, request.Slippage, request.ExchangeType, request.Chain)
	body, err := s.services.Http.GetRequest(swapUrlV1)
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote Proxy request  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var swapResponse *rpc.ExchangeQuoteResponseV1
	err = json.Unmarshal(body, &swapResponse)
	if err != nil {
		s.logger.Errorf("Error for Exchange Quote Proxy request  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}

	var successError rpc.V1ErrorResponse
	err = json.Unmarshal(body, &successError)
	if err != nil {
		return nil, status.Errorf(codes.Internal, string(body), "json unmarshalling error")
	}
	if successError.Error.Message != "" {
		return nil, status.Errorf(codes.NotFound, successError.Error.Message, successError.Error.Message)
	}

	response := &pb.ExchangeQuoteResponse{
		ResAmount:            swapResponse.ResAmount,
		PriceImpact:          swapResponse.PriceImpact,
		ResPricePerToToken:   swapResponse.ResPricePerToToken,
		ResPricePerFromToken: swapResponse.ResPricePerFromToken,
		FromTokenPrice:       swapResponse.FromTokenPrice,
		ToTokenPrice:         swapResponse.ToTokenPrice,
		MinimumReceived:      swapResponse.MinimumReceived,
		ApproveAddress:       swapResponse.ApproveAddress,
	}
	return response, nil
}

// exchangeTokensV1Proxy exchange tokens routed via V1 endpoint
func (s *Swap) exchangeTokensV1Proxy(request *pb.ExchangeTokenRequest, proxyUrl string) (*pb.ExchangeTokenResponse, error) {
	swapUrlV1 := fmt.Sprintf(proxyUrl+"/v1/exchange/tokens?&exchangeType=%v&chain=%v", request.ExchangeType, request.Chain)

	body, err := s.services.Http.GetRequest(swapUrlV1)
	if err != nil {
		var errorRes rpc.V1FailureResponse
		jsonErr := json.Unmarshal(body, &errorRes)
		if jsonErr != nil {
			s.logger.Errorf("Error for Exchange Tokens Proxy request  is : %v", err.Error())
			return nil, status.Errorf(codes.Internal, jsonErr.Error(), "json unmarshalling error")
		}
		if errorRes.Message != "" {
			return nil, status.Errorf(codes.NotFound, errorRes.Message, errorRes.Message)
		}
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	var swapResponse rpc.ExchangeTokenResponseV1
	err = json.Unmarshal(body, &swapResponse)
	if err != nil {
		var successError rpc.V1ErrorResponse
		err = json.Unmarshal(body, &successError)
		if err != nil {
			s.logger.Errorf("Error for Exchange Tokens Proxy request  is : %v", err.Error())
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
		}
		if successError.Error.Message != "" {
			return nil, status.Errorf(codes.NotFound, successError.Error.Message, "Issue at Proxy API")
		}
		s.logger.Errorf("Error for Exchange Tokens Proxy request  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshalling error")
	}
	var response pb.ExchangeTokenResponse
	if len(swapResponse) > 0 {
		for _, item := range swapResponse {
			response.ExchangeTokens = append(response.ExchangeTokens, &pb.ExchangeTokenInfo{
				TokenAddress:  item.TokenAddress,
				TokenDecimals: item.TokenDecimals,
				TokenSymbol:   item.TokenSymbol,
				TokenName:     item.TokenName,
				TokenLogoUrl:  item.TokenLogoURL,
				LogoUrl:       item.LogoURL,
			})
		}
	}

	return &response, nil
}

func CheckExistsWithValueData(itemTobeFoundForExistence string, keyToGetValue string, rVal reflect.Value, rType reflect.Type, keys ...string) (bool, string) {
	g := rVal.Elem()
	returnVal := false
	strValue := ""
	if len(keys) > 0 {
		if rType.Kind() == reflect.Slice {
			for i := 0; i < g.Len(); i++ {
				if len(keys) == 2 {
					if strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).FieldByName(keys[1]).String()) == itemTobeFoundForExistence {
						returnVal = true
						strValue = strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).FieldByName(keyToGetValue).String())
						break
					}
				} else {
					if strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).String()) == itemTobeFoundForExistence {
						returnVal = true
						strValue = strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keyToGetValue).String())
						break
					}
				}

			}
		}
	}
	return returnVal, strValue
}

func (s *Swap) GetDecimalMapperQuote(requestObj *pb.ExchangeMultiQuoteRequest, exchangeTokens *pb.ExchangeTokenResponse) (map[string]string, map[string]string) {
	decimalMapperSrc := map[string]string{}
	decimalMapperDst := map[string]string{}
	for _, item := range requestObj.MultiChainRequests {
		_, srcTokenDecimals := CheckExistsWithValueData(strings.ToLower(item.SellToken), "TokenDecimals", reflect.ValueOf(&exchangeTokens.ExchangeTokens), reflect.TypeOf(exchangeTokens.ExchangeTokens), "TokenAddress")
		_, dstTokenDecimals := CheckExistsWithValueData(strings.ToLower(item.BuyToken), "TokenDecimals", reflect.ValueOf(&exchangeTokens.ExchangeTokens), reflect.TypeOf(exchangeTokens.ExchangeTokens), "TokenAddress")
		if srcTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(item.SellToken, requestObj.Chain)
			srcTokenDecimals = strconv.Itoa(decimals)
		}
		if dstTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(item.BuyToken, requestObj.Chain)
			dstTokenDecimals = strconv.Itoa(decimals)
		}
		decimalMapperSrc[item.SellToken] = srcTokenDecimals
		decimalMapperDst[item.BuyToken] = dstTokenDecimals

	}
	return decimalMapperSrc, decimalMapperDst
}

func (s *Swap) GetDecimalMapperSwap(requestObj *pb.ExchangeMultiSwapRequest, exchangeTokens *pb.ExchangeTokenResponse) (map[string]string, map[string]string) {
	decimalMapperSrc := map[string]string{}
	decimalMapperDst := map[string]string{}
	for _, item := range requestObj.MultiChainRequests {
		_, srcTokenDecimals := CheckExistsWithValueData(strings.ToLower(item.SellToken), "TokenDecimals", reflect.ValueOf(&exchangeTokens.ExchangeTokens), reflect.TypeOf(exchangeTokens.ExchangeTokens), "TokenAddress")
		_, dstTokenDecimals := CheckExistsWithValueData(strings.ToLower(item.BuyToken), "TokenDecimals", reflect.ValueOf(&exchangeTokens.ExchangeTokens), reflect.TypeOf(exchangeTokens.ExchangeTokens), "TokenAddress")
		if srcTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(item.SellToken, requestObj.Chain)
			srcTokenDecimals = strconv.Itoa(decimals)
		}
		if dstTokenDecimals == "" {
			decimals, _ := s.getEthTokenDecimals(item.BuyToken, requestObj.Chain)
			dstTokenDecimals = strconv.Itoa(decimals)
		}
		decimalMapperSrc[item.SellToken] = srcTokenDecimals
		decimalMapperDst[item.BuyToken] = dstTokenDecimals

	}
	return decimalMapperSrc, decimalMapperDst
}

func (s *Swap) GetSwapSource(chain string) config.ChainData {
	for _, c := range s.config.Swap.ChainData {
		if c.ChainName == chain {
			return c
		}
	}
	return config.ChainData{}
}

func (s *Swap) GetMultipleSwapData(request *pb.ExchangeMultiSwapRequest, spender string) (pb.ExchangeMultipleSwapResponse, error) {
	var multipleSwapResponse pb.ExchangeMultipleSwapResponse
	nonceRequest := &pb.NonceRequest{
		Address: request.TakerAddress,
		Chain:   request.Chain,
	}
	nonceResponse, err1 := s.evmCore.GetNonce(nonceRequest)
	if err1 != nil {
		s.logger.Errorf("Error for getting nonce response : %v", err1.Error())
	}
	multipleSwapResponse.QuoteValue = nonceResponse.QuoteValue
	multipleSwapResponse.OpL1Fee = nonceResponse.OpL1Fee
	multipleSwapResponse.GasPrice = nonceResponse.GasPrice
	var contractArray []string
	for _, multiChainRequest := range request.MultiChainRequests {
		contractArray = append(contractArray, multiChainRequest.SellToken)
	}
	allowanceContracts := strings.Join(contractArray, ",")
	allowanceRequest := &pb.AllowanceRequest{
		Chain:    request.Chain,
		Owner:    request.TakerAddress,
		Spender:  spender,
		Contract: allowanceContracts,
	}
	allowanceResonse, err2 := s.evmCore.BulkAllowance(allowanceRequest)
	if err2 != nil {
		s.logger.Errorf("Error for getting allowances response : %v", err1.Error())
	}
	var approvalArray []string
	for index, responseItem := range allowanceResonse.Response {
		if responseItem.Allowance == "0" {
			approvalArray = append(approvalArray, contractArray[index])
		}
	}
	approvalContracts := strings.Join(approvalArray, ",")
	if len(approvalContracts) > 0 {
		approvalRequest := &pb.ApprovalRequest{
			Target: spender,
			Chain:  request.Chain,
			Token:  approvalContracts,
		}
		approvalResponse, err3 := s.evmCore.BulkApproval(approvalRequest)
		if err3 != nil {
			s.logger.Errorf("Error for getting approval response : %v", err1.Error())
		}
		multipleSwapResponse.Approval = approvalResponse.Response
	} else {
		multipleSwapResponse.Approval = []*pb.ApprovalResponse{}
	}
	return multipleSwapResponse, nil
}
