package debankAPI

import (
	"strconv"
)

// LendingInfo get action info
func (d *DebankAPIService) LendingInfo(item PortfolioItem, actionable bool) ([]*Lending, string) {
	var data []*Lending
	var netUsdValue string
	data = append(data, &Lending{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.LendingTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.LendingTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.LendingTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.LendingTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.LendingTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.LendingTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// LiquidityPoolInfo get action info
func (d *DebankAPIService) LiquidityPoolInfo(item PortfolioItem, actionable bool) ([]*LiquidityPool, string) {
	var data []*LiquidityPool
	var netUsdValue string
	data = append(data, &LiquidityPool{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.LiquidityPoolTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.LiquidityPoolTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.LiquidityPoolTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.LiquidityPoolTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.LiquidityPoolTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.LiquidityPoolTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// FarmingInfo get action info
func (d *DebankAPIService) FarmingInfo(item PortfolioItem, actionable bool) ([]*Farm, string) {
	var data []*Farm
	var netUsdValue string
	data = append(data, &Farm{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.FarmTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.FarmTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.FarmTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.FarmTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.FarmTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.FarmTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// DepositInfo get action info
func (d *DebankAPIService) DepositInfo(item PortfolioItem, actionable bool) ([]*Deposit, string) {
	var data []*Deposit
	var netUsdValue string
	data = append(data, &Deposit{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.DepositTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.DepositTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.DepositTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.DepositTokenDetail(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.DepositTokenDetail(item.Detail.StrikeToken),
		CollateralTokenList: d.DepositTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// InvestmentInfo get action info
func (d *DebankAPIService) InvestmentInfo(item PortfolioItem, actionable bool) ([]*Investment, string) {
	var data []*Investment
	var netUsdValue string
	data = append(data, &Investment{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.InvestmentTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.InvestmentTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.InvestmentTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.InvestmentTokenDetail(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.InvestmentTokenDetail(item.Detail.StrikeToken),
		CollateralTokenList: d.InvestmentTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// OptionSellerInfo get action info
func (d *DebankAPIService) OptionSellerInfo(item PortfolioItem, actionable bool) ([]*OptionsSeller, string) {
	var data []*OptionsSeller
	var netUsdValue string
	data = append(data, &OptionsSeller{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.OptionsSellerTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.OptionsSellerTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.OptionsSellerTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.OptionsSellerTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.OptionsSellerTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.OptionsSellerTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// RewardsInfo get action info
func (d *DebankAPIService) RewardsInfo(item PortfolioItem, actionable bool) ([]*Rewards, string) {
	var data []*Rewards
	var netUsdValue string
	data = append(data, &Rewards{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.RewardsTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.RewardsTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.RewardsTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.RewardsTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.RewardsTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.RewardsTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// VestingInfo get action info
func (d *DebankAPIService) VestingInfo(item PortfolioItem, actionable bool) ([]*Vesting, string) {
	var data []*Vesting
	var netUsdValue string
	data = append(data, &Vesting{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.VestingTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.VestingTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.VestingTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.VestingTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.VestingTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.VestingTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// YieldInfo get action info
func (d *DebankAPIService) YieldInfo(item PortfolioItem, actionable bool) ([]*Yield, string) {
	var data []*Yield
	var netUsdValue string
	data = append(data, &Yield{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.YieldTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.YieldTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.YieldTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.YieldTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.YieldTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.YieldTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// StakedInfo get action info
func (d *DebankAPIService) StakedInfo(item PortfolioItem, actionable bool) ([]*Staked, string) {
	var data []*Staked
	var netUsdValue string
	data = append(data, &Staked{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.StakedTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.StakedTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.StakedTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.StakedTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.StakedTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.StakedTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// LockedInfo get action info
func (d *DebankAPIService) LockedInfo(item PortfolioItem, actionable bool) ([]*Locked, string) {
	var data []*Locked
	var netUsdValue string
	data = append(data, &Locked{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.LockedTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.LockedTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.LockedTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.LockedTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.LockedTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.LockedTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// LeaveragedFarmingInfo get action info
func (d *DebankAPIService) LeaveragedFarmingInfo(item PortfolioItem, actionable bool) ([]*LeaveragedFarming, string) {
	var data []*LeaveragedFarming
	var netUsdValue string
	data = append(data, &LeaveragedFarming{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.LeaveragedFarmingTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.LeaveragedFarmingTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.LeaveragedFarmingTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.LeaveragedFarmingTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.LeaveragedFarmingTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.LeaveragedFarmingTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// OptionBuyerInfo get action info
func (d *DebankAPIService) OptionBuyerInfo(item PortfolioItem, actionable bool) ([]*OptionsBuyer, string) {
	var data []*OptionsBuyer
	var netUsdValue string
	data = append(data, &OptionsBuyer{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.OptionsBuyerTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.OptionsBuyerTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.OptionsBuyerTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.OptionsBuyerTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.OptionsBuyerTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.OptionsBuyerTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// InsuranceSellerInfo get action info
func (d *DebankAPIService) InsuranceSellerInfo(item PortfolioItem, actionable bool) ([]*InsuranceSeller, string) {
	var data []*InsuranceSeller
	var netUsdValue string
	data = append(data, &InsuranceSeller{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.InsuranceSellerTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.InsuranceSellerTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.InsuranceSellerTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.InsuranceSellerTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.InsuranceSellerTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.InsuranceSellerTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// InsuranceBuyerInfo get action info
func (d *DebankAPIService) InsuranceBuyerInfo(item PortfolioItem, actionable bool) ([]*InsuranceBuyer, string) {
	var data []*InsuranceBuyer
	var netUsdValue string
	data = append(data, &InsuranceBuyer{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.InsuranceBuyerTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.InsuranceBuyerTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.InsuranceBuyerTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.InsuranceBuyerTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.InsuranceBuyerTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.InsuranceBuyerTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}

// PerpetualsInfo get action info
func (d *DebankAPIService) PerpetualsInfo(item PortfolioItem, actionable bool) ([]*Perpetuals, string) {
	var data []*Perpetuals
	var netUsdValue string
	data = append(data, &Perpetuals{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.PerpetualsTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.PerpetualsTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.PerpetualsTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.PerpetualsTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.PerpetualsTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.PerpetualsTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		assetValueFloat := d.helper.ConvertStringToFloat64(value.AssetValue)
		netUsdValue = strconv.FormatFloat(d.helper.ConvertStringToFloat64(netUsdValue)+assetValueFloat, 'f', -1, 64)
	}
	return data, netUsdValue
}

// OthersInfo get action info
func (d *DebankAPIService) OthersInfo(item PortfolioItem, actionable bool) ([]*Others, string) {
	var data []*Others
	var netUsdValue string
	data = append(data, &Others{
		//IsActionable: 	 supportedActions[responseItem.ID],
		IsActionable: actionable,
		Type:         item.Name,
		AssetValue:   strconv.FormatFloat(item.Stats.AssetUsdValue, 'f', -1, 64),
		//DebtValue:       item.Stats.DebtUsdValue,
		NetValue:            strconv.FormatFloat(item.Stats.NetUsdValue, 'f', -1, 64),
		TokensSupplied:      d.OthersTokenList(item.Detail.SupplyTokenList),
		BorrowTokenList:     d.OthersTokenList(item.Detail.BorrowTokenList),
		RewardTokenList:     d.OthersTokenList(item.Detail.RewardTokenList),
		UnderLyingTokenList: d.OthersTokenDetails(item.Detail.UnderlyingToken),
		StrikeTokenList:     d.OthersTokenDetails(item.Detail.StrikeToken),
		CollateralTokenList: d.OthersTokenList(item.Detail.CollateralTokenList),
		//HealthRate:      item.Detail.HealthRate,
	})
	for _, value := range data {
		netUsdValue = netUsdValue + value.NetValue
	}
	return data, netUsdValue
}
