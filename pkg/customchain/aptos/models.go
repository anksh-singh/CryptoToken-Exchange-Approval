package aptos

type ListTransactions []struct {
	Version             string      `json:"version"`
	Hash                string      `json:"hash"`
	StateChangeHash     string      `json:"state_change_hash"`
	EventRootHash       string      `json:"event_root_hash"`
	StateCheckpointHash interface{} `json:"state_checkpoint_hash"`
	GasUsed             string      `json:"gas_used"`
	Success             bool        `json:"success"`
	VMStatus            string      `json:"vm_status"`
	AccumulatorRootHash string      `json:"accumulator_root_hash"`
	Changes             []struct {
		Address      string `json:"address,omitempty"`
		StateKeyHash string `json:"state_key_hash"`
		Data         struct {
			Type string `json:"type"`
			Data struct {
				Coin struct {
					Value string `json:"value"`
				} `json:"coin"`
				DepositEvents struct {
					Counter string `json:"counter"`
					GUID    struct {
						ID struct {
							Addr        string `json:"addr"`
							CreationNum string `json:"creation_num"`
						} `json:"id"`
					} `json:"guid"`
				} `json:"deposit_events"`
				Frozen         bool `json:"frozen"`
				WithdrawEvents struct {
					Counter string `json:"counter"`
					GUID    struct {
						ID struct {
							Addr        string `json:"addr"`
							CreationNum string `json:"creation_num"`
						} `json:"id"`
					} `json:"guid"`
				} `json:"withdraw_events"`
			} `json:"data"`
		} `json:"data"`
		Type   string `json:"type"`
		Handle string `json:"handle,omitempty"`
		Key    string `json:"key,omitempty"`
		Value  string `json:"value,omitempty"`
	} `json:"changes"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Payload                 struct {
		Function      string        `json:"function"`
		TypeArguments []interface{} `json:"type_arguments"`
		Arguments     []string      `json:"arguments"`
		Type          string        `json:"type"`
	} `json:"payload"`
	Signature struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
		Type      string `json:"type"`
	} `json:"signature"`
	Events []struct {
		GUID struct {
			CreationNumber string `json:"creation_number"`
			AccountAddress string `json:"account_address"`
		} `json:"guid"`
		SequenceNumber string `json:"sequence_number"`
		Type           string `json:"type"`
		Data           struct {
			Amount string `json:"amount"`
		} `json:"data"`
	} `json:"events"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type TxByHash struct {
	Version             string      `json:"version"`
	Hash                string      `json:"hash"`
	StateChangeHash     string      `json:"state_change_hash"`
	EventRootHash       string      `json:"event_root_hash"`
	StateCheckpointHash interface{} `json:"state_checkpoint_hash"`
	GasUsed             string      `json:"gas_used"`
	Success             bool        `json:"success"`
	VMStatus            string      `json:"vm_status"`
	AccumulatorRootHash string      `json:"accumulator_root_hash"`
	Changes             []struct {
		Address      string `json:"address,omitempty"`
		StateKeyHash string `json:"state_key_hash"`
		Data         struct {
			Type string `json:"type"`
			Data struct {
				Coin struct {
					Value string `json:"value"`
				} `json:"coin"`
				DepositEvents struct {
					Counter string `json:"counter"`
					GUID    struct {
						ID struct {
							Addr        string `json:"addr"`
							CreationNum string `json:"creation_num"`
						} `json:"id"`
					} `json:"guid"`
				} `json:"deposit_events"`
				Frozen         bool `json:"frozen"`
				WithdrawEvents struct {
					Counter string `json:"counter"`
					GUID    struct {
						ID struct {
							Addr        string `json:"addr"`
							CreationNum string `json:"creation_num"`
						} `json:"id"`
					} `json:"guid"`
				} `json:"withdraw_events"`
			} `json:"data"`
		} `json:"data"`
		Type   string `json:"type"`
		Handle string `json:"handle,omitempty"`
		Key    string `json:"key,omitempty"`
		Value  string `json:"value,omitempty"`
	} `json:"changes"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Payload                 struct {
		Function      string        `json:"function"`
		TypeArguments []interface{} `json:"type_arguments"`
		Arguments     []string      `json:"arguments"`
		Type          string        `json:"type"`
	} `json:"payload"`
	Signature struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
		Type      string `json:"type"`
	} `json:"signature"`
	Events []struct {
		GUID struct {
			CreationNumber string `json:"creation_number"`
			AccountAddress string `json:"account_address"`
		} `json:"guid"`
		SequenceNumber string `json:"sequence_number"`
		Type           string `json:"type"`
		Data           struct {
			Amount string `json:"amount"`
		} `json:"data"`
	} `json:"events"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type Blocks struct {
	BlockHeight    string      `json:"block_height"`
	BlockHash      string      `json:"block_hash"`
	BlockTimestamp string      `json:"block_timestamp"`
	FirstVersion   string      `json:"first_version"`
	LastVersion    string      `json:"last_version"`
	Transactions   interface{} `json:"transactions"`
}

type TxPayload struct {
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Payload                 struct {
		Type          string   `json:"type"`
		Function      string   `json:"function"`
		TypeArguments []string `json:"type_arguments"`
		Arguments     []string `json:"arguments"`
	} `json:"payload"`
	Signature struct {
		Type      string `json:"type"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	} `json:"signature"`
}

type SendTxResponse struct {
	Hash                    string `json:"hash"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Payload                 struct {
		Type          string        `json:"type"`
		Function      string        `json:"function"`
		TypeArguments []string      `json:"type_arguments"`
		Arguments     []interface{} `json:"arguments"`
	} `json:"payload"`
	Signature struct {
		Type      string `json:"type"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	} `json:"signature"`
}

type BalanceResp struct {
	Data []struct {
		CoinType string  `json:"coin_type"`
		Value    int     `json:"value"`
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		Decimals int     `json:"decimals"`
		Amount   float64 `json:"amount"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
