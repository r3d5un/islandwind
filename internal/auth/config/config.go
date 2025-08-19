package config

import (
	"log/slog"
)

type Config struct {
	SigningSecret string `json:"signingSecret"`
	TokenIssuer   string `json:"tokenIssuer"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(slog.String("signingSecret", "omitted"))
}

// BasicAuthConfig contains the username and password used in the basic authentication
// for the HTTP server.
type BasicAuthConfig struct {
	// Username is the admin username.
	//
	// Set through the ISLANDWIND_SERVER_AUTHENTICATION_USERNAME environment variable.
	Username string `json:"username"`
	// Password is the admin password.
	//
	// Field is safe for logging as the [BasicAuthConfig] contains a custom [BasicAuthConfig.LogValue] method.
	//
	// Set through the ISLANDWIND_SERVER_AUTHENTICATION_PASSWORD environment variable.
	Password string `json:"password"`
}

func (c BasicAuthConfig) LogValue() slog.Value {
	return slog.GroupValue(slog.String("username", c.Username))
}
