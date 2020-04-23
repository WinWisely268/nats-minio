package log

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

// Log instantiates new zero logger
type Logger *zerolog.Logger

// NewLogger instantiates new zerolog
func NewLogger(cmd string, pkgName string) *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	z := zerolog.New(output).With().Str(cmd, pkgName).Timestamp().Logger()
	return Logger(z)
}
