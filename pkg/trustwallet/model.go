package trustwallet

type ExchangeToken struct {
	Name        string `json:"name"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Explorer    string `json:"explorer"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	Status      string `json:"status"`
	ID          string `json:"id"`
	Links       []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"links"`
	Tags []string `json:"tags"`
}
