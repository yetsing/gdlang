package object

func NewStringIterator(s *String) *StringIterator {
	return &StringIterator{
		value: []rune(s.Value),
		index: 0,
	}
}

type StringIterator struct {
	value []rune
	index int
}

func (si *StringIterator) Type() ObjectType {
	return LIST_ITERATOR_OBJ
}

func (si *StringIterator) TypeIs(objectType ObjectType) bool {
	return si.Type() == objectType
}

func (si *StringIterator) TypeNotIs(objectType ObjectType) bool {
	return si.Type() != objectType
}

func (si *StringIterator) String() string {
	return "<str_iterator>"
}

func (si *StringIterator) Next() Object {
	if si.index >= len(si.value) {
		return StopIteration
	}
	e := si.value[si.index]
	tp := NewTuple([]Object{NewInteger(int64(si.index)), NewString(string(e))})
	si.index++
	return tp
}
