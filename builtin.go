package banglisp

import "fmt"

type builtinFunctionType func(env *Environment, args []*Object) (*Object, error)

type BuiltinFunction struct {
	code     builtinFunctionType
	arity    int
	variadic bool
}

func newBuiltinFunction(code builtinFunctionType, arity int, variadic bool) *Object {
	bf := &BuiltinFunction{
		code:     code,
		arity:    arity,
		variadic: variadic,
	}
	return newObject(BuiltinFunctionType, bf)
}

func installBuiltinFunction(name string, code builtinFunctionType, arity int, variadic bool) {
	sym := newSymbol(name)
	v := sym.value.(*Symbol)
	v.function = newBuiltinFunction(code, arity, variadic)
}

func builtinNot(_ *Environment, args []*Object) (*Object, error) {
	// (not obj)
	if args[0] == nilObj {
		return tObj, nil
	}

	return tObj, nil
}

func builtinEq(_ *Environment, args []*Object) (*Object, error) {
	// (eq a b)
	if objectEqual(args[0], args[1]) {
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
	if args[1] == nilObj {
		return cons(args[0], emptyList), nil
	}

	return cons(args[0], args[1]), nil
}

func builtinPrint(_ *Environment, args []*Object) (*Object, error) {
	fmt.Printf("%v\n", args[0])
	return nilObj, nil
}

func builtinFuncall(env *Environment, args []*Object) (*Object, error) {
	switch fn := args[0].value.(type) {
	case *Symbol:
		if isNull(fn.function) {
			return nil, fmt.Errorf("void function: %v", *fn)
		}
		switch fn.function.kind {
		case BuiltinFunctionType:
			fn := fn.function.value.(*BuiltinFunction)
			return fn.code(env, args[1:])
		case ClosureType:
			fn := fn.function.value.(*Closure)
			return fn.apply(env, args[1:])
		default:
			return nil, &ErrUnsupportedArgumentType{"funcall", fn.function}
		}
	case *BuiltinFunction:
		return fn.code(env, args[1:])
	case *Closure:
		return fn.apply(env, args[1:])
	default:
		return nil, &ErrUnsupportedArgumentType{"funcall", args[0]}
	}
}

func builtinSymbolName(_ *Environment, args []*Object) (*Object, error) {
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"symbol-name", args[0]}
	}

	return sym.name, nil
}

func builtinSymbolValue(_ *Environment, args []*Object) (*Object, error) {
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"symbol-value", args[0]}
	}

	return sym.value, nil
}

func builtinSymbolFunction(_ *Environment, args []*Object) (*Object, error) {
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"symbol-function", args[0]}
	}

	return sym.function, nil
}

func builtinSymbolPlist(_ *Environment, args []*Object) (*Object, error) {
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"symbol-plist", args[0]}
	}

	return sym.plist, nil
}

func builtinSymbolPackage(_ *Environment, args []*Object) (*Object, error) {
	sym, ok := args[0].value.(*Symbol)
	if !ok {
		return nil, &ErrUnsupportedArgumentType{"symbol-package", args[0]}
	}

	return sym.package_, nil
}

func initBuiltinFunctions() {
	installBuiltinFunction("eq", builtinEq, 2, false)
	installBuiltinFunction("not", builtinNot, 1, false)
	installBuiltinFunction("null", builtinNull, 1, false)
	installBuiltinFunction("atom", builtinAtom, 1, false)

	// cons cell operations
	installBuiltinFunction("car", builtinCar, 1, false)
	installBuiltinFunction("first", builtinCar, 1, false)
	installBuiltinFunction("cdr", builtinCdr, 1, false)
	installBuiltinFunction("rest", builtinCdr, 1, false)
	installBuiltinFunction("cons", builtinCons, 2, false)

	// utility
	installBuiltinFunction("print", builtinPrint, 1, false)
	installBuiltinFunction("funcall", builtinFuncall, 0, true)

	// symbol functions
	installBuiltinFunction("symbol-name", builtinSymbolName, 1, false)
	installBuiltinFunction("symbol-value", builtinSymbolValue, 1, false)
	installBuiltinFunction("symbol-function", builtinSymbolFunction, 1, false)
	installBuiltinFunction("symbol-plist", builtinSymbolPlist, 1, false)
	installBuiltinFunction("symbol-package", builtinSymbolPackage, 1, false)
}
