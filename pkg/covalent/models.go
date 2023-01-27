package covalent

import "time"

type AssetsResponseForCovalent struct {
	Data struct {
		Address       string    `json:"address"`
		UpdatedAt     time.Time `json:"updated_at"`
		NextUpdateAt  time.Time `json:"next_update_at"`
		QuoteCurrency string    `json:"quote_currency"`
		ChainID       int       `json:"chain_id"`
		Items         []struct {
			ContractDecimals     int         `json:"contract_decimals"`
			ContractName         string      `json:"contract_name"`
			ContractTickerSymbol string      `json:"contract_ticker_symbol"`
			ContractAddress      string      `json:"contract_address"`
			SupportsErc          interface{} `json:"supports_erc"`
			LogoURL              string      `json:"logo_url"`
			LastTransferredAt    interface{} `json:"last_transferred_at"`
			Type                 string      `json:"type"`
			Balance              string      `json:"balance"`
			Balance24H           string      `json:"balance_24h"`
			QuoteRate            float64     `json:"quote_rate"`
			QuoteRate24H         float64     `json:"quote_rate_24h"`
			Quote                float64     `json:"quote"`
			Quote24H             float64     `json:"quote_24h"`
			NftData              interface{} `json:"nft_data"`
		} `json:"items"`
		Pagination interface{} `json:"pagination"`
	} `json:"data"`
	Error        bool        `json:"error"`
	ErrorMessage interface{} `json:"error_message"`
	ErrorCode    interface{} `json:"error_code"`
}

type ListTransactionListForCovalent struct {
	Data struct {
		Address       string    `json:"address"`
		UpdatedAt     time.Time `json:"updated_at"`
		NextUpdateAt  time.Time `json:"next_update_at"`
		QuoteCurrency string    `json:"quote_currency"`
		ChainID       int       `json:"chain_id"`
		Items         []struct {
			BlockSignedAt    time.Time   `json:"block_signed_at"`
			BlockHeight      int         `json:"block_height"`
			TxHash           string      `json:"tx_hash"`
			TxOffset         int         `json:"tx_offset"`
			Successful       bool        `json:"successful"`
			FromAddress      string      `json:"from_address"`
			FromAddressLabel interface{} `json:"from_address_label"`
			ToAddress        string      `json:"to_address"`
			ToAddressLabel   interface{} `json:"to_address_label"`
			Value            string      `json:"value"`
			ValueQuote       float64     `json:"value_quote"`
			GasOffered       int         `json:"gas_offered"`
			GasSpent         int         `json:"gas_spent"`
			GasPrice         int64       `json:"gas_price"`
			GasQuote         float64     `json:"gas_quote"`
			GasQuoteRate     float64     `json:"gas_quote_rate"`
			LogEvents        []struct {
				BlockSignedAt              time.Time   `json:"block_signed_at"`
				BlockHeight                int         `json:"block_height"`
				TxOffset                   int         `json:"tx_offset"`
				LogOffset                  int         `json:"log_offset"`
				TxHash                     string      `json:"tx_hash"`
				RawLogTopics               []string    `json:"raw_log_topics"`
				SenderContractDecimals     int         `json:"sender_contract_decimals"`
				SenderName                 string      `json:"sender_name"`
				SenderContractTickerSymbol string      `json:"sender_contract_ticker_symbol"`
				SenderAddress              string      `json:"sender_address"`
				SenderAddressLabel         interface{} `json:"sender_address_label"`
				SenderLogoURL              string      `json:"sender_logo_url"`
				RawLogData                 string      `json:"raw_log_data"`
				Decoded                    struct {
					Name      string `json:"name"`
					Signature string `json:"signature"`
					Params    []struct {
						Name    string `json:"name"`
						Type    string `json:"type"`
						Indexed bool   `json:"indexed"`
						Decoded bool   `json:"decoded"`
						Value   string `json:"value"`
					} `json:"params"`
				} `json:"decoded"`
			} `json:"log_events"`
		} `json:"items"`
		Pagination struct {
			HasMore    bool        `json:"has_more"`
			PageNumber int         `json:"page_number"`
			PageSize   int         `json:"page_size"`
			TotalCount interface{} `json:"total_count"`
		} `json:"pagination"`
	} `json:"data"`
	Error        bool        `json:"error"`
	ErrorMessage interface{} `json:"error_message"`
	ErrorCode    interface{} `json:"error_code"`
}
