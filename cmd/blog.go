package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb/importer/blog"
)

func init() {
	command := &cobra.Command{
		Use:   "blog-download",
		Short: "Download blog data via wordpress API",
		Run: func(cmd *cobra.Command, args []string) {
			blog.Download()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "blog-analyze",
		Short: "Analyze downloaded blog data",
		Run: func(cmd *cobra.Command, args []string) {
			blog.Analyze()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "blog-import",
		Short: "Import blog data into MDB",
		Run: func(cmd *cobra.Command, args []string) {
			blog.Import()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "blog-latest",
		Short: "Import latest blog posts into MDB",
		Run: func(cmd *cobra.Command, args []string) {
			blog.ImportLatest()
		},
	}
	RootCmd.AddCommand(command)
}
