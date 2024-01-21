package cmd

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashnode-cli",
	Short: "Hashnode in the command line",
	Long: `Hashnode CLI lets you read, post and surf articles seamlessly right
from the command line`,
}

func Execute() {
	utils.SetupConfig()
	if err := rootCmd.Execute(); err != nil {
		utils.Exit("Unable to execute command")
	}
}
