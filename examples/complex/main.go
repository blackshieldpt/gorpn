package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/blackshieldpt/gorpn"
)

func main() {
	// Create a github.com/blackshieldpt/gorpn instance
	calc := gorpn.New()

	// Register custom functions
	registerCustomFunctions(calc)

	// Set some default variables
	calc.SetVariables(map[string]float64{
		"gravity": 9.81,
		"c":       299792458, // speed of light m/s
	})

	fmt.Println("=== Calculator Demo ===")
	fmt.Println()

	// Demo 1: Basic operations with precedence
	demo1(calc)

	// Demo 2: Variables and constants
	demo2(calc)

	// Demo 3: Built-in functions
	demo3(calc)

	// Demo 4: Custom functions
	demo4(calc)

	// Demo 5: Variadic functions
	demo5(calc)

	// Demo 6: Real-world formulas
	demo6(calc)

	// Demo 7: Interactive REPL
	fmt.Println("\n=== Interactive Mode ===")
	fmt.Println("Enter expressions (type 'quit' to exit, 'vars' to see variables)")
	fmt.Println("Set variables: x = 5")
	fmt.Println()
	repl(calc)
}

func registerCustomFunctions(calc *gorpn.Calculator) {
	// Fahrenheit to Celsius
	calc.RegisterFunction("f2c", 1, func(args []float64) (float64, error) {
		return (args[0] - 32) * 5 / 9, nil
	})

	// Celsius to Fahrenheit
	calc.RegisterFunction("c2f", 1, func(args []float64) (float64, error) {
		return args[0]*9/5 + 32, nil
	})

	// Degrees to radians
	calc.RegisterFunction("deg", 1, func(args []float64) (float64, error) {
		return args[0] * math.Pi / 180, nil
	})

	// Radians to degrees
	calc.RegisterFunction("rad", 1, func(args []float64) (float64, error) {
		return args[0] * 180 / math.Pi, nil
	})

	// Clamp value between min and max
	calc.RegisterFunction("clamp", 3, func(args []float64) (float64, error) {
		value, minVal, maxVal := args[0], args[1], args[2]
		if value < minVal {
			return minVal, nil
		}
		if value > maxVal {
			return maxVal, nil
		}
		return value, nil
	})

	// Linear interpolation
	calc.RegisterFunction("lerp", 3, func(args []float64) (float64, error) {
		a, b, t := args[0], args[1], args[2]
		return a + (b-a)*t, nil
	})

	// Map value from one range to another
	calc.RegisterFunction("map", 5, func(args []float64) (float64, error) {
		value, inMin, inMax, outMin, outMax := args[0], args[1], args[2], args[3], args[4]
		return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin, nil
	})

	// Factorial (for small integers)
	calc.RegisterFunction("fact", 1, func(args []float64) (float64, error) {
		n := int(args[0])
		if n < 0 {
			return 0, fmt.Errorf("factorial of negative number")
		}
		if n > 20 {
			return 0, fmt.Errorf("factorial overflow for n > 20")
		}
		result := 1.0
		for i := 2; i <= n; i++ {
			result *= float64(i)
		}
		return result, nil
	})

	// Combination (n choose k)
	calc.RegisterFunction("comb", 2, func(args []float64) (float64, error) {
		n, k := int(args[0]), int(args[1])
		if k > n || k < 0 {
			return 0, nil
		}
		if k > n-k {
			k = n - k
		}
		result := 1.0
		for i := 0; i < k; i++ {
			result = result * float64(n-i) / float64(i+1)
		}
		return result, nil
	})

	// Distance between two 2D points
	calc.RegisterFunction("dist2d", 4, func(args []float64) (float64, error) {
		x1, y1, x2, y2 := args[0], args[1], args[2], args[3]
		return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)), nil
	})

	// Distance between two 3D points
	calc.RegisterFunction("dist3d", 6, func(args []float64) (float64, error) {
		x1, y1, z1, x2, y2, z2 := args[0], args[1], args[2], args[3], args[4], args[5]
		return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) + (z2-z1)*(z2-z1)), nil
	})

	// Quadratic formula (returns one root, use -b for other)
	calc.RegisterFunction("quadp", 3, func(args []float64) (float64, error) {
		a, b, c := args[0], args[1], args[2]
		discriminant := b*b - 4*a*c
		if discriminant < 0 {
			return 0, fmt.Errorf("no real roots")
		}
		return (-b + math.Sqrt(discriminant)) / (2 * a), nil
	})

	calc.RegisterFunction("quadm", 3, func(args []float64) (float64, error) {
		a, b, c := args[0], args[1], args[2]
		discriminant := b*b - 4*a*c
		if discriminant < 0 {
			return 0, fmt.Errorf("no real roots")
		}
		return (-b - math.Sqrt(discriminant)) / (2 * a), nil
	})
}

