package util

import (
	"net/http"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// 根据桶别名，从配置文件中加载信息，创建客户端
func NewClient(config *Config, param *Param, bucketName string) *cos.Client {
	secretID := config.Base.SecretID
	secretKey := config.Base.SecretKey
	secretToken := config.Base.SessionToken
	if param.SecretID != "" {
		secretID = param.SecretID
	}
	if param.SecretKey != "" {
		secretKey = param.SecretKey
	}
	if param.SessionToken != "" {
		secretToken = param.SessionToken
	}
	if bucketName == "" { // 不指定 bucket，则创建用于发送 Service 请求的客户端
		return cos.NewClient(nil, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     secretID,
				SecretKey:    secretKey,
				SessionToken: secretToken,
			},
		})
	} else {
		return cos.NewClient(GenURL(config, param, bucketName), &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     secretID,
				SecretKey:    secretKey,
				SessionToken: secretToken,
			},
		})
	}
}

// 根据函数参数创建客户端
func CreateClient(config *Config, param *Param, bucketIDName string) *cos.Client {
	secretID := config.Base.SecretID
	secretKey := config.Base.SecretKey
	if param.SecretID != "" {
		secretID = param.SecretID
	}
	if param.SecretKey != "" {
		secretKey = param.SecretKey
	}
	return cos.NewClient(CreateURL(bucketIDName, config.Base.Protocol, param.Endpoint), &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
}

func NewCdnClient(config *Config, param *Param) *cdn.Client {
	secretID := config.Base.SecretID
	secretKey := config.Base.SecretKey
	if param.SecretID != "" {
		secretID = param.SecretID
	}
	if param.SecretKey != "" {
		secretKey = param.SecretKey
	}

	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(secretID, secretKey)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cdn.NewClient(credential, "", cpf)

	return client
}

func NewClbClient(config *Config, param *Param, region string) *clb.Client {
	secretID := config.Base.SecretID
	secretKey := config.Base.SecretKey
	if param.SecretID != "" {
		secretID = param.SecretID
	}
	if param.SecretKey != "" {
		secretKey = param.SecretKey
	}

	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(secretID, secretKey)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "clb.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := clb.NewClient(credential, region, cpf)

	return client
}
