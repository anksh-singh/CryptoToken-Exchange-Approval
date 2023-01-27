package near

type BalanceRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  BalanceParams `json:"params"`
	Id      int           `json:"id"`
}

type BalanceParams struct {
	RequestType string `json:"request_type"`
	Finality    string `json:"finality"`
	AccountId   string `json:"account_id"`
}

type BalanceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Amount        string `json:"amount"`
		BlockHash     string `json:"block_hash"`
		BlockHeight   int    `json:"block_height"`
		CodeHash      string `json:"code_hash"`
		Locked        string `json:"locked"`
		StoragePaidAt int    `json:"storage_paid_at"`
		StorageUsage  int    `json:"storage_usage"`
	} `json:"result"`
	Id int `json:"id"`
}
type HistoryResponse struct {
	BlockHash      string `json:"block_hash"`
	BlockTimestamp string `json:"block_timestamp"`
	Hash           string `json:"hash"`
	ActionIndex    int    `json:"action_index"`
	SignerId       string `json:"signer_id"`
	ReceiverId     string `json:"receiver_id"`
	ActionKind     string `json:"action_kind"`
	Args           struct {
		AccessKey struct {
			Nonce      int `json:"nonce"`
			Permission struct {
				PermissionKind    string `json:"permission_kind"`
				PermissionDetails struct {
					Allowance   interface{}   `json:"allowance"`
					ReceiverId  string        `json:"receiver_id"`
					MethodNames []interface{} `json:"method_names"`
				} `json:"permission_details,omitempty"`
			} `json:"permission"`
		} `json:"access_key,omitempty"`
		PublicKey string `json:"public_key,omitempty"`
		Deposit   string `json:"deposit,omitempty"`
	} `json:"args"`
}

type NonceRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	Id      int    `json:"id"`
}

type Params struct {
	RequestType string `json:"request_type"`
	Finality    string `json:"finality"`
	AccountId   string `json:"account_id"`
	PublicKey   string `json:"public_key"`
}

type NonceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockHash   string `json:"block_hash"`
		BlockHeight int    `json:"block_height"`
		Nonce       int64  `json:"nonce"`
		Permission  string `json:"permission"`
	} `json:"result"`
	Id int `json:"id"`
}

type SendTransactionRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type SendTrandsactionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      string `json:"id"`
}

type TokenPriceItem struct {
	Price   string `json:"price"`
	Decimal int    `json:"decimal"`
	Symbol  string `json:"symbol"`
}

type HistoryResponseNew struct {
	BlockHash      string `json:"block_hash"`
	BlockTimestamp string `json:"block_timestamp"`
	Hash           string `json:"hash"`
	ActionIndex    int    `json:"action_index"`
	SignerId       string `json:"signer_id"`
	ReceiverId     string `json:"receiver_id"`
	ActionKind     string `json:"action_kind"`
	Args           struct {
		Gas      int64  `json:"gas,omitempty"`
		Deposit  string `json:"deposit,omitempty"`
		ArgsJson struct {
			Amount           string `json:"amount,omitempty"`
			ReceiverId       string `json:"receiver_id,omitempty"`
			AccountId        string `json:"account_id,omitempty"`
			RegistrationOnly bool   `json:"registration_only,omitempty"`
		} `json:"args_json,omitempty"`
		ArgsBase64 string `json:"args_base64,omitempty"`
		MethodName string `json:"method_name,omitempty"`
		AccessKey  struct {
			Nonce      int `json:"nonce"`
			Permission struct {
				PermissionKind string `json:"permission_kind"`
			} `json:"permission"`
		} `json:"access_key,omitempty"`
		PublicKey string `json:"public_key,omitempty"`
	} `json:"args"`
}

type SendTxError struct {
	Jsonrpc string `json:"jsonrpc"`
	Error   struct {
		Name  string `json:"name"`
		Cause struct {
			Info struct {
			} `json:"info"`
			Name string `json:"name"`
		} `json:"cause"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TxExecutionError struct {
				InvalidTxError string `json:"InvalidTxError"`
			} `json:"TxExecutionError"`
		} `json:"data"`
	} `json:"error"`
	Id int `json:"id"`
}
