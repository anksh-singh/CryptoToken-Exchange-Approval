package v1Proxy

type BridgeChains struct {
	Chains []struct {
		ChainID   int64  `json:"chain_id"`
		ChainType string `json:"chain_type"`
		Coin      string `json:"coin"`
		LogoURL   string `json:"logo_url"`
		Mainnet   bool   `json:"mainnet"`
		Name      string `json:"name"`
		Symbol    string `json:"symbol"`
	} `json:"chains"`
}

type BridgeExchangeTokens struct {
	Tokens []struct {
		TokenAddress  string `json:"token_address"`
		TokenDecimals int64  `json:"token_decimals"`
		TokenLogoURL  string `json:"token_logo_url"`
		TokenName     string `json:"token_name"`
		TokenSymbol   string `json:"token_symbol"`
	} `json:"tokens"`
}

type BridgeQuote struct {
	BridgeFee struct {
		Amount          string `json:"amount"`
		AmountUsd       string `json:"amount_usd"`
		ContractAddress string `json:"contract_address"`
		Symbol          string `json:"symbol"`
		TokenDecimals   int64  `json:"token_decimals"`
	} `json:"bridge_fee"`
	Estimate struct {
		ApproveAddress    string  `json:"approve_address"`
		ExecutionDuration float64 `json:"execution_duration"`
		FromAmount        string  `json:"from_amount"`
		FromAmountUsd     string  `json:"from_amount_usd"`
		FromTokenDecimals int64   `json:"from_token_decimals"`
		ToAmount          string  `json:"to_amount"`
		ToAmountMin       string  `json:"to_amountMin"`
		ToAmountUsd       string  `json:"to_amount_usd"`
		ToTokenDecimals   int64   `json:"to_token_decimals"`
	} `json:"estimate"`
	Tool        string `json:"tool"`
	ToolDetails struct {
		Key     string `json:"key"`
		LogoURL string `json:"logo_url"`
		Name    string `json:"name"`
	} `json:"tool_details"`
}

type BridgeTransaction struct {
	Tool        string `json:"tool"`
	ToolDetails struct {
		Key     string `json:"key"`
		LogoURL string `json:"logo_url"`
		Name    string `json:"name"`
	} `json:"tool_details"`
	TransactionRequest struct {
		Data     string `json:"data"`
		GasLimit string `json:"gas_limit"`
		To       string `json:"to"`
		Value    string `json:"value"`
	} `json:"transaction_request"`
}

type BridgeTransactionStatus struct {
	IsSuccess bool   `json:"isSuccess"`
	Msg       string `json:"msg"`
	Status    string `json:"status"`
	TxHash    string `json:"txHash"`
}

type SuccessErrorV1 struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type ErrorResponse struct {
	Errors  []interface{} `json:"errors"`
	Message string        `json:"message"`
	Name    string        `json:"name"`
}
