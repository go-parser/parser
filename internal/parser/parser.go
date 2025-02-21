package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Parse is the main entry point for parsing an input string
func Parse(input string) (string, error) {
	parser := &parser{
		input: input,
	}
	return parser.parse()
}

// TokenType represents different types of tokens in the expression
type TokenType int

// Token types for operators, literals, and other symbols
const (
	Literal    TokenType = iota // String, number literals
	Add                         // +
	Sub                         // -
	Mul                         // *
	Div                         // /
	Mod                         // %
	OpenParen                   // (
	CloseParen                  // )
	At                          // @ for function calls
	Dollar                      // $ for variables
	Identifier                  // Variable/function names
	Comma                       // ,
	And                         // &&
	Or                          // ||
	Not                         // !
	Gt                          // >
	Gte                         // >=
	Lt                          // <
	Lte                         // <=
	Eq                          // ==
	Ne                          // !=
)

// Token represents a single token with its type and value
type Token struct {
	Type  TokenType
	Value string
}

// parser holds the state for parsing expressions
type parser struct {
	input  string  // Input string to parse
	tokens []Token // Tokenized input
	pos    int     // Current position in tokens
}

// parse tokenizes the input and starts parsing the expression
func (p *parser) parse() (string, error) {
	tokens, err := tokenize(p.input)
	if err != nil {
		return "", err
	}
	p.tokens = tokens
	return p.parseLogicalExpression()
}

// parseLogicalExpression handles logical operators (AND, OR)
func (p *parser) parseLogicalExpression() (string, error) {
	left, err := p.parseComparisonExpression()
	if err != nil {
		return "", err
	}

	for p.pos < len(p.tokens) {
		switch p.tokens[p.pos].Type {
		case And:
			p.pos++
			right, err := p.parseComparisonExpression()
			if err != nil {
				return "", err
			}
			left = fmt.Sprintf("and(%s,%s)", left, right)
		case Or:
			p.pos++
			right, err := p.parseComparisonExpression()
			if err != nil {
				return "", err
			}
			left = fmt.Sprintf("or(%s,%s)", left, right)
		default:
			return left, nil
		}
	}
	return left, nil
}

