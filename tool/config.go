package tool

import "github.com/spf13/viper"

func ReadConfig(configPath, filename, fileType string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(filename)
	viper.SetConfigType(fileType)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
