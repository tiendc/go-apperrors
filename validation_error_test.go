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

type test3rdPartyVldErr struct {
	errStr string
}

func (e *test3rdPartyVldErr) Error() string { return e.errStr }

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
			vldErr := &test3rdPartyVldErr{}
			errors.As(e, &vldErr)
			message := vldErr.errStr
			if buildCfg.TranslationFunc != nil {
				message, _ = buildCfg.TranslationFunc(buildCfg.Language, e.Error(), nil)
			}
			return &InfoBuilderResult{
				ErrorInfo: &ErrorInfo{
					Status:  http.StatusBadRequest,
					Code:    vldErr.errStr,
					Message: message,
				},
			}
		}

		assert.Nil(t, NewValidationErrorWithInfoBuilder(nil))
		vldErr := NewValidationErrorWithInfoBuilder(infoBuilder,
			&test3rdPartyVldErr{errStr: "3rdPartyVldErr1"},
			&test3rdPartyVldErr{errStr: "3rdPartyVldErr2"},
			&test3rdPartyVldErr{errStr: "3rdPartyVldErr3"})

		result := vldErr.Build(LanguageEn)
		errInfo := result.ErrorInfo

		assert.Equal(t, 0, len(result.TransMissingKeys))
		assert.Equal(t, "(ErrValidation)-in-en", errInfo.Message)
		assert.Equal(t, 3, len(errInfo.InnerErrors))

		inErr0 := errInfo.InnerErrors[0]
		assert.Equal(t, http.StatusBadRequest, inErr0.Status)
		assert.Equal(t, "3rdPartyVldErr1", inErr0.Code)
		assert.Equal(t, "(3rdPartyVldErr1)-in-en", inErr0.Message)

		inErr1 := errInfo.InnerErrors[1]
		assert.Equal(t, http.StatusBadRequest, inErr1.Status)
		assert.Equal(t, "3rdPartyVldErr2", inErr1.Code)
		assert.Equal(t, "(3rdPartyVldErr2)-in-en", inErr1.Message)

		inErr2 := errInfo.InnerErrors[2]
		assert.Equal(t, http.StatusBadRequest, inErr2.Status)
		assert.Equal(t, "3rdPartyVldErr3", inErr2.Code)
		assert.Equal(t, "(3rdPartyVldErr3)-in-en", inErr2.Message)
	})

	t.Run("success: implement AppError interface", func(t *testing.T) {
		initConfig(okConfig)

		assert.Nil(t, NewValidationError())
		vldErr := NewValidationError(
			&testVldErr{defaultAppError: New(err3rdPartyVld1).(*defaultAppError)},
			&testVldErr{defaultAppError: New(err3rdPartyVld2).(*defaultAppError)},
			&testVldErr{defaultAppError: New(err3rdPartyVld3).(*defaultAppError)},
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

		inErr2 := errInfo.InnerErrors[2]
		assert.Equal(t, http.StatusBadRequest, inErr2.Status)
		assert.Equal(t, "ErrValidationItem", inErr2.Code)
		assert.Equal(t, "(ErrValidation3)-in-en", inErr2.Message)
	})
}
