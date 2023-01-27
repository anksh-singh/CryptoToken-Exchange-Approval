package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Web                    WebConfig
	NonEVMConfig           NonEVMConfig `yaml:"nonEVMConfig"`
	Fantom                 FantomConfig
	Arbitrum               ArbitrumConfig
	Coingecko              CoingeckoConfig
	Covalent               CovalentConfig
	BlockNative            BlockNativeConfig
	Unmarshall             UnmarshallConfig
	TrustWallet            TrustWalletConfig
	Celo                   CeloConfig
	OpenAPI                OpenAPIConfig
	Logger                 Log
	EVM                    EVM
	Bridge                 BridgeConfig
	Cosmos                 CosmosConfig
	Datadog                Datadog
	Socket                 SocketConfig
	Debridge               DebridgeConfig
	Alchemy                AlchemyConfig
	Proxies                ProxiesConfig
	Swap                   SwapConfig
	Tenderly               TenderlyConfig
	SignAssist             SignAssistConfig
	BlowFish               BlowFishConfig
	DebankAPI              DebankAPIConfig
	ClientCodes            map[string]string
	WEB_DATADOG_SERVICE    string
	DATADOG_SERVICE        string
	NONEVM_CLUSTER         string
	NONEVM_DATADOG_SERVICE string
	PROXIES_ENDPOINT       string
	BRIDGE_DATADOG_SERVICE string
	COSMOS_DATADOG_SERVICE string
	EVM_DATADOG_SERVICE    string
	NONEVM_GRPC_ENDPOINT   string
	COSMOS_GRPC_ENDPOINT   string
	BRIDGE_GRPC_ENDPOINT   string
	EVM_GRPC_ENDPOINT      string
	LOG_ENCODING_FORMAT    string
}

type BlowFishConfig struct {
	AccessKey string
	EndPoint  string
}

type SignAssistConfig struct {
	AccessKey string
	Endpoint  string
}

type DebankAPIConfig struct {
	EndPoint  string
	AccessKey string
}

type TenderlyConfig struct {
	AccessKey string
	Project   string
	UserName  string
}

type AlchemyConfig struct {
	APIKey string `yaml:"APIKey"`
}

type ProxiesConfig struct {
	ExchangeTypes []string `yaml:"exchangeTypes"`
}

type CeloConfig struct {
	GrpcClientEndPoint   string
	ServerPort           string
	ChainTokensUrl       string
	ZeroXUrl             string
	CeloCoinGeckoChainId string
	OpenAPIChainId       string
	LogFile              string
}

//type NearConfig struct {
//	ServerPort             string
//	SolanaTokenListUrl     string
//	JupiterApi             string
//	JupiterApiTokenListUrl string
//	SolanaLogoUrl          string
//	LogFile                string
//	Datadog                Datadog `yaml:"datadog"`
//}

type NonEVMConfig struct {
	NonEVMWallet []*NonEVMChainInfo `yaml:"nonEVMWallet"`
	ServerPort   string             `yaml:"serverPort"`
	LogFile      string             `yaml:"logFile"`
	Datadog      Datadog            `yaml:"datadog"`
}

type NonEVMChainInfo struct {
	ChainName string       `yaml:"chainName"`
	Solana    SolanaConfig `yaml:"solanaConfig"`
	Near      NearConfig   `yaml:"nearConfig"`
	Aptos     AptosConfig  `yaml:"aptosConfig"`
}

type SolanaConfig struct {
	ChainName              string `yaml:"chainName"`
	SolanaTokenListUrl     string `yaml:"solanaTokenListUrl"`
	JupiterApi             string `yaml:"jupiterApi"`
	JupiterApiTokenListUrl string `yaml:"jupiterApiTokenListUrl"`
	SolanaLogoUrl          string `yaml:"solanaLogoUrl"`
}

type NearConfig struct {
	ChainName        string `yaml:"chainName"`
	NearTokenListUrl string `yaml:"nearTokenListUrl"`
	//JupiterApi             string
	//JupiterApiTokenListUrl string
	Rpc         string `yaml:"rpc"`
	NearLogoUrl string `yaml:"nearLogoUrl"`
}

type AptosConfig struct {
	ChainName    string `yaml:"chainName"`
	Rpc          string `yaml:"rpc"`
	AptosLogoUrl string `yaml:"aptosLogoUrl"`
}

