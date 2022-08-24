package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig `mapstructure:"ServerConfig"`
	DBConfig     DBConfig     `mapstructure:"DBConfig"`
	JWTConfig    JWTConfig    `mapstructure:"JWTConfig"`
	Logger       Logger       `mapstructure:"Logger"`
}

type ServerConfig struct {
	AppVersion       string `mapstructure:"AppVersion"`
	Mode             string `mapstructure:"Mode"`
	RoutePrefix      string `mapstructure:"RoutePrefix"`
	Debug            bool   `mapstructure:"Debug"`
	Port             string `mapstructure:"Port"`
	TimeoutSecs      int64  `mapstructure:"TimeoutSecs"`
	ReadTimeoutSecs  int64  `mapstructure:"ReadTimeoutSecs"`
	WriteTimeoutSecs int64  `mapstructure:"WriteTimeoutSecs"`
}

type DBConfig struct {
	DataSourceName string `mapstructure:"DataSourceName"`
	MaxOpen        int    `mapstructure:"MaxOpen"`
	MaxIdle        int    `mapstructure:"MaxIdle"`
	MaxLifetime    int    `mapstructure:"MaxLifeTime"`
}

type JWTConfig struct {
	AccessSessionTime  int    `mapstructure:"AccessSessionTime"`
	RefreshSessionTime int    `mapstructure:"RefreshSessionTime"`
	SecretKey          string `mapstructure:"SecretKey"`
}

type Logger struct {
	Development bool   `mapstructure:"Development"`
	Encoding    string `mapstructure:"Encoding"`
	Level       string `mapstructure:"Level"`
}

// LoadConfig sets viper settings
func LoadConfig() (*Config, error) {
	vp := viper.New()
	vp.SetConfigName("../pkg/base/config/configLocal")
	vp.SetConfigType("json")
	// AddConfigPath adds a path for Viper to search for the config file in.
	vp.AddConfigPath(".")
	// AutomaticEnv makes Viper check if environment variables match any of the existing keys
	// (config, default or flags). If matching env vars are found, they are loaded into Viper.
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		if res, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New(res.Error())
		}
		return nil, err
	}
	var c Config
	err := vp.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
