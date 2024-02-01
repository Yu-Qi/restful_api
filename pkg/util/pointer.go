package util

// Ptr returns a pointer to the given value.
func Ptr[T any](t T) *T {
	return &t
}
