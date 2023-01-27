// Package coingecko fetches coingecko related url responses
package coingecko

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"bridge-allowance/utils"
	"encoding/json"
	"errors"
	"fmt"
	currency2 "github.com/bojanz/currency"
	web32 "github.com/chenzhijie/go-web3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var currencySymbol = make(map[string]string)
var currencyCode = make(map[string]string)

type ICoinGecko interface {
	GetTokenExchange(currency string, chain string) (*pb.TokenPriceResponse, error)
	GetTokenInfo(*pb.TokenInfoRequest) (*TokenInfoResponse, error)
	GetTokenDetail(in *pb.TokenDetailRequest) (*pb.TokenDetailResponse, error)
	GetTokenExchangeForContract(chain string, contractAddress string, currency string) (*pb.TokenPriceResponse, error)
}

type CoinGecko struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewCoinGecko(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *CoinGecko {
	Init(env, logger, httpRequest)
	return &CoinGecko{env: env, logger: logger, httpRequest: httpRequest}
}

func Init(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) {
	var coingeckoCodes []string
	url := fmt.Sprintf(env.Coingecko.EndPoint+"/simple/supported_vs_currencies?x_cg_pro_api_key=%s", env.Coingecko.APIkey)
	body, err := httpRequest.GetRequest(url)
	if err != nil {
		logger.Error(err)
		return
	}
	err = json.Unmarshal(body, &coingeckoCodes)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, code := range coingeckoCodes {
		gotSymbol, ok := currency2.GetSymbol(strings.ToUpper(code), currency2.NewLocale("en"))
		if ok {
			currencySymbol[gotSymbol] = code
			currencyCode[code] = gotSymbol
		}
	}
	return
}

func (c *CoinGecko) GetTokenExchangeForContract(chain string, contractAddress string, currency string) (*pb.TokenPriceResponse, error) {
	coinGeckoId := assetPlatformsCoingeckoId[chain]
	currency = strings.ToLower(currency)
	url := fmt.Sprintf(c.env.Coingecko.EndPoint+"/simple/token_price/%s?contract_addresses=%s&vs_currencies=%s&x_cg_pro_api_key=%s",
		coinGeckoId, contractAddress, currency, c.env.Coingecko.APIkey)
	c.logger.Info(url)
	//TODO:Handling the incorrect currency symbol
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var unmarshalTokenExchange map[string]map[string]float64
	err = json.Unmarshal(body, &unmarshalTokenExchange)
	if len(unmarshalTokenExchange[contractAddress]) == 0 {
		//Adding Currency as Invalid Basing upon given coingecko id of map is right
		c.logger.Error("Invalid Currency")
		// cant do err.Error() -- as err is nil
		return nil, status.Errorf(codes.InvalidArgument, errors.New("currency not supported by coingecko").Error(), "Invalid Currency")
	}
	return &pb.TokenPriceResponse{
		Price: unmarshalTokenExchange[contractAddress][currency],
	}, nil
}

// unshift moves the given value to index 0 of slice
// familiar to javascript unshift method
func unshift(slice [][]interface{}, value []interface{}) interface{} {
	var typ = reflect.TypeOf(slice)
	if typ.Kind() == reflect.Slice {
		var vv = reflect.ValueOf(slice)
		var tmp = reflect.MakeSlice(typ, vv.Len()+1, vv.Cap()+1)
		tmp.Index(0).Set(reflect.ValueOf(value))
		var dst = tmp.Slice(1, tmp.Len())
		reflect.Copy(dst, vv)
		return tmp.Interface()
	} else {
		return nil
	}
}

func (c *CoinGecko) getTokenInfo(coingeckoId string, platformId string, blockChainSite string,
	in *pb.TokenInfoRequest) (*TokenInfoResponse, error) {
	var res TokenInfoResponse
	var token TokenInfo
	var chart TokenChart
	tokenDefaultAddress := TokenAddressMapForTokenInfo[in.Chain]
	if tokenDefaultAddress != "" {
		tokenInLowerCase := strings.ToLower(in.Token)
		if tokenInLowerCase == tokenDefaultAddress {
			url := c.env.Coingecko.EndPoint + "/coins/" + coingeckoId + "?x_cg_pro_api_key=" + c.env.Coingecko.APIkey
			body, err := c.httpRequest.GetRequest(url)
			if err != nil {
				c.logger.Error(err)
				return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
			}
			if err := json.Unmarshal(body, &token); err != nil {
				c.logger.Error(err)
				return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
			}
		} else {
			//contract address is case sensitive
			url := fmt.Sprintf("%v/coins/%v/contract/%v?&x_cg_pro_api_key=%s", c.env.Coingecko.EndPoint,
				platformId, in.Token, c.env.Coingecko.APIkey)
			body, err := c.httpRequest.GetRequest(url)
			if err != nil {
				c.logger.Error(err)
				return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
			}
			if err := json.Unmarshal(body, &token); err != nil {
				c.logger.Error(err)
				return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
			}
			if blockChainSite != "" {
				token.Links.BlockchainSite[0] = fmt.Sprintf(blockChainSite, in.Token)
			}
		}
	} else {
		url := fmt.Sprintf("%v/coins/%v/contract/%v?&x_cg_pro_api_key=%s", c.env.Coingecko.EndPoint,
			platformId, in.Token, c.env.Coingecko.APIkey)
		body, err := c.httpRequest.GetRequest(url)
		if err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
		}
		if err := json.Unmarshal(body, &token); err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
		}
		if blockChainSite != "" {
			token.Links.BlockchainSite[0] = fmt.Sprintf(blockChainSite, in.Token)
		}
	}
	tokenId := token.ID
	tokenInLowerCase := strings.ToLower(in.Token)
	if tokenInLowerCase == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" ||
		tokenInLowerCase == "0x0000000000000000000000000000000000001010" ||
		tokenInLowerCase == "kteeeeeeeeeeeeeeeeeeeeeeee" ||
		tokenInLowerCase == "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0000" ||
		tokenInLowerCase == "xdceeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" ||
		tokenInLowerCase == "0x4200000000000000000000000000000000000006" ||
		tokenInLowerCase == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {

		chartInfoUrl := fmt.Sprintf("%v/coins/%v/market_chart?vs_currency=usd&days=%v&x_cg_pro_api_key=%s",
			c.env.Coingecko.EndPoint, tokenId, in.Range, c.env.Coingecko.APIkey)
		body, err := c.httpRequest.GetRequest(chartInfoUrl)
		if err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
		}
		if err = json.Unmarshal(body, &chart); err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "json unmarshall error")
		}
		res.Token = token
		res.Chart = chart
		if res.Token.ContractAddress == nil {
			res.Token.ContractAddress = in.Token
		}
		if coingeckoId == "evmos" {
			res.Token.Links.BlockchainSite[0] = "https://evm.evmos.org/"
		}
		//rangeN, _ := strconv.ParseInt(in.Range, 10, 0)
		//currentDate := time.Now()
		//pastDate := currentDate.AddDate(0, 0, -(int(rangeN)))
		//latestAvailableDate := chart.Prices[0][0]
		//beforelatestAvailableDate := chart.Prices[1][0]
		//chartWindow := beforelatestAvailableDate.(int64) - latestAvailableDate.(int64)
		//if latestAvailableDate.(int64)-pastDate.Unix() >= chartWindow {
		//	dates := latestAvailableDate.(int64)
		//	for ok := true; ok; ok = dates > pastDate.Unix() {
		//		chart.Prices  = unshift(chart.Prices,[]interface{}{dates,0}).([][]interface{})
		//		dates = dates - chartWindow
		//	}
		//}
	} else {
		chartInfoUrl := fmt.Sprintf("%v/coins/%s/contract/%s/market_chart?vs_currency=usd&days=%v&x_cg_pro_api_key=%s",
			c.env.Coingecko.EndPoint, platformId, in.Token, in.Range, c.env.Coingecko.APIkey)
		body, err := c.httpRequest.GetRequest(chartInfoUrl)
		if err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}
		if err := json.Unmarshal(body, &chart); err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
		}

		res.Token = token
		res.Chart = chart
		if res.Token.ContractAddress == nil {
			res.Token.ContractAddress = in.Token
		}
		//rangeN, _ := strconv.ParseInt(in.Range, 10, 0)
		//currentDate := time.Now()
		//pastDate := currentDate.AddDate(0, 0, -(int(rangeN)))
		//if rangeN <= 0{
		//	c.logger.Error("Given Range is 0")
		//	return nil, status.Errorf(codes.Unavailable,"Chart information couldn't be fetched for this token Id")
		//}
		//latestAvailableDate := chart.Prices[0][0]
		//beforelatestAvailableDate := chart.Prices[1][0]
		//chartWindow := beforelatestAvailableDate.(float64) - latestAvailableDate.(float64)
		//if latestAvailableDate.(float64)- float64(pastDate.Unix()) >= chartWindow {
		//	dates := latestAvailableDate.(float64)
		//	for ok := true; ok; ok = dates > float64(pastDate.Unix()) {
		//		//chart.Prices = append([][]interface{}{{dates,0}},chart.Prices...)
		//		chart.Prices  = unshift(chart.Prices,[]interface{}{dates,0}).([][]interface{})
		//		dates = dates - chartWindow
		//	}
		//}
	}
	utility := make([]interface{}, 0)
	res.Utility = utility
	return &res, nil
}

