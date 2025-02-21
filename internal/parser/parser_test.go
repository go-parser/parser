package parser

import "testing"

func TestConvertExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				expr: "1+2",
			},
			want: "add(1:int,2:int)",
		},
		{
			name: "test2",
			args: args{
				expr: "1+2*3",
			},
			want: "add(1:int,multi(2:int,3:int))",
		},
		{
			name: "test3",
			args: args{
				expr: "1+2*3-4",
			},
			want: "sub(add(1:int,multi(2:int,3:int)),4:int)",
		},
		{
			name: "test4",
			args: args{
				expr: "1+2*3-4/5",
			},
			want: "sub(add(1:int,multi(2:int,3:int)),div(4:int,5:int))",
		},
		{
			name: "test5",
			args: args{
				expr: "1+2*3-4/5+@sum(1,2,3)",
			},
			want: "add(sub(add(1:int,multi(2:int,3:int)),div(4:int,5:int)),sum(1:int,2:int,3:int))",
		},
		{
			name: "test6",
			args: args{
				expr: "1+2*@funA($a,2)-3",
			},
			want: "sub(add(1:int,multi(2:int,funA($a,2:int))),3:int)",
		},
		{
			name: "test7",
			args: args{
				expr: "1+2*@funA($a,2)-@funB(3,4)",
			},
			want: "sub(add(1:int,multi(2:int,funA($a,2:int))),funB(3:int,4:int))",
		},
		{
			name: "test8",
			args: args{
				expr: "1+2*@funA($a)-@funB(3,4)+@funC(5,6)",
			},
			want: "add(sub(add(1:int,multi(2:int,funA($a))),funB(3:int,4:int)),funC(5:int,6:int))",
		},
		{
			name: "test9",
			args: args{
				expr: `$a+"s"`,
			},
			want: "add($a,s:str)",
		},
		{
			name: "test10",
			args: args{
				expr: `$a+"s"+$b`,
			},
			want: "add(add($a,s:str),$b)",
		},
		{
			name: "test11",
			args: args{
				expr: `$a+"s"+$b+"t"`,
			},
			want: "add(add($a,s:str),add($b,t:str))",
		},

		{
			name: "test12",
			args: args{
				expr: `$a+"s"+$b+"t"+$c`,
			},
			want: "add(add(add($a,s:str),add($b,t:str)),$c)",
		},
		{
			name: "test13",
			args: args{
				expr: `@funA($a+1,2)`,
			},
			want: "funA(add($a,1:int),2:int)",
		},
		{
			name: "test14",
			args: args{
				expr: "1+2*@funA($a+1,$b)",
			},
			want: "add(1:int,multi(2:int,funA(add($a,1:int),$b)))",
		},
		{
			name: "test15",
			args: args{
				expr: "@funA($a+1,$b)",
			},
			want: "funA(add($a,1:int),$b)",
		},
		{
			name: "test16",
			args: args{
				expr: "1*(2+3)",
			},
			want: "multi(1:int,add(2:int,3:int))",
		},
		{
			name: "test17",
			args: args{
				expr: "$a > 1 && $b < 2",
			},
			want: "and(gt($a,1:int),lt($b,2:int))",
		},
		{
			name: "test18",
			args: args{
				expr: "($a > 1 && $b < 2) || $c == 3",
			},
			want: "or(and(gt($a,1:int),lt($b,2:int)),eq($c,3:int))",
		},
		{
			name: "test19",
			args: args{
				expr: "!($a > 1 && $b < 2)",
			},
			want: "not(and(gt($a,1:int),lt($b,2:int)))",
		},
		{
			name: "test20",
			args: args{
				expr: "($a+1)>=10",
			},
			want: "gte(add($a,1:int),10:int)",
		},
		{
			name: "test21",
			args: args{
				expr: `($stock>100 && $stock<200) && $mfr=="motorola"`,
			},
			want: "and(and(gt($stock,100:int),lt($stock,200:int)),eq($mfr,motorola:str))",
		},
		{
			name: "test22",
			args: args{
				expr: `($stock>100 && $stock<200) || $mfr=="motorola"`,
			},
			want: "or(and(gt($stock,100:int),lt($stock,200:int)),eq($mfr,motorola:str))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.expr)
			if err != nil {
				t.Errorf("ConvertExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ConvertExpression() = %s", tt.args.expr)
				t.Errorf("ConvertExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
