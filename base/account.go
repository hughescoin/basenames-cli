package base

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BasenamesContract struct {
	Address common.Address
	ABI     abi.ABI
	Client  *ethclient.Client
}

func (c *Client) NewBasenamesContract() (*BasenamesContract, error) {

	client, err := ethclient.Dial(c.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	basenamesABI, err := abi.JSON(strings.NewReader(BasenamesABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse basenames ABI: %v", err)
	}

	// Hardcode the Basenames contract address
	contractAddress := common.HexToAddress("0x03c4738Ee98aE44591e1A4A4F3CaB6641d95DD9a")

	return &BasenamesContract{
		Address: contractAddress,
		ABI:     basenamesABI,
		Client:  client,
	}, nil
}

func (c *Client) GetBlock() (string, error) {
	url := c.RpcURL

	// Create the JSON request payload
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	//set the headers of the client's http client
	c.setHeaders(req)

	// Send the request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print the response
	fmt.Println(string(body))
	return string(body), nil

}

// get the balance of the account
func (c *Client) GetBalance(address string) (string, error) {
	// Dial the Ethereum client
	client, err := ethclient.Dial(c.RpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Convert address string to common.Address
	account := common.HexToAddress(address)

	// Get the balance
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %v", err)
	}

	// Convert balance to string
	balanceStr := balance.String()

	fmt.Printf("Balance of %s: %s wei\n", address, balanceStr)
	return balanceStr, nil
}

func (c *Client) ReadContract(to common.Address, data []byte) ([]byte, error) {
	client, err := ethclient.Dial(c.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	return result, nil
}

func (c *Client) WriteContract(to common.Address, data []byte, value *big.Int) (string, error) {
	client, err := ethclient.Dial(c.RpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(c.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	if value == nil {
		value = big.NewInt(0) // Default value is zero
	}

	gasLimit := uint64(21000) // You might want to estimate this
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %v", err)
	}

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	ts := types.Transactions{signedTx}
	b := new(bytes.Buffer)
	ts.EncodeIndex(0, b)
	rawTxBytes := b.Bytes()
	rawTxHex := hex.EncodeToString(rawTxBytes)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", rawTxHex)
	return rawTxHex, nil
}
