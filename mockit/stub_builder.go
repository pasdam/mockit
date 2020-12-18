package mockit

import "reflect"

type stubBuilder struct {
	args      []reflect.Value
	mock      *instanceMock
	completed bool
}

func (b *stubBuilder) CallRealMethod() {
	b.assertUncompleted()

	b.mock.mockedCalls.Add(b.args, nil)
}

func (b *stubBuilder) Return(values ...interface{}) {
	b.assertUncompleted()

	typeOf := b.mock.target.Type()
	out := convertToValuesAndVerifies(b.mock.t, values, typeOf.NumOut(), typeOf.Out)

	b.mock.mockedCalls.Add(b.args, out)
}

func (b *stubBuilder) ReturnDefaults() {
	b.assertUncompleted()

	b.mock.mockedCalls.Add(b.args, b.mock.defaultOut)
}

func (b *stubBuilder) assertUncompleted() {
	if b.completed {
		b.mock.t.Error("The stub is already configured, please create a new one")
	}
}
