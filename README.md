
As the famous programmer [Stephen Stills once sang](https://en.wikipedia.org/wiki/Love_the_One_You%27re_With):
> "🎶 And if you can't test with the code you love honey<br>
Love the tests you code, Love the tests you code,<br>
Love the tests you code, Love the tests you code.🎵"

This started from me cribbing notes from [this excellent article](https://symflower.com/en/company/blog/2022/better-table-driven-testing/)
and the [Reddit Comments](https://www.reddit.com/r/golang/comments/t3hhh6/take_on_a_better_unit_test_style/).

### Example of testing some Go code
Let us assume you have something simple to test like this:
```
// IsSorted reports whether data is sorted.
func IsSorted(data []int) bool {
	n := len(data)
	if n == 0 || n == 1 {
		return true
	}
	i := 0
	for i < n-1 && data[i] <= data[i+1] {
		i = i + 1
	}
	return i == n-1
}
```
How do you unit test this in Go?

### Beginner SideNote: AAA pattern
The AAA (Arrange, Act, Assert) pattern is a common way of writing unit tests for a method under test.

+ The **Arrange** section of a unit test method initializes objects and sets the value of the data that is passed to the method under test.
+ The **Act** section invokes the method under test with the arranged parameters.
+ The **Assert** section verifies that the action of the method under test behaves as expected. 

Sometimes it is more natural to think of these as `Given`/`When`/`Then`, especially if you are more familiar with Behavior Driven Development (BDD).
+ Given = Arrange
+ When = Act
+ Then = Assert

### Basic Function Per Test Case Style

This is the most basic test style. Since every test case gets its own test function, there is a lot of redundant code, and the specific behavior a test case should check is also not described.

Disadvantages:
+ Redundant validation code
+ Overview of existing tests is not great
+ Long and hard to read test function names

<details>
  <summary>Click to expand!</summary>

```
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
```
</details>

Characteristics of a Golang test function:

+ The first and only parameter needs to be `t *testing.T`
+ It begins with the word `Test` followed by a word or phrase starting with a capital letter.
(usually the method under test i.e. `TestValidateClient`)
+ Calls `t.Error` or `t.Fail` to indicate a failure (or let `testify`'s require do that)
+ `t.Log` can be used to provide non-failing debug information
+ Must be saved in a file name ending with `_test.go` such as: `addition_test.go`

Helpful Testing methods:
+ 💨 `t.Run` to give your test cases subtests
+ ⏭ `t.Skip`, for when we only want to run a test sometimes
+ 🧹 `t.Cleanup`, for cleaning up state in between tests
+ 🙈 `t.Helper`, mark this function to be skipped when printing stacktraces.

### Standard Go Table-Driven style
This is the style you see most often in Go projects. It’s already table-driven, which is a huge improvement to the function-per-test-case style, since it reduces the amount of redundant code. However, there are still disadvantages:

+ No relevant stack traces for failing tests
+ No description of the behavior the test cases should check
+ Missing field names make tests often hard to follow
+ More fields and data make these tests often very hard to read and maintain

<details>
  <summary>Click to expand!</summary>

```
package testdemo

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Standard Table Driven Tests
func TestStdGoIsSorted1(t *testing.T) {
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

```

</details>


### Symflower-style table-driven unit tests
By calling a validation function directly instead of using a for loop, you can pinpoint problematic cases with stacktraces that have accurate line numbers.
Finding the test case that failed in a test case that is 1000+ lines can be difficult if test names are similar. But with this style you get an exact line number for where the test case is defined.

Characteristics:
+ Table-driven
+ Validation function
+ Named test cases
+ Structured test cases

Benefits:
+ Meaningful stack traces for failing tests 
+ Improved readability and understandability by named and structured test cases
+ Reuse of validation code

<details>
  <summary>Click to expand!</summary>

```
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
		Expected: false, // actually true, but we want to see failures
	})
	validate(t, testCase{Name: "Two elements unsorted",
		Array:    []int{1, 0},
		Expected: true, // actually false, but we want to see failures
	})
}
```

</details>

The Symflower-Style is an interesting testing pattern that also appears in https://github.com/hashicorp/consul 
with only slight variations. Some examples in the wild can be seen here:

+ [consul/leader_connect_test.go#L1251-L1267](https://github.com/hashicorp/consul/blob/42ec34d/agent/consul/leader_connect_test.go#L1251-L1267) - uses the same ordering (testCase type, then function, then cases), but not the function call to replace the slice and for loop

+ [config/runtime_test.go#L65-L84](https://github.com/hashicorp/consul/blob/7b0548dd8d78f3a02c5763c11be1232c21b06643/agent/config/runtime_test.go#L62-L84) - uses the function call to replace the slice and for loop, but keeps the type outside of the case, since it is re-used.

+ [state/catalog_events_test.go#L304-L312](https://github.com/hashicorp/consul/blob/7b0548dd8d78f3a02c5763c11be1232c21b06643/agent/consul/state/catalog_events_test.go#L304-L312) - an example of both together

The small difference in consul is that in the original Symflower article they needed to fork Testify to get good stacktraces, but Consul has 
a trick to remove the need to do that. They call `t.Run` from a helper outside the test function. In the examples a helper that calls `t.Run` is called from either the for loop, or from a wrapper that handles printing the right test line.

In consul, they use [runCase](https://github.com/hashicorp/consul/blob/7b0548dd8d78f3a02c5763c11be1232c21b06643/agent/config/runtime_test.go#L5272-L5279). This prints the line number of where the test case is called, so you can easily jump to either the line that failed in the test function, or the case that failed.

```
func runCase(t *testing.T, name string, fn func(t *testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		t.Log("case:", name)
		fn(t)
	})
}
```

### Testify Testing Framework
##### Testify Suites
A `testify` [suite](https://pkg.go.dev/github.com/stretchr/testify/suite) works by taking in a `*testing.T` value and running each suite method whose name starts with `Test` as a subtest.

+ Testify `suite` definitions use struct embedding to define the suite, and absorb the built-in basic suite functionality from testify
  ```
  type ExampleTestSuite struct {
      suite.Suite
      // ... add shared state here for all the tests in the suite
    }
  ```
+ All your suite methods that begin with "Test" are run as tests within a suite.
+ In order for 'go test' to run this suite, we need to create a normal test function and pass our suite to `suite.Run`
  ```
  func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
  }
  ```
+ Lifecycle methods - order can be [seen here](https://go.dev/play/p/PUzY9YjnC15)
  + `SetupSuite` - useful only in cases where the setup code is time-consuming and isn't modified in any of the tests
  + `SetupTest` - each individual test function runs with a clean environment.
  + `BeforeTest` - mostly for logging as it executes right before the test starts and receives the suite and test names as input
  + `AfterTest` - Good for cleanup
+ `suite.T()` - Get the test context (`t *testing.T`) to use standard Go Test methods like `Run`, `Skip`, `Cleanup`,`Helper`, `Log`

##### Testify Assertions (Require)
The `require` package provides helpful functions for asserting the expected outcome of a test case. Optionally, you can also provide a helpful failure description.

99% of the time you want to use `require` instead of `assert`. Require will stop testing on a failure, `assert` will just continue along.
```
    // these next two lines are equivalent btw:
	require.Equal(suite.T(), want, got)
	suite.Require().Equal(want, got)
```
Most useful (see [here](https://vyskocil.org/blog/testify-make-go-testing-easy-2/)) assertions:
+ `Equal` / `NotEqual`
+ `Error` / `NoError`
+ `JSONEq` / `YAMLEq`
+ `ElementsMatch`

The complete list is [fairly exhaustive](https://pkg.go.dev/github.com/stretchr/testify/assert#pkg-functions)

### Running a single subtest
Go test has a `-run` flag that takes a slash-separated list of regular expressions that match each name element in turn. For example:

```
$ go test -run TestExampleTestSuite/^TestExample$/Two_elements_unsorted ./...
--- FAIL: TestExampleTestSuite (0.00s)
    --- FAIL: TestExampleTestSuite/TestExample (0.00s)
        --- FAIL: TestExampleTestSuite/TestExample/Two_elements_unsorted (0.00s)
            testify_suite_test.go:56: case: Two elements unsorted
            testify_suite_test.go:56:
                	Error Trace:	testify_suite_test.go:41
                	Error:      	Not equal:
                	            	expected: true
                	            	actual  : false
                	Test:       	TestExampleTestSuite/TestExample/Two_elements_unsorted
FAIL
FAIL	github.com/StevenACoffman/testdemo	0.008s
FAIL

```
Here you can see the `testify` suite is itself a top-level test, while the suite's test methods are [subtests][3], and anything below that are then sub-sub-tests. Some other examples:
>     go test -run Foo     # Run top-level tests matching "Foo".
>     go test -run Foo/A=  # Run subtests of Foo matching "A=".
>     go test -run /A=1    # Run all subtests of a top-level test matching "A=1".

How does this help you in table driven tests? The names of subtests are `string` values, which can be generated on-the-fly, e.g.:

    for i, c := range cases {
        name := fmt.Sprintf("C=%d", i)
        t.Run(name, func(t *testing.T) {
            if res := myfn(c.arg); res != c.expected {
                t.Errorf("myfn(%q) should return %q, but it returns %q",
                    c.arg, c.expected, res)
            }
        })
    }

To run the case at index `2`, you could start it like

    go test -run /C=2

or

    go test -run TestName/C=2

  [3]: https://tip.golang.org/doc/go1.7#testing
  [4]: https://tip.golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks
