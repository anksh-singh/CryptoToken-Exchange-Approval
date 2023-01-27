package utils

import (
	"bridge-allowance/config"
)

const COSMOS = "cosmos_network"

var CosmosEVMCompatibleChains = map[string]string{
	"cronos-cosmos": "cronos",
	"evmos-cosmos":  "evmos",
	"terra-classic": "terra",
	"terra":         "Unsupported format",
	"cosmos":        "Unsupported format",
}

// GetCosmosChains Load configured evm chains from config
// Returns an empty string array if no configurations are found
func (u *UtilConf) GetCosmosChains() []string {
	var cosmosChainList []string
	if len(u.conf.Cosmos.Cfg.Wallets) < 1 {
		return cosmosChainList
	}
	for _, w := range u.conf.Cosmos.Cfg.Wallets {
		cosmosChainList = append(cosmosChainList, w.ChainName)
	}
	u.log.Info("CosmosChainList:= ", cosmosChainList)
	return cosmosChainList
}

func (u *UtilConf) GetCosmosWalletSource(chain string) config.Source {
	for _, w := range u.conf.Cosmos.Cfg.Wallets {
		if w.ChainName == chain {
			return w.Source
		}
		u.log.Info(w.ChainName)
	}
	return config.Source{}
}

func (u *UtilConf) IsCosmos(str string) bool {
	if val, ok := CosmosEVMCompatibleChains[str]; ok {
		str = val
	}
	for _, v := range u.GetCosmosChains() {
		if v == str {
			u.log.Info("IsCosmos: ", true)
			return true
		}
	}
	u.log.Info("IsCosmos: ", false)
	return false
}

func (u *UtilConf) GetCosmosWalletInfo(chain string) config.CosmosWallets {
	for _, w := range u.conf.Cosmos.Cfg.Wallets {
		if w.ChainName == chain {
			return w
		}
		u.log.Debug(w.ChainName)
	}
	return config.CosmosWallets{}
}

func (u *UtilConf) RenameEVMCosmosCompatibleChains(chain string) string {
	if val, ok := CosmosEVMCompatibleChains[chain]; ok {
		chain = val
	}
	return chain
}
