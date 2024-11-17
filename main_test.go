package goapperrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	t.Run("panic on nil input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "expect panic")
			}
		}()
		Init(nil)
	})
}

func Test_Add(t *testing.T) {
	e1 := errors.New("ErrTokenInvalid")
	e2 := Add(e1, &ErrorConfig{})
	assert.Equal(t, e1, e2)

	t.Run("panic on nil input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "expect panic")
			}
		}()
		_ = Add(nil, &ErrorConfig{})
	})
}

func Test_Create(t *testing.T) {
	e := Create("ErrTokenInvalid", &ErrorConfig{})
	assert.Equal(t, "ErrTokenInvalid", e.Error())

	t.Run("panic on nil input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "expect panic")
			}
		}()
		_ = Create("ErrBadProductSKU", nil)
	})
}

func Test_Remove(t *testing.T) {
	Remove(nil)
	e := Create("ErrBadProductSKU", &ErrorConfig{})
	assert.NotNil(t, GetErrorConfig(e))
	Remove(e)
	assert.Nil(t, GetErrorConfig(e))
}

func Test_Build(t *testing.T) {
	t.Run("builds direct app errors", func(t *testing.T) {
		initConfig(okConfig)

		ae1 := NewAppError(errTest1)
		ae2 := NewMultiError(NewAppError(errTest1), NewAppError(errTest2)).
			WithCustomConfig(&ErrorConfig{
				Status: 1234,
				Code:   "Err1234",
			})

		res1 := Build(ae1, LanguageDe)
		info1 := res1.ErrorInfo
		assert.Equal(t, 0, len(res1.TransMissingKeys))
		assert.Equal(t, 500, info1.Status)
		assert.Equal(t, "ErrTest1", info1.Code)
		assert.Equal(t, "(ErrTest1)-in-de", info1.Message)

		res2 := Build(ae2, LanguageFr)
		info2 := res2.ErrorInfo
		assert.Equal(t, 0, len(res2.TransMissingKeys))
		assert.Equal(t, "(Err1234)-in-fr", info2.Message)
		assert.Equal(t, 1234, info2.Status)
		assert.Equal(t, "Err1234", info2.Code)
		assert.Equal(t, 2, len(info2.InnerErrors))
	})

	t.Run("builds indirect app errors", func(t *testing.T) {
		initConfig(okConfig)

		ae1 := NewAppError(errTest1)
		ae2 := NewMultiError(NewAppError(errTest1), NewAppError(errTest2)).
			WithCustomConfig(&ErrorConfig{
				Status: 1234,
				Code:   "Err1234",
			})

		res1 := Build(Wrap(ae1), LanguageDe)
		info1 := res1.ErrorInfo
		assert.Equal(t, 0, len(res1.TransMissingKeys))
		assert.Equal(t, 500, info1.Status)
		assert.Equal(t, "ErrTest1", info1.Code)
		assert.Equal(t, "(ErrTest1)-in-de", info1.Message)

		res2 := Build(fmt.Errorf("%w", ae2), LanguageFr)
		info2 := res2.ErrorInfo
		assert.Equal(t, 0, len(res2.TransMissingKeys))
		assert.Equal(t, 1234, info2.Status)
		assert.Equal(t, "Err1234", info2.Code)
		assert.Equal(t, "(Err1234)-in-fr", info2.Message)
		assert.Equal(t, 2, len(info2.InnerErrors))
	})

	t.Run("builds non app error", func(t *testing.T) {
		initConfig(okConfig)
		defer initErrorMapping(errTest2, &ErrorConfig{
			Status: 1234,
			Code:   "Err1234",
		})()

		res1 := Build(errTest1, LanguageEn)
		info1 := res1.ErrorInfo
		assert.Equal(t, 0, len(res1.TransMissingKeys))
		assert.Equal(t, 500, info1.Status)
		assert.Equal(t, "ErrTest1", info1.Code)
		assert.Equal(t, "(ErrTest1)-in-en", info1.Message)

		res2 := Build(errTest2, LanguageEn)
		info2 := res2.ErrorInfo
		assert.Equal(t, 0, len(res2.TransMissingKeys))
		assert.Equal(t, 1234, info2.Status)
		assert.Equal(t, "Err1234", info2.Code)
		assert.Equal(t, "(Err1234)-in-en", info2.Message)
	})
}
