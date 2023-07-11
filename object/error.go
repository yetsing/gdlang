package object

import "fmt"

type Error struct {
	Message string
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func WrongNumberArgument(got, want int) *Error {
	return NewError("wrong number of arguments. got=%d, want=%d", got, want)
}

func WrongNumberArgument2(got, min, max int) *Error {
	return NewError("wrong number of arguments. got=%d, want=%d-%d", got, min, max)
}

func wrongArgumentType(got ObjectType) *Error {
	return NewError("wrong argument type: '%s'", got)
}

func wrongArgumentTypeAt(got ObjectType, at int) *Error {
	return NewError("wrong argument type: '%s' at %d", got, at)
}

func attributeError(otype, name string) *Error {
	return NewError("'%s' object has not attribute '%s'", otype, name)
}

func Unreachable(msg string) *Error {
	return NewError("unreachable %s", msg)
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) TypeIs(objectType ObjectType) bool {
	return e.Type() == objectType
}

func (e *Error) TypeNotIs(objectType ObjectType) bool {
	return e.Type() != objectType
}

func (e *Error) String() string {
	return fmt.Sprintf("Error: %s", e.Message)
}
