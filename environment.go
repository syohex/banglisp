package banglisp

type bindPair struct {
	name  *Object
	value *Object
}

type Frame struct {
	bindings []bindPair
}

func (f *Frame) addBinding(name *Object, value *Object) {
	f.bindings = append(f.bindings, bindPair{name, value})
}

type Environment struct {
	frames []*Frame
}

func newEmptyEnvironment() *Environment {
	e := &Environment{}
	e.frames = make([]*Frame, 0)
	return e
}

func (e *Environment) pushFrame(f *Frame) {
	e.frames = append([]*Frame{f}, e.frames...)
}

func (e *Environment) popFrame(count int) {
	e.frames = e.frames[count:]
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