type FantomConfig struct {
	GrpcClientEndPoint      string
	CovalentEndPoint        string
	GraphqlClientEndPoint   string
	ServerPort              string
	FantomCovalentId        string
	FantomCoinGeckoChainId  string
	FantomUnmarshallChainId string
	ChainTokensUrl          string
	ZeroXUrl                string
	OpenAPIChainId          string
	LogFile                 string
}

type ArbitrumConfig struct {
	GrpcClientEndPoint string
	ServerPort         string
	LogFile            string
}

type WebConfig struct {
	Host    string
	Port    string
	LogFile string
	Datadog Datadog `yaml:"datadog"`
}

type CoingeckoConfig struct {
	APIkey             string
	EndPoint           string
	TokenDetailLogoUrl string
}
type TrustWalletConfig struct {
	EndPoint string
}

type UnmarshallConfig struct {
	APIkey     string
	EndPoint   string
	EndPointV2 string
}

type CovalentConfig struct {
	EndPoint            string
	APIKey              string
	CovalentBalancesAPI string
}

type OpenAPIConfig struct {
	EndPoint string
	APIKey   string
}

type LIFIConfig struct {
	EndPoint string
}

type BlockNativeConfig struct {
	EndPoint   string
	AuthHeader string
}

type Log struct {
	LogLevel string
	LogPath  string
}

type BridgeConfig struct {
	Datadog   Datadog   `yaml:"datadog"`
	BridgeCfg BridgeCfg `yaml:"bridgeCfg"`
}

type EVM struct {
	Cfg     Cfg     `yaml:"cfg"`
	Datadog Datadog `yaml:"datadog"`
}

type Nomenclature struct {
	ChainNameFull  string `yaml:"chainNameFull"`
	ChainNameShort string `yaml:"chainNameShort"`
	Unmarshal      string `yaml:"unmarshal"`
	Debank         string `yaml:"debank"`
}

type BridgeNomenclature struct {
	Service  string `yaml:"service"`
	ChainKey string `yaml:"chainKey"`
	ChainId  string `yaml:"chainId"`
}

type Source struct {
	TokenPriceSource string `yaml:"tokenPriceSource"`
	BalanceSource    string `yaml:"balanceSource"`
	HistorySource    string `yaml:"historySource"`
	UserDataSource   string `yaml:"userDataSource"`
	NonceSource      string `yaml:"nonceSource"`
	GasLimitSource   string `yaml:"gasLimitSource"`
	SendTxSource     string `yaml:"sendTxSource"`
}
type Wallets struct {
	ChainName       string             `yaml:"chainName"`
	ChainID         int                `yaml:"chainId"`
	RPC             string             `yaml:"rpc"`
	Decimal         int                `yaml:"decimal"`
	RestAPI         string             `yaml:"restAPI"`
	CurrencySymbol  string             `yaml:"currencySymbol"`
	Nomenclature    Nomenclature       `yaml:"nomenclature"`
	Source          Source             `yaml:"source"`
	SocketSupport   bool               `yaml:"socketSupport"`
	NativeTokenInfo NativeTokenInfo    `yaml:"nativeTokenInfo"`
	GasLimitFactor  GasLimitFactor     `yaml:"gasLimitFactor"`
	Bridge          BridgeNomenclature `yaml:"bridge"`
	DebridgeSupport bool               `yaml:"debridgeSupport"`
	ChainInfo       ChainInfo          `yaml:"chainInfo"`
}

type GasLimitFactor struct {
	Zerox    float64 `yaml:"zeroX"`
	Dodo     float64 `yaml:"dodo"`
	Lifi     float64 `yaml:"lifi"`
	ZeroSwap float64 `yaml:"zeroSwap"`
}

type NativeTokenInfo struct {
	Name     string `yaml:"name"`
	Symbol   string `yaml:"symbol"`
	Address  string `yaml:"address"`
	ChainId  string `yaml:"chainId"`
	Decimals string `yaml:"decimals"`
	LogoURI  string `yaml:"logoURI"`
}

type Cfg struct {
	ServerPort           string    `yaml:"serverPort"`
	LogFile              string    `yaml:"LogFile"`
	ZeroxAffilateAddress string    `yaml:"zeroxAffilateAddress"`
	Wallets              []Wallets `yaml:"wallets"`
}

