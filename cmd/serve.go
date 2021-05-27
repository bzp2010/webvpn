package cmd

import (
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/server"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

// root command of run server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run WebVPN server",
}

func init() {
	// run only public server
	serveCmd.AddCommand(&cobra.Command{
		Use: "public",
		Short: "Run public server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			err := server.NewServer(&server.Options{
				Public: true,
				Admin:  false,
			})
			if err != nil {
				g.Log().Error(err)
			}
		},
	})

	// run only admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "admin",
		Short: "Run admin server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			err := server.NewServer(&server.Options{
				Public: false,
				Admin:  true,
			})
			if err != nil {
				g.Log().Error(err)
			}
		},
	})

	// run public and admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "all",
		Short: "Run all server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			err := server.NewServer(&server.Options{
				Public: true,
				Admin:  true,
			})
			if err != nil {
				g.Log().Error(err)
			}
		},
	})
}