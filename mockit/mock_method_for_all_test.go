package mockit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MockMethodForAll_ShouldReturnExpectedValueForTheReadmeExample(t *testing.T) {
	// NOTE: if this fails (i.e. the contract changed), please update the README as well
	err := errors.New("some-real-error")
	m := MockMethodForAll(t, err, err.Error)
	m.With().Return("some-mocked-value")
	assert.Equal(t, "some-mocked-value", err.Error())
}

func Test_MockMethodForAll_ShouldMockAllInstances(t *testing.T) {
	err1 := errors.New("some-real-error")
	m := MockMethodForAll(t, err1, err1.Error)
	m.With().Return("some-mocked-value")
	assert.Equal(t, "some-mocked-value", err1.Error())

	err2 := errors.New("some-other-error")
	assert.Equal(t, "some-mocked-value", err2.Error())
}
