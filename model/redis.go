package model

type RedisMapConfig struct {
	List map[string]*RedisConfig `mapstructure:"list"`
}

type RedisConfig struct {
	ProxyList    string `mapstructure:"proxy_list"`
	Password     string `mapstructure:"password"`
	Prefix       string `mapstructure:"prefix"`
	Db           int    `mapstructure:"db"`
	MaxIdle      int    `mapstructure:"max_idle"`
	MaxActive    int    `mapstructure:"max_active"`
	ConnTimeout  int    `mapstructure:"conn_timeout"`
	IdelTimeout  int    `mapstructure:"idle_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}
