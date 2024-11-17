package goapperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is designed to be used as base error type for any error in an application.
// An AppError can carry much extra information such as `cause`, `debug log`, and
// stack trace. It also supports translating the message into a specific language.
type AppError interface {
	error

	// Params gets non-translating param map
	Params() map[string]any
	// TransParams gets translating param map (values of this map will be translated)
	TransParams() map[string]string
	// Cause gets cause of the error
	Cause() error
	// Debug gets debug message
	Debug() string
	// Config returns the custom config if set, otherwise returns the global mapping one
	Config() *ErrorConfig
	// CustomConfig gets custom config associated with the error
	CustomConfig() *ErrorConfig
	// CustomBuilder gets custom info builder
	CustomBuilder() InfoBuilderFunc

	// WithParam sets a custom param
	WithParam(k string, v any) AppError
	// WithTransParam sets a custom param with value to be translated when build info
	WithTransParam(k string, v string) AppError
	// WithCause sets cause of the error
	WithCause(err error) AppError
	// WithDebug sets debug message (used for debug purpose)
	WithDebug(format string, args ...any) AppError
	// WithCustomConfig sets custom config for the error
	WithCustomConfig(*ErrorConfig) AppError
	// WithCustomBuilder sets custom info builder
	WithCustomBuilder(InfoBuilderFunc) AppError

	// Build builds error info
	Build(Language, ...InfoBuilderOption) *InfoBuilderResult
}

// DefaultAppError implements AppError interface
type DefaultAppError struct {
	err           error
	cause         error
	params        map[string]any
	transParams   map[string]string
	debug         string
	customConfig  *ErrorConfig
	customBuilder InfoBuilderFunc

	disallowGlobalConfigMapping bool
}

// Error implements `error` interface
func (e *DefaultAppError) Error() string {
	return e.err.Error()
}

// Is implementation used by errors.Is().
// This function returns true if either the inner error or the cause satisfies.
func (e *DefaultAppError) Is(err error) bool {
	if errors.Is(e.err, err) {
		return true
	}
	if e.cause != nil && errors.Is(e.cause, err) {
		return true
	}
	return false
}

// Unwrap implementation used by errors.Unwrap() and errors.Is()
func (e *DefaultAppError) Unwrap() error {
	return e.err
}

func (e *DefaultAppError) Params() map[string]any {
	return e.params
}

func (e *DefaultAppError) TransParams() map[string]string {
	return e.transParams
}

func (e *DefaultAppError) Cause() error {
	return e.cause
}

func (e *DefaultAppError) Debug() string {
	return e.debug
}

func (e *DefaultAppError) CustomConfig() *ErrorConfig {
	return e.customConfig
}

func (e *DefaultAppError) CustomBuilder() InfoBuilderFunc {
	return e.customBuilder
}

func (e *DefaultAppError) WithParam(k string, v any) AppError {
	e.params[k] = v
	return e
}

func (e *DefaultAppError) WithTransParam(k string, v string) AppError {
	e.transParams[k] = v
	return e
}

func (e *DefaultAppError) WithCause(cause error) AppError {
	e.cause = cause
	return e
}

func (e *DefaultAppError) WithDebug(format string, args ...any) AppError {
	if !globalConfig.Debug {
		return e
	}
	msg := fmt.Sprintf(format, args...)
	if e.debug == "" {
		e.debug = msg
	} else {
		e.debug = e.debug + globalConfig.MultiErrorSeparator + msg
	}
	return e
}

func (e *DefaultAppError) WithCustomConfig(cfg *ErrorConfig) AppError {
	e.customConfig = cfg
	return e
}

func (e *DefaultAppError) WithCustomBuilder(infoBuilder InfoBuilderFunc) AppError {
	e.customBuilder = infoBuilder
	return e
}

// Config returns the custom config if set, otherwise returns the global mapping one
func (e *DefaultAppError) Config() *ErrorConfig {
	if e.disallowGlobalConfigMapping {
		return e.customConfig
	}
	if e.customConfig != nil {
		return e.customConfig
	}
	return GetErrorConfig(e.err)
}

