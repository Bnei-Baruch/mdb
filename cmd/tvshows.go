package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/tvshows"
)

func init() {
	command := &cobra.Command{
		Use:   "tvshows",
		Short: "Import TV Shows to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			tvshows.ImportTVShows()
		},
	}
	RootCmd.AddCommand(command)
}
