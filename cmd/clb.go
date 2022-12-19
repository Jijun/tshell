package cmd

import (
	"github.com/spf13/cobra"
)

var clbCmd = &cobra.Command{
	Use:   "clb",
	Short: "负载均衡",
	Long:  "负载均衡",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(clbCmd)
}
