package constants

// TransactionType represents the type of transaction
type transactionType struct {
	NativeTransfer    string
	TransferToken     string
	SmartContractCall string
}

// TransactionType represents the type of transaction
var TransactionType = &transactionType{
	NativeTransfer:    "native_transfer",
	TransferToken:     "transfer_token",
	SmartContractCall: "smart_contract_call",
}
