package ubeswap

import "time"

type ExchangeToken struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Version   Version   `json:"version"`
	LogoURI   string    `json:"logoURI"`
	Keywords  []string  `json:"keywords"`
	Tokens    []*Tokens `json:"tokens"`
}
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type Tokens struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	ChainID  int    `json:"chainId"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}
