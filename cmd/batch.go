package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

func init() {
	command := &cobra.Command{
		Use:   "batch",
		Short: "run some batch command",
		Run: func(cmd *cobra.Command, args []string) {
			batch.RenameUnits()
			//batch.ReadRequestsLog()
		},
	}
	RootCmd.AddCommand(command)
}
