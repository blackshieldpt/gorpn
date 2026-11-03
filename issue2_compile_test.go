package gorpn

import (
	"math"
	"sync"
	"testing"
)

// Claim 4: an expression can be parsed once and evaluated many times, with
// the same results as EvaluateWithVars.
func TestIssue2_CompileEvaluatesManyTimes(t *testing.T) {
	calc := New()
	calc.SetVariable("k", 10)
	calc.RegisterFunction("double", 1, func(args []float64) (float64, error) { return args[0] * 2, nil })

	expr, err := calc.Compile("double(x) + k * pi")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	for _, x := range []float64{0, 1, -2.5, 100} {
		vars := map[string]float64{"x": x}
		got, err := expr.Eval(vars)
		if err != nil {
			t.Fatalf("Eval(x=%v): %v", x, err)
		}
		want, err := calc.EvaluateWithVars("double(x) + k * pi", vars)
		if err != nil {
			t.Fatalf("EvaluateWithVars(x=%v): %v", x, err)
		}
		if math.IsNaN(got) || got != want {
			t.Errorf("x=%v: Eval = %v, EvaluateWithVars = %v", x, got, want)
		}
	}
}

// Claim 4: compile and evaluation errors carry the same prefixes and messages
// as EvaluateWithVars.
func TestIssue2_CompileErrorsMatchEvaluate(t *testing.T) {
	compileErrs := []struct {
		expr    string
		wantErr string
	}{
		{"2 +", "tokenization error: unexpected end of expression"},
		{"(2 + 3", "parsing error: mismatched parentheses: extra '('"},
	}
	for _, tt := range compileErrs {
		t.Run(tt.expr, func(t *testing.T) {
			expr, err := Compile(tt.expr)
			if err == nil || err.Error() != tt.wantErr || expr != nil {
				t.Errorf("got %v, %v; expected nil, %q", expr, err, tt.wantErr)
			}
		})
	}

	expr, err := Compile("x + 1")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	_, err = expr.Eval(nil)
	want := "evaluation error: undefined variable: x"
	if err == nil || err.Error() != want {
		t.Errorf("got error %v, expected %q", err, want)
	}
}

// Claim 4 with G1: the package-level Compile uses built-ins only, and call
// variables still override constants.
func TestIssue2_PackageCompile(t *testing.T) {
	expr, err := Compile("pi * 2")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	if got, err := expr.Eval(nil); err != nil || math.Abs(got-2*math.Pi) > 1e-9 {
		t.Errorf("Eval(nil) = %v, %v; expected 2π", got, err)
	}
	if got, err := expr.Eval(map[string]float64{"pi": 3}); err != nil || got != 6 {
		t.Errorf("Eval(pi=3) = %v, %v; expected 6", got, err)
	}
}

// Each Eval sees the Calculator's variables and functions as they are when
// it starts, so later changes reach a compiled expression.
func TestIssue2_CompiledExprSeesCalculatorStateAtEval(t *testing.T) {
	calc := New()
	expr, err := calc.Compile("pi + f(1)")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	calc.RegisterFunction("f", 1, func(args []float64) (float64, error) { return args[0], nil })
	if got, err := expr.Eval(nil); err != nil || math.Abs(got-(math.Pi+1)) > 1e-9 {
		t.Errorf("got %v, %v; expected π+1", got, err)
	}
	calc.SetVariable("pi", 3)
	if got, err := expr.Eval(nil); err != nil || got != 4 {
		t.Errorf("got %v, %v; expected 4 after overriding pi", got, err)
	}
}

// Claims 3 and 4: a compiled expression's evaluation cost doesn't grow with
// the variables in scope.
func TestIssue2_CompiledEvalAllocsDoNotGrowWithVars(t *testing.T) {
	calc := New()
	calc.SetVariables(manyVars(1000))
	expr, err := calc.Compile(issue2Expr)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	few, many := manyVars(4), manyVars(1000)
	small := allocsOf(t, func() (float64, error) { return expr.Eval(few) })
	large := allocsOf(t, func() (float64, error) { return expr.Eval(many) })
	if large != small {
		t.Errorf("allocations grow with variables: %v with 4, %v with 1000", small, large)
	}
}

// Claim 4: Eval does not parse again. It must allocate less than tokenizing
// and parsing the expression alone, which a re-parse on every Eval can't.
func TestIssue2_EvalDoesNotReparse(t *testing.T) {
	expr, err := Compile(issue2Expr)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	parse := testing.AllocsPerRun(20, func() {
		tokens, _ := newTokenizer(issue2Expr).Tokenize()
		_, _, _ = shuntingYardWithArgCounts(tokens)
	})
	vars := manyVars(4)
	eval := allocsOf(t, func() (float64, error) { return expr.Eval(vars) })
	if eval >= parse {
		t.Errorf("Eval allocates %v, no fewer than parsing alone (%v)", eval, parse)
	}
}

// An Expr can be shared between goroutines while the Calculator is not being
// changed. Run with -race.
func TestIssue2_ConcurrentEval(t *testing.T) {
	calc := New()
	calc.SetVariable("k", 2)
	expr, err := calc.Compile("x * k")
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				x := float64(g*1000 + i)
				got, err := expr.Eval(map[string]float64{"x": x})
				if err != nil || got != 2*x {
					t.Errorf("Eval(x=%v) = %v, %v; expected %v", x, got, err, 2*x)
					return
				}
			}
		}(g)
	}
	wg.Wait()
}

func BenchmarkCompiledEval_258vars(b *testing.B) {
	vars := manyVars(258)
	expr, err := New().Compile(issue2Expr)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := expr.Eval(vars); err != nil {
			b.Fatal(err)
		}
	}
}
