package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/ffprobe"
)

func init() {
	command := &cobra.Command{
		Use:   "ffprobe-analyze",
		Short: "Analyze bulk ffprobe data to be imported",
		Run: func(cmd *cobra.Command, args []string) {
			ffprobe.Analyze()
		},
	}
	RootCmd.AddCommand(command)
}
