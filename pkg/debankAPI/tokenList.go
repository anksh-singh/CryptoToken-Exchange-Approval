package debankAPI

import "fmt"

func (d *DebankAPIService) LiquidityPoolTokenList(tokenListInfo TokenList) []*LiquidityPoolTokenList {
	var tokenList []*LiquidityPoolTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*LiquidityPoolTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LendingTokenList(tokenListInfo TokenList) []*LendingTokenList {
	var tokenList []*LendingTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*LendingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) YieldTokenList(tokenListInfo TokenList) []*YieldTokenList {
	var tokenList []*YieldTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*YieldTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) StakedTokenList(tokenListInfo TokenList) []*StakedTokenList {
	var tokenList []*StakedTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*StakedTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) FarmTokenList(tokenListInfo TokenList) []*FarmTokenList {
	var tokenList []*FarmTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*FarmTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) DepositTokenDetail(item TokenDetail) []*DepositTokenList {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),

		})

	if len(tokenList) == 0 {
		return []*DepositTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) DepositTokenList(tokenListInfo TokenList) []*DepositTokenList {
	var tokenList []*DepositTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*DepositTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LockedTokenList(tokenListInfo TokenList) []*LockedTokenList {
	var tokenList []*LockedTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*LockedTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) LeaveragedFarmingTokenList(tokenListInfo TokenList) []*LeaveragedFarmingTokenList {
	var tokenList []*LeaveragedFarmingTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*LeaveragedFarmingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) VestingTokenList(tokenListInfo TokenList) []*VestingTokenList {
	var tokenList []*VestingTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*VestingTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) RewardsTokenList(tokenListInfo TokenList) []*RewardsTokenList {
	var tokenList []*RewardsTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
		})
	}
	if len(tokenList) == 0 {
		return []*RewardsTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InvestmentTokenDetail(item TokenDetail) []*InvestmentTokenList {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),

		})

	if len(tokenList) == 0 {
		return []*InvestmentTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InvestmentTokenList(tokenListInfo TokenList) []*InvestmentTokenList {
	var tokenList []*InvestmentTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*InvestmentTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OptionsSellerTokenList(tokenListInfo TokenList) []*OptionsSellerTokenList {
	var tokenList []*OptionsSellerTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*OptionsSellerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OptionsBuyerTokenList(tokenListInfo TokenList) []*OptionsBuyerTokenList {
	var tokenList []*OptionsBuyerTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*OptionsBuyerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InsuranceSellerTokenList(tokenListInfo TokenList) []*InsuranceSellerTokenList {
	var tokenList []*InsuranceSellerTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*InsuranceSellerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) InsuranceBuyerTokenList(tokenListInfo TokenList) []*InsuranceBuyerTokenList {
	var tokenList []*InsuranceBuyerTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*InsuranceBuyerTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) PerpetualsTokenList(tokenListInfo TokenList) []*PerpetualsTokenList {
	var tokenList []*PerpetualsTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*PerpetualsTokenList{}
	}
	return tokenList
}

func (d *DebankAPIService) OthersTokenList(tokenListInfo TokenList) []*OthersTokenList {
	var tokenList []*OthersTokenList
	var tokenAddress string
	for _, item := range tokenListInfo {
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
			Balance: fmt.Sprintf("%f", item.Amount),
			QuoteRate: fmt.Sprintf("%f", item.Price),
			QuotePrice: fmt.Sprintf("%f", quotePrice),
			
		})
	}
	if len(tokenList) == 0 {
		return []*OthersTokenList{}
	}
	return tokenList
}

