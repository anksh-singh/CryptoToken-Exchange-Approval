package openApi

type GasPriceOpenAPIStruct struct {
	Level string  `json:"level"`
	Price float32 `json:"price"`
}

type GasPriceStruct struct {
	Slow   float32
	Fast   float32
	Normal float32
}
