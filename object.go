package banglisp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type objectType int

const (
	FixnumType objectType = iota + 1
	FloatType
	StringType
	SymbolType
	PackageType
	ConsCellType
	SpecialFormType
	BuiltinFunctionType
	ClosureType
)

type Object struct {
	id    int
	kind  objectType
	value interface{}
}

type Symbol struct {
	name     *Object
	value    *Object
	function *Object
	plist    *Object
	package_ *Object
}

type Package struct {
	name  *Object
	table map[string]*Object
}

type ConsCell struct {
	car *Object
	cdr *Object
}

var objectID = 0

func isAtom(obj *Object) bool {
	switch obj.kind {
	case FixnumType, FloatType, StringType, SymbolType:
		return true
	default:
		return false
	}
}

func (o objectType) String() string {
	switch o {
	case FixnumType:
		return "Fixnum"
	case FloatType:
		return "Float"
	case StringType:
		return "String"
	case SymbolType:
		return "Symbol"
	case PackageType:
		return "Package"
	case ConsCellType:
		return "ConsCell"
	case BuiltinFunctionType:
		return "BuiltinFunction"
	case SpecialFormType:
		return "SpecialForm"
	case ClosureType:
		return "ClosureType"
	default:
		return "UNKNOWN_TYPE"
	}
}

func (obj *Object) isSelfEvaluated() bool {
	switch obj.kind {
	case FixnumType, FloatType, StringType:
		return true
	default:
		return false
	}
}

func (obj *Object) Eval(env *Environment) (*Object, error) {
	if obj.isSelfEvaluated() {
		return obj, nil
	}

	switch obj.kind {
	case SymbolType:
		val, ok := env.lookupSymbol(obj)
		if !ok {
			v := obj.value.(*Symbol)
			if v.value == nil {
				name := v.name.value.(string)
				return nil, &ErrUnboundVariable{name}
			}

			val = v.value
		}

		return val, nil
	case ConsCellType:
		v := obj.value.(*ConsCell)

		car, ok := v.car.value.(*Symbol)
		if !ok {
			return nil, fmt.Errorf("first element of cons cell is not list: %v(%v)", *obj, obj.kind)
		}

		if isNull(car.function) {
			return nil, fmt.Errorf("symbol '%v' does not have function", *car.name)
		}

		return v.car.apply(v.cdr, env)
	default:
		return nil, fmt.Errorf("unsupported eval type")
	}
}

func (obj *Object) apply(args *Object, env *Environment) (*Object, error) {
	switch obj.kind {
	case SymbolType:
		car := obj.value.(*Symbol)
		if car.function == nil {
			return nil, fmt.Errorf("symbol '%v' does not have function", *car.name)
		}

		return car.function.apply(args, env)
	case SpecialFormType:
		form := obj.value.(*SpecialForm)
		formArgs := noEvalArguments(args)
		if form.variadic {
			if len(formArgs) < form.arity {
				return nil, &ErrWrongNumberArguments{true, form.arity, len(formArgs)}
			}
		} else {
			if len(formArgs) != form.arity {
				return nil, &ErrWrongNumberArguments{false, form.arity, len(formArgs)}
			}
		}
		return form.code(env, formArgs)
	case BuiltinFunctionType:
		fn := obj.value.(*BuiltinFunction)
		fnArgs, err := evalArguments(args, env)
		if err != nil {
			return nil, err
		}

		if fn.variadic {
			if len(fnArgs) < fn.arity {
				return nil, &ErrWrongNumberArguments{true, fn.arity, len(fnArgs)}
			}
		} else {
			if len(fnArgs) != fn.arity {
				return nil, &ErrWrongNumberArguments{false, fn.arity, len(fnArgs)}
			}
		}

		return fn.code(env, fnArgs)
	case ClosureType:
		fn := obj.value.(*Closure)
		fnArgs, err := evalArguments(args, env)
		if err != nil {
			return nil, err
		}

		return fn.apply(env, fnArgs)
	default:
		return nil, fmt.Errorf("first element of cons cell is not list")
	}
}

