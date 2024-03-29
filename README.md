# Effective Go

[Learn Go Programming](https://blog.learngoprogramming.com/)

# Testing

In Go, a _unit_ is a package, and a _unit test_ verifies the behavior of a single package.

The `*testing.T` type uses signaling methods to communicate when a test fails.

A common way to check test failures is with `t.Logf()` followed by `t.Fail()`. `t.Fail()` marks the test as failed, but the test keeps running. For example:

```go
if err != nil {
    t.Logf("Parse(%q) err = %q, want nil", rawurl, err)
    t.Fail()
}
```

To make a test fail and stop execution, use the `t.FailNow()` method.

## External and internal tests

Put both internal and external tests in the same folder as the code they test.

_Internal tests_ verify code from the same package. They are called 'white-box tests' and can test exported and unexported identifiers.

_External tests_ verify code from another package. They are called 'black-box tests'. External tests use the `_test` suffix for the package name. For example, `package url_test`.

If you want to test an unexported function from an external test, you must export the function from the external test package. For example, the `parseScheme` function is unexported:

```go
// url.go
package url

func parseScheme() {...}
```

Create a new file in the same package that assigns the unexported function to an exported function:

```go
// export_test.go
package url

var ParseScheme = parseScheme
```

Finally, test the function in an external test file:

```go
// parse_scheme_test.go
package url_test

func TestParseScheme(t *testing.T) {...}
```


## Test coverage

View how much of your code is covered by tests:
```shell
$ go test -coverprofile cover.out`
```
The coverage output file is optional. By convention, the coverage profile file is named `cover.out`.

Use a coverage output file to view test coverage by function:
```shell
$ go tool cover -func=cover.out
```

After you create the `coverprofile` file, use the `cover` go tool to generate HTML output to view what code is and is not covered. This command opens the coverage output in the browser:
```shell
$ go tool cover -html=cover.out
```

To create the HTML file, but not open it in the browser automatically:
```shell
$ go tool cover -html=cover.out -o coverage.html
```
The cover tool uses three colors to identify code coverage:
- grey: not tracked by the coverage tool
- green: sufficiently tested
- red: not covered by tests

## Benchmarks

Benchmark your methods to determine their efficiency. Benchmark functions use the `BenchmarkXxx(b *testing.B)` signature. Place them in the `x_test.go` file with the other test functions.

To write a benchmark function, call the method that you want to benchmark within a `for` loop in the `BenchmarkXxx` function. The `for` loop uses `b.N` as its upper bound. `b.N` helps adjust the test runner to properly measure performance. In addition, you can run the `b.ReportAllocs()` function to see how many memory allocations your code makes.

For example, the following function benchmarks the `String()` method on the `URL` type:

```go
func BenchmarkURLString(b *testing.B) {
	b.ReportAllocs()
	b.Logf("Loop %d times\n", b.N)

	u := &URL{Scheme: "https", Host: "foo.com", Path: "go"}

	for i := 0; i < b.N; i++ {
		u.String()
	}
}
```
To run the test, use the `-bench` flag. Use a dot (`.`) to run every benchmark in the package:

```shell
$ go test -bench .
...
BenchmarkURLString-12    	 8868142	       153.8 ns/op	      64 B/op	       4 allocs/op
--- BENCH: BenchmarkURLString-12
    ...
PASS
ok  	url/url	1.506s
```

The `B/op` column indicates that there were 64 bytes allocated in each operation. The `allocs/op` value indicates the number of memory allocations that made by the code in the benchmark.

When you run benchmarks with the `-bench` flag, the regular tests run as well (use the `-v` flag to verify). If you want to run only the benchmark tests, use the `-run` flag with the `^$` regular expression:

```shell
$ go test -run=^$ -bench .
```

The `^$` regex tells the runner to ignore tests other than the benchmarks.

### Sub-benchmarks

You can run sub-benchmarks, just as you can run subtests:

```go
func BenchmarkURLString(b *testing.B) {
    var benchmarks = []*URL{
        {Scheme: "https"},
        {Scheme: "https", Host: "foo.com"},
        {Scheme: "https", Host: "foo.com", Path: "go"},
    }
    for _, u := range benchmarks {
        b.Run(u.String(), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                u.String()
            }
        })
    }
}
```

### Comparing benchmarks

1. Save the current benchmark result of the method:
```shell
go test -bench . -count 10 > old.txt
```
The `-count` flag runs the benchmark the number of times that you pass to it. There is no recommendation for the number you pass to `count`--`10` is a random number.

2. Refactor your code.
3. Run the benchmarks again and compare with `benchstat`.
   First, install `benchstat`:
   ```shell
   $ go install golang.org/x/perf/cmd/benchstat@latest
   ```
   Next, compare the `old.txt` and `new.txt` files:
   ```shell
   $ benchstat old.txt new.txt 
   goos: linux
   goarch: amd64
   pkg: url/url
   cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
                │   old.txt    │               new.txt                │
                │    sec/op    │    sec/op     vs base                │
   URLString-12   138.85n ± 5%   99.70n ± 14%  -28.19% (p=0.000 n=10)

                │  old.txt   │              new.txt               │
                │    B/op    │    B/op     vs base                │
   URLString-12   64.00 ± 0%   56.00 ± 0%  -12.50% (p=0.000 n=10)

                │  old.txt   │              new.txt               │
                │ allocs/op  │ allocs/op   vs base                │
   URLString-12   4.000 ± 0%   3.000 ± 0%  -25.00% (p=0.000 n=10)
   ```

## Methods

The most useful are `t.Errorf()` and `t.Fatalf()`. The following table describes all available `t.*` test methods:

| Method          | Description |
|-----------------|:------------|
| `t.Log()`        | Log a message. |
| `t.Logf()`       | Log a formatted message.|
| `t.Fail()`       | Mark the test as failed, but continue test execution. |
| `t.FailNow()`    | Mark the test as failed, and immediately stop execution. |
| `t.Error()`      | Combination of `Log()` and `Fail()`. |
| `t.Errorf()`     | Combination of `Logf()` and `Fail()`. |
| `t.Fatal()`      | Combination of `Log()` and `FailNow()`. |
| `t.Fatalf()`     | Combination of `Logf()` and `FailNow()`. |

## Table-driven tests

Also called data-driven and parameterized tests. They verify code with varying inputs. You can also implement subtests that run tests in isolation.

Imagine table-driven tests as actual tables, where the headers are struct fields, and the rows become individual slices in the test cases:

| product     | rating  | price |
|:------------|---------|-------|
| prod one    | 5       | 20    |
| prod two    | 10      | 30    |
| prod three  | 15      | 40    |

You can represent this in a test as follows:

```go
func TestTable(t *testing.T) {
    type product struct {
        product string
        rating  int
        price   float64
    }
    testCases := []product {
        {"prod one", 5, 20},
        {"prod two", 10, 30},
        {"prod three", 15, 40},
    }
}
```
## Testing the main method

The `main` method uses globals, which you cannot easily test. The trick is to put the logic into the `run()` function. You can inject the globals to the `run` function, and you can mimic those values during tests.

The flag package uses `NewFlagSet` function to create the default flag set, `CommandLine`.

## Build tags 

Add build tags at the top of test files so you can specify whether or not you want to run them in your test commands. For example, the following build tag designates the file for go tests that use the `cli` tag:

```go
//go:build cli
```
To run tests in this file, run the following command:
```shell
$ go test -tags=cli
```
For additional details, see [Build constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints).

# Packages

A package name should describe what it provides, not what it does.

When you write external tests, use a `_test` suffix. For example, a package that contains external tests for the `url` package is `url_test`.

## Import external packages

When you import an external pacakge, you list the module name in the `go.mod` file, followed by the path to the specific library from the project root. For example, if the `go.mod` file contains the following:

```go 
module url
...
```

Then you import the module as follows:
```go
import "url/path/to/library"
```

Commonly, packages are publically available in repositories, and the module name is the path to the root of the repository:

```go 
module github.com/rjs/url-parser
...
```

In this case, the import statement for the `parser` package within this repo is as follows:

```go
import "github.com/rjs/url-parser/parser"
```


# Formatting verbs

| verb | Definition | Usage |
|------|:-----------|:------|
| %q   | Wraps the given string in double quotes. | |
| %#v  | Prints the Go syntax representation of the value. | t.Errorf("%#v.String()\ngot  %q\nwant %q", u, got, want) |
| %t   | Boolean values. |

# nil

You can execute a method on a `nil` type. A method is a function that takes the receiver as a hidden first parameter. So, when you have a `nil` type, Go can find the method function to run, but it does not have anything to execute it on. For example, the Go compiler does the following when calling the `String()` method on a `nil` typ:

```go
var u *URL
u.String() // (*url.URL).String(u)
```

# Documentation 

## Testable examples

[Blog post](https://go.dev/blog/examples)

A _testable example_ is live documentation for code. You write a testable example to demonstrate the package API to other developers. The API includes the exported identifiers, such as functions, methods, etc. A testable example never goes out of date.

The testing package runs testable examples and checks their results, but it does not report successes or failures:

```go
func ExampleURL_fields() {
	u, err := url.Parse("https://foo.com/go")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Scheme)
	fmt.Println(u.Host)
	fmt.Println(u.Path)
	fmt.Println(u)
	// Output:
	// https
	// foo.com
	// go
	// https://foo.com/go
}
```

### Naming conventions

Testable examples use the following naming conventions:

| Signature                  | Description |
|:---------------------------|:------------|
| `func Example()`             | Example for the entire package. |
| `func ExampleParse()`        | Example for the `Parse` function. |
| `func ExampleURL()`          | Example for the `URL` type. |
| `func ExampleURL_Hostname()` | Example for the `Hostname` method on the `URL` type. |

## godoc server

You can generate docs that include your [testable examples](#testable-examples) with `godoc`. The following command installs the latest version:

```shell
$ go install golang.org/x/tools/cmd/godoc@latest
```

To view any `ExampleXxx` functions as Go documentation, run the `go doc` server with the following command:
```shell
$ godoc -play -http ":6060"
```
To show additional examples of the same type, use the `_xxx()` suffix on the function name. For example:
```go
func ExampleURL(){...}
func ExampleURL_fields(){...}
```


# Interfaces

When a type satisfies an interface, you say _type X is a Y_. For example, _URL is a Stringer_ or _Parser is a Reader_.

## Empty interface

Go versions prior to 1.18 used the empty interface: `interface{}`. This is an interface that did not implement any methods, so any type satisfied it. In Go 1.18 and later, `interface{}` was replaced with `any`. 

## Stringer

`Stringer` prints the string representation of the object. The `fmt.Print[x]` packages detect when a type has a `Stringer` method, so it calls that method for proper formatting.

The `Stringer` interface:
```go
type Stringer interface {
    String() string
}
```

Implementation example:
```go
func (u *URL) String() string {
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
}
```
### testString

Create a `testString()` function to return a concise string representation value for tests. For example, the following is the test version of the `String()` function in the previous section:

```go
func (u *URL) testString() string {
	return fmt.Sprintf("scheme=%q, host=%q, path=%q", u.Scheme, u.Host, u.Path)
}
```

# Misc

## strconv.ParseInt

Prefer `strconv.ParseInt()` over `strconv.Atoi()` because the former can parse hex and numbers with underscores (`1_000`).

## Short-if declaration

`if variable := value; condition`

For example:
```go
if err := json.Marshal(&val); err != nil {
    // handle error
}
```

## Named return values

Name the returned values of a function so other developers can see what it returns without having to read the code.

## ok return value

If you are writing a helper function, do NOT return an error from the helper--return `ok bool`. This allows you to check the `ok` value in the caller and return the error there. For example:

```go
func Parse(rawurl string) (*URL, error) {

	scheme, rest, ok := parseScheme(rawurl)
	if !ok {
		return nil, errors.New("missing scheme")
	}
    ...
}

