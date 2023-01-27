package debankAPI

type LiquidityPool struct {
	IsActionable        bool                      `json:"is_actionable"`
	Type                string                    `json:"type"`
	AssetValue          string                    `json:"asset_value"`
	NetValue            string                    `json:"net_value"`
	PairAddress         string                    `json:"pair_address,omitempty"`
	PairName            string                    `json:"pair_name,omitempty"`
	PairBalance         string                    `json:"pair_balance,omitempty"`
	PoolShare           string                    `json:"pool_share,omitempty"`
	TokensSupplied      []*LiquidityPoolTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*LiquidityPoolTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*LiquidityPoolTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*LiquidityPoolTokenList `json:"underlying_token,omitempty"`
	StrikeTokenList     []*LiquidityPoolTokenList `json:"strike_token,omitempty"`
	CollateralTokenList []*LiquidityPoolTokenList `json:"collateral_token_list,omitempty"`
}

type Lending struct {
	IsActionable        bool                `json:"is_actionable"`
	Type                string              `json:"type"`
	AssetValue          string              `json:"asset_value"`
	NetValue            string              `json:"net_value"`
	TokensSupplied      []*LendingTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*LendingTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*LendingTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*LendingTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*LendingTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*LendingTokenList `json:"collateral_token_list,omitempty"`
}

type Yield struct {
	IsActionable        bool              `json:"is_actionable"`
	Type                string            `json:"type"`
	AssetValue          string            `json:"asset_value"`
	NetValue            string            `json:"net_value"`
	TokensSupplied      []*YieldTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*YieldTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*YieldTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*YieldTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*YieldTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*YieldTokenList `json:"collateral_token_list,omitempty"`
}

type Staked struct {
	IsActionable        bool               `json:"is_actionable"`
	Type                string             `json:"type"`
	AssetValue          string             `json:"asset_value"`
	NetValue            string             `json:"net_value"`
	TokensSupplied      []*StakedTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*StakedTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*StakedTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*StakedTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*StakedTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*StakedTokenList `json:"collateral_token_list,omitempty"`
}

type Farm struct {
	IsActionable        bool             `json:"is_actionable"`
	Type                string           `json:"type"`
	AssetValue          string           `json:"asset_value"`
	NetValue            string           `json:"net_value"`
	TokensSupplied      []*FarmTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*FarmTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*FarmTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*FarmTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*FarmTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*FarmTokenList `json:"collateral_token_list,omitempty"`
}

type Deposit struct {
	IsActionable        bool                `json:"is_actionable"`
	Type                string              `json:"type"`
	AssetValue          string              `json:"asset_value"`
	NetValue            string              `json:"net_value"`
	TokensSupplied      []*DepositTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*DepositTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*DepositTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*DepositTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*DepositTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*DepositTokenList `json:"collateral_token_list,omitempty"`
}

type Locked struct {
	IsActionable        bool               `json:"is_actionable"`
	Type                string             `json:"type"`
	AssetValue          string             `json:"asset_value"`
	NetValue            string             `json:"net_value"`
	TokensSupplied      []*LockedTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*LockedTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*LockedTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*LockedTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*LockedTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*LockedTokenList `json:"collateral_token_list,omitempty"`
}

type LeaveragedFarming struct {
	IsActionable        bool                          `json:"is_actionable"`
	Type                string                        `json:"type"`
	AssetValue          string                        `json:"asset_value"`
	NetValue            string                        `json:"net_value"`
	TokensSupplied      []*LeaveragedFarmingTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*LeaveragedFarmingTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*LeaveragedFarmingTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*LeaveragedFarmingTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*LeaveragedFarmingTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*LeaveragedFarmingTokenList `json:"collateral_token_list,omitempty"`
}

type Vesting struct {
	IsActionable        bool                `json:"is_actionable"`
	Type                string              `json:"type"`
	AssetValue          string              `json:"asset_value"`
	NetValue            string              `json:"net_value"`
	TokensSupplied      []*VestingTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*VestingTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*VestingTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*VestingTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*VestingTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*VestingTokenList `json:"collateral_token_list,omitempty"`
}

