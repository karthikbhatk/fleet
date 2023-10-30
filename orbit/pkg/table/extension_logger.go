package table

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is a wrapper around zerolog, which we use for tables
// using the go-kit logger
type Logger struct {
	zerolog.Logger
}

// Fatal logs a fatal message, just like go-kit's logutil.Fatal.
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal().Msg(fmt.Sprint(args...))
}

// Log logs a message, implementing log.Logger interface
func (l *Logger) Log(keyValuePairs ...interface{}) error {
	log.Logger.Info().Msg(fmt.Sprint(keyValuePairs...))
	return nil
}

// SetLevelKey sets the level key, just like go-kit's logutil.SetLevelKey.
func (l *Logger) SetLevelKey(key interface{}) *Logger {
	// Note: zerolog doesn't support changing the level key dynamically like logutil, so we just log it as a warning.
	l.Logger.Warn().Msgf("Attempted to set level key to: %v", key)
	return l
}

// NewOsqueryLogger returns the Logger struct.
func NewOsqueryLogger() *Logger {
	// Return a Logger struct with our global logger, and use the existing global config for the log level.
	return &Logger{log.Logger}
}
