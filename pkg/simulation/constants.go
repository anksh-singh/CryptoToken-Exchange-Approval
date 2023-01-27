package simulation

var simulationSupportedChains = map[string]string{
	"ethereum":  "ethereum",
	"polygon":   "polygon",
	"bsc":       "bsc",
	"avalanche": "avalanche",
	"arbitrum":  "arbitrum",
	"gnosis":    "gnosis",
	"fantom":    "fantom",
	"optimism":  "optimism",
	"solana":    "solana",
	"1":         "ethereum",
	"137":       "polygon",
	"56":        "bsc",
	"43114":     "avalanche",
	"42161":     "arbitrum",
	"100":       "gnosis",
	"250":       "fantom",
	"10":        "optimism",
}

var BlowFishSupportedChains = map[string]string{
	"solana":   "solana",
	"polygon":  "polygon",
	"ethereum": "ethereum",
	"137":      "polygon",
	"1":        "ethereum",
}

var SignAssistSupportedChains = map[string]string{
	"ethereum":  "1",
	"polygon":   "137",
	"bsc":       "56",
	"avalanche": "43114",
	"1":         "1",
	"137":       "137",
	"56":        "56",
	"43114":     "43114",
}

var TenderlyChainNetworkId = map[string]string{
	"ethereum":  "1",
	"polygon":   "137",
	"gnosis":    "100",
	"bsc":       "56",
	"optimism":  "10",
	"arbitrum":  "42161",
	"fantom":    "250",
	"avalanche": "43114",
	"rsk":       "30",
	"kiln":      "1337802",
	"poa":       "99",
	"1":         "1",
	"137":       "137",
	"100":       "100",
	"56":        "56",
	"10":        "10",
	"42161":     "42161",
	"250":       "250",
	"43114":     "43114",
	"30":        "30",
	"1337802":   "1337802",
	"99":        "99",
}

const (
	tokenABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [{"name": "","type": "string"}],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [{ "name": "", "type": "uint8" }],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [{ "name": "", "type": "string" }],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}]`

	decimals        = "decimals"
	name            = "name"
	symbol          = "symbol"
	approveHexValue = "0x0000000000000000000000000000000000000000ffffffffffffffffffffffff"
)

var ChainIds = map[string]string{
	"ethereum":  "1",
	"polygon":   "137",
	"gnosis":    "100",
	"bsc":       "56",
	"optimism":  "10",
	"arbitrum":  "42161",
	"fantom":    "250",
	"avalanche": "43114",
	//"solana":    "solana",
}
