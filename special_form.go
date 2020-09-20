package banglisp

func specialQuote(_ *Environment, args []*Object) (*Object, error) {
	// (quote exp)
	return args[0], nil
}

func specialIf(env *Environment, args []*Object) (*Object, error) {
	// (if cond then else)
	cond, err := args[0].Eval(env)
	if err != nil {
		return nil, err
	}

	if isNull(cond) {
		if len(args) == 2 {
			return nilObj, nil
		}

		var val *Object
		var err error
		for _, arg := range args[2:] {
			val, err = arg.Eval(env)
			if err != nil {
				return nil, err
			}
		}

		return val, nil
	}

	// then
	return args[1].Eval(env)
}

func specialSetq(env *Environment, args []*Object) (*Object, error) {
	// (setq sym value)
	if _, ok := env.lookupSymbol(args[0]); !ok {
		sym := args[0].value.(*Symbol)
		sym.value = args[1]
		return args[1], nil
	}

	// change value of local variable
	env.updateValue(args[0], args[1])
	return args[1], nil
}

func initSpecialForm() {
	installSpecialForm("quote", specialQuote, 1, false)
	installSpecialForm("if", specialIf, 2, true)
	installSpecialForm("setq", specialSetq, 2, false)
}
