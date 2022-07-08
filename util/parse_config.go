package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var ConfigPath string // 配置文件夹地址
var Env string        // 环境名称

// GetConfigPath 获取配置文件路径
func GetConfigPath(fileName string) string {
	return ConfigPath + "/" + fileName + ".toml"
}

// ParseConfigPath 解析配置文件
func ParseConfigPath(configPath string) error {
	path := strings.Split(configPath, "/") // 解析路径
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfigPath = prefix
	Env = path[len(path)-2]
	return nil
}

// ParseConfig 解析配置
func ParseConfig(path string, config interface{}) error {
	file, err := os.Open(path) // 打开配置文件
	if err != nil {
		return fmt.Errorf("open config file %v failed: %v ", path, err)
	}
	data, err := ioutil.ReadAll(file) // 读取配置文件
	if err != nil {
		return fmt.Errorf("read config %v failed: %v ", path, err)
	}
	v := viper.New() // 使用第三方扩展 Viper 读取配置文件
	v.SetConfigType("toml")
	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return fmt.Errorf("viper read config faild, config: %v, err: %v ", string(data), err)
	}
	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("viper Parse config faild, config: %v, err: %v ", string(data), err)
	}
	return nil
}
