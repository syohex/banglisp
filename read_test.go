package banglisp

import (
	"math/big"
	"strings"
	"testing"
)

func TestReadFixnum(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    int64
		wantErr bool
	}{
		{
			name:    "positive number",
			expr:    "1234",
			want:    1234,
			wantErr: false,
		},
		{
			name:    "one digit",
			expr:    "9",
			want:    9,
			wantErr: false,
		},
		{
			name:    "negative number",
			expr:    "-12345",
			want:    -12345,
			wantErr: false,
		},
		{
			name:    "invalid number",
			expr:    "-123xyz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.expr)
			got, err := Read(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v", err)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Errorf("no error but return value is nil")
				return
			}

			if got.kind != FIXNUM {
				t.Errorf("got invalid type object [input]: %s -> %v", tt.expr, got.kind)
				return
			}

			if v, ok := got.value.(int64); !(ok && v == tt.want) {
				t.Errorf("got invalid value object [input]: %s -> Got: %d, Expected: %d", tt.expr, v, tt.want)
				return
			}
		})
	}
}

func TestReadFloat(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    *big.Float
		wantErr bool
	}{
		{
			name:    "positive floating number",
			expr:    "123.5",
			want:    big.NewFloat(123.5),
			wantErr: false,
		},
		{
			name:    "negative floating number",
			expr:    "-123.5",
			want:    big.NewFloat(-123.5),
			wantErr: false,
		},
		{
			name:    "end with point",
			expr:    "-123.",
			want:    big.NewFloat(-123),
			wantErr: false,
		},
		{
			name:    "more than one point",
			expr:    "123.45.67",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.expr)
			got, err := Read(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v", err)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Errorf("no error but return value is nil")
				return
			}

			if got.kind != FLOAT {
				t.Errorf("got invalid type object [input]: %s -> %v", tt.expr, got.kind)
				return
			}

			v, ok := got.value.(float64)
			if !ok {
				t.Errorf("got invalid value object")
				return
			}

			if tt.want.Cmp(big.NewFloat(v)) != 0 {
				t.Errorf("got: %g, expected: %v", v, tt.want)
				return
			}
		})
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    string
		wantErr bool
	}{
		{
			name:    "text",
			expr:    `"Hello World"`,
			want:    "Hello World",
			wantErr: false,
		},
		{
			name:    "string number",
			expr:    `"1234"`,
			want:    "1234",
			wantErr: false,
		},
		{
			name:    "contain escaped character",
			expr:    `"foo\nbar"`,
			want:    "foo\nbar",
			wantErr: false,
		},
		{
			name:    "unterminated string",
			expr:    `"unterminated`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.expr)
			got, err := Read(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v", err)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Errorf("no error but return value is nil")
				return
			}

			if got.kind != STRING {
				t.Errorf("got invalid type object [input]: %s -> %v", tt.expr, got.kind)
				return
			}

			if v, ok := got.value.(string); !(ok && v == tt.want) {
				t.Errorf("got invalid value object [input]: %s -> Got: %s, Expected: %s", tt.expr, v, tt.want)
				return
			}
		})
	}
}
