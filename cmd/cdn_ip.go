package cmd

import (
	"os"
	"tshell/util"

	"github.com/olekukonko/tablewriter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP归属查询",
	Long: ` IP归属查询,提供查询 IP 是否为腾讯云 CDN 加速节点的 IP，以及 IP 所在省份及运营商

Format:
  ./tshell cdn ip [ipv4 or ipv6]
  
Example:
  ./tshell cdn ip  1.2.3.4
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip2geo(args)
	},
}

func init() {
	cdnCmd.AddCommand(ipCmd)
}

func ip2geo(args []string) {
	req := cdn.NewDescribeCdnIpRequest()
	req.Ips = common.StringPtrs(args)
	c := util.NewCdnClient(&config, &param)
	res, err := c.DescribeCdnIp(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	datas := res.Response.Ips
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ip", "City", "Platform", "Location", "Area"})
	for _, v := range datas {
		table.Append([]string{*v.Ip, *v.City, *v.Platform, *v.Location, *v.Area})
	}
	table.SetBorder(false)
	table.Render()

}
