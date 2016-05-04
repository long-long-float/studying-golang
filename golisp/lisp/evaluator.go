package lisp

import (
	"fmt"
)

func Evaluate(exprs []Expression) error {
	rootEnv := &Environment{nil, VariableTable{}, nil}
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
				_, err := tail.Each(func(arg Expression) (Expression, error) {
					v, err := evalExpression(arg, current)
					if err != nil {
						return nil, err
					}
					fmt.Println(v.Pretty())
					return nil, nil
				})

				if err != nil {
					return nil, err
				}

				return &Cons{}, nil

			case "atom":
				v, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
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
				arg, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
				switch arg := arg.(type) {
				case *Cons:
					return arg.car, nil
				default:
					return nil, fmt.Errorf("cannot fetch car from %s", arg.Pretty())
				}

			case "cdr":
				arg, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
				switch arg := arg.(type) {
				case *Cons:
					return arg.cdr, nil
				default:
					return nil, fmt.Errorf("cannot fetch cdr from %s", arg.Pretty())
				}

			case "cons":
				return evalExpression(tail, current)

			case "thread/run":
				arg, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
				lambda, ok := arg.(*Lambda)
				if !ok {
					return nil, fmt.Errorf("arguments of thread/run must be Lambda")
				}

				thread := NewThread(lambda)

				go func() {
					result, err := applyLambda(&Cons{lambda, tail.cdr}, current)
					thread.NotifyFinishing(result, err)
				}()

				return thread, nil

			case "thread/wait":
				arg, err := evalExpression(tail.car, current)
				if err != nil {
					return nil, err
				}
				thread, ok := arg.(*Thread)
				if !ok {
					return nil, fmt.Errorf("arguments of thread/wait must be Thread")
				}

				result := <-thread.finishCh

				if result.err != nil {
					return nil, result.err
				} else {
					return result.expr, nil
				}

			// special forms
			case "cond":
				result, err := tail.Each(func(cond Expression) (Expression, error) {
					switch pair := cond.(type) {
					case *Cons:
						f, s := pair.car, pair.cdr.car
						v, err := evalExpression(f, current)
						if err != nil {
							return nil, err
						}
						if v == True {
							return s, nil
						}
						return nil, nil
					}
					// TODO: Consではなかった時にエラーを返す
					return nil, nil
				})

				if err != nil {
					return nil, err
				}

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

				var err error
				current.vtable[string(id.name)], err = evalExpression(val, current)

				if err != nil {
					return nil, err
				}

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
			return nil, fmt.Errorf("undefined variable or function %s", string(expr.name))
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
	car, err := evalExpression(expr.car, current)
	if err != nil {
		return nil, err
	}
	switch head := car.(type) {
	case *Lambda:
		lambda := head
		vtable := VariableTable{}
		args := expr.cdr

		// TODO: lambda.argsとargsの長さの比較をする

		argsCurrent := args
		_, err := lambda.args.Each(func(expr Expression) (Expression, error) {
			id, ok := expr.(*Identifier)
			if !ok {
				return nil, fmt.Errorf("arguments of lambda must be Identifier")
			}

			vtable[string(id.name)], _ = evalExpression(argsCurrent.car, current)
			argsCurrent = argsCurrent.cdr

			return nil, nil
		})

		if err != nil {
			return nil, err
		}

		env := &Environment{current, vtable, lambda}

		env.parent = current

		var retVal Expression = &Cons{}
		_, err = lambda.body.Each(func(expr Expression) (Expression, error) {
			var err error
			retVal, err = evalExpression(expr, env)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

		if err != nil {
			return nil, err
		}
		return retVal, nil

	default:
		return nil, fmt.Errorf("cannot apply %s", head.Pretty())
	}
}
