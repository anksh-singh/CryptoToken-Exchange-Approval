package utils

type TokenDetail struct {
	TokenName         string `json:"token_name"`
	TokenSymbol       string `json:"token_symbol"`
	TokenDecimals     string `json:"token_decimals"`
	TokenAddress      string `json:"token_address"`
	TokenLogoURL      string `json:"token_logo_url"`
	TokenListedCount  string `json:"token_listed_count"`
	TokenPrice        string `json:"token_price"`
	TokenLastActivity string `json:"token_last_activity"`
	TokenWebsite      string `json:"token_website"`
	TokenTradeVolume  string `json:"token_trade_volume"`
	LogoURL           string `json:"logo_url"`
}
