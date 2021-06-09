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

var createSourceFromKiteiMakorCmd = &cobra.Command{
	Use:   "kitei_makor_to_source",
	Short: "create Source From Kitei Makor",
	Run:   createSourceFromKiteiMakorCmdFn,
}

func init() {
	RootCmd.AddCommand(buildCUSourcesCmd)
	RootCmd.AddCommand(buildCUSourcesValidatorCmd)
	RootCmd.AddCommand(removeFilesByFileNameCmd)
	RootCmd.AddCommand(createSourceFromKiteiMakorCmd)
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

func createSourceFromKiteiMakorCmdFn(cmd *cobra.Command, args []string) {
	log.SetLevel(log.DebugLevel)
	//executor := new(cusource.KiteiMakorPrintWithDoc)
	executor := new(cusource.KiteiMakorCompare)
	executor.Run()
}
