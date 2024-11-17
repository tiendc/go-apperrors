package goapperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_setDefault(t *testing.T) {
	config := &Config{}
	config.setDefault()

	assert.Nil(t, config.WrapFunc)
	assert.Equal(t, defaultMaxStackDepth, config.MaxStackDepth)
	assert.Equal(t, defaultLanguage, config.DefaultLanguage)
	assert.Nil(t, config.TranslationFunc)
	assert.False(t, config.FallbackToErrorContentOnMissingTranslation)
	assert.Equal(t, defaultErrorSeparator, config.MultiErrorSeparator)
	assert.Equal(t, defaultErrorStatus, config.DefaultErrorStatus)
	assert.Equal(t, defaultValidationErrorStatus, config.DefaultValidationErrorStatus)
	assert.Equal(t, defaultValidationErrorCode, config.DefaultValidationErrorCode)
}
