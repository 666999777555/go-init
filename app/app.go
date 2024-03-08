package app

import (
	"github.com/666999777555/go-init/config"
	"github.com/666999777555/go-init/mysql"
)

func Init(
	serviceName string,
	naocsIP, nacosPort string,
	app ...string,
) error {
	if err := config.InitNacos(naocsIP, nacosPort); err != nil {
		return err
	}
	for _, val := range app {
		switch val {
		case "mysql":
			err := mysql.InitMysql(serviceName)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
