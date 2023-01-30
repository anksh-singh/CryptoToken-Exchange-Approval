package utils

import (
	"bridge-allowance/config"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	hrp "github.com/harmony-one/go-sdk/pkg/address"
	ioTex "github.com/iotexproject/iotex-core/pkg/util/addrutil"
	"github.com/shopspring/decimal"
	// "github.com/unstoppabledomains/resolution-go/v2"
	"github.com/wealdtech/go-ens/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"strings"
)

const EVM = "evm"

const (
	Wei             = 1
	GWei            = 1e9
	Ether           = 1e18
	HrpAddrPrefix   = "one"
	IoAddrPrefix    = "io"
	ZeroXAddrPrefix = "0x"
	XdcAddrPrefix   = "xdc"
)

// EnsDomains Supported ENS Deployments - https://docs.ens.domains/ens-deployments
var EnsDomains = [5]string{".eth", ".xyz", ".luxe", ".kred", ".art"}

func (u *UtilConf) GetEVMBridgeData(service string) []*config.BridgeNomenclature {
	var bridgeList []*config.BridgeNomenclature
	for _, w := range u.conf.EVM.Cfg.Wallets {
		if w.Bridge.Service == service {
			bridgeList = append(bridgeList, &config.BridgeNomenclature{
				ChainId:  w.Bridge.ChainId,
				ChainKey: w.Bridge.ChainKey,
				Service:  w.Bridge.Service,
			})
		}
	}
	return bridgeList
}

// GetEVMChains Load configured evm chains from config
// Returns an empty string array if no configurations are found
func (u *UtilConf) GetEVMChains() []string {
	var evmChainsList []string
	if len(u.conf.EVM.Cfg.Wallets) < 1 {
		return evmChainsList
	}
	for _, w := range u.conf.EVM.Cfg.Wallets {
		evmChainsList = append(evmChainsList, w.ChainName)
	}
	u.log.Info("EVMChainList:= ", evmChainsList)
	return evmChainsList
}
func (u *UtilConf) GetNonEVMChains() []string {
	var nonevmChainsList []string
	if len(u.conf.NonEVMConfig.NonEVMWallet) < 1 {
		return nonevmChainsList
	}
	for _, w := range u.conf.NonEVMConfig.NonEVMWallet {
		nonevmChainsList = append(nonevmChainsList, w.ChainName)
	}
	u.log.Info("nonEVMChainList:= ", nonevmChainsList)
	return nonevmChainsList
}
func (u *UtilConf) GetWalletSource(chain string) config.Source {
	for _, w := range u.conf.EVM.Cfg.Wallets {
		if w.ChainName == chain {
			u.log.Debugf("Wallet source for chain: %v = %v", w.ChainName, w.Source)
			return w.Source
		}
	}
	return config.Source{}
}

func (u *UtilConf) GetWalletInfo(chain string) config.Wallets {
	for _, w := range u.conf.EVM.Cfg.Wallets {
		if strings.ToLower(w.ChainName) == strings.ToLower(chain) ||
			strings.ToLower(w.NativeTokenInfo.ChainId) == strings.ToLower(chain) ||
			strings.ToLower(w.Bridge.ChainKey) == strings.ToLower(chain) {
			return w
		}
		u.log.Debug(w.ChainName)
	}
	return config.Wallets{}
}

func (u *UtilConf) IsEVM(str string) bool {
	for _, v := range u.GetEVMChains() {
		if v == str {
			u.log.Info("IsEVM: ", true)
			return true
		}
	}
	u.log.Info("IsEVM: ", false)
	return false
}


// ResolveENSAddress resolve an ENS(Ethereum Naming Service) address into a 0x address
func (u *UtilConf) ResolveENSAddress(domain string) (string, error) {
	hasSuffix := false
	//Pre-checks to Fail-fast
	for _, ensDomain := range EnsDomains {
		if strings.HasSuffix(domain, ensDomain) {
			hasSuffix = true
		}
	}
	if !hasSuffix {
		return domain, errors.New(fmt.Sprintf("domain: %v, not found", domain))
	}
	//TODO:Cache ENS addresses for a faster lookup
	client, err := ethclient.Dial(u.GetWalletInfo("ethereum").RPC)
	if err != nil {
		return domain, status.Errorf(codes.Internal, err.Error(), "Internal Error")
	}
	address, err := ens.Resolve(client, domain)
	return address.Hex(), err
}

// ResolveBech32Address resolves bech 32 address to 0x address
func (u *UtilConf) ResolveBech32Address(address string) string {
	if !strings.HasPrefix(address, HrpAddrPrefix) {
		return address
	}
	zeroXAddr, err := hrp.Bech32ToAddress(address)
	if err != nil {
		u.log.Errorf("Unable to resolve bech32 address: %v", address)
	}
	return zeroXAddr.Hex()
}

// ResolveIoAddress resolves IoTex address to 0x address
func (u *UtilConf) ResolveIoAddress(address string) string {
	if !strings.HasPrefix(address, IoAddrPrefix) {
		return address
	}
	zeroXAddr, err := ioTex.IoAddrToEvmAddr(address)
	if err != nil {
		u.log.Errorf("Unable to resolve ioTex address: %v", address)
	}
	return zeroXAddr.Hex()
}

// ResolveXDCAddress resolves xdc address to 0x address
func (u *UtilConf) ResolveXDCAddress(address string) string {
	if strings.HasPrefix(address, XdcAddrPrefix) {
		return strings.Replace(address, XdcAddrPrefix, ZeroXAddrPrefix, 1)
	}
	return address
}


// ToDecimal wei to decimals
func (u *UtilConf) ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)
	return result
}



