package banglisp

import (
	"testing"
)

func TestNewFixnum(t *testing.T) {
	num1 := NewFixnum(10)

	if !Eq(num1, num1) {
		t.Error("Eq(obj, obj) must always return true")
		return
	}

	num2 := NewFixnum(20)
	if Eq(num1, num2) {
		t.Error("Eq(obj1, obj2) must always return false")
		return
	}

	v, ok := num2.value.(int64)
	if !ok || v != 20 {
		t.Error("invalid value")
		return
	}
}

func TestEvalFixnum(t *testing.T) {
	num1 := NewFixnum(10)
	ev, err := num1.Eval()
	if err != nil {
		t.Error("failed to evaluate fixnum object")
		return
	}

	v, ok := ev.value.(int64)
	if !ok || v != 10 {
		t.Error("fixnum must be self-evaluated type")
		return
	}
}
