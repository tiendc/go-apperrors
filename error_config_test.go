package goapperrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetErrorConfig(t *testing.T) {
	t.Run("found: direct mapping", func(t *testing.T) {
		initConfig(okConfig)
		errCfg := &ErrorConfig{
			Status: 1234,
		}
		defer initErrorMapping(errTest1, errCfg)()

		gotCfg := GetErrorConfig(errTest1)
		assert.Equal(t, errCfg, gotCfg)
	})

	t.Run("found: indirect mapping", func(t *testing.T) {
		initConfig(okConfig)
		errCfg := &ErrorConfig{
			Status: 1234,
		}
		defer initErrorMapping(errTest1, errCfg)()

		e1 := Wrap(errTest1)
		e2 := Wrap(e1)
		e3 := fmt.Errorf("blah: %w", e2)
		gotCfg := GetErrorConfig(e3)
		assert.Equal(t, errCfg, gotCfg)
	})

	t.Run("not found", func(t *testing.T) {
		initConfig(okConfig)
		defer initErrorMapping(errTest1, &ErrorConfig{
			Status: 1234,
		})()

		assert.Nil(t, GetErrorConfig(errTest2))
	})

	t.Run("not found: error is in multi-error", func(t *testing.T) {
		initConfig(okConfig)
		defer initErrorMapping(errTest1, &ErrorConfig{
			Status: 1234,
		})()

		assert.Nil(t, GetErrorConfig(errors.Join(errTest1, errTest2)))
	})
}
