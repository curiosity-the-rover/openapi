package common

type Builder[T any] interface {
	Build() *T
}

type getSpec[T any] interface {
	getSpec(c Context, visited map[string]struct{}) (*T, error)
}
