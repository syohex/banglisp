package banglisp

import "fmt"

type ErrUnboundVariable struct {
	name string
}

func (e *ErrUnboundVariable) Error() string {
	return fmt.Sprintf("unbound variable: %s", e.name)
}
