package util

type Select[T any] struct {
	left       T
	right      T
	preferLeft bool
}

func SeedSelect[T any](value bool) func(T, T) Select[T] {
	return func(left T, right T) Select[T] {
		return Select[T]{
			left,
			right,
			value,
		}
	}
}

func (s Select[T]) Get() T {
	switch s.left.(type) {
	case string:
		if s.preferLeft && s.right != "" {
			return s.left
		}
		return s.right
	case bool:
		if s.preferLeft {
			return s.left
		}
		return s.right
	}
	panic("Unexpected type")
}
