package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add Personal Access Token to the config file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		viper.Set("token", token)
		viper.WriteConfig()
		fmt.Println(
			lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("#39d927")).
				Render("Token set successfully"),
		)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
