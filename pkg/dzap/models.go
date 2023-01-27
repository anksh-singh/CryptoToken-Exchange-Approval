package dzap

import (
	"github.com/umbracle/ethgo"
	"math/big"
)

type ExchangeTokenResponse struct {
	Contract string `json:"contract"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
	//QuoteRate string `json:"quote_rate"`
	//Verified bool `json:"verified"`
	//Balance  string `json:"balance"`
}

type ExchangePathRequest struct {
	ChainId  int         `json:"chainId"`
	Requests []*Requests `json:"requests"`
}

type Requests struct {
	Amount           string `json:"amount"`
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	Slippage         string `json:"slippage"`
}

type MultiExchangePath struct {
	Status string `json:"status"`
	Data   struct {
		ToTokenAmount   string `json:"toTokenAmount"`
		FromTokenAmount string `json:"fromTokenAmount"`
		EstimatedGas    int    `json:"estimatedGas"`
		FromToken       struct {
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			Address  string `json:"address"`
		} `json:"fromToken"`
		ToToken struct {
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			Address  string `json:"address"`
		} `json:"toToken"`
	} `json:"data"`
}

type ExchangeParamsRequest struct {
	ChainId    int           `json:"chainId"`
	SwapParams []*SwapParams `json:"swapParams"`
}

type SwapParams struct {
	Amount           string `json:"amount"`
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	Slippage         string `json:"slippage"`
}

type ErrorResponseV1 struct {
	Status string `json:"status"`
	Data   struct {
		Errors struct {
			StatusCode  int    `json:"statusCode"`
			Error       string `json:"error"`
			Description string `json:"description"`
			Meta        []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"meta"`
			RequestId string `json:"requestId"`
		} `json:"errors"`
		FromToken struct {
			Address string `json:"address"`
		} `json:"fromToken"`
		ToToken struct {
			Address string `json:"address"`
		} `json:"toToken"`
	} `json:"data"`
}

type ExchangeParamsResponse struct {
	ErcSwapDetails []struct {
		Executor string `json:"executor"`
		Desc     struct {
			Field1 string `json:"0"`
			Field2 string `json:"1"`
			Field3 string `json:"2"`
			Field4 string `json:"3"`
			Field5 struct {
				Type string `json:"type"`
				Hex  string `json:"hex"`
			} `json:"4"`
			Field6 struct {
				Type string `json:"type"`
				Hex  string `json:"hex"`
			} `json:"5"`
			Field7 struct {
				Type string `json:"type"`
				Hex  string `json:"hex"`
			} `json:"6"`
			Field8          string `json:"7"`
			SrcToken        string `json:"srcToken"`
			DstToken        string `json:"dstToken"`
			SrcReceiver     string `json:"srcReceiver"`
			DstReceiver     string `json:"dstReceiver"`
			Amount          string `json:"amount"`
			MinReturnAmount string `json:"minReturnAmount"`
			Flags           string `json:"flags"`
			Permit          string `json:"permit"`
		} `json:"desc"`
		RouteData       string `json:"routeData"`
		Permit          string `json:"permit"`
		MinReturnAmount string `json:"minReturnAmount"`
	} `json:"ercSwapDetails"`
	Value string `json:"value"`
}

type ErcSwapDetailsPayload struct {
	Executor        ethgo.Address          `abi:"executor"`
	Desc            SwapDescriptionPayload `abi:"desc"`
	RouteData       []byte                 `abi:"routeData"`
	Permit          []byte                 `abi:"permit"`
	MinReturnAmount *big.Int               `abi:"minReturnAmount"`
}

