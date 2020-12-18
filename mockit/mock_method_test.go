package mockit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MockMethod_ShouldReturnExpectedValueForTheReadmeExample(t *testing.T) {
	// NOTE: if this fails (i.e. the contract changed), please update the README as well
	err := errors.New("some-error")
	m := MockMethod(t, err, err.Error)
	m.With().Return("some-other-value")
	assert.Equal(t, "some-other-value", err.Error())
}

func Test_MockMethod_ShouldMockOnlyTheSpecifiedInstance(t *testing.T) {
	err1 := errors.New("some-error")
	m := MockMethod(t, err1, err1.Error)
	m.With().Return("some-other-value")
	assert.Equal(t, "some-other-value", err1.Error())

	err2 := errors.New("some-other-error")
	assert.Equal(t, "some-other-error", err2.Error())
}
