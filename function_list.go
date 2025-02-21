package parser

import (
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

var decimalsPlace int32 = 6

type ArgType int

const (
	ArgTypeInt ArgType = iota
	ArgTypeFloat
	ArgTypeString
)

func SetDecimalsPlace(place int32) {
	decimalsPlace = place
}

func isContainDot(s string) bool {
	return strings.Contains(s, ".")
}

func getArgType(a any) ArgType {
	switch a.(type) {
	case int, int64, int32, int16, int8:
		return ArgTypeInt
	case float64, float32:
		return ArgTypeFloat
	default:
		return ArgTypeString
	}
}

func getComputeType(a, b any) ArgType {
	aType := getArgType(a)
	bType := getArgType(b)

	if aType == ArgTypeFloat || bType == ArgTypeFloat {
		return ArgTypeFloat
	}

	if aType == ArgTypeInt && bType == ArgTypeInt {
		return ArgTypeInt
	}

	if isContainDot(cast.ToString(a)) || isContainDot(cast.ToString(b)) {
		return ArgTypeFloat
	}

	return ArgTypeInt
}

// 定义函数映射表
var funcMap = map[string]Function{
	"append": func(args ...any) any {
		if len(args) == 0 {
			return ""
		}

		if len(args) == 1 {
			return args[0]
		}

		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return a + b
	},
	"trim": func(args ...any) any {
		if len(args) == 0 {
			return ""
		}

		if len(args) == 1 {
			return args[0]
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return strings.Trim(a, b)
	},
	"trimInt": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return cast.ToInt64(args[0])
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return cast.ToInt64(strings.Trim(a, b))
	},
	"add": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return args[0]
		}

		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			dec1 := decimal.NewFromFloat(cast.ToFloat64(args[0]))
			dec2 := decimal.NewFromFloat(cast.ToFloat64(args[1]))
			return dec1.Add(dec2).Round(decimalsPlace).InexactFloat64()
		}

		return cast.ToInt64(args[0]) + cast.ToInt64(args[1])
	},
	"sub": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return args[0]
		}

		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			dec1 := decimal.NewFromFloat(cast.ToFloat64(args[0]))
			dec2 := decimal.NewFromFloat(cast.ToFloat64(args[1]))
			return dec1.Sub(dec2).Round(decimalsPlace).InexactFloat64()
		}

		return cast.ToInt64(args[0]) - cast.ToInt64(args[1])
	},

	"multi": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return args[0]
		}

		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			dec1 := decimal.NewFromFloat(cast.ToFloat64(args[0]))
			dec2 := decimal.NewFromFloat(cast.ToFloat64(args[1]))
			return dec1.Mul(dec2).Round(decimalsPlace).InexactFloat64()
		}

		return cast.ToInt64(args[0]) * cast.ToInt64(args[1])
	},
	"div": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return args[0]
		}

		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			dec1 := decimal.NewFromFloat(cast.ToFloat64(args[0]))
			dec2 := decimal.NewFromFloat(cast.ToFloat64(args[1]))
			return dec1.Div(dec2).Round(decimalsPlace).InexactFloat64()
		}

		return cast.ToInt64(args[0]) / cast.ToInt64(args[1])
	},
	"mod": func(args ...any) any {
		if len(args) == 0 {
			return 0
		}

		if len(args) == 1 {
			return args[0]
		}
		m := cast.ToInt64(args[1])
		if m == 0 {
			return 0
		}
		return cast.ToInt64(args[0]) % cast.ToInt64(args[1])
	},
	"eq": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		return cast.ToString(args[0]) == cast.ToString(args[1])
	},
	"ne": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		return cast.ToString(args[0]) != cast.ToString(args[1])
	},
	"gt": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			return cast.ToFloat64(args[0]) > cast.ToFloat64(args[1])
		}
		return cast.ToInt64(args[0]) > cast.ToInt64(args[1])
	},
	"gte": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			return cast.ToFloat64(args[0]) >= cast.ToFloat64(args[1])
		}
		return cast.ToInt64(args[0]) >= cast.ToInt64(args[1])
	},
	"lt": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			return cast.ToFloat64(args[0]) < cast.ToFloat64(args[1])
		}
		return cast.ToInt64(args[0]) < cast.ToInt64(args[1])
	},
	"lte": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		computeType := getComputeType(args[0], args[1])
		if computeType == ArgTypeFloat {
			return cast.ToFloat64(args[0]) <= cast.ToFloat64(args[1])
		}
		return cast.ToInt64(args[0]) <= cast.ToInt64(args[1])
	},
	"hasPrefix": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return strings.HasPrefix(a, b)
	},
	"hasSuffix": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return strings.HasSuffix(a, b)
	},
	"contains": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return strings.Contains(a, b)
	},
	"regexp": func(args ...any) any {
		if len(args) < 2 {
			return false
		}
		a := cast.ToString(args[0])
		b := cast.ToString(args[1])
		return regexp.MustCompile(b).MatchString(a)
	},
	"not": func(args ...any) any {
		if len(args) == 0 {
			return false
		}
		return !cast.ToBool(args[0])
	},
	"and": func(args ...any) any {
		if len(args) == 0 {
			return false
		}
		return cast.ToBool(args[0]) && cast.ToBool(args[1])
	},
	"or": func(args ...any) any {
		if len(args) == 0 {
			return false
		}
		return cast.ToBool(args[0]) || cast.ToBool(args[1])
	},
}
