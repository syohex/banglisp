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

func TestDefunFunction(t *testing.T) {
	input := "(funcall (function +) 1 2 3 4)"
	r := strings.NewReader(input)
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('%s') error=%v", input, err)
		return
	}

	val, err := Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", input, err)
		return
	}

	v, ok := val.value.(int64)
	if !ok {
		t.Errorf("function add does not return fixnum value: %v", *expr)
		return
	}

	if v != 10 {
		t.Errorf("%s return unexpected value: got %d, expected: 10", input, v)
		return
	}
}

func TestDefunSimple(t *testing.T) {
	addFunc := `
(defun add (a b)
  (+ a b))
`

	r := strings.NewReader(addFunc)
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('%s') error=%v", addFunc, err)
		return
	}

	_, err = Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", addFunc, err)
		return
	}

	callExpr := "(add 10 20)"
	expr, err = Read(strings.NewReader(callExpr))
	if err != nil {
		t.Errorf("Read('%s') error=%v", callExpr, err)
		return
	}

	val, err := Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", addFunc, err)
		return
	}

	v, ok := val.value.(int64)
	if !ok {
		t.Errorf("%s does not return fixnum value: %v", callExpr, *expr)
		return
	}

	if v != 30 {
		t.Errorf("%s return unexpected value: got %d, expected: 30", callExpr, v)
		return
	}
}

func TestLambdaSimple(t *testing.T) {
	lambdaExpr := "(funcall (lambda (a b) (* a b)) 5 4)"

	r := strings.NewReader(lambdaExpr)
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('%s') error=%v", lambdaExpr, err)
		return
	}

	val, err := Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", lambdaExpr, err)
		return
	}

	v, ok := val.value.(int64)
	if !ok {
		t.Errorf("return value is not fixnum type: %v", *expr)
		return
	}

	if v != 20 {
		t.Errorf("%s return unexpected value: got %d, expected: 20", lambdaExpr, v)
		return
	}
}

func TestLetSimple(t *testing.T) {
	letExpr := `
(let ((a 10)
      (b (+ 10 10))
      (c 30))
  (* a b c))
`

	r := strings.NewReader(letExpr)
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('%s') error=%v", letExpr, err)
		return
	}

	val, err := Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", letExpr, err)
		return
	}

	v, ok := val.value.(int64)
	if !ok {
		t.Errorf("return value is not fixnum type: %v", *expr)
		return
	}

	if v != 6000 {
		t.Errorf("%s return unexpected value: got %d, expected: 6000", letExpr, v)
		return
	}
}

func TestLetStarSimple(t *testing.T) {
	input := `
(let* ((a 10)
       (b (+ a a a))
       (c (* a b)))
  (* a b c))
`

	r := strings.NewReader(input)
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('%s') error=%v", input, err)
		return
	}

	val, err := Eval(expr)
	if err != nil {
		t.Errorf("could not evaluate %s: %v", input, err)
		return
	}

	v, ok := val.value.(int64)
	if !ok {
		t.Errorf("return value is not fixnum: %v", *expr)
		return
	}

	if v != 90000 {
		t.Errorf("%s return unexpected value: got %d, expected: 6000", input, v)
		return
	}
}

func TestSimpleOrAnd(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want int64
		nil  bool
	}{
		{
			name: "or",
			expr: "(or nil nil nil 99)",
			want: 99,
		},
		{
			name: "or nil",
			expr: "(or nil nil)",
			nil:  true,
		},
		{
			name: "and",
			expr: "(and t t t 42)",
			want: 42,
		},
		{
			name: "and nil",
			expr: "(and nil t t 42)",
			nil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Read(strings.NewReader(tt.expr))
			if err != nil {
				t.Errorf("Read('%s') error=%v", tt.expr, err)
				return
			}

			val, err := Eval(expr)
			if err != nil {
				t.Errorf("could not evaluate %s: %v", tt.expr, err)
				return
			}

			if tt.nil {
				if !isNull(val) {
					t.Errorf("%s does not return nil: %v", tt.expr, *val)
				}
				return
			}

			v, ok := val.value.(int64)
			if !ok {
				t.Errorf("return value is not fixnum value: %v", *expr)
				return
			}

			if v != tt.want {
				t.Errorf("%s return unexpected value: got %d, expected: %d", tt.expr, v, tt.want)
				return
			}
		})
	}
}
