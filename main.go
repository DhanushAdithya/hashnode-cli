package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DhanushAdithya/hashnode-cli/cmd"
	"github.com/spf13/viper"
)

func exit(msgs ...interface{}) {
	fmt.Println(msgs...)
	os.Exit(1)
}

func setupConfig() {
	viper.SetConfigName("hashnode")
	viper.SetConfigType("yaml")
	homeDir, _ := os.UserHomeDir()
	configFile := filepath.Join(homeDir, "hashnode.yaml")
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(configFile); err != nil {
				exit("Unable to create config file:", err)
			}
		}
	}
	viper.AddConfigPath(homeDir)
	fmt.Println("alo", viper.GetString("token"))
	if err := viper.ReadInConfig(); err != nil {
		exit("Unable to read config:", err)
	}
}

func main() {
	setupConfig()
	cmd.Execute()
}
