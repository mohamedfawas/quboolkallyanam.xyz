package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string          `mapstructure:"environment"`
	GRPC        GRPCConfig      `mapstructure:"grpc"`
	Postgres    PostgresConfig  `mapstructure:"postgres"`
	RabbitMQ    RabbitMQConfig  `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig    `mapstructure:"pubsub"`
	Firestore   FirestoreConfig `mapstructure:"firestore"`
	MongoDB     MongoDBConfig   `mapstructure:"mongodb"`
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

type FirestoreConfig struct {
	ProjectID string `mapstructure:"project_id"`
}

type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
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

	v.SetEnvPrefix("CHAT")
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
		"firestore.project_id",
		"mongodb.uri",
		"mongodb.database",
		"mongodb.timeout",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	v.SetDefault("grpc.port", 50054)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "qubool_kallyanam_chat")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "UTC")

	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")

	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")

	v.SetDefault("firestore.project_id", "qubool-kallyanam-chat")

	v.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("mongodb.database", "quboolKallyanam")
	v.SetDefault("mongodb.timeout", 10*time.Second)
}
