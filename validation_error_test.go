package goapperrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	err3rdPartyVld1 = errors.New("ErrValidation1")
	err3rdPartyVld2 = errors.New("ErrValidation2")
	err3rdPartyVld3 = errors.New("ErrValidation3")
)

type testVldErr struct {
	*defaultAppError
}

func (e *testVldErr) Build(lang Language, options ...InfoBuilderOption) *InfoBuilderResult {
	var message string
	buildCfg := e.BuildConfig(lang, options...)
	if buildCfg.TranslationFunc != nil {
		message, _ = buildCfg.TranslationFunc(lang, e.Error(), nil)
	}
	return &InfoBuilderResult{
		ErrorInfo: &ErrorInfo{
			Status:  http.StatusBadRequest,
			Code:    "ErrValidationItem",
			Message: message,
		},
	}
}

func Test_ValidationError(t *testing.T) {
	t.Run("success: provides custom builder", func(t *testing.T) {
		initConfig(okConfig)

		infoBuilder := func(e AppError, buildCfg *InfoBuilderConfig) *InfoBuilderResult {
			var message string
			if buildCfg.TranslationFunc != nil {
				message, _ = buildCfg.TranslationFunc(buildCfg.Language, e.Error(), nil)
			}
			return &InfoBuilderResult{
				ErrorInfo: &ErrorInfo{
					Status:  http.StatusBadRequest,
					Code:    "ErrValidationItem",
					Message: message,
				},
			}
		}

		assert.Nil(t, NewValidationErrorWithInfoBuilder(nil))
		vldErr := NewValidationErrorWithInfoBuilder(infoBuilder,
			err3rdPartyVld1, err3rdPartyVld2, err3rdPartyVld3)

		result := vldErr.Build(LanguageEn)
		errInfo := result.ErrorInfo

		assert.Equal(t, 0, len(result.TransMissingKeys))
		assert.Equal(t, "(ErrValidation)-in-en", errInfo.Message)
		assert.Equal(t, 3, len(errInfo.InnerErrors))

		inErr0 := errInfo.InnerErrors[0]
		assert.Equal(t, http.StatusBadRequest, inErr0.Status)
		assert.Equal(t, "ErrValidationItem", inErr0.Code)
		assert.Equal(t, "(ErrValidation1)-in-en", inErr0.Message)

		inErr1 := errInfo.InnerErrors[1]
		assert.Equal(t, http.StatusBadRequest, inErr1.Status)
		assert.Equal(t, "ErrValidationItem", inErr1.Code)
		assert.Equal(t, "(ErrValidation2)-in-en", inErr1.Message)
	})

	t.Run("success: implement AppError interface", func(t *testing.T) {
		initConfig(okConfig)

		assert.Nil(t, NewValidationError())
		vldErr := NewValidationError(
			&testVldErr{defaultAppError: NewAppError(err3rdPartyVld1).(*defaultAppError)},
			&testVldErr{defaultAppError: NewAppError(err3rdPartyVld2).(*defaultAppError)},
			&testVldErr{defaultAppError: NewAppError(err3rdPartyVld3).(*defaultAppError)},
		)

		result := vldErr.Build(LanguageEn)
		errInfo := result.ErrorInfo

		assert.Equal(t, 0, len(result.TransMissingKeys))
		assert.Equal(t, "(ErrValidation)-in-en", errInfo.Message)
		assert.Equal(t, 3, len(errInfo.InnerErrors))

		inErr0 := errInfo.InnerErrors[0]
		assert.Equal(t, http.StatusBadRequest, inErr0.Status)
		assert.Equal(t, "ErrValidationItem", inErr0.Code)
		assert.Equal(t, "(ErrValidation1)-in-en", inErr0.Message)

		inErr1 := errInfo.InnerErrors[1]
		assert.Equal(t, http.StatusBadRequest, inErr1.Status)
		assert.Equal(t, "ErrValidationItem", inErr1.Code)
		assert.Equal(t, "(ErrValidation2)-in-en", inErr1.Message)
	})
}
