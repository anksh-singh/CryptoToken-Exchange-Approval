package debankAPI

// GetLendingEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetLendingEmptyActionInfo(actionInfo []*Lending, actionType string, actionable bool) []*Lending {
	var data []*Lending
	if len(actionInfo) == 0 {
		data = append(data, &Lending{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*LendingTokenList{},
			BorrowTokenList:     []*LendingTokenList{},
			RewardTokenList:     []*LendingTokenList{},
			UnderLyingTokenList: []*LendingTokenList{},
			StrikeTokenList:     []*LendingTokenList{},
			CollateralTokenList: []*LendingTokenList{},
		})
	}
	return data
}

// GetLiquidityEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetLiquidityEmptyActionInfo(actionInfo []*LiquidityPool, actionType string, actionable bool) []*LiquidityPool {
	var data []*LiquidityPool
	if len(actionInfo) == 0 {
		data = append(data, &LiquidityPool{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*LiquidityPoolTokenList{},
			BorrowTokenList:     []*LiquidityPoolTokenList{},
			RewardTokenList:     []*LiquidityPoolTokenList{},
			UnderLyingTokenList: []*LiquidityPoolTokenList{},
			StrikeTokenList:     []*LiquidityPoolTokenList{},
			CollateralTokenList: []*LiquidityPoolTokenList{},
		})
	}
	return data
}

// GetDepositEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetDepositEmptyActionInfo(actionInfo []*Deposit, actionType string, actionable bool) []*Deposit {
	var data []*Deposit
	if len(actionInfo) == 0 {
		data = append(data, &Deposit{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*DepositTokenList{},
			BorrowTokenList:     []*DepositTokenList{},
			RewardTokenList:     []*DepositTokenList{},
			UnderLyingTokenList: []*DepositTokenList{},
			StrikeTokenList:     []*DepositTokenList{},
			CollateralTokenList: []*DepositTokenList{},
		})
	}
	return data
}

// GetFarmEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetFarmEmptyActionInfo(actionInfo []*Farm, actionType string, actionable bool) []*Farm {
	var data []*Farm
	if len(actionInfo) == 0 {
		data = append(data, &Farm{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*FarmTokenList{},
			BorrowTokenList:     []*FarmTokenList{},
			RewardTokenList:     []*FarmTokenList{},
			UnderLyingTokenList: []*FarmTokenList{},
			StrikeTokenList:     []*FarmTokenList{},
			CollateralTokenList: []*FarmTokenList{},
		})
	}
	return data
}

// GetInvestmentEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetInvestmentEmptyActionInfo(actionInfo []*Investment, actionType string, actionable bool) []*Investment {
	var data []*Investment
	if len(actionInfo) == 0 {
		data = append(data, &Investment{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*InvestmentTokenList{},
			BorrowTokenList:     []*InvestmentTokenList{},
			RewardTokenList:     []*InvestmentTokenList{},
			UnderLyingTokenList: []*InvestmentTokenList{},
			StrikeTokenList:     []*InvestmentTokenList{},
			CollateralTokenList: []*InvestmentTokenList{},
		})
	}
	return data
}

// GetOptionsSellerEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetOptionsSellerEmptyActionInfo(actionInfo []*OptionsSeller, actionType string, actionable bool) []*OptionsSeller {
	var data []*OptionsSeller
	if len(actionInfo) == 0 {
		data = append(data, &OptionsSeller{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*OptionsSellerTokenList{},
			BorrowTokenList:     []*OptionsSellerTokenList{},
			RewardTokenList:     []*OptionsSellerTokenList{},
			UnderLyingTokenList: []*OptionsSellerTokenList{},
			StrikeTokenList:     []*OptionsSellerTokenList{},
			CollateralTokenList: []*OptionsSellerTokenList{},
		})
	}
	return data
}

// GetRewardsEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetRewardsEmptyActionInfo(actionInfo []*Rewards, actionType string, actionable bool) []*Rewards {
	var data []*Rewards
	if len(actionInfo) == 0 {
		data = append(data, &Rewards{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*RewardsTokenList{},
			BorrowTokenList:     []*RewardsTokenList{},
			RewardTokenList:     []*RewardsTokenList{},
			UnderLyingTokenList: []*RewardsTokenList{},
			StrikeTokenList:     []*RewardsTokenList{},
			CollateralTokenList: []*RewardsTokenList{},
		})
	}
	return data
}

// GetVestingEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetVestingEmptyActionInfo(actionInfo []*Vesting, actionType string, actionable bool) []*Vesting {
	var data []*Vesting
	if len(actionInfo) == 0 {
		data = append(data, &Vesting{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*VestingTokenList{},
			BorrowTokenList:     []*VestingTokenList{},
			RewardTokenList:     []*VestingTokenList{},
			UnderLyingTokenList: []*VestingTokenList{},
			StrikeTokenList:     []*VestingTokenList{},
			CollateralTokenList: []*VestingTokenList{},
		})
	}
	return data
}

// GetYieldEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetYieldEmptyActionInfo(actionInfo []*Yield, actionType string, actionable bool) []*Yield {
	var data []*Yield
	if len(actionInfo) == 0 {
		data = append(data, &Yield{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*YieldTokenList{},
			BorrowTokenList:     []*YieldTokenList{},
			RewardTokenList:     []*YieldTokenList{},
			UnderLyingTokenList: []*YieldTokenList{},
			StrikeTokenList:     []*YieldTokenList{},
			CollateralTokenList: []*YieldTokenList{},
		})
	}
	return data
}

// GetStakedEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetStakedEmptyActionInfo(actionInfo []*Staked, actionType string, actionable bool) []*Staked {
	var data []*Staked
	if len(actionInfo) == 0 {
		data = append(data, &Staked{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*StakedTokenList{},
			BorrowTokenList:     []*StakedTokenList{},
			RewardTokenList:     []*StakedTokenList{},
			UnderLyingTokenList: []*StakedTokenList{},
			StrikeTokenList:     []*StakedTokenList{},
			CollateralTokenList: []*StakedTokenList{},
		})
	}
	return data
}

// GetLockedEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetLockedEmptyActionInfo(actionInfo []*Locked, actionType string, actionable bool) []*Locked {
	var data []*Locked
	if len(actionInfo) == 0 {
		data = append(data, &Locked{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*LockedTokenList{},
			BorrowTokenList:     []*LockedTokenList{},
			RewardTokenList:     []*LockedTokenList{},
			UnderLyingTokenList: []*LockedTokenList{},
			StrikeTokenList:     []*LockedTokenList{},
			CollateralTokenList: []*LockedTokenList{},
		})
	}
	return data
}

// GetLeaveragedFarmingEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetLeaveragedFarmingEmptyActionInfo(actionInfo []*LeaveragedFarming, actionType string, actionable bool) []*LeaveragedFarming {
	var data []*LeaveragedFarming
	if len(actionInfo) == 0 {
		data = append(data, &LeaveragedFarming{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*LeaveragedFarmingTokenList{},
			BorrowTokenList:     []*LeaveragedFarmingTokenList{},
			RewardTokenList:     []*LeaveragedFarmingTokenList{},
			UnderLyingTokenList: []*LeaveragedFarmingTokenList{},
			StrikeTokenList:     []*LeaveragedFarmingTokenList{},
			CollateralTokenList: []*LeaveragedFarmingTokenList{},
		})
	}
	return data
}

// GetInsuranceSellerEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetInsuranceSellerEmptyActionInfo(actionInfo []*InsuranceSeller, actionType string, actionable bool) []*InsuranceSeller {
	var data []*InsuranceSeller
	if len(actionInfo) == 0 {
		data = append(data, &InsuranceSeller{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*InsuranceSellerTokenList{},
			BorrowTokenList:     []*InsuranceSellerTokenList{},
			RewardTokenList:     []*InsuranceSellerTokenList{},
			UnderLyingTokenList: []*InsuranceSellerTokenList{},
			StrikeTokenList:     []*InsuranceSellerTokenList{},
			CollateralTokenList: []*InsuranceSellerTokenList{},
		})
	}
	return data
}

// GetOptionsBuyerEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetOptionsBuyerEmptyActionInfo(actionInfo []*OptionsBuyer, actionType string, actionable bool) []*OptionsBuyer {
	var data []*OptionsBuyer
	if len(actionInfo) == 0 {
		data = append(data, &OptionsBuyer{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*OptionsBuyerTokenList{},
			BorrowTokenList:     []*OptionsBuyerTokenList{},
			RewardTokenList:     []*OptionsBuyerTokenList{},
			UnderLyingTokenList: []*OptionsBuyerTokenList{},
			StrikeTokenList:     []*OptionsBuyerTokenList{},
			CollateralTokenList: []*OptionsBuyerTokenList{},
		})
	}
	return data
}

// GetInsuranceBuyerEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetInsuranceBuyerEmptyActionInfo(actionInfo []*InsuranceBuyer, actionType string, actionable bool) []*InsuranceBuyer {
	var data []*InsuranceBuyer
	if len(actionInfo) == 0 {
		data = append(data, &InsuranceBuyer{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*InsuranceBuyerTokenList{},
			BorrowTokenList:     []*InsuranceBuyerTokenList{},
			RewardTokenList:     []*InsuranceBuyerTokenList{},
			UnderLyingTokenList: []*InsuranceBuyerTokenList{},
			StrikeTokenList:     []*InsuranceBuyerTokenList{},
			CollateralTokenList: []*InsuranceBuyerTokenList{},
		})
	}
	return data
}

// GetPerpetualsEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetPerpetualsEmptyActionInfo(actionInfo []*Perpetuals, actionType string, actionable bool) []*Perpetuals {
	var data []*Perpetuals
	if len(actionInfo) == 0 {
		data = append(data, &Perpetuals{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*PerpetualsTokenList{},
			BorrowTokenList:     []*PerpetualsTokenList{},
			RewardTokenList:     []*PerpetualsTokenList{},
			UnderLyingTokenList: []*PerpetualsTokenList{},
			StrikeTokenList:     []*PerpetualsTokenList{},
			CollateralTokenList: []*PerpetualsTokenList{},
		})
	}
	return data
}

// GetOthersEmptyActionInfo get empty action info for portfolio
func (d *DebankAPIService) GetOthersEmptyActionInfo(actionInfo []*Others, actionType string, actionable bool) []*Others {
	var data []*Others
	if len(actionInfo) == 0 {
		data = append(data, &Others{
			IsActionable:        actionable,
			Type:                actionType,
			AssetValue:          "0",
			NetValue:            "0",
			TokensSupplied:      []*OthersTokenList{},
			BorrowTokenList:     []*OthersTokenList{},
			RewardTokenList:     []*OthersTokenList{},
			UnderLyingTokenList: []*OthersTokenList{},
			StrikeTokenList:     []*OthersTokenList{},
			CollateralTokenList: []*OthersTokenList{},
		})
	}
	return data
}
