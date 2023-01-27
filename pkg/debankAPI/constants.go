package debankAPI

//var supportedChainIds = "ethereum,bsc,gnosis,polygon,fantom,heco,avalanche,optimism,arbitrum,celo,moonriver,cronos,boba,metis,bttc,aurora,moonbeam,fuse,harmony,klaytn,astar,iotex,evmos"
var supportedChainIds = []string{"ethereum", "bsc", "gnosis", "polygon", "fantom", "heco", "avalanche", "optimism", "arbitrum", "celo", "moonriver", "cronos", "boba", "metis", "bttc", "aurora", "moonbeam", "fuse", "harmony", "klaytn", "astar", "iotex", "evmos"}

const (
	DEPOSIT           = "Deposit"
	FARMING           = "Farming"
	INVESTMENT        = "Investment"
	LENDING           = "Lending"
	LIQUIDITY_POOL    = "Liquidity Pool"
	OPTIONS_SELLER    = "Options Seller"
	REWARDS           = "Rewards"
	VESTING           = "Vesting"
	YIELD             = "Yield"
	STAKED            = "Staked"
	LOCKED            = "Locked"
	LEAVERAGEDFARMING = "Leaveraged Farming"
	OPTIONSBUYER      = "Options Buyer"
	INSURANCESELLER   = "Insurance Seller"
	INSURANCEBUYER    = "Insurance Buyer"
	PERPETUALS        = "Perpetuals"
	OTHERS            = "Others"
)

var DebankChainNetworkId = map[string]string{
	"ethereum":  "eth",
	"bsc":       "bsc",
	"gnosis":    "xdai",
	"polygon":   "matic",
	"fantom":    "ftm",
	"heco":      "heco",
	"avalanche": "avax",
	"optimism":  "op",
	"arbitrum":  "arb",
	"celo":      "celo",
	"moonriver": "movr",
	"cronos":    "cro",
	"boba":      "boba",
	"metis":     "metis",
	"bttc":      "btt",
	"aurora":    "aurora",
	"moonbeam":  "mobm",
	"fuse":      "fuse",
	"harmony":   "hmy",
	"klaytn":    "klay",
	"astar":     "astar",
	"iotex":     "iotx",
	"evmos":     "evmos",

	//"OEC": "okt",
	//"SmartBch": "sbch",
	//"Shiden": "sdn",
	//"Palm": "palm",
	//"RSK": "rsk",
	//"Wanchain": "wan",
	//"KCC": "kcc",
	//"Songbird": "sgb",
	//"DFK": "dfk",
	//"Telos": "tlos",
	//"Swimmer": "swm",
	//"Arbitrum Nova": "nova",
	//"Canto": "canto",
	//"Dogechain": "doge",
	//"Kava": "kava",
	//"Step": "step",
	//"Milkomeda": "mada",
	//"Conflux": "cfx",
	//"Bitgert": "brise",
}

var supportedActions = map[string]bool{
	"lido":               true,
	"aave2":              true,
	"matic_aave":         true,
	"avax_aave":          true,
	"granary":            true,
	"avax_granary":       true,
	"ftm_granary ":       true,
	"op_granary":         true,
	"uwulend":            true,
	"ftm_geist":          true,
	"arb_radiantcapital": true,
	"uniswap2":           true,
	"sushiswap":          true,
	"apifi":              true,
	"bsc_pancakeswap":    true,
	"bsc_sushiswap":      true,
	"bsc_biswap":         true,
	"bsc_mdex":           true,
	"bsc_apeswap":        true,
	"bsc_babyswap":       true,
	"matic_quickswap":    true,
	"matic_sushiswap":    true,
	"matic_mmf":          true,
	"matic_apeswap":      true,
	"matic_dfyn":         true,
	"ftm_spookyswap":     true,
	"ftm_sushiswap":      true,
	"arb_sushiswap":      true,
	"avax_sushiswap":     true,
	"avax_traderjoexyz":  true,
	"avax_pangolin":      true,
	"movr_sushiswap":     true,
	"celo_sushiswap":     true,
	"xdai_sushiswap":     true,
	"metis_granary":      true,
	"hmy_granary":        true,
	"hmy_sushiswap":      true,
	"aurora_trisolaris":  true,

	// Missing from Debank
	//"sushiswap" : true,
}

// Need to implement
var actions = map[string]interface{}{
	"lido":       `{"Deposit": true, "Withdraw": true}`,
	"aave2":      `{"Deposit": true, "Withdraw": true, "Borrow": true, "Repay": true}`,
	"matic_aave": `{"Deposit": true, "Withdraw": true, "Borrow": true, "Repay": true}`,
	"avax_aave":  `{"Deposit": true, "Withdraw": true, "Borrow": true, "Repay": true}`,
}

var supportedProtocolActions = map[string]interface{}{
	//"matic_quickswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"arb_radiantcapital": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"matic_aave": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"bsc_biswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"bsc_sushiswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"bsc_pancakeswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"matic_sushiswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	////"movr_sushiswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"movr_sushiswap": `{"Deposit": true, "Farming": false, "Investment": false, "Lending": false, "Liquidity Pool": false, "Options Seller": false, "Rewards": false, "Vesting": false, "Yield": false, "Staked": false, "Locked": false, "Leaveraged Farming": false, "Options Buyer": false, "Insurance Seller": false, "Insurance Buyer": false, "Perpetuals": false, "Others": false}`,
	//"ftm_spookyswap": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"ftm_granary": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
	//"ftm_geist": `{"Deposit": true, "Farming": true, "Investment": true, "Lending": true, "Liquidity Pool": true, "Options Seller": true, "Rewards": true, "Vesting": true, "Yield": true, "Staked": true, "Locked": true, "Leaveraged Farming": true, "Options Buyer": true, "Insurance Seller": true, "Insurance Buyer": true, "Perpetuals": true, "Others": true}`,
}
