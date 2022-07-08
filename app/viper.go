package app

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Raven0520/Golang/Scaffold/util"
	"github.com/spf13/viper"
)

// InitViperConfig 初始化配置文件
func InitViperConfig() error {
	file, err := os.Open(util.ConfigPath + "/")
	if err != nil {
		return err
	}
	fileList, err := file.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f := range fileList {
		// 跳过系统文件
		if f.Name() == ".DS_Store" {
			continue
		}
		if !f.IsDir() {
			bts, err := ioutil.ReadFile(util.ConfigPath + "/" + f.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("toml")
			err = v.ReadConfig(bytes.NewBuffer(bts))
			if err != nil {
				return err
			}
			pathArr := strings.Split(f.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	return nil
}
