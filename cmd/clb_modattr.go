package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
	"tshell/util"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var modifyCmd = &cobra.Command{
	Use:   "modattr",
	Short: "modify clb attribute",
	Long: ` 修改clb属性, ssl证书, http2
Format:
 ./tshell clb modattr -lb -lbl -d -pub -pri   -http2 -default

 Example:
 ./tshell clb modattr -lb lb-23sm4ft7 -lbl lbl-cje1ii4l -d *.eqxiu.com -pub /tmp/a.pub -pri /tmp/a.pri
`,
	Args: func(cmd *cobra.Command, args []string) error {
		//if err := cobra.ExactArgs(5)(cmd, args); err != nil {
		//	return err
		//}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lb, _ := cmd.Flags().GetString("lb")
		lbl, _ := cmd.Flags().GetString("lbl")
		domain, _ := cmd.Flags().GetString("domain")
		pubKey, _ := cmd.Flags().GetString("pub")
		priKey, _ := cmd.Flags().GetString("pri")
		http2, _ := cmd.Flags().GetBool("http2")
		defServer, _ := cmd.Flags().GetBool("default")

		modify(lb, lbl, domain, pubKey, priKey, http2, defServer)
	},
}

func init() {
	clbCmd.AddCommand(modifyCmd)
	modifyCmd.Flags().String("lb", "", "loadBalancer id")
	modifyCmd.Flags().String("lbl", "", "loadBalancer listenner id")
	modifyCmd.Flags().StringP("domain", "d", "", "domain")
	modifyCmd.Flags().String("pub", "", "public keyfile")
	modifyCmd.Flags().String("pri", "", "private keyfile")
	modifyCmd.Flags().BoolP("http2", "", false, "enable http2")
	modifyCmd.Flags().BoolP("default", "", false, "enable default server")
}

func modify(lb, lbl, domain, pubKey, priKey string, http2, defaultServer bool) {
	req := clb.NewModifyDomainAttributesRequest()
	req.Domain = common.StringPtr(domain)
	req.LoadBalancerId = common.StringPtr(lb)
	req.ListenerId = common.StringPtr(lbl)
	req.DefaultServer = common.BoolPtr(defaultServer)
	req.Http2 = common.BoolPtr(http2)
	pubSlice, _ := util.ReadLine2Slice(pubKey)
	priSlice, _ := util.ReadLine2Slice(priKey)

	pubKeyStr := strings.Join(pubSlice, "\\n")
	priKeyStr := strings.Join(priSlice, "\\n")
	req.Certificate.CertContent = common.StringPtr(pubKeyStr)
	req.Certificate.SSLMode = common.StringPtr("UNIDIRECTIONAL")
	req.Certificate.CertKey = common.StringPtr(priKeyStr)
	req.Certificate.CertName = common.StringPtr(fmt.Sprintf("%s-%s", domain, time.Now().Format("2006-01-02")))

	c := util.NewClbClient(&config, &param)
	res, err := c.ModifyDomainAttributes(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}

	logger.Infof("RequestID: %s ", *res.Response.RequestId)
}
