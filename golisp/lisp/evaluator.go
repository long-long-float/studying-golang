package lisp

import (
	"fmt"
)

func Evaluate(exprs []Expression) error {
	rootEnv := &Environment{nil, VariableTable{}}
	for _, expr := range exprs {
		if _, err := evalExpression(expr, rootEnv); err != nil {
			return err
		}
	}
	return nil
}

func evalExpression(iexpr Expression, current *Environment) (Expression, error) {
	switch expr := iexpr.(type) {
	case *Cons:
		if expr.IsNull() {
			return expr, nil
		}

		tail := expr.cdr
		switch head := expr.car.(type) {
		case *Identifier:
			switch string(head.name) {
			// functions
			case "print":
				tail.Each(func(arg Expression) Expression {
					v, _ := evalExpression(arg, current)
					fmt.Println(v.Pretty())
					return nil
				})

				return &Cons{}, nil

			case "atom":
				v, _ := evalExpression(tail.car, current)
				if isAtom(v) {
					return True, nil
				} else {
					return &Cons{}, nil
				}

			case "eq":
				f, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
				s, err := evalExpression(tail.cdr.car, current)
				if err != nil {
					return nil, err
				}
				if f.Equals(s) {
					return True, nil
				} else {
					return &Cons{}, nil
				}

			case "car":
				arg, _ := evalExpression(tail.car, current)
				switch arg := arg.(type) {
				case *Cons:
					return arg.car, nil
				default:
					return nil, fmt.Errorf("cannot fetch car from %s", arg.Pretty())
				}

			case "cdr":
				arg, _ := evalExpression(tail.car, current)
				switch arg := arg.(type) {
				case *Cons:
					return arg.cdr, nil
				default:
					return nil, fmt.Errorf("cannot fetch cdr from %s", arg.Pretty())
				}

			case "cons":
				return evalExpression(tail, current)

			// special forms
			case "cond":
				result := tail.Each(func(cond Expression) Expression {
					switch pair := cond.(type) {
					case *Cons:
						f, s := pair.car, pair.cdr.car
						v, _ := evalExpression(f, current)
						if v == True {
							return s
						}
						return nil
					}
					return nil
				})

				if result != nil {
					return evalExpression(result, current)
				} else {
					return &Cons{}, nil
				}
			case "quote":
				return tail, nil
			case "lambda":
				args, ok := tail.car.(*Cons)
				if !ok {
					return nil, fmt.Errorf("arguments of lambda must be Cons")
				}
				return &Lambda{current, args, tail.cdr, expr}, nil
			case "define":
				id, ok := tail.car.(*Identifier)
				if !ok {
					return nil, fmt.Errorf("name of define must be Identifier")
				}
				val := tail.cdr.car

				current.vtable[string(id.name)], _ = evalExpression(val, current)

				return val, nil

			default:
				return applyLambda(expr, current)
			}

		default:
			return applyLambda(expr, current)
		}

	case *Identifier:
		val := current.find(string(expr.name))
		if val == nil {
			return nil, fmt.Errorf("undefined variable %s", string(expr.name))
		}
		return val, nil
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

func applyLambda(expr *Cons, current *Environment) (Expression, error) {
	car, _ := evalExpression(expr.car, current)
	switch head := car.(type) {
	case *Lambda:
		lambda := head
		vtable := VariableTable{}
		args := expr.cdr

		// TODO: lambda.argsとargsの長さの比較をする

		argsCurrent := args
		lambda.args.Each(func(expr Expression) Expression {
			// TODO: Identifierではなかった時のエラー処理
			id, ok := expr.(*Identifier)
			if !ok {
				return nil
			}

			vtable[string(id.name)], _ = evalExpression(argsCurrent.car, current)
			argsCurrent = argsCurrent.cdr

			return nil
		})

		env := &Environment{current, vtable}

		env.parent = current

		var retVal Expression = &Cons{}
		lambda.body.Each(func(expr Expression) Expression {
			// TODO: エラー処理
			retVal, _ = evalExpression(expr, env)
			return nil
		})

		return retVal, nil

	default:
		return nil, fmt.Errorf("cannot apply %s", head.Pretty())
	}
}
