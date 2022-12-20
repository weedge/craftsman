package gorm

import (
	"time"
)

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	kvLogger IkvLogger
	// slowThreshold slow sql exec time threshold
	slowThreshold time.Duration
	// traceLogLevel just for gorm trace log map kitex log level to print
	traceLogLevel int
}

// defaultConfig default config
func defaultConfig() *config {
	return &config{
		kvLogger:      DefaultLogger,
		slowThreshold: 200 * time.Millisecond,
		traceLogLevel: int(Debug),
	}
}

// WithKvLogger IkvLogger impl
func WithKvLogger(kvLogger IkvLogger) Option {
	return option(func(cfg *config) {
		cfg.kvLogger = kvLogger
	})
}

// WithSlowThreshold slow sql exec time threshold
func WithSlowThreshold(slowThreshold time.Duration) Option {
	return option(func(cfg *config) {
		cfg.slowThreshold = slowThreshold
	})
}

// WithTraceLogLevel just for gorm trace log map kitex log level to print
func WithTraceLogLevel(traceLogLevel int) Option {
	return option(func(cfg *config) {
		cfg.traceLogLevel = traceLogLevel
	})
}
