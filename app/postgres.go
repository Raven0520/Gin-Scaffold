package app

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Raven0520/Golang/Scaffold/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	DriverName      string `mapstructure:"driver_name"`        // 数据库驱动
	DataSourceName  string `mapstructure:"data_source_name"`   // 数据源
	MaxOpenConn     int    `mapstructure:"max_open_conn"`      // 最大连接数
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`      // 空闲连接池的最大连接数
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"` // 可以重用的最长连接时间
}

type PostgresMapConfig struct {
	List map[string]*PostgresConfig `mapstructure:"list"`
}

var logLevel = logger.Info
var PostgresPool map[string]*gorm.DB

// InitPostgresPool 初始化数据库连接 gorm 方式
func InitPostgresPool(path string, level string) error {
	SetPgSQLLogLevel(level) // 设置日志等级
	DBConfigMap := &PostgresMapConfig{}
	err := util.ParseConfig(path, DBConfigMap)
	if err != nil {
		return err
	}
	if len(DBConfigMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(util.DateTimeFormat), " empty postgres config.")
	}
	PostgresPool = map[string]*gorm.DB{}
	for configName, config := range DBConfigMap.List {
		dialector := postgres.New(postgres.Config{
			DSN: config.DataSourceName,
		})
		DBGorm, err := gorm.Open(dialector, &gorm.Config{
			QueryFields: true,
			Logger: logger.New(
				Writer{},
				logger.Config{
					SlowThreshold: time.Second * 2, // 满 SQL 阀值
					LogLevel:      logLevel,        // Log Level
					Colorful:      true,            // 禁用彩色打印
				},
			),
		})
		if err != nil {
			return err
		}
		PQ, err := DBGorm.DB()
		if err == nil {
			PQ.SetMaxOpenConns(config.MaxOpenConn)
			PQ.SetMaxIdleConns(config.MaxIdleConn)
			PQ.SetConnMaxLifetime(time.Duration(config.MaxConnLifeTime) * time.Second)
			PostgresPool[configName] = DBGorm
		} else {
			return err
		}
	}
	return nil
}

// GetPgSQLPool GetGormPool 获取数据库连接
func GetPgSQLPool(name string) (*gorm.DB, error) {
	if db, ok := PostgresPool[name]; ok {
		return db, nil
	}
	return nil, errors.New("get pool error")
}

// CloseDB 关闭数据库
func ClosePgSQLDB() error {
	for _, pool := range PostgresPool {
		db, err := pool.DB()
		if err != nil {
			return err
		}
		err = db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// SetLogLevel 设置日志级别
func SetPgSQLLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "SILENT": // 静默
		logLevel = logger.Silent
	case "ERROR":
		logLevel = logger.Error
	case "WARNING":
		logLevel = logger.Warn
	case "INFO":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}
}
