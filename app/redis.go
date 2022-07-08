package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/Raven0520/Gin-Scaffold/model"
	"github.com/Raven0520/Gin-Scaffold/util"
	"github.com/gomodule/redigo/redis"
)

// ConfigRedisMap 全局变量
var ConfigRedisMap *model.RedisMapConfig
var RedisPool map[string]*redis.Pool

// InitRedisConfig 加载 Redis 配置
func InitRedisConfig(path string) error {
	RedisConfigMap := &model.RedisMapConfig{}
	err := util.ParseConfig(path, RedisConfigMap)
	if err != nil {
		return err
	}
	if len(RedisConfigMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(util.DateTimeFormat), " empty redis config.")
	}
	RedisPool = map[string]*redis.Pool{}
	for configName, config := range RedisConfigMap.List {
		dialector := &redis.Pool{
			MaxIdle:         config.MaxIdle,   // 最大空闲连接数
			MaxActive:       config.MaxActive, // 分配的最大连接数
			MaxConnLifetime: time.Duration(config.ConnTimeout),
			IdleTimeout:     time.Duration(config.IdelTimeout),
			Wait:            true,
			Dial: func() (redis.Conn, error) {
				options := redis.DialDatabase(config.Db)
				password := redis.DialPassword(config.Password)
				// **重要** 设置读写超时
				readTimeout := redis.DialReadTimeout(time.Second * time.Duration(config.ReadTimeout))
				writeTimeout := redis.DialReadTimeout(time.Second * time.Duration(config.WriteTimeout))
				conTimeout := redis.DialConnectTimeout(time.Second * time.Duration(config.ConnTimeout))
				c, err := redis.Dial("tcp", config.ProxyList, options, password, readTimeout, writeTimeout, conTimeout)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		}
		RedisPool[configName] = dialector
	}
	return nil
}

// GetRedisPool 获取 Redis 数据库连接
func GetRedisPool(name string) (*redis.Pool, error) {
	if pool, ok := RedisPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("GetRedisPoolError") // 获取 Redis 连接池错误
}

// CloseRedisDB 关闭 Redis 数据库
func CloseRedisDB() error {
	for _, pool := range RedisPool {
		err := pool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
