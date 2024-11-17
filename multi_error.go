package goapperrors

import (
	"errors"
	"strings"
)

type MultiError interface {
	AppError

	InnerErrors() AppErrors
}

type DefaultMultiError struct {
	*DefaultAppError
}

// Unwrap - implementation used by errors.Is
func (e *DefaultMultiError) Unwrap() []error {
	return e.InnerErrors().Unwrap()
}

// InnerErrors returns all wrapped errors as AppErrors
func (e *DefaultMultiError) InnerErrors() AppErrors {
	err := e.err
	for {
		if inErrs, ok := err.(AppErrors); ok { //nolint:errorlint
			return inErrs
		}
		if err = errors.Unwrap(err); err == nil {
			return nil
		}
	}
}

// WithParam - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithParam(k string, v any) AppError {
	_ = e.DefaultAppError.WithParam(k, v)
	return e
}

// WithTransParam - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithTransParam(k string, v string) AppError {
	_ = e.DefaultAppError.WithTransParam(k, v)
	return e
}

// WithCause - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithCause(cause error) AppError {
	_ = e.DefaultAppError.WithCause(cause)
	return e
}

// WithDebug - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithDebug(format string, args ...any) AppError {
	_ = e.DefaultAppError.WithDebug(format, args...)
	return e
}

// WithCustomConfig - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithCustomConfig(cfg *ErrorConfig) AppError {
	_ = e.DefaultAppError.WithCustomConfig(cfg)
	return e
}

// WithCustomBuilder - re-defines to make sure the returning points to this error object
func (e *DefaultMultiError) WithCustomBuilder(infoBuilder InfoBuilderFunc) AppError {
	_ = e.DefaultAppError.WithCustomBuilder(infoBuilder)
	return e
}

// Build implements Build function
func (e *DefaultMultiError) Build(lang Language, options ...InfoBuilderOption) *InfoBuilderResult {
	buildCfg := e.BuildConfig(lang, options...)
	buildResult := e.build(buildCfg)
	errInfo := buildResult.ErrorInfo

	inErrs := e.InnerErrors()
	errInfo.InnerErrors = make([]*ErrorInfo, 0, len(inErrs))
	for _, inErr := range inErrs {
		inResult := inErr.Build(lang)
		buildResult.TransMissingKeys = append(buildResult.TransMissingKeys, inResult.TransMissingKeys...)
		errInfo.InnerErrors = append(errInfo.InnerErrors, inResult.ErrorInfo)
	}

	if buildResult.TransMissingMainKey {
		var sb strings.Builder
		for i, info := range errInfo.InnerErrors {
			if i > 0 {
				sb.WriteString(buildCfg.ErrorSeparator)
			}
			sb.WriteString(info.Message)
		}
		errInfo.Message = sb.String()
	}

	return buildResult
}

// NewMultiError creates a MultiError with wrapping the given errors
func NewMultiError(errs ...AppError) MultiError {
	if len(errs) == 0 {
		return nil
	}
	e := &DefaultMultiError{
		DefaultAppError: newDefaultAppError(AppErrors(errs)),
	}
	// MultiError does not allow using global config mapping
	// If you need to set custom config, sets via the field `customConfig`
	e.disallowGlobalConfigMapping = true
	return e
}

// AsMultiError converts AppError to MultiError
func AsMultiError(err AppError) MultiError {
	e, _ := err.(MultiError)
	return e
}
