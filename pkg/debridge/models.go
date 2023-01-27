package debridge

type Chains struct {
	ID               string `json:"id"`
	CommunityID      int    `json:"community_id"`
	Name             string `json:"name"`
	NativeTokenID    string `json:"native_token_id"`
	LogoURL          string `json:"logo_url"`
	WrappedTokenID   string `json:"wrapped_token_id"`
	IsSupportPreExec bool   `json:"is_support_pre_exec"`
}

type Quote struct {
	Estimation struct {
		SrcChainTokenIn struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			Amount   string `json:"amount"`
		} `json:"srcChainTokenIn"`
		SrcChainTokenOut struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"srcChainTokenOut"`
		DstChainTokenIn struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"dstChainTokenIn"`
		DstChainTokenOut struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"dstChainTokenOut"`
		ExecutionFee struct {
			Token struct {
				Address  string `json:"address"`
				Name     string `json:"name"`
				Symbol   string `json:"symbol"`
				Decimals int    `json:"decimals"`
			} `json:"token"`
			RecommendedAmount string `json:"recommendedAmount"`
			ActualAmount      string `json:"actualAmount"`
		} `json:"executionFee"`
	} `json:"estimation"`
	Tx struct {
		AllowanceTarget string `json:"allowanceTarget"`
	} `json:"tx"`
}

type Transaction struct {
	Estimation struct {
		SrcChainTokenIn struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			Amount   string `json:"amount"`
		} `json:"srcChainTokenIn"`
		SrcChainTokenOut struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"srcChainTokenOut"`
		DstChainTokenIn struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"dstChainTokenIn"`
		DstChainTokenOut struct {
			Address   string `json:"address"`
			Name      string `json:"name"`
			Symbol    string `json:"symbol"`
			Decimals  int    `json:"decimals"`
			Amount    string `json:"amount"`
			MinAmount string `json:"minAmount"`
		} `json:"dstChainTokenOut"`
		ExecutionFee struct {
			Token struct {
				Address  string `json:"address"`
				Name     string `json:"name"`
				Symbol   string `json:"symbol"`
				Decimals int    `json:"decimals"`
			} `json:"token"`
			RecommendedAmount string `json:"recommendedAmount"`
			ActualAmount      string `json:"actualAmount"`
		} `json:"executionFee"`
	} `json:"estimation"`
	Tx struct {
		To    string `json:"to"`
		Data  string `json:"data"`
		Value string `json:"value"`
	} `json:"tx"`
}

type TransactionStatus struct {
	Send struct {
		ID                 string  `json:"id"`
		Type               int     `json:"type"`
		EventOriginChainID int     `json:"eventOriginChainId"`
		Amount             string  `json:"amount"`
		Receiver           string  `json:"receiver"`
		DebridgeID         string  `json:"debridgeId"`
		ChainToID          int     `json:"chainToId"`
		TokenAddress       string  `json:"tokenAddress"`
		TokenDecimals      int     `json:"tokenDecimals"`
		AmountDecimals     int     `json:"amountDecimals"`
		TokenName          string  `json:"tokenName"`
		TokenSymbol        string  `json:"tokenSymbol"`
		SubmissionApproved bool    `json:"submissionApproved"`
		IsExecuted         bool    `json:"isExecuted"`
		Nonce              int     `json:"nonce"`
		BlockNumber        int     `json:"blockNumber"`
		BlockTimeStamp     int     `json:"blockTimeStamp"`
		TrackedTimeStamp   int     `json:"trackedTimeStamp"`
		TransactionHash    string  `json:"transactionHash"`
		SubmissionID       string  `json:"submissionId"`
		ReferralCode       int     `json:"referralCode"`
		ExecutionFee       string  `json:"executionFee"`
		Flags              string  `json:"flags"`
		FallbackAddress    string  `json:"fallbackAddress"`
		Data               string  `json:"data"`
		RawAutoparams      string  `json:"rawAutoparams"`
		NativeSender       string  `json:"nativeSender"`
		ConfirmationsCount int     `json:"confirmationsCount"`
		ReceivedAmount     string  `json:"receivedAmount"`
		FixFee             string  `json:"fixFee"`
		TransferFee        string  `json:"transferFee"`
		UseAssetFee        bool    `json:"useAssetFee"`
		LogoURI            string  `json:"logoURI"`
		Sender             string  `json:"sender"`
		AmountUsd          float64 `json:"amountUsd"`
		ProtocolFeeUsd     float64 `json:"protocolFeeUsd"`
		ExecutionFeeUsd    float64 `json:"executionFeeUsd"`
	} `json:"send"`
	Claim struct {
		ID                 string `json:"id"`
		EventOriginChainID int    `json:"eventOriginChainId"`
		Amount             string `json:"amount"`
		ChainFromID        int    `json:"chainFromId"`
		Receiver           string `json:"receiver"`
		Nonce              int    `json:"nonce"`
		DebridgeID         string `json:"debridgeId"`
		TokenAddress       string `json:"tokenAddress"`
		TokenDecimals      int    `json:"tokenDecimals"`
		TokenName          string `json:"tokenName"`
		TokenSymbol        string `json:"tokenSymbol"`
		BlockNumber        int    `json:"blockNumber"`
		BlockTimeStamp     int    `json:"blockTimeStamp"`
		TrackedTimeStamp   int    `json:"trackedTimeStamp"`
		TransactionHash    string `json:"transactionHash"`
		SubmissionID       string `json:"submissionId"`
		ExecutionFee       string `json:"executionFee"`
		Flags              string `json:"flags"`
		FallbackAddress    string `json:"fallbackAddress"`
		Data               string `json:"data"`
		Type               int    `json:"type"`
		NativeSender       string `json:"nativeSender"`
		RawAutoparams      string `json:"rawAutoparams"`
	} `json:"claim"`
}

type CoinGeckoPrice struct {
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
}

type NativeTokenInfo struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	ChainID  string `json:"chainId"`
	Decimals int64  `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}
