package lisp

import (
	"fmt"
)

type Expression interface {
	String() string
}

type Cons struct {
	car Expression
	cdr *Cons
}

func (self *Cons) String() string {
	if self.car == nil && self.cdr == nil {
		return "nil"
	} else {
		return fmt.Sprintf("cons(%s %s)", self.car, self.cdr)
	}
}

type Identifier struct {
	name []rune
}

func (self *Identifier) String() string {
	return string(self.name)
}

type Integer struct {
	value int
}

func (self *Integer) String() string {
	return fmt.Sprint(self.value)
}
