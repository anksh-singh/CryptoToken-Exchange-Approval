package zeroswap

type ZeroSwapModel struct {
	FromToken    SellTokenInfo `json:"sellTokenInfo"`
	ToToken      BuyTokenInfo  `json:"buyTokenInfo"`
	BuyAmount    string        `json:"buyAmount"`
	SellAmount   string        `json:"sellAmount"`
	EstimatedGas string        `json:"estimatedGas"`
	// Tx           TxData          `json:"tx"`
	GasPrice             string `json:"gasPrice"`
	To                   string `json:"to"`
	Price                string `json:"price"`
	Data                 string `json:"data"`
	EstimatedPriceImpact string `json:"estimatedPriceImpact"`
}

type SellTokenInfo struct {
	Symbol   string   `json:"symbol"`
	Name     string   `json:"name"`
	Decimals int      `json:"decimals"`
	Address  string   `json:"address"`
	LogoURI  string   `json:"logoURI"`
	Tags     []string `json:"tags"`
}

type BuyTokenInfo struct {
	Symbol   string   `json:"symbol"`
	Name     string   `json:"name"`
	Decimals int      `json:"decimals"`
	Address  string   `json:"address"`
	LogoURI  string   `json:"logoURI"`
	Eip2612  bool     `json:"eip2612"`
	Tags     []string `json:"tags"`
}

type Protocols struct {
	Name             string `json:"name"`
	Part             int    `json:"part"`
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
}

type TxData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	Gas      int    `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

type InputRequest struct {
	Chain        string `json:"chain"`
	TakerAddress string `json:"takerAddress"`
	SellToken    string `json:"sellToken"`
	BuyToken     string `json:"buyToken"`
	SellAmount   string `json:"sellAmount"`
	Slippage     string `json:"slippage"`
	ExchangeType string `json:"exchangeType"`
}

type TradeCountInfo struct {
	ChainId        int    `json:"chainId"`
	ChainName      string `json:"chainName"`
	Account        string `json:"account"`
	FreeTradeCount string `json:"freeTradeCount"`
}

type SignatureRequest struct {
	Account           string `json:"account"`
	SellTokenAddress  string `json:"sellTokenAddress"`
	BuyTokenAddress   string `json:"buyTokenAddress"`
	SlippageTolerance string `json:"slippageTolerance"`
	SellAmount        string `json:"sellAmount"`
}

type WebClient struct {
	ChainName string
	RPC       string
}

type ChainInfo struct {
	Name         string
	ChainId      int
	TokenSymbol  string
	TokenAddress string
}

type ZeroSwapPayload struct {
	ChainId             string `json:"chainId"`
	Signature           string `json:"signature"`
	BuyToken            string `json:"buyToken"`
	SellToken           string `json:"sellToken"`
	BuyTokenAddress     string `json:"buyTokenAddress"`
	SellTokenAddress    string `json:"sellTokenAddress"`
	SellAmount          string `json:"sellAmount"`
	Recipient           string `json:"recipient"`
	SlippageTolerance   string `json:"slippageTolerance"`
	TransactionDeadline string `json:"transactionDeadline"`
	AffiliateAddress    string `json:"affiliateAddress"`
}

type RespError struct {
	Name    string        `json:"name"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors"`
}
