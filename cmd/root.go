package cmd

import (
	"fmt"
	"log"
	"os"
	_ "tshell/logger"
	"tshell/util"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var config util.Config
var param util.Param
var cmdCnt int //控制某些函数在一个命令中被调用的次数

var rootCmd = &cobra.Command{
	Use:   "tshell",
	Short: "Welcome to use tshell",
	Long:  "Welcome to use tshell!",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	Version: "v1-beta",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-path", "c", "", "config file path(default is $HOME/.tshell.yaml)")
	rootCmd.PersistentFlags().StringVarP(&param.SecretID, "secret-id", "i", "", "config secretId")
	rootCmd.PersistentFlags().StringVarP(&param.SecretKey, "secret-key", "k", "", "config secretKey")
	rootCmd.PersistentFlags().StringVarP(&param.SessionToken, "token", "", "", "config sessionToken")
	rootCmd.PersistentFlags().StringVarP(&param.Endpoint, "endpoint", "e", "", "config endpoint")
}

func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	viper.SetConfigType("yaml")
	if cfgFile != "" {
		if cfgFile[0] == '~' {
			cfgFile = home + cfgFile[1:]
		}
		viper.SetConfigFile(cfgFile)
	} else {
		_, err = os.Stat(home + "/.tshell.yaml")
		if os.IsNotExist(err) {
			log.Println("Welcome to tshell!\nWhen you use tshell for the first time, you need to input some necessary information to generate the default configuration file of tshell.")
			initConfigFile(false)
			cmdCnt++
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".tshell")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		if err := viper.UnmarshalKey("tshell", &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if config.Base.Protocol == "" {
			config.Base.Protocol = "https"
		}
		// 尝试解码secretId/secretKey/session, 能解开就是加密的，否则就不解
		secretKey, err := util.DecryptSecret(config.Base.SecretKey)
		if err == nil {
			config.Base.SecretKey = secretKey
		}
		secretId, err := util.DecryptSecret(config.Base.SecretID)
		if err == nil {
			config.Base.SecretID = secretId
		}
		sessionToken, err := util.DecryptSecret(config.Base.SessionToken)
		if err == nil {
			config.Base.SessionToken = sessionToken
		}

	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
