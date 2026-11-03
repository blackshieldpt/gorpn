package gorpn

import (
	"fmt"
	"unicode"
)

type tokenType int

const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenLeftParen
	tokenRightParen
	tokenFunction
	tokenVariable
	tokenComma
)

type token struct {
	Type  tokenType
	Value string
}

func (t token) String() string {
	return fmt.Sprintf("{%v: %s}", t.Type, t.Value)
}

type tokenizer struct {
	input string
	pos   int
}

func newTokenizer(input string) *tokenizer {
	return &tokenizer{input: input, pos: 0}
}

func (t *tokenizer) peek() rune {
	if t.pos >= len(t.input) {
		return 0
	}
	return rune(t.input[t.pos])
}

func (t *tokenizer) advance() rune {
	ch := t.peek()
	t.pos++
	return ch
}

func (t *tokenizer) skipWhitespace() {
	for t.pos < len(t.input) && unicode.IsSpace(rune(t.input[t.pos])) {
		t.pos++
	}
}

func (t *tokenizer) readNumber() string {
	start := t.pos
	hasDecimal := false

	for t.pos < len(t.input) {
		ch := rune(t.input[t.pos])
		if ch == '.' && !hasDecimal {
			hasDecimal = true
			t.pos++
		} else if unicode.IsDigit(ch) {
			t.pos++
		} else {
			break
		}
	}

	// Handle scientific notation (e.g., 1e10, 1.5e-3)
	if t.pos < len(t.input) && (t.input[t.pos] == 'e' || t.input[t.pos] == 'E') {
		t.pos++
		if t.pos < len(t.input) && (t.input[t.pos] == '+' || t.input[t.pos] == '-') {
			t.pos++
		}
		for t.pos < len(t.input) && unicode.IsDigit(rune(t.input[t.pos])) {
			t.pos++
		}
	}

	return t.input[start:t.pos]
}

func (t *tokenizer) readIdentifier() string {
	start := t.pos
	for t.pos < len(t.input) {
		ch := rune(t.input[t.pos])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			t.pos++
		} else {
			break
		}
	}
	return t.input[start:t.pos]
}

// opNegate is the operator value for unary minus. The tokenizer creates it;
// it never appears in input text, so it can't clash with a binary operator.
const opNegate = "neg"

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^' || ch == '%'
}

func (t *tokenizer) Tokenize() ([]token, error) {
	var tokens []token
	// Operands and operators must alternate. expectOperand is true where a
	// number, identifier, '(' or unary sign must come next.
	expectOperand := true
	// One entry per open '(': true when it opens a function call
	var parens []bool
	// True right after a call's '(' with nothing since; only then may ')' follow
	callJustOpened := false

	for {
		t.skipWhitespace()
		if t.pos >= len(t.input) {
			break
		}

		ch := t.peek()
		start := t.pos
		afterCallOpen := callJustOpened
		callJustOpened = false

		// Number
		if unicode.IsDigit(ch) || (ch == '.' && t.pos+1 < len(t.input) && unicode.IsDigit(rune(t.input[t.pos+1]))) {
			if !expectOperand {
				return nil, fmt.Errorf("unexpected number at position %d", start)
			}
			tokens = append(tokens, token{Type: tokenNumber, Value: t.readNumber()})
			expectOperand = false
			continue
		}

		// Identifier (function or variable)
		if unicode.IsLetter(ch) || ch == '_' {
			if !expectOperand {
				return nil, fmt.Errorf("unexpected identifier at position %d", start)
			}
			ident := t.readIdentifier()
			t.skipWhitespace()
			// If followed by '(', it's a function; the '(' is still to come
			if t.peek() == '(' {
				tokens = append(tokens, token{Type: tokenFunction, Value: ident})
			} else {
				tokens = append(tokens, token{Type: tokenVariable, Value: ident})
				expectOperand = false
			}
			continue
		}

		// Operators
		if isOperator(ch) {
			t.advance()
			if expectOperand {
				// Only a sign can start an operand; unary plus is a no-op
				if ch == '-' {
					tokens = append(tokens, token{Type: tokenOperator, Value: opNegate})
				} else if ch != '+' {
					return nil, fmt.Errorf("unexpected operator '%c' at position %d", ch, start)
				}
				continue
			}
			tokens = append(tokens, token{Type: tokenOperator, Value: string(ch)})
			expectOperand = true
			continue
		}

		// Parentheses
		if ch == '(' {
			if !expectOperand {
				return nil, fmt.Errorf("unexpected '(' at position %d", start)
			}
			isCall := len(tokens) > 0 && tokens[len(tokens)-1].Type == tokenFunction
			parens = append(parens, isCall)
			tokens = append(tokens, token{Type: tokenLeftParen, Value: string(t.advance())})
			callJustOpened = isCall
			continue
		}
		if ch == ')' {
			// Only a call may be empty: f()
			if expectOperand && !afterCallOpen {
				return nil, fmt.Errorf("unexpected ')' at position %d", start)
			}
			// An unmatched ')' is reported by the parser
			if len(parens) > 0 {
				parens = parens[:len(parens)-1]
			}
			tokens = append(tokens, token{Type: tokenRightParen, Value: string(t.advance())})
			expectOperand = false
			continue
		}

		// Comma (function argument separator)
		if ch == ',' {
			if expectOperand {
				return nil, fmt.Errorf("unexpected ',' at position %d", start)
			}
			if len(parens) == 0 || !parens[len(parens)-1] {
				return nil, fmt.Errorf("misplaced comma at position %d: not inside a function call", start)
			}
			tokens = append(tokens, token{Type: tokenComma, Value: string(t.advance())})
			expectOperand = true
			continue
		}

		return nil, fmt.Errorf("unexpected character '%c' at position %d", ch, t.pos)
	}

	if expectOperand {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	return tokens, nil
}
