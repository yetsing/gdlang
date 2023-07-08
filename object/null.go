package object

var (
	NULL = &Null{}
)

type Null struct {
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) TypeIs(objectType ObjectType) bool {
	return n.Type() == objectType
}

func (n *Null) TypeNotIs(objectType ObjectType) bool {
	return n.Type() != objectType
}

func (n *Null) String() string {
	return "null"
}

func (n *Null) HashKey() HashKey {
	return HashKey{
		Type:  n.Type(),
		Value: 0,
	}
}
