package banglisp

import (
	"fmt"
	"os"
)

type objectType int

const (
	FIXNUM objectType = iota + 1
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
	default:
		return "UNKNOWN_TYPE"
	}
}

func (obj *Object) isSelfEvaluated() bool {
	switch obj.kind {
	case FIXNUM:
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

func NewFixnum(val int64) *Object {
	obj := &Object{
		id:    newID(),
		kind:  FIXNUM,
		value: val,
	}

	return obj
}
