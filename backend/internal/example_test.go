package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example demonstrates using testify's assert and require in tests.
func TestExample_Sum(t *testing.T) {
	sum := func(a, b int) int { return a + b }

	result := sum(2, 3)
	assert.Equal(t, 5, result, "sum(2,3) should return 5") // continue on failure

	require.NotEqual(t, 6, result, "sum(2,3) should not return 6") // aborts on failure
}