func demo1(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 1: Operator Precedence ---")
	expressions := []string{
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"2 ^ 3 ^ 2",
		"10 - 6 / 2 + 1",
		"-5 + 3 * -2",
	}
	for _, expr := range expressions {
		result, _ := calc.Evaluate(expr)
		fmt.Printf("  %s = %g\n", expr, result)
	}
	fmt.Println()
}

func demo2(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 2: Variables and Constants ---")

	// Built-in constants
	result, _ := calc.Evaluate("pi")
	fmt.Printf("  pi = %g\n", result)
	result, _ = calc.Evaluate("e")
	fmt.Printf("  e = %g\n", result)

	// User-defined variables
	result, _ = calc.Evaluate("gravity")
	fmt.Printf("  gravity = %g\n", result)
	result, _ = calc.Evaluate("c")
	fmt.Printf("  c (speed of light) = %g\n", result)

	// Expression with variables
	calc.SetVariable("mass", 70)
	result, _ = calc.Evaluate("mass * gravity")
	fmt.Printf("  mass * gravity (weight of 70kg) = %g N\n", result)

	// Inline variables
	result, _ = calc.EvaluateWithVars("0.5 * m * v^2", map[string]float64{"m": 10, "v": 5})
	fmt.Printf("  kinetic energy (m=10, v=5) = %g J\n", result)
	fmt.Println()
}

func demo3(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 3: Built-in Functions ---")
	expressions := []string{
		"sin(pi / 2)",
		"cos(0)",
		"sqrt(144)",
		"log10(1000)",
		"pow(2, 10)",
		"abs(-42)",
		"round(3.7)",
		"floor(3.7)",
		"ceil(3.2)",
		"exp(1)",
		"hypot(3, 4)",
	}
	for _, expr := range expressions {
		result, _ := calc.Evaluate(expr)
		fmt.Printf("  %s = %g\n", expr, result)
	}
	fmt.Println()
}

func demo4(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 4: Custom Functions ---")

	// Temperature conversions
	result, _ := calc.Evaluate("f2c(212)")
	fmt.Printf("  f2c(212) = %g°C (boiling point)\n", result)
	result, _ = calc.Evaluate("c2f(0)")
	fmt.Printf("  c2f(0) = %g°F (freezing point)\n", result)

	// Angle conversions
	result, _ = calc.Evaluate("sin(deg(90))")
	fmt.Printf("  sin(deg(90)) = %g\n", result)
	result, _ = calc.Evaluate("rad(pi)")
	fmt.Printf("  rad(pi) = %g degrees\n", result)

	// Clamp and lerp
	result, _ = calc.Evaluate("clamp(15, 0, 10)")
	fmt.Printf("  clamp(15, 0, 10) = %g\n", result)
	result, _ = calc.Evaluate("lerp(0, 100, 0.25)")
	fmt.Printf("  lerp(0, 100, 0.25) = %g\n", result)

	// Map value
	result, _ = calc.Evaluate("map(50, 0, 100, 0, 1)")
	fmt.Printf("  map(50, 0, 100, 0, 1) = %g\n", result)

	// Factorial and combinations
	result, _ = calc.Evaluate("fact(5)")
	fmt.Printf("  fact(5) = %g\n", result)
	result, _ = calc.Evaluate("comb(10, 3)")
	fmt.Printf("  comb(10, 3) = %g\n", result)

	// Distance functions
	result, _ = calc.Evaluate("dist2d(0, 0, 3, 4)")
	fmt.Printf("  dist2d(0, 0, 3, 4) = %g\n", result)

	// Quadratic formula: x^2 - 5x + 6 = 0 -> roots: 2, 3
	result, _ = calc.Evaluate("quadp(1, -5, 6)")
	fmt.Printf("  quadp(1, -5, 6) = %g (x² - 5x + 6 = 0)\n", result)
	result, _ = calc.Evaluate("quadm(1, -5, 6)")
	fmt.Printf("  quadm(1, -5, 6) = %g\n", result)
	fmt.Println()
}

func demo5(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 5: Variadic Functions ---")
	expressions := []string{
		"min(5, 2, 8, 1, 9)",
		"max(5, 2, 8, 1, 9)",
		"sum(1, 2, 3, 4, 5)",
		"avg(10, 20, 30)",
		"min(abs(-5), abs(-3), abs(-10))",
		"max(sqrt(4), sqrt(9), sqrt(16))",
	}
	for _, expr := range expressions {
		result, _ := calc.Evaluate(expr)
		fmt.Printf("  %s = %g\n", expr, result)
	}
	fmt.Println()
}

