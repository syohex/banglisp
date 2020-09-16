package banglisp

type objectType int

const (
	NUMBER objectType = iota + 1
)

type Object struct {
	id    int
	kind  objectType
	value interface{}
}

var objectID = 0

func (o *Object) isSelfEvaluated() bool {
	switch o.kind {
	case NUMBER:
		return true
	default:
		return false
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

func NewNumber(val int) *Object {
	obj := &Object{
		id:    newID(),
		kind:  NUMBER,
		value: val,
	}

	return obj
}
