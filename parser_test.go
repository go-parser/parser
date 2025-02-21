package parser

import (
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/spf13/cast"
)

var vars = map[string]any{"stock": 2}

func newFunc() *FunctionCall {
	f, _ := ParseExpression(`@trimInt($stock,"stock:")*100+5`)
	return f
}

func TestRunConditions(t *testing.T) {
	type args struct {
		expression string
		vars       map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test2",
			args: args{
				expression: `$stock>100`,
				vars:       map[string]any{"stock": 120},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ParseExpression(tt.args.expression)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			got := f.Execute(tt.args.vars)
			if got != tt.want {
				t.Errorf("RunConditions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpression_Execute(t *testing.T) {
	type fields struct {
		If        string
		Then      string
		Otherwise string
	}
	type args struct {
		vars map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		{
			name: "test1",
			fields: fields{
				If:        `$MIN_ORD_QTY>0`, // MIN_ORD_QTY >= 0
				Then:      `$MIN_ORD_QTY`,   // MIN_ORD_QTY
				Otherwise: "0",
			},
			args: args{
				vars: map[string]any{"MIN_ORD_QTY": 120},
			},
			want: 120,
		},
		{
			name: "test2",
			fields: fields{
				If:        `($MIN_ORD_QTY+1)>=10`, // MIN_ORD_QTY+1 >= 10
				Then:      `$MIN_ORD_QTY`,         // MIN_ORD_QTY
				Otherwise: "0",
			},
			args: args{
				vars: map[string]any{"MIN_ORD_QTY": 5},
			},
			want: 0,
		},
		{
			name: "test3",
			fields: fields{
				If:        `$stock>100`, // stock > 100
				Then:      `$stock`,     // stock
				Otherwise: "0",
			},
			args: args{
				vars: map[string]any{"stock": "200"},
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expression{
				If:        tt.fields.If,
				Then:      tt.fields.Then,
				Otherwise: tt.fields.Otherwise,
			}
			e.Parse()
			if got, _ := e.Eval(tt.args.vars); cast.ToInt(got) != cast.ToInt(tt.want) {
				t.Errorf("Expression.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	type args struct {
		expression string
		vars       map[string]any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test1",
			args: args{
				expression: `@trimInt($stock,"stock:")*100+5`,
				vars:       vars,
			},
			want: int64(205),
		},
		{
			name: "test2",
			args: args{
				expression: `@trimInt($stock,"stock:") % 100 + 5`,
				vars:       vars,
			},
			want: int64(7),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := ParseExpression(tt.args.expression)
			if got := executeFunctionCall(f, tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testTrim(args ...any) any {
	a := cast.ToString(args[0])
	b := cast.ToString(args[1])
	return strings.Trim(a, b)
}

func testMulti(args ...any) any {
	a := cast.ToInt(args[0])
	b := cast.ToInt(args[1])
	return a * b
}

func testAdd(args ...any) any {
	a := cast.ToInt(args[0])
	b := cast.ToInt(args[1])
	return a + b
}

func hardCodeFunc(vars map[string]any) any {
	stock := vars["$stock"]

	stock = testTrim(stock, "stock:")

	stock = testMulti(stock, 100)
	stock = testAdd(stock, 5)

	return stock
}

func BenchmarkHardCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hardCodeFunc(vars)
	}
}

func BenchmarkHardCodeParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hardCodeFunc(vars)
		}
	})
}

func BenchmarkGoParserExecute(b *testing.B) {
	f := newFunc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeFunctionCall(f, vars)
	}
	b.StopTimer()
}

func BenchmarkGoParserExecuteParallel(b *testing.B) {
	f := newFunc()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			executeFunctionCall(f, vars)
		}
	})
	b.StopTimer()
}

func BenchmarkExprLangExecute(b *testing.B) {
	code := `trim2Int($stock,"stock:")*100+5`
	env := map[string]any{
		"$stock": "stock:200",
		"trim2Int": func(args ...any) (any, error) {
			if len(args) == 0 {
				return 0, nil
			}

			if len(args) == 1 {
				return cast.ToInt64(args[0]), nil
			}
			a := cast.ToString(args[0])
			b := cast.ToString(args[1])
			return cast.ToInt64(strings.Trim(a, b)), nil
		},
	}

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		b.Errorf("Compile() error = %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr.Run(program, env)
	}
	b.StopTimer()
}

func BenchmarkExprLangExecuteParallel(b *testing.B) {
	code := `trim2Int($stock,"stock:")*100+5`
	env := map[string]any{
		"$stock": "stock:200",
		"trim2Int": func(args ...any) (any, error) {
			if len(args) == 0 {
				return 0, nil
			}

			if len(args) == 1 {
				return cast.ToInt64(args[0]), nil
			}
			a := cast.ToString(args[0])
			b := cast.ToString(args[1])
			return cast.ToInt64(strings.Trim(a, b)), nil
		},
	}

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		b.Errorf("Compile() error = %v", err)
		return
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expr.Run(program, env)
		}
	})
	b.StopTimer()
}

func TestMain(m *testing.M) {
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()
	m.Run()
}
