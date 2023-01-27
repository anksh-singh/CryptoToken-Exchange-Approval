package debankAPI

import "fmt"

func (d *DebankAPIService) LiquidityPoolTokenDetails(item TokenDetail) []*LiquidityPoolTokenList {
	var tokenList []*LiquidityPoolTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &LiquidityPoolTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*LiquidityPoolTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LendingTokenDetails(item TokenDetail) []*LendingTokenList {
	var tokenList []*LendingTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &LendingTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*LendingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) YieldTokenDetails(item TokenDetail) []*YieldTokenList {
	var tokenList []*YieldTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &YieldTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*YieldTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) StakedTokenDetails(item TokenDetail) []*StakedTokenList {
	var tokenList []*StakedTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &StakedTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*StakedTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) FarmTokenDetails(item TokenDetail) []*FarmTokenList {
	var tokenList []*FarmTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &FarmTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*FarmTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) DepositTokenDetails(item TokenDetail) []*DepositTokenList {
	var tokenList []*DepositTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &DepositTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*DepositTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LockedTokenDetails(item TokenDetail) []*LockedTokenList {
	var tokenList []*LockedTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &LockedTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*LockedTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LeaveragedFarmingTokenDetails(item TokenDetail) []*LeaveragedFarmingTokenList {
	var tokenList []*LeaveragedFarmingTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &LeaveragedFarmingTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*LeaveragedFarmingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) VestingTokenDetails(item TokenDetail) []*VestingTokenList {
	var tokenList []*VestingTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &VestingTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*VestingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) RewardsTokenDetails(item TokenDetail) []*RewardsTokenList {
	var tokenList []*RewardsTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &RewardsTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*RewardsTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InvestmentTokenDetails(item TokenDetail) []*InvestmentTokenList {
	var tokenList []*InvestmentTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &InvestmentTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*InvestmentTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OptionsSellerTokenDetails(item TokenDetail) []*OptionsSellerTokenList {
	var tokenList []*OptionsSellerTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &OptionsSellerTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*OptionsSellerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OptionsBuyerTokenDetails(item TokenDetail) []*OptionsBuyerTokenList {
	var tokenList []*OptionsBuyerTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &OptionsBuyerTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*OptionsBuyerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InsuranceSellerTokenDetails(item TokenDetail) []*InsuranceSellerTokenList {
	var tokenList []*InsuranceSellerTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &InsuranceSellerTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*InsuranceSellerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InsuranceBuyerTokenDetails(item TokenDetail) []*InsuranceBuyerTokenList {
	var tokenList []*InsuranceBuyerTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &InsuranceBuyerTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*InsuranceBuyerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) PerpetualsTokenDetails(item TokenDetail) []*PerpetualsTokenList {
	var tokenList []*PerpetualsTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &PerpetualsTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*PerpetualsTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OthersTokenDetails(item TokenDetail) []*OthersTokenList {
	var tokenList []*OthersTokenList
	var tokenAddress string

	if item.ID == item.Chain {
		tokenAddress = d.GetNativeTokenAddress(item.Chain)
	} else {
		tokenAddress = item.ID
	}
	quotePrice := item.Amount * item.Price
	tokenList = append(tokenList, &OthersTokenList{
		TokenAddress:  tokenAddress,
		TokenName:     item.Name,
		TokenSymbol:   item.Symbol,
		TokenDecimals: int32(item.Decimals),
		LogoUrl:       item.LogoURL,
		Balance:       fmt.Sprintf("%f", item.Amount),
		QuoteRate:     fmt.Sprintf("%f", item.Price),
		QuotePrice:    fmt.Sprintf("%f", quotePrice),
	})
	if len(tokenList) == 0 {
		return []*OthersTokenList{}
	}
	return tokenList
}
