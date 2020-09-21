package banglisp

import (
	"fmt"
	"math"
)

func floatValue(obj *Object) (float64, bool, error) {
	switch v := obj.value.(type) {
	case int64:
		return float64(v), false, nil
	case float64:
		return v, true, nil
	default:
		return 0, false, fmt.Errorf("unsupported type")
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

func builtinEqual(_ *Environment, args []*Object) (*Object, error) {
	v1, t1, err1 := floatValue(args[0])
	if err1 != nil {
		return nil, err1
	}

	v2, t2, err2 := floatValue(args[1])
	if err2 != nil {
		return nil, err2
	}

	const epsilon = 0.00000001
	if t1 || t2 {
		if math.Abs(v1-v2) < epsilon {
			return tObj, nil
		}

		return nilObj, nil
	}

	if int64(v1) == int64(v2) {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinLessThan(_ *Environment, args []*Object) (*Object, error) {
	v1, _, err1 := floatValue(args[0])
	if err1 != nil {
		return nil, err1
	}

	v2, _, err2 := floatValue(args[1])
	if err2 != nil {
		return nil, err2
	}

	if v1 < v2 {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinLessThanEqual(_ *Environment, args []*Object) (*Object, error) {
	v1, _, err1 := floatValue(args[0])
	if err1 != nil {
		return nil, err1
	}

	v2, _, err2 := floatValue(args[1])
	if err2 != nil {
		return nil, err2
	}

	if v1 < v2 {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinGreaterThan(_ *Environment, args []*Object) (*Object, error) {
	v1, _, err1 := floatValue(args[0])
	if err1 != nil {
		return nil, err1
	}

	v2, _, err2 := floatValue(args[1])
	if err2 != nil {
		return nil, err2
	}

	if v1 > v2 {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinGreaterThanEqual(_ *Environment, args []*Object) (*Object, error) {
	v1, _, err1 := floatValue(args[0])
	if err1 != nil {
		return nil, err1
	}

	v2, _, err2 := floatValue(args[1])
	if err2 != nil {
		return nil, err2
	}

	if v1 >= v2 {
		return tObj, nil
	}

	return nilObj, nil
}

func builtinSin(_ *Environment, args []*Object) (*Object, error) {
	v, _, err := floatValue(args[0])
	if err != nil {
		return nil, err
	}

	return newFloat(math.Sin(v)), nil
}

func builtinCos(_ *Environment, args []*Object) (*Object, error) {
	v, _, err := floatValue(args[0])
	if err != nil {
		return nil, err
	}

	return newFloat(math.Cos(v)), nil
}

func builtinTan(_ *Environment, args []*Object) (*Object, error) {
	v, _, err := floatValue(args[0])
	if err != nil {
		return nil, err
	}

	return newFloat(math.Tan(v)), nil
}

func initNumberFunctions() {
	installBuiltinFunction("+", builtinAdd, 0, true)
	installBuiltinFunction("-", builtinMinus, 1, true)
	installBuiltinFunction("*", builtinMul, 0, true)
	installBuiltinFunction("/", builtinDiv, 1, true)
	installBuiltinFunction("mod", builtinMod, 1, true)

	installBuiltinFunction("=", builtinEqual, 2, false)
	installBuiltinFunction("<", builtinLessThan, 2, false)
	installBuiltinFunction("<=", builtinLessThanEqual, 2, false)
	installBuiltinFunction(">", builtinGreaterThan, 2, false)
	installBuiltinFunction(">=", builtinGreaterThanEqual, 2, false)

	installBuiltinFunction("sin", builtinSin, 1, false)
	installBuiltinFunction("cos", builtinCos, 1, false)
	installBuiltinFunction("tan", builtinTan, 1, false)
}
