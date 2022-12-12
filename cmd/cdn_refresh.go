package cmd

import (
	"fmt"
	"os"
	"strings"

	"tshell/util"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var cdnrefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "刷新CDN缓存,可访问到最新资源",
	Long: ` 您的源站有资源更新/需要清理违规资源/域名有配置变更，为避免全网用户受节点缓存影响仍访问到旧的资源/受旧配置的影响，可提交刷新任务，保证全网用户可访问到最新资源或正常访问

Format:
  ./tshell refresh http(s)://xxx.com/somekey

Example:
  ./tshell refresh -r http://example.com/test
  ./tshell refresh  http://example.com/test/a.png
  ./tshell refresh  -f /tmp/a.txt
  `,
	Args: func(cmd *cobra.Command, args []string) error {

		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return fmt.Errorf("need specify -f xxx.txt  or %s ", err.Error())
			}
			for _, arg := range args {
				hasHttpPrefix := strings.HasPrefix(arg, "http")
				if !hasHttpPrefix {
					return fmt.Errorf("Invalid arguments!, url must begin with http or https ")
				}
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		urlEncode, _ := cmd.Flags().GetBool("urlencode")
		recursive, _ := cmd.Flags().GetBool("recursive")
		area, _ := cmd.Flags().GetString("area")
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			lines, err := util.ReadLine2Slice(file)
			if err != nil {
				logger.Infoln(err)
				os.Exit(1)
			}

			args = append(args, lines...)
		}

		flushType, _ := cmd.Flags().GetString("flushType")
		if recursive {
			flushPath(args, urlEncode, area, flushType)
		} else {
			flushUrls(args, urlEncode, area)
		}
	},
}

func init() {
	cdnCmd.AddCommand(cdnrefreshCmd)

	cdnrefreshCmd.Flags().BoolP("urlencode", "u", false, "是否对中文字符进行编码后刷新")
	cdnrefreshCmd.Flags().StringP("area", "a", "", "刷新区域, 无此参数时，默认刷新加速域名所在加速区域, 可选项: mainland, overseas")
	cdnrefreshCmd.Flags().StringP("file", "f", "", "支持 txt 格式文件, URL 必须包含 http:// 或 https://，例如 http://www.test.com/test.html，一行一个 ")
	cdnrefreshCmd.Flags().BoolP("recursive", "r", false, "刷新目录")
	cdnrefreshCmd.Flags().StringP("flushType", "t", "flush", " 刷新类型, flush：刷新产生更新的资源, delete：刷新全部资源")
}

func flushPath(args []string, urlEncode bool, area, flushType string) {

	var urls []string
	for _, arg := range args {
		urls = append(urls, arg)
	}
	request := cdn.NewPurgePathCacheRequest()
	request.UrlEncode = common.BoolPtr(urlEncode)
	request.Paths = common.StringPtrs(urls)
	request.Area = common.StringPtr(area)
	request.FlushType = common.StringPtr(flushType)
	c := util.NewCdnClient(&config, &param)
	response, err := c.PurgePathCache(request)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logger.Infof("RequestID: %s TaskId: %s ", *response.Response.RequestId, *response.Response.TaskId)

}

func flushUrls(args []string, urlEncode bool, area string) {

	var urls []string
	for _, arg := range args {
		urls = append(urls, arg)
	}
	request := cdn.NewPurgeUrlsCacheRequest()
	request.UrlEncode = common.BoolPtr(urlEncode)
	request.Urls = common.StringPtrs(urls)
	request.Area = common.StringPtr(area)
	c := util.NewCdnClient(&config, &param)
	response, err := c.PurgeUrlsCache(request)
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logger.Infof("RequestID: %s TaskId: %s ", *response.Response.RequestId, *response.Response.TaskId)
}
