package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var ffprobeCmd = &cobra.Command{
	Use:   "ffprobe",
	Short: "Import ffprobe metadata",
	Run:   ffprobeFn,
}

func init() {
	batchCmd.AddCommand(ffprobeCmd)
}

func ffprobeFn(cmd *cobra.Command, args []string) {
	batch.ImportFFprobeMetadata()
}
