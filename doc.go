// Package gorpn evaluates infix maths expressions such as "2 * sqrt(x) + 1".
//
// An expression is tokenized, converted to Reverse Polish Notation with the
// shunting-yard algorithm, and evaluated on a stack. Operators, from tightest
// to loosest binding: ^ (right-associative), unary minus, then * / %, then
// + and -. Unary + and - may start any operand, so 2 * -x and -2 ^ 2 (= -4)
// work as expected.
//
// Names are ASCII letters, digits and underscores, not starting with a digit.
// A name resolves to a variable passed to the call, then a variable set on the
// Calculator, then a built-in constant (pi, e, phi); a custom function
// overrides a built-in function of the same name. Each evaluation looks names
// up once, when it starts.
//
// Use Eval or EvalWithFunctions for one-off evaluation, a Calculator to keep
// variables and custom functions between evaluations, and Compile to parse an
// expression once and evaluate it many times.
//
// Errors are prefixed with the stage that failed: "tokenization error: ",
// "parsing error: " or "evaluation error: ". A result of NaN or ±Inf is an
// error unless a variable or custom function supplied a NaN or ±Inf value.
package gorpn
