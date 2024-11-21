package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
)

// InitNaCos 初始化 NaCos 客户端并获取配置内容。
// 该函数通过创建一个 NaCos 配置客户端，连接到指定的 NaCos 服务并从中获取配置。
//
// 参数:
// - Rename: 命名空间 ID，用于区分不同的命名空间。
// - host: NaCos 服务器的 IP 地址。
// - port: NaCos 服务器的端口号。
// - dataId: 配置的 DataId，用于标识配置。
// - group: 配置的 Group，通常用于区分不同环境或不同模块的配置。
//
// 返回值:
// - 返回 NaCos 配置的内容（字符串）和可能发生的错误。
// 如果配置获取成功，返回配置内容；否则，返回空字符串和错误信息。
func InitNaCos(Rename string, host string, port uint64, dataId string, group string) (string, error) {
	// 初始化文件夹，存储日志和缓存
	// 创建日志目录
	err := os.MkdirAll("./tmp/nacos/log", 0777)
	if err != nil {
		return "", err
	}
	// 创建缓存目录
	err = os.MkdirAll("./tmp/nacos/cache", 0777)
	if err != nil {
		return "", err
	}

	// 配置与 NaCos 服务的连接
	clientConfig := constant.ClientConfig{
		NamespaceId:         Rename,              // 命名空间 ID，用于区分不同的配置空间
		TimeoutMs:           5000,                // 配置请求超时时间，单位毫秒
		NotLoadCacheAtStart: true,                // 启动时不加载缓存，避免读取过时数据
		LogDir:              "./tmp/nacos/log",   // 日志文件夹路径
		CacheDir:            "./tmp/nacos/cache", // 缓存文件夹路径
		LogLevel:            "debug",             // 日志等级，调试级别
	}

	// 配置 NaCos 服务的连接信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      host,     // NaCos 服务器 IP 地址
			ContextPath: "/nacos", // NaCos 服务的上下文路径
			Port:        port,     // NaCos 服务器的端口
			Scheme:      "http",   // 使用的协议
		},
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs, // NaCos 服务器配置
		"clientConfig":  clientConfig,  // 客户端配置
	})
	if err != nil {
		return "", err // 返回错误
	}

	// 从 NaCos 获取配置内容
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId, // 配置的 DataId
		Group:  group,  // 配置的 Group
	})

	// 返回获取的配置内容或错误信息
	return content, nil
}
