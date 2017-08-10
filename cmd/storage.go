package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/storage"
)

func init() {
	command := &cobra.Command{
		Use:   "storage",
		Short: "Import storage locations status of files in MDB",
		Run: func(cmd *cobra.Command, args []string) {
			storage.ImportStorageStatus()
		},
	}
	RootCmd.AddCommand(command)
}
