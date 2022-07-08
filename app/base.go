package app

import "github.com/Raven0520/Golang/Scaffold/util"

// BaseConf 全局变量
var BaseConf *BaseConfig

type BaseConfig struct {
	Base Base       `mapstructure:"base"`
	Http HttpConfig `mapstructure:"http"`
}

// BaseConfig 基础配置结构体
type Base struct {
	Env          string `mapstructure:"env"`
	DebugMode    string `mapstructure:"debug_mode"`
	LogLevel     string `mapstructure:"log_level"`
	TimeLocation string `mapstructure:"time_location"`
}

// HttpConfig HTTP 服务配置
type HttpConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	ReadTimeout    string `mapstructure:"read_timeout"`
	WriteTimeout   string `mapstructure:"write_timeout"`
	MaxHeaderBytes string `mapstructure:"max_header_bytes"`
}

func InitBaseConfig(path string) error {
	BaseConf = &BaseConfig{}
	if err := util.ParseConfig(path, BaseConf); err != nil {
		return err
	}
	// 设置 Env 默认值
	if BaseConf.Base.Env == "" {
		BaseConf.Base.Env = "Dev"
	}
	// 设置 Debug 默认值
	if BaseConf.Base.DebugMode == "" {
		BaseConf.Base.DebugMode = "debug"
	}
	// 设置 LogLevel 默认值
	if BaseConf.Base.LogLevel == "" {
		BaseConf.Base.LogLevel = "trace"
	}
	// 设置 默认时区
	if BaseConf.Base.TimeLocation == "" {
		BaseConf.Base.TimeLocation = "Asia/Shanghai"
	}
	return nil
}

// GetEnv 获取环境名称
func GetEnv() string {
	Env := "Dev" // 默认环境
	if BaseConf.Base.Env != "" {
		Env = BaseConf.Base.Env
	}
	return Env
}

// GetDebugMode Debug 模式
func GetDebugMode() string {
	Mode := "debug"
	if BaseConf.Base.DebugMode != "" {
		Mode = BaseConf.Base.DebugMode
	}
	return Mode
}
