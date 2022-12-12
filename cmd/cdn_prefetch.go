package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"tshell/util"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var cdnprefetchCmd = &cobra.Command{
	Use:   "prefetch",
	Short: "cdn缓存预热",
	Long: ` 当您有运营活动或安装包/升级包发布等，可提交预热任务，提前将静态资源预热至 CDN 加速节点，降低源站压力，提升用户服务可用性和用户体验

Format:
  ./tshell cdn prefetch http(s)://xxx.com/somekey

Example:
  ./tshell cdn prefetch  http://example.com/test/a.png
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
		//recursive, _ := cmd.Flags().GetBool("recursive")
		area, _ := cmd.Flags().GetString("area")
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			f, err := os.Open(file)
			if err != nil {
				logger.Fatalln(err)
				return
			}
			defer f.Close()
			reader := bufio.NewReader(f)
			for {
				line, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}
				hasHttpPrefix := strings.HasPrefix(string(line), "http")
				if !hasHttpPrefix {
					logger.Fatal("Invalid arguments!, url must begin with http or https ")
				}
				args = append(args, string(line))
			}
		}

		prefetch(args, urlEncode, area)
	},
}

func init() {
	cdnCmd.AddCommand(cdnprefetchCmd)

	cdnprefetchCmd.Flags().BoolP("urlencode", "u", false, "是否对中文字符进行编码后刷新")
	cdnprefetchCmd.Flags().StringP("area", "a", "", "刷新区域, 无此参数时，默认刷新加速域名所在加速区域, 可选项: global,  mainland, overseas")
	cdnprefetchCmd.Flags().StringP("file", "f", "", "支持 txt 格式文件, URL 必须包含 http:// 或 https://，例如 http://www.test.com/test.html，一行一个 ")
}

func prefetch(args []string, urlEncode bool, area string) {

	var urls []string
	for _, arg := range args {
		urls = append(urls, arg)
	}
	fmt.Println(urls)
	request := cdn.NewPushUrlsCacheRequest()
	request.UrlEncode = common.BoolPtr(urlEncode)
	request.Urls = common.StringPtrs(urls)
	request.Area = common.StringPtr(area)
	c := util.NewCdnClient(&config, &param)
	response, err := c.PushUrlsCache(request)

	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	logger.Infof("RequestID: %s TaskId: %s ", *response.Response.RequestId, *response.Response.TaskId)

}
