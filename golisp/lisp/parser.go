package lisp

import (
	"fmt"
	"regexp"
	"strconv"
)

type state struct {
	src      []rune
	position int
}

func (self *state) current() rune {
	return self.src[self.position]
}

func (self *state) currentAsString() string {
	if self.isEOF() {
		return "EOF"
	} else {
		return string(self.current())
	}
}

func (self *state) expect(char rune) error {
	if self.current() != char {
		return fmt.Errorf("unexpected %s, expect %s", self.currentAsString(), string(char))
	} else {
		return nil
	}
}

func (self *state) consume() rune {
	cur := self.src[self.position]
	if !self.isEOF() {
		self.position++
	}
	return cur
}

func (self *state) currentIsSpace() bool {
	cur := self.current()
	return cur == ' ' || cur == '\r' || cur == '\n'
}

func (self *state) skipSpaces() {
	for self.currentIsSpace() && !self.isEOF() {
		self.consume()
	}
}

func (self *state) isEOF() bool {
	return self.position >= len(self.src)-1
}

func Parse(src []rune) ([]Expression, error) {
	exprs := []Expression{}
	state := &state{src, 0}

	for !state.isEOF() {
		expr, err := parseExpression(state)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)

		state.skipSpaces()
	}

	return exprs, nil
}

func parseExpression(state *state) (Expression, error) {
	cur := state.current()
	switch {
	case ('a' <= cur && cur <= 'z') || ('A' <= cur && cur <= 'Z'):
		return parseIdentifier(state)
	case (cur == '-') || ('0' <= cur && cur <= '9'):
		return parseInteger(state)
	case cur == '\'':
		return parseChar(state)
	case cur == '"':
		return parseString(state)
	case cur == '(':
		return parseList(state)
	}

	return nil, fmt.Errorf("unexpected %s at %d", state.currentAsString(), state.position)
}

func parseList(state *state) (Expression, error) {
	if err := state.expect('('); err != nil {
		return nil, err
	}
	state.consume()

	state.skipSpaces()

	head := &Cons{}
	current := head
	for state.current() != ')' && !state.isEOF() {
		var err error
		current.car, err = parseExpression(state)
		if err != nil {
			return nil, err
		}
		current.cdr = &Cons{}

		current = current.cdr

		state.skipSpaces()
	}

	if err := state.expect(')'); err != nil {
		return nil, err
	}
	state.consume()

	return head, nil
}

func parseIdentifier(state *state) (Expression, error) {
	firstReg := regexp.MustCompile(`[a-zA-Z]`)
	reg := regexp.MustCompile(`[a-zA-Z0-9\-/]`)
	name := parseWhile(state, firstReg, reg)
	switch string(name) {
	case "t":
		return True, nil
	case "nil":
		return &Cons{}, nil
	default:
		return &Identifier{name}, nil
	}
}

func parseInteger(state *state) (Expression, error) {
	firstReg := regexp.MustCompile(`[0-9\-]`)
	reg := regexp.MustCompile(`[0-9]`)
	value, _ := strconv.Atoi(string(parseWhile(state, firstReg, reg)))
	return &Integer{value}, nil
}

func parseChar(state *state) (Expression, error) {
	if err := state.expect('\''); err != nil {
		return nil, err
	}
	state.consume()

	reg := regexp.MustCompile(`[^']`)
	if !reg.MatchString(string(state.current())) {
		return nil, fmt.Errorf("unexpected ' at %d", state.position)
	}
	value := state.consume()

	if err := state.expect('\''); err != nil {
		return nil, err
	}
	state.consume()

	return &Char{value}, nil
}

func parseString(state *state) (Expression, error) {
	if err := state.expect('"'); err != nil {
		return nil, err
	}
	state.consume()

	reg := regexp.MustCompile(`[^"]`)
	value := parseWhile(state, reg, reg)

	if err := state.expect('"'); err != nil {
		return nil, err
	}
	state.consume()

	head := &Cons{&Identifier{[]rune("quote")}, &Cons{}}
	current := head.cdr
	for _, cur := range value {
		current.car = &Char{rune(cur)}
		current.cdr = &Cons{}

		current = current.cdr
	}

	return head, nil
}

func parseWhile(state *state, firstReg, reg *regexp.Regexp) []rune {
	value := []rune{}
	if firstReg.MatchString(string(state.current())) {
		value = append(value, state.consume())
	}
	for reg.MatchString(string(state.current())) {
		value = append(value, state.consume())
	}

	return value
}
