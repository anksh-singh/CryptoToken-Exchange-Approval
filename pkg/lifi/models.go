package lifi

type ChainsInfo struct {
	Chains []*Chains `json:"chains"`
}

type Chains struct {
	Key              string   `json:"key"`
	ChainType        string   `json:"chainType"`
	Name             string   `json:"name"`
	Coin             string   `json:"coin"`
	ID               int      `json:"id"`
	Mainnet          bool     `json:"mainnet"`
	LogoURI          string   `json:"logoURI"`
	TokenlistURL     string   `json:"tokenlistUrl"`
	MulticallAddress string   `json:"multicallAddress"`
	Metamask         MetaMask `json:"metamask"`
	FaucetUrls       []string `json:"faucetUrls,omitempty"`
}

type MetaMask struct {
	ChainID           string         `json:"chainId"`
	BlockExplorerUrls []string       `json:"blockExplorerUrls"`
	ChainName         string         `json:"chainName"`
	NativeCurrency    NativeCurrency `json:"nativeCurrency"`
	RPCUrls           []string       `json:"rpcUrls"`
}

type NativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

type ChainsToken struct {
	Connections []*Connections `json:"connections"`
}

type Connections struct {
	FromChainID int           `json:"fromChainId"`
	ToChainID   int           `json:"toChainId"`
	FromTokens  []*FromTokens `json:"fromTokens"`
	ToTokens    []*ToTokens   `json:"toTokens"`
}

type FromTokens struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
	ChainID  int    `json:"chainId"`
	CoinKey  string `json:"coinKey"`
	Name     string `json:"name"`
	LogoURI  string `json:"logoURI"`
	PriceUSD string `json:"priceUSD"`
}

type ToTokens struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
	ChainID  int    `json:"chainId"`
	CoinKey  string `json:"coinKey,omitempty"`
	Name     string `json:"name"`
	LogoURI  string `json:"logoURI,omitempty"`
	PriceUSD string `json:"priceUSD,omitempty"`
}

type ToolDetails struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	LogoURI string `json:"logoURI"`
}

type Action struct {
	FromChainID int        `json:"fromChainId"`
	FromAmount  string     `json:"fromAmount"`
	FromToken   FromTokens `json:"fromToken"`
	FromAddress string     `json:"fromAddress"`
	ToChainID   int        `json:"toChainId"`
	ToToken     ToTokens   `json:"toToken"`
	ToAddress   string     `json:"toAddress"`
	Slippage    float64    `json:"slippage"`
}

type GasCosts struct {
	Type      string     `json:"type"`
	Price     string     `json:"price"`
	Estimate  string     `json:"estimate"`
	Limit     string     `json:"limit"`
	Amount    string     `json:"amount"`
	AmountUSD string     `json:"amountUSD"`
	Token     FromTokens `json:"token"`
}

type Data struct {
	RelayFee      string `json:"relayFee"`
	CallData      string `json:"callData"`
	TransactionID string `json:"transactionId"`
}

type Estimate struct {
	FromAmount        string      `json:"fromAmount"`
	ToAmount          string      `json:"toAmount"`
	ToAmountMin       string      `json:"toAmountMin"`
	ApprovalAddress   string      `json:"approvalAddress"`
	GasCosts          []*GasCosts `json:"gasCosts"`
	ExecutionDuration float64     `json:"executionDuration"`
	FromAmountUSD     string      `json:"fromAmountUSD"`
	ToAmountUSD       string      `json:"toAmountUSD"`
	Data              Data        `json:"data"`
}

