package jupiter

type SolanaTokenList struct {
	Name    string `json:"name"`
	LogoURI string `json:"logoURI"`
	Tokens   TokenList  `json:"tokens"`
}

type TokenList []struct {
	ChainId int  `json:"chainID"`
	Address string  `json:"address"`
	Symbol string   `json:"symbol"`
	Name   string    `json:"name"`
	Decimals int     `json:"decimals"`
	LogoURI  string    `json:"logoURI"`
	Tags    []string   `json:"tags"`
}


type ExchangeQuoteRes struct {
	Data []ExchangeQuoteData `json:"data"`
}

type  ExchangeQuoteData struct {
	InAmount int  `json:"inAmount"`
	OutAmount int  `json:"outAmount"`
	Amount  int  `json:"amount"`
	OutAmountWithSlippage int `json:"outAmountWithSlippage"`
	OtherAmountThreshold int `json:"otherAmountThreshold"`
	SwapMode  string `json:"swapMode"`
	PriceImpactPct float64  `json:"priceImpactPct"`
	MarketInfos []MarketInfoData `json:"marketInfos"`
}

type MarketInfoData struct {
	Id string `json:"id"`
	Label string `json:"label"`
	InputMint string 	`json:"inputMint"`
	OutputMint string `json:"outputMint"`
	NotEnoughLiquidity bool `json:"notEnoughLiquidity"`
	InAmount int  `json:"inAmount"`
	OutAmount int  `json:"outAmount"`
	PriceImpactPct float64  `json:"priceImpactPct"`
	LpFee  LpFee  `json:"lpFee"`
	PlatformFee PlatformFee `json:"platformFee"`
}

type LpFee struct {
	Amount float64 `json:"amount"`
	Mint string `json:"mint"`
	Pct float64  `json:"pct"`
}

type PlatformFee struct {
	Amount float64 `json:"amount"`
	Mint string `json:"mint"`
	Pct float64  `json:"pct"`
}

type SwapRequest struct {
	Route         ExchangeQuoteData `json:"route"`
	UserPublicKey string            `json:"userPublicKey"`
	WrapUnwrapSOL bool              `json:"wrapUnwrapSOL"`
	FeeAccount    string            `json:"feeAccount"`
}

type SwapResponse struct {
	SetupTransaction string `json:"setupTransaction"`
	SwapTransaction  string  `json:"swapTransaction"`
	CleanupTransaction string `json:"cleanupTransaction"`
}