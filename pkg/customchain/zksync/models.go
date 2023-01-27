package zksync

import "time"

type TokenQuoteResponse struct {
	L1Address string `json:"l1Address"`
	L2Address string `json:"l2Address"`
	Address   string `json:"address"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Decimals  int    `json:"decimals"`
	UsdPrice  string `json:"usdPrice"`
}

type BalanceResponse struct {
	Message string `json:"message"`
	Result  []struct {
		Balance         string `json:"balance"`
		ContractAddress string `json:"contractAddress"`
		Decimals        string `json:"decimals"`
		Name            string `json:"name"`
		Symbol          string `json:"symbol"`
		Type            string `json:"type"`
	} `json:"result"`
	Status string `json:"status"`
}

type ListTransactionResponse struct {
	List []struct {
		TransactionHash string `json:"transactionHash"`
		Data            struct {
			ContractAddress string      `json:"contractAddress"`
			Calldata        string      `json:"calldata"`
			Value           string      `json:"value"`
			FactoryDeps     interface{} `json:"factoryDeps"`
		} `json:"data"`
		IsL1Originated   bool      `json:"isL1Originated"`
		Status           string    `json:"status"`
		Fee              string    `json:"fee"`
		Nonce            int       `json:"nonce"`
		BlockNumber      int       `json:"blockNumber"`
		L1BatchNumber    int       `json:"l1BatchNumber"`
		BlockHash        string    `json:"blockHash"`
		IndexInBlock     int       `json:"indexInBlock"`
		InitiatorAddress string    `json:"initiatorAddress"`
		ReceivedAt       time.Time `json:"receivedAt"`
		EthCommitTxHash  string    `json:"ethCommitTxHash"`
		EthProveTxHash   string    `json:"ethProveTxHash"`
		EthExecuteTxHash string    `json:"ethExecuteTxHash"`
		Erc20Transfers   []struct {
			TokenInfo struct {
				L1Address string `json:"l1Address"`
				L2Address string `json:"l2Address"`
				Address   string `json:"address"`
				Symbol    string `json:"symbol"`
				Name      string `json:"name"`
				Decimals  int    `json:"decimals"`
				UsdPrice  string `json:"usdPrice"`
			} `json:"tokenInfo"`
			From   string `json:"from"`
			To     string `json:"to"`
			Amount string `json:"amount"`
		} `json:"erc20Transfers"`
		Transfer struct {
			TokenInfo struct {
				L1Address string `json:"l1Address"`
				L2Address string `json:"l2Address"`
				Address   string `json:"address"`
				Symbol    string `json:"symbol"`
				Name      string `json:"name"`
				Decimals  int    `json:"decimals"`
				UsdPrice  string `json:"usdPrice"`
			} `json:"tokenInfo"`
			From   string `json:"from"`
			To     string `json:"to"`
			Amount string `json:"amount"`
		} `json:"transfer"`
		BalanceChanges []struct {
			TokenInfo struct {
				L1Address string `json:"l1Address"`
				L2Address string `json:"l2Address"`
				Address   string `json:"address"`
				Symbol    string `json:"symbol"`
				Name      string `json:"name"`
				Decimals  int    `json:"decimals"`
				UsdPrice  string `json:"usdPrice"`
			} `json:"tokenInfo"`
			From   string `json:"from"`
			To     string `json:"to"`
			Amount string `json:"amount"`
			Type   string `json:"type"`
		} `json:"balanceChanges"`
	} `json:"list"`
	Total int `json:"total"`
}

type TokenPriceUSD struct {
	Request struct {
		Network    string `json:"network"`
		ApiVersion string `json:"apiVersion"`
		Resource   string `json:"resource"`
		Args       struct {
			TokenLike string `json:"token_like"`
			Currency  string `json:"currency"`
		} `json:"args"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"request"`
	Status string      `json:"status"`
	Error  interface{} `json:"error"`
	Result struct {
		TokenId     int    `json:"tokenId"`
		TokenSymbol string `json:"tokenSymbol"`
		PriceIn     string `json:"priceIn"`
		Decimals    int    `json:"decimals"`
		Price       string `json:"price"`
	} `json:"result"`
}
