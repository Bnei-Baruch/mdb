package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var convertPrepareCmd = &cobra.Command{
	Use:   "convert_prepare",
	Short: "Prepare db table queue of files to convert",
	Run:   convertPrepareFn,
}

func init() {
	batchCmd.AddCommand(convertPrepareCmd)
}

func convertPrepareFn(cmd *cobra.Command, args []string) {
	batch.PrepareFilesForConvert()
}
