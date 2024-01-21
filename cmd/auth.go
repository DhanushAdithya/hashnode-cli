package cmd

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add Personal Access Token to the config file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		utils.SetToken(token)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
