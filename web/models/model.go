package models

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HTTPResponse struct {
	Err HTTPError `json:"error"`
}

type GetTokenPrice struct {
	Chain    string `json:"ids"`
	Currency string `json:"vs_currencies"`
}

type SendTxBody struct {
	Message string `json:"message"`
}

type MultiSwapRequestBody struct {
	SwapParams []struct {
		SellAmount string `json:"sell_amount"`
		SellToken  string `json:"sell_token"`
		BuyToken   string `json:"buy_token"`
		Slippage   string `json:"slippage"`
	} `json:"swap_params"`
}

type SimulateTxBody struct {
	TxBytes string `json:"tx_bytes"`
}

type UnsErrorObject struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type GasLessSwapBody struct {
	TakerAddress          string `json:"takerAddress"`
	SellToken             string `json:"sellToken"`
	BuyToken              string `json:"buyToken"`
	SellAmount            string `json:"sellAmount"`
	Signature             string `json:"signature"`
	Slippage              string `json:"slippage"`
	Receiver              string `json:"receiver"`
	BuyAmount             string `json:"buyAmount"`
	ValidTo               int    `json:"validTo"`
	AppData               string `json:"appData"`
	FeeAmount             string `json:"feeAmount"`
	Kind                  string `json:"kind"`
	PartiallyFillable     bool   `json:"partiallyFillable"`
	SellTokenBalance      string `json:"sellTokenBalance"`
	BuyTokenBalance       string `json:"buyTokenBalance"`
	SigningScheme         string `json:"signingScheme"`
	From                  string `json:"from"`
	QuoteId               int    `json:"quoteId"`
	TransactionalDeadline string `json:"transactionalDeadline"`
}

type CosmosSendTxRequest struct {
	TxBytes string `json:"tx_bytes"`
	Mode    string `json:"mode"`
}

type CosmosSendTxBody struct {
	Message string `json:"message"`
}

type BluzelleSendTxRequest struct {
	Tx   interface{} `json:"tx"`
	Mode string      `json:"mode"`
}

type GetPositionRequest []struct {
	Chain       string   `json:"chain"`
	Address     string   `json:"address"`
	ProtocolIds []string `json:"protocol_ids"`
}


