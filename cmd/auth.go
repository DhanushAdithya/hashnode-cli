package cmd

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add Personal Access Token to the config file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		isSearchToken := viper.GetBool("search-token")
		utils.SetToken(token, isSearchToken)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().BoolP("search-token", "s", false, "Search token for searching articles")
	viper.BindPFlag("search-token", authCmd.Flags().Lookup("search-token"))
	viper.SetDefault("search-token", false)
}