func (c *CoinGecko) GetTokenInfo(in *pb.TokenInfoRequest) (*TokenInfoResponse, error) {
	var tokenInfo *TokenInfoResponse
	//Following CoingeckoIds are fetched from /coins/list api
	switch in.Chain {
	case "bsc":
		tokenInfoRes, err := c.getTokenInfo("binancecoin", "binance-smart-chain",
			"https://bscscan.com/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "matic", "polygon":
		tokenInfoRes, err := c.getTokenInfo("matic-network", "polygon-pos",
			"https://polygonscan.com/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "celo":
		tokenInfoRes, err := c.getTokenInfo("", "celo",
			"https://explorer.celo.org/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "optimism":
		tokenInfoRes, err := c.getTokenInfo("", "optimistic-ethereum",
			"https://optimistic.etherscan.io/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "ethereum":
		tokenInfoRes, err := c.getTokenInfo("ethereum", "ethereum",
			"https://etherscan.com/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "avalanche":
		tokenInfoRes, err := c.getTokenInfo("avalanche-2", "avalanche",
			"https://snowtrace.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "solana":
		tokenInfoRes, err := c.getTokenInfo("solana", "solana",
			"https://explorer.solana.com/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "tezos":
		tokenInfoRes, err := c.getTokenInfo("tezos", "tezos",
			"https://tzkt.io/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "tomochain":
		tokenInfoRes, err := c.getTokenInfo("tomochain", "tomochain",
			"https://tomoscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "fantom":
		tokenInfoRes, err := c.getTokenInfo("fantom", "fantom",
			"https://tomoscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "arbitrum":
		tokenInfoRes, err := c.getTokenInfo("", "arbitrum-one",
			"https://tomoscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "algorand":
		tokenInfoRes, err := c.getTokenInfo("algorand", "algorand",
			"https://algoexplorer.io/asset/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "harmony":
		tokenInfoRes, err := c.getTokenInfo("harmony", "harmony-shard-0",
			"https://explorer.harmony.one/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "harmonyone":
		tokenInfoRes, err := c.getTokenInfo("harmony", "harmony-shard-0",
			"https://explorer.harmony.one/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "zilliqa":
		tokenInfoRes, err := c.getTokenInfo("zilliqa", "zilliqa",
			"https://viewblock.io/zilliqa/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "metis":
		tokenInfoRes, err := c.getTokenInfo("metis", "metis-andromeda",
			"", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "aurora":
		tokenInfoRes, err := c.getTokenInfo("aurora", "aurora",
			"", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "elrond":
		tokenInfoRes, err := c.getTokenInfo("elrond-erd-2", "elrond",
			"https://explorer.elrond.com/tokens/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "gnosis":
		tokenInfoRes, err := c.getTokenInfo("xdai", "xdai",
			"https://blockscout.com/xdai/mainnet/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "moonriver":
		tokenInfoRes, err := c.getTokenInfo("moonriver", "moonriver",
			"https://moonriver.moonscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "moonbeam":
		tokenInfoRes, err := c.getTokenInfo("moonbeam", "moonbeam",
			"https://moonscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "klaytn":
		tokenInfoRes, err := c.getTokenInfo("klay-token", "klay-token",
			"https://scope.klaytn.com/account/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "boba":
		tokenInfoRes, err := c.getTokenInfo("ethereum", "boba",
			"https://blockexplorer.boba.network/tokens/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "xinfin":
		tokenInfoRes, err := c.getTokenInfo("xdce-crowd-sale", "xinfin",
			"https://explorer.xinfin.network/tokens/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "heco":
		tokenInfoRes, err := c.getTokenInfo("huobi-token", "huobi-token",
			"https://hecoinfo.com/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "iotex":
		tokenInfoRes, err := c.getTokenInfo("iotex", "iotex", "https://iotexscan.io/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "bttc":
		tokenInfoRes, err := c.getTokenInfo("bittorrent", "bittorrent", "https://bttcscan.com/token/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "evmos":
		tokenInfoRes, err := c.getTokenInfo("evmos", "evmos", "https://evm.evmos.org/address/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "fuse":
		tokenInfoRes, err := c.getTokenInfo("fuse", "fuse", "https://explorer.fuse.io/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "cronos":
		tokenInfoRes, err := c.getTokenInfo("cronos", "cronos", "https://cronoscan.com/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	case "astar":
		tokenInfoRes, err := c.getTokenInfo("astar", "astar", "https://astar.subscan.io/%s", in)
		if err != nil {
			return nil, err
		}
		tokenInfo = tokenInfoRes
	default:
		return nil, status.Errorf(codes.Unavailable, "Chain not found", "Chain not found")
	}
	return tokenInfo, nil
}

func (c *CoinGecko) GetTokenDetail(in *pb.TokenDetailRequest) (*pb.TokenDetailResponse, error) {
	coinGeckoId := assetPlatformsCoingeckoId[in.Chain]
	url := fmt.Sprintf(c.env.Coingecko.EndPoint+"/coins/%s/contract/%s?&x_cg_pro_api_key=%s", coinGeckoId, strings.ToLower(in.ContractAddress), c.env.Coingecko.APIkey)
	var tokenDetail TokenDetail
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		tokenName, tokenSymbol, tokenDecimal, err := c.TokenDetails(in.ContractAddress, in.Chain)
		if err != nil {
			c.logger.Error(err)
			return nil, status.Errorf(codes.Internal, "Could not find the token with the given Token contract")
		}
		return &pb.TokenDetailResponse{
			TokenName:     tokenName,
			TokenSymbol:   strings.ToUpper(tokenSymbol),
			TokenDecimals: tokenDecimal,
			TokenAddress:  in.ContractAddress,
			TokenPrice:    "0",
		}, nil
	}
	err = json.Unmarshal(body, &tokenDetail)
	if err != nil {
		c.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Json Unmarshalling")
	}
	_, _, decimal, err := c.TokenDetails(in.ContractAddress, in.Chain)
	if err != nil {
		c.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	var tokenTradeVolume string
	tradeVolumeStr := fmt.Sprintf("%v", tokenDetail.MarketData.TotalVolume.Usd)
	if strings.Contains(tradeVolumeStr, "e") || strings.Contains(tradeVolumeStr, "E") {
		tokenTradeVolume = fmt.Sprintf("$%d", int(tokenDetail.MarketData.TotalVolume.Usd))
	} else {
		tokenTradeVolume = fmt.Sprintf("$%v", tokenDetail.MarketData.TotalVolume.Usd)
	}
	var tokenWebsite string
	if len(tokenDetail.Links.Homepage) > 0 {
		tokenWebsite = tokenDetail.Links.Homepage[0]
	}
	return &pb.TokenDetailResponse{
		TokenName:         tokenDetail.Name,
		TokenSymbol:       strings.ToUpper(tokenDetail.Symbol),
		TokenDecimals:     decimal,
		TokenAddress:      tokenDetail.ContractAddress,
		TokenLogoUrl:      tokenDetail.Image.Large,
		TokenListedCount:  fmt.Sprint(len(tokenDetail.Tickers)),
		TokenPrice:        fmt.Sprintf("$%v", tokenDetail.MarketData.CurrentPrice.Usd),
		TokenLastActivity: tokenDetail.LastUpdated.Format(time.RFC3339Nano),
		TokenWebsite:      tokenWebsite,
		TokenTradeVolume:  tokenTradeVolume,
		LogoUrl:           c.env.Coingecko.TokenDetailLogoUrl,
	}, nil
}

func (c *CoinGecko) GetTokenExchange(currency string, chain string) (*pb.TokenPriceResponse, error) {
	coinGeckoId := coingeckoCoinsListMap[chain]
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]`).MatchString
	var symbol string
	const DefaultSymbol = "$"
	const DefaultCurrency = "USD"
	if !IsLetter(currency) { //Check if currency is symbol
		if symbolCode, ok := currencySymbol[currency]; ok {
			symbol = currency
			currency = symbolCode
		} else {
			//Defaults to USD currency & $symbol
			symbol = DefaultSymbol
			currency = DefaultCurrency
		}
	} else {
		if curSymbol, ok := currencyCode[strings.ToLower(currency)]; ok {
			symbol = curSymbol
		} else {
			//Defaults to $symbol
			symbol = DefaultSymbol
			currency = DefaultCurrency
		}
	}
	currency = strings.ToLower(currency)
	url := fmt.Sprintf(c.env.Coingecko.EndPoint+"/simple/price?ids="+
		"%s&vs_currencies=%s&x_cg_pro_api_key=%s", coinGeckoId, currency, c.env.Coingecko.APIkey)
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Error(err)
		return nil, status.Errorf(codes.Internal, string(body), "Internal Error")
	}
	var unmarshalTokenExchange map[string]map[string]float64
	err = json.Unmarshal(body, &unmarshalTokenExchange)
	//Set default price to 0.0 to handle a failed 3P call
	var price = 0.0
	if len(unmarshalTokenExchange[coinGeckoId]) != 0 {
		price = unmarshalTokenExchange[coinGeckoId][currency]
	}
	return &pb.TokenPriceResponse{
		Price:          price,
		CurrencyCode:   strings.ToUpper(currency),
		CurrencySymbol: symbol,
	}, nil
}

func (c *CoinGecko) GetTokenExchangeByCoingeckoId(currency string, coinGeckoId string) (float64, error) {
	currency = strings.ToLower(currency)
	url := fmt.Sprintf(c.env.Coingecko.EndPoint+"/simple/price?ids="+
		"%s&vs_currencies=%s&x_cg_pro_api_key=%s", coinGeckoId, currency, c.env.Coingecko.APIkey)
	c.logger.Info(url)
	body, err := c.httpRequest.GetRequest(url)
	if err != nil {
		c.logger.Error(err)
		return 0, status.Errorf(codes.Internal, string(body), "Internal Error")
	}
	var unmarshalTokenExchange map[string]map[string]float64
	err = json.Unmarshal(body, &unmarshalTokenExchange)
	//Set default price to 0.0 to handle a failed 3P call
	var price = 0.0
	if len(unmarshalTokenExchange[coinGeckoId]) != 0 {
		price = unmarshalTokenExchange[coinGeckoId][currency]
	}
	return price, nil
}

// GetTokenExchangeV2 returns the GetTokenExchange's price response in string format
func (c *CoinGecko) GetTokenExchangeV2(currency string, chain string) (*pb.TokenPriceResponseV2, error) {
	tokenPrice, err := c.GetTokenExchange(currency, chain)
	if err != nil {
		return nil, err
	}
	tokenPriceStr := strconv.FormatFloat(tokenPrice.Price, 'f', -1, 64)
	return &pb.TokenPriceResponseV2{
		Price:          tokenPriceStr,
		CurrencyCode:   tokenPrice.CurrencyCode,
		CurrencySymbol: tokenPrice.CurrencySymbol,
	}, nil
}

func (c *CoinGecko) GetRPC(chain string) (string, error) {
	var rpc string
	var err error
	for _, item := range c.env.EVM.Cfg.Wallets {
		if chain == item.ChainName {
			rpc = item.RPC
			return rpc, err
		}
	}
	return rpc, err
}

func (c *CoinGecko) TokenDetails(address string, chain string) (string, string, string, error) {
	rpc, err := c.GetRPC(chain)
	web3, err := web32.NewWeb3(rpc)
	client, err := web3.Eth.NewContract(tokenABI, address)

	tName, err := client.Call(name)
	if err != nil {
		c.logger.Error(err)
		return "", "", "", err
	}

	tSymbol, err := client.Call(symbol)
	if err != nil {
		c.logger.Error(err)
		return "", "", "", err
	}

	tDecimals, err := client.Call(decimals)
	if err != nil {
		c.logger.Error(err)
		return "", "", "", err
	}

	tokenName := fmt.Sprintf("%v", tName)
	tokenSymbol := fmt.Sprintf("%v", tSymbol)
	tokenDecimals := fmt.Sprintf("%v", tDecimals)

	return tokenName, tokenSymbol, tokenDecimals, nil
}

func (c *CoinGecko) GetCosmosTokenMarketInfo(chain string) (QuoteData, error) {
	var quoteRateChange24h string
	var quoteRateCPctchange24h float64
	reqUrl := fmt.Sprintf(c.env.Coingecko.EndPoint+"/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&page=1&"+
		"sparkline=false&x_cg_pro_api_key=%s",
		chain, c.env.Coingecko.APIkey)
	c.logger.Info(reqUrl)
	body, err := c.httpRequest.GetRequest(reqUrl)
	if err != nil {
		c.logger.Error(err)
		return QuoteData{}, status.Errorf(codes.Internal, err.Error())
	}
	var tokenMarketInfo CoinMarketInfo
	err = json.Unmarshal(body, &tokenMarketInfo)
	if err != nil {
		c.logger.Error(err)
		return QuoteData{}, status.Errorf(codes.Internal, err.Error())
	}
	if len(tokenMarketInfo) > 0 {
		quoteRateChange24h = strconv.FormatFloat(tokenMarketInfo[0].PriceChange24H, 'f', -1, 64)
		quoteRateCPctchange24h = tokenMarketInfo[0].PriceChangePercentage24H
	}
	return QuoteData{
		QuoteRateChange24h:    quoteRateChange24h,
		QuoteRatePctChange24h: quoteRateCPctchange24h,
	}, nil

}
