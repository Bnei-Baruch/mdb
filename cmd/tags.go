package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/tags"
)

func init() {
	command := &cobra.Command{
		Use:   "tags",
		Short: "Import tags to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			tags.ImportTags()
		},
	}
	RootCmd.AddCommand(command)
}
