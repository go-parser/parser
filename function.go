package parser

import (
	"errors"
	"strconv"
	"strings"
)

const (
	varPrefix = '$'
)

// Define function type
type Function func(args ...any) any

func RegisterFunc(name string, f Function) {
	funcMap[name] = f
}

type FunctionArg struct {
	FunctionCall *FunctionCall
	Variable     string
	Const        any
}

// Define function call type
type FunctionCall struct {
	Expression   string
	Function     Function
	FunctionName string
	Args         []*FunctionArg // Arguments can be another function call or a constant/variable
	Variable     string
	Const        any
}

func (f *FunctionCall) Execute(vars map[string]any) any {
	return executeFunctionCall(f, vars)
}

// ParseFunctionExpression parses expression, function call format is funcName(arg1,arg2,...)
func ParseFunctionExpression(expr string) (*FunctionCall, error) {
	// Find function name and parameter list
	left := strings.Index(expr, "(")
	right := strings.LastIndex(expr, ")")
	if left == -1 || right == -1 {
		// If no brackets found, it's a variable or constant
		fVar := parseVar(expr)
		return &FunctionCall{
			Expression: expr,
			Function:   nil,
			Const:      fVar.Const,
			Variable:   fVar.Variable,
		}, nil
	}
	funcName := strings.TrimSpace(expr[:left])
	paramList := expr[left+1 : right]

	if _, ok := funcMap[funcName]; !ok {
		return nil, errors.New("function not found: " + funcName)
	}

	// Create function call
	call := FunctionCall{
		Expression:   expr,
		Function:     funcMap[funcName],
		FunctionName: funcName,
	}

	// Parse parameter list
	level := 0
	start := 0
	for i, ch := range paramList {
		switch ch {
		case '(':
			level++
		case ')':
			level--
		case ',':
			if level == 0 {
				arg, err := parseArgument(strings.TrimSpace(paramList[start:i]))
				if err != nil {
					return nil, err
				}
				call.Args = append(call.Args, arg)
				start = i + 1
			}
		}
	}

	// Add last argument
	arg, err := parseArgument(strings.TrimSpace(paramList[start:]))
	if err != nil {
		return nil, err
	}

	call.Args = append(call.Args, arg)

	return &call, nil
}

// Parse argument
func parseArgument(arg string) (*FunctionArg, error) {
	if strings.Contains(arg, `(`) {
		f, err := ParseFunctionExpression(arg)
		if err != nil {
			return nil, err
		}
		return &FunctionArg{
			FunctionCall: f,
		}, nil
	}

	if arg == "" {
		return &FunctionArg{Const: arg}, nil
	}

	if arg[0] == varPrefix {
		return &FunctionArg{Variable: arg[1:]}, nil
	}

	index := strings.LastIndex(arg, ":")

	typ := arg[index+1:]
	if index == -1 {
		typ = "int"
	} else {
		arg = arg[:index]
	}

	switch typ {
	case "int":
		val, _ := strconv.ParseInt(arg, 10, 64)
		return &FunctionArg{Const: val}, nil
	case "float":
		val, _ := strconv.ParseFloat(arg, 64)
		return &FunctionArg{Const: val}, nil
	case "str":
		return &FunctionArg{Const: arg}, nil
	default:
		return nil, errors.New("invalid argument: " + arg)
	}
}

// Parse variable
func parseVar(arg string) *FunctionArg {
	if arg[0] == varPrefix {
		return &FunctionArg{Variable: arg[1:]}
	}

	if arg == "" {
		return &FunctionArg{Const: arg}
	}

	index := strings.LastIndex(arg, ":")

	typ := arg[index+1:]
	if index == -1 {
		typ = "int"
	} else {
		arg = arg[:index]
	}

	switch typ {
	case "int":
		val, _ := strconv.ParseInt(arg, 10, 64)
		return &FunctionArg{Const: val}
	case "float":
		val, _ := strconv.ParseFloat(arg, 64)
		return &FunctionArg{Const: val}
	default:
		return &FunctionArg{Const: arg}
	}
}

// Execute function call
func executeFunctionCall(call *FunctionCall, vars map[string]any) any {
	if call.Function == nil {
		// If function is nil, it's a variable or constant
		if call.Variable != "" {
			if val, ok := vars[call.Variable]; ok {
				return val
			}
			return call.Variable
		}

		return call.Const
	}

	// Parse and execute all arguments
	args := make([]any, len(call.Args))
	for i, arg := range call.Args {
		// If argument is a function call, execute it
		if arg.FunctionCall != nil {
			args[i] = executeFunctionCall(arg.FunctionCall, vars)
			continue
		}

		// If argument is a variable, replace it
		if arg.Variable != "" {
			if val, ok := vars[arg.Variable]; ok {
				args[i] = val
			} else {
				args[i] = arg.Variable
			}
			continue
		}

		args[i] = arg.Const
	}

	// Execute function
	return call.Function(args...)
}
