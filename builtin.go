package banglisp

import "fmt"

func builtinEq(_ *Environment, args []*Object) (*Object, error) {
	// (eq a b)
	if Eq(args[0], args[1]) {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinNull(_ *Environment, args []*Object) (*Object, error) {
	// (null a)
	if isNull(args[0]) {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinAtom(_ *Environment, args []*Object) (*Object, error) {
	// (atom a)
	if isAtom(args[0]) {
		return tObj, nil
	}

	return nilObj, nil
}

func floatValue(obj *Object) (float64, bool, error) {
	switch v := obj.value.(type) {
	case int64:
		return float64(v), false, nil
	case float64:
		return v, true, nil
	default:
		return 0, false, fmt.Errorf("unsuppored type")
	}
}

func builtinAdd(_ *Environment, args []*Object) (*Object, error) {
	// (+ n1 n2 ....)
	if len(args) == 0 {
		return newFixnum(0), nil
	}

	hasFloat := false
	ret := 0.0

	for _, arg := range args {
		f, isFloat, err := floatValue(arg)
		if err != nil {
			return nil, &ErrUnsupportedArgumentType{"+", arg}
		}

		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret += f
	}

	if hasFloat {
		return newFloat(ret), nil
	}

	return newFixnum(int64(ret)), nil
}

func builtinMinus(_ *Environment, args []*Object) (*Object, error) {
	// (- n1 n2 ....)
	ret, hasFloat, err := floatValue(args[0])
	if err != nil {
		return nil, &ErrUnsupportedArgumentType{"-", args[0]}
	}

	if len(args) == 1 {
		if hasFloat {
			return newFloat(ret * -1), nil
		}

		return newFixnum(int64(ret) * -1), nil
	}

	for _, arg := range args[1:] {
		f, isFloat, err := floatValue(arg)
		if err != nil {
			return nil, &ErrUnsupportedArgumentType{"-", arg}
		}

		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret -= f
	}

	if hasFloat {
		return newFloat(ret), nil
	}

	return newFixnum(int64(ret)), nil
}

func builtinMul(_ *Environment, args []*Object) (*Object, error) {
	// (* n1 n2 ....)
	if len(args) == 0 {
		return newFixnum(1), nil
	}

	hasFloat := false
	ret := 1.0

	for _, arg := range args {
		f, isFloat, err := floatValue(arg)
		if err != nil {
			return nil, &ErrUnsupportedArgumentType{"*", arg}
		}

		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret *= f
	}

	if hasFloat {
		return newFloat(ret), nil
	}

	return newFixnum(int64(ret)), nil
}

func builtinDiv(_ *Environment, args []*Object) (*Object, error) {
	// (/ n1 n2 ....)
	ret, hasFloat, err := floatValue(args[0])
	if err != nil {
		return nil, &ErrUnsupportedArgumentType{"/", args[0]}
	}

	for _, arg := range args[1:] {
		f, isFloat, err := floatValue(arg)
		if err != nil {
			return nil, &ErrUnsupportedArgumentType{"/", arg}
		}

		if !hasFloat && isFloat {
			hasFloat = true
		}

		ret /= f
	}

	if hasFloat {
		return newFloat(ret), nil
	}

	return newFixnum(int64(ret)), nil
}

func builtinMod(_ *Environment, args []*Object) (*Object, error) {
	// (% n1 n2 ....)
	var ret int64
	switch v := args[0].value.(type) {
	case int64:
		ret = v
	default:
		return nil, &ErrUnsupportedArgumentType{"mod", args[0]}
	}

	for _, arg := range args[1:] {
		switch v := arg.value.(type) {
		case int64:
			ret %= v
		default:
			return nil, &ErrUnsupportedArgumentType{"mod", arg}
		}
	}

	return newFixnum(ret), nil
}

func builtinCar(_ *Environment, args []*Object) (*Object, error) {
	c, ok := args[0].value.(*ConsCell)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"car", args[0]}
	}
	return c.car, nil
}

func builtinCdr(_ *Environment, args []*Object) (*Object, error) {
	c, ok := args[0].value.(*ConsCell)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"cdr", args[0]}
	}
	return c.cdr, nil
}

func builtinCons(_ *Environment, args []*Object) (*Object, error) {
	return cons(args[0], args[1]), nil
}

func builtinPrint(_ *Environment, args []*Object) (*Object, error) {
	fmt.Printf("%v\n", args[0])
	return nilObj, nil
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

	// utility
	installBuiltinFunction("print", builtinPrint, 1, false)
}
