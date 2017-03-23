package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/sources"
)

func init() {
	command := &cobra.Command{
		Use:   "sources",
		Short: "Import study materials sources to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			sources.ImportSources()
		},
	}
	RootCmd.AddCommand(command)
}
