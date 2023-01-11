package conf

import (
	"binance-wails/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

const (
	filePerm = 0666
	dirPerm  = 0755
)

var (
	configFile    = "config.yaml"
	configFileDir = "./"
	configContent = `window:
  width: 1024
  height: 768

proxy: "http://localhost:7890"

binancews:
  addr: "localhost:2303"
  path: "/ws"
`
)

func Init() {
	checkConfigFile()

	viper.SetConfigFile(configFile)
	viper.AddConfigPath(configFileDir)
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

func checkConfigFile() error {
	exist, err := utils.PathExists(configFileDir)
	if err != nil {
		return err
	}

	if !exist {
		os.MkdirAll(configFileDir, dirPerm)
	}

	configFilePath := filepath.Join(configFileDir, configFile)
	exist, err = utils.PathExists(configFilePath)
	if !exist {
		file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE, filePerm)
		if err != nil {
			return err
		}
		defer file.Close()

		file.WriteString(configContent)
	}

	return nil
}

func setProxy() {
	os.Setenv("HTTP_PROXY", viper.GetString("proxy"))
	os.Setenv("HTTPS_PROXY", viper.GetString("proxy"))
}
