package testdemo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Standard Table Driven Tests
func TestStdGoIsSorted(t *testing.T) {
	var tests = []struct {
		input []int
		want  bool
	}{
		{[]int(nil), true},
		{[]int{0}, true},
		{[]int{0, -9223372036854775808}, false},
		{[]int{0, 0}, true},
	}
	for _, test := range tests {
		got := IsSorted(test.input)
		require.Equal(t, test.want, got)
	}
}
