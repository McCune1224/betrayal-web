package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertNoError(t testing.TB, err error, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NoError(t, err, msgAndArgs...)
}

func AssertError(t testing.TB, err error, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Error(t, err, msgAndArgs...)
}

func RequireNoError(t testing.TB, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
