package testdemo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPerFunctionEmptyIsSorted(t *testing.T) {
	data := []int(nil)
	actual := IsSorted(data)
	expected := true
	require.Equal(t, expected, actual)
}
func TestPerFunctionOneElementIsSorted(t *testing.T) {
	data := []int{0}
	actual := IsSorted(data)
	expected := true
	require.Equal(t, expected, actual)
}
func TestPerFunctionUnsortedIsNotSorted(t *testing.T) {
	data := []int{0, -9223372036854775808}
	actual := IsSorted(data)
	expected := false
	require.Equal(t, expected, actual)
}
func TestPerFunctionTwoEqualIsSorted(t *testing.T) {
	data := []int{0, 0}
	actual := IsSorted(data)
	expected := true
	require.Equal(t, expected, actual)
}
