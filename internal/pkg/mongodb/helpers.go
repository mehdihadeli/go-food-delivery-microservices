package mongodb

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

	"emperror.dev/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	cursor, err := collection.Find(
		ctx,
		filter,
		&options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var items []T

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/cursor/#retrieve-all-documents
	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, err
	}

	return utils.NewListResult[T](
		items,
		listQuery.GetSize(),
		listQuery.GetPage(),
		count,
	), nil
}
