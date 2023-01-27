package simulation

import (
	"math/big"
	"time"
)

type SimulateTxRequest struct {
	Chain     string   `json:"chain"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	InputData []string `json:"inputData"`
	Gas       string   `json:"gas"`
	GasPrice  string   `json:"gasPrice"`
	Value     string   `json:"value"`
	Website   string   `json:"website"`
}

type SimulateTxResponse struct {
	ChainId           string `json:"chain_id"`
	ProviderId        string `json:"provider_id"`
	SimulationResult  bool   `json:"simulation_result"`
	Type              string `json:"type,omitempty"`
	Action            string `json:"action,omitempty"`
	ActionMessage     string `json:"action_message,omitempty"`
	GasUsed           string `json:"gas_used"`
	HumanReadableForm string `json:"human_readable_form,omitempty"`
	SimulationData    Data   `json:"simulation_data,omitempty"`
}

type Data struct {
	AssetsSent     []AssetsSent     `json:"assets_sent"`
	AssetsReceived []AssetsReceived `json:"assets_received"`
	ApprovalAssets []ApprovalAssets `json:"approval_assets"`
}

type AssetsSent struct {
	ContractAddress string `json:"contract_address"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Decimals        string `json:"decimals"`
}

type AssetsReceived struct {
	ContractAddress string `json:"contract_address"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Decimals        string `json:"decimals"`
}

type ApprovalAssets struct {
	ContractAddress string `json:"contract_address"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Decimals        string `json:"decimals"`
}

//BlowFish struct

type BlowFishEVMTxBody struct {
	TxObject    BlowFishTxObject `json:"txObject"`
	Metadata    BlowFishMetadata `json:"metadata"`
	UserAccount string           `json:"userAccount"`
}

type BlowFishTxObject struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
}

type BlowFishMetadata struct {
	Origin string `json:"origin"`
}

// BlowFishSolanaTxBody Struct
type BlowFishSolanaTxBody struct {
	Transactions []string `json:"transactions"`
	Metadata     struct {
		Origin string `json:"origin"`
	} `json:"metadata"`
	UserAccount string `json:"user_account"`
}

type BlowFishSolanaResponse struct {
	Status   string `json:"status"`
	Action   string `json:"action"`
	Warnings []struct {
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Kind     string `json:"kind"`
	} `json:"warnings"`
	SimulationResults struct {
		IsRecentBlockhashExpired bool `json:"isRecentBlockhashExpired"`
		ExpectedStateChanges     []struct {
			HumanReadableDiff string `json:"humanReadableDiff"`
			SuggestedColor    string `json:"suggestedColor"`
			RawInfo           struct {
				Kind string `json:"kind"`
				Data struct {
					Symbol   string `json:"symbol"`
					Name     string `json:"name"`
					Decimals int    `json:"decimals"`
					Diff     struct {
						Sign   string `json:"sign"`
						Digits int    `json:"digits"`
					} `json:"diff"`
				} `json:"data"`
			} `json:"rawInfo"`
		} `json:"expectedStateChanges"`
		Error struct {
			Kind               string `json:"kind"`
			HumanReadableError string `json:"humanReadableError"`
		} `json:"error"`
		Raw struct {
			Err struct {
				InstructionError []interface{} `json:"InstructionError"`
			} `json:"err"`
			Logs     []string `json:"logs"`
			Accounts []struct {
				Lamports   int64    `json:"lamports"`
				Data       []string `json:"data"`
				Owner      string   `json:"owner"`
				Executable bool     `json:"executable"`
				RentEpoch  int      `json:"rentEpoch"`
			} `json:"accounts"`
			UnitsConsumed int         `json:"unitsConsumed"`
			ReturnData    interface{} `json:"returnData"`
		} `json:"raw"`
	} `json:"simulationResults"`
}

