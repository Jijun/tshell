package cmd

import (
	"os"
	"tshell/util"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Used to delete an existing bucket",
	Long: `Used to delete an existing bucket

Format:
  ./tshell config delete -a <alias> [-c <config-file-path>]

Example:
  ./tshell config delete -a example`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteBucketConfig(cmd)
	},
}

func init() {
	configCmd.AddCommand(configDeleteCmd)

	configDeleteCmd.Flags().StringP("alias", "a", "", "Bucket alias")

	_ = configDeleteCmd.MarkFlagRequired("alias")
}

func deleteBucketConfig(cmd *cobra.Command) {
	alias, _ := cmd.Flags().GetString("alias")
	b, i, err := util.FindBucket(&config, alias)
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}

	if i < 0 {
		logger.Fatalln("Bucket not exist in config file!")
	}
	config.Buckets = append(config.Buckets[:i], config.Buckets[i+1:]...)

	viper.Set("tshell.buckets", config.Buckets)
	if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logger.Infof("Delete succeccfully! name: %s, endpoint: %s, alias: %s", b.Name, b.Endpoint, b.Alias)
}
