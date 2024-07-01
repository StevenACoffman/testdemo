package testdemo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsSortedF(t *testing.T) {
	f := func(array []int, expected bool) {
		t.Helper()
		actual := IsSorted(array)
		require.Equal(t, expected, actual)

	}

	f([]int{}, true)
	f([]int{0}, true)
	f([]int{0, 1}, true)  // actually true, but we want to see failures
	f([]int{1, 0}, false) // actually false, but we want to see failures
}
