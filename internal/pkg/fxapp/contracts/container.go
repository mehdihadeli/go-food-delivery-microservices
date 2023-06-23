package contracts

type Container interface {
	ResolveFunc(function interface{})
	ResolveFuncWithParamTag(function interface{}, paramTagName string)
}
