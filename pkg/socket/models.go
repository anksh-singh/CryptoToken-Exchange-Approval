package socket

type ChainsSocket struct {
	Result []struct {
		ChainID  int64 `json:"chainId"`
		Currency struct {
			Decimals                int64  `json:"decimals"`
			Icon                    string `json:"icon"`
			MinNativeCurrencyForGas string `json:"minNativeCurrencyForGas"`
			Name                    string `json:"name"`
			Symbol                  string `json:"symbol"`
		} `json:"currency"`
		Explorers        []string `json:"explorers"`
		Icon             string   `json:"icon"`
		IsL1             bool     `json:"isL1"`
		Name             string   `json:"name"`
		ReceivingEnabled bool     `json:"receivingEnabled"`
		Refuel           struct {
			ReceivingEnabled bool `json:"receivingEnabled"`
			SendingEnabled   bool `json:"sendingEnabled"`
		} `json:"refuel"`
		Rpcs           []string `json:"rpcs"`
		SendingEnabled bool     `json:"sendingEnabled"`
	} `json:"result"`
	Success bool `json:"success"`
}

type ChainTokensSocket struct {
	Result []struct {
		Address         string `json:"address"`
		ChainAgnosticID string `json:"chainAgnosticId"`
		ChainID         int64  `json:"chainId"`
		Decimals        int64  `json:"decimals"`
		Icon            string `json:"icon"`
		LogoURI         string `json:"logoURI"`
		Name            string `json:"name"`
		Symbol          string `json:"symbol"`
	} `json:"result"`
	Success bool `json:"success"`
}

type QuoteSocketResponse struct {
	Result struct {
		BridgeRouteErrors struct {
			Across struct {
				Status string `json:"status"`
			} `json:"across"`
			Anyswap_router_v4 struct {
				AvailableLiquidity string `json:"availableLiquidity"`
				Status             string `json:"status"`
			} `json:"anyswap-router-v4"`
			Arbitrum_bridge struct {
				Status string `json:"status"`
			} `json:"arbitrum-bridge"`
			Hop struct {
				Status string `json:"status"`
			} `json:"hop"`
			Optimism_bridge struct {
				Status string `json:"status"`
			} `json:"optimism-bridge"`
			Polygon_bridge struct {
				Status string `json:"status"`
			} `json:"polygon-bridge"`
			Refuel_bridge struct {
				MaxAmount string `json:"maxAmount"`
				Status    string `json:"status"`
			} `json:"refuel-bridge"`
		} `json:"bridgeRouteErrors"`
		FromAsset struct {
			Address         string `json:"address"`
			ChainAgnosticID string `json:"chainAgnosticId"`
			ChainID         int64  `json:"chainId"`
			Decimals        int64  `json:"decimals"`
			Icon            string `json:"icon"`
			LogoURI         string `json:"logoURI"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
		} `json:"fromAsset"`
		FromChainID int64 `json:"fromChainId"`
		Routes      []struct {
			ChainGasBalances struct {
				One37 struct {
					HasGasBalance bool   `json:"hasGasBalance"`
					MinGasBalance string `json:"minGasBalance"`
				} `json:"137"`
				Five6 struct {
					HasGasBalance bool   `json:"hasGasBalance"`
					MinGasBalance string `json:"minGasBalance"`
				} `json:"56"`
			} `json:"chainGasBalances"`
			FromAmount         string `json:"fromAmount"`
			IsOnlySwapRoute    bool   `json:"isOnlySwapRoute"`
			MaxServiceTime     int64  `json:"maxServiceTime"`
			MinimumGasBalances struct {
				One37 string `json:"137"`
				Five6 string `json:"56"`
			} `json:"minimumGasBalances"`
			Recipient         string    `json:"recipient"`
			RouteID           string    `json:"routeId"`
			Sender            string    `json:"sender"`
			ServiceTime       int64     `json:"serviceTime"`
			ToAmount          string    `json:"toAmount"`
			TotalGasFeesInUsd float64   `json:"totalGasFeesInUsd"`
			TotalUserTx       int64     `json:"totalUserTx"`
			UsedBridgeNames   []string  `json:"usedBridgeNames"`
			UserTxs           []UserTxs `json:"userTxs"`
		} `json:"routes"`
		ToAsset struct {
			Address         string `json:"address"`
			ChainAgnosticID string `json:"chainAgnosticId"`
			ChainID         int64  `json:"chainId"`
			Decimals        int64  `json:"decimals"`
			Icon            string `json:"icon"`
			LogoURI         string `json:"logoURI"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
		} `json:"toAsset"`
		ToChainID int64 `json:"toChainId"`
	} `json:"result"`
	Success bool `json:"success"`
}

