package utils

import (
	conf "bridge-allowance/config"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

var testlogger = SetupLogger("info", "test.log", "json")

type expect struct {
	name string
	path string
	want []string
}

var expectargs = []expect{
	expect{
		name: "config",
		path: "../config/test",
		want: []string{"arbitrum", "ethereum", "fantom", "harmony", "matic", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc"},
	},
	expect{
		name: "configinvalid",
		path: "../config/test",
		want: []string{},
	},
	expect{
		name: "default_dev",
		path: "../config",
		want: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
	},
	expect{
		name: "default_production",
		path: "../config",
		want: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
	},
	expect{
		name: "default_staging",
		path: "../config",
		want: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
	},
}

func TestGetEVMChains(t *testing.T) {

	for _, arg := range expectargs {
		t.Run(arg.name, func(t *testing.T) {
			//want  := []string {"arbitrum", "ethereum", "fantom", "harmony", "matic", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc"};

			configobj := conf.LoadConfig(arg.name, arg.path)
			u := NewUtils(testlogger, configobj)
			got := u.GetEVMChains()
			var configuration conf.Config
			data := getConfig(arg.name)
			err := yaml.Unmarshal([]byte(data), &configuration)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if !reflect.DeepEqual(got, arg.want) && !(len(got) == 0 && len(arg.want) == 0) {
				t.Errorf("GetEVMChains() got = %v, want %v", got, arg.want)
				t.Errorf("GetEVMChains() got = %v, want %v", len(got), len(arg.want))
			}
		})
	}
}

func TestGetEVMChainsFromInvalidConfig(t *testing.T) {
	want := []string{}

	configobj := conf.LoadConfig("configinvalid", "../config/test")
	u := NewUtils(testlogger, configobj)

	got := u.GetEVMChains()
	if !(len(got) == 0 && len(want) == 0) {
		t.Errorf("TestGetEVMChainsFromInvalidConfig() got = %v, want %v", got, want)
	}
}

func TestALLGetWalletSource(t *testing.T) {
	//var chainName = "arbitrum1"
	//configs := []string {"config","configinvalid"};
	//chains  := []string {"arbitrum", "ethereum", "fantom", "harmony", "matic", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc"};

	for _, arg := range expectargs {
		var configuration conf.Config
		data := getConfig(arg.name)
		err := yaml.Unmarshal([]byte(data), &configuration)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		configobj := conf.LoadConfig(arg.name, arg.path)
		util := NewUtils(testlogger, configobj)

		for _, name := range arg.want {
			t.Run(arg.name+"/"+name, func(t *testing.T) {
				want := GetWalletSourceFromTestData(name, configuration)

				got := util.GetWalletSource(name)

				if !reflect.DeepEqual(got, want) {
					t.Errorf("TestALLGetWalletSource() got = %v, want %v", got, want)
				}
			})
		}
	}

}

func TestALLGetWalletInfo(t *testing.T) {
	//configs := []string {"config","configinvalid"};
	//chains  := []string {"arbitrum", "ethereum", "fantom", "harmony", "matic", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc"};
	for _, arg := range expectargs {
		var configuration conf.Config
		data := getConfig(arg.name)
		err := yaml.Unmarshal([]byte(data), &configuration)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		configobj := conf.LoadConfig(arg.name, arg.path)
		util := NewUtils(testlogger, configobj)

		for _, name := range arg.want {
			t.Run(arg.name+"/"+name, func(t *testing.T) {
				want := GetWalletInfoFromTestData(name, configuration)

				got := util.GetWalletInfo(name)

				if !reflect.DeepEqual(got, want) {
					t.Errorf("TestALLGetWalletInfo() got = %v, want %v", got, want)
				}
			})
		}
	}
}

func TestALLIsEVMFromTestData(t *testing.T) {

	type expect struct {
		name   string
		path   string
		chains []string
		want   bool
	}

	var expectargs = []expect{
		expect{
			name:   "config",
			path:   "../config/test",
			want:   true,
			chains: []string{"arbitrum", "ethereum", "fantom", "harmony", "matic", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc"},
		},
		expect{
			name:   "configinvalid",
			path:   "../config/test",
			want:   false,
			chains: []string{},
		},
		expect{
			name:   "default_dev",
			path:   "../config",
			want:   true,
			chains: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
		},
		expect{
			name:   "default_production",
			path:   "../config",
			want:   true,
			chains: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
		},
		expect{
			name:   "default_staging",
			path:   "../config",
			want:   true,
			chains: []string{"arbitrum", "ethereum", "fantom", "harmony", "polygon", "celo", "optimism", "xinfin", "metis", "avalanche", "aurora", "bsc", "boba", "zilliqa", "heco", "xdai", "cronos", "gnosis", "moonriver", "moonbeam", "klaytn", "bttc", "iotex"},
		},
	}

	for _, expectargument := range expectargs {
		configobj := conf.LoadConfig(expectargument.name, expectargument.path)
		util := NewUtils(testlogger, configobj)
		got := util.GetEVMChains()
		for _, name := range expectargument.chains {
			t.Run(expectargument.name+"/"+name, func(t *testing.T) {
				res := IsEVMFromTestData(name, got)
				if res != expectargument.want {
					t.Errorf("TestALLIsEVMFromTestData() not EVM = %v ", name)
				}
			})
		}
	}
}

func GetWalletSourceFromTestData(chain string, configuration conf.Config) conf.Source {
	for _, w := range configuration.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			return w.Source
		}
	}
	return conf.Source{}
}

func GetWalletInfoFromTestData(chain string, configuration conf.Config) conf.Wallets {
	for _, w := range configuration.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			return w
		}
	}
	return conf.Wallets{}
}

func IsEVMFromTestData(str string, chains []string) bool {
	for _, v := range chains {
		if v == str {
			return true
		}
	}
	return false
}

func getConfig(typ string) string {
	switch typ {
	case "config":
		return getConfigTestData()
	case "configinvalid":
		return getInValidConfigTestData()
	case "default_dev":
		return getDefault_devTestData()
	case "default_production":
		return getDefault_productionTestData()
	case "default_staging":
		return getDefault_stagingTestData()
	}
	return getConfigTestData()
}

/*
//ConvertXdcAddressTo0x replace xdc prefix with 0x
func (u *UtilConf) ConvertXdcAddressTo0xFromTestData(address string) string {
	if len(address) > 3 && "xdc" == address[:3] {
		return fmt.Sprint("0x", address[3:])
	}
	return address
}*/

