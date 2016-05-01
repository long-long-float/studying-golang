package lisp

import (
	"fmt"
)

var True = &T{}

type Expression interface {
	String() string
	Pretty() string
	Equals(expr Expression) bool
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

func (self *Cons) Equals(expr Expression) bool {
	if cons, ok := expr.(*Cons); ok {
		if self.IsNull() && self.IsNull() {
			return true
		} else if self.IsNull() || self.IsNull() {
			return false
		}
		return self.car.Equals(cons.car) && self.cdr.Equals(cons.cdr)
	} else {
		return false
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

func (self *Identifier) Equals(expr Expression) bool {
	panic("cannot compare Identifier")
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

func (self *Integer) Equals(expr Expression) bool {
	if i, ok := expr.(*Integer); ok {
		return self.value == i.value
	} else {
		return false
	}
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

func (self *Char) Equals(expr Expression) bool {
	if char, ok := expr.(*Char); ok {
		return self.value == char.value
	} else {
		return false
	}
}

type T struct{}

func (self *T) String() string {
	return "T"
}

func (self *T) Pretty() string {
	return self.String()
}

func (self *T) Equals(expr Expression) bool {
	_, ok := expr.(*T)
	return ok
}

type Lambda struct {
	parent *Environment

	args *Cons
	body *Cons

	self *Cons
}

func (self *Lambda) String() string {
	return self.self.String()
}

func (self *Lambda) Pretty() string {
	return "<lambda: " + self.String() + ">"
}

func (self *Lambda) Equals(expr Expression) bool {
	// TODO: これだと引数の名前まで一致していないと等しくならないので修正する
	if lambda, ok := expr.(*Lambda); ok {
		return self.self.Equals(lambda.self)
	}
	return false
}

type Thread struct {
	lambda *Lambda

	finishCh chan Expression
}

func (self *Thread) String() string {
	return "<thread: " + self.lambda.String() + ">"
}

func (self *Thread) Pretty() string {
	return self.String()
}

func (self *Thread) Equals(expr Expression) bool {
	return self == expr
}
