package config

import (
	"github.com/spf13/viper"
	"time"
)

type Log struct {
	Path      string `mapstructure:"path"`
	Level     string `mapstructure:"level"`
	MaxSize   int    `mapstructure:"maxsize"`
	MaxBackup int    `mapstructure:"maxbackup"`
	MaxAge    int    `mapstructure:"maxage"`
	Compress  bool   `mapstructure:"compress"`
}

type App struct {
	Salt         string `mapstructure:"salt" json:"-" yaml:"salt"`
	PageSize     int    `mapstructure:"page-size" json:"PageSize" yaml:"page-size"`
	JwtTokenName string `mapstructure:"jwt-token-name" json:"JwtTokenName" yaml:"jwt-token-name"`
	JwtSecret    string `mapstructure:"jwt-secret" json:"JwtSecret" yaml:"jwt-secret"`
	AdminSecret  string `mapstructure:"admin-secret" json:"-" yaml:"jwt-secret"`
}

type Server struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout"`
}

type DataBase struct {
	Type        string `mapstructure:"type" json:"type" yaml:"type"`
	Host        string `mapstructure:"host" json:"host" yaml:"host"`
	User        string `mapstructure:"user" json:"user" yaml:"user"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	Name        string `mapstructure:"name" json:"name" yaml:"name"`
	TablePrefix string `mapstructure:"table-prefix" json:"TablePrefix" yaml:"table-prefix"`
}

type Elasticsearch struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Timeout  string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Index    string `mapstructure:"index" json:"index" yaml:"index"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Config struct {
	RunMode string  `mapstructure:"run_mode"`
	Log     *Log    `mapstructure:"log"`
	Server  *Server `mapstructure:"server"`
}

var config *Config

func ParseConfigFile(file string) (*Config, error) {
	config = new(Config)

	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
