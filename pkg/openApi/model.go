package openApi

type PositionalDataModel struct {
	TotalUsdValue float64      `json:"total_usd_value"`
	ChainList     []*ChainData `json:"chain_list"`
}

type ChainData struct {
	ID                     string  `json:"id"`
	CommunityID            int     `json:"community_id"`
	Name                   string  `json:"name"`
	NativeTokenID          string  `json:"native_token_id"`
	LogoURL                string  `json:"logo_url"`
	WrappedTokenID         string  `json:"wrapped_token_id"`
	IsSupportBalanceChange bool    `json:"is_support_balance_change"`
	UsdValue               float64 `json:"usd_value"`
}
