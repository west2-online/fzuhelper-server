package internal

type Option[T any] interface {
	Apply(T)
}

type ApplyOption[T any] struct {
	f func(T)
}

func (x *ApplyOption[T]) Apply(do T) {
	x.f(do)
}

func NewApplyOption[T any](f func(T)) *ApplyOption[T] {
	return &ApplyOption[T]{
		f: f,
	}
}