func getDefault_stagingTestData() string {
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-staging"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Staging"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  cluster: "staging"
  solanaTokenListUrl: "https://raw.githubusercontent.com/nonevm-labs/token-list/main/src/tokens/nonevm.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  jupiterApiTokenListUrl: "https://cache.jup.ag/tokens"
  solanaLogoUrl: "https://cdn.jsdelivr.net/gh/trustwallet/assets@master/blockchains/nonevm/info/logo.png"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-staging"
    version: "1.0"

trustWallet:
  endPoint: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains"


# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"
  tokenDetailLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"


unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

lifi:
  endpoint: "https://li.quest/v1"

bridge:
  bridgeCfg:
    bridgeExchangeChains: "liFiBridge"
    bridgeExchangeChainTokens: "liFiBridge"
    bridgeExchangeQuote: "liFiBridge"
    bridgeExchangeTransaction: "liFiBridge"
    grpcClientEndPoint: "unifront-adapter-bridge:80"
    serverPort: "8085"
    LogFile: "bridge.log"
  datadog:
    env: "dev"
    service: "bridge-dev"
    version: "1.0"


# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "unifront-adapter-cosmos:80"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-staging"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: "unifront-adapter-evm:80"
    serverPort: "8083"
    LogFile: "evm.log"
    wallets:
      - chainName: "arbitrum"
        chainId: 42161
        rpc: "https://arb1.arbitrum.io/rpc"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/arbitrum.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        bridge:
          service: "lifi"
          chainKey: "ARB"
          chainId: "42161"
        nomenclature:
          chainNameFull: "arbitrum"
          chainNameShort: "arb"
          unmarshal: "arbitrum"
          debank: "arb"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "ethereum"
        chainId: 1
        rpc: "https://eth-mainnet.alchemyapi.io/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://config.dodoex.io/tokens/mainnet.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "ETH"
          chainId: "1"
        nomenclature:
          chainNameFull: "ethereum"
          chainNameShort: "eth"
          unmarshal: "ethereum"
          debank: "eth"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "fantom"
        chainId: 250
        rpc: "https://rpcapi.fantom.network"
        decimal: 18
        currencySymbol: "FTM"
        exchangeTokenUrl: "https://raw.githubusercontent.com/SpookySwap/spooky-info/master/src/constants/token/spookyswap.json"
        zeroxUrl: "https://fantom.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "FTM"
          chainId: "250"
        nomenclature:
          chainNameFull: "Fantom Opera"
          chainNameShort: "ftm"
          unmarshal: "fantom"
          debank: "ftm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Fantom"
          symbol: "FTM"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "250"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4001/large/Fantom.png?1558015016"
      - chainName: "harmony"
        chainId: 1666600000
        rpc: "https://api.s0.t.hmny.io"
        decimal: 18
        currencySymbol: "ONE"
        bridge:
          service: "lifi"
          chainKey: "ONE"
          chainId: 1666600000
        nomenclature:
          chainNameFull: "harmony one"
          chainNameShort: "one"
          unmarshal: "harmony"
          debank: "hmy"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "polygon"
        chainId: 137
        rpc: "https://polygon-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "MATIC"
        exchangeTokenUrl: "https://unpkg.com/quickswap-default-token-list@1.2.9/build/quickswap-default.tokenlist.json"
        zeroxUrl: "https://polygon.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "POL"
          chainId: "137"
        nomenclature:
          chainNameFull: "polygon"
          chainNameShort: "matic"
          unmarshal: "matic"
          debank: "matic"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Matic Token"
          symbol: "MATIC"
          address: "0x0000000000000000000000000000000000001010"
          chainId: "137"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4713/large/matic-token-icon.png?1624446912"
      - chainName: "celo"
        chainId: 42220
        rpc: "https://forno.celo.org"
        decimal: 18
        currencySymbol: "CELO"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
        zeroxUrl: "https://celo.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor : 3
        bridge:
          service: "lifi"
          chainKey: "CEL"
          chainId: "42220"
        nomenclature:
          chainNameFull: "celo"
          chainNameShort: "CELO"
          unmarshal: "celo"
          debank: "celo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "optimism"
        chainId: 10
        rpc: "https://opt-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ/"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://static.optimism.io/optimism.tokenlist.json"
        zeroxUrl: "https://optimism.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "OPT"
          chainId: "10"
        nomenclature:
          chainNameFull: "optimism"
          chainNameShort: "opt"
          unmarshal: "optimism"
          debank: "op"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "xinfin"
        chainId: 50
        rpc: "https://rpc.xinfin.network/"
        decimal: 18
        currencySymbol: "XDC"
        nomenclature:
          chainNameFull: "xinfin"
          chainNameShort: "xdc"
          unmarshal: "xinfin"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "metis"
        chainId: 1088
        rpc: "https://andromeda.metis.io/?owner=1088"
        decimal: 18
        currencySymbol: "METIS"
        nomenclature:
          chainNameFull: "metis"
          chainNameShort: "metis"
          unmarshal: "metis"
          debank: "metis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "avalanche"
        chainId: 43114
        rpc: "https://api.avax.network/ext/bc/C/rpc"
        decimal: 18
        currencySymbol: "AVAX"
        exchangeTokenUrl: "https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json"
        zeroxUrl: "https://avalanche.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "AVA"
          chainId: "43114"
        nomenclature:
          chainNameFull: "avalanche c-chain"
          chainNameShort: "avax"
          unmarshal: "avalanche"
          debank: "avax"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "AVAX Token"
          symbol: "AVAX"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/12559/large/coin-round-red.png?1604021818"
      - chainName: "aurora"
        chainId: 1313161554
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/trisolaris-labs/tokens/master/lists/1313161554/list.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "aurora"
          chainNameShort: "aurora"
          unmarshal: "aurora"
          debank: "aurora"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "bsc"
        chainId: 56
        rpc: "https://bsc-dataseed.binance.org"
        decimal: 18
        currencySymbol: "BNB"
        exchangeTokenUrl: "https://tokens.pancakeswap.finance/pancakeswap-extended.json"
        zeroxUrl: "https://bsc.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "BSC"
          chainId: "56"
        nomenclature:
          chainNameFull: "binance smart chain"
          chainNameShort: "bsc"
          unmarshal: "bsc"
          debank: "bsc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "BNB Token"
          symbol: "BNB"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "56"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xeEEEeeEeEeEeeeeEeeEEeEeEeeeEEEEeeEeEeeef.png"
      #TODO: Add a way to fetch multiple RPCs
      - chainName: "boba"
        chainId: 288
        rpc: "https://lightning-replica.boba.network/"
        decimal: 18
        currencySymbol: "BOBA"
        exchangeTokenUrl: "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/src/tokens/boba.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "Boba Network"
          chainNameShort: "boba"
          unmarshal: "boba"
          debank: "boba"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          sendTxSource: "" #TODO: Boba Network has a separate RPC for write transaction
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0x4200000000000000000000000000000000000006"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "zilliqa"
        chainId: 288
        rpc: "https://api.zilliqa.com"
        decimal: 18
        currencySymbol: "ZIL"
        nomenclature:
          chainNameFull: "zilliqa"
          chainNameShort: "zilliqa"
          unmarshal: "zilliqa"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "heco"
        chainId: 128
        rpc: "https://http-mainnet-node.huobichain.com"
        decimal: 18
        currencySymbol: "HT"
        nomenclature:
          chainNameFull: "huobi eco chain"
          chainNameShort: "heco"
          unmarshal: "heco"
          debank: "heco"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "xdai"
        chainId: 100
        rpc: "https://rpc.gnosischain.com"
        decimal: 18
        currencySymbol: "XDAI"
        bridge:
          service: "lifi"
          chainKey: "DAI"
          chainId: "100"
        nomenclature:
          chainNameFull: "gnosis chain"
          chainNameShort: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "cronos"
        chainId: 25
        rpc: "https://evm.cronos.org"
        decimal: 18
        currencySymbol: "CRO"
        bridge:
          service: "lifi"
          chainKey: "CRO"
          chainId: "25"
        nomenclature:
          chainNameFull: "cronos"
          chainNameShort: "cro"
          unmarshal: ""
          debank: "cro"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          userDataSource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      - chainName: "gnosis"
        chainId: 25
        rpc: "https://rpc.ankr.com/gnosis"
        decimal: 18
        currencySymbol: "xDai"
        nomenclature:
          chainNameFull: "gnosis"
          chainNameShort: "gnosis"
          unmarshal: "gnosis"
          debank: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "moonriver"
        chainId: 1285
        rpc: "https://moonriver.public.blastapi.io"
        decimal: 18
        currencySymbol: "movr"
        bridge:
          service: "lifi"
          chainKey: "MOR"
          chainId: "1285"
        nomenclature:
          chainNameFull: "moonriver"
          chainNameShort: "movr"
          unmarshal: "moonriver"
          debank: "movr"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/moonriver.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "moonbeam"
        chainId: 1284
        rpc: "https://rpc.api.moonbeam.network"
        decimal: 18
        currencySymbol: "glmr"
        bridge:
          service: "lifi"
          chainKey: "MOO"
          chainId: 1284
        nomenclature:
          chainNameFull: "moonbeam"
          chainNameShort: "moonbeam"
          unmarshal: "moonbeam"
          debank: "mobm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
      - chainName: "klaytn"
        chainId: 8217
        rpc: "https://public-node-api.klaytnapi.com/v1/cypress"
        decimal: 18
        currencySymbol: "klay"
        nomenclature:
          chainNameFull: "klaytn"
          chainNameShort: "klaytn"
          unmarshal: "klaytn"
          debank: "klay"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "bttc"
        chainId: 199
        rpc: "https://rpc.bt.io"
        decimal: 18
        currencySymbol: "btt"
        nomenclature:
          chainNameFull: "bittorrent"
          chainNameShort: "bittorrent"
          unmarshal: "bttc"
          debank: "btt"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "iotex"
        chainId: 4689
        rpc: "https://iotexrpc.com"
        decimal: 18
        currencySymbol: "iotx"
        nomenclature:
          chainNameFull: "iotex"
          chainNameShort: "iotex"
          unmarshal: "iotex"
          debank: "iotx"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""

  datadog:
    env: "dev"
    service: "evm-staging"
    version: "1.0"
`
	return ymltestdata
}

func getDefault_devTestData() string {
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-dev"
    version: "1.0"


datadog:
  env: "dev"
  service: "bridge-allowance-Dev"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  cluster: "dev"
  solanaTokenListUrl: "https://raw.githubusercontent.com/nonevm-labs/token-list/main/src/tokens/nonevm.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  jupiterApiTokenListUrl: "https://cache.jup.ag/tokens"
  solanaLogoUrl: "https://cdn.jsdelivr.net/gh/trustwallet/assets@master/blockchains/nonevm/info/logo.png"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-dev"
    version: "1.0"


coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"
  tokenDetailLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"

unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

trustWallet:
  endPoint: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains"


lifi:
  endpoint: "https://li.quest/v1"

bridge:
  bridgeCfg:
    bridgeExchangeChains: "liFiBridge"
    bridgeExchangeChainTokens: "liFiBridge"
    bridgeExchangeQuote: "liFiBridge"
    bridgeExchangeTransaction: "liFiBridge"
    grpcClientEndPoint: "unifront-adapter-bridge:80"
    serverPort: "8085"
    LogFile: "bridge.log"
  datadog:
    env: "dev"
    service: "bridge-dev"
    version: "1.0"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "unifront-adapter-cosmos:80"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-dev"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: "unifront-adapter-evm:80"
    serverPort: "8083"
    LogFile: "evm.log"
    wallets:
      - chainName: "arbitrum"
        chainId: 42161
        rpc: "https://arb1.arbitrum.io/rpc"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/arbitrum.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        bridge:
          service: "lifi"
          chainKey: "ARB"
          chainId: "42161"
        nomenclature:
          chainNameFull: "arbitrum"
          chainNameShort: "arb"
          unmarshal: "arbitrum"
          debank: "arb"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "ethereum"
        chainId: 1
        rpc: "https://eth-mainnet.alchemyapi.io/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://config.dodoex.io/tokens/mainnet.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "ETH"
          chainId: "1"
        nomenclature:
          chainNameFull: "ethereum"
          chainNameShort: "eth"
          unmarshal: "ethereum"
          debank: "eth"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "fantom"
        chainId: 250
        rpc: "https://rpcapi.fantom.network"
        decimal: 18
        currencySymbol: "FTM"
        exchangeTokenUrl: "https://raw.githubusercontent.com/SpookySwap/spooky-info/master/src/constants/token/spookyswap.json"
        zeroxUrl: "https://fantom.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "FTM"
          chainId: "250"
        nomenclature:
          chainNameFull: "Fantom Opera"
          chainNameShort: "ftm"
          unmarshal: "fantom"
          debank: "ftm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Fantom"
          symbol: "FTM"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "250"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4001/large/Fantom.png?1558015016"
      - chainName: "harmony"
        chainId: 1666600000
        rpc: "https://api.s0.t.hmny.io"
        decimal: 18
        currencySymbol: "ONE"
        bridge:
          service: "lifi"
          chainKey: "ONE"
          chainId: 1666600000
        nomenclature:
          chainNameFull: "harmony one"
          chainNameShort: "one"
          unmarshal: "harmony"
          debank: "hmy"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "polygon"
        chainId: 137
        rpc: "https://polygon-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "MATIC"
        exchangeTokenUrl: "https://unpkg.com/quickswap-default-token-list@1.2.9/build/quickswap-default.tokenlist.json"
        zeroxUrl: "https://polygon.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "POL"
          chainId: "137"
        nomenclature:
          chainNameFull: "polygon"
          chainNameShort: "matic"
          unmarshal: "matic"
          debank: "matic"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Matic Token"
          symbol: "MATIC"
          address: "0x0000000000000000000000000000000000001010"
          chainId: "137"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4713/large/matic-token-icon.png?1624446912"
      - chainName: "celo"
        chainId: 42220
        rpc: "https://forno.celo.org"
        decimal: 18
        currencySymbol: "CELO"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
        zeroxUrl: "https://celo.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor : 3
        bridge:
          service: "lifi"
          chainKey: "CEL"
          chainId: "42220"
        nomenclature:
          chainNameFull: "celo"
          chainNameShort: "CELO"
          unmarshal: "celo"
          debank: "celo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "optimism"
        chainId: 10
        rpc: "https://opt-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ/"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://static.optimism.io/optimism.tokenlist.json"
        zeroxUrl: "https://optimism.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "OPT"
          chainId: "10"
        nomenclature:
          chainNameFull: "optimism"
          chainNameShort: "opt"
          unmarshal: "optimism"
          debank: "op"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "xinfin"
        chainId: 50
        rpc: "https://rpc.xinfin.network/"
        decimal: 18
        currencySymbol: "XDC"
        nomenclature:
          chainNameFull: "xinfin"
          chainNameShort: "xdc"
          unmarshal: "xinfin"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "metis"
        chainId: 1088
        rpc: "https://andromeda.metis.io/?owner=1088"
        decimal: 18
        currencySymbol: "METIS"
        nomenclature:
          chainNameFull: "metis"
          chainNameShort: "metis"
          unmarshal: "metis"
          debank: "metis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "avalanche"
        chainId: 43114
        rpc: "https://api.avax.network/ext/bc/C/rpc"
        decimal: 18
        currencySymbol: "AVAX"
        exchangeTokenUrl: "https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json"
        zeroxUrl: "https://avalanche.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "AVA"
          chainId: "43114"
        nomenclature:
          chainNameFull: "avalanche c-chain"
          chainNameShort: "avax"
          unmarshal: "avalanche"
          debank: "avax"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "AVAX Token"
          symbol: "AVAX"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/12559/large/coin-round-red.png?1604021818"
      - chainName: "aurora"
        chainId: 1313161554
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/trisolaris-labs/tokens/master/lists/1313161554/list.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "aurora"
          chainNameShort: "aurora"
          unmarshal: "aurora"
          debank: "aurora"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "bsc"
        chainId: 56
        rpc: "https://bsc-dataseed.binance.org"
        decimal: 18
        currencySymbol: "BNB"
        exchangeTokenUrl: "https://tokens.pancakeswap.finance/pancakeswap-extended.json"
        zeroxUrl: "https://bsc.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "BSC"
          chainId: "56"
        nomenclature:
          chainNameFull: "binance smart chain"
          chainNameShort: "bsc"
          unmarshal: "bsc"
          debank: "bsc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "BNB Token"
          symbol: "BNB"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "56"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xeEEEeeEeEeEeeeeEeeEEeEeEeeeEEEEeeEeEeeef.png"
      #TODO: Add a way to fetch multiple RPCs
      - chainName: "boba"
        chainId: 288
        rpc: "https://lightning-replica.boba.network/"
        decimal: 18
        currencySymbol: "BOBA"
        exchangeTokenUrl: "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/src/tokens/boba.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "Boba Network"
          chainNameShort: "boba"
          unmarshal: "boba"
          debank: "boba"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          sendTxSource: "" #TODO: Boba Network has a separate RPC for write transaction
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0x4200000000000000000000000000000000000006"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "zilliqa"
        chainId: 288
        rpc: "https://api.zilliqa.com"
        decimal: 18
        currencySymbol: "ZIL"
        nomenclature:
          chainNameFull: "zilliqa"
          chainNameShort: "zilliqa"
          unmarshal: "zilliqa"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "heco"
        chainId: 128
        rpc: "https://http-mainnet-node.huobichain.com"
        decimal: 18
        currencySymbol: "HT"
        nomenclature:
          chainNameFull: "huobi eco chain"
          chainNameShort: "heco"
          unmarshal: "heco"
          debank: "heco"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "xdai"
        chainId: 100
        rpc: "https://rpc.gnosischain.com"
        decimal: 18
        currencySymbol: "XDAI"
        bridge:
          service: "lifi"
          chainKey: "DAI"
          chainId: "100"
        nomenclature:
          chainNameFull: "gnosis chain"
          chainNameShort: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "cronos"
        chainId: 25
        rpc: "https://evm.cronos.org"
        decimal: 18
        currencySymbol: "CRO"
        bridge:
          service: "lifi"
          chainKey: "CRO"
          chainId: "25"
        nomenclature:
          chainNameFull: "cronos"
          chainNameShort: "cro"
          unmarshal: ""
          debank: "cro"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          userDataSource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      - chainName: "gnosis"
        chainId: 25
        rpc: "https://rpc.ankr.com/gnosis"
        decimal: 18
        currencySymbol: "xDai"
        nomenclature:
          chainNameFull: "gnosis"
          chainNameShort: "gnosis"
          unmarshal: "gnosis"
          debank: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "moonriver"
        chainId: 1285
        rpc: "https://moonriver.public.blastapi.io"
        decimal: 18
        currencySymbol: "movr"
        bridge:
          service: "lifi"
          chainKey: "MOR"
          chainId: "1285"
        nomenclature:
          chainNameFull: "moonriver"
          chainNameShort: "movr"
          unmarshal: "moonriver"
          debank: "movr"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/moonriver.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "moonbeam"
        chainId: 1284
        rpc: "https://rpc.api.moonbeam.network"
        decimal: 18
        currencySymbol: "glmr"
        bridge:
          service: "lifi"
          chainKey: "MOO"
          chainId: 1284
        nomenclature:
          chainNameFull: "moonbeam"
          chainNameShort: "moonbeam"
          unmarshal: "moonbeam"
          debank: "mobm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
      - chainName: "klaytn"
        chainId: 8217
        rpc: "https://public-node-api.klaytnapi.com/v1/cypress"
        decimal: 18
        currencySymbol: "klay"
        nomenclature:
          chainNameFull: "klaytn"
          chainNameShort: "klaytn"
          unmarshal: "klaytn"
          debank: "klay"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "bttc"
        chainId: 199
        rpc: "https://rpc.bt.io"
        decimal: 18
        currencySymbol: "btt"
        nomenclature:
          chainNameFull: "bittorrent"
          chainNameShort: "bittorrent"
          unmarshal: "bttc"
          debank: "btt"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "iotex"
        chainId: 4689
        rpc: "https://iotexrpc.com"
        decimal: 18
        currencySymbol: "iotx"
        nomenclature:
          chainNameFull: "iotex"
          chainNameShort: "iotex"
          unmarshal: "iotex"
          debank: "iotx"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""

  datadog:
    env: "dev"
    service: "evm-dev"
    version: "1.0"`
	return ymltestdata
}

