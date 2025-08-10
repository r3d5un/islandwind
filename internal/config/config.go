package config

import (
	"context"
	"log/slog"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `json:"app"`
}

type AppConfig struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

func New() (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/islandwind")

	viper.SetDefault("app.name", "islandwind")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("server.port", 4000)

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
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