func evalArguments(args *Object, env *Environment) ([]*Object, error) {
	var ret []*Object
	next := args
	for {
		if next == emptyList {
			break
		}

		v := next.value.(*ConsCell)
		ev, err := v.car.Eval(env)
		if err != nil {
			return nil, err
		}

		ret = append(ret, ev)
		next = v.cdr
	}

	return ret, nil
}

func noEvalArguments(args *Object) []*Object {
	var ret []*Object
	next := args
	for {
		if next == emptyList {
			break
		}

		v := next.value.(*ConsCell)
		ret = append(ret, v.car)
		next = v.cdr
	}

	return ret
}

func stringConsCell(sb *strings.Builder, obj *Object) {
	first := true
	next := obj
	for {
		if next == emptyList {
			break
		}

		v := next.value.(*ConsCell)
		if !first {
			first = true
			sb.WriteByte(' ')
		}

		sb.WriteString(v.car.String())

		if v.cdr.kind != ConsCellType {
			sb.WriteString(" . ")
			sb.WriteString(v.cdr.String())
			break
		}

		first = false
		next = v.cdr
	}
}

func (obj Object) String() string {
	switch obj.kind {
	case FixnumType:
		v := obj.value.(int64)
		return strconv.FormatInt(v, 10)
	case FloatType:
		v := obj.value.(float64)
		return strconv.FormatFloat(v, 'E', -1, 64)
	case StringType:
		v := obj.value.(string)
		return fmt.Sprintf(`"%s"`, v)
	case SymbolType:
		v := obj.value.(*Symbol)
		n := v.name.value.(string)
		return n
	case PackageType:
		v := obj.value.(*Package)
		n := v.name.value.(string)
		return n
	case ConsCellType:
		var sb strings.Builder
		sb.WriteByte('(')
		stringConsCell(&sb, &obj)
		sb.WriteByte(')')
		return sb.String()
	case BuiltinFunctionType:
		return "#<builtin>"
	case ClosureType:
		v := obj.value.(*Closure)
		if v.name != nil {
			return fmt.Sprintf("#<function %v>", *v.name)
		} else {
			return "#<function lambda>"
		}
	default:
		return "error: unsupported print type"
	}
}

func (p *Package) setSymbol(s *Object) {
	if s.kind != SymbolType {
		fmt.Fprintf(os.Stderr, "Invalid value passing: %v", s)
		panic("invalid")
	}

	sym := s.value.(*Symbol)
	symName := sym.name.value.(string)
	p.table[symName] = s
}

func newID() int {
	n := objectID
	objectID++
	return n
}

func objectEqual(a *Object, b *Object) bool {
	return a.id == b.id
}

func isNull(v *Object) bool {
	return objectEqual(v, nilObj)
}

func cons(car *Object, cdr *Object) *Object {
	return newConsCell(car, cdr)
}

func intern(name *Object, pack *Object) *Object {
	if pack == nil {
		pack = defaultPackage
	}

	n := name.value.(string)
	p := pack.value.(*Package)

	if v, ok := p.table[n]; ok {
		return v
	}

	newSym := newSymbolInternal(n)
	sym := newSym.value.(*Symbol)
	sym.package_ = pack
	p.table[n] = newSym

	return newSym
}

func newObject(kind objectType, val interface{}) *Object {
	obj := &Object{
		id:    newID(),
		kind:  kind,
		value: val,
	}

	return obj
}

func newFixnum(val int64) *Object {
	return newObject(FixnumType, val)
}

func newFloat(val float64) *Object {
	return newObject(FloatType, val)
}

func newString(val string) *Object {
	return newObject(StringType, val)
}

func newSymbolInternal(val string) *Object {
	s := &Symbol{
		name:     newString(val),
		function: nilObj,
		plist:    nilObj,
	}

	return newObject(SymbolType, s)
}

func newSymbol(val string) *Object {
	return intern(newString(val), nil)
}

func newPackage(name string) *Object {
	p := &Package{
		name:  newString(name),
		table: make(map[string]*Object),
	}

	return newObject(PackageType, p)
}

func newConsCell(car *Object, cdr *Object) *Object {
	c := &ConsCell{
		car: car,
		cdr: cdr,
	}

	return newObject(ConsCellType, c)
}
