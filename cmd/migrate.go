package cmd

import (
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/server"
	"github.com/bzp2010/webvpn/model"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database scheme",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.NewServer(&server.Options{
			Public: false,
			Admin:  false,
		})
		if err != nil {
			g.Log().Errorf("webvpn server create failed: %s", err.Error())
			return
		}

		err = s.DB.AutoMigrate(&model.Service{})
		if err != nil {
			g.Log().Errorf("database migrate failed: %s", err.Error())
			return
		}
	},
}