package core

import "time"

type SendTxRpcRes struct {
	Code      string `json:"code"`
	Data      string `json:"data"`
	Log       string `json:"log"`
	Codespace string `json:"codespace"`
	Hash      string `json:"hash"`
}

type RpcErrorRes struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}
type CosmosAssets struct {
	ChainName string `json:"chain_name"`
	Assets    []struct {
		Description string `json:"description"`
		DenomUnits  []struct {
			Denom    string `json:"denom"`
			Exponent int    `json:"exponent"`
		} `json:"denom_units"`
		Base     string `json:"base"`
		Name     string `json:"name"`
		Display  string `json:"display"`
		Symbol   string `json:"symbol"`
		LogoURIs struct {
			Svg string `json:"svg"`
			Png string `json:"png"`
		} `json:"logo_URIs"`
		CoingeckoID string `json:"coingecko_id"`
	} `json:"assets"`
}

type CosmostationAssets []struct {
	Denom       string `json:"denom"`
	Type        string `json:"type"`
	BaseDenom   string `json:"base_denom"`
	BaseType    string `json:"base_type"`
	DpDenom     string `json:"dp_denom"`
	OriginChain string `json:"origin_chain"`
	Decimal     int    `json:"decimal"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image"`
	Path        string `json:"path,omitempty"`
	Channel     string `json:"channel,omitempty"`
	Contract    string `json:"contract,omitempty"`
}

