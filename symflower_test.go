package testdemo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsSorted(t *testing.T) {
	type testCase struct {
		Name     string
		Array    []int
		Expected bool
	}
	validate := func(t *testing.T, tc testCase) {
		t.Helper()
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			t.Log("case:", tc.Name)
			actual := IsSorted(tc.Array)
			require.Equal(t, tc.Expected, actual)
		})
	}
	validate(t, testCase{Name: "Empty",
		Array:    []int{},
		Expected: true,
	})
	validate(t, testCase{Name: "Single element",
		Array:    []int{0},
		Expected: true,
	})
	validate(t, testCase{Name: "Two elements",
		Array:    []int{0, 1},
		Expected: true, // actually true, but we want to see failures
	})
	validate(t, testCase{Name: "Two elements unsorted",
		Array:    []int{1, 0},
		Expected: false, // actually false, but we want to see failures
	})
}
