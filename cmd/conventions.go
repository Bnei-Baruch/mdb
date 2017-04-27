package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/conventions"
)

func init() {
	command := &cobra.Command{
		Use:   "conventions",
		Short: "Insert convetions to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			conventions.ImportConvetions()
		},
	}
	RootCmd.AddCommand(command)
}
