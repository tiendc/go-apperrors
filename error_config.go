package goapperrors

import "errors"

// LogLevel represents log level set for an error.
// You can use LogLevel to report the level of an error to external
// services such as Sentry or Rollbar.
type LogLevel string

const (
	LogLevelNone  LogLevel = ""
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warning"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// ErrorConfig configuration of an error to be used when build error info
type ErrorConfig struct {
	Status   int
	Code     string
	Title    string
	LogLevel LogLevel
	TransKey string
	Extra    any
}

// GetErrorConfig gets global mapping config of an error if set
func GetErrorConfig(err error) *ErrorConfig {
	if err == nil {
		return nil
	}
	if cfg := getValueInErrorMap(err); cfg != nil {
		return cfg
	}
	return GetErrorConfig(errors.Unwrap(err))
}

// getValueInErrorMap returns the value for the error key in the map.
// If the error key is unhashable, getting value from a map will panic.
// In that situation this func will recover from panic and return `zero` value.
func getValueInErrorMap(err error) *ErrorConfig {
	defer func() {
		_ = recover()
	}()
	if cfg, ok := mapError[err]; ok {
		return cfg
	}
	return nil
}
