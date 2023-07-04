package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/keycloak"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exportKCCmd = &cobra.Command{
	Use:   "kc",
	Short: "Fill users table with kc account id",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(keycloak.ExportKC).Run()
	},
}

func init() {
	RootCmd.AddCommand(exportKCCmd)
}