type BridgeCfg struct {
	ServerPort                      string `yaml:"serverPort"`
	BridgeExchangeChains            string `yaml:"bridgeExchangeChains"`
	BridgeExchangeChainTokens       string `yaml:"bridgeExchangeChainTokens"`
	BridgeExchangeQuote             string `yaml:"bridgeExchangeQuote"`
	BridgeExchangeTransaction       string `yaml:"bridgeExchangeTransaction"`
	BridgeExchangeTransactionStatus string `yaml:"bridgeExchangeTransactionStatus"`
	LogFile                         string `yaml:"LogFile"`
}

type CosmosConfig struct {
	Cfg     CosmosConf `yaml:"cfg"`
	Datadog Datadog    `yaml:"datadog"`
}

type CosmosConf struct {
	GrpcClientEndPoint   string  `yaml:"grpcClientEndPoint"`
	ServerPort           string  `yaml:"serverPort"`
	LogFile              string  `yaml:"LogFile"`
	CosmostationImageUrl string  `yaml:"cosmostationImageUrl"`
	DenomInfo            string  `yaml:"denomInfo"`
	IBCInfo              string  `yaml:"ibcInfo"`
	MinGasFee            float32 `yaml:"minGasFee"`
	BackUpUrls           map[string]string
	Wallets              []CosmosWallets `yaml:"wallets"`
}

type CosmosCurrencyConfig struct {
	Denom       string `yaml:"denom"`
	DisplayName string `yaml:"displayName"`
	Decimal     int    `yaml:"decimal"`
	LogoUrl     string `yaml:"logoUrl"`
	Symbol      string `yaml:"symbol"`
}

type CosmosWallets struct {
	ChainName              string  `yaml:"chainName"`
	ChainID                string  `yaml:"chainId"`
	Decimals               int64   `yaml:"decimals"`
	Denom                  string  `yaml:"denom"`
	Prefix                 string  `yaml:"prefix"`
	CosmostationAssetsInfo string  `yaml:"cosmostationAssetsInfo"`
	BaseFee                float32 `yaml:"baseFee"`
	RPC                    string  `yaml:"rpc"`
	REST                   string  `yaml:"rest"`
	Logo                   string  `yaml:"logo"`
	Source                 Source  `yaml:"source"`
}

type Datadog struct {
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
}

type SocketConfig struct {
	EndPoint string `yaml:"endPoint"`
	APIKey   string `yaml:"APIKey"`
	LogoUrl  string `yaml:"logoUrl"`
}

type DebridgeConfig struct {
	EndPoint            string `yaml:"endPoint"`
	TransactionEndPoint string `yaml:"transactionEndPoint"`
	LogoUrl             string `yaml:"logoUrl"`
}

type ChainInfo struct {
	ID             string `yaml:"id"`
	ChainID        int    `yaml:"chainId"`
	Name           string `yaml:"name"`
	NativeTokenID  string `yaml:"nativeTokenID"`
	LogoUrl        string `yaml:"logoUrl"`
	WrappedTokenID string `yaml:"wrappedTokenID"`
}

type SwapConfig struct {
	ChainData        []ChainData `yaml:"chainData"`
	ZeroxLogoUrl     string      `yaml:"zeroxLogoUrl"`
	DodoLogoUrl      string      `yaml:"dodoLogoUrl"`
	LIFILogoUrl      string      `yaml:"lifiLogoUrl"`
	OneInchLogoUrl   string      `yaml:"oneInchLogoUrl"`
	ZeroSwapLogoUrl  string      `yaml:"zeroSwapLogoUrl"`
	DodoEndpoint     string      `yaml:"dodoEndpoint"`
	LIFIEndpoint     string      `yaml:"lIFIEndpoint"`
	OneInchEndpoint  string      `yaml:"oneInchEndpoint"`
	ZeroSwapEndpoint string      `yaml:"zeroSwapEndpoint"`
	ZeroSwapApiKey   string      `yaml:"zeroSwapApiKey"`
	CowSwapUrl       string      `yaml:"cowSwapUrl"`
	CowSwapLogoUrl   string      `yaml:"cowSwapLogoUrl"`
	DZapUrl          string      `yaml:"DZapUrl"`
	DZapLogoUrl      string      `yaml:"DZapLogoUrl"`
}

