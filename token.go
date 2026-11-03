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

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^' || ch == '%'
}

func (t *tokenizer) Tokenize() ([]token, error) {
	var tokens []token

	for {
		t.skipWhitespace()
		if t.pos >= len(t.input) {
			break
		}

		ch := t.peek()

		// Number
		if unicode.IsDigit(ch) || (ch == '.' && t.pos+1 < len(t.input) && unicode.IsDigit(rune(t.input[t.pos+1]))) {
			tokens = append(tokens, token{Type: tokenNumber, Value: t.readNumber()})
			continue
		}

		// Identifier (function or variable)
		if unicode.IsLetter(ch) || ch == '_' {
			ident := t.readIdentifier()
			t.skipWhitespace()
			// If followed by '(', it's a function
			if t.peek() == '(' {
				tokens = append(tokens, token{Type: tokenFunction, Value: ident})
			} else {
				tokens = append(tokens, token{Type: tokenVariable, Value: ident})
			}
			continue
		}

		// Operators
		if isOperator(ch) {
			// Handle unary minus/plus
			if (ch == '-' || ch == '+') && (len(tokens) == 0 ||
				tokens[len(tokens)-1].Type == tokenOperator ||
				tokens[len(tokens)-1].Type == tokenLeftParen ||
				tokens[len(tokens)-1].Type == tokenComma) {

				sign := t.advance()
				t.skipWhitespace()

				// Check if followed by a number
				if unicode.IsDigit(t.peek()) || t.peek() == '.' {
					num := t.readNumber()
					if sign == '-' {
						num = "-" + num
					}
					tokens = append(tokens, token{Type: tokenNumber, Value: num})
				} else if unicode.IsLetter(t.peek()) || t.peek() == '_' {
					// Unary operator before variable/function - insert 0 for subtraction
					if sign == '-' {
						tokens = append(tokens, token{Type: tokenNumber, Value: "0"})
						tokens = append(tokens, token{Type: tokenOperator, Value: "-"})
					}
					// For unary +, just continue (it's a no-op)
				} else if t.peek() == '(' {
					// Unary operator before parenthesis
					if sign == '-' {
						tokens = append(tokens, token{Type: tokenNumber, Value: "-1"})
						tokens = append(tokens, token{Type: tokenOperator, Value: "*"})
					}
				} else {
					return nil, fmt.Errorf("unexpected character after unary operator at position %d", t.pos)
				}
				continue
			}
			tokens = append(tokens, token{Type: tokenOperator, Value: string(t.advance())})
			continue
		}

		// Parentheses
		if ch == '(' {
			tokens = append(tokens, token{Type: tokenLeftParen, Value: string(t.advance())})
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{Type: tokenRightParen, Value: string(t.advance())})
			continue
		}

		// Comma (function argument separator)
		if ch == ',' {
			tokens = append(tokens, token{Type: tokenComma, Value: string(t.advance())})
			continue
		}

		return nil, fmt.Errorf("unexpected character '%c' at position %d", ch, t.pos)
	}

	return tokens, nil
}
