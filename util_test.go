package goapperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseRequestLangAsStr(t *testing.T) {
	t.Run("any lang", func(t *testing.T) {
		langs, err := ParseAcceptLanguageAsStr("*")
		assert.NoError(t, err)
		assert.Equal(t, []string{"mul"}, langs)
	})

	t.Run("simple language string", func(t *testing.T) {
		langs, err := ParseAcceptLanguageAsStr("de")
		assert.NoError(t, err)
		assert.Equal(t, []string{"de"}, langs)
	})

	t.Run("full locale string", func(t *testing.T) {
		langs, err := ParseAcceptLanguageAsStr("ja-JP")
		assert.NoError(t, err)
		assert.Equal(t, []string{"ja-JP"}, langs)
	})

	t.Run("complex accept-language", func(t *testing.T) {
		langs, err := ParseAcceptLanguageAsStr("fr-CH, fr;q=0.9, en;q=0.8, en, de;q=0.7, de, *;q=0.5")
		assert.NoError(t, err)
		assert.Equal(t, []string{"fr-CH", "en", "de", "fr", "en", "de", "mul"}, langs)
	})

	t.Run("invalid accept-language", func(t *testing.T) {
		_, err := ParseAcceptLanguageAsStr("abc123")
		assert.NotNil(t, err)
	})
}
