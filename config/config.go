package config

import (
	"fmt"
	"home_control_hub/internal/utils"
	"os"

	"gopkg.in/yaml.v3"
)

// 全局配置变量
var GlobalConfig Config

type Config struct {
	NasConfig       `yaml:"nas"`
	Ip2regionConfig `yaml:"ip2region"`
}

type NasConfig struct {
	Url      string `yaml:"url"`
	Mac      string `yaml:"mac"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Ip2regionConfig struct {
	Path      string `yaml:"path"`
	MakerPath string `yaml:"maker_path"`
	DbPath    string `yaml:"db_path"`
}

func init() {
	configFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println("[config]读取配置文件失败：", err)
		return
	}

	err = yaml.Unmarshal(configFile, &GlobalConfig)
	if err != nil {
		fmt.Println("[config]解析配置文件失败：", err)
		return
	}

	ip2regionHandler()

	fmt.Println("[config]配置初始化成功")
}

func ip2regionHandler() {
	// 获取当前工作目录
	root := utils.GetPwd()
	GlobalConfig.Ip2regionConfig.DbPath = root + "/assets/ip2region.xdb"
	GlobalConfig.Ip2regionConfig.MakerPath = GlobalConfig.Path + GlobalConfig.Ip2regionConfig.MakerPath
}
