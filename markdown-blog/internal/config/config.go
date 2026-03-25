package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig    `mapstructure:"server"`
	Site     SiteConfig      `mapstructure:"site"`
	Database DatabaseConfig  `mapstructure:"database"`
	Redis    RedisConfig     `mapstructure:"redis"`
	Log      LogConfig       `mapstructure:"log"`
	Github   GithubConfig    `mapstructure:"github"`
	Content  ContentConfig   `mapstructure:"content"`
	JWT      JWTConfig       `mapstructure:"jwt"`
}

type SiteConfig struct {
	Title       string `mapstructure:"title"`
	Author      string `mapstructure:"author"`
	Description string `mapstructure:"description"`
	Beian       string `mapstructure:"beian"`
	Keywords    string `mapstructure:"keywords"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
}

type GithubConfig struct {
	ClientID    string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL string `mapstructure:"callback_url"`
}

type ContentConfig struct {
	Path   string `mapstructure:"path"`
	BaseURL string `mapstructure:"base_url"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int64  `mapstructure:"expiration"`
	Issuer     string `mapstructure:"issuer"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Enable environment variable override
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// GetDSN returns MySQL DSN
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// GetRedisAddr returns Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
