package object

import (
	"bytes"
	"fmt"
	"strings"
)

type HashPair struct {
	Key   Object
	Value Object
}

type Dict struct {
	Pairs map[HashKey]HashPair
}

func NewDict(pairs map[HashKey]HashPair) *Dict {
	return &Dict{Pairs: pairs}
}

func (d *Dict) Type() ObjectType {
	return DICT_OBJ
}

func (d *Dict) TypeIs(objectType ObjectType) bool {
	return d.Type() == objectType
}

func (d *Dict) TypeNotIs(objectType ObjectType) bool {
	return d.Type() != objectType
}

func (d *Dict) String() string {
	var out bytes.Buffer

	var elements []string
	for _, pair := range d.Pairs {
		elements = append(elements, fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}
