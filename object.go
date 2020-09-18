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

func (c *ConsCell) IsNil() bool {
	return isNull(c.car) && isNull(c.cdr)
}

var objectID = 0

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

func (obj *Object) Eval() (*Object, error) {
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
	default:
		return nil, fmt.Errorf("unsupported eval type")
	}
}

func stringConsCell(sb strings.Builder, obj *Object) {
	v := obj.value.(*ConsCell)
	for {
		sb.WriteString(v.car.String())
		if v.IsNil() {
			break
		}

		if v.cdr.kind == ConsCellType {
			sb.WriteByte(' ')
			stringConsCell(sb, v.cdr)
		} else {
			sb.WriteString(" . ")
			sb.WriteString(v.cdr.String())
		}
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
		stringConsCell(sb, &obj)
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
		name:  newString(val),
		plist: nilObj,
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
