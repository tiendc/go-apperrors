[![Go Version][gover-img]][gover] [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![GoReport][rpt-img]][rpt]

# Manipulating application errors with ease

## Functionalities

- Provides AppError type to be used for wrapping any kind of errors in an application
- Supports a centralized definition of errors
- AppError can carry extra information such as `cause`, `debug log`, and stack trace
- AppError supports translating its error message into a specific language
- AppError can be transformed to a JSON structure which is friendly to client side
- Provides MultiError type to handle multiple AppError(s) with a common use case of validation errors

## Installation

```shell
go get github.com/tiendc/go-apperrors
```

## Usage

### General usage

**Initializes `go-apperrors` at program startup**

```go
import gae "github.com/tiendc/go-apperrors"

func main() {
    ...
    gae.Init(&gae.Config{
        Debug: ENV == "development",
        DefaultLogLevel: gae.LogLevelInfo,
        TranslationFunc: func (lang gae.Language, key string, params map[string]any) {
            // Provides your implementation to translate message
        },
    })
    ...
}
```

**Defines your errors**

```go
// It is recommended to add a new directory for placing app errors.
// In this example, I use `apperrors/errors.go`.

import gae "github.com/tiendc/go-apperrors"

// Some base errors
var (
    ErrInternalServer = gae.Create("ErrInternalServer", &gae.ErrorConfig{
        Status: http.StatusInternalServerError,
        LogLevel: gae.LogLevelError, // this indicates an unexpected error
    })
    ErrUnauthorized = gae.Create("ErrUnauthorized", &gae.ErrorConfig{Status: http.StatusUnauthorized})
    ErrNotFound = gae.Create("ErrNotFound", &gae.ErrorConfig{Status: http.StatusNotFound})
    ...
)

// Errors from external libs
var (
    ErrRedisKeyNotFound = gae.Add(redis.Nil, &gae.ErrorConfig{Status: http.StatusNotFound})
)

// Some more detailed errors
var (
    ErrUserNotInProject = gae.Create("ErrUserNotInProject", &gae.ErrorConfig{Status: http.StatusForbidden})
    ...
)
```

**Handles errors in your main processing code**

```go
// There are some use cases as below.

// 1. You get an unexpected error
// Just wrap it and return. This will result in error 500 returned to client.
resp, err := updateProject(project)
if err != nil {
    return gae.New(err).WithDebug("extra info: project %s", project.ID)
    // OR `return gae.Wrap(err)` if you don't need to add extra info
}

// 2. You get an error which may be expected
resp, err := deleteProject(project)
if err != nil {
    if errors.Is(err, DBNotFound) { // this error can be expected
        return gae.New(gae.ErrNotFound).WithCause(err)
    }
    return gae.Wrap(err) // unexpected error
}

// 3. You want to return an error when a condition isn't satisfied
if `user.ID` is not in `project.userIDs` {
    // This will return error Forbidden to client as we configured previously
    return gae.New(gae.ErrUserNotInProject)
}
```

**Handles validation errors**

```go
// Validation is normally performed when you parse requests from client.
// You may use an external lib for the validation. That's why you need to make
// `adapter` code to transform the validation errors to `AppError`s.

// Add a new file to the above directory, says `apperrors/validation_errors.go`.
// This building function will be used later to create validation error.
func ValidationErrorInfoBuilder(err AppError, buildCfg *gae.InfoBuilderConfig) *gae.InfoBuilderResult {
    // Extracts the inner error and casts it to the validation error type you use
    vldErr := &thirdPartyLib.Error{}
    if !errors.As(err, &vldErr) {
        // panic, this should not happen
    }
    return &gae.InfoBuilderResult{
        // Transform the error from the 3rd party lib to ErrorInfo struct
        ErrorInfo: &gae.ErrorInfo{
            Message: buildCfg.TranslationFunc(buildCfg.Language, vldErr.getMessage(), err.Params()),
            Source: vldErr.getSource(),
            ...
        }
    }
}

// When parse your request
func (req UpdateProjectReq) Validate() gae.ValidationError {
    vldErrors := validateReq(req)
    return gae.NewValidationErrorWithInfoBuilder(apperrors.ValidationErrorInfoBuilder, vldErrors...)
}
```

**Handles errors before returning them to client**

```go
// In the base handler, implements function `RenderError()`
func RenderError(err error) {
    // Gets language from request, you can use `gae.ParseAcceptLanguage()`
    lang := parseLanguageFromRequest()

    // Call goapperrors.Build
    buildResult := gae.Build(err, lang)

    // Logs error to Sentry or a similar service
    if buildResult.ErrorInfo.LogLevel != gae.LogLevelNone {
        logErrorToSentry(err, buildResult.ErrorInfo.LogLevel)
    }

    // Sends the error as JSON to client
    response.SendJSON(buildResult.ErrorInfo)
}

// In your specific handler, for instance, project handler
func (h ProjectHandler) UpdateProject() {
    updateProjectReq, err := parseAndValidateRequest(httpReq)
    if err != nil {
        baseHandler.RenderError(err)
        return
    }

    resp, err := useCase.UpdateProject(updateProjectReq)
    if err != nil {
        baseHandler.RenderError(err)
        return
    }

    // Send result to client
    response.SendJSON(resp)
}
```

### Global configuration

[See the full code](config.go)

#### Debug (default: `false`)

You should synchronize this option with the ENV value. If it is `true`, error building will
return the fields `Cause` and `Debug` which is convenient for development. If it is `false`,
they will not be returned which is more secured in `production` env as their values can
have sensitive information.

#### WrapFunc (default: `nil`)

If this value is default (which is `nil`), `go-apperrors` will use the lib
[github.com/go-errors/errors](https://github.com/go-errors/errors) to wrap and attach stack trace to errors.
If you don't want to attach stack trace, provide a self-implementation on initialization.

```go
Init(&Config{
    WrapFunc: func (err error) error { return err }
})
```

#### TranslationFunc (default: `nil`)

Sets this option by providing a function to help the lib translate error messages. Otherwise,
the translation will be disabled.

```go
Init(&Config{
    TranslationFunc: func (lang Language, key string, params map[string]any) {
        // Provides your implementation to translate message
    },
})
```

#### FallbackToErrorContentOnMissingTranslation (default: `true`)

When translation fails, if this flag is `true`, the error content will be used to assign to
the output `Message` field. Otherwise, the output message will be empty.

NOTE: turn off this flag if you don't want to reveal sensitive information on building.

## Contributing

- You are welcome to make pull requests for new functions and bug fixes.

## License

- [MIT License](LICENSE)

[doc-img]: https://pkg.go.dev/badge/github.com/tiendc/go-apperrors
[doc]: https://pkg.go.dev/github.com/tiendc/go-apperrors
[gover-img]: https://img.shields.io/badge/Go-%3E%3D%201.20-blue
[gover]: https://img.shields.io/badge/Go-%3E%3D%201.20-blue
[ci-img]: https://github.com/tiendc/go-apperrors/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/tiendc/go-apperrors/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/tiendc/go-apperrors/branch/main/graph/badge.svg
[cov]: https://codecov.io/gh/tiendc/go-apperrors
[rpt-img]: https://goreportcard.com/badge/github.com/tiendc/go-apperrors
[rpt]: https://goreportcard.com/report/github.com/tiendc/go-apperrors
