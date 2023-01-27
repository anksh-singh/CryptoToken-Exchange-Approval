package debridge

import (
	"bridge-allowance/pkg/grpc/proto/pb"
	"encoding/json"
	"fmt"
	web32 "github.com/chenzhijie/go-web3"
	"github.com/onrik/ethrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func (d *DeBridge) GetRPC(chain int) (string, error) {
	var rpc string
	var err error
	for _, item := range d.env.EVM.Cfg.Wallets {
		if chain == item.ChainID {
			rpc = item.RPC
			return rpc, err
		}
	}
	return rpc, err
}

func (d *DeBridge) GetNativeTokenInfo(chain string) NativeTokenInfo {
	var tokenInfo NativeTokenInfo
	for _, item := range d.env.EVM.Cfg.Wallets {
		chainId, err := strconv.Atoi(chain)
		if err != nil {
			d.logger.Error(err)
		}
		if chainId == item.ChainID {
			tokenInfo.Name = item.NativeTokenInfo.Name
			tokenInfo.Symbol = item.NativeTokenInfo.Symbol
			tokenDecimals, err := strconv.Atoi(item.NativeTokenInfo.Decimals)
			if err != nil {
				d.logger.Error(err)
			}
			tokenInfo.Decimals = int64(tokenDecimals)
			tokenInfo.ChainID = item.NativeTokenInfo.ChainId
			tokenInfo.LogoURI = item.NativeTokenInfo.LogoURI
		}
	}
	return tokenInfo
}

func (d *DeBridge) GetEstimatedGas(request *pb.BridgeTransactionRequest, data string) (string, error) {
	transaction := ethrpc.T{
		From: request.FromAddress,
		To:   request.ToAddress,
		Data: data,
	}
	chain, err := strconv.ParseUint(request.FromChain, 10, 64)
	rpc, err := d.GetRPC(int(chain))
	client := ethrpc.New(rpc)
	gasLimit, err := client.EthEstimateGas(transaction)
	if err != nil {
		d.logger.Error("Error fetching gas estimate")
		gasLimit = 8000000 //TODO:To be refactored
		err = nil
	}
	gasLimit = gasLimit * 31
	return strconv.Itoa(gasLimit), err
}

func (d *DeBridge) FindChainNameAndValidate(fromChainId string) (string, bool) {
	isDeBridgeSupport := false
	for _, item := range d.env.EVM.Cfg.Wallets {
		if strings.ToLower(fromChainId) == strings.ToLower(item.ChainName) {
			if item.DebridgeSupport == true {
				isDeBridgeSupport = true
				return item.ChainName, isDeBridgeSupport
			}

		} else {
			chainId := strconv.FormatInt(int64(item.ChainID), 10)
			if fromChainId == chainId {
				if item.DebridgeSupport == true {
					isDeBridgeSupport = true
					return item.ChainName, isDeBridgeSupport
				}
			}
		}

	}
	return "", isDeBridgeSupport
}

func (d *DeBridge) GetProtocolFee(chain string) (string, error) {
	chainId, err := strconv.ParseUint(chain, 10, 64)
	rpc, err := d.GetRPC(int(chainId))
	web3, err := web32.NewWeb3(rpc)
	deBridgeClient, err := web3.Eth.NewContract(deBridgeABI, deBridgeContractAddress)

	protocolFee, err := deBridgeClient.Call(protocolFeeMethod)
	if err != nil {
		d.logger.Error(err)
		return "", status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}

	platformFee := fmt.Sprintf("%v", protocolFee)
	return platformFee, err
}

func (d *DeBridge) GetProtocolFeeInUSD(request *pb.BridgeQuoteRequest) (string, error) {
	var chainName string
	fromChain := d.FindChainKey(request.FromChain)
	if request.FromToken == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" || request.FromToken != "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		for _, item := range d.env.EVM.Cfg.Wallets {
			chain, err := strconv.ParseInt(fromChain, 10, 64)
			if err != nil {
				d.logger.Error(err)
				return "", status.Errorf(codes.Internal, err.Error(), "type conversion error")
			}
			if int(chain) == item.ChainID {
				chainName = item.ChainName
			}
		}
	}

	coinGeckoId := coingeckoCoinsListMap[chainName]
	currency := "usd"
	url := fmt.Sprintf(d.env.Coingecko.EndPoint+"/simple/price?ids=%s&vs_currencies=usd&x_cg_pro_api_key=%s", coinGeckoId, d.env.Coingecko.APIkey)
	body, err := d.httpRequest.GetRequest(url)
	if err != nil {
		d.logger.Error(err)
		return "", status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var unmarshalTokenExchange map[string]map[string]float64
	err = json.Unmarshal(body, &unmarshalTokenExchange)
	//Set default price to 0.0 to handle a failed 3P call
	var price = 0.0
	if len(unmarshalTokenExchange[coinGeckoId]) != 0 {
		price = unmarshalTokenExchange[coinGeckoId][currency]
	}
	return strconv.FormatFloat(price, 'f', -1, 64), nil
}

