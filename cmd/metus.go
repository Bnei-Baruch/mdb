package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/metus"
)

func init() {
	command := &cobra.Command{
		Use:   "metus-analyze",
		Short: "Metus analyze data",
		Run: func(cmd *cobra.Command, args []string) {
			metus.Analyze()
		},
	}
	RootCmd.AddCommand(command)

}
