# mockit

- [mockit](#mockit)
  - [Usage](#usage)
    - [Argument matcher](#argument-matcher)
    - [Pausing and restoring a mock](#pausing-and-restoring-a-mock)
  - [Development](#development)
    - [TODOs](#todos)
    - [Contributing](#contributing)
    - [Internals](#internals)
  - [Credits](#credits)

Mockit is a library to use during testing for Go application, and aim to make mocking of functions/methods easy.

## Usage

To mock a function:

```go
m := NewFuncMock(t, filepath.Base).(*funcMock)
m.With([]interface{}{"some-argument"}).Return([]interface{}{"result"})
```

This will make sure that when `filepath.Base` is called with the argument `some-argument`, it will return `result`.

When a method is mocked and a matching call is not found (i.e. arguments are different) it will return the zero values.

It is possible to make the mock call the real method:

```go
m.With([]interface{}{"some-argument"}).CallRealMethod()
```

or return zero values:

```go
m.With([]interface{}{"some-argument"}).ReturnDefaults()
```

Mocks are matched in order, which means that:

```go
m.With([]interface{}{"some-argument"}).CallRealMethod()
m.With([]interface{}{"some-argument"}).ReturnDefaults()
```

will make `filepath.Base("some-argument")` call the real method.

### Argument matcher

It is also possible to use argument matchers, to have generic mocks. At the moment there only one matcher implemented, and it matches any argument:

```go
m := NewFuncMock(t, filepath.Base).(*funcMock)
m.With([]interface{}{argument.any}).Return([]interface{}{"result"})
```

This will make `filepath.Base` return `result` for any input.

### Pausing and restoring a mock

It is possible to temporary disable a mock:

```go
m := NewFuncMock(t, filepath.Base).(*funcMock)
m.With([]interface{}{"matching-argument"}).Return([]interface{}{"some-out"})

// ... Do something with the mock

m.Disable()
```

At this point the mock is disabled and the real implementation is used.

To enable the mock again, just use:

```go
m.Enable()
```

## Development

### TODOs

This are (not in a particular order) the missing features that are going to be implemented in a not well defined future (patches are welcome):

- [ ] [Verify in order calls](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#in_order_verification)
- [ ] [Verifying exact number of invocations / at least x / never](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#at_least_verification)
- [Arguments matcher](https://site.mockito.org/javadoc/current/index.html?org/mockito/ArgumentMatcher.html)
  - [ ] IsA: to match for specific types
  - [ ] NotNil: to match any not nil value
- [ ] [Stubbing consecutive calls](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#stubbing_consecutive_calls)
- [ ] [Stubbing with callbacks](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#answer_stubs)
- [ ] Mock [variadic function](https://gobyexample.com/variadic-functions)
- [ ] Override existing mock, i.e. change return values of a stub
- [ ] Mock struct methods
  - [ ] [Making sure interaction(s) never happened on mock](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#never_verification)
  - [ ] [Finding redundant invocations](https://site.mockito.org/javadoc/current/org/mockito/Mockito.html#finding_redundant_invocations)

### Contributing

Rules to contribute to the repo:

1. Define ine identifier per file, which means that each go file contains either a struct (with related methods), an interface, or a function. Constants should be declared in the file `<package>.go`, i.e. in `mockit.go` for the `mockit` package.
2. Write unit test for each method/function, in order to keep the coverage to 100%.

### Internals

This library uses [monkey](https://github.com/bouk/monkey), a package for [monkey patching](https://en.wikipedia.org/wiki/Monkey_patch) in Go.

## Credits

All of this was possible only because of [bouk](https://github.com/bouk), as this library is basically a wrapper around [monkey](https://github.com/bouk/monkey), so kudos to him.