func getDefault_productionTestData() string {
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-production"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Production"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  cluster: "production"
  solanaTokenListUrl: "https://raw.githubusercontent.com/nonevm-labs/token-list/main/src/tokens/nonevm.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  jupiterApiTokenListUrl: "https://cache.jup.ag/tokens"
  solanaLogoUrl: "https://cdn.jsdelivr.net/gh/trustwallet/assets@master/blockchains/nonevm/info/logo.png"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-production"
    version: "1.0"


trustWallet:
  endPoint: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains"

lifi:
  endpoint: "https://li.quest/v1"

bridge:
  bridgeCfg:
    bridgeExchangeChains: "liFiBridge"
    bridgeExchangeChainTokens: "liFiBridge"
    bridgeExchangeQuote: "liFiBridge"
    bridgeExchangeTransaction: "liFiBridge"
    grpcClientEndPoint: "unifront-adapter-bridge:80"
    serverPort: "8085"
    LogFile: "bridge.log"
  datadog:
    env: "dev"
    service: "bridge-dev"
    version: "1.0"


# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"
  tokenDetailLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"


unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "unifront-adapter-cosmos:80"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-production"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: "unifront-adapter-evm:80"
    serverPort: "8083"
    LogFile: "evm.log"
    wallets:
      - chainName: "arbitrum"
        chainId: 42161
        rpc: "https://arb1.arbitrum.io/rpc"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/arbitrum.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        bridge:
          service: "lifi"
          chainKey: "ARB"
          chainId: "42161"
        nomenclature:
          chainNameFull: "arbitrum"
          chainNameShort: "arb"
          unmarshal: "arbitrum"
          debank: "arb"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "ethereum"
        chainId: 1
        rpc: "https://eth-mainnet.alchemyapi.io/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://config.dodoex.io/tokens/mainnet.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "ETH"
          chainId: "1"
        nomenclature:
          chainNameFull: "ethereum"
          chainNameShort: "eth"
          unmarshal: "ethereum"
          debank: "eth"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "fantom"
        chainId: 250
        rpc: "https://rpcapi.fantom.network"
        decimal: 18
        currencySymbol: "FTM"
        exchangeTokenUrl: "https://raw.githubusercontent.com/SpookySwap/spooky-info/master/src/constants/token/spookyswap.json"
        zeroxUrl: "https://fantom.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "FTM"
          chainId: "250"
        nomenclature:
          chainNameFull: "Fantom Opera"
          chainNameShort: "ftm"
          unmarshal: "fantom"
          debank: "ftm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Fantom"
          symbol: "FTM"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "250"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4001/large/Fantom.png?1558015016"
      - chainName: "harmony"
        chainId: 1666600000
        rpc: "https://api.s0.t.hmny.io"
        decimal: 18
        currencySymbol: "ONE"
        bridge:
          service: "lifi"
          chainKey: "ONE"
          chainId: 1666600000
        nomenclature:
          chainNameFull: "harmony one"
          chainNameShort: "one"
          unmarshal: "harmony"
          debank: "hmy"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "polygon"
        chainId: 137
        rpc: "https://polygon-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "MATIC"
        exchangeTokenUrl: "https://unpkg.com/quickswap-default-token-list@1.2.9/build/quickswap-default.tokenlist.json"
        zeroxUrl: "https://polygon.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "POL"
          chainId: "137"
        nomenclature:
          chainNameFull: "polygon"
          chainNameShort: "matic"
          unmarshal: "matic"
          debank: "matic"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Matic Token"
          symbol: "MATIC"
          address: "0x0000000000000000000000000000000000001010"
          chainId: "137"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4713/large/matic-token-icon.png?1624446912"
      - chainName: "celo"
        chainId: 42220
        rpc: "https://forno.celo.org"
        decimal: 18
        currencySymbol: "CELO"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
        zeroxUrl: "https://celo.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor : 3
        bridge:
          service: "lifi"
          chainKey: "CEL"
          chainId: "42220"
        nomenclature:
          chainNameFull: "celo"
          chainNameShort: "CELO"
          unmarshal: "celo"
          debank: "celo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "optimism"
        chainId: 10
        rpc: "https://opt-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ/"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://static.optimism.io/optimism.tokenlist.json"
        zeroxUrl: "https://optimism.api.0x.org"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "OPT"
          chainId: "10"
        nomenclature:
          chainNameFull: "optimism"
          chainNameShort: "opt"
          unmarshal: "optimism"
          debank: "op"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "xinfin"
        chainId: 50
        rpc: "https://rpc.xinfin.network/"
        decimal: 18
        currencySymbol: "XDC"
        nomenclature:
          chainNameFull: "xinfin"
          chainNameShort: "xdc"
          unmarshal: "xinfin"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "metis"
        chainId: 1088
        rpc: "https://andromeda.metis.io/?owner=1088"
        decimal: 18
        currencySymbol: "METIS"
        nomenclature:
          chainNameFull: "metis"
          chainNameShort: "metis"
          unmarshal: "metis"
          debank: "metis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "avalanche"
        chainId: 43114
        rpc: "https://api.avax.network/ext/bc/C/rpc"
        decimal: 18
        currencySymbol: "AVAX"
        exchangeTokenUrl: "https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json"
        zeroxUrl: "https://avalanche.api.0x.org"
        dodoExUrl: ""
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 2
        bridge:
          service: "lifi"
          chainKey: "AVA"
          chainId: "43114"
        nomenclature:
          chainNameFull: "avalanche c-chain"
          chainNameShort: "avax"
          unmarshal: "avalanche"
          debank: "avax"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "AVAX Token"
          symbol: "AVAX"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/12559/large/coin-round-red.png?1604021818"
      - chainName: "aurora"
        chainId: 1313161554
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/trisolaris-labs/tokens/master/lists/1313161554/list.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "aurora"
          chainNameShort: "aurora"
          unmarshal: "aurora"
          debank: "aurora"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "bsc"
        chainId: 56
        rpc: "https://bsc-dataseed.binance.org"
        decimal: 18
        currencySymbol: "BNB"
        exchangeTokenUrl: "https://tokens.pancakeswap.finance/pancakeswap-extended.json"
        zeroxUrl: "https://bsc.api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        gasLimitFactor: 3
        bridge:
          service: "lifi"
          chainKey: "BSC"
          chainId: "56"
        nomenclature:
          chainNameFull: "binance smart chain"
          chainNameShort: "bsc"
          unmarshal: "bsc"
          debank: "bsc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "BNB Token"
          symbol: "BNB"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "56"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xeEEEeeEeEeEeeeeEeeEEeEeEeeeEEEEeeEeEeeef.png"
      #TODO: Add a way to fetch multiple RPCs
      - chainName: "boba"
        chainId: 288
        rpc: "https://lightning-replica.boba.network/"
        decimal: 18
        currencySymbol: "BOBA"
        exchangeTokenUrl: "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/src/tokens/boba.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nomenclature:
          chainNameFull: "Boba Network"
          chainNameShort: "boba"
          unmarshal: "boba"
          debank: "boba"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          sendTxSource: "" #TODO: Boba Network has a separate RPC for write transaction
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0x4200000000000000000000000000000000000006"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "zilliqa"
        chainId: 288
        rpc: "https://api.zilliqa.com"
        decimal: 18
        currencySymbol: "ZIL"
        nomenclature:
          chainNameFull: "zilliqa"
          chainNameShort: "zilliqa"
          unmarshal: "zilliqa"
          debank: ""
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "heco"
        chainId: 128
        rpc: "https://http-mainnet-node.huobichain.com"
        decimal: 18
        currencySymbol: "HT"
        nomenclature:
          chainNameFull: "huobi eco chain"
          chainNameShort: "heco"
          unmarshal: "heco"
          debank: "heco"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "xdai"
        chainId: 100
        rpc: "https://rpc.gnosischain.com"
        decimal: 18
        currencySymbol: "XDAI"
        bridge:
          service: "lifi"
          chainKey: "DAI"
          chainId: "100"
        nomenclature:
          chainNameFull: "gnosis chain"
          chainNameShort: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "cronos"
        chainId: 25
        rpc: "https://evm.cronos.org"
        decimal: 18
        currencySymbol: "CRO"
        bridge:
          service: "lifi"
          chainKey: "CRO"
          chainId: "25"
        nomenclature:
          chainNameFull: "cronos"
          chainNameShort: "cro"
          unmarshal: ""
          debank: "cro"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          userDataSource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      - chainName: "gnosis"
        chainId: 25
        rpc: "https://rpc.ankr.com/gnosis"
        decimal: 18
        currencySymbol: "xDai"
        nomenclature:
          chainNameFull: "gnosis"
          chainNameShort: "gnosis"
          unmarshal: "gnosis"
          debank: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "moonriver"
        chainId: 1285
        rpc: "https://moonriver.public.blastapi.io"
        decimal: 18
        currencySymbol: "movr"
        bridge:
          service: "lifi"
          chainKey: "MOR"
          chainId: "1285"
        nomenclature:
          chainNameFull: "moonriver"
          chainNameShort: "movr"
          unmarshal: "moonriver"
          debank: "movr"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/moonriver.json"
        zeroxUrl: ""
        dodoExUrl: "https://route-api.dodoex.io"
        zeroxLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        dodoLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "moonbeam"
        chainId: 1284
        rpc: "https://rpc.api.moonbeam.network"
        decimal: 18
        currencySymbol: "glmr"
        bridge:
          service: "lifi"
          chainKey: "MOO"
          chainId: 1284
        nomenclature:
          chainNameFull: "moonbeam"
          chainNameShort: "moonbeam"
          unmarshal: "moonbeam"
          debank: "mobm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
      - chainName: "klaytn"
        chainId: 8217
        rpc: "https://public-node-api.klaytnapi.com/v1/cypress"
        decimal: 18
        currencySymbol: "klay"
        nomenclature:
          chainNameFull: "klaytn"
          chainNameShort: "klaytn"
          unmarshal: "klaytn"
          debank: "klay"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "bttc"
        chainId: 199
        rpc: "https://rpc.bt.io"
        decimal: 18
        currencySymbol: "btt"
        nomenclature:
          chainNameFull: "bittorrent"
          chainNameShort: "bittorrent"
          unmarshal: "bttc"
          debank: "btt"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "iotex"
        chainId: 4689
        rpc: "https://iotexrpc.com"
        decimal: 18
        currencySymbol: "iotx"
        nomenclature:
          chainNameFull: "iotex"
          chainNameShort: "iotex"
          unmarshal: "iotex"
          debank: "iotx"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          userDataSource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""

  datadog:
    env: "dev"
    service: "evm-production"
    version: "1.0"
`
	return ymltestdata
}

func getConfigTestData() string {
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-dev"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Dev"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  solanaTokenListUrl: "https://raw.githubusercontent.com/nonevm-labs/token-list/main/src/tokens/nonevm.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-dev"
    version: "1.0"

# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"

unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "localhost:8082"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-dev"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: "localhost:8083"
    serverPort: "8083"
    LogFile: "evm.log"
    wallets:
      - chainName: "arbitrum"
        chainId: 42161
        rpc: "https://arb1.arbitrum.io/rpc"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/arbitrum.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "arbitrum"
          chainNameShort: "arb"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "ethereum"
        chainId: 1
        rpc: "https://eth-mainnet.alchemyapi.io/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://config.dodoex.io/tokens/mainnet.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "ethereum"
          chainNameShort: "eth"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
      - chainName: "fantom"
        chainId: 250
        rpc: "https://rpcapi.fantom.network"
        decimal: 18
        currencySymbol: "FTM"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
        zeroxUrl: "https://fantom.api.0x.org"
        nomenclature:
          chainNameFull: "Fantom Opera"
          chainNameShort: "ftm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Fantom"
          symbol: "FTM"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "250"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4001/large/Fantom.png?1558015016"
      - chainName: "harmony"
        chainId: 1666600000
        rpc: "https://api.s0.t.hmny.io"
        decimal: 18
        currencySymbol: "ONE"
        nomenclature:
          chainNameFull: "harmony one"
          chainNameShort: "one"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "matic"
        chainId: 137
        rpc: "https://polygon-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "MATIC"
        exchangeTokenUrl: "https://unpkg.com/quickswap-default-token-list@1.2.9/build/quickswap-default.tokenlist.json"
        zeroxUrl: "https://polygon.api.0x.org"
        nomenclature:
          chainNameFull: "polygon"
          chainNameShort: "matic"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Matic Token"
          symbol: "MATIC"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "137"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4713/large/matic-token-icon.png?1624446912"
      - chainName: "celo"
        chainId: 42220
        rpc: "https://forno.celo.org"
        decimal: 18
        currencySymbol: "CELO"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
        zeroxUrl: "https://celo.api.0x.org"
        nomenclature:
          chainNameFull: "celo"
          chainNameShort: "CELO"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "optimism"
        chainId: 10
        rpc: "https://opt-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ/"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://static.optimism.io/optimism.tokenlist.json"
        zeroxUrl: "https://optimism.api.0x.org"
        nomenclature:
          chainNameFull: "optimism"
          chainNameShort: "opt"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "xinfin"
        chainId: 50
        rpc: "https://rpc.xinfin.network/"
        decimal: 18
        currencySymbol: "XDC"
        nomenclature:
          chainNameFull: "xinfin"
          chainNameShort: "xdc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "metis"
        chainId: 1088
        rpc: "https://andromeda.metis.io/?owner=1088"
        decimal: 18
        currencySymbol: "METIS"
        nomenclature:
          chainNameFull: "metis"
          chainNameShort: "metis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "avalanche"
        chainId: 43114
        rpc: "https://api.avax.network/ext/bc/C/rpc"
        decimal: 18
        currencySymbol: "AVAX"
        exchangeTokenUrl: "https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json"
        nomenclature:
          chainNameFull: "avalanche c-chain"
          chainNameShort: "avax"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "aurora"
        chainId: 1313161554
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/trisolaris-labs/tokens/master/lists/1313161554/list.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "aurora"
          chainNameShort: "aurora"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "bsc"
        chainId: 56
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "BNB"
        exchangeTokenUrl: "https://tokens.pancakeswap.finance/pancakeswap-extended.json"
        zeroxUrl: "https://bsc.api.0x.org"
        nomenclature:
          chainNameFull: "binance smart chain"
          chainNameShort: "bsc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "BNB Token"
          symbol: "BNB"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "56"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xeEEEeeEeEeEeeeeEeeEEeEeEeeeEEEEeeEeEeeef.png"
      #TODO: Add a way to fetch multiple RPCs
      - chainName: "boba"
        chainId: 288
        rpc: "https://lightning-replica.boba.network/"
        decimal: 18
        currencySymbol: "BOBA"
        exchangeTokenUrl: "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/src/tokens/boba.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "Boba Network"
          chainNameShort: "boba"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          sendTxSource: "" #TODO: Boba Network has a separate RPC for write transaction
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "zilliqa"
        chainId: 288
        rpc: "https://api.zilliqa.com"
        decimal: 18
        currencySymbol: "ZIL"
        nomenclature:
          chainNameFull: "zilliqa"
          chainNameShort: "zilliqa"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "heco"
        chainId: 128
        rpc: "https://http-mainnet-node.huobichain.com"
        decimal: 18
        currencySymbol: "HT"
        nomenclature:
          chainNameFull: "huobi eco chain"
          chainNameShort: "heco"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "xdai"
        chainId: 100
        rpc: "https://rpc.gnosischain.com"
        decimal: 18
        currencySymbol: "XDAI"
        nomenclature:
          chainNameFull: "gnosis chain"
          chainNameShort: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "cronos"
        chainId: 25
        rpc: "https://evm.cronos.org"
        decimal: 18
        currencySymbol: "CRO"
        nomenclature:
          chainNameFull: "cronos"
          chainNameShort: "cro"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      - chainName: "gnosis"
        chainId: 25
        rpc: "https://rpc.ankr.com/gnosis"
        decimal: 18
        currencySymbol: "xDai"
        nomenclature:
          chainNameFull: "gnosis"
          chainNameShort: "gnosis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "moonriver"
        chainId: 1285
        rpc: "https://moonriver.public.blastapi.io"
        decimal: 18
        currencySymbol: "movr"
        nomenclature:
          chainNameFull: "moonriver"
          chainNameShort: "movr"
        exchangeTokenUrl: "https://raw.githubusercontent.com/viaprotocol/tokenlists/main/tokenlists/moonriver.json"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
      - chainName: "moonbeam"
        chainId: 1284
        rpc: "https://rpc.api.moonbeam.network"
        decimal: 18
        currencySymbol: "glmr"
        nomenclature:
          chainNameFull: "moonbeam"
          chainNameShort: "moonbeam"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
      - chainName: "klaytn"
        chainId: 8217
        rpc: "https://public-node-api.klaytnapi.com/v1/cypress"
        decimal: 18
        currencySymbol: "klay"
        nomenclature:
          chainNameFull: "klaytn"
          chainNameShort: "klaytn"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "bttc"
        chainId: 199
        rpc: "https://rpc.bt.io"
        decimal: 18
        currencySymbol: "btt"
        nomenclature:
          chainNameFull: "bittorrent"
          chainNameShort: "bittorrent"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""

  datadog:
    env: "dev"
    service: "evm-dev"
    version: "1.0"
  `
	return ymltestdata
}

