package cmd

import (
	"fmt"
	"os"
	"time"

	"tshell/util"

	"github.com/olekukonko/tablewriter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var cdnHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "cdn操作历史查询",
	Long: ` 预热历史查询, 刷新历史查询

Format:
  ./tshell cdn history --start "2022-12-01 00:00:00" --end "2022-12-01 23:59:59" 
  
Example:
  ./tshell cdn history --start "2022-12-01 00:00:00" --end "2022-12-01 23:59:59" --keyword "asset.eqh5.com" --limit 200 --offset 1
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt64("limit")
		offset, _ := cmd.Flags().GetInt64("offset")
		flushType, _ := cmd.Flags().GetString("type")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		keyword, _ := cmd.Flags().GetString("keyword")
		taskId, _ := cmd.Flags().GetString("taskId")
		if flushType == "url" || flushType == "path" {
			queryRefresh(args, keyword, taskId, start, end, flushType, limit, offset)
		} else if flushType == "push" {
			queryPush(args, keyword, taskId, start, end, limit, offset)
		}
	},
}

func init() {
	cdnCmd.AddCommand(cdnHistoryCmd)

	cdnHistoryCmd.Flags().Int64("limit", 20, "分页查询限制数目，默认为 20")
	cdnHistoryCmd.Flags().Int64("offset", 0, "分页查询偏移量，默认为 0 ")
	cdnHistoryCmd.Flags().StringP("type", "t", "url", "指定刷新类型查询 url：url 刷新记录 path：目录刷新记录 push: 预热 ")
	cdnHistoryCmd.Flags().String("keyword", "", "支持域名过滤，或 http(s):// 开头完整 URL 过滤")
	cdnHistoryCmd.Flags().String("taskId", "", "根据任务 ID 查询时，填充任务 ID 查询时任务 ID 与起始时间必须填充一项 ")
	cdnHistoryCmd.Flags().String("start", "", "根据时间区间查询时，填充开始时间，如 2018-08-08 00:00:00")
	cdnHistoryCmd.Flags().String("end", "", "根据时间区间查询时，填充结束时间，如 2018-08-08 00:00:00")

}

func queryRefresh(args []string, keyword, taskId, start, end, flushType string, limit, offset int64) {
	req := cdn.NewDescribePurgeTasksRequest()
	if taskId != "" {
		req.TaskId = common.StringPtr(taskId)
	}
	if keyword != "" {
		req.Keyword = common.StringPtr(keyword)
	}
	if start == "" {
		start = time.Now().Format("2006-01-02") + " 00:00:00"
	}

	if end == "" {
		end = time.Now().Format("2006-01-02") + " 23:59:59"
	}

	req.StartTime = common.StringPtr(start)
	req.EndTime = common.StringPtr(end)
	req.Limit = common.Int64Ptr(limit)
	req.Offset = common.Int64Ptr(offset)
	req.PurgeType = common.StringPtr(flushType)

	c := util.NewCdnClient(&config, &param)
	res, err := c.DescribePurgeTasks(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logs := res.Response.PurgeLogs
	totalCount := *res.Response.TotalCount
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TaskId", "Url", "Time", "Status"})
	for _, v := range logs {
		table.Append([]string{*v.TaskId, *v.Url, *v.CreateTime, *v.Status})
	}
	table.SetFooter([]string{"", "", fmt.Sprintf("offset %d Limit %d", offset, limit), fmt.Sprintf("Total %d", totalCount)})
	table.SetBorder(false)
	table.Render()

}

func queryPush(args []string, keyword, taskId, start, end string, limit, offset int64) {
	req := cdn.NewDescribePushTasksRequest()
	if taskId != "" {
		req.TaskId = common.StringPtr(taskId)
	}
	if keyword != "" {
		req.Keyword = common.StringPtr(keyword)
	}
	req.StartTime = common.StringPtr(start)
	req.EndTime = common.StringPtr(end)
	req.Limit = common.Int64Ptr(limit)
	req.Offset = common.Int64Ptr(offset)

	c := util.NewCdnClient(&config, &param)
	res, err := c.DescribePushTasks(req)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logs := res.Response.PushLogs
	totalCount := *res.Response.TotalCount
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Log", "Time", "Status"})
	for _, v := range logs {
		table.Append([]string{*v.Url, *v.CreateTime, *v.Status})
	}
	table.SetFooter([]string{"", "", fmt.Sprintf("offset %d Limit %d Total %d", offset, limit, totalCount)})
	table.SetBorder(false)
	table.Render()

}
