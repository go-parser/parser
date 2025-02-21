package parser

import (
	"errors"
	"fmt"

	"github.com/go-parser/parser/internal/parser"
	"github.com/spf13/cast"
)

type Action struct {
	Expression string `json:"expression"`
	execute    *FunctionCall
}

type Expression struct {
	If              string  `json:"if,omitempty"`
	Then            string  `json:"then,omitempty"`
	Otherwise       string  `json:"otherwise,omitempty"`
	ifAction        *Action `json:"-"`
	thenAction      *Action `json:"-"`
	otherwiseAction *Action `json:"-"`
}

func (e *Expression) String() string {
	s := ""
	if e.If != "" {
		s += fmt.Sprintf("if: %s [%s], ", e.If, e.ifAction.execute.Expression)
	}
	if e.Then != "" {
		s += fmt.Sprintf("then: %s [%s], ", e.Then, e.thenAction.execute.Expression)
	}
	if e.Otherwise != "" {
		s += fmt.Sprintf("otherwise: %s [%s]", e.Otherwise, e.otherwiseAction.execute.Expression)
	}
	return s
}

func ParseExpression(expr string) (*FunctionCall, error) {
	expr, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}

	return ParseFunctionExpression(expr)
}

func ParseAndExecute(expr string, vars map[string]any) (any, error) {
	f, err := ParseExpression(expr)
	if err != nil {
		return nil, err
	}
	return f.Execute(vars), nil
}

func (e *Expression) Parse() error {
	if e.Then == "" {
		return fmt.Errorf("then is required")
	}
	// Parse condition
	if e.If != "" {
		expr, err := ParseExpression(e.If)
		if err != nil {
			return err
		}

		e.ifAction = &Action{
			Expression: e.If,
			execute:    expr,
		}
	}
	// Parse Then
	if e.Then != "" {
		execute, err := ParseExpression(e.Then)
		if err != nil {
			return err
		}

		e.thenAction = &Action{
			Expression: e.Then,
			execute:    execute,
		}
	}
	// Parse Otherwise
	if e.Otherwise != "" {
		execute, err := ParseExpression(e.Otherwise)
		if err != nil {
			return err
		}

		e.otherwiseAction = &Action{
			Expression: e.Otherwise,
			execute:    execute,
		}
	}
	return nil
}

func (e *Expression) Eval(vars map[string]any) (any, error) {
	// Execute condition
	condition := true
	if e.ifAction != nil {
		conditionRes := e.ifAction.execute.Execute(vars)
		if !cast.ToBool(conditionRes) {
			condition = false
		}
	}

	// Execute Then
	if condition && e.Then != "" && e.thenAction != nil {
		return executeFunctionCall(e.thenAction.execute, vars), nil
	}
	// Execute Otherwise
	if !condition && e.Otherwise != "" && e.otherwiseAction != nil {
		return executeFunctionCall(e.otherwiseAction.execute, vars), nil
	}
	return nil, errors.New("invalid expression")
}