func getInValidConfigTestData() string {
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-dev"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Dev"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  solanaTokenListUrl: "https://raw.githubusercontent.com/nonevm-labs/token-list/main/src/tokens/nonevm.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-dev"
    version: "1.0"

# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"

unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "localhost:8082"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-dev"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: ""
    serverPort: ""
    LogFile: "evm.log"
    wallets:

  datadog:
    env: "dev"
    service: "evm-dev"
    version: "1.0"
`
	return ymltestdata

}

/*func GetValidConfigTestData() string{
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-dev"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Dev"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  solanaTokenListUrl: "https://raw.githubusercontent.com/solana-labs/token-list/main/src/tokens/solana.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-dev"
    version: "1.0"

# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"

unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "localhost:8082"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-dev"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: "localhost:8083"
    serverPort: "8083"
    LogFile: "evm.log"
    wallets:
      - chainName: "arbitrum"
        chainId: 42161
        rpc: "https://arb1.arbitrum.io/rpc"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/sushiswap/list/master/lists/token-lists/default-token-list/tokens/arbitrum.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "arbitrum"
          chainNameShort: "arb"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "ethereum"
        chainId: 1
        rpc: "https://eth-mainnet.alchemyapi.io/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://config.dodoex.io/tokens/mainnet.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "ethereum"
          chainNameShort: "eth"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
      - chainName: "fantom"
        chainId: 250
        rpc: "https://rpcapi.fantom.network"
        decimal: 18
        currencySymbol: "FTM"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
        zeroxUrl: "https://fantom.api.0x.org"
        nomenclature:
          chainNameFull: "Fantom Opera"
          chainNameShort: "ftm"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Fantom"
          symbol: "FTM"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "250"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4001/large/Fantom.png?1558015016"
      - chainName: "harmony"
        chainId: 1666600000
        rpc: "https://api.s0.t.hmny.io"
        decimal: 18
        currencySymbol: "ONE"
        nomenclature:
          chainNameFull: "harmony one"
          chainNameShort: "one"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "matic"
        chainId: 137
        rpc: "https://polygon-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ"
        decimal: 18
        currencySymbol: "MATIC"
        exchangeTokenUrl: "https://unpkg.com/quickswap-default-token-list@1.2.9/build/quickswap-default.tokenlist.json"
        zeroxUrl: "https://polygon.api.0x.org"
        nomenclature:
          chainNameFull: "polygon"
          chainNameShort: "matic"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "Matic Token"
          symbol: "MATIC"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "137"
          decimals: "18"
          logoURI: "https://assets.coingecko.com/coins/images/4713/large/matic-token-icon.png?1624446912"
      - chainName: "celo"
        chainId: 42220
        rpc: "https://forno.celo.org"
        decimal: 18
        currencySymbol: "CELO"
        exchangeTokenUrl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
        zeroxUrl: "https://celo.api.0x.org"
        nomenclature:
          chainNameFull: "celo"
          chainNameShort: "CELO"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "optimism"
        chainId: 10
        rpc: "https://opt-mainnet.g.alchemy.com/v2/GVkrt_8cLHv1Yi04m7lqZ2dbteVprcjQ/"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://static.optimism.io/optimism.tokenlist.json"
        zeroxUrl: "https://optimism.api.0x.org"
        nomenclature:
          chainNameFull: "optimism"
          chainNameShort: "opt"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
      - chainName: "xinfin"
        chainId: 50
        rpc: "https://rpc.xinfin.network/"
        decimal: 18
        currencySymbol: "XDC"
        nomenclature:
          chainNameFull: "xinfin"
          chainNameShort: "xdc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "metis"
        chainId: 1088
        rpc: "https://andromeda.metis.io/?owner=1088"
        decimal: 18
        currencySymbol: "METIS"
        nomenclature:
          chainNameFull: "metis"
          chainNameShort: "metis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "avalanche"
        chainId: 43114
        rpc: "https://api.avax.network/ext/bc/C/rpc"
        decimal: 18
        currencySymbol: "AVAX"
        exchangeTokenUrl: "https://raw.githubusercontent.com/pangolindex/tokenlists/main/pangolin.tokenlist.json"
        nomenclature:
          chainNameFull: "avalanche c-chain"
          chainNameShort: "avax"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "aurora"
        chainId: 1313161554
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "ETH"
        exchangeTokenUrl: "https://raw.githubusercontent.com/trisolaris-labs/tokens/master/lists/1313161554/list.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "aurora"
          chainNameShort: "aurora"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "bsc"
        chainId: 56
        rpc: "https://mainnet.aurora.dev"
        decimal: 18
        currencySymbol: "BNB"
        exchangeTokenUrl: "https://tokens.pancakeswap.finance/pancakeswap-extended.json"
        zeroxUrl: "https://bsc.api.0x.org"
        nomenclature:
          chainNameFull: "binance smart chain"
          chainNameShort: "bsc"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "ubeswap"
          swapExchangeQuote: "zeroX"
          swapExchangeSwap: "zeroX"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/565/large/0x-protocol.png?1596623034"
        nativeTokenInfo:
          name: "BNB Token"
          symbol: "BNB"
          address: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
          chainId: "56"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xeEEEeeEeEeEeeeeEeeEEeEeEeeeEEEEeeEeEeeef.png"
      #TODO: Add a way to fetch multiple RPCs
      - chainName: "boba"
        chainId: 288
        rpc: "https://lightning-replica.boba.network/"
        decimal: 18
        currencySymbol: "BOBA"
        exchangeTokenUrl: "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/src/tokens/boba.json"
        zeroxUrl: "https://api.0x.org"
        dodoExUrl: "https://route-api.dodoex.io"
        nomenclature:
          chainNameFull: "Boba Network"
          chainNameShort: "boba"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          swapExchangeTokenSource: "cocoswap"
          swapExchangeQuote: "dodoEx"
          swapExchangeSwap: "dodoEx"
          sendTxSource: "" #TODO: Boba Network has a separate RPC for write transaction
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
        nativeTokenInfo:
          name: "Ether"
          symbol: "ETH"
          address: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
          decimals: "18"
          logoURI: "https://assets.unmarshal.io/tokens/0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE.png"
      - chainName: "zilliqa"
        chainId: 288
        rpc: "https://api.zilliqa.com"
        decimal: 18
        currencySymbol: "ZIL"
        nomenclature:
          chainNameFull: "zilliqa"
          chainNameShort: "zilliqa"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "rpc"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      # Limited API support
      - chainName: "heco"
        chainId: 128
        rpc: "https://http-mainnet-node.huobichain.com"
        decimal: 18
        currencySymbol: "HT"
        nomenclature:
          chainNameFull: "huobi eco chain"
          chainNameShort: "heco"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "xdai"
        chainId: 100
        rpc: "https://rpc.gnosischain.com"
        decimal: 18
        currencySymbol: "XDAI"
        nomenclature:
          chainNameFull: "gnosis chain"
          chainNameShort: "xdai"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      # Limited API support
      - chainName: "cronos"
        chainId: 25
        rpc: "https://evm.cronos.org"
        decimal: 18
        currencySymbol: "CRO"
        nomenclature:
          chainNameFull: "cronos"
          chainNameShort: "cro"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: ""
          historySource: ""
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: ""
      - chainName: "gnosis"
        chainId: 25
        rpc: "https://rpc.ankr.com/gnosis"
        decimal: 18
        currencySymbol: "xDai"
        nomenclature:
          chainNameFull: "gnosis"
          chainNameShort: "gnosis"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
      - chainName: "moonriver"
        chainId: 1285
        rpc: "https://moonriver.public.blastapi.io"
        decimal: 18
        currencySymbol: "movr"
        nomenclature:
          chainNameFull: "moonriver"
          chainNameShort: "movr"
        exchangeTokenUrl: "https://raw.githubusercontent.com/viaprotocol/tokenlists/main/tokenlists/moonriver.json"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
          exchangeLogoUrl: "https://assets.coingecko.com/markets/images/588/small/dodoex.png?1601864278"
      - chainName: "moonbeam"
        chainId: 1284
        rpc: "https://rpc.api.moonbeam.network"
        decimal: 18
        currencySymbol: "glmr"
        nomenclature:
          chainNameFull: "moonbeam"
          chainNameShort: "moonbeam"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: "dodoex"
          swapExchangeQuote: "dodoex"
          swapExchangeSwap: "dodoex"
      - chainName: "klaytn"
        chainId: 8217
        rpc: "https://public-node-api.klaytnapi.com/v1/cypress"
        decimal: 18
        currencySymbol: "klay"
        nomenclature:
          chainNameFull: "klaytn"
          chainNameShort: "klaytn"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""
      - chainName: "bttc"
        chainId: 199
        rpc: "https://rpc.bt.io"
        decimal: 18
        currencySymbol: "btt"
        nomenclature:
          chainNameFull: "bittorrent"
          chainNameShort: "bittorrent"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "unmarshal"
          historySource: "unmarshal"
          nonceSource: "debank"
          gasLimitSource: "rpc"
          sendTxSource: "rpc"
          swapExchangeTokenSource: ""
          swapExchangeQuote: ""
          swapExchangeSwap: ""

  datadog:
    env: "dev"
    service: "evm-dev"
    version: "1.0"
  `
return ymltestdata

}

