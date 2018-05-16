package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/exporter/tagging"
)

var exportTaggingTVshowsCmd = &cobra.Command{
	Use:   "tagging_tvshows",
	Short: "Export list of VIDEO_PROGRAM collections with their units",
	Run:   exportTaggingTVshowsFn,
}

func init() {
	exportCmd.AddCommand(exportTaggingTVshowsCmd)
}

func exportTaggingTVshowsFn(cmd *cobra.Command, args []string) {
	tagging.ExportTVShows()
}
