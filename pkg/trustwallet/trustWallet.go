package trustwallet

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var trustWalletChainName = map[string]string{
	"bsc":       "smartchain",
	"ethereum":  "ethereum",
	"polygon":   "polygon",
	"avalanche": "avalanchec",
	"celo":      "celo",
	"fantom":    "fantom",
	"optimism":  "optimism",
	"arbitrum":  "arbitrum",
	"boba":      "boba",
	"aurora":    "aurora",
	"metis":     "metis",
	"band":      "band",
	"kava":      "kava",
	"osmosis":   "osmosis",
	"terra":     "terra",
	"matic":     "polygon",
	"polkadot":  "polkadot",
	"stafi":     "stafi",
	"algorand":  "algorand",
	"heco":      "heco",
	//"edgeware":  "edgeware", //Need to check
	"zilliqa":   "zilliqa",
	"solana":    "solana",
	"harmony":   "harmony",
	"elrond":    "elrond",
	"tomochain": "tomochain",
	//"xinfin":    "xdc", //trust wallet status abandoned
	"gnosis":    "xdai",
	"cronos":    "cronos",
	"moonriver": "moonriver",
	"moonbeam":  "moonbeam",
}

type ITrustWallet interface {
	GetTokenInfo(chain string, contractAddress string) (*ExchangeToken, error)
}
type TrustWallet struct {
	env         *config.Config
	logger      *zap.SugaredLogger
	httpRequest utils.IHttpRequest
}

func NewTrustWallet(env *config.Config, logger *zap.SugaredLogger, httpRequest utils.IHttpRequest) *TrustWallet {
	return &TrustWallet{env: env, logger: logger, httpRequest: httpRequest}
}

func (t *TrustWallet) GetTokenInfo(chain string, contractAddress string) (*ExchangeToken, error) {
	if trustWalletChain, ok := trustWalletChainName[chain]; ok {
		chain = trustWalletChain
	}
	url := fmt.Sprintf(t.env.TrustWallet.EndPoint+"/%s/info/info.json", chain)
	var tokenExchange ExchangeToken
	body, err := t.httpRequest.GetRequest(url)
	err = json.Unmarshal(body, &tokenExchange)
	if err != nil {
		t.logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Json Unmarshalling")
	}
	return &tokenExchange, err
}
