package core

type BluzelleAssetInfo struct {
	CoingeckoId string
	Symbol      string
	LogoUrl     string
	Description string
}

var blzInfo = BluzelleAssetInfo{
	CoingeckoId: "bluzelle",
	Symbol:      "BLZ",
	LogoUrl:     "https://tokens.dharma.io/assets/0x5732046a883704404f284ce41ffadd5b007fd668/icon.png",
	Description: "The Native Token of Bluzelle network",
}
