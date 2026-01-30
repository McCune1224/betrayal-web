package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertNoError(t *testing.T) {
	t.Run("passes when no error", func(t *testing.T) {
		mockT := &testing.T{}
		AssertNoError(mockT, nil, "should not fail")
		assert.False(t, mockT.Failed(), "should not fail when no error")
	})
}

func TestAssertError(t *testing.T) {
	t.Run("passes when error exists", func(t *testing.T) {
		mockT := &testing.T{}
		AssertError(mockT, assert.AnError, "should not fail")
		assert.False(t, mockT.Failed(), "should not fail when error exists")
	})
}

func TestRequireNoError(t *testing.T) {
	t.Run("continues when no error", func(t *testing.T) {
		called := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Expected for fatal
				}
			}()
			RequireNoError(t, nil, "should not fail")
			called = true
		}()
		assert.True(t, called, "should continue when no error")
	})

}
