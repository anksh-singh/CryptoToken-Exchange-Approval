package rpc

import (
	"time"
)

////TODO: Move services to a common model
//type Services struct {
//	Http                     *utils.HttpRequest
//	CoinGecko                *coingecko.CoinGecko
//	Helper                   *utils.Helpers
//	Covalent                 *covalent.CovalentService
//	Unmarshall               *unmarshal.UnmarshallService
//	ZeroX                    *_x.OXService
//	CocoSwapTokenExchange    *cocoswap.TokenExchangeStruct
//	UniSwapTokenExchange     *ubeswap.TokenExchangeStruct
//	DoDoExTokenExchange      *dodoEth.TokenExchangeStruct
//	DoDoExTokenExchangeCache *dodoEth.TokenExchangeStructCache
//	Debank                   *openApi.OpenAPI
//	DodoSwap                 *dodo.ServiceDodo
//	TrustWallet              trustwallet.ITrustWallet
//}

type BlockNativeGasPrice struct {
	System             string `json:"system"`
	Network            string `json:"network"`
	Unit               string `json:"unit"`
	MaxPrice           int    `json:"maxPrice"`
	CurrentBlockNumber int    `json:"currentBlockNumber"`
	MsSinceLastBlock   int    `json:"msSinceLastBlock"`
	BlockPrices        []struct {
		BlockNumber               int     `json:"blockNumber"`
		EstimatedTransactionCount int     `json:"estimatedTransactionCount"`
		BaseFeePerGas             float64 `json:"baseFeePerGas"`
		EstimatedPrices           []struct {
			Confidence           float64     `json:"confidence"`
			Price                interface{} `json:"price"`
			MaxPriorityFeePerGas float64     `json:"maxPriorityFeePerGas"`
			MaxFeePerGas         float64     `json:"maxFeePerGas"`
		} `json:"estimatedPrices"`
	} `json:"blockPrices"`
	EstimatedBaseFees []struct {
		Pending1 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+1,omitempty"`
		Pending2 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+2,omitempty"`
		Pending3 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+3,omitempty"`
		Pending4 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+4,omitempty"`
		Pending5 []struct {
			Confidence int     `json:"confidence"`
			BaseFee    float64 `json:"baseFee"`
		} `json:"pending+5,omitempty"`
	} `json:"estimatedBaseFees"`
}

type TokenAllowance struct {
	Contract string `json:"contract"`
	Owner    string `json:"owner"`
	Spender  string `json:"spender"`
}

type AllowanceRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type ContractABIRequest struct {
	Contract string `json:"contract"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    int64  `json:"value"`
	Data     string `json:"data"`
	Chain    string `json:"chain"`
	Method   string `json:"method"`
}
type TomoWalletBalances []struct {
	Address      string  `json:"address"`
	TokenAddress string  `json:"tokenAddress"`
	Decimals     int     `json:"decimals"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	TotalSupply  string  `json:"totalSupply"`
	Icon         string  `json:"icon"`
	Type         string  `json:"type"`
	UsdPrice     float64 `json:"usdPrice"`
	Balance      string  `json:"balance"`
	Verified     bool    `json:"verified"`
}

type TrcTxs struct {
	ID               string    `json:"_id"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	TransactionHash  string    `json:"transactionHash"`
	Address          string    `json:"address"`
	BlockHash        string    `json:"blockHash"`
	BlockNumber      int       `json:"blockNumber"`
	CreatedAt        time.Time `json:"createdAt"`
	Data             string    `json:"data"`
	Timestamp        time.Time `json:"timestamp"`
	TransactionIndex int       `json:"transactionIndex"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Value            string    `json:"value"`
	ValueNumber      float64   `json:"valueNumber"`
	Symbol           string    `json:"symbol"`
	Decimals         int       `json:"decimals"`
	BlockTime        time.Time `json:"blockTime"`
}

