package base

import (
	"fmt"
	"log"
	"os"
)

var BaseClient *Client

const (
	BASENAMES_RPC_URL     = "BASENAMES_RPC_URL"
	BASENAMES_PRIVATE_KEY = "BASENAMES_PRIVATE_KEY"
)

func ReadEnvCredentials() (*Credentials, error) {
	rpcURL := os.Getenv(BASENAMES_RPC_URL)
	privateKey := os.Getenv(BASENAMES_PRIVATE_KEY)

	if rpcURL == "" {
		return nil, fmt.Errorf("%s not set as an environment variable", BASENAMES_RPC_URL)
	}

	if privateKey == "" {
		return nil, fmt.Errorf("%s not set as an environment variable", BASENAMES_PRIVATE_KEY)
	}

	// Remove "0x" prefix from private key if present
	if len(privateKey) > 2 && privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}

	address, err := GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error deriving address from private key: %v", err)
	}

	return &Credentials{
		RpcUrl:     rpcURL,
		PrivateKey: privateKey,
		Address:    address,
	}, nil
}

func InitClient() error {
	creds, err := ReadEnvCredentials()
	if err != nil {
		log.Fatalf("error reading environmental variables: %s", err)
	}

	BaseClient = NewClient(creds.RpcUrl, creds.PrivateKey, creds.Address)
	fmt.Printf("Client initialized! \nHello: %s\n", creds.Address)
	return nil
}

func IsClientInitialized() bool {
	return BaseClient != nil
}
