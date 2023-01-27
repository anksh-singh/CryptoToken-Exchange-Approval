package models



type OpportunityData struct {
	Apr                         string `json:"apr"`
	Chain                       string      `json:"chain"`
	Logo                        string      `json:"logo"`
	StakeTokenName              string      `json:"stakeTokenName,omitempty"`
	ReceiptTokenName            string      `json:"receiptTokenName,omitempty"`
	ContractDecimals            string      `json:"contractDecimals,omitempty"`
	StakeTokenLogoUrl           string      `json:"stakeTokenLogoUrl,omitempty"`
	StakeTokenContractAddress   string      `json:"stakeTokenContractAddress,omitempty"`
	ReceiptTokenLogoUrl         string      `json:"receiptTokenLogoUrl,omitempty"`
	ReceiptTokenContractAddres  string      `json:"receiptTokenContractAddres,omitempty"`
	ReceiptTokenContractAddress string      `json:"receiptTokenContractAddress,omitempty"`
	StakeToReceiptExchangeRate  float64     `json:"stakeToReceiptExchangeRate,omitempty"`
	ReceiptToStakeExchangeRate  float64     `json:"receiptToStakeExchangeRate,omitempty"`
	QuoteRate                   float64     `json:"quoteRate,omitempty"`
	ReceiptQuoteRate            float64     `json:"receiptQuoteRate,omitempty"`
	StakingType                 string      `json:"stakingType"`
	ProtocolName                string      `json:"protocolName"`
	CoolDownPeriod              string      `json:"coolDownPeriod"`
	MinLockup                   string      `json:"minLockup"`
	RewardSchedule              string      `json:"rewardSchedule"`
	TokenName                   string      `json:"tokenName,omitempty"`
}


type Opportunities struct {
	Current []OpportunityData `json:"current"`
	Others  []OpportunityData `json:"others"`
}