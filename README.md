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

TBD

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
