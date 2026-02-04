package main

import (
	"fmt"
	"log"

	"github.com/blackshieldpt/gorpn"
)

func main() {
	// Basic arithmetic
	result, err := gorpn.Eval("2 + 3 * 4", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("2 + 3 * 4 = %g\n", result)

	// With parentheses
	result, _ = gorpn.Eval("(2 + 3) * 4", nil)
	fmt.Printf("(2 + 3) * 4 = %g\n", result)

	// Using variables
	vars := map[string]float64{
		"x": 10,
		"y": 5,
	}
	result, _ = gorpn.Eval("x * y + 2", vars)
	fmt.Printf("x * y + 2 (x=10, y=5) = %g\n", result)

	// Using built-in functions
	result, _ = gorpn.Eval("sqrt(16) + pow(2, 3)", nil)
	fmt.Printf("sqrt(16) + pow(2, 3) = %g\n", result)

	// Using constants
	result, _ = gorpn.Eval("2 * pi", nil)
	fmt.Printf("2 * pi = %g\n", result)

	// Complex expression with variables
	vars = map[string]float64{"a": 3, "b": 4}
	result, _ = gorpn.Eval("sqrt(a^2 + b^2)", vars)
	fmt.Printf("sqrt(a^2 + b^2) (a=3, b=4) = %g\n", result)
}
