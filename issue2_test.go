package gorpn

import (
	"fmt"
	"math"
	"testing"
)

// manyVars returns n variables named v0..v{n-1}, with values 1..n.
func manyVars(n int) map[string]float64 {
	m := make(map[string]float64, n)
	for i := 0; i < n; i++ {
		m[fmt.Sprintf("v%d", i)] = float64(i + 1)
	}
	return m
}

const issue2Expr = "(v1 - v2) / (v3 - v2)"

// allocsOf reports the allocations of one call to f, after a warm-up call.
func allocsOf(t *testing.T, f func() (float64, error)) float64 {
	t.Helper()
	if _, err := f(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return testing.AllocsPerRun(20, func() { _, _ = f() })
}

// Claim 1: the call's variables are not copied on every evaluation.
func TestIssue2_AllocsDoNotGrowWithCallVars(t *testing.T) {
	calc := New()
	few, many := manyVars(4), manyVars(1000)
	small := allocsOf(t, func() (float64, error) { return calc.EvaluateWithVars(issue2Expr, few) })
	large := allocsOf(t, func() (float64, error) { return calc.EvaluateWithVars(issue2Expr, many) })
	if large != small {
		t.Errorf("allocations grow with call variables: %v with 4, %v with 1000", small, large)
	}
}

// identity is a custom function that returns its argument.
func identity(args []float64) (float64, error) { return args[0], nil }

// issue2FnExpr calls a custom and a built-in function, so function lookup runs.
const issue2FnExpr = "f1(v1) + sqrt(v2)"

// Claim 2: the registered functions are not copied on every evaluation.
func TestIssue2_AllocsDoNotGrowWithCustomFunctions(t *testing.T) {
	few, many := New(), New()
	few.RegisterFunction("f1", 1, identity)
	for i := 0; i < 1000; i++ {
		many.RegisterFunction(fmt.Sprintf("f%d", i), 1, identity)
	}
	vars := manyVars(4)
	small := allocsOf(t, func() (float64, error) { return few.EvaluateWithVars(issue2FnExpr, vars) })
	large := allocsOf(t, func() (float64, error) { return many.EvaluateWithVars(issue2FnExpr, vars) })
	if large != small {
		t.Errorf("allocations grow with registered functions: %v with 1, %v with 1000", small, large)
	}
}

// Claim 3: EvalWithFunctions doesn't copy the functions passed to it either.
func TestIssue2_EvalWithFunctionsAllocsDoNotGrow(t *testing.T) {
	few := map[string]Function{"f1": {ArgCount: 1, Fn: identity}}
	many := map[string]Function{}
	for i := 0; i < 1000; i++ {
		many[fmt.Sprintf("f%d", i)] = Function{ArgCount: 1, Fn: identity}
	}
	vars := manyVars(4)
	small := allocsOf(t, func() (float64, error) { return EvalWithFunctions(issue2FnExpr, vars, few) })
	large := allocsOf(t, func() (float64, error) { return EvalWithFunctions(issue2FnExpr, vars, many) })
	if large != small {
		t.Errorf("allocations grow with passed functions: %v with 1, %v with 1000", small, large)
	}
}

// Claim 3: variables set once on the Calculator are not copied per evaluation
// either, so evaluating many formulas against a shared variable set is cheap.
func TestIssue2_AllocsDoNotGrowWithCalculatorVars(t *testing.T) {
	few, many := New(), New()
	few.SetVariables(manyVars(4))
	many.SetVariables(manyVars(1000))
	small := allocsOf(t, func() (float64, error) { return few.Evaluate(issue2Expr) })
	large := allocsOf(t, func() (float64, error) { return many.Evaluate(issue2Expr) })
	if large != small {
		t.Errorf("allocations grow with calculator variables: %v with 4, %v with 1000", small, large)
	}
}

// G1: user names override built-ins, layer by layer. Callers depend on this.
func TestIssue2_NameResolutionOrder(t *testing.T) {
	calc := New()
	calc.SetVariable("pi", 3)
	calc.SetVariable("x", 1)

	tests := []struct {
		name     string
		eval     func() (float64, error)
		expected float64
	}{
		{"built-in constant", func() (float64, error) { return New().Evaluate("pi") }, math.Pi},
		{"calculator var overrides constant", func() (float64, error) { return calc.Evaluate("pi") }, 3},
		{"call var overrides calculator var", func() (float64, error) { return calc.EvaluateWithVars("x", map[string]float64{"x": 2}) }, 2},
		{"call var overrides constant", func() (float64, error) { return calc.EvaluateWithVars("pi", map[string]float64{"pi": 4}) }, 4},
		{"Eval var overrides constant e", func() (float64, error) { return Eval("e * 2", map[string]float64{"e": 10}) }, 20},
		{"phi stays built in when not overridden", func() (float64, error) { return calc.Evaluate("phi") }, math.Phi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.eval()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.IsNaN(result) || math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

// G1: a custom function overrides a built-in of the same name; other
// built-ins are unaffected.
func TestIssue2_CustomFunctionOverridesBuiltin(t *testing.T) {
	calc := New()
	calc.RegisterFunction("sqrt", 1, func(args []float64) (float64, error) { return 42, nil })

	result, err := calc.Evaluate("sqrt(4)")
	if err != nil || result != 42 {
		t.Errorf("sqrt(4) = %v, %v; expected the override's 42", result, err)
	}
	result, err = calc.Evaluate("abs(-2)")
	if err != nil || result != 2 {
		t.Errorf("abs(-2) = %v, %v; expected the built-in's 2", result, err)
	}
}

// G2: an evaluation sees variables as they were when it started, even if a
// custom function changes them before they are read.
func TestIssue2_EvaluationSeesVariableSnapshot(t *testing.T) {
	calc := New()
	calc.SetVariable("x", 1)
	vars := map[string]float64{"y": 2}
	calc.RegisterFunction("bump", 0, func(args []float64) (float64, error) {
		calc.SetVariable("x", 100)
		vars["y"] = 200
		return 0, nil
	})

	// bump() runs before x and y are read
	result, err := calc.EvaluateWithVars("bump() + x + y", vars)
	if err != nil || result != 3 {
		t.Errorf("got %v, %v; expected 3 (x and y as they were at the start)", result, err)
	}
	// The next evaluation sees the new values
	result, err = calc.EvaluateWithVars("x + y", vars)
	if err != nil || result != 300 {
		t.Errorf("got %v, %v; expected 300", result, err)
	}
}

// G2: changing then clearing the Calculator's variables mid-evaluation
// doesn't affect that evaluation.
func TestIssue2_EvaluationSurvivesClearVariables(t *testing.T) {
	calc := New()
	calc.SetVariable("x", 5)
	calc.RegisterFunction("wipe", 0, func(args []float64) (float64, error) {
		calc.SetVariable("x", 100)
		calc.ClearVariables()
		return 0, nil
	})

	result, err := calc.Evaluate("wipe() + x")
	if err != nil || result != 5 {
		t.Errorf("got %v, %v; expected 5", result, err)
	}
}

// G2: functions registered or replaced mid-evaluation are not seen by it.
func TestIssue2_EvaluationSeesFunctionSnapshot(t *testing.T) {
	calc := New()
	calc.RegisterFunction("val", 0, func(args []float64) (float64, error) { return 1, nil })
	calc.RegisterFunction("swap", 0, func(args []float64) (float64, error) {
		calc.RegisterFunction("val", 0, func(args []float64) (float64, error) { return 100, nil })
		return 0, nil
	})

	result, err := calc.Evaluate("swap() + val()")
	if err != nil || result != 1 {
		t.Errorf("got %v, %v; expected 1 (val as registered at the start)", result, err)
	}
}

// G2: a function first registered mid-evaluation is unknown to it.
func TestIssue2_FunctionRegisteredMidEvaluationIsNotSeen(t *testing.T) {
	calc := New()
	calc.RegisterFunction("reg", 0, func(args []float64) (float64, error) {
		calc.RegisterFunction("late", 0, func(args []float64) (float64, error) { return 7, nil })
		return 0, nil
	})
	_, err := calc.Evaluate("reg() + late()")
	want := "evaluation error: undefined function: late"
	if err == nil || err.Error() != want {
		t.Errorf("got error %v, expected %q", err, want)
	}
}

// G2 applied to plan 004's rule: whether a function counts as custom is
// decided when the evaluation starts. Overriding exp mid-evaluation must not
// let the built-in exp's overflow through as if a custom function produced it.
func TestIssue2_CustomFlagIsSnapshot(t *testing.T) {
	calc := New()
	calc.RegisterFunction("m", 0, func(args []float64) (float64, error) {
		calc.RegisterFunction("exp", 1, func(args []float64) (float64, error) { return math.Inf(1), nil })
		return 0, nil
	})
	_, err := calc.Evaluate("m() + exp(1000)")
	want := "evaluation error: result is not a finite number: +Inf"
	if err == nil || err.Error() != want {
		t.Errorf("got error %v, expected %q", err, want)
	}
}

// G3: undefined names are reported when evaluation reaches them, so the first
// error in evaluation order wins.
func TestIssue2_ErrorOrderUnchanged(t *testing.T) {
	tests := []struct {
		expr    string
		wantErr string
	}{
		{"sqrt(-1) + nope", "evaluation error: error in function sqrt: sqrt of negative number"},
		{"nope + foo(1)", "evaluation error: undefined variable: nope"},
		{"foo(1) + nope", "evaluation error: undefined function: foo"},
		{"1e400 + nope", "evaluation error: invalid number: 1e400"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			_, err := Eval(tt.expr, nil)
			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("got error %v, expected %q", err, tt.wantErr)
			}
		})
	}
}

// G4: caller-supplied NaN/Inf still passes through, whichever layer it comes from.
func TestIssue2_NonFiniteFromEveryLayerPassesThrough(t *testing.T) {
	calc := New()
	calc.SetVariable("big", math.Inf(1))
	calc.RegisterFunction("exp", 1, func(args []float64) (float64, error) { return math.Inf(1), nil })

	tests := []struct {
		name string
		eval func() (float64, error)
	}{
		{"call var", func() (float64, error) { return New().EvaluateWithVars("x * 2", map[string]float64{"x": math.Inf(1)}) }},
		{"calculator var", func() (float64, error) { return calc.Evaluate("big * 2") }},
		{"custom override of a built-in", func() (float64, error) { return calc.Evaluate("exp(1)") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.eval()
			if err != nil || !math.IsInf(result, 1) {
				t.Errorf("got %v, %v; expected +Inf, nil", result, err)
			}
		})
	}
}

// G4: a built-in that no custom function overrides is still checked.
func TestIssue2_BuiltinStillCheckedWhenNotOverridden(t *testing.T) {
	calc := New()
	calc.RegisterFunction("other", 1, identity)
	_, err := calc.Evaluate("exp(1000)")
	want := "evaluation error: result is not a finite number: +Inf"
	if err == nil || err.Error() != want {
		t.Errorf("got error %v, expected %q", err, want)
	}
}

// G5: variables passed to one call don't persist into the Calculator.
func TestIssue2_CallVarsDoNotPersist(t *testing.T) {
	calc := New()
	if _, err := calc.EvaluateWithVars("x", map[string]float64{"x": 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := calc.Evaluate("x")
	want := "evaluation error: undefined variable: x"
	if err == nil || err.Error() != want {
		t.Errorf("got error %v, expected %q", err, want)
	}
}

func benchEvaluateWithVars(b *testing.B, n int) {
	vars := manyVars(n)
	calc := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := calc.EvaluateWithVars(issue2Expr, vars); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvaluateWithVars_4vars(b *testing.B)   { benchEvaluateWithVars(b, 4) }
func BenchmarkEvaluateWithVars_258vars(b *testing.B) { benchEvaluateWithVars(b, 258) }
