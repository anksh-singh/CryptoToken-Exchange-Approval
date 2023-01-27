package utils

import (
	"bridge-allowance/config"
	"strings"
)

func (u *UtilConf) GetNonEVMWalletInfo(chain string) *config.NonEVMChainInfo {
	for _, w := range u.conf.NonEVMConfig.NonEVMWallet {
		if strings.ToLower(w.ChainName) == strings.ToLower(chain) {
			return w
		}
		u.log.Debug(w.ChainName)
	}
	return &config.NonEVMChainInfo{}
}

func (u *UtilConf) IsNonEVM(str string) bool {
	for _, v := range u.GetNonEVMChains() {
		if v == str {
			u.log.Info("IsNonEVM: ", true)
			return true
		}
	}
	u.log.Info("IsNonEVM: ", false)
	return false
}
