package cmd

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/tui"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckToken()
		tui.GetFeed()
	},
}

func init() {
	rootCmd.AddCommand(feedCmd)

	feedCmd.Flags().IntP("min-read", "m", 0, "Minimum read time")
	feedCmd.Flags().IntP("max-read", "M", 0, "Maximum read time")
	feedCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags")
	feedCmd.Flags().StringP("type", "T", "FEATURED", "Feed type")

	viper.BindPFlag("min-read", feedCmd.Flags().Lookup("min-read"))
	viper.BindPFlag("max-read", feedCmd.Flags().Lookup("max-read"))
	viper.BindPFlag("tags", feedCmd.Flags().Lookup("tags"))
	viper.BindPFlag("type", feedCmd.Flags().Lookup("type"))

	viper.SetDefault("min-read", 0)
	viper.SetDefault("max-read", 0)
	viper.SetDefault("tags", []string{})
	viper.SetDefault("type", "FEATURED")
}
