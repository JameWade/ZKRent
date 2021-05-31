package main

import (
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/cmd/VehicleUser-client/commands"
	"os"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.GenKeyCmd,
		commands.SendtxCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
