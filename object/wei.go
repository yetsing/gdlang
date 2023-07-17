package object

type wei struct {
	store map[string]Object
}

func (w *wei) Type() ObjectType {
	return WEI_OBJ
}

func (w *wei) TypeIs(objectType ObjectType) bool {
	return w.Type() == objectType
}

func (w *wei) TypeNotIs(objectType ObjectType) bool {
	return w.Type() != objectType
}

func (w *wei) String() string {
	return "wei"
}

//goland:noinspection GoUnusedParameter
func (w *wei) SetAttribute(name string, value Object) Object {
	return NewError("undefined assignment: 'wei.%s", name)
}

func (w *wei) GetAttribute(name string) Object {
	if value, ok := w.store[name]; ok {
		return value
	}
	return NewError("undefined: 'wei.%s'", name)
}

func (w *wei) Add(name string, value Object) {
	w.store[name] = value
}

func newWei() *wei {
	store := make(map[string]Object)
	return &wei{store: store}
}
