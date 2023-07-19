package object

type Tuple struct {
	Elements []Object
}

func NewTuple(elements []Object) *Tuple {
	return &Tuple{
		Elements: elements,
	}
}

func (t *Tuple) Type() ObjectType {
	return TUPLE_OBJ
}

func (t *Tuple) TypeIs(objectType ObjectType) bool {
	return t.Type() == objectType
}

func (t *Tuple) TypeNotIs(objectType ObjectType) bool {
	return t.Type() != objectType
}

func (t *Tuple) String() string {
	visited := make(map[Object]bool)
	return objectString(t, visited)
}
