package goapperrors

import "errors"

// LogLevel represents log level set for an error
type LogLevel int8

const (
	LogLevelNone  LogLevel = iota
	LogLevelDebug LogLevel = iota
	LogLevelInfo  LogLevel = iota
	LogLevelWarn  LogLevel = iota
	LogLevelError LogLevel = iota
	LogLevelFatal LogLevel = iota
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
