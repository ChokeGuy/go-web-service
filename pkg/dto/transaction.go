package dto

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type NativeTransferDto struct {
	TransactionType string  `json:"txType" validate:"required"`
	PrivateKey      string  `json:"privateKey" validate:"required"`
	From            string  `json:"from" validate:"required"`
	To              string  `json:"to" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	Decimal         int     `json:"decimal" validate:"required"`
}

type TransferTokenDto struct {
	TransactionType string  `json:"txType" validate:"required"`
	PrivateKey      string  `json:"privateKey" validate:"required"`
	From            string  `json:"from" validate:"required"`
	To              string  `json:"to" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	Token           string  `json:"token" validate:"required"`
	Decimal         int     `json:"decimal" validate:"required"`
}

type SmartContractCallDto struct {
	TransactionType string `json:"txType" validate:"required"`
	PrivateKey      string `json:"privateKey" validate:"required"`
	From            string `json:"from" validate:"required"`
	To              string `json:"to" validate:"required"`
	Amount          int64  `json:"amount" validate:"required"`
	GasLimit        int64  `json:"gasLimit" validate:"required"`
	GasPrice        int64  `json:"gasPrice" validate:"required"`
	MetaData        string `json:"data" validate:"required"`
}

type SendTransactionDto struct {
	Transaction *types.Transaction `json:"transaction" validate:"required"`
}
