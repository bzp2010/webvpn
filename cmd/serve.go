package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/core"
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
	// initialize database server
	_, err := core.Database()
	if err != nil {
		core.Log().Errorf("database initialize failed: %s", err.Error())
	}

	// initialize webvpn server
	s, err := core.Server(&core.Options{
		Public: isStartPublic,
		Admin:  isStartAdmin,
	})
	if err != nil {
		core.Log().Errorf("webvpn server initialize failed: %s", err.Error())
		return
	}

	// start webvpn server
	err = s.Start()
	if err != nil {
		core.Log().Errorf("webvpn server start failed: %s", err.Error())
		return
	}
}