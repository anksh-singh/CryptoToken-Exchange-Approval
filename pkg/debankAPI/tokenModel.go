package debankAPI

type LiquidityPoolTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
	PooledValue   string `json:"pooled_value,omitempty"`
}

type LendingTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type YieldTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type StakedTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type FarmTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type DepositTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type LockedTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type LeaveragedFarmingTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type VestingTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type RewardsTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type InvestmentTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type OptionsSellerTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type OptionsBuyerTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type InsuranceSellerTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type InsuranceBuyerTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type PerpetualsTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}

type OthersTokenList struct {
	TokenAddress  string `json:"token_address"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
	LogoUrl       string `json:"logo_url"`
	Balance       string `json:"balance"`
	QuoteRate     string `json:"quote_rate"`
	QuotePrice    string `json:"quote_price"`
}