func (d *DeBridge) GetFrontierNativeToken(chain string, tokenAddress string) string {
	if tokenAddress != "" {
		for _, list := range d.env.EVM.Cfg.Wallets {
			if chain == list.Bridge.ChainId {
				tokenAddress = strings.ToLower(list.NativeTokenInfo.Address)
			}
		}
	}
	return tokenAddress
}

func (d *DeBridge) GetDeBridgeNativeToken(chain string, tokenAddress string) string {
	for _, item := range d.env.EVM.Cfg.Wallets {
		if chain == item.Bridge.ChainId {
			if strings.ToLower(tokenAddress) == strings.ToLower(item.NativeTokenInfo.Address) {
				return "0x0000000000000000000000000000000000000000"
			} else {
				return tokenAddress
			}
		}
	}
	return tokenAddress
}

func (d *DeBridge) FindChainKey(fromChain string) string {
	if fromChain != "" {
		for _, item := range d.env.EVM.Cfg.Wallets {
			if fromChain == item.Bridge.ChainId {
				return fromChain
			} else if strings.ToLower(fromChain) == strings.ToLower(item.ChainName) || strings.ToLower(fromChain) == strings.ToLower(item.Bridge.ChainKey) {
				return item.Bridge.ChainId
			}
		}
	}
	return ""
}

func (d *DeBridge) GetAmountUSDInDecimal(amount string, decimal int64) string {
	amountInDecimal := d.helper.CalculateRateWithDecimal(amount, decimal)
	value := strconv.FormatFloat(amountInDecimal, 'f', -1, 64)
	return value
}

func (d *DeBridge) GetTokenAmountInUSD(chain string, token string) float64 {
	var price = 0.0
	var chainName string
	var currency = "usd"
	if token == "0x0000000000000000000000000000000000001010" || token == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		for _, item := range d.env.EVM.Cfg.Wallets {
			if token == item.NativeTokenInfo.Address {
				chain, err := strconv.ParseInt(chain, 10, 64)
				if err != nil {
					d.logger.Error(err)
					return 0
				}
				if int(chain) == item.ChainID {
					chainName = item.ChainName
				}
			}
		}

		coinGeckoId := coingeckoCoinsListMap[chainName]

		url := fmt.Sprintf(d.env.Coingecko.EndPoint+"/simple/price?ids=%s&vs_currencies=usd&x_cg_pro_api_key=%s", coinGeckoId, d.env.Coingecko.APIkey)
		body, err := d.httpRequest.GetRequest(url)
		if err != nil {
			d.logger.Error(err)
			return 0
		}
		var unmarshalTokenExchange map[string]map[string]float64
		err = json.Unmarshal(body, &unmarshalTokenExchange)
		if len(unmarshalTokenExchange[coinGeckoId]) != 0 {
			price = unmarshalTokenExchange[coinGeckoId][currency]
		}
	} else {
		for _, item := range d.env.EVM.Cfg.Wallets {
			if token != item.NativeTokenInfo.Address {
				chain, err := strconv.ParseInt(chain, 10, 64)
				if err != nil {
					d.logger.Error(err)
					return 0
				}
				if int(chain) == item.ChainID {
					chainName = item.ChainName
				}
			}
		}

		coinGeckoId := coingeckoCoinsListMapping[chainName]
		currency := "usd"
		url := fmt.Sprintf(d.env.Coingecko.EndPoint+"/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd&x_cg_pro_api_key=%s", coinGeckoId, token, d.env.Coingecko.APIkey)
		body, err := d.httpRequest.GetRequest(url)
		if err != nil {
			d.logger.Error(err)
			return 0
		}
		var unmarshalTokenExchange map[string]map[string]float64
		err = json.Unmarshal(body, &unmarshalTokenExchange)
		if len(unmarshalTokenExchange[token]) != 0 {
			price = unmarshalTokenExchange[token][currency]
		}
	}
	return price
}
