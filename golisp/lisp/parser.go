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

func (self *state) expect(char rune) error {
	if self.current() != char {
		return fmt.Errorf("unexpected %s, expect %s", string(self.current()), string(char))
	} else {
		return nil
	}
}

func (self *state) consume() rune {
	cur := self.src[self.position]
	self.position++
	return cur
}

func (self *state) currentIsSpace() bool {
	cur := self.current()
	return cur == ' ' || cur == '\r' || cur == '\n'
}

func (self *state) skipSpaces() {
	for self.currentIsSpace() {
		self.position++
	}
}

func Parse(src []rune) ([]Expression, error) {
	exprs := []Expression{}
	state := &state{src, 0}

	expr, _ := parseList(state)
	exprs = append(exprs, expr)

	return exprs, nil
}

func parseList(state *state) (Expression, error) {
	if err := state.expect('('); err != nil {
		return nil, err
	}

	state.consume()

	state.skipSpaces()

	id, _ := parseIdentifier(state)
	state.skipSpaces()
	arg, _ := parseInteger(state)

	if err := state.expect(')'); err != nil {
		return nil, err
	}

	return &Cons{id, &Cons{arg, &Cons{nil, nil}}}, nil
}

func parseIdentifier(state *state) (Expression, error) {
	name := []rune{}
	reg := regexp.MustCompile(`[a-zA-Z]`)
	for reg.MatchString(string(state.current())) {
		name = append(name, state.consume())
	}

	return &Identifier{name}, nil
}

func parseInteger(state *state) (Expression, error) {
	num := []rune{}
	reg := regexp.MustCompile(`[0-9]`)
	for reg.MatchString(string(state.current())) {
		num = append(num, state.consume())
	}

	value, _ := strconv.Atoi(string(num))
	return &Integer{value}, nil
}
