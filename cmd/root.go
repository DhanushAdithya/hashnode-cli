package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashnode-cli",
	Short: "Hashnode in the command line",
	Long: `Hashnode CLI lets you read, post and surf articles seamlessly right
from the command line`,
}

func exit(msgs ...interface{}) {
	fmt.Println(msgs...)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exit("Unable to execute command:", err)
	}
}

func init() {}
