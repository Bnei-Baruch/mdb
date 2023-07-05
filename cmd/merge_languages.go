package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/mergelanguages"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mergeLanguagesCmd = &cobra.Command{
	Use:   "merge_languages",
	Short: "Change language of units and files",
	Run:   mergeLanguagesCmdFn,
}

func init() {
	RootCmd.AddCommand(mergeLanguagesCmd)
}

func mergeLanguagesCmdFn(cmd *cobra.Command, args []string) {
	log.SetLevel(log.DebugLevel)
	m := new(mergelanguages.MergeLanguages)
	if len(args) < 2 {
		println("not enough params")
		return
	}
	m.Init(args[0], args[1])
	m.Run()
}
