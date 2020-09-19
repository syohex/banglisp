package banglisp

import "fmt"

type ErrUnboundVariable struct {
	name string
}

func (e ErrUnboundVariable) Error() string {
	return fmt.Sprintf("unbound variable: %s", e.name)
}

type ErrWrongNumberArguments struct {
	variadic bool
	expected int
	got      int
}

func (e ErrWrongNumberArguments) Error() string {
	if e.variadic {
		return fmt.Sprintf("expected more than %d arguments, but got %d arguments", e.expected, e.got)
	} else {
		return fmt.Sprintf("expected %d arguments, but got %d arguments", e.expected, e.got)
	}
}

type ErrUnsupportedArgumentType struct {
	function string
	argument *Object
}

func (e ErrUnsupportedArgumentType) Error() string {
	return fmt.Sprintf("%s does not accept %v", e.function, *e.argument)
}
