package register

import (
	"fmt"
	"github.com/666999777555/go-init/config"
	"google.golang.org/grpc"
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
	ConfigInfo, err := config.GetConfig("DEFAULT_GROUP", serviceName)
	if err != nil {
		return nil, err
	}
	cnf := new(Config)
	err = yaml.Unmarshal([]byte(ConfigInfo), cnf)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func Register(serviceName string, res func(s *grpc.Server)) error {
	cof, err := getConfig(serviceName)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cof.App.Ip, cof.App.Port))
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	g := grpc.NewServer()
	log.Printf("sever listening at %v", lis.Addr())
	if err := g.Serve(lis); err != nil {
		log.Fatalf(err.Error())
		return err
	}
	return err
}