type CosmosBalanceRes struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type CosmosListTransactionRes struct {
	Header struct {
		ID        int       `json:"id"`
		ChainID   string    `json:"chain_id"`
		BlockID   int       `json:"block_id"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"header"`
	Data struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		Codespace string `json:"codespace"`
		Code      int    `json:"code"`
		Data      string `json:"data"`
		RawLog    string `json:"raw_log"`
		Logs      []struct {
			MsgIndex int    `json:"msg_index"`
			Log      string `json:"log"`
			Events   []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"logs"`
		Info      string `json:"info"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Tx        struct {
			Type string `json:"@type"`
			Body struct {
				Messages []struct {
					Type             string `json:"@type,omitempty"`
					FromAddress      string `json:"from_address,omitempty"`
					ToAddress        string `json:"to_address,omitempty"`
					DelegatorAddress string `json:"delegator_address,omitempty"`
					ValidatorAddress string `json:"validator_address,omitempty"`
					ProposalId       string `json:"proposal_id,omitempty"`
					Voter            string `json:"voter,omitempty"`
					Sender           string `json:"sender,omitempty"`
					Receiver         string `json:"receiver,omitempty"`
					Token            struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"token,omitempty"`
					Amount []struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"amount,omitempty"`
					DelegateAmount struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"amount,omitempty"` //amount has different structure based on tx type
				} `json:"messages"`
				Memo                        string        `json:"memo"`
				TimeoutHeight               string        `json:"timeout_height"`
				ExtensionOptions            []interface{} `json:"extension_options"`
				NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						Type string `json:"@type"`
						Key  string `json:"key"`
					} `json:"public_key"`
					ModeInfo struct {
						Single struct {
							Mode string `json:"mode"`
						} `json:"single"`
					} `json:"mode_info"`
					Sequence string `json:"sequence"`
				} `json:"signer_infos"`
				Fee struct {
					Amount []struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"amount"`
					GasLimit string `json:"gas_limit"`
					Payer    string `json:"payer"`
					Granter  string `json:"granter"`
				} `json:"fee"`
			} `json:"auth_info"`
			Signatures []string `json:"signatures"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"data,omitempty"`
}

type IBCData struct {
	Chain  string `json:"chain"`
	Origin struct {
		Denom interface{} `json:"denom"`
		Chain interface{} `json:"chain"`
	} `json:"origin"`
}

//type IBCTokenInfo struct {
//	IbcTokens []IBCData `json:"ibc_tokens"`
//}

type DenomInfo struct {
	Chain       string `json:"chain"`
	Name        string `json:"name"`
	Denom       string `json:"denom"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	Description string `json:"description"`
	CoingeckoID string `json:"coingecko_id"`
	Logos       struct {
		Png string `json:"png"`
	} `json:"logos"`
}

type CosmosAccountInfo struct {
	Account struct {
		Type    string `json:"@type"`
		Address string `json:"address"`
		PubKey  struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"pub_key"`
		BaseAccount struct { //For different model
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"base_account"`
		AccountNumber string `json:"account_number"`
		Sequence      string `json:"sequence"`
	} `json:"account"`
}

type CosmosScan struct {
	Txs []struct {
		Body struct {
			Messages []struct {
				Type             string `json:"@type"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"messages"`
			Memo                        string        `json:"memo"`
			TimeoutHeight               string        `json:"timeout_height"`
			ExtensionOptions            []interface{} `json:"extension_options"`
			NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
		} `json:"body"`
		AuthInfo struct {
			SignerInfos []struct {
				PublicKey struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"public_key"`
				ModeInfo struct {
					Single struct {
						Mode string `json:"mode"`
					} `json:"single"`
				} `json:"mode_info"`
				Sequence string `json:"sequence"`
			} `json:"signer_infos"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				GasLimit string `json:"gas_limit"`
				Payer    string `json:"payer"`
				Granter  string `json:"granter"`
			} `json:"fee"`
		} `json:"auth_info"`
		Signatures []string `json:"signatures"`
	} `json:"txs"`
	TxResponses []struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		Codespace string `json:"codespace"`
		Code      int    `json:"code"`
		Data      string `json:"data"`
		RawLog    string `json:"raw_log"`
		Logs      []struct {
			MsgIndex int    `json:"msg_index"`
			Log      string `json:"log"`
			Events   []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"logs"`
		Info      string `json:"info"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Tx        struct {
			Type string `json:"@type"`
			Body struct {
				Messages []struct {
					Type             string `json:"@type"`
					ValidatorAddress string `json:"validator_address"`
				} `json:"messages"`
				Memo                        string        `json:"memo"`
				TimeoutHeight               string        `json:"timeout_height"`
				ExtensionOptions            []interface{} `json:"extension_options"`
				NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						Type string `json:"@type"`
						Key  string `json:"key"`
					} `json:"public_key"`
					ModeInfo struct {
						Single struct {
							Mode string `json:"mode"`
						} `json:"single"`
					} `json:"mode_info"`
					Sequence string `json:"sequence"`
				} `json:"signer_infos"`
				Fee struct {
					Amount []struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"amount"`
					GasLimit string `json:"gas_limit"`
					Payer    string `json:"payer"`
					Granter  string `json:"granter"`
				} `json:"fee"`
			} `json:"auth_info"`
			Signatures []string `json:"signatures"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"tx_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
type CosmosAmount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type CosmosTxReceipt struct {
	Tx struct {
		Body struct {
			Messages []struct {
				Type             string      `json:"@type"`
				FromAddress      string      `json:"from_address"`
				ToAddress        string      `json:"to_address"`
				DelegatorAddress string      `json:"delegator_address"`
				ValidatorAddress string      `json:"validator_address"`
				Amount           interface{} `json:"amount"`
			} `json:"messages"`
			Memo                        string        `json:"memo"`
			TimeoutHeight               string        `json:"timeout_height"`
			ExtensionOptions            []interface{} `json:"extension_options"`
			NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
		} `json:"body"`
		AuthInfo struct {
			SignerInfos []struct {
				PublicKey struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"public_key"`
				ModeInfo struct {
					Single struct {
						Mode string `json:"mode"`
					} `json:"single"`
				} `json:"mode_info"`
				Sequence string `json:"sequence"`
			} `json:"signer_infos"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				GasLimit string `json:"gas_limit"`
				Payer    string `json:"payer"`
				Granter  string `json:"granter"`
			} `json:"fee"`
		} `json:"auth_info"`
		Signatures []string `json:"signatures"`
	} `json:"tx"`
	TxResponse struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		Codespace string `json:"codespace"`
		Code      int    `json:"code"`
		Data      string `json:"data"`
		RawLog    string `json:"raw_log"`
		Logs      []struct {
			MsgIndex int    `json:"msg_index"`
			Log      string `json:"log"`
			Events   []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"logs"`
		Info      string `json:"info"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Tx        struct {
			Type string `json:"@type"`
			Body struct {
				Messages []struct {
					Type             string `json:"@type"`
					ValidatorAddress string `json:"validator_address"`
				} `json:"messages"`
				Memo                        string        `json:"memo"`
				TimeoutHeight               string        `json:"timeout_height"`
				ExtensionOptions            []interface{} `json:"extension_options"`
				NonCriticalExtensionOptions []interface{} `json:"non_critical_extension_options"`
			} `json:"body"`
			AuthInfo struct {
				SignerInfos []struct {
					PublicKey struct {
						Type string `json:"@type"`
						Key  string `json:"key"`
					} `json:"public_key"`
					ModeInfo struct {
						Single struct {
							Mode string `json:"mode"`
						} `json:"single"`
					} `json:"mode_info"`
					Sequence string `json:"sequence"`
				} `json:"signer_infos"`
				Fee struct {
					Amount []struct {
						Denom  string `json:"denom"`
						Amount string `json:"amount"`
					} `json:"amount"`
					GasLimit string `json:"gas_limit"`
					Payer    string `json:"payer"`
					Granter  string `json:"granter"`
				} `json:"fee"`
			} `json:"auth_info"`
			Signatures []string `json:"signatures"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"tx_response"`
}
type CosmosValidatorsLaunchPad struct {
	Height string `json:"height"`
	Result []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey string `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          int    `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"result"`
}
type CosmosValidators struct {
	Validators []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string `json:"unbonding_height"`
		UnbondingTime   string `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime string `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"validators"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
type CosmosPoolLaunchPad struct {
	Height string `json:"height"`
	Result struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"result"`
}
type CosmosPool struct {
	Pool struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"pool"`
}
type CosmosAnnualProvision struct {
	AnnualProvisions string `json:"annual_provisions"`
}
type CosmosDistributionParams struct {
	Params struct {
		CommunityTax        string `json:"community_tax"`
		BaseProposerReward  string `json:"base_proposer_reward"`
		BonusProposerReward string `json:"bonus_proposer_reward"`
		WithdrawAddrEnabled bool   `json:"withdraw_addr_enabled"`
	} `json:"params"`
}
type CosmosAnnualProvisionLaunchPad struct {
	AnnualProvisions string `json:"annual_provisions"`
}
type CosmosDistributionParamsLaunchPad struct {
	Height string `json:"height"`
	Result struct {
		CommunityTax        string `json:"community_tax"`
		BaseProposerReward  string `json:"base_proposer_reward"`
		BonusProposerReward string `json:"bonus_proposer_reward"`
		WithdrawAddrEnabled bool   `json:"withdraw_addr_enabled"`
	} `json:"result"`
}
type CosmosInflationLaunchPad struct {
	Inflation string `json:"inflation"`
}
type CosmosInflation struct {
	Inflation string `json:"inflation"`
}

type ValidatorImageURL struct {
	Status struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"status"`
	Them []struct {
		ID       string `json:"id"`
		Pictures struct {
			Primary struct {
				URL    string      `json:"url"`
				Source interface{} `json:"source"`
			} `json:"primary"`
		} `json:"pictures"`
	} `json:"them"`
}
type CosmosDelegationsLaunchPad struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Shares           string `json:"shares"`
		Balance          struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"result"`
}
type CosmosDelegations struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
type CosmosValidatorsForDelegatorLaunchPad struct {
	Height string `json:"height"`
	Result []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey string `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          int    `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string    `json:"unbonding_height"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime time.Time `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"result"`
}
type CosmosValidatorsForDelegator struct {
	Validators []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string `json:"unbonding_height"`
		UnbondingTime   string `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime string `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation string `json:"min_self_delegation"`
	} `json:"validators"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
type AutoGeneratedLaunchPad struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		InitialBalance   string `json:"initial_balance"`
		Balance          string `json:"balance"`
		CreationHeight   int    `json:"creation_height"`
		MinTime          int    `json:"min_time"`
	} `json:"result"`
}
type CosmosUnboundDelegations struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string `json:"creation_height"`
			CompletionTime string `json:"completion_time"`
			InitialBalance string `json:"initial_balance"`
			Balance        string `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
type CosmosRewardsLaunchPad struct {
	Height string `json:"height"`
	Result struct {
		Rewards []struct {
			ValidatorAddress string `json:"validator_address"`
			Reward           []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"reward"`
		} `json:"rewards"`
		Total []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total"`
	} `json:"result"`
}
type CosmosRewards struct {
	Rewards []struct {
		ValidatorAddress string `json:"validator_address"`
		Reward           []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"reward"`
	} `json:"rewards"`
	Total []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
}

type CosmosSendTxRequest struct {
	TxBytes string `json:"tx_bytes"`
	Mode    string `json:"mode"`
}

type CosmosSendTxRes struct {
	TxResponse struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		Codespace string `json:"codespace"`
		Code      int    `json:"code"`
		Data      string `json:"data"`
		RawLog    string `json:"raw_log"`
		Logs      []struct {
			MsgIndex int    `json:"msg_index"`
			Log      string `json:"log"`
			Events   []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"logs"`
		Info      string `json:"info"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Tx        struct {
			TypeURL string `json:"type_url"`
			Value   string `json:"value"`
		} `json:"tx"`
		Timestamp string `json:"timestamp"`
	} `json:"tx_response"`
}
type CosmosCDPParameterResp struct {
	Height string `json:"height"`
	Result struct {
		CollateralParams []struct {
			Denom            string `json:"denom"`
			Type             string `json:"type"`
			LiquidationRatio string `json:"liquidation_ratio"`
			DebtLimit        struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"debt_limit"`
			StabilityFee                     string `json:"stability_fee"`
			AuctionSize                      string `json:"auction_size"`
			LiquidationPenalty               string `json:"liquidation_penalty"`
			SpotMarketID                     string `json:"spot_market_id"`
			LiquidationMarketID              string `json:"liquidation_market_id"`
			KeeperRewardPercentage           string `json:"keeper_reward_percentage"`
			CheckCollateralizationIndexCount string `json:"check_collateralization_index_count"`
			ConversionFactor                 string `json:"conversion_factor"`
		} `json:"collateral_params"`
		DebtParam struct {
			Denom            string `json:"denom"`
			ReferenceAsset   string `json:"reference_asset"`
			ConversionFactor string `json:"conversion_factor"`
			DebtFloor        string `json:"debt_floor"`
		} `json:"debt_param"`
		GlobalDebtLimit struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"global_debt_limit"`
		SurplusAuctionThreshold string `json:"surplus_auction_threshold"`
		SurplusAuctionLot       string `json:"surplus_auction_lot"`
		DebtAuctionThreshold    string `json:"debt_auction_threshold"`
		DebtAuctionLot          string `json:"debt_auction_lot"`
	} `json:"result"`
}
type CosmosBankResponse struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}
type CosmosIncentivesResponse struct {
	Height string `json:"height"`
	Result []struct {
		BaseClaim struct {
			Owner  string `json:"owner"`
			Reward struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"reward"`
		} `json:"base_claim"`
		RewardIndexes []struct {
			CollateralType string `json:"collateral_type"`
			RewardFactor   string `json:"reward_factor"`
		} `json:"reward_indexes"`
	} `json:"result"`
}
type CosmosCDPMarketPriceInfo struct {
	Price struct {
		MarketID string `json:"market_id"`
		Price    string `json:"price"`
	} `json:"price"`
}
type CosmosCDPResponse struct {
	Cdp struct {
		ID         string `json:"id"`
		Owner      string `json:"owner"`
		Type       string `json:"type"`
		Collateral struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"collateral"`
		Principal struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"principal"`
		AccumulatedFees struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"accumulated_fees"`
		FeesUpdated     time.Time `json:"fees_updated"`
		InterestFactor  string    `json:"interest_factor"`
		CollateralValue struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"collateral_value"`
		CollateralizationRatio string `json:"collateralization_ratio"`
	} `json:"cdp"`
}

