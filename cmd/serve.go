package cmd

import (
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/server"
)

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
			startServer(true, false)
		},
	})

	// run only admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "admin",
		Short: "Run admin server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(false, true)
		},
	})

	// run public and admin server
	serveCmd.AddCommand(&cobra.Command{
		Use: "all",
		Short: "Run all server of WebVPN",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(true, true)
		},
	})

	rootCmd.AddCommand(serveCmd)
}

func startServer(isStartPublic, isStartAdmin bool)  {
	s, err := server.NewServer(&server.Options{
		Public: isStartPublic,
		Admin:  isStartAdmin,
	})

	if err != nil {
		g.Log().Errorf("webvpn server create failed: %s", err.Error())
		return
	}

	err = s.Start()
	if err != nil {
		g.Log().Errorf("webvpn server start failed: %s", err.Error())
		return
	}
}