// BuildConfig builds config for building info from the error
func (e *DefaultAppError) BuildConfig(lang Language, options ...InfoBuilderOption) *InfoBuilderConfig {
	errCfg := e.Config()
	if errCfg == nil {
		errCfg = &ErrorConfig{}
	}
	if errCfg.Status == 0 {
		errCfg.Status = globalConfig.DefaultErrorStatus
	}
	if errCfg.Code == "" {
		errCfg.Code = UnwrapToRoot(e.err).Error()
	}
	buildCfg := &InfoBuilderConfig{
		ErrorConfig:     *errCfg,
		InfoBuilderFunc: e.customBuilder,
		Language:        lang,
		ErrorSeparator:  globalConfig.MultiErrorSeparator,
		TranslationFunc: globalConfig.TranslationFunc,
		FallbackToErrorContentOnMissingTranslation: globalConfig.FallbackToErrorContentOnMissingTranslation,
	}
	for _, opt := range options {
		opt(buildCfg)
	}
	return buildCfg
}

// Build builds error info
func (e *DefaultAppError) Build(lang Language, options ...InfoBuilderOption) *InfoBuilderResult {
	return e.build(e.BuildConfig(lang, options...))
}

// build builds error info using the given building config
func (e *DefaultAppError) build(buildCfg *InfoBuilderConfig) *InfoBuilderResult {
	if buildCfg.InfoBuilderFunc != nil {
		return buildCfg.InfoBuilderFunc(e, buildCfg)
	}

	errInfo := &ErrorInfo{
		AssociatedError: e,
	}
	buildResult := &InfoBuilderResult{
		ErrorInfo: errInfo,
	}

	errCfg := buildCfg.ErrorConfig
	errInfo.Status = errCfg.Status
	errInfo.Code = errCfg.Code
	errInfo.LogLevel = errCfg.LogLevel

	message, title := e.buildMessage(buildCfg, buildResult)
	errInfo.Message = message
	errInfo.Title = title

	// In non-debug mode, output fields `Debug` and `Cause` are set empty
	if globalConfig.Debug {
		errInfo.Debug = e.debug
		if e.cause != nil {
			errInfo.Cause = e.cause.Error()
		}
	}

	return buildResult
}

// buildMessage builds detailed message of the error
func (e *DefaultAppError) buildMessage(buildCfg *InfoBuilderConfig, result *InfoBuilderResult) (msg, title string) {
	title = buildCfg.ErrorConfig.Title
	if title == "" {
		title = http.StatusText(result.ErrorInfo.Status)
	}

	if buildCfg.TranslationFunc == nil {
		return e.Error(), title
	}

	params := e.buildParams(buildCfg, result)

	transKey := buildCfg.ErrorConfig.TransKey
	if transKey == "" {
		transKey = buildCfg.ErrorConfig.Code
	}
	if transKey == "" {
		transKey = UnwrapToRoot(e.err).Error()
	}

	msg, err := buildCfg.TranslationFunc(buildCfg.Language, transKey, params)
	if err != nil {
		result.TransMissingMainKey = true
		result.TransMissingKeys = append(result.TransMissingKeys, transKey)
		if buildCfg.FallbackToErrorContentOnMissingTranslation {
			msg = e.Error()
		}
	}

	if title != "" {
		transTitle, err := buildCfg.TranslationFunc(buildCfg.Language, title, params)
		if err != nil {
			result.TransMissingKeys = append(result.TransMissingKeys, title)
		}
		title = transTitle
	}
	return msg, title
}

// buildParams builds param map from params and translating params
func (e *DefaultAppError) buildParams(buildCfg *InfoBuilderConfig, result *InfoBuilderResult) map[string]any {
	params := e.params
	for k, v := range e.transParams {
		if translated, err := buildCfg.TranslationFunc(buildCfg.Language, v, nil); err != nil {
			result.TransMissingKeys = append(result.TransMissingKeys, v)
			params[k] = v
		} else {
			params[k] = translated
		}
	}
	return params
}

// newDefaultAppError creates *DefaultAppError
func newDefaultAppError(err error) *DefaultAppError {
	return &DefaultAppError{
		err:         Wrap(err),
		params:      map[string]any{},
		transParams: map[string]string{},
	}
}

// NewAppError creates an AppError containing the given error
func NewAppError(err error) AppError {
	if err == nil {
		return nil
	}
	return newDefaultAppError(err)
}
