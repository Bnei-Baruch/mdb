package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/lessons_series"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lessonsSeriesCmd = &cobra.Command{
	Use:   "lessons_series",
	Short: "build lessons series collections by sources/likutims",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(lessons_series.LessonsSeries).Run()
	},
}

var insertLikutPropCmd = &cobra.Command{
	Use:   "insert_likut",
	Short: "add prop likut for lesson series collection",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		new(lessons_series.LessonsSeries).RunAddLikutProp()
	},
}

func init() {
	RootCmd.AddCommand(lessonsSeriesCmd)
	RootCmd.AddCommand(insertLikutPropCmd)
}
