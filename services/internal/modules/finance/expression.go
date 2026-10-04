package finance

import (
	"fmt"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"strconv"
	"strings"
)

// EvaluateExpression evaluates a spreadsheet-style arithmetic expression such
// as "=12.5 + 3*2". It supports + - * /, parentheses and unary minus.
func EvaluateExpression(input string) (float64, error) {
	p := &exprParser{src: strings.TrimPrefix(strings.TrimSpace(input), "=")}
	if strings.TrimSpace(p.src) == "" {
		return 0, fmt.Errorf("%w: empty expression", core.ErrInvalidInput)
	}
	v, err := p.parseSum()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.pos != len(p.src) {
		return 0, fmt.Errorf("%w: unexpected %q at position %d", core.ErrInvalidInput, p.src[p.pos], p.pos)
	}
	return v, nil
}

type exprParser struct {
	src string
	pos int
}

func (p *exprParser) skipSpaces() {
	for p.pos < len(p.src) && (p.src[p.pos] == ' ' || p.src[p.pos] == '\t') {
		p.pos++
	}
}

func (p *exprParser) peek() byte {
	p.skipSpaces()
	if p.pos >= len(p.src) {
		return 0
	}
	return p.src[p.pos]
}

func (p *exprParser) parseSum() (float64, error) {
	left, err := p.parseProduct()
	if err != nil {
		return 0, err
	}
	for {
		op := p.peek()
		if op != '+' && op != '-' {
			return left, nil
		}
		p.pos++
		right, err := p.parseProduct()
		if err != nil {
			return 0, err
		}
		if op == '+' {
			left += right
		} else {
			left -= right
		}
	}
}

func (p *exprParser) parseProduct() (float64, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	for {
		op := p.peek()
		if op != '*' && op != '/' {
			return left, nil
		}
		p.pos++
		right, err := p.parseUnary()
		if err != nil {
			return 0, err
		}
		if op == '*' {
			left *= right
		} else {
			if right == 0 {
				return 0, fmt.Errorf("%w: division by zero", core.ErrInvalidInput)
			}
			left /= right
		}
	}
}

func (p *exprParser) parseUnary() (float64, error) {
	switch p.peek() {
	case '-':
		p.pos++
		v, err := p.parseUnary()
		return -v, err
	case '+':
		p.pos++
		return p.parseUnary()
	}
	return p.parseAtom()
}

func (p *exprParser) parseAtom() (float64, error) {
	c := p.peek()
	if c == '(' {
		p.pos++
		v, err := p.parseSum()
		if err != nil {
			return 0, err
		}
		if p.peek() != ')' {
			return 0, fmt.Errorf("%w: missing closing parenthesis", core.ErrInvalidInput)
		}
		p.pos++
		return v, nil
	}

	start := p.pos
	for p.pos < len(p.src) && (p.src[p.pos] >= '0' && p.src[p.pos] <= '9' || p.src[p.pos] == '.') {
		p.pos++
	}
	if start == p.pos {
		return 0, fmt.Errorf("%w: expected number at position %d", core.ErrInvalidInput, start)
	}
	v, err := strconv.ParseFloat(p.src[start:p.pos], 64)
	if err != nil {
		return 0, fmt.Errorf("%w: bad number %q", core.ErrInvalidInput, p.src[start:p.pos])
	}
	return v, nil
}