// parseComparisonExpression handles comparison operators and NOT operations
func (p *parser) parseComparisonExpression() (string, error) {
	if p.pos >= len(p.tokens) {
		return "", errors.New("unexpected end of input")
	}

	if p.tokens[p.pos].Type == Not {
		p.pos++
		expr, err := p.parseComparisonExpression()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("not(%s)", expr), nil
	}

	if p.tokens[p.pos].Type == OpenParen {
		p.pos++
		startPos := p.pos
		expr, err := p.parseLogicalExpression()
		if err != nil {
			p.pos = startPos
			expr, err = p.additive()
			if err != nil {
				return "", err
			}
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != CloseParen {
			return "", errors.New("expected right parenthesis")
		}
		p.pos++

		if p.pos < len(p.tokens) {
			switch p.tokens[p.pos].Type {
			case Eq, Ne, Gt, Gte, Lt, Lte:
				left := expr
				operator := p.tokens[p.pos]
				p.pos++
				right, err := p.additive()
				if err != nil {
					return "", err
				}
				switch operator.Type {
				case Eq:
					return fmt.Sprintf("eq(%s,%s)", left, right), nil
				case Ne:
					return fmt.Sprintf("ne(%s,%s)", left, right), nil
				case Gt:
					return fmt.Sprintf("gt(%s,%s)", left, right), nil
				case Gte:
					return fmt.Sprintf("gte(%s,%s)", left, right), nil
				case Lt:
					return fmt.Sprintf("lt(%s,%s)", left, right), nil
				case Lte:
					return fmt.Sprintf("lte(%s,%s)", left, right), nil
				}
			}
		}
		return expr, nil
	}

	left, err := p.additive()
	if err != nil {
		return "", err
	}

	if p.pos < len(p.tokens) {
		switch p.tokens[p.pos].Type {
		case Eq:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("eq(%s,%s)", left, right), nil
		case Ne:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("ne(%s,%s)", left, right), nil
		case Gt:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("gt(%s,%s)", left, right), nil
		case Gte:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("gte(%s,%s)", left, right), nil
		case Lt:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("lt(%s,%s)", left, right), nil
		case Lte:
			p.pos++
			right, err := p.additive()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("lte(%s,%s)", left, right), nil
		}
	}
	return left, nil
}

// additive handles addition and subtraction operations
func (p *parser) additive() (string, error) {
	term, err := p.term()
	if err != nil {
		return "", err
	}
	for p.pos < len(p.tokens) {
		switch p.tokens[p.pos].Type {
		case Add:
			p.pos++
			right, err := p.term()
			if err != nil {
				return "", err
			}
			term = fmt.Sprintf("add(%s,%s)", term, right)
		case Sub:
			p.pos++
			right, err := p.term()
			if err != nil {
				return "", err
			}
			term = fmt.Sprintf("sub(%s,%s)", term, right)
		default:
			return term, nil
		}
	}
	return term, nil
}

// term handles multiplication, division, and modulo operations
func (p *parser) term() (string, error) {
	factor, err := p.factor()
	if err != nil {
		return "", err
	}
	for p.pos < len(p.tokens) {
		switch p.tokens[p.pos].Type {
		case Mul:
			p.pos++
			right, err := p.factor()
			if err != nil {
				return "", err
			}
			factor = fmt.Sprintf("multi(%s,%s)", factor, right)
		case Div:
			p.pos++
			right, err := p.factor()
			if err != nil {
				return "", err
			}
			factor = fmt.Sprintf("div(%s,%s)", factor, right)
		case Mod:
			p.pos++
			right, err := p.factor()
			if err != nil {
				return "", err
			}
			factor = fmt.Sprintf("mod(%s,%s)", factor, right)
		default:
			return factor, nil
		}
	}
	return factor, nil
}

// factor handles parentheses, function calls, variables, and literals
func (p *parser) factor() (string, error) {
	if p.tokens[p.pos].Type == OpenParen {
		p.pos++
		expr, err := p.additive()
		if err != nil {
			return "", err
		}
		if p.tokens[p.pos].Type != CloseParen {
			return "", errors.New("expected )")
		}
		p.pos++
		return expr, nil
	}
	if p.tokens[p.pos].Type == At {
		p.pos++
		if p.tokens[p.pos].Type != Identifier {
			return "", errors.New("expected function name")
		}
		name := p.tokens[p.pos].Value
		p.pos++
		if p.tokens[p.pos].Type != OpenParen {
			return "", errors.New("expected (")
		}
		p.pos++
		args := []string{}
		for p.tokens[p.pos].Type != CloseParen {
			arg, err := p.parseLogicalExpression()
			if err != nil {
				return "", err
			}
			args = append(args, arg)
			if p.tokens[p.pos].Type == CloseParen {
				break
			}
			if p.tokens[p.pos].Type != Comma {
				return "", errors.New("expected ,")
			}
			p.pos++
		}
		if p.tokens[p.pos].Type != CloseParen {
			return "", errors.New("expected )")
		}
		p.pos++
		return fmt.Sprintf("%s(%s)", name, strings.Join(args, ",")), nil
	}
	if p.tokens[p.pos].Type == Dollar {
		if p.tokens[p.pos+1].Type != Identifier {
			return "", errors.New("expected variable name")
		}
		name := p.tokens[p.pos].Value + p.tokens[p.pos+1].Value
		p.pos += 2
		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == Add {
			p.pos++
			right, err := p.term()
			if err != nil {
				return "", err
			}
			name = fmt.Sprintf("add(%s,%s)", name, right)
		}
		return name, nil
	}
	if p.tokens[p.pos].Type == Literal {
		value := p.tokens[p.pos].Value
		p.pos++
		return value, nil
	}
	return "", errors.New("expected literal")
}

// tokenize converts the input string into a sequence of tokens
// Handles operators, numbers, strings, identifiers, and special characters
func tokenize(input string) ([]Token, error) {
	tokens := []Token{}
	for i := 0; i < len(input); i++ {
		switch {
		case input[i] == '+':
			tokens = append(tokens, Token{Add, "+"})
		case input[i] == '-':
			tokens = append(tokens, Token{Sub, "-"})
		case input[i] == '*':
			tokens = append(tokens, Token{Mul, "*"})
		case input[i] == '/':
			tokens = append(tokens, Token{Div, "/"})
		case input[i] == '%':
			tokens = append(tokens, Token{Mod, "%"})
		case input[i] == '(':
			tokens = append(tokens, Token{OpenParen, "("})
		case input[i] == ')':
			tokens = append(tokens, Token{CloseParen, ")"})
		case input[i] == '@':
			tokens = append(tokens, Token{At, "@"})
		case input[i] == '$':
			tokens = append(tokens, Token{Dollar, "$"})
		case input[i] == ',':
			tokens = append(tokens, Token{Comma, ","})
		case input[i] == '&':
			if input[i+1] == '&' {
				tokens = append(tokens, Token{And, "&&"})
				i++
			} else {
				return nil, errors.New("expected &&")
			}
		case input[i] == '|':
			if input[i+1] == '|' {
				tokens = append(tokens, Token{Or, "||"})
				i++
			} else {
				return nil, errors.New("expected ||")
			}
		case input[i] == '!':
			if input[i+1] == '=' {
				tokens = append(tokens, Token{Ne, "!="})
				i++
			} else {
				tokens = append(tokens, Token{Not, "!"})
			}
		case input[i] == '>':
			if input[i+1] == '=' {
				tokens = append(tokens, Token{Gte, ">="})
				i++
			} else {
				tokens = append(tokens, Token{Gt, ">"})
			}
		case input[i] == '<':
			if input[i+1] == '=' {
				tokens = append(tokens, Token{Lte, "<="})
				i++
			} else {
				tokens = append(tokens, Token{Lt, "<"})
			}
		case input[i] == '=':
			if i+1 < len(input) && input[i+1] == '=' {
				tokens = append(tokens, Token{Eq, "=="})
				i++
			} else {
				return nil, errors.New("expected ==")
			}
		case input[i] == '"':
			j := i + 1
			for j < len(input) && input[j] != '"' {
				j++
			}
			if j == len(input) {
				return nil, errors.New(`expected "`)
			}
			tokens = append(tokens, Token{Literal, input[i+1:j] + `:str`})
			i = j
		case unicode.IsLetter(rune(input[i])):
			j := i
			for j < len(input) && (unicode.IsLetter(rune(input[j])) || unicode.IsDigit(rune(input[j])) || input[j] == '_') {
				j++
			}
			tokens = append(tokens, Token{Identifier, input[i:j]})
			i = j - 1
		case unicode.IsDigit(rune(input[i])):
			j := i
			for j < len(input) && (unicode.IsDigit(rune(input[j])) || input[j] == '.') {
				j++
			}

			token := input[i:j]
			if strings.Contains(token, ".") {
				tokens = append(tokens, Token{Literal, token + `:float`})
			} else {
				tokens = append(tokens, Token{Literal, token + `:int`})
			}
			i = j - 1
		case input[i] == ' ':
			continue
		default:
			return nil, fmt.Errorf("unexpected character: %c", input[i])
		}
	}
	return tokens, nil
}