type SwapDescriptionPayload struct {
	SrcToken        ethgo.Address `abi:"srcToken"`
	DstToken        ethgo.Address `abi:"dstToken"`
	SrcReceiver     ethgo.Address `abi:"srcReceiver"`
	DstReceiver     ethgo.Address `abi:"dstReceiver"`
	Amount          *big.Int      `abi:"amount"`
	MinReturnAmount *big.Int      `abi:"minReturnAmount"`
	Flags           *big.Int      `abi:"flags"`
	Permit          []byte        `abi:"permit"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Config  struct {
		Transitional struct {
			SilentJSONParsing   bool `json:"silentJSONParsing"`
			ForcedJSONParsing   bool `json:"forcedJSONParsing"`
			ClarifyTimeoutError bool `json:"clarifyTimeoutError"`
		} `json:"transitional"`
		TransformRequest  []interface{} `json:"transformRequest"`
		TransformResponse []interface{} `json:"transformResponse"`
		Timeout           int           `json:"timeout"`
		XsrfCookieName    string        `json:"xsrfCookieName"`
		XsrfHeaderName    string        `json:"xsrfHeaderName"`
		MaxContentLength  int           `json:"maxContentLength"`
		MaxBodyLength     int           `json:"maxBodyLength"`
		Env               struct {
		} `json:"env"`
		Headers struct {
			Accept    string `json:"Accept"`
			UserAgent string `json:"User-Agent"`
		} `json:"headers"`
		Method string `json:"method"`
		Url    string `json:"url"`
	} `json:"config"`
	Code   string `json:"code"`
	Status int    `json:"status"`
}

const CONTRACT_ABI = `[
              {
                "inputs": [
                  {
                    "internalType": "address[]",
                    "name": "routers_",
                    "type": "address[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "bool",
                        "name": "isSupported",
                        "type": "bool"
                      },
                      {
                        "internalType": "uint256",
                        "name": "fees",
                        "type": "uint256"
                      }
                    ],
                    "internalType": "struct Router[]",
                    "name": "routerDetails_",
                    "type": "tuple[]"
                  },
                  {
                    "internalType": "enum FeeType[]",
                    "name": "feeTypes_",
                    "type": "uint8[]"
                  },
                  {
                    "internalType": "uint256[]",
                    "name": "fees_",
                    "type": "uint256[]"
                  },
                  {
                    "internalType": "address",
                    "name": "aggregationRouter_",
                    "type": "address"
                  },
                  {
                    "internalType": "address",
                    "name": "wNative_",
                    "type": "address"
                  },
                  {
                    "internalType": "address",
                    "name": "feeVault_",
                    "type": "address"
                  },
                  {
                    "internalType": "address",
                    "name": "feeDiscountNft_",
                    "type": "address"
                  }
                ],
                "stateMutability": "nonpayable",
                "type": "constructor"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": false,
                    "internalType": "enum FeeType[]",
                    "name": "feeTypes",
                    "type": "uint8[]"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256[]",
                    "name": "fees",
                    "type": "uint256[]"
                  }
                ],
                "name": "FeeUpdated",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": false,
                    "internalType": "address",
                    "name": "feeVault",
                    "type": "address"
                  }
                ],
                "name": "FeeVaultUpdated",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "sender",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "token",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      }
                    ],
                    "indexed": false,
                    "internalType": "struct Input[]",
                    "name": "input",
                    "type": "tuple[]"
                  },
                  {
                    "indexed": false,
                    "internalType": "address",
                    "name": "outputLp",
                    "type": "address"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "returnAmount",
                    "type": "uint256"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "feeBps",
                    "type": "uint256"
                  }
                ],
                "name": "LiquidityAdded",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "sender",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                  },
                  {
                    "components": [
                      {
                        "components": [
                          {
                            "internalType": "address",
                            "name": "token",
                            "type": "address"
                          },
                          {
                            "internalType": "uint256",
                            "name": "amount",
                            "type": "uint256"
                          }
                        ],
                        "internalType": "struct Input[]",
                        "name": "lpInput",
                        "type": "tuple[]"
                      },
                      {
                        "components": [
                          {
                            "internalType": "address",
                            "name": "to",
                            "type": "address"
                          },
                          {
                            "internalType": "uint256",
                            "name": "returnAmount",
                            "type": "uint256"
                          }
                        ],
                        "internalType": "struct Output[]",
                        "name": "lpOutput",
                        "type": "tuple[]"
                      }
                    ],
                    "indexed": false,
                    "internalType": "struct LPSwapInfo",
                    "name": "swapInfo",
                    "type": "tuple"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "feeBps",
                    "type": "uint256"
                  }
                ],
                "name": "LpSwapped",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "previousOwner",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "newOwner",
                    "type": "address"
                  }
                ],
                "name": "OwnershipTransferred",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": false,
                    "internalType": "address[]",
                    "name": "routers_",
                    "type": "address[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "bool",
                        "name": "isSupported",
                        "type": "bool"
                      },
                      {
                        "internalType": "uint256",
                        "name": "fees",
                        "type": "uint256"
                      }
                    ],
                    "indexed": false,
                    "internalType": "struct Router[]",
                    "name": "details",
                    "type": "tuple[]"
                  }
                ],
                "name": "RoutersUpdated",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "to",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "token",
                    "type": "address"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                  }
                ],
                "name": "TokensRescued",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "sender",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                  },
                  {
                    "components": [
                      {
                        "internalType": "contract IERC20",
                        "name": "srcToken",
                        "type": "address"
                      },
                      {
                        "internalType": "contract IERC20",
                        "name": "dstToken",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      },
                      {
                        "internalType": "uint256",
                        "name": "returnAmount",
                        "type": "uint256"
                      }
                    ],
                    "indexed": false,
                    "internalType": "struct ErcSwapInfo[]",
                    "name": "swapInfo",
                    "type": "tuple[]"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "feeBps",
                    "type": "uint256"
                  }
                ],
                "name": "TokensSwapped",
                "type": "event"
              },
              {
                "anonymous": false,
                "inputs": [
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "sender",
                    "type": "address"
                  },
                  {
                    "indexed": true,
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "token",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      }
                    ],
                    "indexed": false,
                    "internalType": "struct Input[]",
                    "name": "inputDetails",
                    "type": "tuple[]"
                  },
                  {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "feeBps",
                    "type": "uint256"
                  }
                ],
                "name": "TokensTransferred",
                "type": "event"
              },
              {
                "inputs": [],
                "name": "AGGREGATION_ROUTER",
                "outputs": [
                  {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "BPS_DENOMINATOR",
                "outputs": [
                  {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "address",
                    "name": "recipient_",
                    "type": "address"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "token",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      }
                    ],
                    "internalType": "struct Input[]",
                    "name": "details_",
                    "type": "tuple[]"
                  },
                  {
                    "internalType": "bytes[]",
                    "name": "permit_",
                    "type": "bytes[]"
                  },
                  {
                    "internalType": "uint256",
                    "name": "nftId_",
                    "type": "uint256"
                  }
                ],
                "name": "batchTransfer",
                "outputs": [],
                "stateMutability": "payable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "uint256",
                    "name": "amountA_",
                    "type": "uint256"
                  },
                  {
                    "internalType": "uint256",
                    "name": "amountB_",
                    "type": "uint256"
                  },
                  {
                    "internalType": "uint256",
                    "name": "reserveA_",
                    "type": "uint256"
                  },
                  {
                    "internalType": "uint256",
                    "name": "reserveB_",
                    "type": "uint256"
                  },
                  {
                    "internalType": "address",
                    "name": "router_",
                    "type": "address"
                  }
                ],
                "name": "calculateOptimalSwapAmount",
                "outputs": [
                  {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "enum FeeType",
                    "name": "",
                    "type": "uint8"
                  }
                ],
                "name": "fee",
                "outputs": [
                  {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "feeDiscountNft",
                "outputs": [
                  {
                    "internalType": "contract DZapDiscountNftToken",
                    "name": "",
                    "type": "address"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "feeVault",
                "outputs": [
                  {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "owner",
                "outputs": [
                  {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "renounceOwnership",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "contract IERC20",
                    "name": "token_",
                    "type": "address"
                  },
                  {
                    "internalType": "address",
                    "name": "to_",
                    "type": "address"
                  },
                  {
                    "internalType": "uint256",
                    "name": "amount_",
                    "type": "uint256"
                  }
                ],
                "name": "rescueFunds",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                  }
                ],
                "name": "routers",
                "outputs": [
                  {
                    "internalType": "bool",
                    "name": "isSupported",
                    "type": "bool"
                  },
                  {
                    "internalType": "uint256",
                    "name": "fees",
                    "type": "uint256"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "router",
                        "type": "address"
                      },
                      {
                        "internalType": "address",
                        "name": "token",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      },
                      {
                        "internalType": "bytes",
                        "name": "permit",
                        "type": "bytes"
                      },
                      {
                        "internalType": "address[]",
                        "name": "tokenAToPath",
                        "type": "address[]"
                      },
                      {
                        "internalType": "address[]",
                        "name": "tokenBToPath",
                        "type": "address[]"
                      }
                    ],
                    "internalType": "struct LpSwapDetails[]",
                    "name": "lpSwapDetails_",
                    "type": "tuple[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "router",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "sizeBps",
                        "type": "uint256"
                      },
                      {
                        "internalType": "uint256",
                        "name": "minOutputAmount",
                        "type": "uint256"
                      },
                      {
                        "internalType": "address[]",
                        "name": "nativeToOutputPath",
                        "type": "address[]"
                      }
                    ],
                    "internalType": "struct WethSwapDetails[]",
                    "name": "wEthSwapDetails_",
                    "type": "tuple[]"
                  },
                  {
                    "internalType": "address",
                    "name": "recipient_",
                    "type": "address"
                  },
                  {
                    "internalType": "uint256",
                    "name": "nftId_",
                    "type": "uint256"
                  }
                ],
                "name": "swapLpToTokens",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "components": [
                      {
                        "internalType": "contract IAggregationExecutor",
                        "name": "executor",
                        "type": "address"
                      },
                      {
                        "components": [
                          {
                            "internalType": "contract IERC20",
                            "name": "srcToken",
                            "type": "address"
                          },
                          {
                            "internalType": "contract IERC20",
                            "name": "dstToken",
                            "type": "address"
                          },
                          {
                            "internalType": "address payable",
                            "name": "srcReceiver",
                            "type": "address"
                          },
                          {
                            "internalType": "address payable",
                            "name": "dstReceiver",
                            "type": "address"
                          },
                          {
                            "internalType": "uint256",
                            "name": "amount",
                            "type": "uint256"
                          },
                          {
                            "internalType": "uint256",
                            "name": "minReturnAmount",
                            "type": "uint256"
                          },
                          {
                            "internalType": "uint256",
                            "name": "flags",
                            "type": "uint256"
                          },
                          {
                            "internalType": "bytes",
                            "name": "permit",
                            "type": "bytes"
                          }
                        ],
                        "internalType": "struct SwapDescription",
                        "name": "desc",
                        "type": "tuple"
                      },
                      {
                        "internalType": "bytes",
                        "name": "routeData",
                        "type": "bytes"
                      },
                      {
                        "internalType": "bytes",
                        "name": "permit",
                        "type": "bytes"
                      },
                      {
                        "internalType": "uint256",
                        "name": "minReturnAmount",
                        "type": "uint256"
                      }
                    ],
                    "internalType": "struct ErcSwapDetails[]",
                    "name": "data_",
                    "type": "tuple[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "router",
                        "type": "address"
                      },
                      {
                        "internalType": "address",
                        "name": "token",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "amount",
                        "type": "uint256"
                      },
                      {
                        "internalType": "bytes",
                        "name": "permit",
                        "type": "bytes"
                      },
                      {
                        "internalType": "address[]",
                        "name": "tokenAToPath",
                        "type": "address[]"
                      },
                      {
                        "internalType": "address[]",
                        "name": "tokenBToPath",
                        "type": "address[]"
                      }
                    ],
                    "internalType": "struct LpSwapDetails[]",
                    "name": "lpSwapDetails_",
                    "type": "tuple[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "address",
                        "name": "router",
                        "type": "address"
                      },
                      {
                        "internalType": "address",
                        "name": "lpToken",
                        "type": "address"
                      },
                      {
                        "internalType": "uint256",
                        "name": "minReturnAmount",
                        "type": "uint256"
                      },
                      {
                        "internalType": "address[]",
                        "name": "nativeToToken0",
                        "type": "address[]"
                      },
                      {
                        "internalType": "address[]",
                        "name": "nativeToToken1",
                        "type": "address[]"
                      }
                    ],
                    "internalType": "struct OutputLp",
                    "name": "outputLpDetails_",
                    "type": "tuple"
                  },
                  {
                    "internalType": "address",
                    "name": "recipient_",
                    "type": "address"
                  },
                  {
                    "internalType": "uint256",
                    "name": "nftId_",
                    "type": "uint256"
                  }
                ],
                "name": "swapTokenToLp",
                "outputs": [],
                "stateMutability": "payable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "components": [
                      {
                        "internalType": "contract IAggregationExecutor",
                        "name": "executor",
                        "type": "address"
                      },
                      {
                        "components": [
                          {
                            "internalType": "contract IERC20",
                            "name": "srcToken",
                            "type": "address"
                          },
                          {
                            "internalType": "contract IERC20",
                            "name": "dstToken",
                            "type": "address"
                          },
                          {
                            "internalType": "address payable",
                            "name": "srcReceiver",
                            "type": "address"
                          },
                          {
                            "internalType": "address payable",
                            "name": "dstReceiver",
                            "type": "address"
                          },
                          {
                            "internalType": "uint256",
                            "name": "amount",
                            "type": "uint256"
                          },
                          {
                            "internalType": "uint256",
                            "name": "minReturnAmount",
                            "type": "uint256"
                          },
                          {
                            "internalType": "uint256",
                            "name": "flags",
                            "type": "uint256"
                          },
                          {
                            "internalType": "bytes",
                            "name": "permit",
                            "type": "bytes"
                          }
                        ],
                        "internalType": "struct SwapDescription",
                        "name": "desc",
                        "type": "tuple"
                      },
                      {
                        "internalType": "bytes",
                        "name": "routeData",
                        "type": "bytes"
                      },
                      {
                        "internalType": "bytes",
                        "name": "permit",
                        "type": "bytes"
                      },
                      {
                        "internalType": "uint256",
                        "name": "minReturnAmount",
                        "type": "uint256"
                      }
                    ],
                    "internalType": "struct ErcSwapDetails[]",
                    "name": "data_",
                    "type": "tuple[]"
                  },
                  {
                    "internalType": "address",
                    "name": "recipient_",
                    "type": "address"
                  },
                  {
                    "internalType": "uint256",
                    "name": "nftId_",
                    "type": "uint256"
                  }
                ],
                "name": "swapTokensToTokens",
                "outputs": [],
                "stateMutability": "payable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "address",
                    "name": "newOwner",
                    "type": "address"
                  }
                ],
                "name": "transferOwnership",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "enum FeeType[]",
                    "name": "feeTypes_",
                    "type": "uint8[]"
                  },
                  {
                    "internalType": "uint256[]",
                    "name": "fees_",
                    "type": "uint256[]"
                  }
                ],
                "name": "updateFee",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "address",
                    "name": "feeVault_",
                    "type": "address"
                  }
                ],
                "name": "updateFeeVault",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [
                  {
                    "internalType": "address[]",
                    "name": "routers_",
                    "type": "address[]"
                  },
                  {
                    "components": [
                      {
                        "internalType": "bool",
                        "name": "isSupported",
                        "type": "bool"
                      },
                      {
                        "internalType": "uint256",
                        "name": "fees",
                        "type": "uint256"
                      }
                    ],
                    "internalType": "struct Router[]",
                    "name": "details_",
                    "type": "tuple[]"
                  }
                ],
                "name": "updateRouters",
                "outputs": [],
                "stateMutability": "nonpayable",
                "type": "function"
              },
              {
                "inputs": [],
                "name": "wNative",
                "outputs": [
                  {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                  }
                ],
                "stateMutability": "view",
                "type": "function"
              },
              {
                "stateMutability": "payable",
                "type": "receive"
              }
            ]`
