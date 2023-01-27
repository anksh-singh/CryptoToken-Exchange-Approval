package oneinch

type ExchangeOneInchModel struct {
	FromToken FromTokenInfo `json:"fromToken"`
	ToToken ToTokenInfo `json:"toToken"`
	ToTokenAmount   string `json:"toTokenAmount"`
	FromTokenAmount string `json:"fromTokenAmount"`
	Protocols       [][][]Protocols `json:"protocols"`
	EstimatedGas int `json:"estimatedGas"`
	Tx TxData `json:"tx"`
}

type FromTokenInfo struct {
	Symbol   string   `json:"symbol"`
	Name     string   `json:"name"`
	Decimals int      `json:"decimals"`
	Address  string   `json:"address"`
	LogoURI  string   `json:"logoURI"`
	Tags     []string `json:"tags"`
}

type ToTokenInfo struct {
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
	Chain         string `json:"chain"`
	TakerAddress  string `json:"takerAddress"`
	SellToken     string `json:"sellToken"`
	BuyToken      string `json:"buyToken"`
	SellAmount    string `json:"sellAmount"`
	Slippage      string `json:"slippage"`
	ExchangeType  string `json:"exchangeType"`
}

type ErrorResponse struct {
	StatusCode  int    `json:"statusCode"`
	Error       string `json:"error"`
	Description string `json:"description"`
	Meta        []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"meta"`
	RequestID string `json:"requestId"`
}