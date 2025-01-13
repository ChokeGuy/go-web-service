package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"web-service/pkg/constants"
	"web-service/pkg/dto"
	"web-service/pkg/eth"
	"web-service/pkg/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/gorilla/mux"
)

func getBlockNumber(w http.ResponseWriter, r *http.Request) utils.Response {
	ethClient := eth.EthClient

	blockNumber, err := ethClient.GetBlockNumber()
	if err != nil {
		return utils.BadRequestError(err.Error(), nil)
	}

	return utils.SuccessResponse("Get block number successfully", blockNumber)
}

// signTransaction handles signing transactions based on different types
func signTransaction(w http.ResponseWriter, r *http.Request) utils.Response {
	// Read body once
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.BadRequestError("Failed to read request body", nil)
	}

	// Parse transaction type first
	var txType struct {
		TransactionType string `json:"txType" validate:"required"`
	}
	if err := json.Unmarshal(body, &txType); err != nil {
		return utils.BadRequestError(fmt.Sprintf("Invalid transaction type: %s", err.Error()), nil)
	}

	// Handle different types of transactions based on 'txType'
	switch txType.TransactionType {
	case constants.TransactionType.NativeTransfer:
		var nativePayload dto.NativeTransferDto
		if err := json.Unmarshal(body, &nativePayload); err != nil {
			return utils.BadRequestError(fmt.Sprintf("Invalid native transfer payload: %s", err.Error()), nil)
		}
		return handleNativeTransfer(r, nativePayload)
	case constants.TransactionType.TransferToken:
		var tokenPayload dto.TransferTokenDto
		if err := json.Unmarshal(body, &tokenPayload); err != nil {
			return utils.BadRequestError(fmt.Sprintf("Invalid token transfer payload: %s", err.Error()), nil)
		}
		return handleTransferToken(r, tokenPayload)
	case constants.TransactionType.SmartContractCall:
		var contractPayload dto.SmartContractCallDto
		if err := json.Unmarshal(body, &contractPayload); err != nil {
			return utils.BadRequestError(fmt.Sprintf("Invalid contract call payload: %s", err.Error()), nil)
		}
		return handleSmartContractCall(r, contractPayload)
	default:
		return utils.BadRequestError("Unsupported transaction type", nil)
	}
}

// handleNativeTransfer handles native ETH transfer
func handleNativeTransfer(r *http.Request, payload dto.NativeTransferDto) utils.Response {
	// Convert amount to Wei
	amount := utils.GetAmountInWei(payload.Decimal, payload.Amount)
	return signAndRespond(r, payload.From, payload.To, amount, payload.PrivateKey, nil)
}

// handleTransferToken handles ERC20 token transfer
func handleTransferToken(r *http.Request, payload dto.TransferTokenDto) utils.Response {
	amount := utils.GetAmountInWei(payload.Decimal, payload.Amount)

	// Create signature for transfer(address,uint256) function
	transferFnSignature := []byte("transfer(address,uint256)")
	methodID := utils.GetMethodID(transferFnSignature)
	paddedToAddress := common.LeftPadBytes(common.HexToAddress(payload.To).Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var signature []byte
	signature = append(signature, methodID...)
	signature = append(signature, paddedToAddress...)
	signature = append(signature, paddedAmount...)

	return signAndRespond(r, payload.From, payload.Token, big.NewInt(0), payload.PrivateKey, signature)
}

// handleSmartContractCall handles smart contract function calls
func handleSmartContractCall(r *http.Request, payload dto.SmartContractCallDto) utils.Response {
	var metadata struct {
		Method string        `json:"method"` // Function name to call in the contract
		Params []interface{} `json:"params"` // Function parameters
	}
	if err := json.Unmarshal([]byte(payload.MetaData), &metadata); err != nil {
		return utils.BadRequestError("Invalid metadata for smart_contract_call", nil)
	}

	// Create signature for contract function
	methodSignature := []byte(metadata.Method)
	methodID := utils.GetMethodID(methodSignature)
	params := utils.EncodeParameters(metadata.Params) // Encode parameters

	var signature []byte

	signature = append(signature, methodID...)
	signature = append(signature, params...)

	return signAndRespond(r, payload.From, payload.To, big.NewInt(payload.Amount), payload.PrivateKey, signature)
}

func signAndRespond(r *http.Request, from, to string, amount *big.Int, privateKey string, signature []byte) utils.Response {
	ethClient := eth.EthClient
	client := ethClient.GetClient()

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	nonce, err := client.PendingNonceAt(r.Context(), fromAddress)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Nonce error : %s", err.Error()), nil)
	}

	gasPrice, err := client.SuggestGasPrice(r.Context())
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Gas Price error : %s", err.Error()), nil)
	}

	gasLimit, err := client.EstimateGas(r.Context(), ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		GasPrice: gasPrice,
		Data:     signature,
	})
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Gas Limit error : %s", err.Error()), nil)
	}

	tx := types.NewTx(&types.LegacyTx{
		To:       &toAddress,
		Value:    amount,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Data:     signature,
		Nonce:    nonce,
	})

	signedTx, err := ethClient.SignTransaction(privateKey, tx)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Sign Transaction error : %s", err.Error()), nil)
	}

	return utils.SuccessResponse("Transaction signed successfully", signedTx)
}

func sendTransaction(w http.ResponseWriter, r *http.Request) utils.Response {
	ethClient := eth.EthClient
	// client := ethClient.GetClient()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.BadRequestError("Failed to read request body", nil)
	}

	var payload dto.SendTransactionDto

	if err := json.Unmarshal(body, &payload); err != nil {
		return utils.BadRequestError(err.Error(), nil)
	}
	signedTx := payload.Transaction

	if err := ethClient.SendTransaction(signedTx); err != nil {
		return utils.BadRequestError(err.Error(), nil)
	}

	return utils.SuccessResponse("Transaction sent successfully", signedTx)
}

func SignTransactionRoutes(r *mux.Router) {
	r.HandleFunc("/block-number", utils.WrapHandler(getBlockNumber)).Methods(http.MethodGet)
	r.HandleFunc("/sign-transaction", utils.WrapHandler(signTransaction)).Methods(http.MethodPost)
	r.HandleFunc("/sign-transfer-token", utils.WrapHandler(signTransaction)).Methods(http.MethodPost)
	r.HandleFunc("/send-transaction", utils.WrapHandler(sendTransaction)).Methods(http.MethodPost)
}
