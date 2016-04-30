package lisp

type VariableTable map[string]Expression

type Environment struct {
	parent *Environment
	vtable VariableTable
}

func (self *Environment) find(name string) Expression {
	if val, ok := self.vtable[name]; ok {
		return val
	} else {
		if self.parent != nil {
			return self.parent.find(name)
		}
		return nil
	}
}
