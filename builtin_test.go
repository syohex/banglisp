package banglisp

import (
	"math"
	"strings"
	"testing"
)

func TestFunctionCallError(t *testing.T) {
	r := strings.NewReader("(-)")
	expr, err := Read(r)
	if err != nil {
		t.Errorf("Read('(-)') error=%v", err)
		return
	}

	_, err = Eval(expr)
	if _, ok := err.(*ErrWrongNumberArguments); err == nil || !ok {
		t.Error("could not get wrong number argument error")
		return
	}
}

func TestBuiltinArithmeticOps(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    float64
		isFloat bool
	}{
		// Add +
		{
			name:    "add zero arguments",
			expr:    "(+)",
			want:    0,
			isFloat: false,
		},
		{
			name:    "add one arguments",
			expr:    "(+ 1)",
			want:    1,
			isFloat: false,
		},
		{
			name:    "add more than two arguments",
			expr:    "(+ 1 2 3 4 5 6 7 8 9 10)",
			want:    55,
			isFloat: false,
		},
		{
			name:    "add float arguments",
			expr:    "(+ 1.2 3.4)",
			want:    4.6,
			isFloat: true,
		},
		// Minus -
		{
			name:    "minus one argument",
			expr:    "(- 2)",
			want:    -2,
			isFloat: false,
		},
		{
			name:    "minus arguments",
			expr:    "(- 1 2 3 4 5 6 7 8 9 10)",
			want:    -53,
			isFloat: false,
		},
		{
			name:    "minus float arguments",
			expr:    "(- 10.9 5.5 2.3)",
			want:    3.1,
			isFloat: true,
		},
		// Multiply *
		{
			name:    "mul zero argument",
			expr:    "(*)",
			want:    1,
			isFloat: false,
		},
		{
			name:    "mul arguments",
			expr:    "(* 1 2 3 4 5 6 7 8 9 10)",
			want:    3628800,
			isFloat: false,
		},
		{
			name:    "mul float arguments",
			expr:    "(* 123.45 67.89)",
			want:    8381.0205,
			isFloat: true,
		},
		// Div /
		{
			name:    "div arguments",
			expr:    "(/ 10 5 2)",
			want:    1,
			isFloat: false,
		},
		{
			name:    "div float arguments",
			expr:    "(/ 99.99 3)",
			want:    33.33,
			isFloat: true,
		},
		// Mod %
		{
			name:    "mod arguments",
			expr:    "(mod 982 3)",
			want:    1,
			isFloat: false,
		},
		{
			name:    "mod arguments 2",
			expr:    "(mod 254 255)",
			want:    254,
			isFloat: false,
		},
	}

	const epsilon = 0.000001
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.expr)
			expr, err := Read(r)
			if err != nil {
				t.Errorf("Read('%s') error=%v", tt.expr, err)
				return
			}

			obj, err := Eval(expr)
			if err != nil {
				t.Errorf("Eval('%v') error=%v", *expr, err)
				return
			}

			if tt.isFloat {
				fv := obj.value.(float64)
				if math.Abs(fv-tt.want) >= epsilon {
					t.Errorf("%s => got: %g, expected %g", tt.expr, fv, tt.want)
					return
				}
			} else {
				nv := obj.value.(int64)
				if nv != int64(tt.want) {
					t.Errorf("%s => got: %d, expected %d", tt.expr, nv, int64(tt.want))
					return
				}
			}
		})
	}
}
