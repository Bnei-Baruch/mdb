package cmd

import (
	"github.com/Bnei-Baruch/mdb/monitor"
	_ "github.com/Bnei-Baruch/mdb/monitor/plugins/inputs/all"
	_ "github.com/Bnei-Baruch/mdb/monitor/plugins/outputs/all"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor mdb flows, processes, data integrity, collects and analyses BB archive mdb data, checks integrity and alerts about abnormal behaviour.",
		Run: func(cmd *cobra.Command, args []string) {
			monitor.ProcessMonitoring()
		},
	}
	RootCmd.AddCommand(command)
}
