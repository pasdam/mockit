package argument

// Captor is an argument matcher that stores the received value
type Captor struct {

	// Value is the captured value
	Value interface{}
}

// Capture is the argument matcher that capture the value
func (c *Captor) Capture(arg interface{}) bool {
	c.Value = arg
	return true
}
