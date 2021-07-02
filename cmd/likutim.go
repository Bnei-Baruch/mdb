package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/likutim"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var likutimCmd = &cobra.Command{
	Use:   "likutim",
	Short: "by default run create units and create tar",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(likutim.CreateUnits).Run()
		new(likutim.CreateTar).Run()
	},
}

var tarCmd = &cobra.Command{
	Use:   "tar",
	Short: "Create tar folder name - unit uid type likutin with files of this unit",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(likutim.CreateTar).Run()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create units likutim by uniq files from kitvei makor",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(likutim.CreateUnits).Run()
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "compare all kitvei makor files and find doubles",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(likutim.Compare).Run()
	},
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Only print duplicate kitvei makor without any changes",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(likutim.PrintWithDoc).Run()
	},
}

func init() {
	RootCmd.AddCommand(likutimCmd)
	likutimCmd.AddCommand(tarCmd)
	likutimCmd.AddCommand(createCmd)
	likutimCmd.AddCommand(compareCmd)
	likutimCmd.AddCommand(printCmd)
}
