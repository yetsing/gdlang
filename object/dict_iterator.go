package object

type DictIterator struct {
	pairs []HashPair
	index int
}

func NewDictIterator(d *Dict) *DictIterator {
	var pairs []HashPair
	for _, pair := range d.Pairs {
		pairs = append(pairs, pair)
	}
	return &DictIterator{
		pairs: pairs,
		index: 0,
	}
}

func (di *DictIterator) Type() ObjectType {
	return DICT_ITERATOR_OBJ
}

func (di *DictIterator) TypeIs(objectType ObjectType) bool {
	return di.Type() == objectType
}

func (di *DictIterator) TypeNotIs(objectType ObjectType) bool {
	return di.Type() != objectType
}

func (di *DictIterator) String() string {
	return "<dict_iterator>"
}

func (di *DictIterator) Iter() Iterator {
	return di
}

func (di *DictIterator) Next() Object {
	length := len(di.pairs)
	if di.index == length {
		return StopIteration
	}
	pair := di.pairs[di.index]
	di.index++
	return NewTuple([]Object{pair.Key, pair.Value})
}