type RoutePayload struct {
	Route struct {
		FromAmount        string    `json:"fromAmount"`
		IsOnlySwapRoute   bool      `json:"isOnlySwapRoute"`
		MaxServiceTime    int64     `json:"maxServiceTime"`
		Recipient         string    `json:"recipient"`
		RouteID           string    `json:"routeId"`
		Sender            string    `json:"sender"`
		ServiceTime       int64     `json:"serviceTime"`
		ToAmount          string    `json:"toAmount"`
		TotalGasFeesInUsd float64   `json:"totalGasFeesInUsd"`
		TotalUserTx       int64     `json:"totalUserTx"`
		UsedBridgeNames   []string  `json:"usedBridgeNames"`
		UserTxs           []UserTxs `json:"userTxs"`
	} `json:"route"`
}

type Route struct {
	FromAmount        string    `json:"fromAmount"`
	IsOnlySwapRoute   bool      `json:"isOnlySwapRoute"`
	MaxServiceTime    int64     `json:"maxServiceTime"`
	Recipient         string    `json:"recipient"`
	RouteID           string    `json:"routeId"`
	Sender            string    `json:"sender"`
	ServiceTime       int64     `json:"serviceTime"`
	ToAmount          string    `json:"toAmount"`
	TotalGasFeesInUsd float64   `json:"totalGasFeesInUsd"`
	TotalUserTx       int64     `json:"totalUserTx"`
	UsedBridgeNames   []string  `json:"usedBridgeNames"`
	UserTxs           []UserTxs `json:"userTxs"`
}

type UserTxs struct {
	ApprovalData struct {
		AllowanceTarget       string `json:"allowanceTarget"`
		ApprovalTokenAddress  string `json:"approvalTokenAddress"`
		MinimumApprovalAmount string `json:"minimumApprovalAmount"`
		Owner                 string `json:"owner"`
	} `json:"approvalData"`
	//BridgeSlippage int64 `json:"bridgeSlippage"`
	ChainID int64 `json:"chainId"`
	GasFees struct {
		Asset struct {
			Address         string      `json:"address"`
			ChainAgnosticID interface{} `json:"chainAgnosticId"`
			ChainID         int64       `json:"chainId"`
			Decimals        int64       `json:"decimals"`
			Icon            string      `json:"icon"`
			LogoURI         string      `json:"logoURI"`
			Name            string      `json:"name"`
			Symbol          string      `json:"symbol"`
		} `json:"asset"`
		FeesInUsd float64 `json:"feesInUsd"`
		GasAmount string  `json:"gasAmount"`
		GasLimit  int64   `json:"gasLimit"`
	} `json:"gasFees"`
	MaxServiceTime int64   `json:"maxServiceTime"`
	Recipient      string  `json:"recipient"`
	RoutePath      string  `json:"routePath"`
	Sender         string  `json:"sender"`
	ServiceTime    int64   `json:"serviceTime"`
	StepCount      int64   `json:"stepCount"`
	Steps          []Steps `json:"steps"`
	ToAmount       string  `json:"toAmount"`
	MinAmountOut   string  `json:"minAmountOut"`
	ToAsset        struct {
		Address         string `json:"address"`
		ChainAgnosticID string `json:"chainAgnosticId"`
		ChainID         int64  `json:"chainId"`
		Decimals        int64  `json:"decimals"`
		Icon            string `json:"icon"`
		LogoURI         string `json:"logoURI"`
		Name            string `json:"name"`
		Symbol          string `json:"symbol"`
	} `json:"toAsset"`
	TxType      string `json:"txType"`
	UserTxIndex int64  `json:"userTxIndex"`
	UserTxType  string `json:"userTxType"`
}

