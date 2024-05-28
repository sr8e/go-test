package structs

type empty struct{}

var dummy = empty{}

type setable interface {
	int | string
}

type Set[T setable] map[T]empty

func NewSet[T setable](values ...T) Set[T] {
	s := make(Set[T])
	s.Add(values...)
	return s
}

func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		s[v] = dummy
	}
}

func (s Set[T]) Clone() Set[T] {
	newSet := NewSet[T]()
	for k := range s {
		newSet[k] = dummy
	}
	return newSet
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	unionSet := s.Clone()
	for k := range other {
		unionSet[k] = dummy
	}
	return unionSet
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	interSet := NewSet[T]()
	for k := range other {
		if _, ok := s[k]; ok {
			interSet[k] = dummy
		}
	}
	return interSet
}

func (s Set[T]) ToSlice() []T {
	keys := make([]T, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}
