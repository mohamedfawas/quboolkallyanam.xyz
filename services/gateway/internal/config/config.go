package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	HTTP        HTTPConfig     `mapstructure:"http"`
	Services    ServicesConfig `mapstructure:"services"`
	Auth        AuthConfig     `mapstructure:"auth"`
	// Tracing config @TODO: add tracing config
}

type HTTPConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

type ServicesConfig struct {
	AuthServicePort    string `mapstructure:"auth_port"`
	AdminServicePort   string `mapstructure:"admin_port"`
	UserServicePort    string `mapstructure:"user_port"`
	ChatServicePort    string `mapstructure:"chat_port"`
	PaymentServicePort string `mapstructure:"payment_port"`
}

type AuthConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
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
		"http.port",
		"http.read_timeout",
		"http.write_timeout",
		"http.idle_timeout",
		"services.auth_port",
		"services.admin_port",
		"services.user_port",
		"services.chat_port",
		"services.payment_port",
		"auth.jwt.secret_key",
		"auth.jwt.access_token_minutes",
		"auth.jwt.refresh_token_days",
		"auth.jwt.issuer",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	v.SetDefault("http.port", "8080")
	v.SetDefault("http.read_timeout", 10)
	v.SetDefault("http.write_timeout", 10)
	v.SetDefault("http.idle_timeout", 60)

	v.SetDefault("services.auth_port", "5001")
	v.SetDefault("services.admin_port", "5002")
	v.SetDefault("services.user_port", "5003")
	v.SetDefault("services.chat_port", "5004")
	v.SetDefault("services.payment_port", "5005")

	v.SetDefault("auth.jwt.secret_key", "defaultsecret")
	v.SetDefault("auth.jwt.access_token_minutes", 15)
	v.SetDefault("auth.jwt.refresh_token_days", 7)
	v.SetDefault("auth.jwt.issuer", "qubool-kallyanam")
}