func parseScheme(rawurl string) (scheme, rest string, ok bool) {
	i := strings.Index(rawurl, "://")
	if i < 1 {
		return "", "", false
	}
	return rawurl[:i], rawurl[i+3:], true
}
```

## Naked returns

You can return from a function with just the `return` keyword. This is called a _naked return_. A naked return returns the current state of the result values.

Generally, **do NOT** use naked returns because they impact readability.

## blank identifiers

You can use blank identifiers in function definitions:
```go
handler := func(_ http.ResponseWriter, _ *http.Request) {
	// logic
}
```

# Cross-compilation

You need to know the `GOOS` and `GOOARCH` values to compile the correct binaries.

Next, you can cross compile for multiple operating systems with a Makefile. Create a make target that compiles multiple binaries and places them in the `/bin` directory:

```makefile
compile:
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/hit_linux_amd64 ./cmd/hit
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o ./bin/hit_darwin_amd64 ./cmd/hit
	# windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/hit_win_amd64.exe ./cmd/hit
```
> Make sure that you add the `/bin` directory to the `.gitignore` file.


# CLI tools

## flag package

--- moved to docs
Go provides flag definition functions for common primitive types (`string`, `int`, etc.). A flag definition contains information about the flag such as defaults and usage information. The flag package parses command line flags with this flag definition.

The flag definition function can create internal variables for the flag and return a pointer to that variable, or it can use a variable that you define. For example, if you are defining a `string` flag, you use the `flag.String(...)` function to return a pointer to an internal variable, or you can use the `flag.StringVar(&userVar, ...)` to provide your own variable for the flag definition. `flag.*Var()` functions provide more control over variable definitions.

Each flag definition is saved in a structure called `*Flagset` for tracking. The flag package uses the `CommandLine` flag set when you define a flag.

The `Parse()` function extracts each command line flag in the `*Flagset` and creates name/value pairs, where the name is the flag name, and the value is the argument provided to the flag. Next, it updates any command line flag's internal variable.

---

### Changing flag usage type

You can replace the type that displays beside the flag in usage. In the usage string, enclose the replacement word in backticks (\`\`). For example, the following `flag.StringVar()` function accepts a `string` type by default. You can change that to a `URL` type with backticks:
```go
flag.StringVar(&f.url, "url", "", "HTTP server `URL` to make requests (required)")
```


### Manual implementation

```go
type flags struct {
	url  string
	n, c int
}

