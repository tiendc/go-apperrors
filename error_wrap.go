package goapperrors

import (
	"errors"
	"fmt"
	"runtime"

	goerrors "github.com/go-errors/errors"
)

// Wrap wraps an error with adding stack trace if configured
func Wrap(err error) error {
	if globalConfig.WrapFunc != nil {
		return globalConfig.WrapFunc(err)
	}
	return goerrors.Wrap(err, 1)
}

// Wrapf wraps an error by calling fmt.Errorf and adds stack trace if configured
func Wrapf(format string, args ...any) error {
	return Wrap(fmt.Errorf(format, args...)) //nolint:err113
}

// GetStackTrace gets stack trace stored in the error if there is
func GetStackTrace(err error) []runtime.Frame {
	if err == nil {
		return nil
	}
	// If it is *go-errors.Error, get the stack trace from it
	var gErr *goerrors.Error
	if !errors.As(err, &gErr) {
		return nil
	}
	callers := gErr.Callers()
	frames := runtime.CallersFrames(callers)
	frameList := make([]runtime.Frame, 0, len(callers))
	for {
		next, hasMore := frames.Next()
		frameList = append(frameList, next)
		if !hasMore {
			break
		}
	}
	return frameList
}

// UnwrapToRoot keeps unwrapping until the root error
func UnwrapToRoot(err error) error {
	lastErr := err
	for {
		errs := UnwrapMulti(lastErr)
		if len(errs) == 0 {
			return lastErr
		}
		lastErr = errs[0]
	}
}

// UnwrapMulti unwraps en error to a slice of errors.
// If the error implements `Unwrap() []error`, the result of func call is returned.
// If the error implements `Unwrap() error`, the func is called and the result slice
// has only one element.
// If no `Unwrap` func is implemented in the error, `nil` is returned.
func UnwrapMulti(err error) []error {
	if u, ok := err.(interface{ Unwrap() []error }); ok {
		return u.Unwrap()
	}
	if u, ok := err.(interface{ Unwrap() error }); ok {
		e := u.Unwrap()
		if e == nil {
			return nil
		}
		return []error{e}
	}
	return nil
}
