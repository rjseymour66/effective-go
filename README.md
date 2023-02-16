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

# Packages

A package name should describe what it provides, not what it does.


# Formatting verbs

| verb | Definition | Usage |
|------|:-----------|:------|
| %q   | Wraps the given string in double quotes.

# Misc

## Empty interface

Go versions prior to 1.18 used the empty interface: `interface{}`. This is an interface that did not implement any methods, so any type satisfied it. In Go 1.18 and later, `interface{}` was replaced with `any`. 

## Shortif declaration

`if variable := value; condition`

For example:
```go
if err := json.Marshal(&val); err != nil {
    // handle error
}
```