// parseFunc is a command-line flag parser function
type parseFunc func(string) error

func (f *flags) parse() (err error) {
	// map of flag names and parsers
	parsers := map[string]parseFunc{
		"url": f.urlVar(&f.url),
		"n":   f.intVar(&f.n),
		"c":   f.intVar(&f.c),
	}

	for _, arg := range os.Args[1:] {
		n, v, ok := strings.Cut(arg, "=")
		if !ok {
			continue // can't parse the flag
		}
		parse, ok := parsers[strings.TrimPrefix(n, "-")]
		if !ok {
			continue // can't find parser
		}
		if err := parse(v); err != nil {
			err = fmt.Errorf("invalid value %q for flag %s: %w", v, n, err)
			break
		}
	}
	return err
}

func (f *flags) urlVar(p *string) parseFunc {
	return func(s string) error {
		_, err := url.Parse(s)
		*p = s
		return err
	}
}

func (f *flags) intVar(p *int) parseFunc {
	return func(s string) (err error) {
		*p, err = strconv.Atoi(s)
		return err
	}
}
```

## Custom flag types

First, create a new type that satisfies the `Value` interface:
```go
type Value interface {
    Set(string)
    String() string
}
```

Next, register the type to the default flag set with `Var()`. Then, `Parse` can handle the flag.

## Positional arguments

Define `flag.Usage` as a function that prints usage text that is defined as a variable, and then the usage messages for the optional arguments:

```go
const usageText = `
Usage:
  hit [options] url
Options:`

