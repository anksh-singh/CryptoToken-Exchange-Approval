package rpc

const (
	SONAR_WATCH_DUMMY_RESPONSE = `{
	  "positions": [
		{
		  "chain": "solana",
		  "address": "tEsT1vjsJeKHw9GH5HpnQszn2LWmjR6q1AVCDCj51nd",
		  "protocols": [
			{
			  "protocol_id": "solend",
			  "name": "solend",
			  "site_url": "https://solend.fi/",
			  "logo_url": "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/SLNDpmoWTVADgEdndyvWzroNL7zSi1dF9PC3xHGtPwp/logo.png",
			  "is_actionable": false,
			  "portfolio": {
				"net_usd_value": "0.2857823093475726",
				"lending": [
				  {
					"is_actionable": false,
					"type": "lending",
					"asset_value": "0.004999999",
					"net_value": "0.2857823093475726",
					"tokens_supplied": [
					  {
						"token_address": "SoLEao8wTzSfqhuou8rcYsVoLjthVmiXuEjzdNPMnCz",
						"token_name": "",
						"token_symbol": "sbrLP",
						"token_decimals": 9,
						"logo_url": "",
						"balance": "0.004999999",
						"quote_rate": "14.1855014099523",
						"quote_price": ""
					  }
					],
					"underlying_token_list": [
					  {
						"token_address": "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So",
						"token_name": "",
						"token_symbol": "mSOL",
						"token_decimals": 9,
						"logo_url": "",
						"balance": "0.0008489197586675896",
						"quote_rate": "14.78",
						"quote_price": ""
					  },
					  {
						"token_address": "So11111111111111111111111111111111111111112",
						"token_name": "",
						"token_symbol": "SOL",
						"token_decimals": 9,
						"logo_url": "",
						"balance": "0.004270699272490446",
						"quote_rate": "13.67",
						"quote_price": ""
					  }
					]
				  }
				]
			  }
			},
			{
			  "protocol_id": "raydium",
			  "name": "raydium",
			  "site_url": "https://raydium.io/",
			  "logo_url": "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R/logo.png",
			  "is_actionable": false,
			  "portfolio": {
				"net_usd_value": "0.007",
				"liquidity_pool": [
				  {
					"is_actionable": false,
					"type": "liquidity_pool",
					"asset_value": "0.003426303116848952",
					"net_value": "0.007",
					"tokens_supplied": [
					  {
						"token_address": "J2hGHwbkpj2SVo6Bs4X2Houy7n6oauydhbh9D6HpKBU4",
						"token_name": "",
						"token_symbol": "rayLP",
						"token_decimals": 9,
						"logo_url": "",
						"balance": "0.003426303116848952",
						"quote_rate": "0.48947187383556456",
						"quote_price": ""
					  }
					],
					"underlying_token_list": [
					  {
						"token_address": "GePFQaZKHcWE5vpxHfviQtH5jgxokSs51Y5Q4zgBiMDs",
						"token_name": "",
						"token_symbol": "JFI",
						"token_decimals": 9,
						"logo_url": "",
						"balance": "0.05998711191940201",
						"quote_rate": "0.028558660412360683",
						"quote_price": ""
					  },
					  {
						"token_address": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
						"token_name": "",
						"token_symbol": "USDC",
						"token_decimals": 6,
						"logo_url": "",
						"balance": "0.001713151558424476",
						"quote_rate": "1",
						"quote_price": ""
					  }
					]
				  }
				]
			  }
			}
		  ]
		}
	  ]
	}`
)
