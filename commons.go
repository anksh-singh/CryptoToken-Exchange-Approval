package bridge_allowance

import (
	"bridge-allowance/config"
	"errors"
	"fmt"
)

const EVM = "evm"

// GetEVMChains Load configured evm chains from config
// Returns an empty string array if no configurations are found
func GetEVMChains() []string {
	config := config.LoadConfig("", "")
	var evmChainsList []string
	if len(config.EVM.Cfg.Wallets) < 1 {
		return evmChainsList
	}
	for _, w := range config.EVM.Cfg.Wallets {
		evmChainsList = append(evmChainsList, w.ChainName)
	}
	return evmChainsList
}

// GetWalletChainInfo get chain info for a chain name
func GetWalletChainInfo(cfg config.Config, chain string) config.Wallets {
	for _, w := range cfg.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			return w
		}
	}
	return config.Wallets{}
}

// GetChainId get chain id for chain name
func GetChainId(cfg config.Config, chain string) (int, error) {
	for _, w := range cfg.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			return w.ChainID, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Chain: %v not supported", chain))
}

// GetDebankId Get debank equavalent nomenclature for a chain name
func GetDebankId(cfg config.Config, chainName string) (string, error) {
	for _, w := range cfg.EVM.Cfg.Wallets {
		if w.ChainName == chainName {
			return w.Nomenclature.Debank, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Chain: %v not supported", chainName))
}

// GetUnmarshalId Get unmarshal equavalent nomenclature for a chain name
func GetUnmarshalId(cfg config.Config, chainName string) (string, error) {
	for _, w := range cfg.EVM.Cfg.Wallets {
		if w.ChainName == chainName {
			return w.Nomenclature.Unmarshal, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Chain: %v not supported", chainName))
}

// GetWalletSource get wallet source for a chain name
func GetWalletSource(cfg config.Config, chain string) config.Source {
	for _, w := range cfg.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			return w.Source
		}
	}
	return config.Source{}
}

func IsEVM(str string) bool {
	for _, v := range GetEVMChains() {
		if v == str {
			return true
		}
	}
	return false
}
