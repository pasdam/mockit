package mockit

func configureMockReturn(f *funcMockData, values ...interface{}) {
	typeOf := f.target.Type()
	f.completeMock(f.convertToValuesAndVerifies(values, typeOf.NumOut(), typeOf.Out))
}
