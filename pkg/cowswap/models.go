package cowswap

import "time"

type ExchangeCowResponse struct {
	Name      string    `json:"name"`
	Timestamp string `json:"timestamp"`
	Version Version `json:"version"`
	Tags Tags `json:"tags"`
	LogoURI  string   `json:"logoURI"`
	Keywords []string `json:"keywords"`
	Tokens Tokens `json:"tokens"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type Tags struct {}

type Tokens []struct {
	ChainID    int    `json:"chainId"`
	Address    string `json:"address"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   int    `json:"decimals"`
	LogoURI    string `json:"logoURI"`
	Extensions Extensions `json:"extensions,omitempty"`
}

type Extensions struct {
	BridgeInfo BridgeInfo `json:"bridgeInfo"`
}

type BridgeInfo struct {
	Num137 Num137 `json:"137"`
	Num42161 Num42161 `json:"42161"`
}

type Num137 struct {
	TokenAddress string `json:"tokenAddress"`
}

type Num42161 struct {
	TokenAddress string `json:"tokenAddress"`
}

type QuoteExchange struct {
	Quote struct {
		SellToken         string `json:"sellToken"`
		BuyToken          string `json:"buyToken"`
		Receiver          string `json:"receiver"`
		SellAmount        string `json:"sellAmount"`
		BuyAmount         string `json:"buyAmount"`
		ValidTo           int32    `json:"validTo"`
		AppData           string `json:"appData"`
		FeeAmount         string `json:"feeAmount"`
		Kind              string `json:"kind"`
		PartiallyFillable bool   `json:"partiallyFillable"`
		SellTokenBalance  string `json:"sellTokenBalance"`
		BuyTokenBalance   string `json:"buyTokenBalance"`
	} `json:"quote"`
	From       string      `json:"from"`
	Expiration time.Time   `json:"expiration"`
	ID         int32 `json:"id"`
}

type InputRequestData struct {
	SellToken string `json:"sellToken"`
	BuyToken  string `json:"buyToken"`
	Receiver  string `json:"receiver"`
	AppData   string `json:"appData"`
	PartiallyFillable bool `json:"partiallyFillable"`
	SellTokenBalance  string `json:"sellTokenBalance"`
	BuyTokenBalance   string `json:"buyTokenBalance"`
	From string `json:"from"`
	PriceQuality string `json:"priceQuality"`
	SigningScheme string `json:"signingScheme"`
	OnchainOrder bool `json:"onchainOrder"`
	Kind string `json:"kind"`
	SellAmountBeforeFee string `json:"sellAmountBeforeFee"`
}

type OrderResponse struct {
	UID string `json:"uid"`
}

type OrderStatusResponse struct {
	CreationDate                 time.Time   `json:"creationDate"`
	Owner                        string      `json:"owner"`
	UID                          string      `json:"uid"`
	AvailableBalance             interface{} `json:"availableBalance"`
	ExecutedBuyAmount            string      `json:"executedBuyAmount"`
	ExecutedSellAmount           string      `json:"executedSellAmount"`
	ExecutedSellAmountBeforeFees string      `json:"executedSellAmountBeforeFees"`
	ExecutedFeeAmount            string      `json:"executedFeeAmount"`
	Invalidated                  bool        `json:"invalidated"`
	Status                       string      `json:"status"`
	SettlementContract           string      `json:"settlementContract"`
	FullFeeAmount                string      `json:"fullFeeAmount"`
	IsLiquidityOrder             bool        `json:"isLiquidityOrder"`
	SellToken                    string      `json:"sellToken"`
	BuyToken                     string      `json:"buyToken"`
	Receiver                     string      `json:"receiver"`
	SellAmount                   string      `json:"sellAmount"`
	BuyAmount                    string      `json:"buyAmount"`
	ValidTo                      int         `json:"validTo"`
	AppData                      string      `json:"appData"`
	FeeAmount                    string      `json:"feeAmount"`
	Kind                         string      `json:"kind"`
	PartiallyFillable            bool        `json:"partiallyFillable"`
	SellTokenBalance             string      `json:"sellTokenBalance"`
	BuyTokenBalance              string      `json:"buyTokenBalance"`
	SigningScheme                string      `json:"signingScheme"`
	Signature                    string      `json:"signature"`
}

type CowTradeResponse []struct {
	BlockNumber          int    `json:"blockNumber"`
	LogIndex             int    `json:"logIndex"`
	OrderUID             string `json:"orderUid"`
	BuyAmount            string `json:"buyAmount"`
	SellAmount           string `json:"sellAmount"`
	SellAmountBeforeFees string `json:"sellAmountBeforeFees"`
	Owner                string `json:"owner"`
	BuyToken             string `json:"buyToken"`
	SellToken            string `json:"sellToken"`
	TxHash               string `json:"txHash"`
}