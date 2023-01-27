package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCheckExists(t *testing.T) {
	h := Helpers{}
	N4 := "198.02566690042823231"
	a, _ := h.ConvertStringValueToBigFloat(N4, "18")
	assert.Equal(t, 0.00000000000000019802566690042824, a)
	priceBigFloat, _ := new(big.Float).SetString("198.02566690042823231")
	guaranteedPriceBigFloat, _ := new(big.Float).SetString("191.88959851561835913")
	priceDifferenceBTWGuarenteedPRice := new(big.Float).Sub(priceBigFloat, guaranteedPriceBigFloat)
	priceDifferenceBTWGuarenteedPRice = new(big.Float).Quo(priceDifferenceBTWGuarenteedPRice, priceBigFloat)
	priceDifferenceBTWGuarenteedPRices, _ := new(big.Float).Mul(priceDifferenceBTWGuarenteedPRice, big.NewFloat(100)).Float64()
	l := fmt.Sprint(priceDifferenceBTWGuarenteedPRices)
	assert.Equal(t, "3.0986227597936717", l)
}
