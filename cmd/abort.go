package cmd

import (
	"context"
	"tshell/util"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var abortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Abort parts",
	Long: `Abort parts

Format:
  ./tshell cos abort cos://<bucket-name>[/<prefix>] [flags]

Example:
  ./tshell cos abort cos://examplebucket/test/`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		include, _ := cmd.Flags().GetString("include")
		exclude, _ := cmd.Flags().GetString("exclude")

		abortParts(args[0], include, exclude)
	},
}

func init() {
	cosCmd.AddCommand(abortCmd)

	abortCmd.Flags().String("include", "", "List files that meet the specified criteria")
	abortCmd.Flags().String("exclude", "", "Exclude files that meet the specified criteria")
}

func abortParts(arg string, include string, exclude string) {
	bucketName, cosPath := util.ParsePath(arg)
	c := util.NewClient(&config, &param, bucketName)

	uploads := util.GetUploadsListRecursive(c, cosPath, 0, include, exclude)

	successCnt, failCnt := 0, 0
	for _, u := range uploads {
		_, err := c.Object.AbortMultipartUpload(context.Background(), u.Key, u.UploadID)
		if err != nil {
			logger.Infoln("Abort fail!    UploadID:", u.UploadID, "Key:", u.Key)
			failCnt++
		} else {
			logger.Infoln("Abort success! UploadID:", u.UploadID, "Key:", u.Key)
			successCnt++
		}
	}
	logger.Infoln("Total:", len(uploads), ",", successCnt, "Success,", failCnt, "Fail")
}
