package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func ReadConfig(fileName, filePath string) error {
	viper.SetConfigFile(fileName)
	viper.AddConfigPath(filePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("no suh config file")
		} else {
			return fmt.Errorf("no suh config file")
		}
	}
	return nil
}
