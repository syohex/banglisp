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
)

type Object struct {
	id    int
	kind  objectType
	value interface{}
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

	return nil, fmt.Errorf("unsupported eval type")
}

func (obj *Object) Print() {
	switch obj.kind {
	case FIXNUM:
		v, _ := obj.value.(int64)
		fmt.Println(v)
	case FLOAT:
		v, _ := obj.value.(float64)
		fmt.Println(v)
	case STRING:
		v, _ := obj.value.(string)
		fmt.Printf("\"%s\"\n", v)
	default:
		fmt.Println("unsupported print type")
		os.Exit(1)
	}
}

func newID() int {
	n := objectID
	objectID++
	return n
}

func Eq(a *Object, b *Object) bool {
	return a.id == b.id
}

func newObject(kind objectType, val interface{}) *Object {
	obj := &Object{
		id:    newID(),
		kind:  kind,
		value: val,
	}

	return obj
}

func NewFixnum(val int64) *Object {
	return newObject(FIXNUM, val)
}

func NewFloat(val float64) *Object {
	return newObject(FLOAT, val)
}

func NewString(val string) *Object {
	return newObject(STRING, val)
}
