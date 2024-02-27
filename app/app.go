package app

import (
	"github.com/666999777555/go-init/mysql"
)

func Init(name string, apps ...string) error {
	var err error
	for _, val := range apps {
		switch val {
		case "mysql":
			err = mysql.InitMysql(name)
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}