type Steps struct {
	Type     string `json:"type"`
	Protocol struct {
		Name            string `json:"name"`
		DisplayName     string `json:"displayName"`
		Icon            string `json:"icon"`
		SecurityScore   int    `json:"securityScore"`
		RobustnessScore int    `json:"robustnessScore"`
	} `json:"protocol"`
	FromChainId  int     `json:"fromChainId"`
	SwapSlippage float64 `json:"swapSlippage"`
	ChainId      int     `json:"chainId"`
	FromAsset    struct {
		ChainId         int    `json:"chainId"`
		Address         string `json:"address"`
		Symbol          string `json:"symbol"`
		Name            string `json:"name"`
		Decimals        int    `json:"decimals"`
		Icon            string `json:"icon"`
		LogoURI         string `json:"logoURI"`
		ChainAgnosticId string `json:"chainAgnosticId"`
	} `json:"fromAsset"`
	ToChainId  int    `json:"toChainId"`
	FromAmount string `json:"fromAmount"`
	ToAsset    struct {
		ChainId         int    `json:"chainId"`
		Address         string `json:"address"`
		Symbol          string `json:"symbol"`
		Name            string `json:"name"`
		Decimals        int    `json:"decimals"`
		Icon            string `json:"icon"`
		LogoURI         string `json:"logoURI"`
		ChainAgnosticId string `json:"chainAgnosticId"`
	} `json:"toAsset"`
	ToAmount     string `json:"toAmount"`
	MinAmountOut string `json:"minAmountOut"`
	GasFees      struct {
		GasAmount string `json:"gasAmount"`
		GasLimit  int    `json:"gasLimit"`
		Asset     struct {
			ChainId         int         `json:"chainId"`
			Address         string      `json:"address"`
			Symbol          string      `json:"symbol"`
			Name            string      `json:"name"`
			Decimals        int         `json:"decimals"`
			Icon            string      `json:"icon"`
			LogoURI         string      `json:"logoURI"`
			ChainAgnosticId interface{} `json:"chainAgnosticId"`
		} `json:"asset"`
		FeesInUsd float64 `json:"feesInUsd"`
	} `json:"gasFees"`
	BridgeSlippage float64 `json:"bridgeSlippage"`
	ProtocolFees   struct {
		Asset struct {
			ChainId         int    `json:"chainId"`
			Address         string `json:"address"`
			Symbol          string `json:"symbol"`
			Name            string `json:"name"`
			Decimals        int    `json:"decimals"`
			Icon            string `json:"icon"`
			LogoURI         string `json:"logoURI"`
			ChainAgnosticId string `json:"chainAgnosticId"`
		} `json:"asset"`
		FeesInUsd int    `json:"feesInUsd"`
		Amount    string `json:"amount"`
	} `json:"protocolFees"`
	ServiceTime    int `json:"serviceTime"`
	MaxServiceTime int `json:"maxServiceTime"`
}

