package cmd

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/tui"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a post",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckToken()
		if viper.GetString("file") == "" {
			utils.Exit("Please specify a file to publish (--file <path-to-file>)")
		}
		tui.Publish()
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringP("file", "f", "", "File to publish")
	publishCmd.Flags().StringP("title", "t", "", "Title of the post")
	publishCmd.Flags().StringSliceP("tags", "T", []string{}, "Tags")
	publishCmd.Flags().StringP("cover-image", "c", "", "Cover image")

	viper.BindPFlag("file", publishCmd.Flags().Lookup("file"))
	viper.BindPFlag("title", publishCmd.Flags().Lookup("title"))
	viper.BindPFlag("tags", publishCmd.Flags().Lookup("tags"))
	viper.BindPFlag("cover-image", publishCmd.Flags().Lookup("cover-image"))

	viper.SetDefault("file", "")
	viper.SetDefault("title", "")
	viper.SetDefault("tags", []string{})
	viper.SetDefault("cover-image", "")
}
