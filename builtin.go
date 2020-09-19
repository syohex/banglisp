package banglisp

func builtinEq(_ *Environment, args []*Object) *Object {
	// (eq a b)
	if Eq(args[0], args[1]) {
		return tObj
	}

	return nilObj
}

func builtinNull(_ *Environment, args []*Object) *Object {
	// (null a)
	if isNull(args[0]) {
		return tObj
	}

	return nilObj
}

func builtinAtom(_ *Environment, args []*Object) *Object {
	// (atom a)
	if isAtom(args[0]) {
		return tObj
	}

	return nilObj
}

func floatValue(obj *Object) (float64, bool) {
	switch v := obj.value.(type) {
	case int64:
		return float64(v), false
	case float64:
		return v, true
	default:
		// error
		panic("unsupported type: '+'")
	}
}

func builtinAdd(_ *Environment, args []*Object) *Object {
	// (+ n1 n2 ....)
	if len(args) == 0 {
		return newFixnum(0)
	}

	hasFloat := false
	ret := 0.0

	for _, arg := range args {
		f, isFloat := floatValue(arg)
		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret += f
	}

	if hasFloat {
		return newFloat(ret)
	}

	return newFixnum(int64(ret))
}

func builtinMinus(_ *Environment, args []*Object) *Object {
	// (- n1 n2 ....)
	ret, hasFloat := floatValue(args[0])
	if len(args) == 1 {
		if hasFloat {
			return newFloat(ret * -1)
		}

		return newFixnum(int64(ret) * -1)
	}

	for _, arg := range args[1:] {
		f, isFloat := floatValue(arg)
		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret -= f
	}

	if hasFloat {
		return newFloat(ret)
	}

	return newFixnum(int64(ret))
}

func builtinMul(_ *Environment, args []*Object) *Object {
	// (* n1 n2 ....)
	if len(args) == 0 {
		return newFixnum(1)
	}

	hasFloat := false
	ret := 1.0

	for _, arg := range args {
		f, isFloat := floatValue(arg)
		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret *= f
	}

	if hasFloat {
		return newFloat(ret)
	}

	return newFixnum(int64(ret))
}

func builtinDiv(_ *Environment, args []*Object) *Object {
	// (/ n1 n2 ....)
	sum, hasFloat := floatValue(args[0])

	for _, arg := range args[1:] {
		switch v := arg.value.(type) {
		case int64:
			sum /= float64(v)
		case float64:
			hasFloat = true
			sum /= v
		default:
			panic("unsupported type: '/'")
		}
	}

	if hasFloat {
		return newFloat(sum)
	}

	return newFixnum(int64(sum))
}

func builtinMod(_ *Environment, args []*Object) *Object {
	// (% n1 n2 ....)
	var ret int64
	switch v := args[0].value.(type) {
	case int64:
		ret = v
	default:
		panic("unsupported type: '%'")
	}

	for _, arg := range args[1:] {
		switch v := arg.value.(type) {
		case int64:
			ret %= v
		default:
			panic(" unsupported type: '%'")
		}
	}

	return newFixnum(ret)
}

func builtinCar(_ *Environment, args []*Object) *Object {
	c := args[0].value.(*ConsCell)
	return c.car
}

func builtinCdr(_ *Environment, args []*Object) *Object {
	c := args[0].value.(*ConsCell)
	return c.cdr
}

func builtinCons(_ *Environment, args []*Object) *Object {
	return cons(args[0], args[1])
}

func initBuiltinFunctions() {
	installBuiltinFunction("eq", builtinEq, 2, false)
	installBuiltinFunction("null", builtinNull, 1, false)
	installBuiltinFunction("atom", builtinAtom, 1, false)

	// arithmetic operators
	installBuiltinFunction("+", builtinAdd, 0, true)
	installBuiltinFunction("-", builtinMinus, 1, true)
	installBuiltinFunction("*", builtinMul, 0, true)
	installBuiltinFunction("/", builtinDiv, 1, true)
	installBuiltinFunction("mod", builtinMod, 1, true)

	// cons cell operations
	installBuiltinFunction("car", builtinCar, 1, false)
	installBuiltinFunction("first", builtinCar, 1, false)
	installBuiltinFunction("cdr", builtinCdr, 1, false)
	installBuiltinFunction("rest", builtinCdr, 1, false)
	installBuiltinFunction("cons", builtinCons, 2, false)
}
