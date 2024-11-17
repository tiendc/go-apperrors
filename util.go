package goapperrors

import "errors"

// UnwrapToRoot keeps unwrapping until the root error
func UnwrapToRoot(err error) error {
	lastErr := err
	for {
		e := errors.Unwrap(lastErr)
		if e == nil {
			return lastErr
		}
		lastErr = e
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
