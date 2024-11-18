package goapperrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewAppError(t *testing.T) {
	initConfig(okConfig)

	ae1 := New(nil)
	ae2 := New(errTest1)
	ae3 := New(ae2)
	ae4 := fmt.Errorf("%w", ae3)
	ae5 := errors.Join(ae4, ae3, errTest2)

	assert.Nil(t, ae1)
	assert.ErrorIs(t, ae2, errTest1)
	assert.ErrorIs(t, ae3, ae2)
	assert.ErrorIs(t, ae3, errTest1)
	assert.ErrorIs(t, ae4, ae3)
	assert.ErrorIs(t, ae4, ae2)
	assert.ErrorIs(t, ae4, errTest1)
	assert.ErrorIs(t, ae5, ae4)
	assert.ErrorIs(t, ae5, ae3)
	assert.ErrorIs(t, ae5, ae2)
	assert.ErrorIs(t, ae5, errTest1)
	assert.ErrorIs(t, ae5, errTest2)
}

func Test_AppError_Common(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initConfig(okConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)

		assert.ErrorIs(t, ae, errTest1)
		assert.ErrorIs(t, ae, errTest2) // via cause
		assert.Equal(t, errTest1.Error(), ae.Error())
		assert.Equal(t, "debug: 123", ae.Debug())
		assert.Equal(t, map[string]any{"k1": "v1", "k2": "v2"}, ae.Params())
		assert.Equal(t, map[string]string{"kk1": "vv1"}, ae.TransParams())
		assert.ErrorIs(t, ae.Cause(), errTest2)
	})

	t.Run("success: extending debug message", func(t *testing.T) {
		initConfig(okConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithDebug("blah blah")

		assert.ErrorIs(t, ae, errTest1)
		assert.Equal(t, errTest1.Error(), ae.Error())
		assert.Equal(t, "debug: 123. blah blah", ae.Debug())
	})

	t.Run("success: non-debug mode", func(t *testing.T) {
		initConfig(nonDebugConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithCause(errTest3)

		assert.ErrorIs(t, ae, errTest1)
		assert.Equal(t, errTest1.Error(), ae.Error())
		assert.Equal(t, "", ae.Debug())
		assert.ErrorIs(t, ae.Cause(), errTest3)
	})
}

func Test_AppError_Build(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initConfig(okConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)

		buildRes := ae.Build(LanguageEn)
		errInfo := buildRes.ErrorInfo
		assert.Equal(t, 0, len(buildRes.TransMissingKeys))
		assert.Equal(t, 500, errInfo.Status)
		assert.Equal(t, "ErrTest1", errInfo.Code)
		assert.Equal(t, "(ErrTest1)-in-en", errInfo.Message)
		assert.Equal(t, LogLevelNone, errInfo.LogLevel)
		assert.Equal(t, "debug: 123", errInfo.Debug)
		assert.Equal(t, errTest2.Error(), errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})

	t.Run("success: with global config mapping", func(t *testing.T) {
		initConfig(okConfig)
		defer initErrorMapping(errTest1, &ErrorConfig{
			Status:   555,
			Code:     "ErrCustom",
			TransKey: "CustomTrans",
		})()
		defer initErrorMapping(errTest2, &ErrorConfig{
			Status: 777,
			Code:   "ErrCustom",
		})()

		ae1 := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)
		ae2 := New(errTest2).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1")

		buildRes1 := ae1.Build(LanguageEn)
		errInfo1 := buildRes1.ErrorInfo
		assert.Nil(t, ae1.CustomConfig())
		assert.Equal(t, 0, len(buildRes1.TransMissingKeys))
		assert.Equal(t, 555, errInfo1.Status)
		assert.Equal(t, "ErrCustom", errInfo1.Code)
		assert.Equal(t, "(CustomTrans)-in-en", errInfo1.Message)
		assert.Equal(t, LogLevelNone, errInfo1.LogLevel)
		assert.Equal(t, "debug: 123", errInfo1.Debug)
		assert.Equal(t, errTest2.Error(), errInfo1.Cause)
		assert.Equal(t, 0, len(errInfo1.InnerErrors))

		buildRes2 := ae2.Build(LanguageZh)
		errInfo2 := buildRes2.ErrorInfo
		assert.Nil(t, ae2.CustomConfig())
		assert.Equal(t, 0, len(buildRes2.TransMissingKeys))
		assert.Equal(t, 777, errInfo2.Status)
		assert.Equal(t, "ErrCustom", errInfo2.Code)
		assert.Equal(t, "(ErrCustom)-in-zh", errInfo2.Message)
		assert.Equal(t, LogLevelNone, errInfo2.LogLevel)
		assert.Equal(t, "", errInfo2.Debug)
		assert.Equal(t, "ErrTest2", errInfo2.Cause)
		assert.Equal(t, 0, len(errInfo2.InnerErrors))
	})

	t.Run("success: with overriding global config mapping", func(t *testing.T) {
		initConfig(okConfig)
		// Set global config mapping
		defer initErrorMapping(errTest1, &ErrorConfig{
			Status:   123,
			Code:     "Code123",
			TransKey: "Trans123",
		})()

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2).
			WithCustomConfig(&ErrorConfig{ // override the global config mapping
				Status:   555,
				Code:     "ErrCustom",
				TransKey: "CustomTrans",
			})

		buildRes := ae.Build(LanguageEn)
		errInfo := buildRes.ErrorInfo
		assert.NotNil(t, ae.CustomConfig())
		assert.Equal(t, 0, len(buildRes.TransMissingKeys))
		assert.Equal(t, 555, errInfo.Status)
		assert.Equal(t, "ErrCustom", errInfo.Code)
		assert.Equal(t, "(CustomTrans)-in-en", errInfo.Message)
		assert.Equal(t, LogLevelNone, errInfo.LogLevel)
		assert.Equal(t, "debug: 123", errInfo.Debug)
		assert.Equal(t, errTest2.Error(), errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})

	t.Run("success: with custom info builder", func(t *testing.T) {
		initConfig(okConfig)

		infoBuilder := func(e AppError, cfg *InfoBuilderConfig) *InfoBuilderResult {
			return &InfoBuilderResult{
				ErrorInfo: &ErrorInfo{
					Status:   1234,
					Code:     "Err1234",
					Message:  "Err1234-Message",
					Cause:    "Err1234-Cause",
					LogLevel: LogLevelInfo,
				},
			}
		}

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)

		// Info builder set in error object
		_ = ae.WithCustomBuilder(infoBuilder)
		buildRes := ae.Build(LanguageEn)
		errInfo := buildRes.ErrorInfo
		assert.Nil(t, ae.CustomConfig())
		assert.NotNil(t, ae.CustomBuilder())
		assert.Equal(t, 0, len(buildRes.TransMissingKeys))
		assert.Equal(t, 1234, errInfo.Status)
		assert.Equal(t, "Err1234", errInfo.Code)
		assert.Equal(t, "Err1234-Message", errInfo.Message)
		assert.Equal(t, LogLevelInfo, errInfo.LogLevel)
		assert.Equal(t, "", errInfo.Debug)
		assert.Equal(t, "Err1234-Cause", errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))

		// Info builder set via a function argument
		_ = ae.WithCustomBuilder(nil)
		buildRes = ae.Build(LanguageEn, func(buildCfg *InfoBuilderConfig) {
			buildCfg.InfoBuilderFunc = infoBuilder
		})
		errInfo = buildRes.ErrorInfo
		assert.Nil(t, ae.CustomConfig())
		assert.Nil(t, ae.CustomBuilder())
		assert.Equal(t, 0, len(buildRes.TransMissingKeys))
		assert.Equal(t, 1234, errInfo.Status)
		assert.Equal(t, "Err1234", errInfo.Code)
		assert.Equal(t, "Err1234-Message", errInfo.Message)
		assert.Equal(t, LogLevelInfo, errInfo.LogLevel)
		assert.Equal(t, "", errInfo.Debug)
		assert.Equal(t, "Err1234-Cause", errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})

	t.Run("success: fails to translate but no fallback to error string", func(t *testing.T) {
		initConfig(failedTransConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)

		buildRes := ae.Build(LanguageEn, InfoBuilderOptionFallbackContent(false))
		errInfo := buildRes.ErrorInfo
		assert.Equal(t, []string{"vv1", "ErrTest1", "Internal Server Error"}, buildRes.TransMissingKeys)
		assert.Equal(t, 500, errInfo.Status)
		assert.Equal(t, "ErrTest1", errInfo.Code)
		assert.Equal(t, "", errInfo.Message) // message is empty (more secured)
		assert.Equal(t, LogLevelNone, errInfo.LogLevel)
		assert.Equal(t, "debug: 123", errInfo.Debug)
		assert.Equal(t, errTest2.Error(), errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})

	t.Run("success: fails to translate but fallback to error string", func(t *testing.T) {
		initConfig(failedTransConfig)

		ae := New(errTest1).
			WithDebug("debug: %v", 123).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1").
			WithCause(errTest2)

		buildRes := ae.Build(LanguageEn, InfoBuilderOptionFallbackContent(true))
		errInfo := buildRes.ErrorInfo
		assert.Equal(t, []string{"vv1", "ErrTest1", "Internal Server Error"}, buildRes.TransMissingKeys)
		assert.Equal(t, 500, errInfo.Status)
		assert.Equal(t, "ErrTest1", errInfo.Code)
		assert.Equal(t, "ErrTest1", errInfo.Message)
		assert.Equal(t, LogLevelNone, errInfo.LogLevel)
		assert.Equal(t, "debug: 123", errInfo.Debug)
		assert.Equal(t, errTest2.Error(), errInfo.Cause)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})

	t.Run("success: translation function unset", func(t *testing.T) {
		initConfig(notransConfig)

		ae := New(errTest1).
			WithParam("k1", "v1").
			WithParam("k2", "v2").
			WithTransParam("kk1", "vv1")

		buildRes := ae.Build(LanguageEn)
		errInfo := buildRes.ErrorInfo
		assert.Nil(t, buildRes.TransMissingKeys)
		assert.Equal(t, 500, errInfo.Status)
		assert.Equal(t, "ErrTest1", errInfo.Code)
		assert.Equal(t, "ErrTest1", errInfo.Message)
		assert.Equal(t, 0, len(errInfo.InnerErrors))
	})
}
