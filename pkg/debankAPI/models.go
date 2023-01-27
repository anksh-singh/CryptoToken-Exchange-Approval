package debankAPI

// DebankPositionResponse Debank Position API Response Structure
type DebankPositionResponse []struct {
	ID                    string          `json:"id"`
	Chain                 string          `json:"chain"`
	Name                  string          `json:"name"`
	SiteURL               string          `json:"site_url"`
	LogoURL               string          `json:"logo_url"`
	HasSupportedPortfolio bool            `json:"has_supported_portfolio"`
	Tvl                   float64         `json:"tvl"`
	PortfolioItemList     []PortfolioItem `json:"portfolio_item_list"`
}

type PortfolioItem struct {
	Stats         Stats              `json:"stats"`
	AssetDict     map[string]float64 `json:"asset_dict"`
	UpdateAt      float64            `json:"update_at"`
	Name          string             `json:"name"`
	DetailTypes   []string           `json:"detail_types"`
	Detail        Detail             `json:"detail"`
	ProxyDetail   ProxyDetail        `json:"proxy_detail"`
	Pool          Pool               `json:"pool"`
	PositionIndex string             `json:"position_index"`
}

type Stats struct {
	AssetUsdValue float64 `json:"asset_usd_value"`
	DebtUsdValue  float64 `json:"debt_usd_value"`
	NetUsdValue   float64 `json:"net_usd_value"`
}

type Detail struct {
	SupplyTokenList     TokenList   `json:"supply_token_list,omitempty"`
	BorrowTokenList     TokenList   `json:"borrow_token_list,omitempty"`
	RewardTokenList     TokenList   `json:"reward_token_list,omitempty"`
	UnderlyingToken     TokenDetail `json:"underlying_token,omitempty"`
	StrikeToken         TokenDetail `json:"strike_token,omitempty"`
	CollateralTokenList TokenList   `json:"collateral_token_list,omitempty"`
	HealthRate          float64     `json:"health_rate"`
}

type ProxyDetail struct {
	Project         Project `json:"project"`
	ProxyContractID string  `json:"proxy_contract_id"`
}

type Project struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	SiteURL string `json:"site_url"`
	LogoURL string `json:"logo_url"`
}

type Pool struct {
	ID         string      `json:"id"`
	Chain      string      `json:"chain"`
	ProjectID  string      `json:"project_id"`
	AdapterID  string      `json:"adapter_id"`
	Controller string      `json:"controller"`
	Index      interface{} `json:"index"`
	TimeAt     int         `json:"time_at"`
}
type TokenDetail struct {
	ID              string  `json:"id"`
	Chain           string  `json:"chain"`
	Amount          float64 `json:"amount"`
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	DisplaySymbol   string  `json:"display_symbol"`
	OptimizedSymbol string  `json:"optimized_symbol"`
	Decimals        int     `json:"decimals"`
	LogoURL         string  `json:"logo_url"`
	ProtocolID      string  `json:"protocol_id"`
	Price           float64 `json:"price"`
	IsVerified      bool    `json:"is_verified"`
	IsCore          bool    `json:"is_core"`
	IsWallet        bool    `json:"is_wallet"`
	TimeAt          float64 `json:"time_at"`
}
type TokenList []struct {
	ID              string  `json:"id"`
	Chain           string  `json:"chain"`
	Amount          float64 `json:"amount"`
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	DisplaySymbol   string  `json:"display_symbol"`
	OptimizedSymbol string  `json:"optimized_symbol"`
	Decimals        int     `json:"decimals"`
	LogoURL         string  `json:"logo_url"`
	ProtocolID      string  `json:"protocol_id"`
	Price           float64 `json:"price"`
	IsVerified      bool    `json:"is_verified"`
	IsCore          bool    `json:"is_core"`
	IsWallet        bool    `json:"is_wallet"`
	TimeAt          float64 `json:"time_at"`
}

type ResponseItem struct {
	ID                    string          `json:"id"`
	Chain                 string          `json:"chain"`
	Name                  string          `json:"name"`
	SiteURL               string          `json:"site_url"`
	LogoURL               string          `json:"logo_url"`
	HasSupportedPortfolio bool            `json:"has_supported_portfolio"`
	Tvl                   float64         `json:"tvl"`
	PortfolioItemList     []PortfolioItem `json:"portfolio_item_list"`
}

// Ends here

type PositionResponse struct {
	Positions []*Positions `json:"positions"`
}

type Positions struct {
	Chain     string       `json:"chain"`
	Address   string       `json:"address"`
	Protocols []*Protocols `json:"protocols"`
}
type Protocols struct {
	ProtocolId   string     `json:"protocol_id"`
	Name         string     `json:"name"`
	SiteUrl      string     `json:"site_url"`
	LogoUrl      string     `json:"logo_url"`
	IsActionable bool       `json:"is_actionable"`
	Portfolio    *Portfolio `json:"portfolio"`
}
type Portfolio struct {
	NetUsdValue       string               `json:"net_usd_value"`
	Lending           []*Lending           `json:"lending,omitempty"`
	LiquidityPool     []*LiquidityPool     `json:"liquidity_pool,omitempty"`
	Yield             []*Yield             `json:"yield,omitempty"`
	Staked            []*Staked            `json:"staked,omitempty"`
	Farm              []*Farm              `json:"farm,omitempty"`
	Deposit           []*Deposit           `json:"deposit,omitempty"`
	Locked            []*Locked            `json:"locked,omitempty"`
	LeaveragedFarming []*LeaveragedFarming `json:"leaveraged_farming,omitempty"`
	Vesting           []*Vesting           `json:"vesting,omitempty"`
	Rewards           []*Rewards           `json:"rewards,omitempty"`
	Investment        []*Investment        `json:"investment,omitempty"`
	OptionsSeller     []*OptionsSeller     `json:"options_seller,omitempty"`
	OptionsBuyer      []*OptionsBuyer      `json:"options_buyer,omitempty"`
	InsuranceSeller   []*InsuranceSeller   `json:"insurance_seller,omitempty"`
	InsuranceBuyer    []*InsuranceBuyer    `json:"insurance_buyer,omitempty"`
	Perpetuals        []*Perpetuals        `json:"perpetuals,omitempty"`
	Others            []*Others            `json:"others,omitempty"`
}

type ProtocolActions struct {
	Deposit           bool `json:"Deposit"`
	Farming           bool `json:"Farming"`
	Investment        bool `json:"Investment"`
	Lending           bool `json:"Lending"`
	LiquidityPool     bool `json:"Liquidity Pool"`
	OptionsSeller     bool `json:"Options Seller"`
	Rewards           bool `json:"Rewards"`
	Vesting           bool `json:"Vesting"`
	Yield             bool `json:"Yield"`
	Staked            bool `json:"Staked"`
	Locked            bool `json:"Locked"`
	LeaveragedFarming bool `json:"Leaveraged Farming"`
	OptionsBuyer      bool `json:"Options Buyer"`
	InsuranceSeller   bool `json:"Insurance Seller"`
	InsuranceBuyer    bool `json:"Insurance Buyer"`
	Perpetuals        bool `json:"Perpetuals"`
	Others            bool `json:"Others"`
}
