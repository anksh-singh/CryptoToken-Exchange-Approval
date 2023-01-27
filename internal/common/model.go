package common

import (
	_x "bridge-allowance/pkg/0x"
	"bridge-allowance/pkg/coingecko"
	"bridge-allowance/pkg/covalent"
	"bridge-allowance/pkg/cowswap"
	"bridge-allowance/pkg/customchain/zksync"
	"bridge-allowance/pkg/debridge"
	"bridge-allowance/pkg/dodo"
	"bridge-allowance/pkg/dzap"
	"bridge-allowance/pkg/lifi"
	"bridge-allowance/pkg/oneinch"
	"bridge-allowance/pkg/openApi"
	"bridge-allowance/pkg/socket"
	"bridge-allowance/pkg/tokenExchange/cocoswap"
	"bridge-allowance/pkg/tokenExchange/dodoEth"
	"bridge-allowance/pkg/tokenExchange/ubeswap"
	"bridge-allowance/pkg/trustwallet"
	"bridge-allowance/pkg/unmarshal"
	"bridge-allowance/pkg/zeroswap"
	"bridge-allowance/utils"
)

// TODO: Move services to a common model
type Services struct {
	Http                     *utils.HttpRequest
	CoinGecko                *coingecko.CoinGecko
	Helper                   *utils.Helpers
	Covalent                 *covalent.CovalentService
	Unmarshall               *unmarshal.UnmarshallService
	ZeroX                    *_x.OXService
	CocoSwapTokenExchange    *cocoswap.TokenExchangeStruct
	UniSwapTokenExchange     *ubeswap.TokenExchangeStruct
	DoDoExTokenExchange      *dodoEth.TokenExchangeStruct
	DoDoExTokenExchangeCache *dodoEth.TokenExchangeStructCache
	Debank                   *openApi.OpenAPI
	DodoSwap                 *dodo.ServiceDodo
	TrustWallet              trustwallet.ITrustWallet
	LiFi                     lifi.ILiFi
	Socket                   socket.ISocket
	OneInch                  *oneinch.OneInchService
	DeBridge                 debridge.IdeBridge
	ZeroSwap                 *zeroswap.ZeroSwapService
	CowSwap                  *cowswap.CowSwapService
	DZap                     dzap.IDZap
	Zksync                   zksync.IZksync
}
