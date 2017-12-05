package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var eventsSubcollectionsCmd = &cobra.Command{
	Use:   "events_subcollections",
	Short: "Organize events lesson parts in subcollections",
	Run:   eventsSubcollectionsFn,
}

func init() {
	batchCmd.AddCommand(eventsSubcollectionsCmd)
}

func eventsSubcollectionsFn(cmd *cobra.Command, args []string) {
	batch.EventsSubcollections()
}
