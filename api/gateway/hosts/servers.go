package hosts

import (
	"backend/tool"
	"github.com/spf13/viper"
)

var (
	Owner string
)

func init() {
	if !viper.IsSet("owner") {
		tool.ReadConfig("configs/gateway", "services", "yaml")
	}
	info := viper.GetStringMapString("owner")
	Owner = info["protocol"] + "://" + info["host"] + ":" + info["port"]
}
