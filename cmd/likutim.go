package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/likutim"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var likutimCmd = &cobra.Command{
	Use:   "likutim",
	Short: "create CU type likutim with .doc from kitvei makor",
	Run:   likutimCmdFn,
}

func init() {
	RootCmd.AddCommand(likutimCmd)
}

func likutimCmdFn(cmd *cobra.Command, args []string) {
	log.SetLevel(log.DebugLevel)
	switch args[0] {
	case "c":
		new(likutim.CreateUnits).Run()
	case "comp":
		new(likutim.Compare).Run()
	case "p":
	default:
		new(likutim.PrintWithDoc).Run()
	}
}
