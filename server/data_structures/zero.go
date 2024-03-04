package datastructures

func Zero[T any]() T {
	return *new(T)
}
