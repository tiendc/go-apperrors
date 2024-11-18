package goapperrors

import (
	"errors"

	goerrors "github.com/go-errors/errors"
)

// Init initializes global config
func Init(cfg *Config) {
	if cfg == nil {
		panic("config must not be nil")
	}
	cfg.setDefault()
	globalConfig = cfg
	// Apply max stack depth to `go-errors` settings
	goerrors.MaxStackDepth = globalConfig.MaxStackDepth
}

// Add adds a global config mapping for a base error, then returns the error.
// This function is recommended for adding mapping for app-external errors.
//
// Example:
//
//	 var ErrRedisKeyNotFound = Add(redis.Nil, &ErrorConfig{
//		      Status: http.StatusNotFound,
//		      Code: "ErrRedisKeyNotFound",
//	 })
func Add(err error, cfg *ErrorConfig) error {
	if err == nil || cfg == nil {
		panic("error and config must not be nil")
	}
	if cfg.Code == "" {
		cfg.Code = err.Error()
	}
	if cfg.TransKey == "" {
		cfg.TransKey = cfg.Code
	}
	mapError[err] = cfg
	return err
}

// Create creates an error for the code with the mapping config, then returns the newly created error.
// This function is recommended for adding mapping for app-internal errors. When use this method,
// you don't need to set custom error code or translation key, they will be the same as the input code.
//
// Example:
//
//	var ErrTokenInvalid = Create("ErrTokenInvalid", &ErrorConfig{Status: http.StatusUnauthorized})
//	var ErrNoPermission = Create("ErrNoPermission", &ErrorConfig{Status: http.StatusForbidden})
func Create(code string, cfg *ErrorConfig) error {
	if code == "" || cfg == nil {
		panic("error key and config must not be nil")
	}
	err := errors.New(code) //nolint:err113
	if cfg.Code == "" {
		cfg.Code = code
	}
	if cfg.TransKey == "" {
		cfg.TransKey = cfg.Code
	}
	mapError[err] = cfg
	return err
}

// Remove removes the error from the global config mappings
func Remove(err error) {
	if err != nil {
		delete(mapError, err)
	}
}

// Build builds error info
func Build(err error, lang Language, options ...InfoBuilderOption) *InfoBuilderResult {
	anErr := err
	for {
		if builder, ok := anErr.(interface {
			Build(Language, ...InfoBuilderOption) *InfoBuilderResult
		}); ok {
			return builder.Build(lang, options...)
		}
		if anErr = errors.Unwrap(anErr); anErr == nil {
			break
		}
	}
	return New(err).Build(lang, options...)
}
