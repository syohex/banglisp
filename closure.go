package banglisp

type Closure struct {
	name   *Object
	params []*Object
	body   []*Object
	env    *Environment
}

func newClosure(name *Object, params []*Object, body []*Object, env *Environment) *Object {
	c := &Closure{
		name:   name,
		params: params,
		body:   body,
		env:    env,
	}
	return newObject(ClosureType, c)
}

func (c *Closure) apply(env *Environment, actualArgs []*Object) (*Object, error) {
	if len(c.params) != len(actualArgs) {
		return nil, &ErrWrongNumberArguments{false, len(c.params), len(actualArgs)}
	}

	frame := &Frame{}
	for i := 0; i < len(c.params); i++ {
		frame.addBinding(c.params[i], actualArgs[i])
	}

	env.pushFrame(frame)
	defer env.popFrame(1)

	ret := nilObj
	var err error
	for _, expr := range c.body {
		ret, err = expr.Eval(env)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
