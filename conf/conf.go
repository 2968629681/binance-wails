package conf

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	setProxy()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(" The configuration file was modified")
	})
}

func setProxy() {
	os.Setenv("HTTP_PROXY", viper.GetString("proxy"))
	os.Setenv("HTTPS_PROXY", viper.GetString("proxy"))
}
