package object

import "hash/fnv"

type String struct {
	Value string
}

func NewString(val string) *String {
	return &String{Value: val}
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) TypeIs(objectType ObjectType) bool {
	return s.Type() == objectType
}

func (s *String) TypeNotIs(objectType ObjectType) bool {
	return s.Type() != objectType
}

func (s *String) String() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}