type TomoScan struct {
	Total       int `json:"total"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	Pages       int `json:"pages"`
	Items       []struct {
		BlockHash         string    `json:"blockHash"`
		BlockNumber       int       `json:"blockNumber"`
		CumulativeGasUsed int       `json:"cumulativeGasUsed"`
		From              string    `json:"from"`
		Gas               int       `json:"gas"`
		GasPrice          string    `json:"gasPrice"`
		GasUsed           int       `json:"gasUsed"`
		Hash              string    `json:"hash"`
		ITx               int       `json:"i_tx"`
		Nonce             int       `json:"nonce"`
		Status            bool      `json:"status"`
		Timestamp         time.Time `json:"timestamp"`
		To                string    `json:"to"`
		TransactionIndex  int       `json:"transactionIndex"`
		Value             string    `json:"value"`
		FromModel         struct {
			AccountName interface{} `json:"accountName"`
		} `json:"from_model"`
		ToModel struct {
			AccountName string `json:"accountName"`
		} `json:"to_model"`
	} `json:"items"`
}
type TomoTxReceipt struct {
	FromModel struct {
		MinedBlock    int       `json:"minedBlock"`
		RewardCount   int       `json:"rewardCount"`
		LogCount      int       `json:"logCount"`
		Status        bool      `json:"status"`
		ID            string    `json:"_id"`
		Hash          string    `json:"hash"`
		Balance       string    `json:"balance"`
		BalanceNumber float64   `json:"balanceNumber"`
		Code          string    `json:"code"`
		CreatedAt     time.Time `json:"createdAt"`
		IsToken       bool      `json:"isToken"`
		UpdatedAt     time.Time `json:"updatedAt"`
	} `json:"from_model"`
	Status            bool          `json:"status"`
	ITx               int           `json:"i_tx"`
	_ID               string        `json:"_id"`
	BlockHash         string        `json:"blockHash"`
	BlockNumber       int           `json:"blockNumber"`
	From              string        `json:"from"`
	Gas               int           `json:"gas"`
	GasPrice          string        `json:"gasPrice"`
	Hash              string        `json:"hash"`
	Input             string        `json:"input"`
	Nonce             int           `json:"nonce"`
	To                string        `json:"to"`
	TransactionIndex  int           `json:"transactionIndex"`
	Value             string        `json:"value"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
	CumulativeGasUsed int           `json:"cumulativeGasUsed"`
	GasUsed           int           `json:"gasUsed"`
	Timestamp         time.Time     `json:"timestamp"`
	ID                string        `json:"id"`
	ToModel           interface{}   `json:"to_model"`
	Trc20Txs          []TrcTxs      `json:"trc20Txs"`
	Trc21Txs          []TrcTxs      `json:"trc21Txs"`
	Trc21FeeFund      int           `json:"trc21FeeFund"`
	Trc721Txs         []interface{} `json:"trc721Txs"`
	LatestBlockNumber int           `json:"latestBlockNumber"`
	InputData         string        `json:"inputData"`
	ExtraInfo         []interface{} `json:"extraInfo"`
}

type ExchangeSwapResponseV1 struct {
	Value    string `json:"value"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasLimit string `json:"gas_limit"`
	Data     string `json:"data"`
	Txlink   string `json:"txlink"`
}

type ExchangeQuoteResponseV1 struct {
	FromTokenPrice       string `json:"fromTokenPrice"`
	MinimumReceived      string `json:"minimumReceived"`
	PriceImpact          string `json:"priceImpact"`
	ResAmount            string `json:"resAmount"`
	ResPricePerFromToken string `json:"resPricePerFromToken"`
	ResPricePerToToken   string `json:"resPricePerToToken"`
	ToTokenPrice         string `json:"toTokenPrice"`
	ApproveAddress       string `json:"approveAddress"`
}

type ExchangeTokenResponseV1 []struct {
	LogoURL       string `json:"logo_url"`
	TokenAddress  string `json:"token_address"`
	TokenDecimals string `json:"token_decimals"`
	TokenLogoURL  string `json:"token_logo_url"`
	TokenName     string `json:"token_name"`
	TokenSymbol   string `json:"token_symbol"`
}

type V1ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
type V1FailureResponse struct {
	Message string `json:"message"`
}


type OpportunityData struct {
	Apr                         interface{} `json:"apr"`
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