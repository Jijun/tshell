package cmd

import (
	"os"
	"strconv"
	"time"

	"tshell/util"

	"github.com/olekukonko/tablewriter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var cdnstatCmd = &cobra.Command{
	Use:   "top",
	Short: "TOP 数据查询",
	Long: ` TOP 数据查询 仅支持 90 天内数据查询

Format:
  ./tshell cdn stat --start "2022-12-01 00:00:00" --end "2022-12-01 23:59:59" 
  
Example:
  ./tshell cdn stat --start "2022-12-01 00:00:00" --end "2022-12-01 23:59:59" --domain "asset.eqh5.com" --limit 200 
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt64("limit")
		domains, _ := cmd.Flags().GetStringSlice("domain")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		filter, _ := cmd.Flags().GetString("filter")
		metric, _ := cmd.Flags().GetString("metric")
		top(domains, filter, metric, start, end, limit)
	},
}

func init() {
	cdnCmd.AddCommand(cdnstatCmd)

	cdnstatCmd.Flags().Int64("limit", 100, "只返回前N条数据，默认为最大值100，metric=url时默认为最大值1000 ")
	cdnstatCmd.Flags().StringSliceP("domain", "d", nil, "指定查询域名列表，最多可一次性查询 30 个加速域名明细")
	cdnstatCmd.Flags().String("filter", "flux", `排序使用的指标名称：
	flux：Metric 为 host 时指代访问流量，originHost 时指代回源流量
	bandwidth：Metric 为 host 时指代访问带宽，originHost 时指代回源带宽
	request：Metric 为 host 时指代访问请求数，originHost 时指代回源请求数
	fluxHitRate：平均流量命中率
	2XX：访问 2XX 状态码
	3XX：访问 3XX 状态码
	4XX：访问 4XX 状态码
	5XX：访问 5XX 状态码
	origin_2XX：回源 2XX 状态码
	origin_3XX：回源 3XX 状态码
	origin_4XX：回源 4XX 状态码
	origin_5XX：回源 5XX 状态码
	statusCode：指定访问状态码统计，在 Code 参数中填充指定状态码
	OriginStatusCode：指定回源状态码统计，在 Code 参数中填充指定状态码
	`)
	cdnstatCmd.Flags().String("metric", "url", `排序对象，支持以下几种形式：
	url：访问 URL 排序（无参数的URL），支持的 Filter 为 flux、request
	district：省份、国家/地区排序，支持的 Filter 为 flux、request
	isp：运营商排序，支持的 Filter 为 flux、request
	host：域名访问数据排序，支持的 Filter 为：flux、request、bandwidth、fluxHitRate、2XX、3XX、4XX、5XX、statusCode
	originHost：域名回源数据排序，支持的 Filter 为 flux、request、bandwidth、origin_2XX、origin_3XX、origin_4XX、origin_5XX、OriginStatusCode`)

	cdnstatCmd.Flags().String("start", "", "根据时间区间查询时，填充开始时间，如 2018-08-08 00:00:00")
	cdnstatCmd.Flags().String("end", "", "根据时间区间查询时，填充结束时间，如 2018-08-08 00:00:00")

}

func top(domains []string, filter, metric, start, end string, limit int64) {
	req := cdn.NewListTopDataRequest()
	req.Filter = common.StringPtr(filter)
	if start == "" {
		start = time.Now().Format("2006-01-02") + " 00:00:00"
	}

	if end == "" {
		end = time.Now().Format("2006-01-02") + " 23:59:59"
	}

	req.StartTime = common.StringPtr(start)
	req.EndTime = common.StringPtr(end)
	req.Limit = common.Int64Ptr(limit)
	req.Domains = common.StringPtrs(domains)
	req.Metric = common.StringPtr(metric)

	c := util.NewCdnClient(&config, &param)
	res, err := c.ListTopData(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	datas := res.Response.Data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{filter, metric, "Value"})
	for _, v := range datas {
		resource := v.Resource
		details := v.DetailData
		for _, d := range details {
			table.Append([]string{*resource, *d.Name, strconv.FormatFloat(*d.Value, 'f', 4, 64)})
		}
	}
	table.SetBorder(false)
	table.Render()

}
