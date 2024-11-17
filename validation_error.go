package goapperrors

type ValidationError MultiError

// NewValidationError creates a validation error for the given validation error items
func NewValidationError(errs ...AppError) ValidationError {
	if len(errs) == 0 {
		return nil
	}
	e := NewMultiError(errs...)
	_ = e.WithCustomConfig(&ErrorConfig{
		Status: globalConfig.DefaultValidationErrorStatus,
		Code:   globalConfig.DefaultValidationErrorCode,
	})
	return e
}

func NewValidationErrorWithInfoBuilder(infoBuilder InfoBuilderFunc, errs ...error) ValidationError {
	if len(errs) == 0 {
		return nil
	}
	appErrs := make(AppErrors, 0, len(errs))
	for _, e := range errs {
		appErrs = append(appErrs, NewAppError(e).WithCustomBuilder(infoBuilder))
	}
	return NewValidationError(appErrs...)
}
