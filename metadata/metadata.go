package metadata

type Metadatable interface {
	Metadata() Metadata
}

type Metadata interface {
	IsExists(key string) bool
	Set(key string, value any)
	Get(key string) any
}
