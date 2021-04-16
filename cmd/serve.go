package cmd

import (
	"github.com/bzp2010/webvpn/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run WebVPN server",
	Long:  `Run WebVPN server`,
	Run: func(cmd *cobra.Command, args []string) {
		core.NewServer()
	},
}