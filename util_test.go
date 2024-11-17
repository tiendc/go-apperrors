package goapperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnwrapToRoot(t *testing.T) {
	assert.Nil(t, UnwrapToRoot(nil))

	e1 := Wrap(errTest1)
	e2 := Wrap(e1)
	e3 := Wrap(e2)
	assert.Equal(t, errTest1, UnwrapToRoot(e3))
}

type err1 struct{}

func (e err1) Error() string { return "err1" }

type err2 struct {
	err error
}

func (e err2) Error() string { return "err2" }
func (e err2) Unwrap() error { return e.err }

type err3 struct {
	errs []error
}

func (e err3) Error() string   { return "err3" }
func (e err3) Unwrap() []error { return e.errs }

func Test_UnwrapMulti(t *testing.T) {
	assert.Nil(t, UnwrapMulti(err1{}))

	assert.Nil(t, UnwrapMulti(err2{}))
	assert.Equal(t, []error{errTest1}, UnwrapMulti(err2{err: errTest1}))

	assert.Nil(t, UnwrapMulti(err3{}))
	assert.Equal(t, []error{errTest1, errTest2}, UnwrapMulti(err3{errs: []error{errTest1, errTest2}}))
}
