package cmd

import (
	"github.com/Bnei-Baruch/mdb/importer/str"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var srtCmd = &cobra.Command{
	Use:   "srt",
	Short: "build srt from mp3",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		lang := "he"
		cts := "20, 29"
		for _, arg := range args {
			_arg := strings.Split(arg, "=")
			switch _arg[0] {
			case "cts":
				cts = _arg[1]
				break
			case "lang":
				lang = _arg[1]
				break
			}
		}

		new(str.BuildSrt).Run(cts, lang)
	},
}

func init() {
	RootCmd.AddCommand(srtCmd)
}
