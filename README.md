# Go Parser

## Introduction

---

    Go Parser is a simple expression parser that supports numeric calculations, string concatenation, function calls, conditional statements, function definitions, function registration, complex calculations, and variable substitution.

#### Definitions

##### 1. Numeric Calculations

Supports +, -, *, /, and (). Numeric expressions must be wrapped in parentheses.

```shell
(1+2)*3 # 9
```

##### 2. String Concatenation

+ operator supports string concatenation. Strings must be wrapped in double quotes.

```shell
"a"+"b" # ab
```

##### 3. Variable Input

Supports variable input. Variables must be prefixed with $. Variables cannot contain special characters except underscores.

```shell
$stock+1 # stock is a variable, value is passed during execution
```

##### 4. Functions

Supports function calls, function definitions, and function registration. Function calls must be prefixed with @.

```shell
@funA($a+1,2) # calling function funA
```

##### 5. Conditional Statements

Supports `>,<,>=,<=,==,!=,&&,||,!,(,)` operators

```shell
($a >= 100 && $a < 200) || @funA($b) == "a123"
```

## Examples

---

#### 1. Writing Expressions

```shell
($a+1)*10 # numeric calculation, $a is a variable

100+($a+1)*10 # numeric calculation, $a is a variable

$a+"s"+$b+"t"+$c # string concatenation

@funA($a+1,2) # calling function funA

@funA($a+1,2) > 12 # calling function funA and comparing with 12
```

For easy identification and to avoid ambiguity, variables must be prefixed with $ and function calls must be prefixed with @

#### 2. Function Definition and Registration

```shell
type Function func(args ...any) any // function definition

RegisterFunc(name string, f Function) // register function
```

#### 3. Expression Execution

Simple configuration, numeric calculations supporting +, -, *, /, (), string concatenation

```shell
($a+1)*10
```

```go
f , err := ParseExpression("($a+1)*10")

result := f.Execute(map[string]any{"a": 10})
```

#### 4. Conditional Statement Execution

```shell
$a > 100 && $a < 200 || @funA($b) == "a123"
```

```go
f, err := ParseExpression(`$a > 100 && $a < 200 || @funA($b) == "a123"`)

result := f.Execute(map[string]any{"a": 120, "b": "a123"})
```

## Performance

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
