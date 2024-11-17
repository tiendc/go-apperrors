package goapperrors

import (
	"errors"
	"strings"
)

// AppErrors is a defined type of slice of AppError
type AppErrors []AppError

// Error implements `error` interface
func (aes AppErrors) Error() string {
	var sb strings.Builder
	for i, err := range aes {
		if i > 0 {
			sb.WriteString(globalConfig.MultiErrorSeparator)
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Is implementation used by errors.Is()
func (aes AppErrors) Is(err error) bool {
	for _, e := range aes {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

// Unwrap implementation used by errors.Is().
// errors.Unwrap() only works with `Unwrap() error`.
func (aes AppErrors) Unwrap() []error {
	ret := make([]error, len(aes))
	for i, err := range aes {
		ret[i] = err
	}
	return ret
}
