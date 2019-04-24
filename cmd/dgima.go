package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/dgima"
)

func init() {
	command := &cobra.Command{
		Use:   "dgima-import",
		Short: "Import capture labels data",
		Run: func(cmd *cobra.Command, args []string) {
			dgima.Import()
		},
	}
	RootCmd.AddCommand(command)
}
