package main

import (
	"log"

	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/diff"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/extract"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/history"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/server"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sibyl",
	Short: "sibyl cmd",
	Long:  `sibyl cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	Execute()
}

// now we build only one big fat execution
func init() {
	extractCmd := extract.NewExtractCmd()
	rootCmd.AddCommand(extractCmd)

	serverCmd := server.NewServerCmd()
	rootCmd.AddCommand(serverCmd)

	uploadCmd := upload.NewUploadCmd()
	rootCmd.AddCommand(uploadCmd)

	diffCmd := diff.NewDiffCommand()
	rootCmd.AddCommand(diffCmd)

	historyCmd := history.NewHistoryCmd()
	rootCmd.AddCommand(historyCmd)
}
