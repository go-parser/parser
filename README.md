# Go Parser

A high-performance Go expression parser that supports numeric calculations, string concatenation, function calls, conditional statements, function definitions, and registration.

## Features

- ✨ Numeric Operations: Supports +, -, *, /, and () operators
- 🔤 String Concatenation: Using + operator
- 📝 Variable Substitution: Variables with $ prefix
- 🎯 Function Calls: Functions with @ prefix
- ⚡ Conditional Statements: Supports >, <, >=, <=, ==, !=, &&, ||, ! operators
- 🚀 High Performance: Excellent concurrent execution performance

## Usage Guide

### 1. Basic Syntax

#### Numeric Calculations
Numeric expressions must be wrapped in parentheses:
```shell
(1+2)*3  # Result: 9
```

#### String Concatenation
Use + operator to concatenate strings, strings must be wrapped in double quotes:
```shell
"hello"+"world"  # Result: helloworld
```

#### Variables
Variables must start with $, underscores are allowed but other special characters are not:
```shell
$price+100  # price is a variable, value is passed during execution
```

#### Function Calls
Function calls must start with @:
```shell
@calculate($price, 100)  # Calling calculate function
```

#### Conditional Expressions
Supports complex conditional logic:
```shell
($price >= 100 && $price < 200) || @isVIP($userId)
```

### 2. Code Examples

#### Function Registration
```go
// Function definition
type Function func(args ...any) any

// Register function
RegisterFunc(name string, f Function)
```

#### Expression Execution
```go
// Simple calculation
expr, err := ParseExpression("($price+100)*0.8")
result := expr.Execute(map[string]any{"price": 200})

// Conditional evaluation
expr, err := ParseExpression(`$price > 100 && @isVIP($userId)`)
result := expr.Execute(map[string]any{
    "price": 150,
    "userId": "user123",
})
```

## Performance Benchmarks

```
go test -bench=. -benchmem -tags -v
goos: darwin
goarch: arm64
pkg: github.com/go-parser/parser
cpu: Apple M1 Pro
BenchmarkHardCode-8                  	 6342241	       174.6 ns/op	     112 B/op	       3 allocs/op
BenchmarkHardCodeParallel-8          	17822409	        68.70 ns/op	     112 B/op	       3 allocs/op
BenchmarkGoParserExecute-8           	 6851533	       174.5 ns/op	     112 B/op	       4 allocs/op
BenchmarkGoParserExecuteParallel-8   	20086538	        63.08 ns/op	     112 B/op	       4 allocs/op
BenchmarkExprLangExecute-8           	 2112973	       558.4 ns/op	     368 B/op	      14 allocs/op
BenchmarkExprLangExecuteParallel-8   	 5690395	       234.4 ns/op	     368 B/op	      14 allocs/op
PASS
ok  	github.com/go-parser/parser	8.986s
```
