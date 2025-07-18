package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	GRPC        GRPCConfig     `mapstructure:"grpc"`
	Postgres    PostgresConfig `mapstructure:"postgres"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Admin       AdminConfig    `mapstructure:"admin"`
	RabbitMQ    RabbitMQConfig `mapstructure:"rabbitmq"`
	PubSub      PubSubConfig   `mapstructure:"pubsub"`
	Auth        AuthConfig     `mapstructure:"auth"`
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

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AdminConfig struct {
	DefaultAdminEmail    string `mapstructure:"default_admin_email"`
	DefaultAdminPassword string `mapstructure:"default_admin_password"`
}

type RabbitMQConfig struct {
	DSN          string `mapstructure:"dsn"`
	ExchangeName string `mapstructure:"exchange_name"`
}

type PubSubConfig struct {
	ProjectID string `mapstructure:"project_id"`
}

type AuthConfig struct {
	PendingRegistrationExpiryHours int       `mapstructure:"pending_registration_expiry_hours"`
	OTPExpiryMinutes               int       `mapstructure:"otp_expiry_minutes"`
	JWT                            JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	SecretKey          string `mapstructure:"secret_key"`
	AccessTokenMinutes int    `mapstructure:"access_token_minutes"`
	RefreshTokenDays   int    `mapstructure:"refresh_token_days"`
	Issuer             string `mapstructure:"issuer"`
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
		"auth.pending_registration_expiry_hours",
		"auth.otp_expiry_minutes",
		"auth.jwt.secret_key",
		"auth.jwt.access_token_minutes",
		"auth.jwt.refresh_token_days",
		"auth.jwt.issuer",
		"postgres.host",
		"postgres.port",
		"postgres.user",
		"postgres.password",
		"postgres.dbname",
		"postgres.sslmode",
		"postgres.timezone",
		"redis.host",
		"redis.port",
		"redis.password",
		"redis.db",
		"admin.default_admin_email",
		"admin.default_admin_password",
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

	v.SetDefault("grpc.port", 50051)

	v.SetDefault("auth.pending_registration_expiry_hours", 1)
	v.SetDefault("auth.otp_expiry_minutes", 15)
	v.SetDefault("auth.jwt.secret_key", "your-256-bit-secret-replace-in-production")
	v.SetDefault("auth.jwt.access_token_minutes", 15)
	v.SetDefault("auth.jwt.refresh_token_days", 7)
	v.SetDefault("auth.jwt.issuer", "qubool-kallyanam")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "qubool_kallyanam_auth")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.timezone", "UTC")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("admin.default_admin_email", "adminquboolkallyanam@gmail.com")
	v.SetDefault("admin.default_admin_password", "Admin@123")

	v.SetDefault("rabbitmq.dsn", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange_name", "qubool_kallyanam_events")

	v.SetDefault("pubsub.project_id", "qubool-kallyanam-events")
}