func getInValidConfigTestData() string{
	ymltestdata := `
web:
  port: "8080"
  logfile: "web.log"
  datadog:
    env: "dev"
    service: "web-dev"
    version: "1.0"

datadog:
  env: "dev"
  service: "bridge-allowance-Dev"
  version: "1.0"

nonevm:
  grpcClientEndpoint: "unifront-adapter-nonevm:80"
  serverPort: "8081"
  solanaTokenListUrl: "https://raw.githubusercontent.com/solana-labs/token-list/main/src/tokens/solana.tokenlist.json"
  jupiterApi: "https://quote-api.jup.ag/v1"
  logfile: "nonevm.log"
  datadog:
    env: "dev"
    service: "nonevm-dev"
    version: "1.0"

# fantom is moved to evm adapter
# fantom:
#  grpcClientEndpoint : "localhost:8082"
#  serverPort: "8082"
#  GraphqlClientEndPoint: "https://xapi.fantom.network"
#  fantomcovalentid: "250"
#  fantomunmarshallchainid : "fantom"
#  fantomcoingeckochainid: "fantom"
#  chaintokensurl: "https://raw.githubusercontent.com/Crocoswap/default-token-list/master/src/tokens/fantom.json"
#  zeroxurl: "https://fantom.api.0x.org"
#  openapichainid: "ftm"
#  logfile: "fantom.log"

# celo is moved to evm adapter
# celo:
#   grpcClientEndpoint : "localhost:8084"
#   serverPort: "8084"
#   chaintokensurl: "https://raw.githubusercontent.com/Ubeswap/default-token-list/master/ubeswap.token-list.json"
#   zeroxurl: "https://celo.api.0x.org"
#   celocoingeckochainid: "celo"
#   logfile: "celo.log"

# arbitrum is moved to evm adapter
# arbitrum:
#   grpcClientEndpoint : "localhost:8083"
#   mainNetEndpoint : "https://arb1.arbitrum.io/rpc"
#   testNetEndpoint : "https://rinkeby.arbitrum.io/rpc"
#   serverPort: "8083"
#   chainId: "42161"
#   logfile: "arbitrum.log"

coingecko:
  apikey : "CG-3rZprwbEEjFtakNBS8mghn8H"
  endpoint : "https://api.coingecko.com/api/v3"

unmarshall:
  apikey: "vHIKeRAOTU20ctcwyGQSt8WHA4OoBQ4O2EUae8c9"
  endpoint: "https://api.unmarshal.com/v1"
  endpointv2: "https://api.unmarshal.com/v2"

covalent:
  apikey: "ckey_3137dd17b50348029a5db413978"
  endpoint: "https://api.covalenthq.com/v1"
  covalentBalancesAPI: "%v/%v/address/%v/balances_v2/?key=%v"

blockNative:
  AuthHeader: "7964ddbd-bc3e-4eba-8175-67fe1134e341"
  endpoint: "https://api.blocknative.com"

logger:
  loglevel: "info"
  logpath: ""

openapi:
  endpoint: "https://pro-openapi.debank.com/v1"
  apikey: "883081dbacfc7664464822b9ffa3c58d19f1cf3b"

# Cosmos(Cosmos Env) Adapter configuration
# Cosmos chains with RPC support can be configured for wallet APIs support
#cosmos:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        rest: <rest endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          chainNameCamelCase: <TBD>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
cosmos:
  cfg:
    grpcClientEndPoint: "localhost:8082"
    serverPort: "8082"
    LogFile: "cosmos.log"
    wallets:
      - chainName: "band"
        chainId: "laozi-mainnet"
        rpc: "http://rpc.laozi1.bandchain.org:80"
        rest: "https://laozi1.bandchain.org/api"
        decimal: 6
        currencySymbol: "BAND"
        nomenclature:
          chainNameFull: "bandchain"
          chainNameShort: "band"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "cosmoshub"
        chainId: "cosmoshub-4"
        rpc: "https://rpc-cosmoshub.blockapsis.com"
        rest: "https://lcd-cosmoshub.blockapsis.com"
        decimal: 6
        currencySymbol: "ATOM"
        nomenclature:
          chainNameFull: "Cosmos"
          chainNameShort: "cosmos"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "kava"
        chainId: "kava-9"
        rpc: "https://rpc.kava.io"
        rest: "https://api.data.kava.io"
        decimal: 6
        currencySymbol: "KAVA"
        nomenclature:
          chainNameFull: "Kava"
          chainNameShort: "kava"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "osmosis"
        chainId: "osmosis-1"
        rpc: "https://osmosis.validator.network"
        rest: "https://lcd-osmosis.blockapsis.com"
        decimal: 6
        currencySymbol: "OSMO"
        nomenclature:
          chainNameFull: "Osmosis"
          chainNameShort: "osmo"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
      - chainName: "terra"
        chainId: "columbus-5"
        rpc: "https://terra-rpc.easy2stake.com:443"
        rest: "https://blockdaemon-terra-lcd.api.bdnodes.net:1317"
        decimal: 6
        currencySymbol: "LUNA"
        nomenclature:
          chainNameFull: "Terra"
          chainNameShort: "terra"
        source:
          tokenPriceSource: "coingecko"
          balanceSource: "rest"
          historySource: "rpc"
          nonceSource: ""
          gasLimitSource: ""
          sendTxSource: "rpc"
  datadog:
    env: "dev"
    service: "cosmos-dev"
    version: "1.0"

# EVM Adapter configuration
# EVM chains with RPC support can be configured for wallet APIs support
#evm:
#  config:
#    grpcClientEndPoint: <GRPC client endpoint>
#    serverPort: <GRPC server port>
#    LogFile: <log file name>
#    wallets:
#      - chainName: <chain name in lower case>
#        chainId: <evm chain ID>
#        rpc: <rpc endpoint>
#        decimal: <number of decimals>
#        currencySymbol: <3 letter currency symbol>
#        nomenclature:
#          chainNameFull: <full chain name case-insensitive>
#          chainNameShort: <chain short form>
#          [Not Supported] debank: <debank equivalent name ref:
#                   https://docs.open.debank.com/en/reference/api-pro-reference/chain#get-supported-chain-list>
#          [Not Supported] coingecko: <coingecko equivalent token "id" ref: https://api.coingecko.com/api/v3/coins/list>
#        source:
#          - tokenPriceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            balanceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            historySource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            nonceSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            gasLimitSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
#            sendTxSource: <3P services - Eg: coingecko, unmarshal, rpc etc>
evm:
  cfg:
    grpcClientEndPoint: ""
    serverPort: ""
    LogFile: "evm.log"
    wallets:

  datadog:
    env: "dev"
    service: "evm-dev"
    version: "1.0"
  `
return ymltestdata

}*/