//type Steps struct {
//	//BridgeSlippage int64  `json:"bridgeSlippage"`
//	FromAmount string `json:"fromAmount"`
//	FromAsset  struct {
//		Address         string `json:"address"`
//		ChainAgnosticID string `json:"chainAgnosticId"`
//		ChainID         int64  `json:"chainId"`
//		Decimals        int64  `json:"decimals"`
//		Icon            string `json:"icon"`
//		LogoURI         string `json:"logoURI"`
//		Name            string `json:"name"`
//		Symbol          string `json:"symbol"`
//	} `json:"fromAsset"`
//	FromChainID int64 `json:"fromChainId"`
//	GasFees     struct {
//		Asset struct {
//			Address         string      `json:"address"`
//			ChainAgnosticID interface{} `json:"chainAgnosticId"`
//			ChainID         int64       `json:"chainId"`
//			Decimals        int64       `json:"decimals"`
//			Icon            string      `json:"icon"`
//			LogoURI         string      `json:"logoURI"`
//			Name            string      `json:"name"`
//			Symbol          string      `json:"symbol"`
//		} `json:"asset"`
//		FeesInUsd float64 `json:"feesInUsd"`
//		GasAmount string  `json:"gasAmount"`
//		GasLimit  int64   `json:"gasLimit"`
//	} `json:"gasFees"`
//	MaxServiceTime int64  `json:"maxServiceTime"`
//	MinAmountOut   string `json:"minAmountOut"`
//	Protocol       struct {
//		DisplayName     string `json:"displayName"`
//		Icon            string `json:"icon"`
//		Name            string `json:"name"`
//		RobustnessScore int64  `json:"robustnessScore"`
//		SecurityScore   int64  `json:"securityScore"`
//	} `json:"protocol"`
//	ProtocolFees struct {
//		//Amount string `json:"amount"`
//		Asset struct {
//			Address         string `json:"address"`
//			ChainAgnosticID string `json:"chainAgnosticId"`
//			ChainID         int64  `json:"chainId"`
//			Decimals        int64  `json:"decimals"`
//			Icon            string `json:"icon"`
//			LogoURI         string `json:"logoURI"`
//			Name            string `json:"name"`
//			Symbol          string `json:"symbol"`
//		} `json:"asset"`
//		FeesInUsd int64 `json:"feesInUsd"`
//	} `json:"protocolFees"`
//	ServiceTime int64  `json:"serviceTime"`
//	ToAmount    string `json:"toAmount"`
//	ToAsset     struct {
//		Address         string `json:"address"`
//		ChainAgnosticID string `json:"chainAgnosticId"`
//		ChainID         int64  `json:"chainId"`
//		Decimals        int64  `json:"decimals"`
//		Icon            string `json:"icon"`
//		LogoURI         string `json:"logoURI"`
//		Name            string `json:"name"`
//		Symbol          string `json:"symbol"`
//	} `json:"toAsset"`
//	ToChainID int64  `json:"toChainId"`
//	Type      string `json:"type"`
//}

type RouteTransactionResponseSocket struct {
	Result struct {
		ApprovalData struct {
			AllowanceTarget       string `json:"allowanceTarget"`
			ApprovalTokenAddress  string `json:"approvalTokenAddress"`
			MinimumApprovalAmount string `json:"minimumApprovalAmount"`
			Owner                 string `json:"owner"`
		} `json:"approvalData"`
		ChainID     int64  `json:"chainId"`
		TxData      string `json:"txData"`
		TxTarget    string `json:"txTarget"`
		TxType      string `json:"txType"`
		UserTxIndex int64  `json:"userTxIndex"`
		UserTxType  string `json:"userTxType"`
		Value       string `json:"value"`
	} `json:"result"`
	Success bool `json:"success"`
}

type CheckAllowanceResponse struct {
	Result struct {
		TokenAddress string `json:"tokenAddress"`
		Value        string `json:"value"`
	} `json:"result"`
	Success bool `json:"success"`
}

type ApprovalTransactionResponse struct {
	Result struct {
		Data string `json:"data"`
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"result"`
	Success bool `json:"success"`
}

type TransactionStatusSocket struct {
	Result struct {
		DestinationTransactionHash string `json:"destinationTransactionHash"`
		DestinationTxStatus        string `json:"destinationTxStatus"`
		FromChainID                int64  `json:"fromChainId"`
		SourceTx                   string `json:"sourceTx"`
		SourceTxStatus             string `json:"sourceTxStatus"`
		ToChainID                  int64  `json:"toChainId"`
	} `json:"result"`
	Success bool `json:"success"`
}

type ErrorSocket struct {
	Success bool `json:"success"`
	Error   struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}
