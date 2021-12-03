package config

import (
	"github.com/spf13/viper"
	"time"
)

type Log struct {
	Filename  string `mapstructure:"filename"`
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
	Type        string `mapstructure:"type"`
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	TablePrefix string `mapstructure:"table-prefix"`
}

type Rbac struct {
	Model  string `mapstructure:"model"`
	Policy string `mapstructure:"policy"`
	Auth   string `mapstructure:"auth"`
}

type Elasticsearch struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Timeout  string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Index    string `mapstructure:"index" json:"index" yaml:"index"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Config struct {
	RunMode  string    `mapstructure:"run_mode"`
	Log      *Log      `mapstructure:"log"`
	Server   *Server   `mapstructure:"server"`
	Database *DataBase `mapstructer:"database"`
	Rbac     *Rbac     `mapstructer:"rbac"`
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
