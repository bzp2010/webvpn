package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bzp2010/webvpn/internal/core"
	"github.com/bzp2010/webvpn/internal/model"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database scheme",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := core.Database()
		if err != nil {
			core.Log().Errorf("database initialize failed: %s", err.Error())
			return
		}

		core.Log().Info("connecting to database")

		err = db.AutoMigrate(&model.Service{})
		if err != nil {
			core.Log().Errorf("database migrate failed: %s", err.Error())
			return
		}
	},
}