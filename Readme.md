Package golang error
=======================
[![Made in Vietnam](https://raw.githubusercontent.com/webuild-community/badge/master/svg/made.svg)](https://webuild.community)
![Project status](https://img.shields.io/badge/version-0.0.1-green.svg)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package golang error handling in [dwarves.foundation](https://dwarves.foundation)

## Installation

Use go get.

```
go get github.com/dwarvesf/gerr
```

Then import the `gerr` package into your own code.

```
import "github.com/dwarvesf/gerr"
```

## Usages
Getting started
- Init an error
- Make HTTP response from error
- Prepare for validation error from [validator](https://github.com/go-playground/validator)

### Init an error

Init with struct

```go
err := gerr.Error{
  TracingID : "abc123",
  Code      : 400,
  Message   : "bad request",
  Errors    : []*gerr.Error{
    {
      Target  : "username",
      Message : "username is required field",
    },
    {
      Target  : "password",
      Message : "password is required field",
    }
  }
}
```

Init with utility function. `gerr` receive almost data types to make error instance.

```go
err := gerr.E(
  "bad request",
  400,
  gerr.TracingID("abc123"),
  gerr.Error{
    Target  : "username",
    Message : "username is required field",
  },
  gerr.Error{
    Target  : "password",
    Message : "password is required field",
  }
)
```

### Make HTTP response error

```go
err := gerr.Error{}
// First way
httpErr := gerr.NewResponseError(err)
// OR
httpErr2 := err.ToResponseError()
```

We receive the error with the nested object in detail

### Prepare for validation error from [validator](https://github.com/go-playground/validator)

In `go`, we usually use `validator` package to validate data. We can:
- Validate our variable and receive error structure with our customize
- Localization supported: 12 languages code (Aug 30th, 2020)

The structure we usually receive from `validator` [FieldError](https://github.com/go-playground/validator/blob/master/errors.go#L83) can be:
```go
type FieldError interface {
  // ...
  // eg. JSON name "User.fname"
	Namespace() string 

  // eq. "User.FirstName" see Namespace for comparison
  StructNamespace() string

  // eq. JSON name "fname"
  Field() string

  // eq.  "FirstName"
  StructField() string
}
```
In `gerr`, we use `Namespace()` to get key with nested object's key is combine into a string.
We make a structure `CombinedError` for this case. Make a utility function to convert to `gerr.Error`
The idea for other cases is to try to convert other error types to `gerr.Error`. That's all.

```go
newErr := gerr.CombinedE(
  "bad request",
  400,
  gerr.CombinedE(gerr.Target("user.name"), "name is required field"),
  gerr.CombinedE(gerr.Target("user.password"), "password is required field")
)

var err *gerr.Error
err = newErr.ToError()
```

## External packages
In `gerr` we use some packages
- [logus](https://github.com/sirupsen/logrus): log library for golang

## Supported features

- [x] Generate the error json format compatible front-end
- [x] Common errors
- [ ] Validation request body
  - [ ] gin
- [x] Log util support customize format
  - [x] simple message
  - [x] json format: `{"instance": "production", "message": "i’m a syslog!"}`
  - [x] loki format: `{instance="production"} 00:00:00 i’m a syslog!`

## How to Contribute

Make a pull request...

## License

MIT &copy; [Dwarves Foundation](https://github.com/dwarvesf)