package banglisp

func specialQuote(_ *Environment, args []*Object) *Object {
	// (quote exp)
	return args[0]
}

func specialIf(env *Environment, args []*Object) *Object {
	// (if cond then else)
	cond, err := args[0].Eval(env)
	if err != nil {
		panic(err)
	}

	if isNull(cond) {
		if len(args) == 2 {
			return nilObj
		}

		els, err := args[2].Eval(env)
		if err != nil {
			panic(err)
		}

		return els
	}

	then, err := args[1].Eval(env)
	if err != nil {
		panic(err)
	}

	return then
}

func specialSetf(env *Environment, args []*Object) *Object {
	// (setf sym value)
	if _, ok := env.lookupSymbol(args[0]); !ok {
		sym := args[0].value.(*Symbol)
		sym.value = args[1]
		return args[1]
	}

	// change value of local variable
	env.updateValue(args[0], args[1])
	return args[1]
}

func initSpecialForm() {
	installSpecialForm("quote", specialQuote)
	installSpecialForm("if", specialIf)
	installSpecialForm("setf", specialSetf)
}