type BlowFishResponse struct {
	Action   string `json:"action"`
	Warnings []struct {
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Kind     string `json:"kind"`
	} `json:"warnings"`
	SimulationResults struct {
		ExpectedStateChanges []struct {
			HumanReadableDiff string `json:"humanReadableDiff"`
			RawInfo           struct {
				Kind string `json:"kind"`
				Data struct {
					Name     string `json:"name"`
					Symbol   string `json:"symbol"`
					Decimals int    `json:"decimals"`
					Spender  struct {
						Kind    string `json:"kind"`
						Address string `json:"address"`
					} `json:"spender"`
					Amount struct {
						Before string `json:"before"`
						After  string `json:"after"`
					} `json:"amount"`
					Contract struct {
						Kind    string `json:"kind"`
						Address string `json:"address"`
					} `json:"contract"`
					Owner struct {
						Kind    string `json:"kind"`
						Address string `json:"address"`
					} `json:"owner"`
					Asset struct {
						Address  string        `json:"address"`
						Symbol   string        `json:"symbol"`
						Name     string        `json:"name"`
						Decimals int           `json:"decimals"`
						Verified bool          `json:"verified"`
						Lists    []interface{} `json:"lists"`
						ImageUrl interface{}   `json:"imageUrl"`
						Price    struct {
							Source              string  `json:"source"`
							LastUpdatedAt       int     `json:"last_updated_at"`
							DollarValuePerToken float64 `json:"dollar_value_per_token"`
						} `json:"price"`
					} `json:"asset"`
				} `json:"data"`
			} `json:"rawInfo"`
		} `json:"expectedStateChanges"`
		Error struct {
			Kind               string `json:"kind"`
			HumanReadableError string `json:"humanReadableError"`
			ParsedErrorMessage string `json:"parsedErrorMessage"`
		} `json:"error"`
		Gas struct {
			GasLimit string `json:"gasLimit"`
		} `json:"gas"`
	} `json:"simulationResults"`
}

// SignAssistResponse struct
type SignAssistResponse struct {
	Success    bool `json:"success"`
	Simulation struct {
		AssetsSent []struct {
			ContractAddress string `json:"contract_address"`
			TokenType       string `json:"token_type"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
			Decimals        int    `json:"decimals"`
			Amount          string `json:"amount"`
			Image           string `json:"image"`
		} `json:"assets_sent"`
		AssetsReceived []struct {
			ContractAddress string `json:"contract_address"`
			TokenType       string `json:"token_type"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
			Decimals        int    `json:"decimals"`
			Amount          string `json:"amount"`
			Image           string `json:"image"`
		} `json:"assets_received"`
		AssetsApprovals []struct {
			ContractAddress string `json:"contract_address"`
			TokenType       string `json:"token_type"`
			Name            string `json:"name"`
			Symbol          string `json:"symbol"`
			Decimals        int    `json:"decimals"`
			Amount          string `json:"amount"`
			Image           string `json:"image"`
		} `json:"assets_approvals"`
		TransactionInfo struct {
			From             string `json:"from"`
			To               string `json:"to"`
			EthValue         int    `json:"eth_value"`
			Gas              int    `json:"gas"`
			GasUsed          int    `json:"gas_used"`
			Method           string `json:"method"`
			SimulationResult struct {
				Status bool `json:"status"`
			} `json:"simulationResult"`
		} `json:"transaction_info"`
	} `json:"simulation"`
}

type SignAssistError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SignAssistTxBody struct {
	NetworkId         int                 `json:"network_id"`
	TransactionParams []TransactionParams `json:"transaction_params"`
}

type TransactionParams struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gas_price"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

//Tenderly struct

type TenderlyTxBody struct {
	NetworkId      string   `json:"network_id"`
	From           string   `json:"from"`
	To             string   `json:"to"`
	Input          string   `json:"input"`
	Gas            int      `json:"gas"`
	GasPrice       string   `json:"gas_price"`
	Value          *big.Int `json:"value"`
	SaveIfFails    bool     `json:"save_if_fails"`
	Save           bool     `json:"save"`
	SimulationType string   `json:"simulation_type"`
}

type TenderlyTxErrorResponse struct {
	Error struct {
		Message string      `json:"message"`
		Slug    string      `json:"slug"`
		Data    interface{} `json:"data"`
	} `json:"error"`
}

type Calls []struct {
	To string `json:"to"`
}

type Logs []struct {
	Raw struct {
		Address string `json:"address"`
		Data    string `json:"data"`
	} `json:"raw"`
}

