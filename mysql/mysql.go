package mysql

import (
	"fmt"
	"github.com/666999777555/go-init/config"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// mysql的配置结构体
type mysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Data     string `yaml:"data"`
}

func InitMysql(severName string) error {
	err := config.InitNacos()
	if err != nil {
		return err
	}
	type Val struct {
		mysql mysqlConfig
	}
	mysqlConfigVal := Val{}
	content, err := config.GetConfig("user", "demo")
	err = yaml.Unmarshal([]byte(content), &mysqlConfigVal)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = config.ListConfig("user", "demo")
	if err != nil {
		return err
	}

	fmt.Println(content)
	fmt.Println(mysqlConfigVal)
	configMysql := mysqlConfigVal.mysql
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		configMysql.Name,
		configMysql.Password,
		configMysql.Host,
		configMysql.Port,
		configMysql.Data,
	)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func WithTx(txFc func(tx *gorm.DB) error) {
	var err error
	tx := Db.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
