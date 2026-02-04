package gorpn

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLeftParen
	TokenRightParen
	TokenFunction
	TokenVariable
	TokenComma
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("{%v: %s}", t.Type, t.Value)
}

type Tokenizer struct {
	input string
	pos   int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: input, pos: 0}
}

func (t *Tokenizer) peek() rune {
	if t.pos >= len(t.input) {
		return 0
	}
	return rune(t.input[t.pos])
}

func (t *Tokenizer) advance() rune {
	ch := t.peek()
	t.pos++
	return ch
}

func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.input) && unicode.IsSpace(rune(t.input[t.pos])) {
		t.pos++
	}
}

func (t *Tokenizer) readNumber() string {
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

func (t *Tokenizer) readIdentifier() string {
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

func (t *Tokenizer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		t.skipWhitespace()
		if t.pos >= len(t.input) {
			break
		}

		ch := t.peek()

		// Number
		if unicode.IsDigit(ch) || (ch == '.' && t.pos+1 < len(t.input) && unicode.IsDigit(rune(t.input[t.pos+1]))) {
			tokens = append(tokens, Token{Type: TokenNumber, Value: t.readNumber()})
			continue
		}

		// Identifier (function or variable)
		if unicode.IsLetter(ch) || ch == '_' {
			ident := t.readIdentifier()
			t.skipWhitespace()
			// If followed by '(', it's a function
			if t.peek() == '(' {
				tokens = append(tokens, Token{Type: TokenFunction, Value: ident})
			} else {
				tokens = append(tokens, Token{Type: TokenVariable, Value: ident})
			}
			continue
		}

		// Operators
		if isOperator(ch) {
			// Handle unary minus/plus
			if (ch == '-' || ch == '+') && (len(tokens) == 0 ||
				tokens[len(tokens)-1].Type == TokenOperator ||
				tokens[len(tokens)-1].Type == TokenLeftParen ||
				tokens[len(tokens)-1].Type == TokenComma) {

				sign := t.advance()
				t.skipWhitespace()

				// Check if followed by a number
				if unicode.IsDigit(t.peek()) || t.peek() == '.' {
					num := t.readNumber()
					if sign == '-' {
						num = "-" + num
					}
					tokens = append(tokens, Token{Type: TokenNumber, Value: num})
				} else if unicode.IsLetter(t.peek()) || t.peek() == '_' {
					// Unary operator before variable/function - insert 0 for subtraction
					if sign == '-' {
						tokens = append(tokens, Token{Type: TokenNumber, Value: "0"})
						tokens = append(tokens, Token{Type: TokenOperator, Value: "-"})
					}
					// For unary +, just continue (it's a no-op)
				} else if t.peek() == '(' {
					// Unary operator before parenthesis
					if sign == '-' {
						tokens = append(tokens, Token{Type: TokenNumber, Value: "-1"})
						tokens = append(tokens, Token{Type: TokenOperator, Value: "*"})
					}
				} else {
					return nil, fmt.Errorf("unexpected character after unary operator at position %d", t.pos)
				}
				continue
			}
			tokens = append(tokens, Token{Type: TokenOperator, Value: string(t.advance())})
			continue
		}

		// Parentheses
		if ch == '(' {
			tokens = append(tokens, Token{Type: TokenLeftParen, Value: string(t.advance())})
			continue
		}
		if ch == ')' {
			tokens = append(tokens, Token{Type: TokenRightParen, Value: string(t.advance())})
			continue
		}

		// Comma (function argument separator)
		if ch == ',' {
			tokens = append(tokens, Token{Type: TokenComma, Value: string(t.advance())})
			continue
		}

		return nil, fmt.Errorf("unexpected character '%c' at position %d", ch, t.pos)
	}

	return tokens, nil
}

func TokensToString(tokens []Token) string {
	var parts []string
	for _, t := range tokens {
		parts = append(parts, t.Value)
	}
	return strings.Join(parts, " ")
}
