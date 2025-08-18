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
