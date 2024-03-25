package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	Logger Logger
	Flood  FloodConfig
}

type ServerConfig struct {
	AppVersion        string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
}

type Logger struct {
	Development bool
	Encoding    string
	Level       string
}

type RedisConfig struct {
	Addr         string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  time.Duration
	Password     string
	DB           int
}

type FloodConfig struct {
	TimeLimit   time.Duration
	MaxRequests int
}

// LoadConfig Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// ParseConfig Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}
	c.normalizeTime()

	return &c, nil
}

func (c *Config) normalizeTime() {
	c.Server.CtxDefaultTimeout *= time.Second
	c.Server.WriteTimeout *= time.Second
	c.Server.ReadTimeout *= time.Second
	c.Redis.PoolTimeout *= time.Second
	c.Flood.TimeLimit *= time.Second
}
