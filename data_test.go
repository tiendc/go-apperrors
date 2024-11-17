package goapperrors

import (
	"errors"
	"fmt"
)

var (
	errTest1        = errors.New("ErrTest1")
	errTest2        = errors.New("ErrTest2")
	errTest3        = errors.New("ErrTest3")
	errMissingTrans = errors.New("MissingTranslation")
)

var (
	okConfig = &Config{
		Debug:           true,
		DefaultLanguage: LanguageEn,
		TranslationFunc: testTranslateOK,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        ". ",
	}

	noStackTraceConfig = &Config{
		Debug:           true,
		DefaultLanguage: LanguageEn,
		WrapFunc:        func(err error) error { return err },
		TranslationFunc: testTranslateOK,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        ". ",
	}

	failedTransConfig = &Config{
		Debug:           true,
		DefaultLanguage: LanguageEn,
		TranslationFunc: testTranslateFail,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        ". ",
	}

	notransConfig = &Config{
		Debug:           true,
		DefaultLanguage: LanguageEn,
		TranslationFunc: nil,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        ". ",
	}

	nonDebugConfig = &Config{
		Debug:           false,
		DefaultLanguage: LanguageEn,
		TranslationFunc: testTranslateOK,
		FallbackToErrorContentOnMissingTranslation: true,
		MultiErrorSeparator:                        ". ",
	}
)

func testTranslateOK(lang Language, key string, params map[string]any) (string, error) {
	return fmt.Sprintf("(%s)-in-%s", key, lang), nil
}

func testTranslateFail(lang Language, key string, params map[string]any) (string, error) {
	return "", fmt.Errorf("%w: %s in %s", errMissingTrans, key, lang)
}

func initConfig(cfg *Config) {
	Init(cfg)
}

func initErrorMapping(err error, cfg *ErrorConfig) (cleanup func()) {
	_ = Add(err, cfg)
	return func() {
		Remove(err)
	}
}
