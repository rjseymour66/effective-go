# Effective Go

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

# Misc

## Docs 

You can generate docs that include your [testable examples](#testable-examples) with `godoc`. The following command installs the latest version:

```shell
$ go install golang.org/x/tools/cmd/godoc@latest
```



## Short-if declaration

`if variable := value; condition`

For example:
```go
if err := json.Marshal(&val); err != nil {
    // handle error
}
```
## Testable examples

A _testable example_ is live documentation for code. You write a testable example to demonstrate the package API to other developers. The API includes the exported identifiers, such as functions, methods, etc. A testable example never goes out of date.

The testing package runs testable examples and checks their results, but it does not report successes or failures.

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