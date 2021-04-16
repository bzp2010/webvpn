package cmd

import (
	"github.com/bzp2010/webvpn/internal/core"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"os"
)

var configFile string
var rootCmd = &cobra.Command{
	Use:   "webvpn",
	Short: "WebVPN is a zero-trust gateway for proxy application",
	Long:  `A flexible and configurable zero-trust gateway, it provides pluggable authentication and authorization for applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.NewServer()
	},
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
	g.Config().SetFileName(configFile)
}
