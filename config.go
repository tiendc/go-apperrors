package goapperrors

import (
	"net/http"
)

// Config to provide global config for the library
type Config struct {
	// Debug flag indicates debug mode (default: `false`).
	// If `false`, app error `debug` string can't be set.
	Debug bool

	// WrapFunc to wrap an error with adding stack trace (default: `nil`).
	// This function is nil by default which means the library will use default value
	// which is function Wrap from `github.com/go-errors/errors`.
	WrapFunc func(error) error
	// MaxStackDepth max stack depth (default: `50`).
	// If WrapFunc is set with custom value, this config has no effect.
	MaxStackDepth int

	// DefaultLanguage default language (default: `LanguageEn`)
	DefaultLanguage Language
	// TranslationFunc function to translate message into a specific language (default: `nil`)
	TranslationFunc TranslationFunc
	// FallbackToErrorContentOnMissingTranslation indicates fallback to error content
	// when translation failed (default: `true`).
	// If `false`, when translation fails, the output message will be empty.
	FallbackToErrorContentOnMissingTranslation bool
	// MultiErrorSeparator separator of multiple error strings (default: `\n`)
	MultiErrorSeparator string

	// DefaultErrorStatus default status for error if unset (default: `500`)
	DefaultErrorStatus int
	// DefaultValidationErrorStatus default status for validation error if unset (default: `400`)
	DefaultValidationErrorStatus int
	// DefaultValidationErrorCode default code for validation error if unset (default: `ErrValidation`)
	DefaultValidationErrorCode string

	// DefaultLogLevel default log level for errors if unset (default: `LogLevelNone`)
	DefaultLogLevel LogLevel
}

func (cfg *Config) setDefault() {
	if cfg.MaxStackDepth == 0 {
		cfg.MaxStackDepth = defaultMaxStackDepth
	}
	if cfg.DefaultLanguage == nil {
		cfg.DefaultLanguage = defaultLanguage
	}
	if cfg.MultiErrorSeparator == "" {
		cfg.MultiErrorSeparator = defaultErrorSeparator
	}
	if cfg.DefaultErrorStatus == 0 {
		cfg.DefaultErrorStatus = defaultErrorStatus
	}
	if cfg.DefaultValidationErrorStatus == 0 {
		cfg.DefaultValidationErrorStatus = defaultValidationErrorStatus
	}
	if cfg.DefaultValidationErrorCode == "" {
		cfg.DefaultValidationErrorCode = defaultValidationErrorCode
	}
	if cfg.DefaultLogLevel == LogLevelNone {
		cfg.DefaultLogLevel = defaultLogLevel
	}
}

const (
	defaultMaxStackDepth         = 50
	defaultLanguage              = LanguageEn
	defaultErrorSeparator        = "\n"
	defaultErrorStatus           = http.StatusInternalServerError
	defaultValidationErrorStatus = http.StatusBadRequest
	defaultValidationErrorCode   = "ErrValidation"
	defaultLogLevel              = LogLevelNone
)

var (
	globalConfig = &Config{
		Debug:         false,
		MaxStackDepth: defaultMaxStackDepth,

		DefaultLanguage: defaultLanguage,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        defaultErrorSeparator,

		DefaultErrorStatus:           defaultErrorStatus,
		DefaultValidationErrorStatus: defaultValidationErrorStatus,
		DefaultValidationErrorCode:   defaultValidationErrorCode,

		DefaultLogLevel: defaultLogLevel,
	}

	mapError = make(map[error]*ErrorConfig, 50) //nolint:mnd
)
