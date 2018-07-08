package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/twitter"
)

func init() {
	command := &cobra.Command{
		Use:   "twitter-import",
		Short: "Import archived tweeter data",
		Run: func(cmd *cobra.Command, args []string) {
			//twitter.Analyze()
			twitter.ImportDumps()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "twitter-latest",
		Short: "Import latest tweets from Twitter for registered accounts",
		Run: func(cmd *cobra.Command, args []string) {
			twitter.ImportLatestTweets()
		},
	}
	RootCmd.AddCommand(command)
}
