package core

type SendTransaction struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type ErrorConfig struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SendTransactionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Accounts      interface{}   `json:"accounts"`
			Err           string        `json:"err"`
			Logs          []interface{} `json:"logs"`
			UnitsConsumed int           `json:"unitsConsumed"`
		} `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
}

type TokenBalance struct {
	ContractName         string  `protobuf:"bytes,1,opt,name=contract_name,json=contractName,proto3" json:"contract_name"`
	ContractTickerSymbol string  `protobuf:"bytes,2,opt,name=contract_ticker_symbol,json=contractTickerSymbol,proto3" json:"contract_ticker_symbol"`
	ContractDecimals     int32   `protobuf:"varint,3,opt,name=contract_decimals,json=contractDecimals,proto3" json:"contract_decimals"`
	ContractAddress      string  `protobuf:"bytes,4,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address"`
	Coin                 int64   `protobuf:"varint,5,opt,name=coin,proto3" json:"coin"`
	Type                 string  `protobuf:"bytes,6,opt,name=type,proto3" json:"type"`
	Balance              string  `protobuf:"bytes,7,opt,name=balance,proto3" json:"balance"`
	Quote                float64 `protobuf:"bytes,8,opt,name=quote,proto3" json:"quote"`
	QuoteRate            float64 `protobuf:"fixed64,9,opt,name=quote_rate,json=quoteRate,proto3" json:"quote_rate"`
	LogoUrl              string  `protobuf:"bytes,10,opt,name=logo_url,json=logoUrl,proto3" json:"logo_url"`
	QuoteRate_24H        string  `protobuf:"bytes,11,opt,name=quote_rate_24h,json=quoteRate24h,proto3" json:"quote_rate_24h"`
	QuotePctChange_24H   float64 `protobuf:"fixed64,12,opt,name=quote_pct_change_24h,json=quotePctChange24h,proto3" json:"quote_pct_change_24h"`
}
