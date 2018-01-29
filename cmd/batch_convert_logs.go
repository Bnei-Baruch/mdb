package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var convertLogsCmd = &cobra.Command{
	Use:   "convert_logs",
	Short: "Parse conversion logs",
	Run:   convertLogsFn,
}

func init() {
	batchCmd.AddCommand(convertLogsCmd)
}

func convertLogsFn(cmd *cobra.Command, args []string) {
	batch.ReadLogs()
}