type TransactionRequest struct {
	Data     string `json:"data"`
	To       string `json:"to"`
	Value    string `json:"value"`
	From     string `json:"from"`
	ChainID  int    `json:"chainId"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
}

type Quote struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Tool          string      `json:"tool"`
	ToolDetails   ToolDetails `json:"toolDetails"`
	Action        Action      `json:"action"`
	Estimate      Estimate    `json:"estimate"`
	IncludedSteps []struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Tool        string `json:"tool"`
		ToolDetails struct {
			Key     string `json:"key"`
			Name    string `json:"name"`
			LogoURI string `json:"logoURI"`
		} `json:"toolDetails"`
		Action struct {
			FromChainID int `json:"fromChainId"`
			ToChainID   int `json:"toChainId"`
			FromToken   struct {
				Address  string `json:"address"`
				Symbol   string `json:"symbol"`
				Decimals int    `json:"decimals"`
				ChainID  int    `json:"chainId"`
				Name     string `json:"name"`
				CoinKey  string `json:"coinKey"`
				PriceUSD string `json:"priceUSD"`
				LogoURI  string `json:"logoURI"`
			} `json:"fromToken"`
			ToToken struct {
				Address  string `json:"address"`
				Decimals int    `json:"decimals"`
				Symbol   string `json:"symbol"`
				ChainID  int    `json:"chainId"`
				CoinKey  string `json:"coinKey"`
				Name     string `json:"name"`
				LogoURI  string `json:"logoURI"`
				PriceUSD string `json:"priceUSD"`
			} `json:"toToken"`
			FromAmount  string  `json:"fromAmount"`
			Slippage    float64 `json:"slippage"`
			FromAddress string  `json:"fromAddress"`
			ToAddress   string  `json:"toAddress"`
		} `json:"action"`
		Estimate struct {
			FromAmount        string  `json:"fromAmount"`
			ToAmount          string  `json:"toAmount"`
			ToAmountMin       string  `json:"toAmountMin"`
			ApprovalAddress   string  `json:"approvalAddress"`
			ExecutionDuration float64 `json:"executionDuration"`
			FeeCosts          []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Percentage  string `json:"percentage"`
				Token       struct {
					Address  string `json:"address"`
					Symbol   string `json:"symbol"`
					Decimals int    `json:"decimals"`
					ChainID  int    `json:"chainId"`
					Name     string `json:"name"`
					CoinKey  string `json:"coinKey"`
					PriceUSD string `json:"priceUSD"`
					LogoURI  string `json:"logoURI"`
				} `json:"token"`
				Amount    string `json:"amount"`
				AmountUSD string `json:"amountUSD"`
			} `json:"feeCosts"`
			GasCosts []struct {
				Type      string `json:"type"`
				Price     string `json:"price"`
				Estimate  string `json:"estimate"`
				Limit     string `json:"limit"`
				Amount    string `json:"amount"`
				AmountUSD string `json:"amountUSD"`
				Token     struct {
					Address  string `json:"address"`
					Symbol   string `json:"symbol"`
					Decimals int    `json:"decimals"`
					ChainID  int    `json:"chainId"`
					Name     string `json:"name"`
					CoinKey  string `json:"coinKey"`
					PriceUSD string `json:"priceUSD"`
					LogoURI  string `json:"logoURI"`
				} `json:"token"`
			} `json:"gasCosts"`
			Data struct {
				Bid struct {
					User                           string `json:"user"`
					Router                         string `json:"router"`
					Initiator                      string `json:"initiator"`
					SendingChainID                 int    `json:"sendingChainId"`
					SendingAssetID                 string `json:"sendingAssetId"`
					Amount                         string `json:"amount"`
					ReceivingChainID               int    `json:"receivingChainId"`
					ReceivingAssetID               string `json:"receivingAssetId"`
					AmountReceived                 string `json:"amountReceived"`
					ReceivingAddress               string `json:"receivingAddress"`
					TransactionID                  string `json:"transactionId"`
					Expiry                         int    `json:"expiry"`
					CallDataHash                   string `json:"callDataHash"`
					CallTo                         string `json:"callTo"`
					EncryptedCallData              string `json:"encryptedCallData"`
					SendingChainTxManagerAddress   string `json:"sendingChainTxManagerAddress"`
					ReceivingChainTxManagerAddress string `json:"receivingChainTxManagerAddress"`
					BidExpiry                      int    `json:"bidExpiry"`
				} `json:"bid"`
				BidSignature           string `json:"bidSignature"`
				GasFeeInReceivingToken string `json:"gasFeeInReceivingToken"`
				TotalFee               string `json:"totalFee"`
				MetaTxRelayerFee       string `json:"metaTxRelayerFee"`
				RouterFee              string `json:"routerFee"`
				ServerSign             bool   `json:"serverSign"`
			} `json:"data"`
			FromAmountUSD string `json:"fromAmountUSD"`
			ToAmountUSD   string `json:"toAmountUSD"`
		} `json:"estimate"`
	} `json:"includedSteps"`
	TransactionRequest TransactionRequest `json:"transactionRequest"`
}

type ExchangeTokens struct {
	Address  string `json:"address"`
	ChainID  int64  `json:"chainId"`
	CoinKey  string `json:"coinKey"`
	Decimals int64  `json:"decimals"`
	LogoURI  string `json:"logoURI"`
	Name     string `json:"name"`
	PriceUSD string `json:"priceUSD"`
	Symbol   string `json:"symbol"`
}
type ErrorResponse struct {
	Errors []struct {
		Action struct {
			FromAddress string `json:"fromAddress"`
			FromAmount  string `json:"fromAmount"`
			FromChainID int64  `json:"fromChainId"`
			FromToken   struct {
				Address  string `json:"address"`
				ChainID  int64  `json:"chainId"`
				CoinKey  string `json:"coinKey"`
				Decimals int64  `json:"decimals"`
				LogoURI  string `json:"logoURI"`
				Name     string `json:"name"`
				PriceUSD string `json:"priceUSD"`
				Symbol   string `json:"symbol"`
			} `json:"fromToken"`
			Slippage  float64 `json:"slippage"`
			ToAddress string  `json:"toAddress"`
			ToChainID int64   `json:"toChainId"`
			ToToken   struct {
				Address  string `json:"address"`
				ChainID  int64  `json:"chainId"`
				CoinKey  string `json:"coinKey"`
				Decimals int64  `json:"decimals"`
				LogoURI  string `json:"logoURI"`
				Name     string `json:"name"`
				PriceUSD string `json:"priceUSD"`
				Symbol   string `json:"symbol"`
			} `json:"toToken"`
		} `json:"action"`
		Code      string `json:"code"`
		ErrorType string `json:"errorType"`
		Message   string `json:"message"`
		Tool      string `json:"tool"`
	} `json:"errors"`
	Message string `json:"message"`
}
