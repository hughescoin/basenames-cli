package base

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// WeiToEth converts wei to ether
func WeiToEth(wei *big.Int) string {
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return eth.Text('f', 18)
}

func GetAddressFromPrivateKey(privateKey string) (string, error) {
	// Convert the private key from hex to bytes
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Create an ECDSA private key from the bytes
	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create ECDSA private key: %v", err)
	}

	// Get the public key from the private key
	publicKey := ecdsaPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to get public key")
	}

	// Get the Ethereum address from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address.Hex(), nil
}
