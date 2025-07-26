package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string           `mapstructure:"environment"`
	GRPC        GRPCConfig       `mapstructure:"grpc"`
	Postgres    PostgresConfig   `mapstructure:"postgres"`
	RabbitMQ    RabbitMQConfig   `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig     `mapstructure:"pubsub"`
	GCPStorage  GCPStorageConfig `mapstructure:"gcp_storage"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type RabbitMQConfig struct {
	DSN          string `mapstructure:"dsn"`
	ExchangeName string `mapstructure:"exchange_name"`
}

type PubSubConfig struct {
	ProjectID string `mapstructure:"project_id"`
}

type GCPStorageConfig struct {
	Bucket          string        `mapstructure:"bucket"`
	CredentialsFile string        `mapstructure:"credentials_file"`
	URLExpiry       time.Duration `mapstructure:"url_expiry"`
	Endpoint        string        `mapstructure:"endpoint"` // For GCS emulator
	ProjectID       string        `mapstructure:"project_id"`
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
		"postgres.host",
		"postgres.port",
		"postgres.user",
		"postgres.password",
		"postgres.dbname",
		"postgres.sslmode",
		"postgres.timezone",
		"rabbitmq.dsn",
		"rabbitmq.exchange_name",
		"pubsub.project_id",
		"gcp_storage.bucket",
		"gcp_storage.credentials_file",
		"gcp_storage.url_expiry",
		"gcp_storage.endpoint",
		"gcp_storage.project_id",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	v.SetDefault("grpc.port", 50053)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "qubool_kallyanam_user")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "UTC")

	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")

	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")

	// GCP Storage defaults
	v.SetDefault("gcp_storage.bucket", "qubool-kallyanam-user-photos")
	v.SetDefault("gcp_storage.url_expiry", 24*time.Hour) // 24 hours
	v.SetDefault("gcp_storage.endpoint", "")             // Empty for production GCS
	v.SetDefault("gcp_storage.project_id", "qubool-kallyanam")
}