type SimulateTxBody struct {
	TxBytes string `json:"tx_bytes"`
}

type SimulateTxResponse struct {
	GasInfo struct {
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
	} `json:"gas_info"`
	Result struct {
		Data   string `json:"data"`
		Log    string `json:"log"`
		Events []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"result"`
}

type BluzelleBalanceRes struct {
	Height string `json:"height"`
	Result []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"result"`
}

type BluzelleAccountInfo struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
			PublicKey struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"public_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

type BluzelleSendTxResponse struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		MsgIndex int    `json:"msg_index"`
		Log      string `json:"log"`
		Events   []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
}

type CosmostationAssetInfo []struct {
	Denom       string `json:"denom"`
	Type        string `json:"type"`
	BaseDenom   string `json:"base_denom"`
	BaseType    string `json:"base_type"`
	DpDenom     string `json:"dp_denom"`
	OriginChain string `json:"origin_chain"`
	Decimal     int    `json:"decimal"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CoinGeckoID string `json:"coinGeckoId"`
}
type CosmosIncentiveRewardsInfo struct {
	Height string `json:"height"`
	Result []struct {
		HardClaims []struct {
			BaseClaim struct {
				Owner  string `json:"owner"`
				Reward []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"reward"`
			} `json:"base_claim"`
			SupplyRewardIndexes []struct {
				CollateralType string `json:"collateral_type"`
				RewardIndexes  []struct {
					CollateralType string `json:"collateral_type"`
					RewardFactor   string `json:"reward_factor"`
				} `json:"reward_indexes"`
			} `json:"supply_reward_indexes"`
			BorrowRewardIndexes []struct {
				CollateralType string `json:"collateral_type"`
				RewardIndexes  []struct {
					CollateralType string `json:"collateral_type"`
					RewardFactor   string `json:"reward_factor"`
				} `json:"reward_indexes"`
			} `json:"borrow_reward_indexes"`
			DelegatorRewardIndexes []struct {
				CollateralType string `json:"collateral_type"`
				RewardFactor   string `json:"reward_factor"`
			} `json:"delegator_reward_indexes"`
		} `json:"hard_claims"`
		UsdxMintingClaims []struct {
			BaseClaim struct {
				Owner  string `json:"owner"`
				Reward struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"reward"`
			} `json:"base_claim"`
			RewardIndexes []struct {
				CollateralType string `json:"collateral_type"`
				RewardFactor   string `json:"reward_factor"`
			} `json:"reward_indexes"`
		} `json:"usdx_minting_claims"`
	} `json:"result"`
}

type BlockHeight struct {
	Block struct {
		Header struct {
			ChainID string `json:"chain_id"`
			Height  string `json:"height"`
		} `json:"header"`
	} `json:"block"`
}

