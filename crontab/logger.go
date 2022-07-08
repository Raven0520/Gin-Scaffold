package crontab

import (
	"context"
	"fmt"
	"time"

	"github.com/Raven0520/Gin-Scaffold/model"
	"github.com/Raven0520/Gin-Scaffold/util"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	SQLLogChan      = make(chan *model.SQLLog, 1000)
	StopWriteSQLLog = make(chan bool, 1)
)

// GetInfluxConfig 获取 Influx 配置
func GetInfluxConfig(config string) (conf *model.InfluxConfig, err error) {
	configs := &model.InfluxConfigs{}
	err = util.ParseConfig("./config/influx.toml", configs)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conf, ok := configs.List[config]; ok {
		return conf, nil
	}
	return nil, err
}

// InitInfluxClient 初始化 Influx 连接
func InitInfluxClient(config string) (c influxdb2.Client, w api.WriteAPIBlocking, err error) {
	conf, err := GetInfluxConfig(config)
	if err != nil {
		return nil, nil, err
	}
	c = influxdb2.NewClient(conf.DataSource, conf.Token)
	w = c.WriteAPIBlocking(conf.Org, conf.Bucket)
	return
}

// WriteSQLLog 写入 SQL 日志
func WriteSQLLog() {
	client, write, err := InitInfluxClient("sqls")
	if err != nil {
		fmt.Println("InitInfluxClient console Failed")
	}
	go func() {
		select {
		default:
			for log := range SQLLogChan {
				fields := map[string]interface{}{"Type": log.Type, "Slow": log.Slow, "Location": log.Location, "Rows": log.Rows, "Cost": log.Cost, "SQL": log.SQL, "Error": log.Error}
				p := influxdb2.NewPoint(log.Type, log.Tag, fields, time.Now())
				err := write.WritePoint(context.Background(), p)
				if err != nil {
					fmt.Println("Write InfluxDB Error :", err)
				}
			}
		case s := <-StopWriteSQLLog:
			if s {
				fmt.Println("Stop Write SQL Log")
				close(SQLLogChan)
				client.Close()
			}
		}
	}()
}
