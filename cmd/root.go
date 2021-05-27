package cmd

import (
	"os"

	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var rootCmd = &cobra.Command{
	Use:   "webvpn",
	Short: "WebVPN is a zero-trust gateway for proxy application",
	Long:  `A flexible and configurable zero-trust gateway, it provides pluggable authentication and authorization for applications.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		g.Log().Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "config file")
}

func initConfig() {
	// set config file path
	viper.SetConfigFile(configFile)

	// reading environments
	viper.AutomaticEnv()

	// reading config file
	if err := viper.ReadInConfig(); err != nil {
		g.Log().Error(err)
	}
}
