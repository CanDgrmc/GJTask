package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use: "gotask",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

}

func initConfig() {
	home, err := homedir.Expand(".")
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("config file:", viper.ConfigFileUsed())
	}
}