type TenderlyTxResponse struct {
	Transaction struct {
		Hash              string      `json:"hash"`
		BlockHash         string      `json:"block_hash"`
		BlockNumber       int         `json:"block_number"`
		From              string      `json:"from"`
		Gas               int         `json:"gas"`
		GasPrice          int         `json:"gas_price"`
		GasFeeCap         int         `json:"gas_fee_cap"`
		GasTipCap         int         `json:"gas_tip_cap"`
		CumulativeGasUsed int         `json:"cumulative_gas_used"`
		GasUsed           int         `json:"gas_used"`
		EffectiveGasPrice int         `json:"effective_gas_price"`
		Input             string      `json:"input"`
		Nonce             int         `json:"nonce"`
		To                string      `json:"to"`
		Index             int         `json:"index"`
		Value             string      `json:"value"`
		AccessList        interface{} `json:"access_list"`
		Status            bool        `json:"status"`
		Addresses         interface{} `json:"addresses"`
		ContractIds       interface{} `json:"contract_ids"`
		NetworkID         string      `json:"network_id"`
		Timestamp         time.Time   `json:"timestamp"`
		FunctionSelector  string      `json:"function_selector"`
		L1BlockNumber     int         `json:"l1_block_number"`
		L1Timestamp       int         `json:"l1_timestamp"`
		TransactionInfo   struct {
			ContractID      string      `json:"contract_id"`
			BlockNumber     int         `json:"block_number"`
			TransactionID   string      `json:"transaction_id"`
			ContractAddress string      `json:"contract_address"`
			Method          interface{} `json:"method"`
			Parameters      interface{} `json:"parameters"`
			IntrinsicGas    int         `json:"intrinsic_gas"`
			RefundGas       int         `json:"refund_gas"`
			CallTrace       struct {
				Hash             string    `json:"hash"`
				ContractName     string    `json:"contract_name"`
				FunctionPc       int       `json:"function_pc"`
				FunctionOp       string    `json:"function_op"`
				AbsolutePosition int       `json:"absolute_position"`
				CallerPc         int       `json:"caller_pc"`
				CallerOp         string    `json:"caller_op"`
				CallType         string    `json:"call_type"`
				From             string    `json:"from"`
				FromBalance      string    `json:"from_balance"`
				To               string    `json:"to"`
				ToBalance        string    `json:"to_balance"`
				Value            string    `json:"value"`
				BlockTimestamp   time.Time `json:"block_timestamp"`
				Gas              int       `json:"gas"`
				GasUsed          int       `json:"gas_used"`
				IntrinsicGas     int       `json:"intrinsic_gas"`
				RefundGas        int       `json:"refund_gas"`
				Input            string    `json:"input"`
				BalanceDiff      []struct {
					Address  string `json:"address"`
					Original string `json:"original"`
					Dirty    string `json:"dirty"`
					IsMiner  bool   `json:"is_miner"`
				} `json:"balance_diff"`
				NonceDiff []struct {
					Address  string `json:"address"`
					Original string `json:"original"`
					Dirty    string `json:"dirty"`
				} `json:"nonce_diff"`
				StateDiff []struct {
					Soltype  interface{} `json:"soltype"`
					Original interface{} `json:"original"`
					Dirty    interface{} `json:"dirty"`
					Raw      []struct {
						Address  string `json:"address"`
						Key      string `json:"key"`
						Original string `json:"original"`
						Dirty    string `json:"dirty"`
					} `json:"raw"`
				} `json:"state_diff"`
				Logs []struct {
					Name      string      `json:"name"`
					Anonymous bool        `json:"anonymous"`
					Inputs    interface{} `json:"inputs"`
					Raw       struct {
						Address string   `json:"address"`
						Topics  []string `json:"topics"`
						Data    string   `json:"data"`
					} `json:"raw"`
				} `json:"logs"`
				Output        string      `json:"output"`
				DecodedOutput interface{} `json:"decoded_output"`
				NetworkID     string      `json:"network_id"`
				Calls         []struct {
					Hash             string      `json:"hash"`
					FunctionPc       int         `json:"function_pc"`
					FunctionOp       string      `json:"function_op"`
					AbsolutePosition int         `json:"absolute_position"`
					CallerPc         int         `json:"caller_pc"`
					CallerOp         string      `json:"caller_op"`
					CallType         string      `json:"call_type"`
					From             string      `json:"from"`
					FromBalance      string      `json:"from_balance"`
					To               string      `json:"to"`
					ToBalance        string      `json:"to_balance"`
					Value            interface{} `json:"value"`
					Caller           struct {
						Address string `json:"address"`
						Balance string `json:"balance"`
					} `json:"caller"`
					BlockTimestamp time.Time   `json:"block_timestamp"`
					Gas            int         `json:"gas"`
					GasUsed        int         `json:"gas_used"`
					RefundGas      int         `json:"refund_gas"`
					Input          string      `json:"input"`
					Output         string      `json:"output"`
					DecodedOutput  interface{} `json:"decoded_output"`
					NetworkID      string      `json:"network_id"`
					Calls          []struct {
						Hash             string      `json:"hash"`
						FunctionPc       int         `json:"function_pc"`
						FunctionOp       string      `json:"function_op"`
						AbsolutePosition int         `json:"absolute_position"`
						CallerPc         int         `json:"caller_pc"`
						CallerOp         string      `json:"caller_op"`
						CallType         string      `json:"call_type"`
						From             string      `json:"from"`
						FromBalance      string      `json:"from_balance"`
						To               string      `json:"to"`
						ToBalance        string      `json:"to_balance"`
						Value            string      `json:"value"`
						BlockTimestamp   time.Time   `json:"block_timestamp"`
						Gas              int         `json:"gas"`
						GasUsed          int         `json:"gas_used"`
						RefundGas        int         `json:"refund_gas"`
						Input            string      `json:"input"`
						Output           string      `json:"output"`
						DecodedOutput    interface{} `json:"decoded_output"`
						NetworkID        string      `json:"network_id"`
						Calls            Calls       `json:"calls"`
					} `json:"calls"`
				} `json:"calls"`
			} `json:"call_trace"`
			StackTrace interface{} `json:"stack_trace"`
			Logs       []struct {
				Name      string      `json:"name"`
				Anonymous bool        `json:"anonymous"`
				Inputs    interface{} `json:"inputs"`
				Raw       struct {
					Address string   `json:"address"`
					Topics  []string `json:"topics"`
					Data    string   `json:"data"`
				} `json:"raw"`
			} `json:"logs"`
			BalanceDiff []struct {
				Address  string `json:"address"`
				Original string `json:"original"`
				Dirty    string `json:"dirty"`
				IsMiner  bool   `json:"is_miner"`
			} `json:"balance_diff"`
			NonceDiff []struct {
				Address  string `json:"address"`
				Original string `json:"original"`
				Dirty    string `json:"dirty"`
			} `json:"nonce_diff"`
			StateDiff []struct {
				Soltype  interface{} `json:"soltype"`
				Original interface{} `json:"original"`
				Dirty    interface{} `json:"dirty"`
				Raw      []struct {
					Address  string `json:"address"`
					Key      string `json:"key"`
					Original string `json:"original"`
					Dirty    string `json:"dirty"`
				} `json:"raw"`
			} `json:"state_diff"`
			RawStateDiff interface{} `json:"raw_state_diff"`
			ConsoleLogs  interface{} `json:"console_logs"`
			CreatedAt    time.Time   `json:"created_at"`
		} `json:"transaction_info"`
		ErrorMessage string      `json:"error_message"`
		Method       string      `json:"method"`
		DecodedInput interface{} `json:"decoded_input"`
		CallTrace    []struct {
			CallType     string `json:"call_type"`
			From         string `json:"from"`
			To           string `json:"to"`
			Gas          int    `json:"gas"`
			GasUsed      int    `json:"gas_used"`
			Value        string `json:"value,omitempty"`
			Subtraces    int    `json:"subtraces,omitempty"`
			Type         string `json:"type"`
			Input        string `json:"input"`
			Output       string `json:"output,omitempty"`
			FromBalance  string `json:"fromBalance,omitempty"`
			ToBalance    string `json:"toBalance,omitempty"`
			TraceAddress []int  `json:"trace_address,omitempty"`
			OutOff       int    `json:"outOff,omitempty"`
			GasIn        int    `json:"gas_in,omitempty"`
			GasCost      int    `json:"gas_cost,omitempty"`
			OutLen       int    `json:"outLen,omitempty"`
		} `json:"call_trace"`
	} `json:"transaction"`
	Simulation struct {
		ID               string      `json:"id"`
		ProjectID        string      `json:"project_id"`
		OwnerID          string      `json:"owner_id"`
		NetworkID        string      `json:"network_id"`
		BlockNumber      int         `json:"block_number"`
		TransactionIndex int         `json:"transaction_index"`
		From             string      `json:"from"`
		To               string      `json:"to"`
		Input            string      `json:"input"`
		Gas              int         `json:"gas"`
		GasPrice         string      `json:"gas_price"`
		Value            string      `json:"value"`
		Status           bool        `json:"status"`
		AccessList       interface{} `json:"access_list"`
		QueueOrigin      string      `json:"queue_origin"`
		BlockHeader      struct {
			Number           string      `json:"number"`
			Hash             string      `json:"hash"`
			StateRoot        string      `json:"stateRoot"`
			ParentHash       string      `json:"parentHash"`
			Sha3Uncles       string      `json:"sha3Uncles"`
			TransactionsRoot string      `json:"transactionsRoot"`
			ReceiptsRoot     string      `json:"receiptsRoot"`
			LogsBloom        string      `json:"logsBloom"`
			Timestamp        string      `json:"timestamp"`
			Difficulty       string      `json:"difficulty"`
			GasLimit         string      `json:"gasLimit"`
			GasUsed          string      `json:"gasUsed"`
			Miner            string      `json:"miner"`
			ExtraData        string      `json:"extraData"`
			MixHash          string      `json:"mixHash"`
			Nonce            string      `json:"nonce"`
			BaseFeePerGas    string      `json:"baseFeePerGas"`
			Size             string      `json:"size"`
			TotalDifficulty  string      `json:"totalDifficulty"`
			Uncles           interface{} `json:"uncles"`
			Transactions     interface{} `json:"transactions"`
		} `json:"block_header"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"simulation"`
	Contracts           []interface{} `json:"contracts"`
	GeneratedAccessList []interface{} `json:"generated_access_list"`
}