type Rewards struct {
	IsActionable        bool                `json:"is_actionable"`
	Type                string              `json:"type"`
	AssetValue          string              `json:"asset_value"`
	NetValue            string              `json:"net_value"`
	TokensSupplied      []*RewardsTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*RewardsTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*RewardsTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*RewardsTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*RewardsTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*RewardsTokenList `json:"collateral_token_list,omitempty"`
}

type Investment struct {
	IsActionable        bool                   `json:"is_actionable"`
	Type                string                 `json:"type"`
	AssetValue          string                 `json:"asset_value"`
	NetValue            string                 `json:"net_value"`
	TokensSupplied      []*InvestmentTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*InvestmentTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*InvestmentTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*InvestmentTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*InvestmentTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*InvestmentTokenList `json:"collateral_token_list,omitempty"`
}

type OptionsSeller struct {
	IsActionable        bool                      `json:"is_actionable"`
	Type                string                    `json:"type"`
	AssetValue          string                    `json:"asset_value"`
	NetValue            string                    `json:"net_value"`
	TokensSupplied      []*OptionsSellerTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*OptionsSellerTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*OptionsSellerTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*OptionsSellerTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*OptionsSellerTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*OptionsSellerTokenList `json:"collateral_token_list,omitempty"`
}

type OptionsBuyer struct {
	IsActionable        bool                     `json:"is_actionable"`
	Type                string                   `json:"type"`
	AssetValue          string                   `json:"asset_value"`
	NetValue            string                   `json:"net_value"`
	TokensSupplied      []*OptionsBuyerTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*OptionsBuyerTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*OptionsBuyerTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*OptionsBuyerTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*OptionsBuyerTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*OptionsBuyerTokenList `json:"collateral_token_list,omitempty"`
}

type InsuranceSeller struct {
	IsActionable        bool                        `json:"is_actionable"`
	Type                string                      `json:"type"`
	AssetValue          string                      `json:"asset_value"`
	NetValue            string                      `json:"net_value"`
	TokensSupplied      []*InsuranceSellerTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*InsuranceSellerTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*InsuranceSellerTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*InsuranceSellerTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*InsuranceSellerTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*InsuranceSellerTokenList `json:"collateral_token_list,omitempty"`
}

type InsuranceBuyer struct {
	IsActionable        bool                       `json:"is_actionable"`
	Type                string                     `json:"type"`
	AssetValue          string                     `json:"asset_value"`
	NetValue            string                     `json:"net_value"`
	TokensSupplied      []*InsuranceBuyerTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*InsuranceBuyerTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*InsuranceBuyerTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*InsuranceBuyerTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*InsuranceBuyerTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*InsuranceBuyerTokenList `json:"collateral_token_list,omitempty"`
}

type Perpetuals struct {
	IsActionable        bool                   `json:"is_actionable"`
	Type                string                 `json:"type"`
	AssetValue          string                 `json:"asset_value"`
	NetValue            string                 `json:"net_value"`
	TokensSupplied      []*PerpetualsTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*PerpetualsTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*PerpetualsTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*PerpetualsTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*PerpetualsTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*PerpetualsTokenList `json:"collateral_token_list,omitempty"`
}

type Others struct {
	IsActionable        bool               `json:"is_actionable"`
	Type                string             `json:"type"`
	AssetValue          string             `json:"asset_value"`
	NetValue            string             `json:"net_value"`
	TokensSupplied      []*OthersTokenList `json:"tokens_supplied,omitempty"`
	BorrowTokenList     []*OthersTokenList `json:"borrow_token_list,omitempty"`
	RewardTokenList     []*OthersTokenList `json:"reward_token_list,omitempty"`
	UnderLyingTokenList []*OthersTokenList `json:"underlying_token_list,omitempty"`
	StrikeTokenList     []*OthersTokenList `json:"strike_token_list,omitempty"`
	CollateralTokenList []*OthersTokenList `json:"collateral_token_list,omitempty"`
}
