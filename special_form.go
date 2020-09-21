package banglisp

type specialFormFunction func(env *Environment, args []*Object) (*Object, error)

type SpecialForm struct {
	code     specialFormFunction
	arity    int
	variadic bool
}

func newSpecialForm(code specialFormFunction, arity int, variadic bool) *Object {
	s := &SpecialForm{
		code:     code,
		arity:    arity,
		variadic: variadic,
	}
	return newObject(SpecialFormType, s)
}

func installSpecialForm(name string, code specialFormFunction, arity int, variadic bool) {
	sym := newSymbol(name)
	v := sym.value.(*Symbol)
	v.function = newSpecialForm(code, arity, variadic)
}

func specialQuote(_ *Environment, args []*Object) (*Object, error) {
	// (quote exp)
	return args[0], nil
}

func specialFunction(_ *Environment, args []*Object) (*Object, error) {
	// (function symbol)
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"function", args[0]}
	}

	return sym.function, nil
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

func specialDefun(env *Environment, args []*Object) (*Object, error) {
	// (defun name (params...) body)
	nameSym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, ErrUnsupportedArgumentType{"defun", args[0]}
	}

	sym := intern(nameSym.name, nil)
	symValue := sym.value.(*Symbol)
	symValue.function = newClosure(args[0], noEvalArguments(args[1]), args[2:], env)
	return args[0], nil
}

func specialLambda(env *Environment, args []*Object) (*Object, error) {
	// (lambda (params...) body)
	fn := newClosure(nil, noEvalArguments(args[0]), args[1:], env)
	return fn, nil
}

func specialLet(env *Environment, args []*Object) (*Object, error) {
	// (let ((var1 val1) (var2 val2)) body)
	frame := &Frame{}
	next := args[0]
	for {
		if next == emptyList {
			break
		}

		iter := next.value.(*ConsCell)

		pair := iter.car.value.(*ConsCell)
		name := pair.car

		valueObj := pair.cdr.value.(*ConsCell)
		value, err := valueObj.car.Eval(env)
		if err != nil {
			return nil, err
		}
		frame.addBinding(name, value)

		next = iter.cdr
	}

	env.pushFrame(frame)
	defer env.popFrame(1)

	var err error
	ret := nilObj
	for _, expr := range args[1:] {
		ret, err = expr.Eval(env)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func specialLetStar(env *Environment, args []*Object) (*Object, error) {
	// (let* ((var1 val1) (var2 val2)) body)
	next := args[0]
	frames := 0
	for {
		if next == emptyList {
			break
		}

		iter := next.value.(*ConsCell)

		pair := iter.car.value.(*ConsCell)
		name := pair.car

		valueObj := pair.cdr.value.(*ConsCell)
		value, err := valueObj.car.Eval(env)
		if err != nil {
			return nil, err
		}

		frame := &Frame{}
		frame.addBinding(name, value)
		env.pushFrame(frame)

		frames++

		next = iter.cdr
	}

	defer env.popFrame(frames)

	var err error
	ret := nilObj
	for _, expr := range args[1:] {
		ret, err = expr.Eval(env)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func initSpecialForm() {
	installSpecialForm("quote", specialQuote, 1, false)
	installSpecialForm("function", specialFunction, 1, false)
	installSpecialForm("if", specialIf, 2, true)
	installSpecialForm("setq", specialSetq, 2, false)
	installSpecialForm("defun", specialDefun, 2, true)
	installSpecialForm("lambda", specialLambda, 1, true)
	installSpecialForm("let", specialLet, 1, true)
	installSpecialForm("let*", specialLetStar, 1, true)
}
