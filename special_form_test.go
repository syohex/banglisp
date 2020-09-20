package banglisp

import (
	"strings"
	"testing"
)

func TestBasicSpecialForm(t *testing.T) {
	tests := []struct {
		name string
		expr string
		kind objectType
		want interface{}
	}{
		// quote
		{
			name: "quote number",
			expr: "(quote 1)",
			kind: FixnumType,
			want: int64(1),
		},
		{
			name: "quote symbol",
			expr: "(quote foo)",
			kind: SymbolType,
			want: "foo",
		},
		{
			name: "quote string",
			expr: `(quote "hello")`,
			kind: StringType,
			want: "hello",
		},
		// if
		{
			name: "if then no else",
			expr: "(if t 10)",
			kind: FixnumType,
			want: int64(10),
		},
		{
			name: "if else no else",
			expr: "(if nil 10)",
			kind: SymbolType,
			want: "nil",
		},
		{
			name: "if else",
			expr: `(if nil 10 20)`,
			kind: FixnumType,
			want: int64(20),
		},
		{
			name: "if else with multiple expressions",
			expr: `(if nil 'foo 'bar 'baz)`,
			kind: SymbolType,
			want: "baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.expr)
			expr, err := Read(r)
			if err != nil {
				t.Errorf("Read('%s') error=%v", tt.expr, err)
				return
			}

			val, err := Eval(expr)
			if err != nil {
				t.Errorf("could not evaluate %s", tt.expr)
				return
			}

			switch tt.kind {
			case FixnumType:
				v, ok := val.value.(int64)
				if !ok {
					t.Errorf("error quoted number: %v", *val)
					return
				}

				expected := tt.want.(int64)
				if v != expected {
					t.Errorf("got unexpected number value: got: %d, expected: %d", v, expected)
					return
				}
			case StringType:
				v, ok := val.value.(string)
				if !ok {
					t.Errorf("error quoted number: %v", val)
					return
				}

				expected := tt.want.(string)
				if v != expected {
					t.Errorf("got unexpected number value: got: %s, expected: %s", v, expected)
					return
				}
			case SymbolType:
				v, ok := val.value.(*Symbol)
				if !ok {
					t.Errorf("error quoted symbol: %v", *val)
					return
				}

				name := v.name.value.(string)
				expected := tt.want.(string)
				if name != expected {
					t.Errorf("got unexpected number value: got: %s, expected: %s", name, expected)
					return
				}
			default:
			}
		})
	}
}
