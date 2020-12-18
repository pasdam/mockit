package mockit

import "reflect"

type callsIndex struct {
	in  [][]reflect.Value
	out [][]reflect.Value
}

func (i *callsIndex) Add(in []reflect.Value, out []reflect.Value) {
	// TODO: search for matching arguments and replace in case
	i.in = append(i.in, in)
	i.out = append(i.out, out)
}

func (i *callsIndex) MockedOutFor(in []reflect.Value) ([]reflect.Value, error) {
	index, err := findCall(i.in, in, func(fromCalls, in []reflect.Value) bool {
		return callsMatch(fromCalls, in, true)
	})
	if err != nil {
		return nil, err
	}

	return i.out[index], nil
}
