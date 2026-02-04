package gorpn

import (
	"fmt"
)

type Associativity int

const (
	LeftAssociative Associativity = iota
	RightAssociative
)

type OperatorInfo struct {
	Precedence    int
	Associativity Associativity
}

var operators = map[string]OperatorInfo{
	"+": {Precedence: 2, Associativity: LeftAssociative},
	"-": {Precedence: 2, Associativity: LeftAssociative},
	"*": {Precedence: 3, Associativity: LeftAssociative},
	"/": {Precedence: 3, Associativity: LeftAssociative},
	"%": {Precedence: 3, Associativity: LeftAssociative},
	"^": {Precedence: 4, Associativity: RightAssociative},
}

// ShuntingYard converts infix tokens to RPN (postfix) notation
// using Dijkstra's Shunting Yard algorithm
func ShuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var operatorStack []Token

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case TokenNumber, TokenVariable:
			output = append(output, token)

		case TokenFunction:
			operatorStack = append(operatorStack, token)

		case TokenComma:
			// Pop operators until we find a left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != TokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses or misplaced comma")
			}

		case TokenOperator:
			op1 := operators[token.Value]
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type != TokenOperator {
					break
				}
				op2 := operators[top.Value]
				if (op1.Associativity == LeftAssociative && op1.Precedence <= op2.Precedence) ||
					(op1.Associativity == RightAssociative && op1.Precedence < op2.Precedence) {
					output = append(output, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)

		case TokenLeftParen:
			operatorStack = append(operatorStack, token)

		case TokenRightParen:
			// Pop operators until we find the matching left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != TokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses: missing '('")
			}
			// Pop the left parenthesis
			operatorStack = operatorStack[:len(operatorStack)-1]

			// If there's a function on top of the stack, pop it to output
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type == TokenFunction {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
		}
	}

	// Pop remaining operators
	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]
		if top.Type == TokenLeftParen {
			return nil, fmt.Errorf("mismatched parentheses: extra '('")
		}
		output = append(output, top)
	}

	return output, nil
}

// ShuntingYardWithArgCounts converts infix tokens to RPN and tracks function argument counts
func ShuntingYardWithArgCounts(tokens []Token) ([]Token, map[int]int, error) {
	var output []Token
	var operatorStack []Token
	argCounts := make(map[int]int) // Maps output index to arg count for functions
	var argCountStack []int        // Stack to track argument counts for nested functions
	var wasValueStack []bool       // Stack to track if we just saw a value (for counting args)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case TokenNumber, TokenVariable:
			output = append(output, token)
			// Mark that we saw a value for argument counting
			if len(wasValueStack) > 0 {
				wasValueStack[len(wasValueStack)-1] = true
			}

		case TokenFunction:
			operatorStack = append(operatorStack, token)
			argCountStack = append(argCountStack, 0)
			wasValueStack = append(wasValueStack, false)

		case TokenComma:
			// Pop operators until we find a left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != TokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, nil, fmt.Errorf("mismatched parentheses or misplaced comma")
			}
			// Increment argument count if we saw a value
			if len(argCountStack) > 0 && len(wasValueStack) > 0 && wasValueStack[len(wasValueStack)-1] {
				argCountStack[len(argCountStack)-1]++
				wasValueStack[len(wasValueStack)-1] = false
			}

		case TokenOperator:
			op1 := operators[token.Value]
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type != TokenOperator {
					break
				}
				op2 := operators[top.Value]
				if (op1.Associativity == LeftAssociative && op1.Precedence <= op2.Precedence) ||
					(op1.Associativity == RightAssociative && op1.Precedence < op2.Precedence) {
					output = append(output, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)

		case TokenLeftParen:
			operatorStack = append(operatorStack, token)
			// If this isn't a function call, push dummy values to stacks
			if len(operatorStack) < 2 || operatorStack[len(operatorStack)-2].Type != TokenFunction {
				// Not a function call, just grouping parentheses
			}

		case TokenRightParen:
			// Pop operators until we find the matching left parenthesis
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type != TokenLeftParen {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, nil, fmt.Errorf("mismatched parentheses: missing '('")
			}
			// Pop the left parenthesis
			operatorStack = operatorStack[:len(operatorStack)-1]

			// If there's a function on top of the stack, pop it to output
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].Type == TokenFunction {
				// Count the last argument if there was a value
				if len(argCountStack) > 0 && len(wasValueStack) > 0 && wasValueStack[len(wasValueStack)-1] {
					argCountStack[len(argCountStack)-1]++
				}

				fn := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]

				outputIdx := len(output)
				output = append(output, fn)

				// Store argument count
				if len(argCountStack) > 0 {
					argCounts[outputIdx] = argCountStack[len(argCountStack)-1]
					argCountStack = argCountStack[:len(argCountStack)-1]
				}
				if len(wasValueStack) > 0 {
					wasValueStack = wasValueStack[:len(wasValueStack)-1]
				}

				// Mark that we saw a value for outer function
				if len(wasValueStack) > 0 {
					wasValueStack[len(wasValueStack)-1] = true
				}
			}
		}
	}

	// Pop remaining operators
	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]
		if top.Type == TokenLeftParen {
			return nil, nil, fmt.Errorf("mismatched parentheses: extra '('")
		}
		output = append(output, top)
	}

	return output, argCounts, nil
}
