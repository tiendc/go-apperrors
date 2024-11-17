package goapperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InfoBuilderOption(t *testing.T) {
	buildConfig := &InfoBuilderConfig{}

	InfoBuilderOptionCustomBuilder(func(AppError, *InfoBuilderConfig) *InfoBuilderResult {
		return &InfoBuilderResult{}
	})(buildConfig)
	assert.NotNil(t, buildConfig.InfoBuilderFunc)

	errConfig := ErrorConfig{}
	InfoBuilderOptionCustomConfig(errConfig)(buildConfig)
	assert.Equal(t, errConfig, buildConfig.ErrorConfig)

	InfoBuilderOptionTranslationFunc(func(Language, string, map[string]any) (string, error) {
		return "", nil
	})(buildConfig)
	assert.NotNil(t, buildConfig.TranslationFunc)

	InfoBuilderOptionSeparator("abc123")(buildConfig)
	assert.Equal(t, "abc123", buildConfig.ErrorSeparator)

	InfoBuilderOptionFallbackContent(true)(buildConfig)
	assert.True(t, buildConfig.FallbackToErrorContentOnMissingTranslation)
}
