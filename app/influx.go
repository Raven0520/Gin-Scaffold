package app

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gitlab.f-fans.cn/components/golang/scaffold/model"
	"gitlab.f-fans.cn/components/golang/scaffold/util"
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
func InitInfluxClient(config string) (c influxdb2.Client, err error) {
	conf, err := GetInfluxConfig(config)
	if err != nil {
		return nil, err
	}
	return influxdb2.NewClient(conf.DataSource, conf.Token), nil
}

// WritePrint 写入 Print
func WritePrint(config, measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	conf, err := GetInfluxConfig(config)
	if err != nil {
		fmt.Println("InfluxDB GetInfluxConfig Error :", err)
		return err
	}
	client := influxdb2.NewClient(conf.DataSource, conf.Token)
	write := client.WriteAPIBlocking(conf.Org, conf.Bucket)
	if err != nil {
		return err
	}
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	err = write.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Println("InfluxDB WritePoint Error :", err)
		client.Close()
		return err
	}
	defer client.Close()
	return
}

// WriteErrorPrint 写入 错误日志
func WriteErrorPrint(measurement string, tags map[string]string, fields map[string]interface{}) (err error) {
	return WritePrint("errors", measurement, tags, fields)
}
