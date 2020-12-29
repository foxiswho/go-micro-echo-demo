package main

import (
	"context"
	nacos "github.com/liangzibo/go-plugins-micro-registry-nacos/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

const (
	NamespaceId = "9d5d3937-27a6-45a4-b300-e30dc3656a90"
	NacosHost   = "192.168.0.254"
	NacosPort   = 8848
	serviceName = "go.micro.web.echo"
)

func main() {
	//服务注册
	registry, err := nacosRegistry()
	if err != nil {
		log.Fatal(err)
	}
	// 创建 service
	service := web.NewService(
		web.Name(serviceName), //服务名称
		web.Address(":8071"),  //端口
		web.Version("0.0.1"),  //版本
		//web.RegisterTTL(time.Second*30),
		//web.RegisterInterval(time.Second*15),
		nacos.WebRegistry(registry), // nacos Metadata 源数据
		web.Handler(NewRouter()),
	)
	//log.Infof("ops=%v", service.Options())
	// 初始化 service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	//
	//运行
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// nacos 服务注册
func nacosRegistry() (registry.Registry, error) {
	//创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NacosHost,
			ContextPath: "/nacos",
			Port:        NacosPort,
			//Scheme:      "http",
		},
	}
	// 创建服务发现客户端
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	r := nacos.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{NacosHost}
		options.Context = context.WithValue(context.Background(), "naming_client", namingClient)
	})
	return r, nil
}
