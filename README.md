### Thoughts on testing in Go

As the famous programmer [Stephen Stills once sang](https://en.wikipedia.org/wiki/Love_the_One_You%27re_With):
> "🎶 And if you can't test all the code you love honey<br>
Love the tests you code, Love the tests you code,<br>
Love the tests you code, Love the tests you code.🎵"

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

### Dependency Injection

A dependency can be anything that effects the behavior or outcome of your logic. A real production application commonly grows to have more than one stateful dependency like:

+ A database (or NoSQL K/V store)
+ A cache
+ One or more HTTP clients
+ A message queue
+ Cloud APIs
+ Secrets (database passwords, API credentials, etc.)
+ Logging, Metrics, Tracing

You want different behavior for each of these types of dependencies in Prod vs. Staging vs. in Unit tests. 

The easiest place to inject those dependencies is just instantiating things in inside of main. It's very straightforward. In your main method, at runtime your program should interrogate which environment it is running in from environment variables.
 Pick an arbitrary but unlikely-to-be-in-use environment variable name like "MY_ORG_ENV" and set it to "prod" vs. "staging" vs. "unit", and panic if it's not one of those.

The easiest form of dependency injection, is just to pass them as function arguments, and then you get all the compiler help for free.

This can start to get cumbersome as the number of dependencies grows beyond more than a handful.

##### Shove all the stateful dependencies into a server Struct
The classic way to deal with lots of stateful dependencies in a HTTP Server application is to make all the HTTP handlers a method of the main Server struct and the handlers can just access the necessary dependencies from the server.
```
type Server struct {
  db Database
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) ([]string, errror) {
users, _ := s.db.getUsers(ctx))
...
}
```
### Alternative Dependency Injection Techniques

Some people don't like just passing things as arguments or putting them on a struct with receiver methods.

If you have a good reason to add extra complexity, then read on.

<details>
  <summary>Click to expand!</summary>

##### Stuff all the stateful dependencies into a "God Object" and pass that as an argument everywhere

```
type Deps struct {
  Db Database
}
func GetUsers(deps *Deps) ([]string, errror) {
  users, _ := deps.Db.getUsers(ctx))
  ...
}
```

##### Go interface wild!
Instead of storing the stateful dependency directly, you can also make methods like `GetDB()` that returns the stateful dependency.
```
s.GetDB().db.getUsers(ctx))
```
This is more verbose, but allows you to add some dynamic behavior, like lazy-loading, if you really need that.

##### Use the context.Context

In an HTTP Server, the `http.Server`, every request can retrieve a `context.Context` by just calling `req.Context()`

Idiomatically, `context.Context` is only for:
+ Cancellation signals
+ Deadlines
+ Request-scoped metadata that does not alter behaviour

However, we *can* put arbitrary `interface{}` values into any context, and any values in the BaseContext of the `http.Server` will be inherited by the `context.Context` for every request.
So we *can* inject shared, stateful dependencies into there if we really want to.

Before `http.Server.BaseContext` existed, Kayle Gishen [described injecting dependencies into the request context using a middleware function](https://www.youtube.com/watch?v=_KrV_VWP2n0) with [source code here](https://github.com/kayleg/yt-dependency-injection)
and someone else [summarized it here](https://www.adityathebe.com/journal/5).

There [are some downsides to using context for dependency injection](https://ahmedalhulaibi.com/blog/go-context-misuse/).
+ Using `context.WithValue()` and `context.Value()` is actively choosing to give up information and type checking at compile time.
+ [Obfuscates input](https://ahmedalhulaibi.com/blog/go-context-misuse/#obfuscated-inputs) when reading method and function signatures.
+ [Creates implicit couplings](https://ahmedalhulaibi.com/blog/go-context-misuse/#implicit-and-unclear-temporal-coupling) which slows down refactoring.
+ [Leads to nil pointer exceptions](https://ahmedalhulaibi.com/blog/go-context-misuse/#nil-pointer-exceptions) causing development delays and service disruptions.
+ Not idiomatic - Bespoke solutions to common problems divorce you from benefiting from the wider ecosystem 

You can mitigate some of these with various tricks.

To get back the compiler type safety, add getter/setter helpers in other packages that define context keys as an *unexported* type. There's no way to set them to the wrong type,
and only one way to retrieve these values:
```
type userCtxKeyType string

const userCtxKey userCtxKeyType = "user"

func WithUser(ctx context.Context, user *User) context.Context {
  return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) *User {
  user, ok := ctx.Value(userCtxKey).(*User)
  if !ok {
    // Log this issue
    return nil
  }
  return user
}
```

To avoid obfuscating function inputs, before we ever actually use data from context values, we write a function to pull data from the context values and then pass that data into a function that explicitly states the data it requires. After doing this, the function that we call should never need to pull additional data out of the context that affects the flow of our application.

##### Combine interfaces and God Object and extend Context into MegaContext!
In GraphQL resolvers where you have a `ctx context.Context` and HTTP handlers where you have a request that can give you the same,
you can "upgrade" to a custom context.
```
var ktx interface {
	customctx.Base
	log.CustomContext
	datastore.CustomContext
} = customctx.Upgrade(ctx)
...
```
You can then use it like this:
```
func GetUsers(ctx interface {
  customctx.Base
  customctx.DB
  customctx.Service
  customctx.Time
  customctx.Log
}) ([]string, errror) {
...
}
```
The benefit here is that you are as explicit about the defining the dependencies of the function as if you had
passed them as individual arguments. While it makes the function definition even more verbose than individual arguments, it is much *less* verbose to *call* `GetUsers(ktx)`.
You are still obscuring the dependencies at the call site.

In order for this to work, you need even more tricks:

```
// We store the CustomContextas a Context.Value-style key in the go-context it
// wraps, so that we can re-extract all the CustomContext goodies after someone
// else wraps it go-context style (as, for example, the HTTP server and
// middleware will do).  This type and key ensure, in the usual way, that
// collisions in context-keys are impossible.
//
// The type of the value is *customContext.
type _customContextKeyType string

const _customContextKey _customContextKeyType = "customctx.customContext"

type customContext struct {
	// NOTE(benkraft): Do NOT replace Context after initialization; use
	// WithContext instead.  (See comments there for why.)
	context.Context
	...
}

// Base contains the functionality to convert a CustomContext to and from an ordinary
go-context: 
//  + it ensures that CustomContexts are valid go-contexts 
//  + it provides WithContext which can convert the other direction.
type Base interface {
	// Embedding context.Context is what allows us to use a KA context as a Go
	// context.
	context.Context

	// WithContext replaces the Go-context in the CustomContext with another.
	//
	// This is useful if you want to create a modified context -- say add a
	// deadline -- but keep using the CustomContext extras.  For example:
	//	// Add a deadline to a CustomContext:
	//	ktx = ktx.WithContext(context.WithDeadline(ktx, deadline))
	//	// Add a context-value to a CustomContext:
	//	ktx = ktx.WithContext(context.WithValue(ktx, key, value))
	//	// Create a new background CustomContext (for a fire-and-forget call):
	//	ktx = ktx.WithContext(context.Background())
	//
	//
	// It's also used internally in Upgrade, to ensure that all the
	// wrapping done as the context is passed through HTTP-land is applied when
	// we re-extract the CustomContext.
	WithContext(ctx context.Context) *customContext

	// These replacers implement the clone-replace pattern used in tests, for
	// cases where we also want to allow it in prod.
	// WithHTTP is useful for clients that want to change cookie- or
	// redirect-handling, for example.
	WithHTTP(*http.Client) *customContext
	// WithMemorystore lets clients install a mock memorystore pool.
	WithMemorystore(*memorystore.Client) *customContext
	// WithDatastore lets clients install a mock datastore client.
	WithDatastore(datastore.Client) *customContext

	// Detach replaces all request-specific context in this
	// CustomContext with something that does not depend on the request.
	//
	// This is intended to be used when starting a go-routine that
	// should last beyond the given request.  We want to make sure
	// that the context isn't canceled when the request finishes.  We
	// also want to ensure its logging is not associated with this
	// request.  Detach does all this.  It returns a new context that
	// is typed as context.Context but can be promoted to whatever
	// type it was before the Detach() call.
	//
	// We want to make sure things using the detached context don't
	// run forever, unless that's what you want, so we ask callers to
	// pass in a timeout.  If you *want* to run forever, pass a zero
	// timeout, and we'll just make the new context cancelable but not
	// set a timeout.  In either case we return a cancel-function so
	// you can cancel this context (and hopefully whatever goroutine
	// is using it) manually.
	Detach(timeout time.Duration) (context.Context, context.CancelFunc)
}
```

</details>

### Testing Terms

+ **Stub** - an object that provides predefined answers to method calls.
+ **Mock** - an object on which you set expectations.
+ **Fake** - an object with limited capabilities (for the purposes of testing), e.g. a fake web service.

Test Double is the general term for stubs, mocks and fakes. A mock is single use, but a fake can be reused.

+ **Test Fixture** -  a well known and fixed environment in which tests are run so that results are repeatable. Some people call this the test context.

##### Examples of fixtures:

+ Preparation of input data and set-up/creation of fake or mock objects
+ Loading a database with a specific, known set of data
+ Erasing a hard disk and installing a known clean operating system installation
+ Copying a specific known set of files

Test fixtures contribute to setting up the system for the testing process by providing it with all the necessary data for initialization. The setup using fixtures is done to satisfy any preconditions there may be for the code under test. Fixtures allow us to reliably and repeatably create the state our code relies on upon without worrying about the details.
### HTTP client testing techniques:
+ [Unit Testing http client in Go](http://hassansin.github.io/Unit-Testing-http-client-in-Go)
+ [a way to test http client in go](https://blog.bullgare.com/2020/02/a-way-to-test-http-client-in-go/)

1.  Using `httptest.Server`:
    `httptest.Server` allows us to create a local HTTP server and listen for any requests. When starting, the server chooses any available open port and uses that. So we need to get the URL of the test server and use it instead of the actual service URL.

2. [Accept a Doer as a parameter](https://www.0value.com/let-the-doer-do-it)
   The Doer is a single-method interface, as is often the case in Go:
   ```
   type Doer interface {
       Do(*http.Request) (*http.Response, error)
   }
   ```
3. By Replacing `http.Transport`:
   Transport specifies the mechanism by which individual HTTP requests are made. Instead of using the default http.Transport, we’ll replace it with our own implementation. To implement a transport, we’ll have to implement http.RoundTripper interface. From the documentation:
   ```
   func Test_Mine(t *testing.T) {
       ...
       client := httpClientWithRoundTripper(http.StatusOK, "OK")
       ...
   }

   type roundTripFunc func(req *http.Request) *http.Response

   func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
       return f(req), nil
   }

   func httpClientWithRoundTripper(statusCode int, response string) *http.Client {
       return &http.Client{
           Transport: roundTripFunc(func(req *http.Request) *http.Response {
               return &http.Response{
                   StatusCode: statusCode,
                   Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
               }
           }),
       }
   }
   ```

