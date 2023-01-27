package debridge

const (
	ethEstimatedDuration       = 180
	bscEstimatedDuration       = 60
	maticEstimatedDuration     = 780
	arbitrumEstimatedDuration  = 60
	avalancheEstimatedDuration = 60
	deBridgeABI                = `[{
      "inputs":[],
      "name":"globalFixedNativeFee",
      "outputs":[
         {
            "internalType":"uint256",
            "name":"",
            "type":"uint256"
         }
      ],
      "stateMutability":"view",
      "type":"function"
   }]`
	deBridgeContractAddress = "0x43dE2d77BF8027e25dBD179B491e8d64f38398aA"
	protocolFeeMethod       = "globalFixedNativeFee"
	nativeTokenAddress      = "0x0000000000000000000000000000000000000000"
)

//TODO:Use coingecko API to get ID mapping
var coingeckoCoinsListMap = map[string]string{
	"bsc":       "binancecoin",
	"matic":     "matic-network",
	"polygon":   "matic-network",
	"avalanche": "avalanche-2",
	"arbitrum":  "ethereum",
	"ethereum":  "ethereum",
}

var coingeckoCoinsListMapping = map[string]string{
	"arbitrum":  "arbitrum-one",
	"bsc":       "binance-smart-chain",
	"matic":     "polygon-pos",
	"polygon":   "polygon-pos",
	"avalanche": "avalanche",
	"ethereum":  "ethereum",
}
