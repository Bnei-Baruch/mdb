package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/kabcoil"
)

func init() {
	command := &cobra.Command{
		Use:   "kabcoil-titles",
		Short: "Import kab.co.il titles",
		Run: func(cmd *cobra.Command, args []string) {
			kabcoil.ImportTitles()
		},
	}
	RootCmd.AddCommand(command)
}
