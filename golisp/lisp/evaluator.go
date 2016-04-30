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
			// functions
			case "print":
				tail.Each(func(arg Expression) Expression {
					v, _ := evalExpression(arg)
					fmt.Println(v.Pretty())
					return nil
				})

				return &Cons{}, nil

			case "atom":
				v, _ := evalExpression(tail.car)
				if isAtom(v) {
					return True, nil
				} else {
					return &Cons{}, nil
				}

			case "car":
				arg, _ := evalExpression(tail.car)
				switch arg := arg.(type) {
				case *Cons:
					return arg.car, nil
				default:
					return nil, fmt.Errorf("cannot fetch car from %s", arg.Pretty())
				}

			case "cdr":
				arg, _ := evalExpression(tail.car)
				switch arg := arg.(type) {
				case *Cons:
					return arg.cdr, nil
				default:
					return nil, fmt.Errorf("cannot fetch cdr from %s", arg.Pretty())
				}

			case "cons":
				return evalExpression(tail)

			// special forms
			case "cond":
				result := tail.Each(func(cond Expression) Expression {
					switch pair := cond.(type) {
					case *Cons:
						f, s := pair.car, pair.cdr.car
						v, _ := evalExpression(f)
						if v == True {
							return s
						}
						return nil
					}
					return nil
				})

				if result != nil {
					return evalExpression(result)
				} else {
					return &Cons{}, nil
				}
			case "quote":
				return expr, nil
			}

		// list
		default:
			return expr, nil
		}
	default:
		return expr, nil
	}
	return nil, nil
}

func isAtom(expr Expression) bool {
	switch expr.(type) {
	case *Cons:
		return false
	default:
		return true
	}
}
