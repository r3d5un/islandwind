package config

import (
	"context"
	"log/slog"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `json:"app"`
	Server ServerConfig `json:"server"`
	DB     db.Config    `json:"db"`
}

type AppConfig struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

// ServerConfig is the configuration used when setting up http.ServeMux
type ServerConfig struct {
	// Port refers to which TCP port the server should use
	Port int `json:"port"`
	// IdleTimeout is the number of seconds to wait for the next request when keep-alives are enabled
	IdleTimeout int `json:"idleTimeout"`
	// ReadTimeout is the maximum number of seconds for reading the entire request.
	ReadTimeout int `json:"readTimeout"`
	// WriteTimeout is the maximum duration before timing out the writes of the response.
	WriteTimeout int `json:"writeTimeout"`
	// BasicAuthConfig is used to set the administrator username and password
	Authentication api.BasicAuthConfig `json:"authentication"`
}

func New() (*Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/islandwind")

	viper.SetDefault("app.name", "islandwind")
	viper.SetDefault("app.environment", "development")
	// Default Server Settings
	viper.SetDefault("server.port", 4000)
	viper.SetDefault("server.idleTimeout", 60)
	viper.SetDefault("server.readTimeout", 5)
	viper.SetDefault("server.writeTimeout", 10)
	viper.SetDefault("server.authentication.username", "islandwind")
	viper.SetDefault("server.authentication.password", "islandwind")
	// Default Database Settings
	viper.SetDefault(
		"db.connstr",
		"postgres://postgres:postgres@localhost:5432/islandwind?sslmode=disable",
	)
	viper.SetDefault("db.maxOpenConns", 15)
	viper.SetDefault("db.idleTimeMinutes", 5)
	viper.SetDefault("db.timeoutSeconds", 5)

	viper.SafeWriteConfig()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("islandwind")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.LogAttrs(
			context.Background(), slog.LevelInfo, "configuration updated", slog.Any("even", e),
		)

		viper.Unmarshal(&config)
	})
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		return &config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}
