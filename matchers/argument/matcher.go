package argument

// Matcher is a function used to match for specific arguments during mocking
type Matcher func(arg interface{}) bool
