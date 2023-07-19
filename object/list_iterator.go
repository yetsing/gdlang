package object

type ListIterator struct {
	l     *List
	index int
}

func NewListIterator(l *List) *ListIterator {
	return &ListIterator{
		l:     l,
		index: 0,
	}
}

func (li *ListIterator) Type() ObjectType {
	return LIST_ITERATOR_OBJ
}

func (li *ListIterator) TypeIs(objectType ObjectType) bool {
	return li.Type() == objectType
}

func (li *ListIterator) TypeNotIs(objectType ObjectType) bool {
	return li.Type() != objectType
}

func (li *ListIterator) String() string {
	return "<list_iterator>"
}

func (li *ListIterator) Iter() Iterator {
	return li
}

func (li *ListIterator) Next() Object {
	length := len(li.l.Elements)
	if li.index == length {
		return StopIteration
	}
	val := li.l.Elements[li.index]
	li.index++
	return NewTuple([]Object{NewInteger(int64(li.index - 1)), val})
}
