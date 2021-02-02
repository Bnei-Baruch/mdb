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

var buildCUSourcesValidatorCmd = &cobra.Command{
	Use:   "cu_sources_validator",
	Short: "Validate script audio_sources units",
	Run:   buildCUSourcesValidatorCmdFn,
}

func init() {
	RootCmd.AddCommand(buildCUSourcesCmd)
	RootCmd.AddCommand(buildCUSourcesValidatorCmd)
}

func buildCUSourcesCmdFn(cmd *cobra.Command, args []string) {
	cusource.InitBuildCUSources()
}

func buildCUSourcesValidatorCmdFn(cmd *cobra.Command, args []string) {
	cusource.Validator()
}
