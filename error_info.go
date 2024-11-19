package goapperrors

// ErrorInfo stores error information which can be used to return to client
type ErrorInfo struct {
	Status      int          `json:"status,omitempty"`
	Code        string       `json:"code,omitempty"`
	Source      any          `json:"source,omitempty"`
	Title       string       `json:"title,omitempty"`
	Message     string       `json:"message,omitempty"`
	Cause       string       `json:"cause,omitempty"`
	Debug       string       `json:"debug,omitempty"`
	LogLevel    LogLevel     `json:"logLevel,omitempty"`
	InnerErrors []*ErrorInfo `json:"errors,omitempty"`

	AssociatedError error `json:"-"`
}

// InfoBuilderFunc custom info builder function
type InfoBuilderFunc func(AppError, *InfoBuilderConfig) *InfoBuilderResult

// InfoBuilderConfig config used to build error info
type InfoBuilderConfig struct {
	ErrorConfig                                ErrorConfig
	InfoBuilderFunc                            InfoBuilderFunc
	Language                                   Language
	ErrorSeparator                             string
	TranslationFunc                            TranslationFunc
	TranslateTitle                             bool
	FallbackToErrorContentOnMissingTranslation bool
}

// InfoBuilderResult result of building process
type InfoBuilderResult struct {
	// ErrorInfo result information of error which can be used to return to client
	ErrorInfo *ErrorInfo
	// TransMissingKeys missing keys when translate
	TransMissingKeys []string
	// TransMissingMainKey is set `true` if the main message key is missing
	TransMissingMainKey bool
}

// InfoBuilderOption config setter for building error info
type InfoBuilderOption func(*InfoBuilderConfig)

// InfoBuilderOptionCustomBuilder sets custom info builder
func InfoBuilderOptionCustomBuilder(infoBuilderFunc InfoBuilderFunc) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.InfoBuilderFunc = infoBuilderFunc
	}
}

// InfoBuilderOptionCustomConfig sets custom config
func InfoBuilderOptionCustomConfig(errorConfig ErrorConfig) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.ErrorConfig = errorConfig
	}
}

// InfoBuilderOptionTranslationFunc sets custom translation function
func InfoBuilderOptionTranslationFunc(translationFunc TranslationFunc) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.TranslationFunc = translationFunc
	}
}

// InfoBuilderOptionTranslateTitle sets flag indicating title translation
func InfoBuilderOptionTranslateTitle(translateTitle bool) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.TranslateTitle = translateTitle
	}
}

// InfoBuilderOptionSeparator sets custom content separator
func InfoBuilderOptionSeparator(errorSeparator string) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.ErrorSeparator = errorSeparator
	}
}

// InfoBuilderOptionFallbackContent sets flag fallback to error content on missing translation
func InfoBuilderOptionFallbackContent(fallbackToContent bool) InfoBuilderOption {
	return func(cfg *InfoBuilderConfig) {
		cfg.FallbackToErrorContentOnMissingTranslation = fallbackToContent
	}
}
