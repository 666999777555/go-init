package register

import (
	"fmt"
	"github.com/666999777555/go-init/config"
	"github.com/666999777555/go-init/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gopkg.in/yaml.v3"
	"log"
	"net"
)

type Config struct {
	App struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"app"`
}

func getConfig(serviceName string) (*Config, error) {
	configInfo, err := config.GetConfig("DEFAULT_GROUP", serviceName)
	if err != nil {
		return nil, err
	}
	cnf := new(Config)
	err = yaml.Unmarshal([]byte(configInfo), cnf)
	fmt.Println(cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}

func Register(nacosGroup, serviceName string, register func(s *grpc.Server)) error {
	cof, err := getConfig(serviceName)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "10.2.171.101", cof.App.Port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
		return err
	}

	err = consul.ServiceRegister(nacosGroup, serviceName, cof.App.Ip, cof.App.Port)
	if err != nil {
		return err
	}
	g := grpc.NewServer()
	//健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	register(g)
	log.Printf("sever listening at %v", lis.Addr())
	if err := g.Serve(lis); err != nil {
		log.Fatalf(err.Error())
		return err
	}
	return err
}
