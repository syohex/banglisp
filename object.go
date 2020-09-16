package banglisp

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

func (o *Object) isSelfEvaluated() bool {
	switch o.kind {
	case FIXNUM:
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

func NewFixnum(val int64) *Object {
	obj := &Object{
		id:    newID(),
		kind:  FIXNUM,
		value: val,
	}

	return obj
}
