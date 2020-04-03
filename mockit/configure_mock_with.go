package mockit

func configureMockWith(f *funcMockData, values ...interface{}) {
	typeOf := f.target.Type()
	f.currentMock = &funcCall{
		in: f.convertToValuesAndVerifies(values, typeOf.NumIn(), typeOf.In),
	}
}
