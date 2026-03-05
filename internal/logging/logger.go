package logging

import (
	"strings"

	"go.uber.org/zap"
)

// NewLogger creates a zap logger based on environment (local/dev => development logger).
func NewLogger(env string) (*zap.Logger, error) {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "", "local", "dev", "development":
		return zap.NewDevelopment()
	default:
		return zap.NewProduction()
	}
}
