package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
)

var client config_client.IConfigClient

func InitNacos(ip, port string) error {
	var err error
	portInt, err := strconv.Atoi(port)
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, uint64(portInt), constant.WithContextPath("/nacos")),
	}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return err
}

func GetConfig(group, dataId string) (string, error) {
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, err
}

// ListConfig 完成mysql的监听
func ListConfig(group, dataId string) error {
	return client.ListenConfig(
		vo.ConfigParam{
			DataId: dataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			},
		},
	)
}
