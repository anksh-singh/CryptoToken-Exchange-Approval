package utils

import (
	"bridge-allowance/config"
	"encoding/hex"
	"github.com/onrik/ethrpc"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
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

// ConvertStringValueToBigFloat
// function to get string value of Token Value(Big Int) the given decimal value
func (h *Helpers) ConvertStringValueToBigFloat(value string, decimal string) (float64, error) {
	parsedValue, _ := new(big.Float).SetString(value)
	baseVal, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return -1, status.Errorf(codes.Internal, err.Error())
	}
	wei := math.Pow(10, baseVal)
	weiBigint := new(big.Int).SetInt64(int64(wei))
	valueNativeToken, _ := new(big.Float).Quo(parsedValue, new(big.Float).SetInt(weiBigint)).Float64()
	return valueNativeToken, err
}

func (h *Helpers) ConvertStringValueToFloatWei(value string, decimal string) (float64, error) {
	parsedValue, _ := new(big.Float).SetString(value)
	baseVal, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return -1, status.Errorf(codes.Internal, err.Error())
	}
	wei := math.Pow(10, baseVal)
	weiBigint := new(big.Int).SetInt64(int64(wei))
	valueNativeToken, _ := new(big.Float).Mul(parsedValue, new(big.Float).SetInt(weiBigint)).Float64()
	return valueNativeToken, err
}
func (h *Helpers) ConvertStringValueToIntWei(value string, decimal string) (int64, error) {
	parsedValue, _ := new(big.Float).SetString(value)
	baseVal, err := strconv.ParseFloat(decimal, 64)
	if err != nil {
		return -1, status.Errorf(codes.Internal, err.Error())
	}
	wei := math.Pow(10, baseVal)
	weiBigint := new(big.Int).SetInt64(int64(wei))
	valueNativeToken, _ := new(big.Float).Mul(parsedValue, new(big.Float).SetInt(weiBigint)).Int64()
	return valueNativeToken, err
}

// CalculateRateWithDecimal
// function to calculate Rate of token in base 10 with decimal value
func (h *Helpers) CalculateRateWithDecimal(value string, decimal int64) float64 {
	parsedValue, ok := new(big.Float).SetString(value)
	if !ok {
		return 0
	}
	baseVal := float64(decimal)
	wei := math.Pow(10, baseVal)
	weiBigint := new(big.Int).SetInt64(int64(wei))
	valueNativeToken, _ := new(big.Float).Quo(parsedValue, new(big.Float).SetInt(weiBigint)).Float64()
	return valueNativeToken
}

// ConvertHexToFloat64WithDecimals function to get float64 value of Token Value in base 16 with 18 decimals
func (h *Helpers) ConvertHexToFloat64WithDecimals(value string) float64 {
	valueBigint := new(big.Int)
	valueBigint.SetString(value, 16)
	wei := math.Pow(10, 18)
	weiBigint := new(big.Int).SetInt64(int64(wei))
	valueNativeToken, _ := new(big.Float).Quo(new(big.Float).SetInt(valueBigint), new(big.Float).SetInt(weiBigint)).Float64()
	return valueNativeToken
}

// ConvertHexToFloat64
// function to get float64 value of Token Value in base 16
func (h *Helpers) ConvertHexToFloat64(value string) float64 {
	valueBigint := new(big.Float)
	valueBigint.SetString(value)
	res, _ := valueBigint.Float64()
	return res
}

// ConvertHexToInt64
// function to get int64 value of Token Value in base 16
func (h *Helpers) ConvertHexToInt64(value string) int64 {
	valueBigint := new(big.Int)
	valueBigint.SetString(value[2:], 16)
	res := valueBigint.Int64()
	return res
}

// ConvertHexFloatString
// function to convert hex string
func (h *Helpers) ConvertHexFloatString(value string) string {
	valueBigint1 := new(big.Float)
	valueBigint1.SetString(value)
	res, _ := valueBigint1.Float64()
	resString := strconv.FormatFloat(res, 'f', -1, 64)
	return resString
}

// ConvertStringToFloat64 function that converts string to float64
func (h *Helpers) ConvertStringToFloat64(stringValue string) float64 {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0
	}
	return value
}

// ConvertStringIntToStringFloat64WithDecimals
// function to get float64 value of Token Value in base 10 with decimals
func (h *Helpers) ConvertStringIntToStringFloat64WithDecimals(value string, decimals int64) (string, error) {
	valueInt, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return "", err
	}
	evalVal := float64(valueInt) / math.Pow(10, float64(decimals))
	returnVal := strconv.FormatFloat(evalVal, 'f', -1, 64)
	return returnVal, err
}

// ConvertStringFloatToIntWithDecimals
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

func (h *Helpers) CheckSumAddress(address string) string {
	address = strings.Replace(strings.ToLower(address), "0x", "", 1)
	hash := sha3.NewLegacyKeccak256()
	_, _ = hash.Write([]byte(address))
	sum := hash.Sum(nil)
	digest := hex.EncodeToString(sum)

	b := strings.Builder{}
	b.WriteString("0x")

	for i := 0; i < len(address); i++ {
		a := address[i]
		if a > '9' {
			d, _ := strconv.ParseInt(digest[i:i+1], 16, 8)

			if d >= 8 {
				// Upper case it
				a -= 'a' - 'A'
				b.WriteByte(a)
			} else {
				// Keep it lower
				b.WriteByte(a)
			}
		} else {
			// Keep it lower
			b.WriteByte(a)
		}
	}
	return b.String()
}

// ConvertInterfaceToString function that converts interface to string
func (h *Helpers) ConvertInterfaceToString(value interface{}) string {
	valueString, _ := value.(string)
	return valueString
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
