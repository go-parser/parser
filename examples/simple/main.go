package main

import (
	"fmt"
	"time"

	"github.com/go-parser/parser"
)

func main() {
	fmt.Println()
	{
		fmt.Println("===============test-0 test end===============")
		expr := parser.Expression{
			If:   `$stock>100`,   // stock>100
			Then: `$stock*100+5`, // stock*100+5
		}

		expr.Parse()
		vars := map[string]any{
			"stock": 120,
		}

		fmt.Println("test-0:", expr.String())
		now := time.Now()
		result, _ := expr.Eval(vars)
		fmt.Println("test-0 res:", result, " , cost:", time.Since(now).Nanoseconds())
		fmt.Println("===============test-0 test end===============")
	}

	fmt.Println()
	{
		fmt.Println("===============test-1 test start===============")
		// Conditional expression (structured)
		expr := parser.Expression{
			If:   `@trim($stock,"stock:")>100`,   // (stock-"stock:")>100
			Then: `@trim($stock,"stock:")*100+5`, // (stock-100)*100+5
		}

		expr.Parse()
		vars := map[string]any{
			"stock": "stock:120",
		}

		fmt.Println("test-1:", expr.String())
		now := time.Now()
		result, _ := expr.Eval(vars)
		fmt.Println("test-1 res:", result, " , cost:", time.Since(now).Nanoseconds())
		fmt.Println("===============test-1 test end===============")
	}

	fmt.Println()
	{
		fmt.Println("===============test-2 test start===============")
		// Conditional expression (structured)
		expr := parser.Expression{
			If:   `@hasPrefix($stock,"stock:")`,  // Check if stock starts with "stock:"
			Then: `@trim($stock,"stock:")*100+5`, // (stock-100)*100+5
		}

		expr.Parse()
		vars := map[string]any{
			"stock": "stock:120",
		}

		fmt.Println("test-2:", expr.String())
		now := time.Now()
		result, _ := expr.Eval(vars)
		fmt.Println("test-2 res:", result, " , cost:", time.Since(now).Nanoseconds())
		fmt.Println("===============test-2 test end===============")
	}

	fmt.Println()
	{
		fmt.Println("===============test-3 test start===============")
		// Conditional expression (structured)
		expr := parser.Expression{
			If:   `($stock>100 && $stock<200) && $mfr=="motorola"`, // (stock>100 && stock<200) && mfr=="motorola"
			Then: `$stock*100+5`,                                   // stock*100+5
		}

		expr.Parse()
		vars := map[string]any{
			"stock": 120,
			"mfr":   "motorola",
		}

		fmt.Println("test-3:", expr.String())
		now := time.Now()
		result, _ := expr.Eval(vars)
		fmt.Println("test-3 res:", result, " , cost:", time.Since(now).Nanoseconds())
		fmt.Println("===============test-3 test end===============")
	}

	fmt.Println()
	{
		fmt.Println("===============test-4 test start===============")
		// Conditional expression (structured)
		expr := parser.Expression{
			If:        `$stock>0`, // stock>0
			Then:      `$stock`,   // stock
			Otherwise: "0",
		}

		expr.Parse()
		vars := map[string]any{
			"stock": "a",
		}

		fmt.Println("test-4:", expr.String())
		now := time.Now()
		result, _ := expr.Eval(vars)
		fmt.Println("test-4 res:", result, " , cost:", time.Since(now).Nanoseconds())
		fmt.Println("===============test-3 test end===============")
	}
}
