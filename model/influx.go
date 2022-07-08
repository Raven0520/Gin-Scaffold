package model

// InfluxConfig Influx 配置
type InfluxConfig struct {
	DriverName string `mapstructure:"driver_name"` // 驱动名称
	DataSource string `mapstructure:"data_source"` // 数据节点
	Org        string `mapstructure:"org"`         // 组织
	Bucket     string `mapstructure:"bucket"`      // 存储桶
	Token      string `mapstructure:"token"`       // 令牌
}

// InfluxConfigs 配置列表
type InfluxConfigs struct {
	List map[string]*InfluxConfig `mapstructure:"list"`
}

// SQLLog SQL 日志
type SQLLog struct {
	Tag      map[string]string `json:"tag"`
	Slow     bool              `json:"slow"`
	Type     string            `json:"type"`
	Location string            `json:"location"`
	Rows     uint64            `json:"rows"`
	Cost     float64           `json:"cost"`
	SQL      string            `json:"sql"`
	Error    string            `json:"error"`
}
