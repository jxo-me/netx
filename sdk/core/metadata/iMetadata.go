package metadata

type IMetaDatable interface {
	Metadata() IMetaData
}

type IMetaData interface {
	IsExists(key string) bool
	Set(key string, value any)
	Get(key string) any
}
