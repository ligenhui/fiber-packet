package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	Viper    *viper.Viper                //viper引擎
	drive    = "viper"                   //支持的驱动 viper
	filePath = "core/config/config.toml" //文件路径
)

func Init() {
	//初始化
	switch drive {
	case "viper":
		initViper()
	default:
		panic(fmt.Errorf("Fatal error config, Drive %s does not exist \n", drive))
	}
}

func SetDrive(driveName string) {
	drive = driveName
}

func SetFilePath(path string) {
	filePath = path
}

func initViper() {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	Viper = viper.GetViper()
}
