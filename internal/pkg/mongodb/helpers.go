package mongodb

import (
	"context"

	"emperror.dev/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

// https://stackoverflow.com/a/23650312/581476

func Paginate[T any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	collection *mongo.Collection,
	filter interface{},
) (*utils.ListResult[T], error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.WrapIf(err, "CountDocuments")
	}

	limit := int64(listQuery.GetLimit())
	skip := int64(listQuery.GetOffset())

	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, errors.WrapIf(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errcheck

	products := make([]T, 0, listQuery.GetSize())

	for cursor.Next(ctx) {
		var prod T
		if err := cursor.Decode(&prod); err != nil {
			return nil, errors.WrapIf(err, "Find")
		}
		products = append(products, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.WrapIf(err, "cursor.Err")
	}

	return utils.NewListResult[T](products, listQuery.GetSize(), listQuery.GetPage(), count), nil
}