type ChainData struct {
	ChainName          string             `yaml:"chainName"`
	ChainId            int                `yaml:"chainId"`
	ZeroxSwapConfig    ZeroxSwapConfig    `yaml:"zeroxSwapConfig"`
	OneInchSwapConfig  OneInchSwapConfig  `yaml:"oneInchSwapConfig"`
	DodoSwapConfig     DodoSwapConfig     `yaml:"dodoSwapConfig"`
	LiFiSwapConfig     LiFiSwapConfig     `yaml:"liFiSwapConfig"`
	ZeroswapSwapConfig ZeroswapSwapConfig `yaml:"zeroswapSwapConfig"`
	CowSwapConfig      CowSwapConfig      `yaml:"cowSwapConfig"`
	DZapSwapConfig     DZapSwapConfig     `yaml:"DZapSwapConfig"`
}

type DZapSwapConfig struct {
	ApproveAddress string `yaml:"approveAddress"`
	IsSupported    bool   `yaml:"isSupported"`
}

type ZeroxSwapConfig struct {
	EndPoint         string `yaml:"endPoint"`
	ExchangeTokenUrl string `yaml:"exchangeTokenUrl"`
	TokenSource      string `yaml:"tokenSource"`
	IsSupported      bool   `yaml:"isSupported"`
}

type OneInchSwapConfig struct {
	ExchangeTokenUrl string `yaml:"exchangeTokenUrl"`
	ApproveAddress   string `yaml:"approveAddress"`
	IsSupported      bool   `yaml:"isSupported"`
}

type DodoSwapConfig struct {
	ExchangeTokenUrl string `yaml:"exchangeTokenUrl"`
	TokenSource      string `yaml:"tokenSource"`
	ApproveAddress   string `yaml:"approveAddress"`
	IsSupported      bool   `yaml:"isSupported"`
}

type LiFiSwapConfig struct {
	IsSupported bool `yaml:"isSupported"`
}

type ZeroswapSwapConfig struct {
	ApproveAddress        string `yaml:"approveAddress"`
	GasLessApproveAddress string `yaml:"gasLessApproveAddress"`
	IsSupported           bool   `yaml:"isSupported"`
}

type CowSwapConfig struct {
	WrappedTokenID string `yaml:"wrappedTokenID"`
	IsSupported    bool   `yaml:"isSupported"`
	ApproveAddress string `yaml:"approveAddress"`
}

type CosmosDirectory struct {
	Name      string `json:"name"`
	Height    int    `json:"height"`
	Available bool   `json:"available"`
	Rest      struct {
		Available bool `json:"available"`
		Height    int  `json:"height"`
		Best      []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"best"`
	} `json:"rest"`
}

func LoadCosmosBackUpEndPoints(config *Config) *Config {
	var mapUrls = make(map[string]string, 0)
	for _, chain := range config.Cosmos.Cfg.Wallets {
		var cosmosDirectory CosmosDirectory
		resp, err := http.Get("https://status.cosmos.directory/" + chain.ChainName)
		if err != nil {
			log.Fatalf(" Error in GET Request  : %v", err.Error())
		}
		if resp.StatusCode == 200 {
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					log.Fatalf(" Error in GET Request  : %v", err.Error())
				}
			}(resp.Body)
			body, ipUtilErr := ioutil.ReadAll(resp.Body)
			if ipUtilErr != nil {
				log.Fatalf(" Error in GET Request  : %v", err.Error())
			}
			err = json.Unmarshal(body, &cosmosDirectory)
			if err != nil {
				log.Fatalf(" Error in GET Request  : %v", err.Error())
			}
			if len(cosmosDirectory.Rest.Best) != 0 {
				mapUrls[chain.ChainName] = cosmosDirectory.Rest.Best[0].Address
			}
		}
	}
	config.Cosmos.Cfg.BackUpUrls = mapUrls
	return config
}

func LoadConfig(filename, path string) *Config {
	var configuration *Config
	var configName string
	configName = "default_config" // single config file
	viper.SetConfigName(configName)
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetConfigFile("default_config.yml")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err.Error())
	}
	err := viper.MergeInConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = viper.UnmarshalExact(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	configuration = LoadCosmosBackUpEndPoints(configuration)
	return configuration
}
