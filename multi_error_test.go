package goapperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MultiError_Common(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initConfig(okConfig)

		ae1 := New(errTest1)
		ae2 := New(errTest2)
		me1 := NewMultiError()
		me2 := AsMultiError(NewMultiError(ae1).
			WithParam("k1", "v1").
			WithTransParam("kk1", "vv1").
			WithCause(errTest3).
			WithDebug("debug").
			WithCustomConfig(&ErrorConfig{}).
			WithCustomBuilder(nil))
		me3 := AsMultiError(NewMultiError(ae1, ae2).
			WithCustomConfig(&ErrorConfig{
				Status:   1234,
				Code:     "Err1234",
				TransKey: "CustomKey",
			}))

		assert.Nil(t, me1)

		assert.Equal(t, AppErrors{ae1}, me2.InnerErrors())
		assert.Equal(t, []error{ae1}, UnwrapMulti(me2))
		assert.Equal(t, map[string]any{"k1": "v1"}, me2.Params())
		assert.Equal(t, map[string]string{"kk1": "vv1"}, me2.TransParams())
		assert.ErrorIs(t, errTest3, me2.Cause())
		assert.Equal(t, "debug", me2.Debug())
		assert.Equal(t, ErrorConfig{}, *me2.CustomConfig())

		assert.Equal(t, AppErrors{ae1, ae2}, me3.InnerErrors())
		assert.Equal(t, []error{ae1, ae2}, UnwrapMulti(me3))
		assert.Equal(t, map[string]any{}, me3.Params())
		assert.Equal(t, map[string]string{}, me3.TransParams())
		assert.Nil(t, me3.Cause())
		assert.Equal(t, "", me3.Debug())
	})
}

func Test_MultiError_Build(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initConfig(okConfig)

		ae1 := New(errTest1)
		ae2 := New(errTest2)
		me1 := AsMultiError(NewMultiError(ae1))
		me2 := AsMultiError(NewMultiError(ae1, ae2).
			WithCustomConfig(&ErrorConfig{
				Status:   1234,
				Code:     "Err1234",
				TransKey: "CustomKey",
			}))

		assert.Equal(t, AppErrors{ae1}, me1.InnerErrors())
		assert.Equal(t, []error{ae1}, UnwrapMulti(me1))
		res1 := me1.Build(LanguageEn)
		info1 := res1.ErrorInfo
		assert.Equal(t, 0, len(res1.TransMissingKeys))
		assert.False(t, res1.TransMissingMainKey)
		assert.Equal(t, 500, info1.Status)
		assert.Equal(t, "ErrTest1", info1.Code)

		assert.Equal(t, AppErrors{ae1, ae2}, me2.InnerErrors())
		assert.Equal(t, []error{ae1, ae2}, UnwrapMulti(me2))
		res2 := me2.Build(LanguageFr)
		info2 := res2.ErrorInfo
		assert.Equal(t, 0, len(res2.TransMissingKeys))
		assert.False(t, res2.TransMissingMainKey)
		assert.Equal(t, 1234, info2.Status)
		assert.Equal(t, "Err1234", info2.Code)
		assert.Equal(t, "(CustomKey)-in-fr", info2.Message)
	})

	t.Run("success: failed translation but fallback to content", func(t *testing.T) {
		initConfig(failedTransConfig)
		defer initErrorMapping(errTest1, &ErrorConfig{
			Status: 1234,
			Code:   "Err1234",
		})()

		ae1 := New(errTest1)
		ae2 := New(errTest2)
		me1 := AsMultiError(NewMultiError(ae1, ae2).
			WithCustomConfig(&ErrorConfig{
				Status:   1234,
				Code:     "Err1234",
				TransKey: "CustomKey",
			}))
		// Global mapping has no effect
		defer initErrorMapping(me1, &ErrorConfig{
			Status: 7777,
			Code:   "Err7777",
		})()

		assert.Equal(t, AppErrors{ae1, ae2}, me1.InnerErrors())
		assert.Equal(t, []error{ae1, ae2}, UnwrapMulti(me1))
		res1 := me1.Build(LanguageFr, InfoBuilderOptionFallbackContent(true))
		info1 := res1.ErrorInfo
		assert.Equal(t, 4, len(res1.TransMissingKeys))
		assert.True(t, res1.TransMissingMainKey)
		assert.Equal(t, 1234, info1.Status)
		assert.Equal(t, "Err1234", info1.Code)
		assert.Equal(t, "ErrTest1. ErrTest2", info1.Message)
	})
}
