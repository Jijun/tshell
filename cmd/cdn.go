package cmd

import (
	"github.com/spf13/cobra"
)

var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "cdn toolkit",
	Long:  "cdn query toolkit",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cdnCmd)
}
