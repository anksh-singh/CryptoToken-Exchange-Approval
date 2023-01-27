package _x

type Exchange0xModel struct {
	ChainID              int        `json:"chainId"`
	Price                string     `json:"price"`
	GuaranteedPrice      string     `json:"guaranteedPrice"`
	EstimatedPriceImpact string     `json:"estimatedPriceImpact"`
	To                   string     `json:"to"`
	Data                 string     `json:"data"`
	Value                string     `json:"value"`
	Gas                  string     `json:"gas"`
	EstimatedGas         string     `json:"estimatedGas"`
	GasPrice             string     `json:"gasPrice"`
	ProtocolFee          string     `json:"protocolFee"`
	MinimumProtocolFee   string     `json:"minimumProtocolFee"`
	BuyTokenAddress      string     `json:"buyTokenAddress"`
	SellTokenAddress     string     `json:"sellTokenAddress"`
	BuyAmount            string     `json:"buyAmount"`
	SellAmount           string     `json:"sellAmount"`
	Sources              []*Sources `json:"sources"`
	Orders               []*Orders  `json:"orders"`
	AllowanceTarget      string     `json:"allowanceTarget"`
	SellTokenToEthRate   string     `json:"sellTokenToEthRate"`
	BuyTokenToEthRate    string     `json:"buyTokenToEthRate"`
}

type Sources struct {
	Proportion        string   `json:"proportion"`
	IntermediateToken string   `json:"intermediateToken,omitempty"`
	Hops              []string `json:"hops,omitempty"`
	Name              string   `json:"name"`
}

type Orders struct {
	MakerToken   string   `json:"makerToken"`
	TakerToken   string   `json:"takerToken"`
	MakerAmount  string   `json:"makerAmount"`
	TakerAmount  string   `json:"takerAmount"`
	FillData     FillData `json:"fillData"`
	Source       string   `json:"source"`
	SourcePathID string   `json:"sourcePathId"`
	Type         int      `json:"type"`
}

type FillData struct {
	PoolID string `json:"poolId"`
	Vault  string `json:"vault"`
}

type ErrorResponse struct {
	Code             int64  `json:"code"`
	Reason           string `json:"reason"`
	ValidationErrors []struct {
		Code   int64  `json:"code"`
		Field  string `json:"field"`
		Reason string `json:"reason"`
	} `json:"validationErrors"`
}
