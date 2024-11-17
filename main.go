package goapperrors

import (
	"errors"
	"fmt"
	"runtime"

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
//		var ErrRedisKeyNotFound = Add(redis.Nil, &ErrorInfo{
//	       Status: http.StatusNotFound,
//	       Code: "ErrRedisKeyNotFound",
//	    })
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
// var ErrTokenInvalid = Create("ErrTokenInvalid", &ErrorInfo{Status: http.StatusUnauthorized})
// var ErrNoPermission = Create("ErrNoPermission", &ErrorInfo{Status: http.StatusForbidden})
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
	return NewAppError(err).Build(lang, options...)
}

// Wrap wraps an error with adding stack trace if configured
func Wrap(err error) error {
	if globalConfig.WrapFunc != nil {
		return globalConfig.WrapFunc(err)
	}
	return goerrors.Wrap(err, 1)
}

// Wrapf wraps an error by calling fmt.Errorf and adds stack trace if configured
func Wrapf(format string, args ...any) error {
	return Wrap(fmt.Errorf(format, args...)) //nolint:err113
}

// GetStackTrace gets stack trace stored in the error if there is
func GetStackTrace(err error) []runtime.Frame {
	if err == nil {
		return nil
	}
	// If it is `go-errors`.*Error, get the stack trace from it
	var gErr *goerrors.Error
	if !errors.As(err, &gErr) {
		return nil
	}
	callers := gErr.Callers()
	frames := runtime.CallersFrames(callers)
	frameList := make([]runtime.Frame, 0, len(callers))
	for {
		next, hasMore := frames.Next()
		frameList = append(frameList, next)
		if !hasMore {
			break
		}
	}
	return frameList
}
