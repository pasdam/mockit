package mockit

// Stub contains method to mock a specific method
type Stub interface {

	// CallRealMethod makes sure that the mock perform a call to the real method
	CallRealMethod()

	// Return makes sure the mock to return the specified values
	Return(values ...interface{})

	// ReturnDefaults makes sure the mock to return the default outputs (zero values)
	ReturnDefaults()
}
