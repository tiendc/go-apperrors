package goapperrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AppErrorSlice(t *testing.T) {
	initConfig(okConfig)

	ae1 := NewAppError(errTest1)
	ae2 := NewAppError(ae1)
	aeSlice := AppErrors{ae1, ae2}

	assert.ErrorIs(t, aeSlice, errTest1)
	assert.ErrorIs(t, aeSlice, ae1)
	assert.ErrorIs(t, aeSlice, ae2)
	assert.False(t, errors.Is(aeSlice, errTest2))
	assert.Equal(t, []error{ae1, ae2}, aeSlice.Unwrap())
	assert.Equal(t, "ErrTest1. ErrTest1", aeSlice.Error())
}
