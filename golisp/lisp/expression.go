package lisp

import (
	"fmt"
)

var True = &T{}

type Expression interface {
	String() string
	Pretty() string
}

type Cons struct {
	car Expression
	cdr *Cons
}

func (self *Cons) IsNull() bool {
	return self.car == nil && self.cdr == nil
}

func (self *Cons) Each(f func(Expression) Expression) Expression {
	current := self
	for !current.IsNull() {
		if ret := f(current.car); ret != nil {
			return ret
		}
		current = current.cdr
	}
	return nil
}

func (self *Cons) String() string {
	if self.IsNull() {
		return "nil"
	} else {
		return fmt.Sprintf("cons(%s %s)", self.car, self.cdr)
	}
}

func (self *Cons) Pretty() string {
	switch self.car.(type) {
	case *Char:
		str := ""
		self.Each(func(expr Expression) Expression {
			str += expr.Pretty()
			return nil
		})
		return str
	default:
		str := "("
		self.Each(func(expr Expression) Expression {
			str += expr.Pretty() + " "
			return nil
		})
		str += ")"
		return str
	}
}

type Identifier struct {
	name []rune
}

func (self *Identifier) String() string {
	return string(self.name)
}

func (self *Identifier) Pretty() string {
	return self.String()
}

type Integer struct {
	value int
}

func (self *Integer) String() string {
	return fmt.Sprint(self.value)
}

func (self *Integer) Pretty() string {
	return self.String()
}

type Char struct {
	value rune
}

func (self *Char) String() string {
	return string(self.value)
}

func (self *Char) Pretty() string {
	return self.String()
}

type T struct{}

func (self *T) String() string {
	return "T"
}

func (self *T) Pretty() string {
	return self.String()
}