...
func funcName() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText[1:])
		flag.PrintDefaults()
	}

	flag.Var(toNumber(&f.n), "n", "Number of requests to make")
	flag.Var(toNumber(&f.c), "c", "Concurrency level")
	flag.Parse()

	f.url = flag.Arg(0)
    
    ...
}
```
In the previous example:
- `url` is the positional argument. It is included in the `usageText` constant.
- `flag.PrintDefaults()` method prints the usage information for the 
- `flag.Arg(0)` stores the first argument after the flag.


## Directory structure

You create a tool that the user interacts with and is responsible for the following:
- Parses the flags
- Validates flags
- Calls the business logic library

Go uses the `cmd` directory for executables (the entry point) such as CLI tools. Within each `cmd/subdirectory`, you can name the entry point `main.go` or the name of the package, such as `hit.go`. Regardless of the file name, it must be in the `main` package because that package is what makes a file executable.

Next, you have to create the tool library that contains the business logic. This is a standalone package, so use the name of the library that you are building.

The following is a simple directory structure for the `hit` tool:

```shell
hit-tool
├── cmd         # Executable directory
│   └── hit     # CLI tool directory
├── go.mod
└── hit         # Library directory

```

# HTTP Clients

The HTTP protocol exchanges plain text messages between a server and client: the client sends a request with a simple text message, and the server returns a response body.

## Design

An HTTP client should have a `Client` type that sends requests, and a `Results` type that models the server response.

Things to consider when designing a client:
- Easy to use
- Hides internal complexity
- Consists of composable parts that users can bring together
- Synchronous by default
- Allow users to fine-tune API behavior

## Client type

[httbin.org](http://httpbin.org/) test server.

### Client Connections

TCP connections are expensive, so the HTTP protocol has a caching mechanism called _keep-alive_ that keeps established client/server connections open until a timeout. Then, a client can use the same connection to send HTTP requests without establishing a new connection.

Go's `DefaultClient` keeps 100 connections open and only allows you to reuse 2. You can optimize the connection pool with the Go [`Transport`](https://pkg.go.dev/net/http#Transport) type.


### Responses

You read a response body incrementally, as a stream of bytes. Create a `bytes.Buffer` and read the stream little by little until you read the entire body.

> Use an `io.Reader` to read any resource, and use an `io.Writer` to write to any resource. You can also use [`io.Copy(w, r)`](https://pkg.go.dev/io#Copy) that writes directly to a writer from a reader. Use the [`Discard`](https://pkg.go.dev/io#Discard) variable (of type `Writer`) to discard anything after you read it. You can treat `Discard` as `/dev/null`. This method preserves resources.

# Testing HTTP Clients 

Create handlers depending on what you want to test. For example, if you want to test successful HTTP requests:
1. Create a handler that responds with HTTP status code 200
2. Launch a test server with `httptest.NewServer(handler)`. For testing criteria:
   - Request: input
   - ResponseWriter: output

The test server requires a handler to handle requests and responses. The `Handler` interface has the following signature:

```go 
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
`Request` is the input, and `ResponseWriter` handles the output. After you create a handler that satisfies this interface, you pass it to the `NewServer` function to launch the test server. `NewServer` returns a `Server` server that contains a `URL` value to send requests.

