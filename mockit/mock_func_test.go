package mockit

import (
	"path/filepath"
	"testing"

	"github.com/pasdam/mockit/matchers/argument"
	"github.com/stretchr/testify/assert"
)

func Test_MockFunc_Example_ShouldReturnExpectedValueForTheReadmeExample(t *testing.T) {
	// NOTE: if this fails (i.e. the contract changed), please update the README as well
	m := MockFunc(t, filepath.Base)
	m.With("some-argument").Return("result")
	assert.Equal(t, "result", filepath.Base("some-argument"))
}

func Test_MockFunc_Example_ShouldCaptureArgument(t *testing.T) {
	// NOTE: if this fails (i.e. the contract changed), please update the README as well
	m := MockFunc(t, filepath.Base)
	c := argument.Captor{}
	m.With(c.Capture).Return("result")
	filepath.Base("some-argument")
	assert.Equal(t, "some-argument", c.Value)
}

func Test_mockFunc_ShouldRemoveTheMockWhenTheTestCompletes(t *testing.T) {
	absPath, _ := filepath.Abs(".")

	t.Run("", func(t *testing.T) {
		m := MockFunc(t, filepath.Base)
		m.With(absPath).Return("mock-value")

		assert.Equal(t, "mock-value", filepath.Base(absPath))
	})

	assert.Equal(t, "mockit", filepath.Base(absPath))
}

func Test_mockFunc_ShouldUseArgumentMatcher(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*instanceMock)
	m.With(argument.Any).Return("result")

	assert.Equal(t, "result", filepath.Base("argument-1"))
	assert.Equal(t, "result", filepath.Base("argument-2"))
	assert.Equal(t, "result", filepath.Base("argument-3"))
	assert.Equal(t, 3, len(m.calls))
	m.Verify("argument-1")
	m.Verify("argument-2")
	m.Verify("argument-3")
}

func Test_mockFunc_ShouldReturnDefaultOutputIfNoMatchingCallIsFound(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*instanceMock)
	m.With("argument-1").Return("out-1")
	m.With("argument-2").Return("out-2")
	m.With("argument-3").Return("out-3")

	assert.Equal(t, "", filepath.Base("non-matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("non-matching-argument")
}

func Test_mockFunc_ShouldReturnAZeroValueIfTheMockArgumentIsNil(t *testing.T) {
	m := MockFunc(t, filepath.Walk).(*instanceMock)
	m.With("arg", nil).Return(nil)

	assert.Nil(t, filepath.Walk("arg", nil))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("arg", nil)
}

func Test_mockFunc_ShouldReturnExpectedOutputIfAMatchingCallIsFound(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*instanceMock)
	m.With("matching-argument").Return("some-out")

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))
	assert.Equal(t, 1, len(m.calls))
	m.Verify("matching-argument")
}

func Test_mockFunc_ShouldDisableAndRestoreAMock(t *testing.T) {
	m := MockFunc(t, filepath.Base).(*instanceMock)
	m.With("matching-argument").Return("some-out")

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	m.Disable()

	assert.Equal(t, "matching-argument", filepath.Base("matching-argument"))

	m.Enable()

	assert.Equal(t, "some-out", filepath.Base("matching-argument"))

	assert.Equal(t, 2, len(m.calls))
	m.Verify("matching-argument")
}
