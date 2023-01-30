package utils

import (
	"bridge-allowance/config"
	"github.com/onrik/ethrpc"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type Helpers struct {
	//ConvertHexToFloat64 float64
}

// ConvertStringToFloat64 function that converts string to float64
func (h *Helpers) ConvertStringToFloat64(stringValue string) float64 {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0
	}
	return value
}

// function to get Int64 value of  Value in base 10 with decimals
func (h *Helpers) ConvertStringFloatToIntWithDecimals(value string, base string) (int, error) {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	baseVal, err := strconv.ParseFloat(base, 64)
	if err != nil {
		return 0, err
	}
	dec := math.Pow(float64(10), baseVal)
	floatVal = (floatVal * dec) + 0.5 // go automatically stores floatval * dec values as XXX.99999....in variable , we need to add 0.5
	//When converting a floating-point number to an integer, the fraction is discarded (truncation towards zero).
	//it wont work for negative values
	return int(floatVal), err
}

// CheckTokenListData
// function to get key Value with a condition to be fulfilled from a slice of struct generic
func (h *Helpers) CheckTokenListData(itemTobeFound string, valuefromkey string, rval reflect.Value, rtype reflect.Type, keys ...string) (bool, string) {
	g := rval.Elem()
	returnval := false
	strvalue := ""
	if len(keys) > 0 {
		if rtype.Kind() == reflect.Slice {
			for i := 0; i < g.Len(); i++ {
				if len(keys) == 2 {
					if strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).FieldByName(keys[1]).String()) == itemTobeFound {
						returnval = true
						strvalue = strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).FieldByName(valuefromkey).String())
						break
					}
				} else {
					if strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(keys[0]).String()) == itemTobeFound {
						returnval = true
						strvalue = strings.ToLower(reflect.Indirect(g.Index(i)).FieldByName(valuefromkey).String())
						break
					}
				}

			}
		}
	}
	return returnval, strvalue
}

// WeiToGwei function that calculates reward per block
func (h *Helpers) WeiToGwei(value int64, fee string) string {
	decimal := float64(value)    // convert to float64
	pow := math.Pow(10, decimal) // convert to power of 10
	feeValue, err := strconv.ParseFloat(fee, 64)
	if err != nil {
		return ""
	}
	platformFee := feeValue / pow
	return strconv.FormatFloat(platformFee, 'f', -1, 64)
}

// TrimDollarSign function that trims dollar sign from string
func (h *Helpers) TrimDollarSign(s string) string {
	return strings.TrimLeft(s, "$")
}

// ConvertToWei convert the amount to wei for any data type and decimals
func (h *Helpers) ConvertToWei(interfaceAmount interface{}, decimals string) (*big.Int, error) {
	decimalInt, err := strconv.Atoi(decimals)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while converting decimals to int: %v", err)
	}
	amount := decimal.NewFromFloat(0)
	switch v := interfaceAmount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimalInt)))
	result := amount.Mul(mul)

	weiAmount := new(big.Int)
	weiAmount.SetString(result.String(), 10)

	return weiAmount, nil
}

// GetGasLimit calculate the gas limit for the transaction
func (h *Helpers) GetGasLimit(info config.Wallets, toAddress string, data string, value *big.Int, fromAddress string) (int, error) {
	client := ethrpc.New(info.RPC)
	// getting gas price
	gasLimit, err := client.EthEstimateGas(ethrpc.T{
		To:    toAddress,
		Data:  data,
		Value: value,
		From:  fromAddress,
	})
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}
