package lisp

import (
	"fmt"
)

func Evaluate(exprs []Expression) error {
	for _, expr := range exprs {
		if _, err := evalExpression(expr); err != nil {
			return err
		}
	}
	return nil
}

func evalExpression(iexpr Expression) (Expression, error) {
	switch expr := iexpr.(type) {
	case *Cons:
		tail := expr.cdr
		switch head := expr.car.(type) {
		case *Identifier:
			switch string(head.name) {
			case "print":
				tail.Each(func(arg Expression) {
					fmt.Println(arg.Pretty())
				})

				return &Cons{}, nil
			}
		}
	}
	return nil, nil
}
