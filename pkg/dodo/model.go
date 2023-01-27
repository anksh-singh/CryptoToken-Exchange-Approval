package dodo

type ExchangeModel struct {
	Status int  `json:"status"`
	Data   Data `json:"data"`
}

type Data struct {
	ResAmount            float64 `json:"resAmount"`
	ResPricePerToToken   float64 `json:"resPricePerToToken"`
	ResPricePerFromToken float64 `json:"resPricePerFromToken"`
	PriceImpact          float64 `json:"priceImpact"`
	UseSource            string  `json:"useSource"`
	TargetDecimals       int     `json:"targetDecimals"`
	TargetApproveAddr    string  `json:"targetApproveAddr"`
	To                   string  `json:"to"`
	Data                 string  `json:"data"`
	RouteData            string  `json:"routeData"`
	MsgError             string  `json:"msgError"`
}
