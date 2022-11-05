package gowait

func pointer[T any](val T) *T {
	return &val
}
