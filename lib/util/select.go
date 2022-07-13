package util

type Select[T comparable] struct {
	left       T
	right      T
	preferLeft bool
}

func SeedSelect[T comparable](value bool) func(T, T) Select[T] {
	return func(left T, right T) Select[T] {
		return Select[T]{
			left,
			right,
			value,
		}
	}
}

func (s Select[T]) Get() (out T) {
	var zero T
	switch any(&out).(type) {
	case *string:
		if s.preferLeft && s.right != zero {
			out = s.left
			return
		}
		out = s.right
		return
	case *bool:
		if s.preferLeft {
			out = s.left
			return
		}
		out = s.right
		return
	}
	panic("Unexpected type")
}
