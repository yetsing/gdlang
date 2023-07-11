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

func (d *Dict) GetItem(key Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("unhashable type: '%s'", key.Type())
	}

	pair, ok := d.Pairs[hashKey.HashKey()]
	if !ok {
		return NewError("key '%s' does not exist", key.String())
	}
	return pair.Value
}

func (d *Dict) SetItem(key, value Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("unhashable type: '%s'", key.Type())
	}
	d.Pairs[hashKey.HashKey()] = HashPair{
		Key:   key,
		Value: value,
	}
	return nil
}