Instead of writing an entirely new type to satisfy the `Handler` interface, Go provides the [`HandlerFunc` type](https://pkg.go.dev/net/http#HandlerFunc--a function that has a `ServeHTTP` method. This means that you can create a function that performs some action with a `Request` and `ResponseWriter`, then you can pass it to `HandlerFunc` to start a test server.

So, Go provides the `HandlerFunc` type:
```go 
type HandlerFunc func(ResponseWriter, *Request)
 
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    // Forwards the call to the converted function.
    f(w, r)
}
```
1. Create a function with the same signature as HandlerFunc:
   ```go
	handler := func(w http.ResponseWriter, r *http.Request) {
		// logic
	}
   ```
2. Convert that function to a `Handler` type with the `HandlerFunc`. `HandlerFunc` is an adapter that creates HTTP handlers from ordinary functions:
   ```go
	httpHandler := http.HandlerFunc(handler)
   ```
3. Pass the new handler to the `http.HandlerFunc(func)` method.
   ```go
	server := httptest.NewServer(httpHandler)
   ```
HTTP handlers in GO are concurrent, so the test server handles each request in its own goroutine. When you are tracking the number of requests, you should use the `atomic` package's concurrency-safe counters.


## httptest server

[httptest](https://pkg.go.dev/net/http/httptest)

Prefer the `.Cleanup(server.Close)` function over `defer` functions when testing. The `Cleanup` function runs after the tests complete rather than the enclosing function. This is useful with test helpers.





# Concurrency

## Channels

Using directional channels (send- or receive-only) prevents bugs and makes code easier to understand.

## Pipelines

## Orchestrator function pattern

This is when you create 2 functions:
- Logic function: contains the core logic
- Orchestrator: runs the producer's main logic in a goroutine

This means that the orchestrator can have channel ownership--create, write to, and close the channel--while the logic function is placed in a goroutine.

## time.Ticker

The [time.Ticker](https://pkg.go.dev/time#Ticker) type holds a channel that delivers ticks at intervals. You create one with the `NewTicker(d Duration)` function. It sends the current time to the channel, and the period of time between each send is the `d` value.

For example, the following function creates a ticker that controls how often a request is sent to the `out` channel:

```go
func Throttle(in <-chan *http.Request, out chan<- *http.Request, delay time.Duration) {
    t := time.NewTicker(delay)
    defer t.Stop()
 
    for r := range in {
        <-t.C
        out <- r
    }
}
```
The `<-t.C` line blocks the `for range` loop until the ticker sends a value every _n_ seconds. When the `t.C` channel unblocks, then the function can send a request from the `in` channel to the `out` channel.

# Context

A function that accepts a context should always check if the context is canceled.

`Background` creates a non-cancelable context. Then, you can derive a cancelable timeout context using the `WithTimeout` function.

The signal package has a NotifyContext function to catch an OS signal.