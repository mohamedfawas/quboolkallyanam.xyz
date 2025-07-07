package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string          `mapstructure:"environment"`
	GRPC        GRPCConfig      `mapstructure:"grpc"`
	Email       EmailConfig     `mapstructure:"email"`
	RabbitMQ    RabbitMQConfig  `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig    `mapstructure:"pubsub"`
	Firestore   FirestoreConfig `mapstructure:"firestore"`
	MongoDB     MongoDBConfig   `mapstructure:"mongodb"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
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
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
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
		"email.smtp_host",
		"email.smtp_port",
		"email.smtp_username",
		"email.smtp_password",
		"email.from_email",
		"email.from_name",
		"rabbitmq.dsn",
		"rabbitmq.exchange_name",
		"pubsub.project_id",
		"firestore.project_id",
		"mongodb.uri",
		"mongodb.database",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	v.SetDefault("grpc.port", 50054)

	v.SetDefault("email.smtp_host", "smtp.gmail.com")
	v.SetDefault("email.smtp_port", 587)
	v.SetDefault("email.smtp_username", "")
	v.SetDefault("email.smtp_password", "")
	v.SetDefault("email.from_email", "noreply@qubool-kallyanam.xyz")
	v.SetDefault("email.from_name", "Qubool Kallyanam")

	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")

	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")

	v.SetDefault("firestore.project_id", "qubool-kallyanam-chat")

	v.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("mongodb.database", "qubool_kallyanam_chat")
}
