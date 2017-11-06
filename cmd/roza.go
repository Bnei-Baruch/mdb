package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/roza"
)

func init() {
	command := &cobra.Command{
		Use:   "roza-load",
		Short: "Import roza storage to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			roza.LoadIndex()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "roza-match",
		Short: "Match directories in roza index to MDB content Units",
		Run: func(cmd *cobra.Command, args []string) {
			roza.MatchUnits()
		},
	}
	RootCmd.AddCommand(command)
}
