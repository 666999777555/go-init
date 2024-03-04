package register

import (
	"fmt"
	"github.com/666999777555/go-init/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

func Client(toService string) (*grpc.ClientConn, error) {
	cnfStr, err := config.GetConfig("demo", toService)
	if err != nil {
		return nil, err
	}
	cnf := new(Config)
	err = yaml.Unmarshal([]byte(cnfStr), &cnf)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", cnf.App.Ip, cnf.App.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
