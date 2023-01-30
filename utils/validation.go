package utils

import (
	// "bridge-allowance/pkg/grpc/proto/pb"
	"encoding/hex"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func (u *UtilConf) ValidateAddress(address string, chainGroup string, chain string) (bool, string, error) {
	switch chainGroup {
	case "nonevm":
		if chain == "solana" {
			pubKey, err := solana.PublicKeyFromBase58(address)
			if err != nil {
				return false, address, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v is not nonevm address", address), "Invalid Solana Address")
			}
			isValid := solana.IsOnCurve(pubKey.Bytes())
			if isValid {
				return isValid, address, nil
			} else {
				return isValid, address, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v is not nonevm address", address), "Invalid Solana Address")

			}
		} else {
			return true, address, nil
		}

	case "evm":
		address, _ = u.ResolveENSAddress(address)
		address = u.ResolveXDCAddress(address)
		address = u.ResolveBech32Address(address)
		address = u.ResolveIoAddress(address)
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
		return true, b.String(), nil
	case "cosmos_network":
		//TODO: address valdation for cosmos chains
		return true, address, nil
	default:
		u.log.Info("Unsupported chain: ", chainGroup)
		return false, "", status.Errorf(codes.InvalidArgument, "Unsupported chain", "Invalid Data")
	}
}
func (u *UtilConf) ValidateCosmosAddress(address string, chainName string) (bool, string, error) {
	walletInfo := u.GetCosmosWalletInfo(chainName)
	if !strings.HasPrefix(strings.ToLower(address), strings.ToLower(walletInfo.Prefix)) {
		return false, address, status.Errorf(codes.InvalidArgument, fmt.Sprintf("%v is not %v address", address, chainName), fmt.Sprintf("Invalid %v Address", chainName))
	}
	return true, address, nil
}


