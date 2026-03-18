package logging

import (
	"strings"

	"go.uber.org/zap"
)

func NewLogger(env string) (*zap.Logger, error) {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "", "local", "dev", "development":
		return zap.NewDevelopment()
	default:
		return zap.NewProduction()
	}
}
