package main

import (
	"fmt"
	"os"

	"github.com/hughescoin/basenames-cli/base"
	"github.com/hughescoin/basenames-cli/cmd"
)

func main() {

	err := base.InitClient()
	if err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		fmt.Println("Please ensure BASENAMES_RPC_URL and BASENAMES_PRIVATE_KEY environment variables are set.")
		os.Exit(1)
	}
	cmd.Execute()
}