func demo6(calc *gorpn.Calculator) {
	fmt.Println("--- Demo 6: Real-world Formulas ---")

	// Pythagorean theorem
	result, _ := calc.EvaluateWithVars("sqrt(a^2 + b^2)", map[string]float64{"a": 3, "b": 4})
	fmt.Printf("  Pythagorean (a=3, b=4): c = %g\n", result)

	// Quadratic formula components
	calc.SetVariables(map[string]float64{"a": 1, "b": -5, "c": 6})
	result, _ = calc.Evaluate("(-b + sqrt(b^2 - 4*a*c)) / (2*a)")
	fmt.Printf("  Quadratic formula x₁ = %g\n", result)
	result, _ = calc.Evaluate("(-b - sqrt(b^2 - 4*a*c)) / (2*a)")
	fmt.Printf("  Quadratic formula x₂ = %g\n", result)

	// Compound interest: A = P(1 + r/n)^(nt)
	vars := map[string]float64{
		"P": 1000, // principal
		"r": 0.05, // annual rate
		"n": 12,   // compounds per year
		"t": 10,   // years
	}
	result, _ = calc.EvaluateWithVars("P * (1 + r/n)^(n*t)", vars)
	fmt.Printf("  Compound interest ($1000, 5%%, 10yr): $%.2f\n", result)

	// BMI calculation
	vars = map[string]float64{"weight": 70, "height": 1.75}
	result, _ = calc.EvaluateWithVars("weight / height^2", vars)
	fmt.Printf("  BMI (70kg, 1.75m): %.1f\n", result)

	// Projectile range: R = v² * sin(2θ) / g
	vars = map[string]float64{"v": 20, "theta": 45}
	result, _ = calc.EvaluateWithVars("v^2 * sin(deg(2 * theta)) / gravity", vars)
	fmt.Printf("  Projectile range (v=20m/s, θ=45°): %.2f m\n", result)

	// Circle area
	result, _ = calc.EvaluateWithVars("pi * r^2", map[string]float64{"r": 5})
	fmt.Printf("  Circle area (r=5): %.4f\n", result)

	// Sphere volume
	result, _ = calc.EvaluateWithVars("(4/3) * pi * r^3", map[string]float64{"r": 3})
	fmt.Printf("  Sphere volume (r=3): %.4f\n", result)

	fmt.Println()
}

func repl(calc *gorpn.Calculator) {
	scanner := bufio.NewScanner(os.Stdin)
	userVars := make(map[string]float64)

	fmt.Print("> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			fmt.Print("> ")
			continue
		}

		if line == "quit" || line == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if line == "vars" {
			fmt.Println("User variables:")
			for k, v := range userVars {
				fmt.Printf("  %s = %g\n", k, v)
			}
			fmt.Print("> ")
			continue
		}

		if line == "help" {
			printHelp()
			fmt.Print("> ")
			continue
		}

		// Check for variable assignment: x = 5
		if strings.Contains(line, "=") && !strings.Contains(line, "==") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				expr := strings.TrimSpace(parts[1])

				result, err := calc.EvaluateWithVars(expr, userVars)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					userVars[varName] = result
					calc.SetVariable(varName, result)
					fmt.Printf("%s = %g\n", varName, result)
				}
				fmt.Print("> ")
				continue
			}
		}

		// Evaluate expression
		result, err := calc.EvaluateWithVars(line, userVars)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("= %g\n", result)
		}
		fmt.Print("> ")
	}
}

func printHelp() {
	fmt.Println(`
Available operations:
  Operators: + - * / ^ %
  Constants: pi, e, phi, gravity, c

Built-in functions:
  Trig: sin, cos, tan, asin, acos, atan, atan2, sinh, cosh, tanh
  Math: sqrt, cbrt, abs, ceil, floor, round, trunc, exp, log, log10, log2, pow
  Stats: min, max, avg, sum (variadic)

Custom functions:
  f2c(f), c2f(c)       - Temperature conversion
  deg(r), rad(d)       - Angle conversion
  clamp(v, min, max)   - Clamp value
  lerp(a, b, t)        - Linear interpolation
  map(v, i1, i2, o1, o2) - Map range
  fact(n)              - Factorial
  comb(n, k)           - Combinations
  dist2d(x1,y1,x2,y2)  - 2D distance
  dist3d(...)          - 3D distance
  quadp(a,b,c)         - Quadratic + root
  quadm(a,b,c)         - Quadratic - root

Commands:
  vars  - Show variables
  help  - Show this help
  quit  - Exit`)
}
