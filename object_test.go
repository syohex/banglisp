package banglisp

import (
	"testing"
)

func TestNewNumber(t *testing.T) {
	num1 := NewFixnum(10)

	if !Eq(num1, num1) {
		t.Error("Eq(obj, obj) must always return true")
	}

	num2 := NewFixnum(20)
	if Eq(num1, num2) {
		t.Error("Eq(obj1, obj2) must always return false")
	}

	v, ok := num2.value.(int64)
	if !ok || v != 20 {
		t.Error("invalid value")
	}
}
