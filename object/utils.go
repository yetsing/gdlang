package object

// convertRange 对 i 进行转换，使其满足 0 <= i <= n
func convertRange(i, n int) int {
	if i < 0 {
		i += n
	}
	if i < 0 {
		return 0
	} else if i > n {
		return n
	} else {
		return i
	}
}

func equal(a, b Object) bool {
	visited := make(map[Object]bool)
	return recursiveEqual(a, b, visited)
}

func recursiveEqual(a, b Object, visited map[Object]bool) bool {
	if a.TypeNotIs(b.Type()) {
		return false
	}
	switch at := a.(type) {
	case *Integer:
		bt := b.(*Integer)
		return at.Value == bt.Value
	case *String:
		bt := b.(*String)
		return at.Value == bt.Value
	case *List:
		bt := b.(*List)
		if len(at.Elements) != len(bt.Elements) {
			return false
		}
		va := visited[a]
		vb := visited[b]
		if va && vb {
			return true
		} else if va || vb {
			return false
		}
		visited[a] = true
		visited[b] = true
		for i, ae := range at.Elements {
			be := bt.Elements[i]
			if !recursiveEqual(ae, be, visited) {
				return false
			}
		}
		return true
	case *Dict:
		bt := b.(*Dict)
		if len(at.Pairs) != len(bt.Pairs) {
			return false
		}
		va := visited[a]
		vb := visited[b]
		if va && vb {
			return true
		} else if va || vb {
			return false
		}
		visited[a] = true
		visited[b] = true
		for key, ap := range at.Pairs {
			bp, ok := bt.Pairs[key]
			if !ok || ap.Key != bp.Key {
				return false
			}
			if !recursiveEqual(ap.Value, bp.Value, visited) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
