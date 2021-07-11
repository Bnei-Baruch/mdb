package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/batch"
)

var httpToHttpsCmd = &cobra.Command{
	Use:   "http_to_https",
	Short: "Replace youtube http to https on blog posts",
	Run:   httpToHttpsFn,
}

func init() {
	batchCmd.AddCommand(httpToHttpsCmd)
}

func httpToHttpsFn(cmd *cobra.Command, args []string) {
	log.SetLevel(log.DebugLevel)
	batch.PostsHttpToHttps()
}