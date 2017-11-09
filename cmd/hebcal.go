package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/hebcal"
)

func init() {
	command := &cobra.Command{
		Use:   "hebcal-load",
		Short: "Load hebcal data into memory",
		Run: func(cmd *cobra.Command, args []string) {
			h := new(hebcal.Hebcal)
			h.Load()
			h.Print()
		},
	}
	RootCmd.AddCommand(command)

}
