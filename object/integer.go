package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) TypeIs(objectType ObjectType) bool {
	return i.Type() == objectType
}

func (i *Integer) TypeNotIs(objectType ObjectType) bool {
	return i.Type() != objectType
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

func NewInteger(val int64) *Integer {
	return &Integer{Value: val}
}
