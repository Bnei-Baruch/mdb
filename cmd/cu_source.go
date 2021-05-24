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

var removeFilesBySHA1Cmd = &cobra.Command{
	Use:   "remove_files_by_sha",
	Short: "remove files by SHA1 ",
	Run:   removeFilesBySHA1CmdFn,
}

func init() {
	RootCmd.AddCommand(buildCUSourcesCmd)
	RootCmd.AddCommand(buildCUSourcesValidatorCmd)
	RootCmd.AddCommand(removeFilesBySHA1Cmd)
}

func buildCUSourcesCmdFn(cmd *cobra.Command, args []string) {
	cusource.InitBuildCUSources()
}

func buildCUSourcesValidatorCmdFn(cmd *cobra.Command, args []string) {
	new(cusource.ComparatorDbVsFolder).Run()
}

func removeFilesBySHA1CmdFn(cmd *cobra.Command, args []string) {
	cusource.RemoveFilesBySHA1()
}
