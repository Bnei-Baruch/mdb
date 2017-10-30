package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var batchGorCmd = &cobra.Command{
	Use:   "gor",
	Short: "analyze requests.log (gor)",
	Run:   batchGorFn,
}

func init() {
	batchCmd.AddCommand(batchGorCmd)
}

func batchGorFn(cmd *cobra.Command, args []string) {
	batch.ReadRequestsLog()
}
