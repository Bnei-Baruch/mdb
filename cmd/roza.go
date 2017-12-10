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

	command = &cobra.Command{
		Use:   "roza-match-mdb",
		Short: "Match MDB content Unitsto directories in roza index",
		Run: func(cmd *cobra.Command, args []string) {
			roza.MatchDirectories()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "roza-master",
		Short: "Master merge of mdb, kmedia & roza",
		Run: func(cmd *cobra.Command, args []string) {
			roza.MatchFiles()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "roza-upload",
		Short: "Prepare files for upload",
		Run: func(cmd *cobra.Command, args []string) {
			roza.PrepareUpoad()
		},
	}
	RootCmd.AddCommand(command)
}
