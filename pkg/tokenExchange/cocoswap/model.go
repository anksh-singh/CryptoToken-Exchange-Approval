package cocoswap

type ExchangeToken struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
	ChainId  int    `json:"chainId"`
}
