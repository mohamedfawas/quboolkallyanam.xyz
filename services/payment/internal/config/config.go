package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	GRPC        GRPCConfig     `mapstructure:"grpc"`
	Postgres    PostgresConfig `mapstructure:"postgres"`
	Razorpay    RazorpayConfig `mapstructure:"razorpay"`
	Email       EmailConfig    `mapstructure:"email"`
	RabbitMQ    RabbitMQConfig `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig   `mapstructure:"pubsub"`
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

type RazorpayConfig struct {
	KeyID     string `mapstructure:"key_id"`
	KeySecret string `mapstructure:"key_secret"`
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
		"razorpay.key_id",
		"razorpay.key_secret",
		"email.smtp_host",
		"email.smtp_port",
		"email.smtp_username",
		"email.smtp_password",
		"email.from_email",
		"email.from_name",
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

	v.SetDefault("grpc.port", 50055)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "qubool_kallyanam_payment")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "UTC")

	v.SetDefault("razorpay.key_id", "")
	v.SetDefault("razorpay.key_secret", "")

	v.SetDefault("email.smtp_host", "smtp.gmail.com")
	v.SetDefault("email.smtp_port", 587)
	v.SetDefault("email.smtp_username", "")
	v.SetDefault("email.smtp_password", "")
	v.SetDefault("email.from_email", "noreply@qubool-kallyanam.xyz")
	v.SetDefault("email.from_name", "Qubool Kallyanam")

	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")

	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")
}
