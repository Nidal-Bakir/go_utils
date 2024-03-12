package go_utils

type Set[T comparable] struct {
	internalMap map[T]bool
}

func (s *Set[T]) AddAll(v []T) {
	if s.internalMap == nil {
		s.internalMap = make(map[T]bool)
	}

	for _, v := range v {
		s.internalMap[v] = true
	}
}

func (s Set[T]) GetSlice() []T {
	if s.internalMap == nil || len(s.internalMap) == 0 {
		return []T{}
	}

	keys := make([]T, 0, len(s.internalMap))
	for k := range s.internalMap {
		keys = append(keys, k)
	}

	return keys
}

func (s *Set[T]) Add(v T) {
	if s.internalMap == nil {
		s.internalMap = make(map[T]bool)
	}

	s.internalMap[v] = true
}

func (s *Set[T]) Remove(v T) {
	if s.internalMap == nil {
		return
	}

	delete(s.internalMap, v)
}

func (s *Set[T]) Clear() {
	clear(s.internalMap)
}

func (s Set[T]) Contains(v T) bool {
	if s.internalMap == nil || len(s.internalMap) == 0 {
		return false
	}

	return s.internalMap[v]
}
