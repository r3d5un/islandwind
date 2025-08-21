package config

import (
	"strings"

	"github.com/r3d5un/islandwind/internal/auth/config"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/spf13/viper"
)

// Config contains the configuration loaded at startup of the application. The configuration
// does not reload environment variables; you will have to use the [viper] package directly
// where you need to read updated configuration variables.
type Config struct {
	App       AppConfig              `json:"app"`
	Server    ServerConfig           `json:"server"`
	DB        db.Config              `json:"db"`
	Auth      config.Config          `json:"auth"`
	BasicAuth config.BasicAuthConfig `json:"basicAuty"`
}

// AppConfig contains the most top-level configuration for the application.
type AppConfig struct {
	// Name is the name of the application
	//
	// Set through the ISLANDWIND_APP_NAME environment variable
	Name string `json:"name"`
	// Environment denotes which environment the application instance is running on
	//
	// Set through the ISLANDWIND_APP_ENVIRONMENT environment variable
	Environment string `json:"environment"`
}

// ServerConfig is the configuration used when setting up http.ServeMux
type ServerConfig struct {
	// Port refers to which TCP port the server should use
	//
	// Set through the ISLANDWIND_SERVER_PORT environment variable
	Port int `json:"port"`
	// IdleTimeout is the number of seconds to wait for the next request when keep-alive are enabled
	//
	// Set through the ISLANDWIND_SERVER_IDLETIMEOUT environment variable
	IdleTimeout int `json:"idleTimeout"`
	// ReadTimeout is the maximum number of seconds for reading the entire request.
	//
	// Set through the ISLANDWIND_SERVER_READTIMEOUT environment variable
	ReadTimeout int `json:"readTimeout"`
	// WriteTimeout is the maximum duration before timing out the writes of the response.
	//
	// Set through the ISLANDWIND_SERVER_WRITETIMEOUT environment variable
	WriteTimeout int `json:"writeTimeout"`
}

func New() (*Config, error) {
	var cfg Config

	viper.SetDefault("app.name", "islandwind")
	viper.SetDefault("app.environment", "development")
	// Default Server Settings
	viper.SetDefault("server.port", 4000)
	viper.SetDefault("server.idleTimeout", 60)
	viper.SetDefault("server.readTimeout", 5)
	viper.SetDefault("server.writeTimeout", 10)
	// Authentication
	viper.SetDefault("basicauth.username", "islandwind")
	viper.SetDefault("basicauth.password", "islandwind")
	viper.SetDefault("auth.signingsecret", "accessTokenSecret")
	viper.SetDefault("auth.refreshsigningsecret", "refreshTokenSecret")
	viper.SetDefault("auth.tokenissuer", "islandwind")
	// Default Database Settings
	viper.SetDefault(
		"db.connstr",
		"postgres://postgres:postgres@localhost:5432/islandwind?sslmode=disable",
	)
	viper.SetDefault("db.maxOpenConns", 15)
	viper.SetDefault("db.idleTimeMinutes", 5)
	viper.SetDefault("db.timeoutSeconds", 5)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("islandwind")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
