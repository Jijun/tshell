package cmd

import (
	"fmt"
	"os"

	"tshell/util"

	"github.com/olekukonko/tablewriter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var domainsCmd = &cobra.Command{
	Use:   "domain",
	Short: "cdn域名查询",
	Long: ` cdn域名列表,回源配置等 

Format:
  ./tshell cdn domain  
  
Example:
  ./tshell cdn domain 
  ./tshell cdn domain asset.eqh5.com
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt64("limit")
		offset, _ := cmd.Flags().GetInt64("offset")
		domains(args, offset, limit)
	},
}

func init() {
	cdnCmd.AddCommand(domainsCmd)

	domainsCmd.Flags().Int64("offset", 0, "分页查询偏移量，默认为 0 ")
	domainsCmd.Flags().Int64("limit", 100, "只返回前N条数据，默认为最大值100，metric=url时默认为最大值1000 ")
}

func domains(args []string, offset, limit int64) {
	c := util.NewCdnClient(&config, &param)
	req := cdn.NewDescribeDomainsRequest()
	req.Limit = common.Int64Ptr(limit)
	req.Offset = common.Int64Ptr(offset)
	var filters []*cdn.DomainFilter
	for _, f := range args {
		filter := &cdn.DomainFilter{
			Name:  common.StringPtr("domain"),
			Value: common.StringPtrs([]string{f}),
			Fuzzy: common.BoolPtr(false),
		}
		filters = append(filters, filter)
	}
	req.Filters = filters

	res, err := c.DescribeDomains(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	datas := res.Response.Domains

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "SrvType", "Status", "Cname", "OriginType", "Area"})
	for _, v := range datas {
		table.Append([]string{*v.Domain, *v.ServiceType, *v.Status, *v.Cname, *v.Origin.OriginType, *v.Area})
	}
	totalCount := *res.Response.TotalNumber
	//table.SetFooter([]string{"offset", fmt.Sprintf("%d", offset), "limit", fmt.Sprintf("%d", limit), "total", fmt.Sprintf("%d", *total)})
	table.SetFooter([]string{"", "", "", "", fmt.Sprintf("offset %d Limit %d", offset, limit), fmt.Sprintf("Total %d", totalCount)})
	table.SetBorder(false)
	table.Render()

}
