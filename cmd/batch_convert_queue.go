package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var convertQueueCmd = &cobra.Command{
	Use:   "convert_queue",
	Short: "Queue files for conversion",
	Run:   convertQueueFn,
}

func init() {
	batchCmd.AddCommand(convertQueueCmd)
}

func convertQueueFn(cmd *cobra.Command, args []string) {
	batch.QueueWork()
}
