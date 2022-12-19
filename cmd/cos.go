package cmd

import (
	"github.com/spf13/cobra"
)

var cosCmd = &cobra.Command{
	Use:   "cos",
	Short: "对象存储",
	Long:  "cos operation",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cosCmd)
}
