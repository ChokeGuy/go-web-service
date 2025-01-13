package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *EthClientConfig

// EthClientConfig stores the configuration of the Ethereum client
type EthClientConfig struct {
	RPCURL  string // Ethereum node RPC address
	ChainID *big.Int
	client  *ethclient.Client
	context context.Context
}

// NewEthClientConfig creates a new EthClientConfig
func NewEthClientConfig(rpcURL string) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("cannot connect to Ethereum node: %w", err)
	}

	ctx := context.Background()

	// Retrieve chain ID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("cannot retrieve Chain ID: %w", err)
	}

	EthClient = &EthClientConfig{
		RPCURL:  rpcURL,
		ChainID: chainID,
		client:  client,
		context: ctx,
	}

	return nil
}

// GetClient returns the Ethereum client
func (config *EthClientConfig) GetClient() *ethclient.Client {
	return config.client
}

// closes the connection and cancels the context
func (config *EthClientConfig) Close() {
	config.client.Close()
}

// get the current block number
func (config *EthClientConfig) GetBlockNumber() (uint64, error) {
	blockNumber, err := config.client.BlockNumber(config.context)
	if err != nil {
		log.Fatalf("Error while getting current block number: %v", err)
		return 0, err
	}
	fmt.Printf("Current block number: %d\n", blockNumber)

	return blockNumber, nil
}

// sign a transaction
func (config *EthClientConfig) SignTransaction(privateKeyHex string, tx *types.Transaction) (*types.Transaction, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)

	if err != nil {
		return nil, fmt.Errorf("error converting private key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(config.ChainID), privateKey)

	if err != nil {
		return nil, fmt.Errorf("error while signing transaction: %v", err)
	}
	return signedTx, nil
}

// sends a transaction
func (config *EthClientConfig) SendTransaction(tx *types.Transaction) error {
	err := config.client.SendTransaction(config.context, tx)
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}
	fmt.Printf("Transaction successfully sent, hash: %s\n", tx.Hash().Hex())
	return nil
}

// subscribe to logs (events)
func (config *EthClientConfig) SubscribeLogs(query ethereum.FilterQuery) {
	logs := make(chan types.Log)
	sub, err := config.client.SubscribeFilterLogs(config.context, query, logs)
	if err != nil {
		log.Fatalf("Error while subscribing to logs: %v", err)
	}
	defer sub.Unsubscribe()

	for log := range logs {
		fmt.Printf("Received log: %v\n", log)
	}
}
