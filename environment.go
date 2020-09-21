package banglisp

type Frame struct {
	bindings []bindPair
}

type Environment struct {
	frames []*Frame
}

func newEmptyEnvironment() *Environment {
	e := &Environment{}
	e.frames = make([]*Frame, 0)
	return e
}

func (e *Environment) lookupSymbol(obj *Object) (*Object, bool) {
	for _, f := range e.frames {
		for _, b := range f.bindings {
			if Eq(obj, b.name) {
				return b.value, true
			}
		}
	}

	return nil, false
}

func (e *Environment) updateValue(variable *Object, value *Object) {
	for _, f := range e.frames {
		for _, b := range f.bindings {
			if Eq(variable, b.name) {
				b.value = value
			}
		}
	}
}
