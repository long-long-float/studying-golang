package lisp

type VariableTable map[string]Expression

type Environment struct {
	parent *Environment
	vtable VariableTable

	lambda *Lambda
}

func (self *Environment) find(name string) Expression {
	if self.lambda != nil {
		if val := self.lambda.parent.find(name); val != nil {
			return val
		}
	}

	if val, ok := self.vtable[name]; ok {
		return val
	} else {
		if self.parent != nil {
			return self.parent.find(name)
		}
		return nil
	}
}
