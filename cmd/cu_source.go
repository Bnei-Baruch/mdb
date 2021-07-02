package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/cusource"
	log "github.com/Sirupsen/logrus"
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

var removeFilesByFileNameCmd = &cobra.Command{
	Use:   "remove_files_by_name",
	Short: "remove files by Fiel name ",
	Run:   removeFilesByFileNameCmdFn,
}

func init() {
	RootCmd.AddCommand(buildCUSourcesCmd)
	RootCmd.AddCommand(buildCUSourcesValidatorCmd)
	RootCmd.AddCommand(removeFilesByFileNameCmd)
}

func buildCUSourcesCmdFn(cmd *cobra.Command, args []string) {
	cusource.InitBuildCUSources()
}

func buildCUSourcesValidatorCmdFn(cmd *cobra.Command, args []string) {
	new(cusource.ComparatorDbVsFolder).Run()
}

func removeFilesByFileNameCmdFn(cmd *cobra.Command, args []string) {
	log.SetLevel(log.InfoLevel)
	cusource.RemoveFilesByFileName()
}
