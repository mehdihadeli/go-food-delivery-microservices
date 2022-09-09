package data

type DataModel[E any] interface {
	ToEntity() E
	FromEntity(entity E)
}
