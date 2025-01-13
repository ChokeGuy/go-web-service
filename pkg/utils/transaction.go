package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

func GetAmountInWei(decimal int, amount float64) *big.Int {
	decimalFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)

	amountFloat := new(big.Float).SetFloat64(amount)
	amountWei := new(big.Float).Mul(amountFloat, new(big.Float).SetInt(decimalFactor))

	amountIn := new(big.Int)
	amountWei.Int(amountIn)

	return amountIn
}

func GetMethodID(signature []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(signature)
	methodID := hash.Sum(nil)[:4]

	return methodID
}

// EncodeParameters encodes the given parameters for a smart contract call
func EncodeParameters(params []interface{}) []byte {
	var buf bytes.Buffer
	for _, param := range params {
		switch v := param.(type) {
		case string:
			buf.Write(common.LeftPadBytes([]byte(v), 32))
		case int:
			buf.Write(common.LeftPadBytes(big.NewInt(int64(v)).Bytes(), 32))
		case *big.Int:
			buf.Write(common.LeftPadBytes(v.Bytes(), 32))
		default:
			fmt.Printf("Unsupported parameter type: %v\n", reflect.TypeOf(param))
		}
	}
	return buf.Bytes()
}
