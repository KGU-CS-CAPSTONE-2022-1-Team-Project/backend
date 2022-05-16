package hosts

import "github.com/spf13/viper"

var (
	Owner string
)

func init() {
	if !viper.IsSet("owner") {
		viper.AddConfigPath("configs/gateway")
		viper.SetConfigType("yaml")
		viper.SetConfigName("services")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	info := viper.GetStringMapString("owner")
	Owner = info["protocol"] + "://" + info["host"] + ":" + info["port"]
}
