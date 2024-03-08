package consul

import (
	"context"
	"fmt"
	"github.com/666999777555/go-init/config"
	"github.com/666999777555/go-init/redis"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
	"strconv"
	"time"
)

const CONSUL_KEY = "consul:node:index"

type ConsulConfig struct {
	Consul struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"consul"`
}

func getConfig(nacosGroup, serviceName string) (*ConsulConfig, error) {
	cnf, err := config.GetConfig(nacosGroup, serviceName)
	if err != nil {
		return nil, err
	}
	consulCnf := new(ConsulConfig)
	err = yaml.Unmarshal([]byte(cnf), consulCnf)
	if err != nil {
		return nil, err
	}
	return consulCnf, err
}

func getIndex(ctx context.Context, serviceName string, indexLn int) (int, error) {
	exist, err := redis.ExistKey(ctx, serviceName, CONSUL_KEY)
	if err != nil {
		return 0, err
	}
	if exist {
		indexStr, err := redis.GetByKey(ctx, serviceName, CONSUL_KEY)
		if err != nil {
			return 0, err
		}
		index, err := strconv.Atoi(indexStr)
		newIndex := index + 1

		if newIndex >= indexLn {
			newIndex = 0
		}
		err = redis.SetKey(ctx, serviceName, CONSUL_KEY, newIndex, time.Duration(0))
		if err != nil {
			return 0, err
		}
		return index, nil
	}
	err = redis.SetKey(ctx, serviceName, "consul:node:index", 0, time.Duration(0))
	if err != nil {
		return 0, err
	}
	return 0, err
}

func Agent(ctx context.Context, name string) (string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", err
	}
	sr, info, err := client.Agent().AgentHealthServiceByName(name)
	if err != nil {
		return "", err
	}
	if sr != "passing" {
		return "", fmt.Errorf("******")
	}
	index, err := getIndex(ctx, name, len(info))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v:%v", info[index].Service.Address, info[index].Service.Port), nil
}

func ServiceRegister(nacosGroup, name string, address, port string) error {
	cof, err := getConfig(nacosGroup, name)
	if err != nil {
		return err
	}
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%v:%v", cof.Consul.Ip, cof.Consul.Port),
	})

	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    "register",
		Tags:    []string{"GRPC"},
		Port:    portInt,
		Address: address,
		//Check: &api.AgentServiceCheck{
		//	GRPC:                           fmt.Sprintf("%v:%v", address, port),
		//	Interval:                       "5s",
		//	DeregisterCriticalServiceAfter: "10s",
		//},
	})
}
