package hypothesis

func ForT[T any](condition func(T) bool) Hypothesis[T] {
	return &hypothesis[T]{condition: condition}
}

func For(typ interface{}, condition func(interface{}) bool) Hypothesis[interface{}] {
	return &hypothesis[interface{}]{condition: condition}
}
