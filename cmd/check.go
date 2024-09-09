package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hughescoin/basenames-cli/base"
	"github.com/spf13/cobra"
)

var tokenId string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check basename availability, expiration, or balance",
}

var availabilityCmd = &cobra.Command{
	Use:   "availability",
	Short: "Check a basename's availability",
	Run: func(cmd *cobra.Command, args []string) {
		if tokenId == "" {
			fmt.Println("Error: tokenId is required for checking availability")
			return
		}
		tokenIdBig, success := new(big.Int).SetString(tokenId, 10)
		if !success {
			fmt.Println("Error: Invalid tokenId format")
			return
		}
		//initiatilize contract
		contract, err := base.BaseClient.NewBasenamesContract()
		if err != nil {
			fmt.Printf("Error creating contract instance: %v\n", err)
			return

		}

		data, err := contract.ABI.Pack("isAvailable", tokenIdBig)
		if err != nil {
			fmt.Printf("Error encoding function call: %v\n", err)
			return
		}

		result, err := base.BaseClient.ReadContract(contract.Address, data)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
			return
		}

		var availability bool
		err = contract.ABI.UnpackIntoInterface(&availability, "isAvailable", result)
		if err != nil {
			fmt.Printf("Error decoding result: %v\n", err)
			return
		}

		if availability == true {
			fmt.Printf("%s is available \n", tokenId)
		} else {
			fmt.Printf("%s is not available \n", tokenId)
		}

	},
}

var expirationCmd = &cobra.Command{
	Use:   "expiration",
	Short: "Check a basename's expiration",
	Run: func(cmd *cobra.Command, args []string) {
		if tokenId == "" {
			fmt.Println("Error: tokenId is required for checking expiration")
			return
		}
		fmt.Printf("Checking expiration for token ID: %s\n", tokenId)

		tokenIdBig, success := new(big.Int).SetString(tokenId, 10)
		if !success {
			fmt.Println("Error: Invalid tokenId format")
			return
		}

		contract, err := base.BaseClient.NewBasenamesContract()
		if err != nil {

		}

		data, err := contract.ABI.Pack("nameExpires", tokenIdBig)
		if err != nil {
			fmt.Printf("Error encoding function call: %v\n", err)
			return
		}

		result, err := base.BaseClient.ReadContract(contract.Address, data)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
			return
		}

		var epochTime big.Int
		unpackResult, err := contract.ABI.Unpack("nameExpires", result)
		if err != nil {
			fmt.Printf("Error decoding result: %v\n", err)
			return
		}

		epochTime = *unpackResult[0].(*big.Int)
		expirationTime := time.Unix(epochTime.Int64(), 0)
		fmt.Printf("Expiration time for token ID %s: %s\n", tokenId, expirationTime.Format(time.RFC3339))
		// TODO: Implement expiration check
	},
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check the balance of the current account",
	Run: func(cmd *cobra.Command, args []string) {
		if base.BaseClient == nil {
			fmt.Println("Error: Client not initialized. Please ensure environment variables are set.")
			return
		}
		accountBalance, err := base.BaseClient.GetBalance(base.BaseClient.Address)
		if err != nil {
			fmt.Printf("Error checking block number: %v\n", err)
			return
		}

		accountBalanceBigInt, success := new(big.Int).SetString(accountBalance, 10)
		if !success {
			fmt.Println("Error: Invalid account balance format")
			return
		}
		accountBalance = base.WeiToEth(accountBalanceBigInt)
		fmt.Printf("%s Account balance: %s ETH\n", base.BaseClient.Address, accountBalance)
	},
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Check the latest block number",
	Run: func(cmd *cobra.Command, args []string) {
		if base.BaseClient == nil {
			fmt.Println("Error: Client not initialized. Please ensure environment variables are set.")
			return
		}
		blockNumber, err := base.BaseClient.GetBlock()
		if err != nil {
			fmt.Printf("Error checking block number: %v\n", err)
			return
		}
		fmt.Printf("Latest block number: %s\n", blockNumber)
	},
}

var ownerCmd = &cobra.Command{
	Use:   "ownerOf",
	Short: "Check the owner of a basename",
	Run: func(cmd *cobra.Command, args []string) {
		if tokenId == "" {
			fmt.Println("Error: tokenId is required for checking owner")
			return
		}

		// Convert tokenId string to big.Int
		tokenIdBig, success := new(big.Int).SetString(tokenId, 10)
		if !success {
			fmt.Println("Error: Invalid tokenId format")
			return
		}

		contract, err := base.BaseClient.NewBasenamesContract()
		if err != nil {
			fmt.Printf("Error creating contract instance: %v\n", err)
			return
		}

		// Encode function call
		data, err := contract.ABI.Pack("ownerOf", tokenIdBig)
		if err != nil {
			fmt.Printf("Error encoding function call: %v\n", err)
			return
		}

		// Call the contract
		result, err := base.BaseClient.ReadContract(contract.Address, data)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
			return
		}

		// Decode the result
		var owner common.Address
		err = contract.ABI.UnpackIntoInterface(&owner, "ownerOf", result)
		if err != nil {
			fmt.Printf("Error decoding result: %v\n", err)
			return
		}

		fmt.Printf("Owner of token %s: %s\n", tokenId, owner.Hex())
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.AddCommand(availabilityCmd)
	checkCmd.AddCommand(expirationCmd)
	checkCmd.AddCommand(balanceCmd)
	checkCmd.AddCommand(blockCmd)
	checkCmd.AddCommand(ownerCmd)

	// Add tokenId flag to the check command, making it available to all subcommands
	checkCmd.PersistentFlags().StringVar(&tokenId, "tokenId", "", "Token ID to check")

	// Require tokenId flag for availability and expiration subcommands
	availabilityCmd.MarkFlagRequired("tokenId")
	expirationCmd.MarkFlagRequired("tokenId")

	ownerCmd.Flags().StringVar(&tokenId, "tokenId", "", "Token ID to check owner")
	ownerCmd.MarkFlagRequired("tokenId")
}
