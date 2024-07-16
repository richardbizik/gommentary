package handlers

func Get[T any](v *T) T {
	if v != nil {
		return *v
	}
	return *new(T)
}
