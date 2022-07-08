package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.f-fans.cn/components/golang/scaffold/crontab"
	"gitlab.f-fans.cn/components/golang/scaffold/model"
	"gitlab.f-fans.cn/components/golang/scaffold/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	DriverName      string `mapstructure:"driver_name"`        // 数据库驱动
	DataSourceName  string `mapstructure:"data_source_name"`   // 数据源
	MaxOpenConn     int    `mapstructure:"max_open_conn"`      // 最大连接数
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`      // 空闲连接池的最大连接数
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"` // 可以重用的最长连接时间
}

type MySQLMapConfig struct {
	List map[string]*MySQLConfig `mapstructure:"list"`
}

// Writer 数据库日志输出
type Writer struct{}

func (w Writer) Printf(f string, args ...interface{}) {
	const (
		LocationIndex = iota
		CostIndex
		RowsIndex
		SQLIndex
	)
	var (
		ok, slow              bool
		cost                  float64
		rows, ErrorIndex      uint64
		location, t, err, SQL string
	)
	Log := &model.SQLLog{}
	t = "logs"
	if len(args) == 5 {
		ErrorIndex = 1
		if err, ok = args[ErrorIndex].(string); ok {
			Log.Error = err
		}
		t = "errors"
	}
	if location, ok = args[LocationIndex].(string); ok {
		Log.Location = location
	}
	if cost, ok = args[CostIndex+ErrorIndex].(float64); ok {
		Log.Cost = cost
		if cost > 2000 {
			slow = true
			t = "slow"
		}
		Log.Slow = slow
	}
	if rows, ok = args[RowsIndex+ErrorIndex].(uint64); ok {
		Log.Rows = rows
	}
	if SQL, ok = args[SQLIndex+ErrorIndex].(string); ok {
		Log.SQL = SQL
	}
	Log.Type = t
	Log.Tag = map[string]string{"slow": strconv.FormatBool(slow), "type": t}
	crontab.SQLLogChan <- Log
}

var MySQLPool map[string]*gorm.DB

// InitMySQLPool 初始化 MySQL 数据库连接池
func InitMySQLPool(path string, level string) error {
	SetMySQLLogLevel(level) // 设置日志等级
	MySQLConfigMap := &MySQLMapConfig{}
	err := util.ParseConfig(path, MySQLConfigMap)
	if err != nil {
		return err
	}
	if len(MySQLConfigMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(util.DateTimeFormat), " Empty MySQL Config.")
	}

	MySQLPool = map[string]*gorm.DB{}
	for configName, config := range MySQLConfigMap.List {
		dialector := mysql.New(mysql.Config{
			DSN: config.DataSourceName,
		})
		DBGorm, err := gorm.Open(dialector, &gorm.Config{
			QueryFields: true,
			Logger: logger.New(
				// log.New(os.Stdout, "\r\n", log.LstdFlags),
				Writer{},
				logger.Config{
					SlowThreshold: time.Second * 2, // 慢 SQL 阀值
					LogLevel:      logLevel,        // Log Level
					Colorful:      false,           // 禁用彩色打印
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
			MySQLPool[configName] = DBGorm
		} else {
			return err
		}
	}
	return nil
}

// GetMySQLPool GetGormPool 获取数据库连接
func GetMySQLPool(name string) (*gorm.DB, error) {
	if db, ok := MySQLPool[name]; ok {
		return db, nil
	}
	return nil, errors.New("get pool error")
}

// CloseMySQLDB 关闭数据库
func CloseMySQLDB() error {
	for _, pool := range MySQLPool {
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

// SetMySQLLogLevel 设置日志级别
func SetMySQLLogLevel(level string) {
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
