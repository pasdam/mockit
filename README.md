# mockit

[![Go Report Card](https://goreportcard.com/badge/github.com/pasdam/mockit)](https://goreportcard.com/report/github.com/pasdam/mockit)
[![CI Status](https://github.com/pasdam/mockit/workflows/Continuous%20integration/badge.svg)](https://github.com/pasdam/mockit/actions)
[![GoDoc](https://godoc.org/github.com/pasdam/mockit?status.svg)](https://godoc.org/github.com/pasdam/mockit)

- [mockit](#mockit)
  - [Notes](#notes)
  - [Usage](#usage)
    - [Argument matcher](#argument-matcher)
      - [Capture argument](#capture-argument)
    - [Pausing and restoring a mock](#pausing-and-restoring-a-mock)
    - [Verify a call](#verify-a-call)
    - [Update the library](#update-the-library)
  - [Development](#development)
    - [TODOs](#todos)
    - [Contributing](#contributing)
    - [Internals](#internals)
  - [Credits](#credits)

Mockit is a library to use during testing for Go application, and aim to make mocking of functions/methods easy.

## Notes

This is still a working in progress so **API might change** before reaching a stable state.

Also please note that the **mocking might not work** in some cases if function inlining is enabled, so it might be necessary to disable it during testing:

```sh
go test -gcflags=-l
```

## Usage

To mock a function:

```go
m := MockFunc(t, filepath.Base)
m.With("some-argument").Return("result")
```

This will make sure that when `filepath.Base` is called with the argument `some-argument`, it will return `result`.

To mock an instance method (at the moment only exported methods are supported):

```go
err := errors.New("some-error")
m := MockMethod(t, err, err.Error)
m.With().Return("some-other-value")
```

When a method is mocked and a matching call is not found (i.e. arguments are different) it will return the zero values.

It is possible to make the mock call the real method:

```go
m.With("some-argument").CallRealMethod()
```

or return zero values:

```go
m.With("some-argument").ReturnDefaults()
```

Mocks are matched in order, which means that:

```go
m.With("some-argument").CallRealMethod()
m.With("some-argument").ReturnDefaults()
```

will make `filepath.Base("some-argument")` call the real method.

Mocks are *automatically removed* when the test is completed.

### Argument matcher

It is also possible to use argument matchers, to implement generic behaviors. At the moment there is only one matcher implemented, and it matches any argument:

```go
m := MockFunc(t, filepath.Base)
m.With(argument.any).Return("result")
```

This will make `filepath.Base` return `result` for any input.

#### Capture argument

To capture the argument of a call:

```go
m := MockFunc(t, filepath.Base)
c := argument.Captor{}
m.With(c.Capture).Return("result")
filepath.Base("some-argument")
```

At this point `c.Value` will be `some argument`.

### Pausing and restoring a mock

It is possible to temporary disable a mock:

```go
m := MockFunc(t, filepath.Base)
m.With("matching-argument").Return("some-out")

// ... Do something with the mock

m.Disable()
```

At this point the mock is disabled and the real implementation is used.

To enable the mock again, just use:

```go
m.Enable()
```

### Verify a call

To verify a specified call happened:

```go
m := MockFunc(t, filepath.Base)

// ... Mock calls
// ... And use mock

m.Verify("matching-argument")
```

The `Verify` method will fail the test if the call didn't happen.

### Update the library

To update the library to the latest version simply run:

```sh
go get -u github.com/pasdam/mockit
```

## Development

### TODOs

This are (not in a particular order) the missing features that are going to be implemented in a not well defined future (patches are welcome):

- [ ] Mock unexported methods
- [ ] Mock interfaces
- [ ] Mock a method for all instances
- [ ] Automatically verify at the end of the test, without having to call `verify` method
- [ ] [Verify in order calls](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#in_order_verification)
- [ ] [Verifying exact number of invocations / at least x / never](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#at_least_verification)
- [Arguments matcher](https://site.mockito.org/javadoc/current/index.html?org/mockito/ArgumentMatcher.html)
  - [ ] IsA: to match for specific types
  - [ ] NotNil: to match any not nil value
- [ ] [Stubbing consecutive calls](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#stubbing_consecutive_calls)
- [ ] [Stubbing with callbacks](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#answer_stubs)
- [ ] Mock [variadic function](https://gobyexample.com/variadic-functions)
- [ ] Override existing mock, i.e. change return values of a stub
  - [ ] [Making sure interaction(s) never happened on mock](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#never_verification)
  - [ ] [Finding redundant invocations](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#finding_redundant_invocations)
- [ ] Improve error messages

### Contributing

Rules to contribute to the repo:

1. Define ine identifier per file, which means that each go file contains either a struct (with related methods), an interface, or a function. Constants should be declared in the file `<package>.go`, i.e. in `mockit.go` for the `mockit` package.
2. Write unit test for each method/function, in order to keep the coverage to 100%.

### Internals

This library uses [monkey](https://github.com/bouk/monkey), a package for [monkey patching](https://en.wikipedia.org/wiki/Monkey_patch) in Go. And of course it inherits the [same limitations](https://github.com/bouk/monkey#notes), in particular:

> Monkey sometimes fails to patch a function if inlining is enabled. Try running your tests with inlining disabled, for example: go test -gcflags=-l.

## Credits

All of this was possible only because of [bouk](https://github.com/bouk), as this library is basically a wrapper around [monkey](https://github.com/bouk/monkey), so kudos to him.
