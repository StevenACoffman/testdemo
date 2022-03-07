package testdemo

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ExampleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ExampleTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ExampleTestSuite) TestExample() {
	// equivalent
	require.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	suite.Require().Equal(5, suite.VariableThatShouldStartAtFive)

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
	validate(suite.T(), testCase{Name: "Empty",
		Array:    []int{},
		Expected: true,
	})
	validate(suite.T(), testCase{Name: "Single element",
		Array:    []int{0},
		Expected: true,
	})
	validate(suite.T(), testCase{Name: "Two elements",
		Array:    []int{0, 1},
		Expected: false, // actually true, but we want to see failures
	})
	validate(suite.T(), testCase{Name: "Two elements unsorted",
		Array:    []int{1, 0},
		Expected: true, // actually false, but we want to see failures
	})

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ExampleTestSuite) TestExampleLogLinesLost() {
	//suite.T().Skip()

	// next two lines are equivalent
	require.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
	suite.Require().Equal(5, suite.VariableThatShouldStartAtFive)

	type testCase struct {
		Name     string
		Array    []int
		Expected bool
	}
	validate := func(suite *ExampleTestSuite, tc testCase) {
		suite.T().Helper()
		suite.Run(tc.Name, func() {
			suite.T().Helper()
			suite.T().Log("case:", tc.Name)
			actual := IsSorted(tc.Array)
			suite.Require().Equal(tc.Expected, actual)
		})
	}

	validate(suite, testCase{Name: "Empty",
		Array:    []int{},
		Expected: true,
	})
	validate(suite, testCase{Name: "Single element",
		Array:    []int{0},
		Expected: true,
	})
	validate(suite, testCase{Name: "Two elements",
		Array:    []int{0, 1},
		Expected: false, // actually true, but we want to see failures
	})
	validate(suite, testCase{Name: "Two elements unsorted",
		Array:    []int{1, 0},
		Expected: true, // actually false, but we want to see failures
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
