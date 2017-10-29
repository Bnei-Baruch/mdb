package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var organizeKiteiMakorCmd = &cobra.Command{
	Use:   "kitei_makor",
	Short: "Organize kitei-makor files in derived units",
	Run:   organizeKiteiMakorFn,
}

func init() {
	batchCmd.AddCommand(organizeKiteiMakorCmd)
}

func organizeKiteiMakorFn(cmd *cobra.Command, args []string) {
	batch.OrganizeKiteiMakor()
}
