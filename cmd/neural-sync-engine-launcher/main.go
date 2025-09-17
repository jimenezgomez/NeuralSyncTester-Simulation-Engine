package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "neural-sync",
	Short: "NeuralSyncTester CLI",
}

func main() {
	Execute() // cobra root
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cliCmd)
	rootCmd.AddCommand(serverCmd)
}
