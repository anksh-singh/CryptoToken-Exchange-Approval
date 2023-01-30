package config

import (
	"github.com/spf13/viper"
	"log"
	
)

type Config struct {
	Web                    WebConfig
	NonEVMConfig           NonEVMConfig `yaml:"nonEVMConfig"`
	Logger                 Log
	EVM                    EVM
	Bridge                 BridgeConfig
	Cosmos                 CosmosConfig
	Datadog                Datadog
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

type WebConfig struct {
	Host    string
	Port    string
	LogFile string
	Datadog Datadog `yaml:"datadog"`
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
type ChainInfo struct {
	ID             string `yaml:"id"`
	ChainID        int    `yaml:"chainId"`
	Name           string `yaml:"name"`
	NativeTokenID  string `yaml:"nativeTokenID"`
	LogoUrl        string `yaml:"logoUrl"`
	WrappedTokenID string `yaml:"wrappedTokenID"`
}

type ChainData struct {
	ChainName          string             `yaml:"chainName"`
	ChainId            int                `yaml:"chainId"`
	DZapSwapConfig     DZapSwapConfig     `yaml:"DZapSwapConfig"`
}

type DZapSwapConfig struct {
	ApproveAddress string `yaml:"approveAddress"`
	IsSupported    bool   `yaml:"isSupported"`
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

func LoadConfig(filename, path string) *Config {
	var configuration Config
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
	// configuration = LoadCosmosBackUpEndPoints(configuration)
	return &configuration
}
