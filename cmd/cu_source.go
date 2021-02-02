package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/cusource"
	"github.com/spf13/cobra"
)

var buildCUSourcesCmd = &cobra.Command{
	Use:   "cu_sources",
	Short: "Build audio_sources units",
	Run:   buildCUSourcesCmdFn,
}

func init() {
	RootCmd.AddCommand(buildCUSourcesCmd)
}

func buildCUSourcesCmdFn(cmd *cobra.Command, args []string) {
	cusource.InitBuildCUSources()
}
