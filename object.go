package banglisp

import (
	"fmt"
	"os"
)

type objectType int

const (
	FIXNUM objectType = iota + 1
	FLOAT
	STRING
	SYMBOL
	PACKAGE_
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

type Package_ struct {
	name  *Object
	table map[string]*Object
}

var objectID = 0

func (o objectType) String() string {
	switch o {
	case FIXNUM:
		return "FIXNUM"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case SYMBOL:
		return "SYMBOL"
	default:
		return "UNKNOWN_TYPE"
	}
}

func (obj *Object) isSelfEvaluated() bool {
	switch obj.kind {
	case FIXNUM, FLOAT, STRING:
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
	case SYMBOL:
		v := obj.value.(*Symbol)
		return v.value, nil
	default:
		return nil, fmt.Errorf("unsupported eval type")
	}
}

func (obj *Object) Print() {
	switch obj.kind {
	case FIXNUM:
		v := obj.value.(int64)
		fmt.Println(v)
	case FLOAT:
		v := obj.value.(float64)
		fmt.Println(v)
	case STRING:
		v := obj.value.(string)
		fmt.Printf("\"%s\"\n", v)
	case SYMBOL:
		v := obj.value.(*Symbol)
		n := v.name.value.(string)
		fmt.Printf("%s\n", n)
	default:
		fmt.Println("unsupported print type")
		os.Exit(1)
	}
}

func (p *Package_) setSymbol(s *Object) {
	if s.kind != SYMBOL {
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

func intern(name *Object, pack *Object) *Object {
	if pack == nil {
		pack = defaultPackage
	}

	n := name.value.(string)
	p := pack.value.(*Package_)

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
	return newObject(FIXNUM, val)
}

func newFloat(val float64) *Object {
	return newObject(FLOAT, val)
}

func newString(val string) *Object {
	return newObject(STRING, val)
}

func newSymbolInternal(val string) *Object {
	s := &Symbol{
		name:  newString(val),
		plist: nilObj,
	}

	return newObject(SYMBOL, s)
}

func newSymbol(val string) *Object {
	return intern(newString(val), nil)
}

func newPackage(name string) *Object {
	p := &Package_{
		name:  newString(name),
		table: make(map[string]*Object),
	}

	return newObject(PACKAGE_, p)
}
