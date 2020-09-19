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

type Frame struct {
	bindings []struct {
		name  *Object
		value *Object
	}
}

type Environment struct {
	frames []*Frame
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

type specialFormFunction func(env *Environment, args []*Object) *Object

type SpecialForm struct {
	code specialFormFunction
}

type builtinFunctionType func(env *Environment, args []*Object) (*Object, error)

type BuiltinFunction struct {
	code     builtinFunctionType
	arity    int
	variadic bool
}

func (c *ConsCell) IsNil() bool {
	return isNull(c.car) && isNull(c.cdr)
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
		v := obj.value.(*Symbol)
		if v.value == nil {
			name := v.name.value.(string)
			return nil, &ErrUnboundVariable{name}
		}

		return v.value, nil
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
		ret := form.code(env, formArgs)
		return ret, nil
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
	default:
		return nil, fmt.Errorf("first element of cons cell is not list")
	}
}

func evalArguments(args *Object, env *Environment) ([]*Object, error) {
	var ret []*Object
	next := args
	for {
		v := next.value.(*ConsCell)
		if v.IsNil() {
			break
		}

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
		v := next.value.(*ConsCell)
		if v.IsNil() {
			break
		}

		ret = append(ret, v.car)
		next = v.cdr
	}

	return ret
}

func stringConsCell(sb *strings.Builder, obj *Object) {
	first := true
	next := obj
	for {
		v := next.value.(*ConsCell)
		if v.IsNil() {
			break
		}

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

		next = v.cdr
		first = false
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
		return v
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

func Eq(a *Object, b *Object) bool {
	return a.id == b.id
}

func isNull(v *Object) bool {
	return Eq(v, nilObj)
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

func newSpecialForm(code specialFormFunction) *Object {
	s := &SpecialForm{
		code: code,
	}
	return newObject(SpecialFormType, s)
}

func installSpecialForm(name string, code specialFormFunction) {
	sym := newSymbol(name)
	v := sym.value.(*Symbol)
	v.function = newSpecialForm(code)
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

func newEmptyEnvironment() *Environment {
	e := &Environment{}
	e.frames = append(e.frames, &Frame{})
	return e
}
