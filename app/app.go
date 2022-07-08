package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.f-fans.cn/components/golang/scaffold/util"
)

var TimeLocation *time.Location // 时区

// Init 初始化系统
func Init(configPath string) error {
	return InitModule(configPath, []string{"base", "swagger", "postgres", "redis"}) // 配置需要解析的配置文件名称
}

// InitModule 初始化模块
func InitModule(configPath string, modules []string) error {
	config := flag.String("config", configPath, "input config file like ./config/develop/") // 返回变量地址
	flag.Parse()                                                                            // 执行解析
	if *config == "" {
		flag.Usage()
		os.Exit(1) // 当前程序以给定的状态代码退出
	}

	log.Println("Start Loading Resources ------------------------------------------------") // 开始加载资源
	log.Printf("[INFO]  Config Path : %s \n", *config)                                      // 打印配置路径

	// 设置ip信息，优先设置便于日志打印
	util.SetLocalIPs()

	// 解析配置文件目录
	if err := util.ParseConfigPath(*config); err != nil {
		return err
	}
	log.Printf("[INFO] %s\n", " Parse Config Path Done.") // 解析配置文件成功

	// 初始化配置文件
	if err := InitViperConfig(); err != nil {
		return err
	}
	log.Printf("[INFO] %s\n", " Viper Config Done.")

	// 加载 Base 设置
	if util.InSliceString("base", modules) {
		if err := InitBaseConfig(util.GetConfigPath("base")); err != nil {
			fmt.Printf("[ERROR] %s  InitBaseConfig: %s\n", time.Now().Format(util.DateTimeFormat), err.Error())
		}
		log.Printf("[INFO] %s\n", " Base Config Done.")
	}

	// 加载redis配置
	if util.InSliceString("redis", modules) {
		if err := InitRedisConfig(util.GetConfigPath("redis")); err != nil {
			fmt.Printf("[ERROR] %s InitRedisConfig: %s\n", time.Now().Format(util.DateTimeFormat), err.Error())
		}
		log.Printf("[INFO] %s\n", " Redis Config Done.")
	}

	// 加载 Mysql 配置 并初始化实例
	if util.InSliceString("mysql", modules) {
		if err := InitMySQLPool(util.GetConfigPath("mysql"), BaseConf.Base.LogLevel); err != nil {
			fmt.Printf("[ERROR] %s InitMySQLPool: %s\n", time.Now().Format(util.DateTimeFormat), err.Error())
		}
		log.Printf("[INFO] %s\n", " MySQL Config Done.")
	}

	// 加载 Postgres 配置 并初始化实例
	if util.InSliceString("postgres", modules) {
		if err := InitPostgresPool(util.GetConfigPath("postgres"), BaseConf.Base.LogLevel); err != nil {
			fmt.Printf("[ERROR] %s InitPostgresPool: %s\n", time.Now().Format(util.DateTimeFormat), err.Error())
		}
		log.Printf("[INFO] %s\n", " Postgres Config Done.")
	}

	// 设置时区
	location, err := time.LoadLocation(BaseConf.Base.TimeLocation)
	if err != nil {
		return err
	}
	TimeLocation = location
	log.Println("--------------------------------------------- Loading Resources Success ") // 加载资源成功
	return nil
}

// Destroy 公共销毁函数
func Destroy() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", "Start Destroy Resources.") // 开始销毁加载的资源
	err := CloseMySQLDB()                                 // 关闭数据库连接
	if err != nil {
		log.Printf("[INFO] %s\n", "Close MySQL Connect Failed.") // 关闭数据库连接失败
		return
	} // 错误没有返回
	log.Printf("[INFO] %s\n", "Close MySQL Connect Success.") // 关闭数据库连接成功
	err = ClosePgSQLDB()                                      // 关闭数据库连接
	if err != nil {
		log.Printf("[INFO] %s\n", "Close PgSQL Connect Failed.") // 关闭数据库连接失败
		return
	} // 错误没有返回
	log.Printf("[INFO] %s\n", "Close PgSQL Connect Success.") // 关闭数据库连接成功
	err = CloseRedisDB()                                      // 关闭 Redis 数据库
	if err != nil {
		log.Printf("[INFO] %s\n", "Close Redis Connect Failed.") // 关闭Redis连接失败
		return
	}
	log.Printf("[INFO] %s\n", "Close Redis Connect Success.") // 关闭Redis连接成功
	log.Printf("[INFO] %s\n", "Destroy Resources Success.")   // 销毁加载资源成功
}
