package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	GRPC        GRPCConfig     `mapstructure:"grpc"`
	RabbitMQ    RabbitMQConfig `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig   `mapstructure:"pubsub"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type RabbitMQConfig struct {
	DSN          string `mapstructure:"dsn"`
	ExchangeName string `mapstructure:"exchange_name"`
}

type PubSubConfig struct {
	ProjectID string `mapstructure:"project_id"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := initViper(configPath)
	bindEnvVars(v)
	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func initViper(path string) *viper.Viper {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if path != "" {
		v.SetConfigFile(path)
		v.SetConfigType("yaml")
		_ = v.ReadInConfig() // Non-fatal if missing
	}

	return v
}

func bindEnvVars(v *viper.Viper) {
	keys := []string{
		"environment",
		"grpc.port",
		"rabbitmq.dsn",
		"rabbitmq.exchange_name",
		"pubsub.project_id",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")
	v.SetDefault("grpc.port", 50052)
	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")
	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")
}
