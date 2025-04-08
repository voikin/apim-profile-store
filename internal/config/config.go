package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/voikin/apim-profile-store/pkg/logger"
)

type Config struct {
	Logger   *logger.Config  `mapstructure:"logger"`
	Server   *ServerConfig   `mapstructure:"server"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	Neo4J    *Neo4JConfig    `mapstructure:"neo4j"`
}

type ServerConfig struct {
	GRPC GRPCConfig `mapstructure:"grpc"`
	HTTP HTTPConfig `mapstructure:"http"`
}

type GRPCConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	MaxConnAgeSeconds int    `mapstructure:"max_conn_age_seconds"`
}

type HTTPConfig struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	ReadTimeoutSecs       int    `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSecs      int    `mapstructure:"write_timeout_seconds"`
	ReadHeaderTimeoutSecs int    `mapstructure:"read_header_timeout_seconds"`
}

type PostgresConfig struct {
	DSN string `mapstructure:"dsn"`
}

type Neo4JConfig struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c GRPCConfig) MaxConnectionAge() time.Duration {
	return time.Duration(c.MaxConnAgeSeconds) * time.Second
}

func (c HTTPConfig) ReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeoutSecs) * time.Second
}

func (c HTTPConfig) WriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeoutSecs) * time.Second
}

func (c HTTPConfig) ReadHeaderTimeout() time.Duration {
	return time.Duration(c.ReadTimeoutSecs) * time.Second
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return &Config{}, fmt.Errorf("viper.ReadInConfig: %w", err)
	}

	config := &Config{}

	err = viper.Unmarshal(config)
	if err != nil {
		return &Config{}, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	return config, nil
}
